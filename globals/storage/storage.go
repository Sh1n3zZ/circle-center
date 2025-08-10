package storage

import (
	"context"
	"errors"
	"fmt"
	"os"
	pathpkg "path"
	"path/filepath"
)

// Storage defines an abstract storage interface for saving and reading files.
// This abstraction allows swapping local storage with other providers (e.g., S3).
type Storage interface {
	// Save writes data under base uploads directory joined with subDir and fileName.
	// Returns the relative path (using forward slashes) from the uploads root (e.g., "avatars/2025/08/uuid_xxx.jpg")
	// and the absolute filesystem path.
	Save(ctx context.Context, data []byte, subDir string, fileName string) (relativePath string, absolutePath string, err error)

	// Read reads the file content by its relative path (using forward slashes) from uploads root.
	Read(relativePath string) ([]byte, error)

	// AbsolutePath resolves an absolute filesystem path from a relative uploads path.
	AbsolutePath(relativePath string) (string, error)
}

// LocalStorage implements Storage using the local filesystem under
// <projectRoot>/storage/uploads.
type LocalStorage struct {
	uploadsRoot string
}

// NewLocalStorage creates a LocalStorage rooted at <projectRoot>/storage/uploads.
func NewLocalStorage() (*LocalStorage, error) {
	root, err := FindProjectRoot()
	if err != nil {
		return nil, err
	}
	uploads := filepath.Join(root, "storage", "uploads")
	if err := os.MkdirAll(uploads, 0755); err != nil {
		return nil, fmt.Errorf("failed to create uploads directory: %w", err)
	}
	return &LocalStorage{uploadsRoot: uploads}, nil
}

// Save implements Storage.Save.
func (s *LocalStorage) Save(_ context.Context, data []byte, subDir string, fileName string) (string, string, error) {
	if fileName == "" {
		return "", "", errors.New("fileName is required")
	}
	// Normalize subDir to use OS path for writing, but ensure returned relative path uses forward slashes.
	cleanedSubDir := filepath.Clean(subDir)
	// Prevent breaking out of uploads root
	if cleanedSubDir == ".." || cleanedSubDir == "." || cleanedSubDir == string(os.PathSeparator) {
		cleanedSubDir = ""
	}
	targetDir := filepath.Join(s.uploadsRoot, cleanedSubDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", "", fmt.Errorf("failed to create target directory: %w", err)
	}

	absPath := filepath.Join(targetDir, fileName)
	if err := os.WriteFile(absPath, data, 0644); err != nil {
		return "", "", fmt.Errorf("failed to write file: %w", err)
	}

	// Build relative path using forward slashes for URLs
	rel := fileName
	if cleanedSubDir != "" {
		rel = pathpkg.Join(filepath.ToSlash(cleanedSubDir), fileName)
	}
	return rel, absPath, nil
}

// Read implements Storage.Read.
func (s *LocalStorage) Read(relativePath string) ([]byte, error) {
	abs, err := s.AbsolutePath(relativePath)
	if err != nil {
		return nil, err
	}
	return os.ReadFile(abs)
}

// AbsolutePath implements Storage.AbsolutePath.
func (s *LocalStorage) AbsolutePath(relativePath string) (string, error) {
	// Normalize and prevent traversal
	cleaned := pathpkg.Clean(filepath.ToSlash(relativePath))
	if cleaned == "." || cleaned == ".." || cleaned == "/" || cleaned == "" {
		return "", errors.New("invalid relative path")
	}
	abs := filepath.Join(s.uploadsRoot, filepath.FromSlash(cleaned))
	return abs, nil
}

// FindProjectRoot walks up from the working directory to locate the directory containing go.mod.
func FindProjectRoot() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir := wd
	for {
		mod := filepath.Join(dir, "go.mod")
		if _, err := os.Stat(mod); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("go.mod not found from working directory: %s", wd)
}
