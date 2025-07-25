package operation

import (
	"fmt"

	"circle-center/globals"
	"circle-center/reader"
)

// FindMissingIcons compares the icons present in a directory with the entries
// listed inside an icon_pack.xml file. It returns a slice of package/icon names
// that are present in the directory (icons) but absent from any <string-array>
// in the icon pack definition.
func FindMissingIcons(iconDir string, appFilterPath string) ([]string, error) {
	// Parse local icons
	localIcons, err := reader.ReadLocalIcons(iconDir)
	if err != nil {
		return nil, fmt.Errorf("read local icons: %w", err)
	}

	// Parse appfilter.xml
	items, err := reader.ReadAppFilter(appFilterPath)
	if err != nil {
		return nil, fmt.Errorf("read appfilter: %w", err)
	}

	// Print all local icon package names
	// fmt.Println("[DEBUG] local icons package names:")
	for _, icon := range localIcons {
		fmt.Println(icon.PackageName)
	}

	// Print all appfilter.xml package names
	// fmt.Println("[DEBUG] appfilter.xml package names:")
	for _, item := range items {
		fmt.Println(item.PackageName)
	}

	// fmt.Println("[DEBUG] local icons count: ", len(localIcons))
	// fmt.Println("[DEBUG] appfilter.xml count: ", len(items))

	// Build a set of all package names present in appfilter.xml
	appFilterSet := make(map[string]struct{})
	for _, item := range items {
		appFilterSet[item.PackageName] = struct{}{}
	}

	// Collect missing items
	missingSet := make(map[string]struct{})
	for _, icon := range localIcons {
		if _, ok := appFilterSet[icon.PackageName]; !ok {
			missingSet[icon.PackageName] = struct{}{}
		}
	}

	// Convert set to slice
	var missing []string
	for name := range missingSet {
		missing = append(missing, name)
	}

	return missing, nil
}

// DiffAppFilters compares two appfilter.xml files and returns the differences.
// It returns three slices:
// - onlyInFirst: items that exist only in the first appfilter.xml
// - onlyInSecond: items that exist only in the second appfilter.xml
// - common: items that exist in both files
func DiffAppFilters(firstPath, secondPath string) (onlyInFirst, onlyInSecond, common []globals.Item, err error) {
	// Parse first appfilter.xml
	firstItems, err := reader.ReadAppFilter(firstPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read first appfilter: %w", err)
	}

	// Parse second appfilter.xml
	secondItems, err := reader.ReadAppFilter(secondPath)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("read second appfilter: %w", err)
	}

	// Build sets for efficient lookup
	firstSet := make(map[string]globals.Item)
	secondSet := make(map[string]globals.Item)

	for _, item := range firstItems {
		firstSet[item.Component] = item
	}

	for _, item := range secondItems {
		secondSet[item.Component] = item
	}

	// Find items only in first
	for _, item := range firstItems {
		if _, exists := secondSet[item.Component]; !exists {
			onlyInFirst = append(onlyInFirst, item)
		}
	}

	// Find items only in second
	for _, item := range secondItems {
		if _, exists := firstSet[item.Component]; !exists {
			onlyInSecond = append(onlyInSecond, item)
		}
	}

	// Find common items
	for _, item := range firstItems {
		if secondItem, exists := secondSet[item.Component]; exists {
			common = append(common, item)
			// Note: we use the first file's item for common items
			_ = secondItem // avoid unused variable warning
		}
	}

	return onlyInFirst, onlyInSecond, common, nil
}
