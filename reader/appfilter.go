package reader

import (
	"fmt"
	"io"
	"os"

	"github.com/beevik/etree"

	"circle-center/globals"
)

// ParseFromReader reads an appfilter.xml file from the provided io.Reader and
// returns the slice of parsed globals.Item structures.
func ParseFromReader(r io.Reader) ([]globals.Item, error) {
	doc := etree.NewDocument()
	if _, err := doc.ReadFrom(r); err != nil {
		return nil, fmt.Errorf("failed to read xml: %w", err)
	}

	root := doc.SelectElement("resources")
	if root == nil {
		return nil, fmt.Errorf("missing <resources> root element")
	}

	var items []globals.Item
	// Iterate through all child tokens of <resources> to capture comment nodes
	// that appear directly before each <item>. The last encountered comment
	// prior to an <item> element is taken as the application's human-readable
	// name (e.g. "Browser").
	var currentComment string
	for _, tok := range root.Child {
		switch t := tok.(type) {
		case *etree.Comment:
			// Store trimmed comment text for the next <item> element.
			currentComment = ParseCommentText(t.Data)
		case *etree.Element:
			if t.Tag == "item" {
				comp := t.SelectAttrValue("component", "")
				drawable := t.SelectAttrValue("drawable", "")

				pkg, act := ParseComponentInfo(comp)

				items = append(items, globals.Item{
					Component:    comp,
					Drawable:     drawable,
					PackageName:  pkg,
					ActivityName: act,
					AppName:      currentComment,
				})

				// Reset comment after use to avoid incorrectly assigning it to
				// subsequent items when they lack an explicit comment.
				currentComment = ""
			}
		}
	}

	return items, nil
}

// ParseAppFilterFile opens an appfilter.xml at the given path and returns the
// parsed items. Caller is responsible for providing a valid file path.
func ParseAppFilterFile(path string) ([]globals.Item, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", path, err)
	}
	defer f.Close()

	return ParseFromReader(f)
}
