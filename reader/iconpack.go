package reader

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"circle-center/globals"
)

// parseIconPackFromReader parses an icon_pack.xml file using the provided io.Reader.
// It returns a fully populated IconPackResources structure that contains all
// <string-array> definitions.
func parseIconPackFromReader(r io.Reader) (globals.IconPackResources, error) {
	var res globals.IconPackResources
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&res); err != nil {
		return globals.IconPackResources{}, fmt.Errorf("failed to decode xml: %w", err)
	}
	return res, nil
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
