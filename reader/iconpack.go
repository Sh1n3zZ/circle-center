package reader

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/beevik/etree"

	"circle-center/globals"
)

// parseIconPackFromReader parses an icon_pack.xml file using the provided io.Reader.
// It returns a fully populated IconPackResources structure that contains all
// <string-array> definitions.
func parseIconPackFromReader(r io.Reader) (globals.IconPackResources, error) {
	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(r); err != nil {
		return globals.IconPackResources{}, fmt.Errorf("failed to read xml: %w", err)
	}

	root := doc.SelectElement("resources")
	if root == nil {
		return globals.IconPackResources{}, fmt.Errorf("missing <resources> root element")
	}

	var arrays []globals.StringArray

	for _, sa := range root.SelectElements("string-array") {
		name := sa.SelectAttrValue("name", "")
		var items []string
		for _, itemEl := range sa.SelectElements("item") {
			items = append(items, strings.TrimSpace(itemEl.Text()))
		}
		arrays = append(arrays, globals.StringArray{
			Name:  name,
			Items: items,
		})
	}

	return globals.IconPackResources{Arrays: arrays}, nil
}

// ParseIconPackFile opens and parses an icon_pack.xml file from the given path.
// It wraps parseIconPackFromReader for convenience.
func ParseIconPackFile(path string) (globals.IconPackResources, error) {
	f, err := os.Open(path)
	if err != nil {
		return globals.IconPackResources{}, fmt.Errorf("cannot open file %s: %w", path, err)
	}
	defer f.Close()

	return parseIconPackFromReader(f)
}
