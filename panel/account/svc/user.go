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

	accountdb "circle-center/repository/sqlc/account"
)

type UserService struct {
	queries *accountdb.Queries
}

func NewUserService(db *sql.DB) *UserService {
	return &UserService{
		queries: accountdb.New(db),
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
	}

	return response, nil
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
