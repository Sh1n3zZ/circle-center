package reader

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"

	"circle-center/globals"
)

// parseFromReader reads an appfilter.xml file from the provided io.Reader and
// returns the slice of parsed globals.Item structures.
func parseFromReader(r io.Reader) ([]globals.Item, error) {
	var res globals.Resources
	dec := xml.NewDecoder(r)
	if err := dec.Decode(&res); err != nil {
		return nil, fmt.Errorf("failed to decode xml: %w", err)
	}

	// Post-process each item to populate package and activity names.
	for idx := range res.Items {
		pkg, act := ParseComponentInfo(res.Items[idx].Component)
		res.Items[idx].PackageName = pkg
		res.Items[idx].ActivityName = act
	}

	return res.Items, nil
}

// ParseAppFilterFile opens an appfilter.xml at the given path and returns the
// parsed items. Caller is responsible for providing a valid file path.
func ParseAppFilterFile(path string) ([]globals.Item, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open file %s: %w", path, err)
	}
	defer f.Close()

	return parseFromReader(f)
}
