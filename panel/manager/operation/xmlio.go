package manager

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	svc "circle-center/panel/manager/svc"
	mutils "circle-center/panel/manager/utils"
)

// XMLIOHandler exposes endpoints for XML parse (preview) and import (confirm)
type XMLIOHandler struct {
	service *svc.XMLIOService
}

// NewXMLIOHandler constructs handler
func NewXMLIOHandler(db *sql.DB) *XMLIOHandler {
	return &XMLIOHandler{service: svc.NewXMLIOService(db)}
}

// parseForm represents incoming payload for parse-preview
type parseForm struct {
	Appfilter string `form:"appfilter" json:"appfilter"`
	Appmap    string `form:"appmap" json:"appmap"`
	Theme     string `form:"theme" json:"theme"`
}

// importForm represents incoming payload for confirm-import
type importForm struct {
	ProjectID  uint64 `json:"projectId" binding:"required"`
	Components []struct {
		Name          string `json:"name"`
		Package       string `json:"pkg"`
		ComponentInfo string `json:"componentInfo"`
		Drawable      string `json:"drawable"`
	} `json:"components" binding:"required"`
}

// ParsePreview handles POST /manager/icons/parse to return parsed components for review
func (h *XMLIOHandler) ParsePreview(c *gin.Context) {
	var form parseForm
	// accept either form or JSON
	if c.ContentType() == "application/json" {
		_ = c.ShouldBindJSON(&form)
	} else {
		_ = c.ShouldBind(&form)
	}
	comps, err := h.service.ParseXMLInputs(form.Appfilter, form.Appmap, form.Theme)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "components": comps})
}

// ConfirmImport handles POST /manager/icons/import to persist reviewed components
func (h *XMLIOHandler) ConfirmImport(c *gin.Context) {
	var form importForm
	if err := c.ShouldBindJSON(&form); err != nil {
		// Also support query param for projectId for flexibility
		if pid := c.Query("projectId"); pid != "" {
			if v, err2 := strconv.ParseUint(pid, 10, 64); err2 == nil {
				form.ProjectID = v
			}
		}
		if form.ProjectID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "invalid payload"})
			return
		}
	}

	components := make([]mutils.IconRequestComponent, 0, len(form.Components))
	for _, x := range form.Components {
		components = append(components, mutils.IconRequestComponent{
			Name:          x.Name,
			Package:       x.Package,
			ComponentInfo: x.ComponentInfo,
			Drawable:      x.Drawable,
		})
	}

	summary, err := h.service.SaveIcons(c.Request.Context(), form.ProjectID, components)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"status": "error", "message": err.Error(), "summary": summary})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "summary": summary})
}
