package globals

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// SetupRouter creates a baseline gin.Engine with default middleware (logger & recovery).
// Other packages should register their own routes against the returned engine.
func SetupRouter() *gin.Engine {
	// gin.Default() installs logger and recovery middleware out-of-the-box.
	r := gin.Default()

	config := cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}
	r.Use(cors.New(config))

	// Define API version group /v1 for all routes.
	r.Group("/v1")

	// Other global middlewares could be added here (CORS, auth, etc.)
	return r
}
