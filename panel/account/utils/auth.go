package utils

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// AuthError represents authentication-related errors
type AuthError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e AuthError) Error() string {
	return e.Message
}

// Common auth error codes
var (
	ErrMissingAuthHeader = AuthError{
		Code:    "MISSING_AUTH_HEADER",
		Message: "Authorization header is required",
	}
	ErrInvalidAuthHeader = AuthError{
		Code:    "INVALID_AUTH_HEADER",
		Message: "Authorization header must start with 'Bearer '",
	}
	ErrInvalidTokenFormat = AuthError{
		Code:    "INVALID_TOKEN_FORMAT",
		Message: "Invalid token format",
	}
)

// ExtractBearerToken extracts the JWT token from the Authorization header
// Expected format: "Bearer <token>"
func ExtractBearerToken(c *gin.Context) (string, error) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		return "", ErrMissingAuthHeader
	}

	return ParseBearerToken(authHeader)
}

// ParseBearerToken parses a bearer token from an authorization header string
func ParseBearerToken(authHeader string) (string, error) {
	const bearerPrefix = "Bearer "
	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", ErrInvalidAuthHeader
	}

	token := strings.TrimSpace(authHeader[len(bearerPrefix):])
	if token == "" {
		return "", ErrInvalidTokenFormat
	}

	return token, nil
}

// RespondWithAuthError sends a standardized authentication error response
func RespondWithAuthError(c *gin.Context, err error) {
	var authErr AuthError
	var statusCode int

	if ae, ok := err.(AuthError); ok {
		authErr = ae
		switch ae.Code {
		case "MISSING_AUTH_HEADER", "INVALID_AUTH_HEADER", "INVALID_TOKEN_FORMAT":
			statusCode = http.StatusBadRequest
		default:
			statusCode = http.StatusUnauthorized
		}
	} else {
		authErr = AuthError{
			Code:    "AUTH_ERROR",
			Message: err.Error(),
		}
		statusCode = http.StatusUnauthorized
	}

	c.JSON(statusCode, gin.H{
		"error":   authErr.Code,
		"message": authErr.Message,
	})
}

// ExtractBearerTokenMiddleware creates a middleware that extracts and validates bearer token
// The token will be stored in the context with the key "token"
func ExtractBearerTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := ExtractBearerToken(c)
		if err != nil {
			RespondWithAuthError(c, err)
			c.Abort()
			return
		}

		c.Set("token", token)
		c.Next()
	}
}

// GetTokenFromContext retrieves the token from gin context
func GetTokenFromContext(c *gin.Context) (string, bool) {
	token, exists := c.Get("token")
	if !exists {
		return "", false
	}

	tokenStr, ok := token.(string)
	return tokenStr, ok
}

// ValidateTokenFormat performs basic token format validation
func ValidateTokenFormat(token string) error {
	if token == "" {
		return ErrInvalidTokenFormat
	}

	// basic JWT format validation (3 parts separated by dots)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return AuthError{
			Code:    "INVALID_JWT_FORMAT",
			Message: "JWT token must have 3 parts separated by dots",
		}
	}

	for i, part := range parts {
		if strings.TrimSpace(part) == "" {
			return AuthError{
				Code:    "INVALID_JWT_FORMAT",
				Message: fmt.Sprintf("JWT part %d is empty", i+1),
			}
		}
	}

	return nil
}
