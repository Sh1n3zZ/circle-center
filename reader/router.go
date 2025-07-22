package reader

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers reader related endpoints under the given router group.
// Typically the passed router should be a version group (e.g., /v1).
func RegisterRoutes(router gin.IRouter) {
	readerGroup := router.Group("/reader")
	readerGroup.POST("/readfile", reader)
}
