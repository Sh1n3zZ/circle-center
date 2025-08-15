package storage

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

// IconStorage provides helpers over a base Storage for icon-specific paths.
type IconStorage struct {
	base Storage
}

// NewIconStorage creates an IconStorage using LocalStorage.
func NewIconStorage() (*IconStorage, error) {
	local, err := NewLocalStorage()
	if err != nil {
		return nil, err
	}
	return &IconStorage{base: local}, nil
}

// SaveIcon saves icon bytes under icons/{project_id}/{drawable}.{format}.
// Returns the relative path (e.g., "icons/123/my_icon.png").
func (s *IconStorage) SaveIcon(ctx context.Context, data []byte, projectID uint64, drawable string, format string) (string, string, error) {
	// Clean and validate drawable name
	cleanDrawable := filepath.Clean(drawable)
	if cleanDrawable == "." || cleanDrawable == ".." || cleanDrawable == "/" || cleanDrawable == "" {
		return "", "", fmt.Errorf("invalid drawable name: %s", drawable)
	}

	// Build subdirectory path: icons/{project_id}
	subDir := filepath.Join("icons", strconv.FormatUint(projectID, 10))

	// Build filename: {drawable}.{format}
	fileName := fmt.Sprintf("%s.%s", cleanDrawable, format)

	rel, abs, err := s.base.Save(ctx, data, subDir, fileName)
	if err != nil {
		return "", "", err
	}

	// Ensure returned relative path uses forward slashes
	return path.Clean(rel), abs, nil
}

// SaveIconWithTimestamp saves icon bytes under icons/{project_id}/{yyyy}/{mm}/{drawable}.{format}.
// This provides better organization for large projects with many icons.
func (s *IconStorage) SaveIconWithTimestamp(ctx context.Context, data []byte, projectID uint64, drawable string, format string, now time.Time) (string, string, error) {
	// Clean and validate drawable name
	cleanDrawable := filepath.Clean(drawable)
	if cleanDrawable == "." || cleanDrawable == ".." || cleanDrawable == "/" || cleanDrawable == "" {
		return "", "", fmt.Errorf("invalid drawable name: %s", drawable)
	}

	y, m, _ := now.Date()
	// Build subdirectory path: icons/{project_id}/{yyyy}/{mm}
	subDir := filepath.Join("icons", strconv.FormatUint(projectID, 10), fmt.Sprintf("%04d", y), fmt.Sprintf("%02d", int(m)))

	// Build filename: {drawable}.{format}
	fileName := fmt.Sprintf("%s.%s", cleanDrawable, format)

	rel, abs, err := s.base.Save(ctx, data, subDir, fileName)
	if err != nil {
		return "", "", err
	}

	// Ensure returned relative path uses forward slashes
	return path.Clean(rel), abs, nil
}

// ReadIcon reads icon by its relative path under uploads root.
func (s *IconStorage) ReadIcon(relativePath string) ([]byte, error) {
	return s.base.Read(relativePath)
}

// AbsolutePath resolves absolute path from relative icon path.
func (s *IconStorage) AbsolutePath(relativePath string) (string, error) {
	return s.base.AbsolutePath(relativePath)
}

// GetIconPath generates the expected relative path for an icon without saving it.
// This is useful for checking if an icon already exists or for generating paths.
func (s *IconStorage) GetIconPath(projectID uint64, drawable string, format string) string {
	cleanDrawable := filepath.Clean(drawable)
	return path.Join("icons", strconv.FormatUint(projectID, 10), fmt.Sprintf("%s.%s", cleanDrawable, format))
}

// GetIconPathWithTimestamp generates the expected relative path for an icon with timestamp organization.
func (s *IconStorage) GetIconPathWithTimestamp(projectID uint64, drawable string, format string, now time.Time) string {
	cleanDrawable := filepath.Clean(drawable)
	y, m, _ := now.Date()
	return path.Join("icons", strconv.FormatUint(projectID, 10), fmt.Sprintf("%04d", y), fmt.Sprintf("%02d", int(m)), fmt.Sprintf("%s.%s", cleanDrawable, format))
}
