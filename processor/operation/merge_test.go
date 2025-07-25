package operation

import (
	"os"
	"path/filepath"
	"testing"
)

func TestMergeAppFilters(t *testing.T) {
	// Create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "merge_test_*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Test cases
	tests := []struct {
		name           string
		firstFile      string
		secondFile     string
		mode           MergeMode
		selectedComps  []string
		mergeIntoFirst bool
		wantErr        bool
		wantMerged     int
	}{
		{
			name:       "merge all items",
			firstFile:  os.Getenv("APPFILTER_PATH"),
			secondFile: os.Getenv("APPFILTER_SECOND_PATH"),
			mode:       MergeAll,
			wantErr:    false,
		},
		{
			name:       "merge selected items",
			firstFile:  os.Getenv("APPFILTER_PATH"),
			secondFile: os.Getenv("APPFILTER_SECOND_PATH"),
			mode:       MergeSelected,
			selectedComps: []string{
				"ComponentInfo{com.example.app/com.example.app.MainActivity}",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.firstFile == "" || tt.secondFile == "" {
				t.Skip("APPFILTER_PATH or APPFILTER_SECOND_PATH not set")
			}

			outputFile := filepath.Join(tmpDir, "merged.xml")
			req := MergeRequest{
				FirstFile:          tt.firstFile,
				SecondFile:         tt.secondFile,
				OutputFile:         outputFile,
				Mode:               tt.mode,
				SelectedComponents: tt.selectedComps,
				MergeIntoFirst:     tt.mergeIntoFirst,
			}

			result, err := MergeAppFilters(req)
			if (err != nil) != tt.wantErr {
				t.Errorf("MergeAppFilters() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				t.Logf("Merged %d items, total items: %d", result.ItemsMerged, result.TotalItems)
				t.Logf("Failed items: %d", len(result.FailedItems))

				// Verify the output file exists and is readable
				if _, err := os.Stat(outputFile); err != nil {
					t.Errorf("Output file not created: %v", err)
				}
			}
		})
	}
}

func TestMergeAppFiltersInMemory(t *testing.T) {
	firstFile := os.Getenv("APPFILTER_PATH")
	secondFile := os.Getenv("APPFILTER_SECOND_PATH")

	if firstFile == "" || secondFile == "" {
		t.Skip("APPFILTER_PATH or APPFILTER_SECOND_PATH not set")
	}

	// Create temporary directory for test files
	tmpDir, err := os.MkdirTemp("", "merge_test_*")
	if err != nil {
		t.Fatalf("create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	req := MergeRequest{
		FirstFile:      firstFile,
		SecondFile:     secondFile,
		OutputFile:     filepath.Join(tmpDir, "temp.xml"),
		Mode:           MergeAll,
		MergeIntoFirst: true,
	}

	result, content, err := MergeAppFiltersInMemory(req)
	if err != nil {
		t.Fatalf("MergeAppFiltersInMemory() error = %v", err)
	}

	t.Logf("Merged %d items, total items: %d", result.ItemsMerged, result.TotalItems)
	t.Logf("Content length: %d bytes", len(content))

	// Verify the content is valid XML
	if len(content) == 0 {
		t.Error("Empty content returned")
	}
	if content[0:5] != "<?xml" {
		t.Error("Content does not start with XML declaration")
	}
}
