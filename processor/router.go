package editor

import (
	"github.com/gin-gonic/gin"
)

// RegisterRoutes registers processor related endpoints under the given router group.
// Typically the passed router should be a version group (e.g., /v1).
func RegisterRoutes(router gin.IRouter) {
	processorGroup := router.Group("/processor")
	processorGroup.POST("/diffappfilters", diffAppFilters)
	processorGroup.POST("/difficons", diffIcons)
}
