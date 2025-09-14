package manager

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	accountsvc "circle-center/panel/account/svc"
	oputils "circle-center/panel/account/utils"
	svc "circle-center/panel/manager/svc"
)

// ProjectHandler exposes HTTP handlers for project operations
type ProjectHandler struct {
	service *svc.ProjectService
}

// NewProjectHandler constructs handler
func NewProjectHandler(db *sql.DB, authClient *accountsvc.AuthClient) *ProjectHandler {
	service := svc.NewProjectService(db, authClient)
	return &ProjectHandler{service: service}
}

// CreateProject handles POST /manager/projects
// Requires Bearer token
type createProjectRequest struct {
	Name        string  `json:"name"`
	Slug        *string `json:"slug,omitempty"`
	PackageName *string `json:"package_name,omitempty"`
	Visibility  *string `json:"visibility,omitempty"`
	Description *string `json:"description,omitempty"`
}

func (h *ProjectHandler) CreateProject(c *gin.Context) {
	token, ok := oputils.GetTokenFromContext(c)
	if !ok {
		var err error
		token, err = oputils.ExtractBearerToken(c)
		if err != nil {
			oputils.RespondWithAuthError(c, err)
			return
		}
	}

	var req createProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST", "message": err.Error()})
		return
	}

	resp, err := h.service.CreateProject(c.Request.Context(), token, &svc.CreateProjectRequest{
		Name:        req.Name,
		Slug:        req.Slug,
		PackageName: req.PackageName,
		Visibility:  req.Visibility,
		Description: req.Description,
	})
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CREATE_PROJECT_FAILED", "message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Project created successfully",
		"data":    resp,
	})
}

// ListProjects handles GET /manager/projects
func (h *ProjectHandler) ListProjects(c *gin.Context) {
	token, ok := oputils.GetTokenFromContext(c)
	if !ok {
		var err error
		token, err = oputils.ExtractBearerToken(c)
		if err != nil {
			oputils.RespondWithAuthError(c, err)
			return
		}
	}

	limit := int32(50)
	offset := int32(0)
	if v := c.Query("limit"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 32); err == nil {
			limit = int32(parsed)
		}
	}
	if v := c.Query("offset"); v != "" {
		if parsed, err := strconv.ParseInt(v, 10, 32); err == nil {
			offset = int32(parsed)
		}
	}

	list, err := h.service.ListProjects(c.Request.Context(), token, limit, offset)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "LIST_PROJECTS_FAILED", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": list})
}

// GetProject handles GET /manager/projects/:id
func (h *ProjectHandler) GetProject(c *gin.Context) {
	token, ok := oputils.GetTokenFromContext(c)
	if !ok {
		var err error
		token, err = oputils.ExtractBearerToken(c)
		if err != nil {
			oputils.RespondWithAuthError(c, err)
			return
		}
	}

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_PROJECT_ID", "message": "project id must be uint"})
		return
	}

	resp, err := h.service.GetProject(c.Request.Context(), token, projectID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "GET_PROJECT_FAILED", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "ok", "data": resp})
}

// UpdateProject handles PUT /manager/projects/:id
func (h *ProjectHandler) UpdateProject(c *gin.Context) {
	token, ok := oputils.GetTokenFromContext(c)
	if !ok {
		var err error
		token, err = oputils.ExtractBearerToken(c)
		if err != nil {
			oputils.RespondWithAuthError(c, err)
			return
		}
	}

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_PROJECT_ID", "message": "project id must be uint"})
		return
	}

	var req svc.UpdateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST", "message": err.Error()})
		return
	}

	resp, err := h.service.UpdateProject(c.Request.Context(), token, projectID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UPDATE_PROJECT_FAILED", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project updated successfully",
		"data":    resp,
	})
}

// DeleteProject handles DELETE /manager/projects/:id
func (h *ProjectHandler) DeleteProject(c *gin.Context) {
	token, ok := oputils.GetTokenFromContext(c)
	if !ok {
		var err error
		token, err = oputils.ExtractBearerToken(c)
		if err != nil {
			oputils.RespondWithAuthError(c, err)
			return
		}
	}

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_PROJECT_ID", "message": "project id must be uint"})
		return
	}

	if err := h.service.DeleteProject(c.Request.Context(), token, projectID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "DELETE_PROJECT_FAILED", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Project deleted successfully",
	})
}

// AssignProjectRole handles POST /manager/projects/:id/roles
func (h *ProjectHandler) AssignProjectRole(c *gin.Context) {
	token, ok := oputils.GetTokenFromContext(c)
	if !ok {
		var err error
		token, err = oputils.ExtractBearerToken(c)
		if err != nil {
			oputils.RespondWithAuthError(c, err)
			return
		}
	}

	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_PROJECT_ID", "message": "project id must be uint"})
		return
	}

	var req svc.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST", "message": err.Error()})
		return
	}

	if err := h.service.AssignProjectRole(c.Request.Context(), token, projectID, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ASSIGN_ROLE_FAILED", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Role assigned successfully",
	})
}

// ListProjectRoles handles GET /manager/projects/:id/roles
func (h *ProjectHandler) ListProjectRoles(c *gin.Context) {
    token, ok := oputils.GetTokenFromContext(c)
    if !ok {
        var err error
        token, err = oputils.ExtractBearerToken(c)
        if err != nil {
            oputils.RespondWithAuthError(c, err)
            return
        }
    }

    projectIDStr := c.Param("id")
    projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_PROJECT_ID", "message": "project id must be uint"})
        return
    }

    list, err := h.service.GetProjectMembersRoles(c.Request.Context(), token, projectID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "GET_ROLES_FAILED", "message": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"success": true, "data": list})
}

// DeleteProjectCollaborator handles DELETE /manager/projects/:id/roles/:userId
func (h *ProjectHandler) DeleteProjectCollaborator(c *gin.Context) {
    token, ok := oputils.GetTokenFromContext(c)
    if !ok {
        var err error
        token, err = oputils.ExtractBearerToken(c)
        if err != nil {
            oputils.RespondWithAuthError(c, err)
            return
        }
    }

    projectIDStr := c.Param("id")
    userIDStr := c.Param("userId")
    projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_PROJECT_ID", "message": "project id must be uint"})
        return
    }
    userID, err := strconv.ParseUint(userIDStr, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_USER_ID", "message": "user id must be uint"})
        return
    }

    if err := h.service.RemoveProjectCollaborator(c.Request.Context(), token, projectID, userID); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "REMOVE_COLLABORATOR_FAILED", "message": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": true})
}
