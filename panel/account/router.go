package account

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"circle-center/globals/mail"
	op "circle-center/panel/account/operation"
)

// RegisterRoutes registers all account-related routes
func RegisterRoutes(r *gin.RouterGroup, db *sql.DB, mailService *mail.MailService) {
	userHandler := op.NewUserHandler(db, mailService)
	verificationHandler := op.NewVerificationHandler(db)

	// Account routes
	account := r.Group("/account")
	{
		// Health check
		account.GET("/health", userHandler.HealthCheck)

		// User registration
		account.POST("/register", userHandler.RegisterUser)

		// User login
		account.POST("/login", userHandler.LoginUser)

		// Resend verification email
		account.POST("/resend-verification", userHandler.ResendVerificationEmail)

		// Email verification (GET for email links if needed, POST for API calls)
		// account.GET("/verify", verificationHandler.VerifyEmailByQuery)
		account.POST("/verify", verificationHandler.VerifyEmail)
	}
}
