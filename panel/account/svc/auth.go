package account

import (
	"context"
	"fmt"
	"time"

	dbpkg "circle-center/globals/db"
)

// AuthClient handles authentication operations with JWT and Redis
type AuthClient struct {
	jwtClient *JWTClient
}

// NewAuthClient creates a new authentication client
func NewAuthClient(jwtClient *JWTClient) *AuthClient {
	return &AuthClient{
		jwtClient: jwtClient,
	}
}

// AuthResult represents the result of authentication operations
type AuthResult struct {
	Token     string `json:"token"`
	ExpiresAt int64  `json:"expires_at"`
}

// UserClaims represents user claims in JWT token
type UserClaims struct {
	UserID   uint64 `json:"user_id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

// GenerateToken generates a JWT token for a user and stores it in Redis
func (a *AuthClient) GenerateToken(ctx context.Context, userID uint64, username, email string) (*AuthResult, error) {
	// Generate JWT token
	tokenString, err := a.jwtClient.GenerateToken(userID, username, email)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}

	// Parse token to get expiration time
	token, err := a.jwtClient.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse generated token: %w", err)
	}

	expirationTime, _ := token.Expiration()

	// Store token in Redis for tracking/revocation
	tokenKey := fmt.Sprintf("auth_token:%d:%s", userID, tokenString)
	err = dbpkg.Set(ctx, tokenKey, "active", time.Until(expirationTime))
	if err != nil {
		return nil, fmt.Errorf("failed to store token in Redis: %w", err)
	}

	// Store user session mapping
	sessionKey := fmt.Sprintf("user_session:%d", userID)
	err = dbpkg.Set(ctx, sessionKey, tokenString, time.Until(expirationTime))
	if err != nil {
		return nil, fmt.Errorf("failed to store user session: %w", err)
	}

	return &AuthResult{
		Token:     tokenString,
		ExpiresAt: expirationTime.Unix(),
	}, nil
}

// ValidateToken validates a JWT token and checks if it's not revoked
func (a *AuthClient) ValidateToken(ctx context.Context, tokenString string) (*UserClaims, error) {
	// First validate the JWT token signature and expiration
	token, err := a.jwtClient.ValidateToken(tokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token: %w", err)
	}

	// Extract user information
	userID, username, email, err := a.jwtClient.ExtractUserInfo(token)
	if err != nil {
		return nil, fmt.Errorf("failed to extract user info: %w", err)
	}

	// Check if token is revoked in Redis
	tokenKey := fmt.Sprintf("auth_token:%d:%s", userID, tokenString)
	status, err := dbpkg.Get(ctx, tokenKey)
	if err != nil {
		// If token not found in Redis, it might be expired or revoked
		return nil, fmt.Errorf("token not found or expired")
	}

	if status != "active" {
		return nil, fmt.Errorf("token has been revoked")
	}

	return &UserClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
	}, nil
}

// RevokeToken revokes a specific token
func (a *AuthClient) RevokeToken(ctx context.Context, tokenString string) error {
	// Validate token to get user ID
	token, err := a.jwtClient.ValidateToken(tokenString)
	if err != nil {
		return fmt.Errorf("invalid token: %w", err)
	}

	userID, _, _, err := a.jwtClient.ExtractUserInfo(token)
	if err != nil {
		return fmt.Errorf("failed to extract user info: %w", err)
	}

	// Mark token as revoked in Redis
	tokenKey := fmt.Sprintf("auth_token:%d:%s", userID, tokenString)
	err = dbpkg.Set(ctx, tokenKey, "revoked", time.Hour*24) // Keep revoked status for 24 hours
	if err != nil {
		return fmt.Errorf("failed to revoke token: %w", err)
	}

	return nil
}

// RevokeAllUserTokens revokes all tokens for a specific user
func (a *AuthClient) RevokeAllUserTokens(ctx context.Context, userID uint64) error {
	// Get current user session token
	sessionKey := fmt.Sprintf("user_session:%d", userID)
	currentToken, err := dbpkg.Get(ctx, sessionKey)
	if err == nil && currentToken != "" {
		// Revoke the current session token
		tokenKey := fmt.Sprintf("auth_token:%d:%s", userID, currentToken)
		err = dbpkg.Set(ctx, tokenKey, "revoked", time.Hour*24)
		if err != nil {
			return fmt.Errorf("failed to revoke user session token: %w", err)
		}
	}

	// Delete user session
	_, err = dbpkg.Del(ctx, sessionKey)
	if err != nil {
		return fmt.Errorf("failed to delete user session: %w", err)
	}

	return nil
}

// RefreshToken generates a new token for a user and revokes the old one
func (a *AuthClient) RefreshToken(ctx context.Context, oldTokenString string) (*AuthResult, error) {
	// Validate old token
	userClaims, err := a.ValidateToken(ctx, oldTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid token for refresh: %w", err)
	}

	// Revoke old token
	err = a.RevokeToken(ctx, oldTokenString)
	if err != nil {
		return nil, fmt.Errorf("failed to revoke old token: %w", err)
	}

	// Generate new token
	return a.GenerateToken(ctx, userClaims.UserID, userClaims.Username, userClaims.Email)
}

// LogoutUser logs out a user by revoking their current session
func (a *AuthClient) LogoutUser(ctx context.Context, tokenString string) error {
	return a.RevokeToken(ctx, tokenString)
}

// LogoutAllSessions logs out a user from all sessions
func (a *AuthClient) LogoutAllSessions(ctx context.Context, userID uint64) error {
	return a.RevokeAllUserTokens(ctx, userID)
}

// IsTokenActive checks if a token is active without full validation
func (a *AuthClient) IsTokenActive(ctx context.Context, userID uint64, tokenString string) (bool, error) {
	tokenKey := fmt.Sprintf("auth_token:%d:%s", userID, tokenString)
	status, err := dbpkg.Get(ctx, tokenKey)
	if err != nil {
		return false, nil // Token not found or expired
	}
	return status == "active", nil
}

// GetUserCurrentSession gets the current session token for a user
func (a *AuthClient) GetUserCurrentSession(ctx context.Context, userID uint64) (string, error) {
	sessionKey := fmt.Sprintf("user_session:%d", userID)
	return dbpkg.Get(ctx, sessionKey)
}

// CleanupExpiredTokens removes expired token records from Redis (utility function)
func (a *AuthClient) CleanupExpiredTokens(ctx context.Context) error {
	// This is a utility function that could be called periodically
	// Redis TTL will automatically expire keys, but we can implement
	// additional cleanup logic here if needed
	return nil
}
