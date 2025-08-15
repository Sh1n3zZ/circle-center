package manager

import (
	"database/sql"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"

	accountsvc "circle-center/panel/account/svc"
	oputils "circle-center/panel/account/utils"
	svc "circle-center/panel/manager/svc"
)

// IconIOHandler exposes HTTP handlers for icon upload and retrieval
type IconIOHandler struct {
	service *svc.IconIOService
}

// NewIconIOHandler constructs a new IconIOHandler
func NewIconIOHandler(db *sql.DB, authClient *accountsvc.AuthClient) *IconIOHandler {
	service, err := svc.NewIconIOService(db, authClient)
	if err != nil {
		panic("Failed to create IconIOService: " + err.Error())
	}
	return &IconIOHandler{service: service}
}

// UploadIcon handles POST /manager/icons/:projectId/upload
// Form fields: component_info (string), file (multipart file)
func (h *IconIOHandler) UploadIcon(c *gin.Context) {
	// Auth: extract token from context or header
	token, ok := oputils.GetTokenFromContext(c)
	if !ok {
		var err error
		token, err = oputils.ExtractBearerToken(c)
		if err != nil {
			oputils.RespondWithAuthError(c, err)
			return
		}
	}

	projectIDStr := c.Param("projectId")
	projectID, err := strconv.ParseUint(strings.TrimSpace(projectIDStr), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_PROJECT_ID", "message": "project id must be uint"})
		return
	}

	componentInfo := strings.TrimSpace(c.PostForm("component_info"))
	if componentInfo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_REQUEST", "message": "component_info is required"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "FILE_REQUIRED", "message": "file is required"})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "FILE_READ_ERROR", "message": err.Error()})
		return
	}

	path, err := h.service.ValidateAndSaveIcon(c.Request.Context(), token, projectID, componentInfo, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UPLOAD_FAILED", "message": err.Error()})
		return
	}

	// Infer content-type for response convenience
	ct := "application/octet-stream"
	if kind, err := filetype.Match(data); err == nil && kind != filetype.Unknown {
		ct = kind.MIME.Value
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Icon uploaded successfully",
		"data": gin.H{
			"path":         path,
			"content_type": ct,
		},
	})
}

// GetIcon handles GET /manager/icon/*relpath
// Requires Bearer token; returns the raw bytes with appropriate content-type
func (h *IconIOHandler) GetIcon(c *gin.Context) {
	// Auth
	token, ok := oputils.GetTokenFromContext(c)
	if !ok {
		var err error
		token, err = oputils.ExtractBearerToken(c)
		if err != nil {
			oputils.RespondWithAuthError(c, err)
			return
		}
	}

	rel := c.Param("relpath")
	if len(rel) > 0 && rel[0] == '/' {
		rel = rel[1:]
	}

	abs, err := h.service.GetIconAbsolutePathSecure(c.Request.Context(), token, rel)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NOT_FOUND", "message": err.Error()})
		return
	}

	if _, err := os.Stat(abs); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "ICON_FILE_NOT_FOUND", 
			"message": "Icon file not uploaded yet",
			"code": "ICON_FILE_NOT_FOUND",
		})
		return
	}

	bytes, err := os.ReadFile(abs)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "READ_FAILED", "message": err.Error()})
		return
	}

	ct := "application/octet-stream"
	if kind, err := filetype.Match(bytes); err == nil && kind != filetype.Unknown {
		ct = kind.MIME.Value
	}

	c.Header("Content-Type", ct)
	c.Header("Cache-Control", "private, max-age=31536000")
	c.Status(http.StatusOK)
	_, _ = c.Writer.Write(bytes)
}
