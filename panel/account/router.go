package account

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	op "circle-center/panel/account/operation"
)

// RegisterRoutes registers all account-related routes
func RegisterRoutes(r *gin.RouterGroup, db *sql.DB) {
	handler := op.NewUserHandler(db)

	// Account routes
	account := r.Group("/account")
	{
		// Health check
		account.GET("/health", handler.HealthCheck)

		// User registration
		account.POST("/register", handler.RegisterUser)

		// User login
		account.POST("/login", handler.LoginUser)
	}
}
