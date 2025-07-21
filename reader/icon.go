package reader

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"circle-center/globals"
)

// ParseIconDirectory walks the provided directory and collects information
// about PNG files. Each PNG file is mapped into a globals.LocalIcon structure
// where PackageName is the file name without the .png extension.
func ParseIconDirectory(dir string) ([]globals.LocalIcon, error) {
	var icons []globals.LocalIcon

	walkFn := func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Skip nested directories; only consider files in root dir.
			if path != dir {
				return filepath.SkipDir
			}
			return nil
		}

		if strings.HasSuffix(strings.ToLower(d.Name()), ".png") {
			base := strings.TrimSuffix(d.Name(), filepath.Ext(d.Name()))
			icons = append(icons, globals.LocalIcon{
				FileName:    d.Name(),
				PackageName: base,
				FilePath:    path,
			})
		}
		return nil
	}

	if err := filepath.WalkDir(dir, walkFn); err != nil {
		return nil, fmt.Errorf("error scanning directory %s: %w", dir, err)
	}

	return icons, nil
}
