package editor

import (
	"fmt"

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
