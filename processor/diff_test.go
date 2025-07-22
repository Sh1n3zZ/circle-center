package editor

import (
	"fmt"
	"os"
	"testing"
)

// TestFindMissingIcons tests FindMissingIcons by reading ICON_DIR_PATH and
// APPFILTER_PATH environment variables. If either is unset, the test is
// skipped. When run, it prints each missing package name on its own line.
func TestFindMissingIcons(t *testing.T) {
	dir := os.Getenv("ICON_DIR_PATH")
	iconPack := os.Getenv("APPFILTER_PATH")

	if dir == "" || iconPack == "" {
		t.Skip("ICON_DIR_PATH or APPFILTER_PATH env vars not set; skipping test")
	}

	missing, err := FindMissingIcons(dir, iconPack)
	if err != nil {
		t.Fatalf("finding missing icons failed: %v", err)
	}

	for _, name := range missing {
		fmt.Println(name)
	}
}
