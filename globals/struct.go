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
	AppName      string `xml:"-"` // app name parsed from XML comment
	PackageName  string `xml:"-"` // parsed package name
	ActivityName string `xml:"-"` // parsed activity name
}

// StringArray represents a <string-array> element with its name attribute and
// nested <item> strings.
type StringArray struct {
	Name  string   `xml:"name,attr"`
	Items []string `xml:"item"`
}

// IconPackResources corresponds to the root <resources> element inside
// icon_pack.xml. All string-array children are collected into the Arrays slice.
type IconPackResources struct {
	XMLName xml.Name      `xml:"resources"`
	Arrays  []StringArray `xml:"string-array"`
}

// LocalIcon represents a PNG icon discovered in an icons directory. The
// FileName includes extension; PackageName is derived from the base name
// (without extension). FilePath holds the full path to the file.
type LocalIcon struct {
	FileName    string
	PackageName string
	FilePath    string
}
