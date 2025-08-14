package manager

import (
	"database/sql"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	svc "circle-center/panel/manager/svc"
)

// TokenHandler wires HTTP to TokenService
type TokenHandler struct {
	service *svc.TokenService
}

// NewTokenHandler builds a token handler
func NewTokenHandler(db *sql.DB) *TokenHandler {
	return &TokenHandler{service: svc.NewTokenService(db)}
}

type createTokenReq struct {
	Name string `json:"name"`
}

// GET /manager/projects/:id/tokens
func (h *TokenHandler) List(c *gin.Context) {
	projectIDStr := c.Param("id")
	pid, _ := strconv.ParseUint(strings.TrimSpace(projectIDStr), 10, 64)
	items, err := h.service.ListTokens(c.Request.Context(), pid)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "LIST_TOKENS_FAILED", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true, "data": items})
}

// POST /manager/projects/:id/tokens
func (h *TokenHandler) Create(c *gin.Context) {
	projectIDStr := c.Param("id")
	pid, _ := strconv.ParseUint(strings.TrimSpace(projectIDStr), 10, 64)
	var body createTokenReq
	_ = c.ShouldBindJSON(&body)
	plain, tokenID, err := h.service.CreateToken(c.Request.Context(), pid, body.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "CREATE_TOKEN_FAILED", "message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"success": true, "data": gin.H{"token_id": tokenID, "token": plain}})
}

// DELETE /manager/projects/:id/tokens/:tokenId
func (h *TokenHandler) Delete(c *gin.Context) {
	projectIDStr := c.Param("id")
	tokenIDStr := c.Param("tokenId")
	pid, _ := strconv.ParseUint(strings.TrimSpace(projectIDStr), 10, 64)
	tid, _ := strconv.ParseUint(strings.TrimSpace(tokenIDStr), 10, 64)
	if err := h.service.DeleteToken(c.Request.Context(), pid, tid); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "DELETE_TOKEN_FAILED", "message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": true})
}
