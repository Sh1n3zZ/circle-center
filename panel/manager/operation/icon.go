package manager

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	svc "circle-center/panel/manager/svc"
)

// IconHandler exposes HTTP handlers for icon management
type IconHandler struct {
	service *svc.IconService
}

// NewIconHandler constructs a new IconHandler
func NewIconHandler(db *sql.DB) *IconHandler {
	return &IconHandler{service: svc.NewIconService(db)}
}

// ListIcons handles GET /manager/projects/:id/icons
func (h *IconHandler) ListIcons(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	// Parse query parameters
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 32)
	offset, _ := strconv.ParseInt(c.DefaultQuery("offset", "0"), 10, 32)
	status := c.Query("status")
	packageName := c.Query("package")
	search := c.Query("search")

	params := svc.ListIconsParams{
		ProjectID: projectID,
		Status:    status,
		Package:   packageName,
		Search:    search,
		Limit:     int32(limit),
		Offset:    int32(offset),
	}

	icons, total, err := h.service.ListIcons(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	totalPages := int64(0)
	if limit > 0 {
		totalPages = (total + int64(limit) - 1) / int64(limit)
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"icons":       icons,
			"total":       total,
			"totalPages":  totalPages,
			"currentPage": (offset / limit) + 1,
			"limit":       limit,
			"offset":      offset,
		},
	})
}

// GetIcon handles GET /manager/projects/:id/icons/:iconId
func (h *IconHandler) GetIcon(c *gin.Context) {
	projectIDStr := c.Param("id")
	iconIDStr := c.Param("iconId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	iconID, err := strconv.ParseUint(iconIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid icon id"})
		return
	}

	icon, err := h.service.GetIcon(c.Request.Context(), projectID, iconID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    icon,
	})
}

// CreateIcon handles POST /manager/projects/:id/icons
func (h *IconHandler) CreateIcon(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	var req svc.CreateIconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	icon, err := h.service.CreateIcon(c.Request.Context(), projectID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"data":    icon,
	})
}

// UpdateIcon handles PUT /manager/projects/:id/icons/:iconId
func (h *IconHandler) UpdateIcon(c *gin.Context) {
	projectIDStr := c.Param("id")
	iconIDStr := c.Param("iconId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	iconID, err := strconv.ParseUint(iconIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid icon id"})
		return
	}

	var req svc.UpdateIconRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	icon, err := h.service.UpdateIcon(c.Request.Context(), projectID, iconID, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    icon,
	})
}

// DeleteIcon handles DELETE /manager/projects/:id/icons/:iconId
func (h *IconHandler) DeleteIcon(c *gin.Context) {
	projectIDStr := c.Param("id")
	iconIDStr := c.Param("iconId")

	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	iconID, err := strconv.ParseUint(iconIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid icon id"})
		return
	}

	err = h.service.DeleteIcon(c.Request.Context(), projectID, iconID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "icon deleted successfully",
	})
}

// GetIconStats handles GET /manager/projects/:id/icons/stats
func (h *IconHandler) GetIconStats(c *gin.Context) {
	projectIDStr := c.Param("id")
	projectID, err := strconv.ParseUint(projectIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid project id"})
		return
	}

	stats, err := h.service.GetIconStats(c.Request.Context(), projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    stats,
	})
}
