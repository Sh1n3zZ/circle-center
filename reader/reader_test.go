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
		// Output format: one line per item
		fmt.Printf("drawable: %s package-name: %s activity-name: %s component: %s\n", it.Drawable, it.PackageName, it.ActivityName, it.Component)
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
