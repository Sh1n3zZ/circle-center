package account

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"circle-center/globals/mail"
	svc "circle-center/panel/account/svc"
	"circle-center/panel/account/utils"
)

type UserHandler struct {
	userService *svc.UserService
}

func NewUserHandler(db *sql.DB, mailService *mail.MailService) *UserHandler {
	jwtClient, err := svc.NewJWTClientFromGlobalKeys()
	if err != nil {
		panic("Failed to create JWT client from global keys: " + err.Error())
	}

	authClient := svc.NewAuthClient(jwtClient)

	return &UserHandler{
		userService: svc.NewUserService(db, mailService, authClient),
	}
}

// RegisterUser handles the user registration HTTP endpoint
// @Summary Register a new user
// @Description Register a new user account with the provided information
// @Tags account
// @Accept json
// @Produce json
// @Param request body svc.RegisterRequest true "User registration information"
// @Success 201 {object} svc.RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/register [post]
func (h *UserHandler) RegisterUser(c *gin.Context) {
	var req svc.RegisterRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	// Additional validation
	if err := h.userService.ValidateRegistrationRequest(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Validation failed",
			Message: err.Error(),
		})
		return
	}

	// Register user
	response, err := h.userService.RegisterUser(c.Request.Context(), &req)
	if err != nil {
		// Handle specific errors
		switch err.Error() {
		case "username already exists":
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Username already exists",
				Message: "The provided username is already taken",
			})
			return
		case "email already exists":
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "Email already exists",
				Message: "The provided email is already registered",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Registration failed",
				Message: err.Error(),
			})
			return
		}
	}

	// Return success response
	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "User registered successfully",
		"data":    response,
	})
}

// LoginUser handles the user login HTTP endpoint
// @Summary Login user
// @Description Authenticate user with email and password
// @Tags account
// @Accept json
// @Produce json
// @Param request body svc.LoginRequest true "User login information"
// @Success 200 {object} svc.LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/login [post]
func (h *UserHandler) LoginUser(c *gin.Context) {
	var req svc.LoginRequest

	// Bind and validate request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	response, err := h.userService.LoginUser(c.Request.Context(), &req)
	if err != nil {
		switch err.Error() {
		case "invalid email or password":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Authentication failed",
				Message: "Invalid email or password",
			})
			return
		case "account is not verified":
			c.JSON(http.StatusForbidden, gin.H{
				"error":   "Account not verified",
				"message": "Your account needs email verification",
				"code":    "ACCOUNT_NOT_VERIFIED",
				"email":   req.Email,
			})
			return
		case "account is temporarily locked":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Account locked",
				Message: "Your account is temporarily locked",
			})
			return
		case "account locked due to too many failed attempts":
			c.JSON(http.StatusUnauthorized, ErrorResponse{
				Error:   "Account locked",
				Message: "Account locked due to too many failed attempts. Please try again later.",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Login failed",
				Message: err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Login successful",
		"data":    response,
	})
}

// ResendVerificationEmail handles the resend verification email HTTP endpoint
// @Summary Resend verification email
// @Description Resend verification email to the specified email address
// @Tags account
// @Accept json
// @Produce json
// @Param request body svc.ResendVerificationEmailRequest true "Resend verification email information"
// @Success 200 {object} svc.ResendVerificationEmailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/resend-verification [post]
func (h *UserHandler) ResendVerificationEmail(c *gin.Context) {
	var req svc.ResendVerificationEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	response, err := h.userService.ResendVerificationEmail(c.Request.Context(), &req)
	if err != nil {
		switch err.Error() {
		case "user not found":
			c.JSON(http.StatusNotFound, ErrorResponse{
				Error:   "User not found",
				Message: "No user found with the provided email address",
			})
			return
		case "user is already verified":
			c.JSON(http.StatusConflict, ErrorResponse{
				Error:   "User already verified",
				Message: "The user account is already verified",
			})
			return
		default:
			c.JSON(http.StatusInternalServerError, ErrorResponse{
				Error:   "Failed to resend verification email",
				Message: err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Verification email sent successfully",
		"data":    response,
	})
}

// LogoutUser handles the user logout HTTP endpoint
// @Summary Logout user
// @Description Logout user by revoking the JWT token
// @Tags account
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/logout [post]
func (h *UserHandler) LogoutUser(c *gin.Context) {
	tokenString, err := utils.ExtractBearerToken(c)
	if err != nil {
		utils.RespondWithAuthError(c, err)
		return
	}

	err = h.userService.LogoutUser(c.Request.Context(), tokenString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Logout failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Logged out successfully",
	})
}

// RefreshToken handles the token refresh HTTP endpoint
// @Summary Refresh JWT token
// @Description Refresh an existing JWT token
// @Tags account
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/refresh [post]
func (h *UserHandler) RefreshToken(c *gin.Context) {
	tokenString, err := utils.ExtractBearerToken(c)
	if err != nil {
		utils.RespondWithAuthError(c, err)
		return
	}

	authResult, err := h.userService.RefreshUserToken(c.Request.Context(), tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Token refresh failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Token refreshed successfully",
		"data":    authResult,
	})
}

// GetUserProfileWithMiddleware shows how to use middleware for token extraction
// This is an alternative implementation that uses middleware
func (h *UserHandler) GetUserProfileWithMiddleware(c *gin.Context) {
	tokenString, exists := utils.GetTokenFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Token not found in context",
			Message: "Middleware token extraction failed",
		})
		return
	}

	userClaims, err := h.userService.ValidateUserToken(c.Request.Context(), tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Token validation failed",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User profile retrieved successfully (via middleware)",
		"data": gin.H{
			"user_id":  userClaims.UserID,
			"username": userClaims.Username,
			"email":    userClaims.Email,
		},
	})
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// HealthCheck provides a simple health check endpoint
// @Summary Health check
// @Description Check if the account service is running
// @Tags account
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /account/health [get]
func (h *UserHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "account",
		"message": "Account service is running",
	})
}
