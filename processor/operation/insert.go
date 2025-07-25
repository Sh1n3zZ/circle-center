package operation

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/beevik/etree"

	"circle-center/globals"
)

// InsertItemsToIconPack inserts the provided drawable names into the specified
// <string-array name="arrayName"> inside the icon_pack.xml file located at
// iconPackPath. If the array does not exist, it will be created. Duplicate
// items are ignored.
func InsertItemsToIconPack(appFilterPath string, arrayName string, newItems []string) error {
	doc := etree.NewDocument()
	if err := doc.ReadFromFile(appFilterPath); err != nil {
		return fmt.Errorf("read icon pack xml: %w", err)
	}

	root := doc.SelectElement("resources")
	if root == nil {
		return fmt.Errorf("no <resources> root in %s", appFilterPath)
	}

	// Locate string-array
	var target *etree.Element
	for _, arr := range root.SelectElements("string-array") {
		if arr.SelectAttrValue("name", "") == arrayName {
			target = arr
			break
		}
	}

	// Create if not exists and insert at beginning of root children
	if target == nil {
		target = etree.NewElement("string-array")
		target.CreateAttr("name", arrayName)
		// Prepend to root children
		root.InsertChildAt(0, target)
	}

	// Build set of existing items texts
	existing := make(map[string]struct{})
	for _, item := range target.SelectElements("item") {
		existing[item.Text()] = struct{}{}
	}

	// Insert new items at beginning of target's child list in given order
	for i := len(newItems) - 1; i >= 0; i-- { // reverse to keep original order at top
		it := newItems[i]
		if _, present := existing[it]; present {
			continue
		}
		node := etree.NewElement("item")
		node.SetText(it)
		target.InsertChildAt(0, node)
	}

	doc.Indent(4)
	if err := doc.WriteToFile(appFilterPath); err != nil {
		return fmt.Errorf("write icon pack xml: %w", err)
	}

	return nil
}

// SyncIconsToAppFilter finds icons that exist in iconDir but are missing from
// appfilter.xml (comparing via FindMissingIcons), sanitises filenames via
// CleanIconFileNames into a "renamed" subdirectory, and inserts corresponding
// <item> elements into appfilter.xml.
func SyncIconsToAppFilter(iconDir, appFilterPath string) error {
	missingPkgs, err := FindMissingIcons(iconDir, appFilterPath)
	if err != nil {
		return err
	}
	if len(missingPkgs) == 0 {
		return nil // nothing to do
	}

	renamedDir := filepath.Join(iconDir, "renamed")
	// CleanIconFileNames copies **all** icons; we'll filter after copy
	cleanedIcons, err := CleanIconFileNames(iconDir, renamedDir)
	if err != nil {
		return err
	}

	// Build map for quick lookup of cleaned icons by package name
	cleanMap := make(map[string]globals.LocalIcon)
	for _, ic := range cleanedIcons {
		cleanMap[ic.PackageName] = ic
	}

	doc := etree.NewDocument()
	if err := doc.ReadFromFile(appFilterPath); err != nil {
		return fmt.Errorf("read appfilter xml: %w", err)
	}

	root := doc.SelectElement("resources")
	if root == nil {
		return fmt.Errorf("no <resources> root in appfilter")
	}

	// Build set of existing components to avoid duplicates
	existing := make(map[string]struct{})
	for _, item := range root.SelectElements("item") {
		comp := item.SelectAttrValue("component", "")
		existing[comp] = struct{}{}
	}

	// Insert new items at beginning (reverse iteration to preserve order)
	for i := len(missingPkgs) - 1; i >= 0; i-- {
		pkg := missingPkgs[i]
		drawableName := strings.ReplaceAll(pkg, ".", "_")
		componentStr := fmt.Sprintf("ComponentInfo{%s/TODO}", pkg)

		if _, ok := existing[componentStr]; ok {
			continue
		}

		newItem := etree.NewElement("item")
		newItem.CreateAttr("component", componentStr)
		newItem.CreateAttr("drawable", drawableName)

		root.InsertChildAt(0, newItem)
	}

	doc.Indent(4)
	if err := doc.WriteToFile(appFilterPath); err != nil {
		return fmt.Errorf("write appfilter xml: %w", err)
	}

	return nil
}
