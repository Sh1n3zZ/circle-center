package operation

import (
	"os"
	"testing"
)

// TestSyncIconsToAppFilter orchestrates the full sync flow using environment
// variables configured in .vscode/settings.json. If any variable is missing, the
// test is skipped.
func TestSyncIconsToAppFilter(t *testing.T) {
	iconDir := os.Getenv("ICON_DIR_PATH")
	appFilter := os.Getenv("APPFILTER_PATH")

	if iconDir == "" || appFilter == "" {
		t.Skip("Required environment variables not set; skipping test")
	}

	if err := SyncIconsToAppFilter(iconDir, appFilter); err != nil {
		t.Fatalf("sync icons to appfilter failed: %v", err)
	}
}
