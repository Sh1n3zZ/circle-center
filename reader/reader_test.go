package reader

import (
	"fmt"
	"os"
	"testing"
)

// TestReadAppFilter validates the ReadAppFilter function. It looks for the
// environment variable APPFILTER_PATH to determine the test file path. If the
// variable is not set, the test is skipped. When set, the parser runs and
// prints each parsed item on its own line in the format requested by the user.
func TestReadAppFilter(t *testing.T) {
	path := os.Getenv("APPFILTER_PATH")
	if path == "" {
		t.Skip("APPFILTER_PATH environment variable not set; skipping test")
	}

	items, err := ReadAppFilter(path)
	if err != nil {
		t.Fatalf("failed to parse appfilter: %v", err)
	}

	for _, it := range items {
		// Output format: one line per item including comment-based app name
		fmt.Printf("app-name: %s drawable: %s package-name: %s activity-name: %s component: %s\n", it.AppName, it.Drawable, it.PackageName, it.ActivityName, it.Component)
	}
}

// TestReadIconPack validates the ReadIconPack function using the ICONPACK_PATH
// environment variable. If the variable is unset, the test is skipped. When
// set, each <string-array> is printed on its own line for inspection.
func TestReadIconPack(t *testing.T) {
	path := os.Getenv("ICONPACK_PATH")
	if path == "" {
		t.Skip("ICONPACK_PATH environment variable not set; skipping test")
	}

	res, err := ReadIconPack(path)
	if err != nil {
		t.Fatalf("failed to parse icon pack: %v", err)
	}

	for _, arr := range res.Arrays {
		fmt.Printf("name: %s items: %v\n", arr.Name, arr.Items)
	}
}

// TestReadLocalIcons validates the ReadLocalIcons function using ICON_DIR_PATH
// environment variable which should point to a directory containing png files.
func TestReadLocalIcons(t *testing.T) {
	dir := os.Getenv("ICON_DIR_PATH")
	if dir == "" {
		t.Skip("ICON_DIR_PATH environment variable not set; skipping test")
	}

	icons, err := ReadLocalIcons(dir)
	if err != nil {
		t.Fatalf("failed to parse local icons: %v", err)
	}

	for _, ic := range icons {
		fmt.Printf("package-name: %s file-name: %s path: %s\n", ic.PackageName, ic.FileName, ic.FilePath)
	}
}
