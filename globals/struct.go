package globals

import "encoding/xml"

// Resources maps the root <resources> element of appfilter.xml.
// Item elements inside will be unmarshalled into the Items slice.
type Resources struct {
	XMLName xml.Name `xml:"resources"`
	Items   []Item   `xml:"item"`
}

// Item represents one <item> element inside <resources>.
// It stores the raw component and drawable attributes from XML.
// PackageName and ActivityName fields are populated after parsing the
// Component attribute with helper utilities and are ignored during
// the XML unmarshalling (xml:"-").
type Item struct {
	Component    string `xml:"component,attr"`
	Drawable     string `xml:"drawable,attr"`
	PackageName  string `xml:"-"` // parsed package name
	ActivityName string `xml:"-"` // parsed activity name
}
