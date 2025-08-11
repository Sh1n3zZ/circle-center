package account

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	"circle-center/globals/mail"
	svc "circle-center/panel/account/svc"
	"circle-center/panel/account/utils"
)

// ProfileHandler handles profile-related HTTP requests
type ProfileHandler struct {
	profileService *svc.ProfileService
}

// NewProfileHandler creates a new profile handler
func NewProfileHandler(db *sql.DB, mailService *mail.MailService) *ProfileHandler {
	jwtClient, err := svc.NewJWTClientFromGlobalKeys()
	if err != nil {
		panic("Failed to create JWT client from global keys: " + err.Error())
	}

	authClient := svc.NewAuthClient(jwtClient)
	profileService := svc.NewProfileService(db, authClient)

	return &ProfileHandler{
		profileService: profileService,
	}
}

// GetUserProfile handles the get user profile HTTP endpoint
// @Summary Get user profile
// @Description Get current user's profile information
// @Tags profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} svc.GetUserProfileResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/profile [get]
func (h *ProfileHandler) GetUserProfile(c *gin.Context) {
	// Get user ID from token
	tokenString, exists := utils.GetTokenFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Token not found in context",
			Message: "Middleware token extraction failed",
		})
		return
	}

	userClaims, err := h.profileService.ValidateUserToken(c.Request.Context(), tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Token validation failed",
			Message: err.Error(),
		})
		return
	}

	// Get user profile
	response, err := h.profileService.GetUserProfile(c.Request.Context(), userClaims.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to get user profile",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User profile retrieved successfully",
		"data":    response,
	})
}

// UpdateUserProfile handles the update user profile HTTP endpoint
// @Summary Update user profile
// @Description Update current user's profile information
// @Tags profile
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param request body svc.UpdateUserProfileRequest true "Profile update information"
// @Success 200 {object} svc.UpdateUserProfileResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/profile [put]
func (h *ProfileHandler) UpdateUserProfile(c *gin.Context) {
	// Get user ID from token
	tokenString, exists := utils.GetTokenFromContext(c)
	if !exists {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Token not found in context",
			Message: "Middleware token extraction failed",
		})
		return
	}

	userClaims, err := h.profileService.ValidateUserToken(c.Request.Context(), tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{
			Error:   "Token validation failed",
			Message: err.Error(),
		})
		return
	}

	// Bind and validate request
	var req svc.UpdateUserProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	// Update user profile
	response, err := h.profileService.UpdateUserProfile(c.Request.Context(), userClaims.UserID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update user profile",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "User profile updated successfully",
		"data":    response,
	})
}
