package manager

import (
	"context"
	"database/sql"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/h2non/filetype"

	"circle-center/globals/storage"
	accountsvc "circle-center/panel/account/svc"
	managerdb "circle-center/repository/sqlc/manager"
)

// IconIOService handles icon file upload and secure retrieval.
type IconIOService struct {
	queries  *managerdb.Queries
	auth     *accountsvc.AuthClient
	storage  *storage.IconStorage
	maxBytes int64
}

// NewIconIOService constructs a new IconIOService.
func NewIconIOService(db *sql.DB, authClient *accountsvc.AuthClient) (*IconIOService, error) {
	st, err := storage.NewIconStorage()
	if err != nil {
		return nil, err
	}
	return &IconIOService{
		queries:  managerdb.New(db),
		auth:     authClient,
		storage:  st,
		maxBytes: 5 * 1024 * 1024, // 5MB default limit for icon uploads
	}, nil
}

// ValidateAndSaveIcon validates auth, ensures the icon record exists, and saves file to storage.
// Returns the stored relative path (e.g., "icons/{project_id}/{drawable}.png").
func (s *IconIOService) ValidateAndSaveIcon(ctx context.Context, token string, projectID uint64, componentInfo string, fileBytes []byte) (string, error) {
	if s.auth == nil {
		return "", fmt.Errorf("auth client not initialized")
	}
	if len(fileBytes) == 0 {
		return "", fmt.Errorf("empty file")
	}
	if int64(len(fileBytes)) > s.maxBytes {
		return "", fmt.Errorf("file too large: max %d bytes", s.maxBytes)
	}

	// Validate token and load claims
	claims, err := s.auth.ValidateToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// Ensure project exists and is owned by current user
	p, err := s.queries.GetProjectByID(ctx, projectID)
	if err != nil {
		return "", fmt.Errorf("project not found")
	}
	if p.OwnerUserID != claims.UserID {
		return "", fmt.Errorf("forbidden")
	}

	comp := strings.TrimSpace(componentInfo)
	if comp == "" {
		return "", fmt.Errorf("component_info is required")
	}

	// Ensure icon row exists for this project/component
	icon, err := s.queries.GetIconByComponent(ctx, managerdb.GetIconByComponentParams{ProjectID: projectID, ComponentInfo: comp})
	if err != nil {
		return "", fmt.Errorf("icon record not found for component_info")
	}

	// Detect image kind and choose extension
	kind, err := filetype.Match(fileBytes)
	if err != nil {
		return "", fmt.Errorf("failed to detect file type: %w", err)
	}
	if kind == filetype.Unknown || kind.MIME.Type != "image" {
		return "", fmt.Errorf("unsupported file type")
	}
	ext := kind.Extension
	if ext == "jpeg" {
		ext = "jpg"
	}

	// Save file using drawable as the filename (overwrite allowed by design)
	rel, _, err := s.storage.SaveIcon(ctx, fileBytes, projectID, icon.Drawable, ext)
	if err != nil {
		return "", fmt.Errorf("failed to save icon: %w", err)
	}
	return rel, nil
}

// GetIconAbsolutePathSecure validates token and ownership based on relpath and returns absolute path.
// relpath must be like: icons/{project_id}/...  This method ensures the requester owns the project.
func (s *IconIOService) GetIconAbsolutePathSecure(ctx context.Context, token string, relpath string) (string, error) {
	if s.auth == nil {
		return "", fmt.Errorf("auth client not initialized")
	}

	claims, err := s.auth.ValidateToken(ctx, token)
	if err != nil {
		return "", fmt.Errorf("invalid token: %w", err)
	}

	// Normalize and parse project id from path
	norm := filepath.ToSlash(strings.TrimSpace(relpath))
	if len(norm) > 0 && norm[0] == '/' {
		norm = norm[1:]
	}
	parts := strings.Split(norm, "/")
	if len(parts) < 2 || parts[0] != "icons" {
		return "", fmt.Errorf("invalid icon path")
	}
	pidStr := parts[1]
	pid, err := strconv.ParseUint(pidStr, 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid project id in path")
	}

	// Verify ownership
	p, err := s.queries.GetProjectByID(ctx, pid)
	if err != nil {
		return "", fmt.Errorf("project not found")
	}
	if p.OwnerUserID != claims.UserID {
		return "", fmt.Errorf("forbidden")
	}

	return s.storage.AbsolutePath(norm)
}
