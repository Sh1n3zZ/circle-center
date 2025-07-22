package reader

import (
	"circle-center/globals"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ReadAppFilter is the public entry point for parsing an appfilter.xml file.
// It delegates the actual work to ParseAppFilterFile located in appfilter.go.
func ReadAppFilter(path string) ([]globals.Item, error) {
	return ParseAppFilterFile(path)
}

// ReadIconPack provides a public API to parse icon_pack.xml files.
func ReadIconPack(path string) (globals.IconPackResources, error) {
	return ParseIconPackFile(path)
}

// ReadLocalIcons parses a directory containing PNG icon files and returns the
// slice of globals.LocalIcon structures.
func ReadLocalIcons(dir string) ([]globals.LocalIcon, error) {
	return ParseIconDirectory(dir)
}

// reader handles POST /readfile which accepts a form-data file field named "file" containing an appfilter.xml.
// It parses the XML and returns the parsed items in JSON format.
func reader(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file field: " + err.Error()})
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open uploaded file: " + err.Error()})
		return
	}
	defer f.Close()

	fileType := c.PostForm("type")

	switch fileType {
	case "icon_pack":
		res, err := parseIconPackFromReader(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"icon_pack": res})
	case "", "appfilter":
		items, err := parseFromReader(f)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": items})
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unknown type: " + fileType})
	}
}
