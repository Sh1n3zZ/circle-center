package account

import (
	"context"
	"fmt"

	dbpkg "circle-center/globals/db"
	accountdb "circle-center/repository/sqlc/account"
)

// VerifyEmailRequest represents the email verification request
type VerifyEmailRequest struct {
	Token string `json:"token" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

// VerifyEmailResponse represents the email verification response
type VerifyEmailResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// VerifyEmailService handles email verification logic
type VerifyEmailService struct {
	queries *accountdb.Queries
}

// NewVerifyEmailService creates a new verification service
func NewVerifyEmailService(queries *accountdb.Queries) *VerifyEmailService {
	return &VerifyEmailService{
		queries: queries,
	}
}

// VerifyEmail handles email verification
func (s *VerifyEmailService) VerifyEmail(ctx context.Context, req *VerifyEmailRequest) (*VerifyEmailResponse, error) {
	tokenKey := fmt.Sprintf("verification_token:%s", req.Token)
	storedEmail, err := dbpkg.Get(ctx, tokenKey)
	if err != nil {
		return &VerifyEmailResponse{
			Success: false,
			Message: "Invalid or expired verification token",
		}, nil
	}

	if storedEmail != req.Email {
		return &VerifyEmailResponse{
			Success: false,
			Message: "Token and email do not match",
		}, nil
	}

	user, err := s.queries.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return &VerifyEmailResponse{
			Success: false,
			Message: "User not found",
		}, nil
	}

	if user.Status {
		return &VerifyEmailResponse{
			Success: true,
			Message: "Account is already verified",
		}, nil
	}

	// update user status to verified (status = true means active)
	err = s.queries.UpdateUserStatus(ctx, accountdb.UpdateUserStatusParams{
		Status: true,
		ID:     user.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to activate user: %w", err)
	}

	// update email verification timestamp
	err = s.queries.VerifyEmail(ctx, user.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to update email verification: %w", err)
	}

	_, err = dbpkg.Del(ctx, tokenKey)
	if err != nil {
		// Log error but don't fail the verification
		// TODO: use Sentry to log errors
		fmt.Printf("Warning: failed to delete verification token from Redis: %v\n", err)
	}

	return &VerifyEmailResponse{
		Success: true,
		Message: "Email verified successfully",
	}, nil
}
