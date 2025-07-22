package editor

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"circle-center/globals"
	"circle-center/reader"
)

// CleanIconFileNames copies PNG files from srcDir to dstDir while replacing
// all dots '.' in the base file name (excluding extension) with underscores
// '_'. The function returns a slice of LocalIcon structures describing the
// newly written files. If dstDir does not exist, it will be created.
func CleanIconFileNames(srcDir, dstDir string) ([]globals.LocalIcon, error) {
	// Ensure destination directory exists
	if err := os.MkdirAll(dstDir, 0o755); err != nil {
		return nil, fmt.Errorf("create destination dir: %w", err)
	}

	// Read source icons
	icons, err := reader.ReadLocalIcons(srcDir)
	if err != nil {
		return nil, fmt.Errorf("read local icons: %w", err)
	}

	var results []globals.LocalIcon

	for _, icon := range icons {
		base := strings.TrimSuffix(icon.FileName, filepath.Ext(icon.FileName))
		cleanedBase := strings.ReplaceAll(base, ".", "_")
		newFileName := cleanedBase + filepath.Ext(icon.FileName) // .png
		srcPath := icon.FilePath
		dstPath := filepath.Join(dstDir, newFileName)

		// Copy file content
		if err := copyFile(srcPath, dstPath); err != nil {
			return nil, fmt.Errorf("copy %s to %s: %w", srcPath, dstPath, err)
		}

		results = append(results, globals.LocalIcon{
			FileName:    newFileName,
			PackageName: cleanedBase,
			FilePath:    dstPath,
		})
	}

	return results, nil
}

// copyFile copies a file from src to dst paths.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	// Ensure data flushed to disk
	if err := out.Sync(); err != nil {
		return err
	}

	return nil
}
