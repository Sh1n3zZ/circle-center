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
	requestHandler := op.NewRequestHandler(db)
	tokenHandler := op.NewTokenHandler(db)
	xmlioHandler := op.NewXMLIOHandler(db)
	iconHandler := op.NewIconHandler(db)
	iconioHandler := op.NewIconIOHandler(db, authClient)

	manager := r.Group("/manager")
	{
		manager.GET("/projects",
			utils.ExtractBearerTokenMiddleware(),
			projectHandler.ListProjects,
		)

		manager.GET("/projects/:id/tokens",
			utils.ExtractBearerTokenMiddleware(),
			tokenHandler.List,
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

		manager.GET("/projects/:id/roles",
			utils.ExtractBearerTokenMiddleware(),
			projectHandler.ListProjectRoles,
		)

		manager.DELETE("/projects/:id/roles/:userId",
			utils.ExtractBearerTokenMiddleware(),
			projectHandler.DeleteProjectCollaborator,
		)

		manager.POST("/projects/:id/tokens",
			utils.ExtractBearerTokenMiddleware(),
			tokenHandler.Create,
		)
		manager.DELETE("/projects/:id/tokens/:tokenId",
			utils.ExtractBearerTokenMiddleware(),
			tokenHandler.Delete,
		)

		manager.POST("/icons/parse",
			utils.ExtractBearerTokenMiddleware(),
			xmlioHandler.ParsePreview,
		)
		manager.POST("/icons/import",
			utils.ExtractBearerTokenMiddleware(),
			xmlioHandler.ConfirmImport,
		)

		// Icon management endpoints
		manager.GET("/projects/:id/icons",
			utils.ExtractBearerTokenMiddleware(),
			iconHandler.ListIcons,
		)
		manager.GET("/projects/:id/icons/stats",
			utils.ExtractBearerTokenMiddleware(),
			iconHandler.GetIconStats,
		)
		manager.GET("/projects/:id/icons/:iconId",
			utils.ExtractBearerTokenMiddleware(),
			iconHandler.GetIcon,
		)
		manager.POST("/projects/:id/icons",
			utils.ExtractBearerTokenMiddleware(),
			iconHandler.CreateIcon,
		)
		manager.PUT("/projects/:id/icons/:iconId",
			utils.ExtractBearerTokenMiddleware(),
			iconHandler.UpdateIcon,
		)
		manager.DELETE("/projects/:id/icons/:iconId",
			utils.ExtractBearerTokenMiddleware(),
			iconHandler.DeleteIcon,
		)

		manager.GET("/icons/*relpath",
			utils.ExtractBearerTokenMiddleware(),
			iconioHandler.GetIcon,
		)
		manager.POST("/icons/:projectId/upload",
			utils.ExtractBearerTokenMiddleware(),
			iconioHandler.UploadIcon,
		)
	}

	request := r.Group("")
	{
		request.POST("/request", requestHandler.UploadRequest)
	}
}
