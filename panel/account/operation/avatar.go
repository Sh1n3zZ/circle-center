package account

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/h2non/filetype"

	svc "circle-center/panel/account/svc"
	imgutils "circle-center/panel/account/utils"
)

// AvatarHandler contains HTTP handlers for avatar operations.
type AvatarHandler struct {
	avatarService *svc.AvatarService
}

// NewAvatarHandler constructs a new AvatarHandler.
func NewAvatarHandler(avatarService *svc.AvatarService) *AvatarHandler {
	return &AvatarHandler{avatarService: avatarService}
}

// UploadAvatar handles authenticated avatar upload.
// Content-Type: multipart/form-data, field name: "file".
func (h *AvatarHandler) UploadAvatar(c *gin.Context) {
	token, ok := c.Get("token")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "TOKEN_MISSING", "message": "token not found"})
		return
	}
	tokenString, _ := token.(string)
	claims, err := h.avatarService.ValidateUserToken(c.Request.Context(), tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "TOKEN_INVALID", "message": err.Error()})
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

	result, err := h.avatarService.ValidateAndSaveAvatar(c.Request.Context(), claims.UserID, data)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "UPLOAD_FAILED", "message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Avatar uploaded successfully",
		"data": gin.H{
			"path": result.Path,
			"url":  result.URL,
		},
	})
}

// GetAvatar serves avatar by stored relative path. Public access.
// Route param: :relpath (can contain directories). Query: size, quality
func (h *AvatarHandler) GetAvatar(c *gin.Context) {
	relative := c.Param("relpath")
	if len(relative) > 0 && relative[0] == '/' {
		relative = relative[1:]
	}

	sizeParam := c.Query("size")
	qualityParam := c.Query("quality")

	targetSize, targetQuality, err := h.avatarService.ParseTargetSizeAndQuality(sizeParam, qualityParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "INVALID_PARAMS", "message": err.Error()})
		return
	}

	absPath, err := h.avatarService.GetAvatarAbsolutePath(relative)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "NOT_FOUND", "message": "avatar not found"})
		return
	}

	raw, err := imgutils.ReadFileBytes(absPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "READ_FAILED", "message": err.Error()})
		return
	}
	processed, err := imgutils.ProcessImage(raw, targetSize, targetQuality)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "PROCESS_FAILED", "message": err.Error()})
		return
	}

	ct := "application/octet-stream"
	if kind, err := filetype.Match(processed); err == nil && kind != filetype.Unknown {
		ct = kind.MIME.Value
	}

	c.Header("Content-Type", ct)
	c.Header("Cache-Control", "public, max-age=31536000")
	c.Header("X-Image-Quality", strconv.Itoa(targetQuality))
	if targetSize > 0 {
		c.Header("X-Image-Size", strconv.Itoa(targetSize))
	}
	c.Status(http.StatusOK)
	_, _ = c.Writer.Write(processed)
}
