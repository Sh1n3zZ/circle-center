package account

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"circle-center/globals/mail"
	op "circle-center/panel/account/operation"
	svc "circle-center/panel/account/svc"
	"circle-center/panel/account/utils"
)

// RegisterRoutes registers all account-related routes
func RegisterRoutes(r *gin.RouterGroup, db *sql.DB, mailService *mail.MailService) {
	userHandler := op.NewUserHandler(db, mailService)
	verificationHandler := op.NewVerificationHandler(db)
	// avatar services and handler
	jwtClient, err := svc.NewJWTClientFromGlobalKeys()
	if err != nil {
		panic("Failed to create JWT client from global keys: " + err.Error())
	}
	authClient := svc.NewAuthClient(jwtClient)
	avatarService, err := svc.NewAvatarService(db, authClient)
	if err != nil {
		panic("Failed to create AvatarService: " + err.Error())
	}
	avatarHandler := op.NewAvatarHandler(avatarService)

	// Account routes
	account := r.Group("/account")
	{
		// Health check
		account.GET("/health", userHandler.HealthCheck)

		// User registration
		account.POST("/register", userHandler.RegisterUser)

		// User login
		account.POST("/login", userHandler.LoginUser)

		// User logout
		account.POST("/logout", userHandler.LogoutUser)

		// Token refresh
		account.POST("/refresh", userHandler.RefreshToken)

		// Resend verification email
		account.POST("/resend-verification", userHandler.ResendVerificationEmail)

		// Email verification (GET for email links if needed, POST for API calls)
		// account.GET("/verify", verificationHandler.VerifyEmailByQuery)
		account.POST("/verify", verificationHandler.VerifyEmail)

		// Get user profile with middleware
		account.GET("/profile",
			utils.ExtractBearerTokenMiddleware(),
			userHandler.GetUserProfileWithMiddleware)

		// Avatar upload (protected)
		account.POST("/avatar",
			utils.ExtractBearerTokenMiddleware(),
			avatarHandler.UploadAvatar)

		// Avatar public get by relative path wildcard
		// Example: GET /v1/account/avatar/avatars/2025/08/uuid_xxx.jpg?size=256&quality=85
		account.GET("/avatar/*relpath", avatarHandler.GetAvatar)
	}
}
