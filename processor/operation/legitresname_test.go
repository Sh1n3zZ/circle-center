package operation

import (
	"fmt"
	"os"
	"testing"
)

// TestCleanIconFileNames tests CleanIconFileNames using RENAME_SRC_DIR and
// RENAME_DST_DIR environment variables. If missing, skip.
func TestCleanIconFileNames(t *testing.T) {
	src := os.Getenv("RENAME_SRC_DIR")
	dst := os.Getenv("RENAME_DST_DIR")

	if src == "" || dst == "" {
		t.Skip("RENAME_SRC_DIR or RENAME_DST_DIR not set; skipping test")
	}

	icons, err := CleanIconFileNames(src, dst)
	if err != nil {
		t.Fatalf("clean icon file names failed: %v", err)
	}

	for _, ic := range icons {
		fmt.Printf("%s -> %s\n", ic.PackageName, ic.FilePath)
	}
}
