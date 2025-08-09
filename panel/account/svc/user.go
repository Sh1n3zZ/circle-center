package account

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"

	configure "circle-center/globals/configure"
	dbpkg "circle-center/globals/db"
	"circle-center/globals/mail"
	accountdb "circle-center/repository/sqlc/account"
)

type UserService struct {
	queries     *accountdb.Queries
	mailService *mail.MailService
}

func NewUserService(db *sql.DB, mailService *mail.MailService) *UserService {
	return &UserService{
		queries:     accountdb.New(db),
		mailService: mailService,
	}
}

// RegisterRequest represents the user registration request
type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=64,alphanum"`
	Email       string `json:"email" binding:"required,email"`
	Password    string `json:"password" binding:"required,min=8"`
	DisplayName string `json:"display_name" binding:"max=100"`
	Phone       string `json:"phone" binding:"omitempty,len=11"`
	Locale      string `json:"locale" binding:"omitempty,len=5"`
	Timezone    string `json:"timezone" binding:"omitempty,max=50"`
}

// RegisterResponse represents the user registration response
type RegisterResponse struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Locale      string `json:"locale"`
	Timezone    string `json:"timezone"`
	CreatedAt   string `json:"created_at"`
	EmailSent   bool   `json:"email_sent"`
	EmailError  string `json:"email_error,omitempty"`
}

// LoginRequest represents the user login request
type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents the user login response
type LoginResponse struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Locale      string `json:"locale"`
	Timezone    string `json:"timezone"`
	Token       string `json:"token,omitempty"`
}

// RegisterUser handles user registration
func (s *UserService) RegisterUser(ctx context.Context, req *RegisterRequest) (*RegisterResponse, error) {
	// Validate username uniqueness
	_, err := s.queries.GetUserByUsername(ctx, req.Username)
	if err == nil {
		return nil, fmt.Errorf("username already exists")
	}

	// Validate email uniqueness
	_, err = s.queries.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, fmt.Errorf("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Set default values
	locale := req.Locale
	if locale == "" {
		locale = "en_US"
	}

	timezone := req.Timezone
	if timezone == "" {
		timezone = "UTC"
	}

	// Create user parameters
	params := accountdb.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		DisplayName:  sql.NullString{String: req.DisplayName, Valid: req.DisplayName != ""},
		Phone:        sql.NullString{String: req.Phone, Valid: req.Phone != ""},
		Locale:       locale,
		Timezone:     timezone,
	}

	// Create user in database
	result, err := s.queries.CreateUser(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Get the created user ID
	userID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get user ID: %w", err)
	}

	// Get the created user
	user, err := s.queries.GetUserByID(ctx, uint64(userID))
	if err != nil {
		return nil, fmt.Errorf("failed to get created user: %w", err)
	}

	// Generate verification token
	verificationToken, err := s.GenerateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Store verification token in Redis with 24-hour expiration
	tokenKey := fmt.Sprintf("verification_token:%s", verificationToken)
	err = dbpkg.Set(ctx, tokenKey, user.Email, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to store verification token: %w", err)
	}

	// Build response
	response := &RegisterResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName.String,
		Phone:       user.Phone.String,
		Locale:      user.Locale,
		Timezone:    user.Timezone,
		CreatedAt:   user.CreatedAt.Format(time.RFC3339),
		EmailSent:   false,
		EmailError:  "",
	}

	// Send verification email
	if s.mailService != nil {
		config := configure.GetConfig()
		verificationURL := fmt.Sprintf("%s/login?verification=true&token=%s&email=%s", config.Frontend.BaseURL, verificationToken, user.Email)
		err = s.mailService.SendVerificationEmail(user.Email, verificationToken, verificationURL)
		if err != nil {
			response.EmailError = err.Error()
		} else {
			response.EmailSent = true
		}
	}

	return response, nil
}

// LoginUser handles user login
func (s *UserService) LoginUser(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// Get user by email
	user, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	// Check if account is active
	if !user.Status {
		return nil, fmt.Errorf("account is not verified")
	}

	// Check if account is locked
	if user.LockedUntil.Valid && user.LockedUntil.Time.After(time.Now()) {
		return nil, fmt.Errorf("account is temporarily locked")
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		// Increment failed attempts
		s.queries.IncrementFailedAttempts(ctx, user.ID)

		// Lock account if too many failed attempts (5 or more)
		if user.FailedAttempts >= 4 {
			lockUntil := time.Now().Add(30 * time.Minute)
			s.queries.LockUser(ctx, accountdb.LockUserParams{
				LockedUntil: sql.NullTime{Time: lockUntil, Valid: true},
				ID:          user.ID,
			})
			return nil, fmt.Errorf("account locked due to too many failed attempts")
		}

		return nil, fmt.Errorf("invalid email or password")
	}

	// Update last login and reset failed attempts
	err = s.queries.UpdateLastLogin(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update login time: %w", err)
	}

	// Generate a simple token (in production, use JWT)
	token, err := s.GenerateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate token: %w", err)
	}

	// Build response
	response := &LoginResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName.String,
		Phone:       user.Phone.String,
		Locale:      user.Locale,
		Timezone:    user.Timezone,
		Token:       token,
	}

	return response, nil
}

// TestDatabaseConnection tests if the database connection is working
func (s *UserService) TestDatabaseConnection(ctx context.Context) error {
	// Try a simple query to test the connection
	_, err := s.queries.CountUsers(ctx)
	return err
}

// ValidateRegistrationRequest validates the registration request
func (s *UserService) ValidateRegistrationRequest(req *RegisterRequest) error {
	// Username validation
	if strings.Contains(req.Username, " ") {
		return fmt.Errorf("username cannot contain spaces")
	}

	// Email validation
	if !strings.Contains(req.Email, "@") {
		return fmt.Errorf("invalid email format")
	}

	// Phone validation (if provided)
	if req.Phone != "" && !strings.HasPrefix(req.Phone, "+") {
		return fmt.Errorf("phone number must be in E.164 format (+country code)")
	}

	return nil
}

// GenerateSecureToken generates a secure random token
func (s *UserService) GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ResendVerificationEmailRequest represents the resend verification email request
type ResendVerificationEmailRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResendVerificationEmailResponse represents the resend verification email response
type ResendVerificationEmailResponse struct {
	EmailSent  bool   `json:"email_sent"`
	EmailError string `json:"email_error,omitempty"`
}

// ResendVerificationEmail handles resending verification email
func (s *UserService) ResendVerificationEmail(ctx context.Context, req *ResendVerificationEmailRequest) (*ResendVerificationEmailResponse, error) {
	// Check if user exists
	user, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	// Check if user is already verified
	if user.Status {
		return nil, fmt.Errorf("user is already verified")
	}

	// Generate new verification token
	verificationToken, err := s.GenerateSecureToken(32)
	if err != nil {
		return nil, fmt.Errorf("failed to generate verification token: %w", err)
	}

	// Store verification token in Redis with 24-hour expiration
	tokenKey := fmt.Sprintf("verification_token:%s", verificationToken)
	err = dbpkg.Set(ctx, tokenKey, user.Email, 24*time.Hour)
	if err != nil {
		return nil, fmt.Errorf("failed to store verification token: %w", err)
	}

	// Build response
	response := &ResendVerificationEmailResponse{
		EmailSent:  false,
		EmailError: "",
	}

	// Send verification email
	if s.mailService != nil {
		config := configure.GetConfig()
		verificationURL := fmt.Sprintf("%s/login?verification=true&token=%s&email=%s", config.Frontend.BaseURL, verificationToken, user.Email)
		err = s.mailService.SendVerificationEmail(user.Email, verificationToken, verificationURL)
		if err != nil {
			response.EmailError = err.Error()
		} else {
			response.EmailSent = true
		}
	}

	return response, nil
}
