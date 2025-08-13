package manager

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	accountsvc "circle-center/panel/account/svc"
	"circle-center/panel/account/utils"
	op "circle-center/panel/manager/operation"
)

// RegisterRoutes registers all manager-related routes
func RegisterRoutes(r *gin.RouterGroup, db *sql.DB, authClient *accountsvc.AuthClient) {
	projectHandler := op.NewProjectHandler(db, authClient)

	manager := r.Group("/manager")
	{
		// Protected routes for project management
		manager.GET("/projects",
			utils.ExtractBearerTokenMiddleware(),
			projectHandler.ListProjects,
		)

		manager.GET("/projects/:id",
			utils.ExtractBearerTokenMiddleware(),
			projectHandler.GetProject,
		)

		manager.POST("/projects",
			utils.ExtractBearerTokenMiddleware(),
			projectHandler.CreateProject,
		)

		manager.PUT("/projects/:id",
			utils.ExtractBearerTokenMiddleware(),
			projectHandler.UpdateProject,
		)

		manager.DELETE("/projects/:id",
			utils.ExtractBearerTokenMiddleware(),
			projectHandler.DeleteProject,
		)

		manager.POST("/projects/:id/roles",
			utils.ExtractBearerTokenMiddleware(),
			projectHandler.AssignProjectRole,
		)
	}
}
