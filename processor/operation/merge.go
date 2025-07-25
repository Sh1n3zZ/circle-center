package operation

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/beevik/etree"

	"circle-center/globals"
	"circle-center/reader"
)

// MergeMode defines how items should be merged
type MergeMode int

const (
	// MergeAll merges all differences from both files
	MergeAll MergeMode = iota
	// MergeSelected merges only selected items based on their components
	MergeSelected
)

// MergeRequest represents a request to merge two appfilter.xml files
type MergeRequest struct {
	// Source file paths
	FirstFile  string
	SecondFile string
	// Output file path for the merged result
	OutputFile string
	// Mode determines how items should be merged
	Mode MergeMode
	// Components to merge when Mode is MergeSelected
	SelectedComponents []string
	// If true, merge items from second file into first file
	// If false, merge items from first file into second file
	MergeIntoFirst bool
}

// MergeResult contains information about the merge operation
type MergeResult struct {
	// Number of items merged from each file
	ItemsMerged int
	// Any items that failed to merge (e.g., duplicates)
	FailedItems []globals.Item
	// Final item count in the output file
	TotalItems int
}

// MergeAppFilters merges two appfilter.xml files according to the specified mode
// and selection criteria. It preserves the XML structure and comments of the
// target file while adding new items from the source file.
func MergeAppFilters(req MergeRequest) (*MergeResult, error) {
	// Read both files
	firstItems, err := reader.ReadAppFilter(req.FirstFile)
	if err != nil {
		return nil, fmt.Errorf("read first file: %w", err)
	}

	secondItems, err := reader.ReadAppFilter(req.SecondFile)
	if err != nil {
		return nil, fmt.Errorf("read second file: %w", err)
	}

	// Determine source and target items based on merge direction
	var sourceItems, targetItems []globals.Item
	if req.MergeIntoFirst {
		sourceItems = secondItems
		targetItems = firstItems
	} else {
		sourceItems = firstItems
		targetItems = secondItems
	}

	// Build set of existing components in target
	existingComponents := make(map[string]struct{})
	for _, item := range targetItems {
		existingComponents[item.Component] = struct{}{}
	}

	// Prepare result
	result := &MergeResult{
		ItemsMerged: 0,
		FailedItems: make([]globals.Item, 0),
	}

	// Create new XML document
	doc := etree.NewDocument()
	targetFile := req.SecondFile
	if req.MergeIntoFirst {
		targetFile = req.FirstFile
	}
	if err := doc.ReadFromFile(targetFile); err != nil {
		return nil, fmt.Errorf("read target file: %w", err)
	}

	root := doc.SelectElement("resources")
	if root == nil {
		return nil, fmt.Errorf("missing resources element in target file")
	}

	// Process source items
	for _, item := range sourceItems {
		// Skip if component already exists in target
		if _, exists := existingComponents[item.Component]; exists {
			result.FailedItems = append(result.FailedItems, item)
			continue
		}

		// Check if this item should be merged based on mode and selection
		shouldMerge := false
		switch req.Mode {
		case MergeAll:
			shouldMerge = true
		case MergeSelected:
			for _, comp := range req.SelectedComponents {
				if comp == item.Component {
					shouldMerge = true
					break
				}
			}
		}

		if !shouldMerge {
			continue
		}

		// Create new item element
		newItem := etree.NewElement("item")
		newItem.CreateAttr("component", item.Component)
		newItem.CreateAttr("drawable", item.Drawable)

		// If there's an app name, add it as a comment before the item
		if item.AppName != "" {
			comment := etree.NewComment(" " + item.AppName + " ")
			root.InsertChild(root.SelectElement("item"), comment)
		}

		// Add the new item at the beginning of resources
		root.InsertChild(root.SelectElement("item"), newItem)
		result.ItemsMerged++
		existingComponents[item.Component] = struct{}{}
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(req.OutputFile), 0o755); err != nil {
		return nil, fmt.Errorf("create output directory: %w", err)
	}

	// Write the merged file
	f, err := os.Create(req.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("create output file: %w", err)
	}
	defer f.Close()

	// Write XML declaration
	if _, err := io.WriteString(f, `<?xml version="1.0" encoding="utf-8"?>`+"\n"); err != nil {
		return nil, fmt.Errorf("write xml declaration: %w", err)
	}

	// Write the document with indentation
	doc.Indent(2)
	if _, err := doc.WriteTo(f); err != nil {
		return nil, fmt.Errorf("write merged content: %w", err)
	}

	result.TotalItems = len(existingComponents)
	return result, nil
}

// MergeAppFiltersInMemory performs the merge operation in memory without writing
// to a file. This is useful for API responses where we want to return the merged
// content directly.
func MergeAppFiltersInMemory(req MergeRequest) (*MergeResult, string, error) {
	result, err := MergeAppFilters(req)
	if err != nil {
		return nil, "", err
	}

	// Read the output file
	content, err := os.ReadFile(req.OutputFile)
	if err != nil {
		return nil, "", fmt.Errorf("read merged file: %w", err)
	}

	// Clean up the temporary output file
	os.Remove(req.OutputFile)

	return result, string(content), nil
}
