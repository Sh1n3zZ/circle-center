package account

import (
	"context"
	"database/sql"
	"fmt"

	accountdb "circle-center/repository/sqlc/account"
)

// ProfileService handles user profile operations
type ProfileService struct {
	queries *accountdb.Queries
	auth    *AuthClient
}

// NewProfileService creates a new profile service
func NewProfileService(db *sql.DB, authClient *AuthClient) *ProfileService {
	return &ProfileService{
		queries: accountdb.New(db),
		auth:    authClient,
	}
}

// GetUserProfileRequest represents the request to get user profile
type GetUserProfileRequest struct {
	UserID uint64 `json:"user_id"`
}

// GetUserProfileResponse represents the user profile response
type GetUserProfileResponse struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Locale      string `json:"locale"`
	Timezone    string `json:"timezone"`
	AvatarUrl   string `json:"avatar_url,omitempty"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// UpdateUserProfileRequest represents the request to update user profile
type UpdateUserProfileRequest struct {
	DisplayName string `json:"display_name,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Locale      string `json:"locale,omitempty"`
	Timezone    string `json:"timezone,omitempty"`
}

// UpdateUserProfileResponse represents the response after updating user profile
type UpdateUserProfileResponse struct {
	ID          uint64 `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	DisplayName string `json:"display_name,omitempty"`
	Phone       string `json:"phone,omitempty"`
	Locale      string `json:"locale"`
	Timezone    string `json:"timezone"`
	AvatarUrl   string `json:"avatar_url,omitempty"`
	UpdatedAt   string `json:"updated_at"`
}

// GetUserProfile retrieves user profile information
func (s *ProfileService) GetUserProfile(ctx context.Context, userID uint64) (*GetUserProfileResponse, error) {
	user, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	response := &GetUserProfileResponse{
		ID:          user.ID,
		Username:    user.Username,
		Email:       user.Email,
		DisplayName: user.DisplayName.String,
		Phone:       user.Phone.String,
		Locale:      user.Locale,
		Timezone:    user.Timezone,
		AvatarUrl:   user.AvatarUrl.String,
		CreatedAt:   user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}

// UpdateUserProfile updates user profile information
func (s *ProfileService) UpdateUserProfile(ctx context.Context, userID uint64, req *UpdateUserProfileRequest) (*UpdateUserProfileResponse, error) {
	// Get current user to preserve existing data
	user, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	// Prepare update parameters, preserving existing values if not provided
	displayName := user.DisplayName
	if req.DisplayName != "" {
		displayName = sql.NullString{String: req.DisplayName, Valid: true}
	}

	phone := user.Phone
	if req.Phone != "" {
		phone = sql.NullString{String: req.Phone, Valid: true}
	}

	locale := user.Locale
	if req.Locale != "" {
		locale = req.Locale
	}

	timezone := user.Timezone
	if req.Timezone != "" {
		timezone = req.Timezone
	}

	// Update user profile
	err = s.queries.UpdateUserProfile(ctx, accountdb.UpdateUserProfileParams{
		DisplayName: displayName,
		AvatarUrl:   user.AvatarUrl, // Preserve existing avatar
		Phone:       phone,
		Locale:      locale,
		Timezone:    timezone,
		ID:          userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	// Get updated user
	updatedUser, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated user: %w", err)
	}

	response := &UpdateUserProfileResponse{
		ID:          updatedUser.ID,
		Username:    updatedUser.Username,
		Email:       updatedUser.Email,
		DisplayName: updatedUser.DisplayName.String,
		Phone:       updatedUser.Phone.String,
		Locale:      updatedUser.Locale,
		Timezone:    updatedUser.Timezone,
		AvatarUrl:   updatedUser.AvatarUrl.String,
		UpdatedAt:   updatedUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	return response, nil
}

// ValidateUserToken validates token and returns user claims, delegating to AuthClient
func (s *ProfileService) ValidateUserToken(ctx context.Context, tokenString string) (*UserClaims, error) {
	if s.auth == nil {
		return nil, fmt.Errorf("auth client not initialized")
	}
	return s.auth.ValidateToken(ctx, tokenString)
}
