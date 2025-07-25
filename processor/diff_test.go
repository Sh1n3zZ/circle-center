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

// TestDiffAppFilters tests DiffAppFilters by reading APPFILTER_PATH and
// APPFILTER_SECOND_PATH environment variables. If either is unset, the test is
// skipped. When run, it prints the differences between the two appfilter.xml files.
func TestDiffAppFilters(t *testing.T) {
	firstPath := os.Getenv("APPFILTER_PATH")
	secondPath := os.Getenv("APPFILTER_SECOND_PATH")

	if firstPath == "" || secondPath == "" {
		t.Skip("APPFILTER_PATH or APPFILTER_SECOND_PATH env vars not set; skipping test")
	}

	onlyInFirst, onlyInSecond, common, err := DiffAppFilters(firstPath, secondPath)
	if err != nil {
		t.Fatalf("diffing appfilters failed: %v", err)
	}

	fmt.Printf("=== Items only in first file (%s) ===\n", firstPath)
	for _, item := range onlyInFirst {
		fmt.Printf("component: %s drawable: %s package: %s activity: %s\n",
			item.Component, item.Drawable, item.PackageName, item.ActivityName)
	}

	fmt.Printf("\n=== Items only in second file (%s) ===\n", secondPath)
	for _, item := range onlyInSecond {
		fmt.Printf("component: %s drawable: %s package: %s activity: %s\n",
			item.Component, item.Drawable, item.PackageName, item.ActivityName)
	}

	fmt.Printf("\n=== Common items (count: %d) ===\n", len(common))
	for _, item := range common {
		fmt.Printf("component: %s drawable: %s package: %s activity: %s\n",
			item.Component, item.Drawable, item.PackageName, item.ActivityName)
	}
}
