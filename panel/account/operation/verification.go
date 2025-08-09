package account

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	svc "circle-center/panel/account/svc"
	accountdb "circle-center/repository/sqlc/account"
)

type VerificationHandler struct {
	verifyService *svc.VerifyEmailService
}

// NewVerificationHandler creates a new verification handler
func NewVerificationHandler(db *sql.DB) *VerificationHandler {
	queries := accountdb.New(db)
	return &VerificationHandler{
		verifyService: svc.NewVerifyEmailService(queries),
	}
}

// VerifyEmail handles the email verification HTTP endpoint
// @Summary Verify user email
// @Description Verify user email with token and email address
// @Tags account
// @Accept json
// @Produce json
// @Param request body svc.VerifyEmailRequest true "Email verification information"
// @Success 200 {object} svc.VerifyEmailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/verify [post]
func (h *VerificationHandler) VerifyEmail(c *gin.Context) {
	var req svc.VerifyEmailRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request data",
			Message: err.Error(),
		})
		return
	}

	response, err := h.verifyService.VerifyEmail(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Verification failed",
			Message: err.Error(),
		})
		return
	}

	if response.Success {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": response.Message,
			"data":    response,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": response.Message,
			"data":    response,
		})
	}
}

// VerifyEmailByQuery handles email verification via query parameters (for email links)
// @Summary Verify user email via query parameters
// @Description Verify user email with token and email from query parameters
// @Tags account
// @Produce json
// @Param token query string true "Verification token"
// @Param email query string true "Email address"
// @Success 200 {object} svc.VerifyEmailResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /account/verify [get]
// func (h *VerificationHandler) VerifyEmailByQuery(c *gin.Context) {
// 	token := c.Query("token")
// 	email := c.Query("email")

// 	if token == "" || email == "" {
// 		c.JSON(http.StatusBadRequest, ErrorResponse{
// 			Error:   "Missing parameters",
// 			Message: "Token and email are required",
// 		})
// 		return
// 	}

// 	req := svc.VerifyEmailRequest{
// 		Token: token,
// 		Email: email,
// 	}

// 	response, err := h.verifyService.VerifyEmail(c.Request.Context(), &req)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, ErrorResponse{
// 			Error:   "Verification failed",
// 			Message: err.Error(),
// 		})
// 		return
// 	}

// 	if response.Success {
// 		c.JSON(http.StatusOK, gin.H{
// 			"success": true,
// 			"message": response.Message,
// 			"data":    response,
// 		})
// 	} else {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"success": false,
// 			"message": response.Message,
// 			"data":    response,
// 		})
// 	}
// }
