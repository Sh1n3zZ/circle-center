package reader

import "circle-center/globals"

// ReadAppFilter is the public entry point for parsing an appfilter.xml file.
// It delegates the actual work to ParseAppFilterFile located in appfilter.go.
func ReadAppFilter(path string) ([]globals.Item, error) {
	return ParseAppFilterFile(path)
}
