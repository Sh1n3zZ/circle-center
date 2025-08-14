package manager

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"

	svc "circle-center/panel/manager/svc"
)

// RequestHandler exposes HTTP handlers for request upload
type RequestHandler struct {
	service *svc.RequestService
}

// NewRequestHandler constructs handler
func NewRequestHandler(db *sql.DB) *RequestHandler {
	return &RequestHandler{service: svc.NewRequestService(db)}
}

// uploadRequestForm models the multipart form fields
type uploadRequestForm struct {
	TokenID string `form:"TokenID" binding:"required"`
}

// UploadRequest handles POST /request (multipart/form-data)
// Fields:
// - Header: TokenID (string) via header or form field
// - Part: apps (string) json
// - Part: archive (file)
func (h *RequestHandler) UploadRequest(c *gin.Context) {
	// Accept token either from header or form
	token := c.GetHeader("TokenID")
	if token == "" {
		var form uploadRequestForm
		_ = c.ShouldBind(&form)
		token = form.TokenID
	}
	if token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing token"})
		return
	}

	apps := c.PostForm("apps")
	fileHeader, _ := c.FormFile("archive")
	archivePath := ""
	if fileHeader != nil {
		// For now, we do not persist the file. A production setup should
		// save to a temp/storage path and pass that absolute path.
		// Keep empty archivePath to satisfy DB nullability.
	}

	resp, err := h.service.ProcessUpload(c.Request.Context(), token, apps, archivePath, nil)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": resp.Status, "message": resp.Message})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": resp.Status, "message": resp.Message})
}
