package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/h2non/filetype"

	configure "circle-center/globals/configure"
	"circle-center/globals/storage"
	accountdb "circle-center/repository/sqlc/account"
)

// AvatarService handles avatar upload and retrieval business logic.
type AvatarService struct {
	queries  *accountdb.Queries
	auth     *AuthClient
	storage  *storage.AvatarStorage
	maxBytes int64
}

// NewAvatarService creates a new AvatarService.
func NewAvatarService(db *sql.DB, authClient *AuthClient) (*AvatarService, error) {
	st, err := storage.NewAvatarStorage()
	if err != nil {
		return nil, err
	}
	cfg := configure.GetConfig()
	maxBytes := int64(2 * 1024 * 1024)
	if cfg != nil && cfg.Avatar.MaxUploadBytes > 0 {
		maxBytes = cfg.Avatar.MaxUploadBytes
	}
	return &AvatarService{
		queries:  accountdb.New(db),
		auth:     authClient,
		storage:  st,
		maxBytes: maxBytes,
	}, nil
}

// UploadResult represents the result after avatar upload.
type UploadResult struct {
	Path string `json:"path"`
	URL  string `json:"url"`
}

// ValidateAndSaveAvatar validates the file and saves it. It also updates user's avatar_url.
func (s *AvatarService) ValidateAndSaveAvatar(ctx context.Context, userID uint64, fileBytes []byte) (*UploadResult, error) {
	if len(fileBytes) == 0 {
		return nil, fmt.Errorf("empty file")
	}
	if int64(len(fileBytes)) > s.maxBytes {
		return nil, fmt.Errorf("file too large: max %d bytes", s.maxBytes)
	}

	// Detect MIME type
	kind, err := filetype.Match(fileBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to detect file type: %w", err)
	}
	if kind == filetype.Unknown {
		return nil, fmt.Errorf("unsupported file type")
	}
	if kind.MIME.Type != "image" {
		return nil, fmt.Errorf("file is not an image")
	}

	// Map to extension
	ext := kind.Extension
	if ext == "jpeg" {
		ext = "jpg"
	}

	// Generate unique filename
	name, err := storage.GenerateUniqueFilename(ext)
	if err != nil {
		return nil, fmt.Errorf("failed to generate filename: %w", err)
	}

	// Save avatar
	now := time.Now()
	rel, _, err := s.storage.SaveAvatar(ctx, fileBytes, name, now)
	if err != nil {
		return nil, fmt.Errorf("failed to save avatar: %w", err)
	}

	// Read current user to preserve other fields
	user, err := s.queries.GetUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to load user: %w", err)
	}

	// Update user profile avatar_url (store relative path) preserving existing fields
	err = s.queries.UpdateUserProfile(ctx, accountdb.UpdateUserProfileParams{
		DisplayName: user.DisplayName,
		AvatarUrl:   sql.NullString{String: rel, Valid: true},
		Phone:       user.Phone,
		Locale:      user.Locale,
		Timezone:    user.Timezone,
		ID:          userID,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to update user profile: %w", err)
	}

	// Build URL for access via public route
	cfg := configure.GetConfig()
	base := ""
	if cfg != nil && cfg.Server.Host != "" && cfg.Server.Port != 0 {
		// Note: External base URL is not guaranteed; leave URL empty or construct relative
		base = ""
	}

	return &UploadResult{
		Path: rel,
		URL:  base,
	}, nil
}

// ParseTargetSizeAndQuality parses size and quality from query, enforced by config.
func (s *AvatarService) ParseTargetSizeAndQuality(sizeParam, qualityParam string) (int, int, error) {
	maxQuality := 90
	cfg := configure.GetConfig()
	if cfg != nil && cfg.Avatar.MaxImageQuality > 0 {
		maxQuality = cfg.Avatar.MaxImageQuality
	}

	// Defaults: size 0 means keep original; quality 0 means default 80 bounded by maxQuality
	targetSize := 0
	targetQuality := 0

	if sizeParam != "" {
		var parsed int
		_, err := fmt.Sscanf(sizeParam, "%d", &parsed)
		if err != nil || parsed < 0 {
			return 0, 0, errors.New("invalid size")
		}
		targetSize = parsed
	}

	if qualityParam != "" {
		var parsed int
		_, err := fmt.Sscanf(qualityParam, "%d", &parsed)
		if err != nil || parsed < 1 || parsed > 100 {
			return 0, 0, errors.New("invalid quality")
		}
		if parsed > maxQuality {
			parsed = maxQuality
		}
		targetQuality = parsed
	}

	if targetQuality == 0 {
		// Set a sensible default under the cap
		q := 80
		if q > maxQuality {
			q = maxQuality
		}
		targetQuality = q
	}

	return targetSize, targetQuality, nil
}

// GetAvatarAbsolutePath resolves absolute path from stored relative path.
func (s *AvatarService) GetAvatarAbsolutePath(rel string) (string, error) {
	return s.storage.AbsolutePath(rel)
}

// ValidateUserToken validates token and returns user claims, delegating to AuthClient.
func (s *AvatarService) ValidateUserToken(ctx context.Context, tokenString string) (*UserClaims, error) {
	if s.auth == nil {
		return nil, fmt.Errorf("auth client not initialized")
	}
	return s.auth.ValidateToken(ctx, tokenString)
}
