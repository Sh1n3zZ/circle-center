package storage

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"time"
)

// AvatarStorage provides helpers over a base Storage for avatar-specific paths.
type AvatarStorage struct {
	base Storage
}

// NewAvatarStorage creates an AvatarStorage using LocalStorage.
func NewAvatarStorage() (*AvatarStorage, error) {
	local, err := NewLocalStorage()
	if err != nil {
		return nil, err
	}
	return &AvatarStorage{base: local}, nil
}

// SaveAvatar saves avatar bytes under avatars/yyyy/mm with the provided filename.
// Returns the relative path (e.g., "avatars/2025/08/filename.jpg").
func (s *AvatarStorage) SaveAvatar(ctx context.Context, data []byte, fileName string, now time.Time) (string, string, error) {
	y, m, _ := now.Date()
	subDir := filepath.Join("avatars", fmt.Sprintf("%04d", y), fmt.Sprintf("%02d", int(m)))
	rel, abs, err := s.base.Save(ctx, data, subDir, fileName)
	if err != nil {
		return "", "", err
	}
	// Ensure returned relative path uses forward slashes
	return path.Clean(rel), abs, nil
}

// ReadAvatar reads avatar by its relative path under uploads root.
func (s *AvatarStorage) ReadAvatar(relativePath string) ([]byte, error) {
	return s.base.Read(relativePath)
}

// AbsolutePath resolves absolute path from relative avatar path.
func (s *AvatarStorage) AbsolutePath(relativePath string) (string, error) {
	return s.base.AbsolutePath(relativePath)
}
