package svc

import (
	"circle-center/globals"
	"circle-center/processor/operation"
	"io"
	"net/http"

	"circle-center/reader"

	"github.com/gin-gonic/gin"
)

// diffAppFilters handles POST /diffappfilters which accepts two form-data file fields:
// "file1" and "file2" containing appfilter.xml files to compare.
// It returns the differences in JSON format.
func DiffAppFilters(c *gin.Context) {
	// Get first file
	file1Header, err := c.FormFile("file1")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file1 field: " + err.Error()})
		return
	}

	file1, err := file1Header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file1: " + err.Error()})
		return
	}
	defer file1.Close()

	// Get second file
	file2Header, err := c.FormFile("file2")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file2 field: " + err.Error()})
		return
	}

	file2, err := file2Header.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open file2: " + err.Error()})
		return
	}
	defer file2.Close()

	// Parse both files
	firstItems, err := reader.ParseFromReader(file1)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parse file1 failed: " + err.Error()})
		return
	}

	// Reset file2 reader position
	file2.Seek(0, io.SeekStart)
	secondItems, err := reader.ParseFromReader(file2)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "parse file2 failed: " + err.Error()})
		return
	}

	// Build sets for efficient lookup
	firstSet := make(map[string]globals.Item)
	secondSet := make(map[string]globals.Item)

	for _, item := range firstItems {
		firstSet[item.Component] = item
	}

	for _, item := range secondItems {
		secondSet[item.Component] = item
	}

	// Find differences
	var onlyInFirst, onlyInSecond, common []globals.Item

	// Items only in first file
	for _, item := range firstItems {
		if _, exists := secondSet[item.Component]; !exists {
			onlyInFirst = append(onlyInFirst, item)
		}
	}

	// Items only in second file
	for _, item := range secondItems {
		if _, exists := firstSet[item.Component]; !exists {
			onlyInSecond = append(onlyInSecond, item)
		}
	}

	// Common items
	for _, item := range firstItems {
		if _, exists := secondSet[item.Component]; exists {
			common = append(common, item)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"only_in_first":  onlyInFirst,
		"only_in_second": onlyInSecond,
		"common":         common,
		"summary": gin.H{
			"first_count":       len(firstItems),
			"second_count":      len(secondItems),
			"only_first_count":  len(onlyInFirst),
			"only_second_count": len(onlyInSecond),
			"common_count":      len(common),
		},
	})
}

// diffIcons handles POST /difficons which accepts form-data with:
// "icon_dir" (directory path) and "icon_pack" (file path) to compare local icons
// with icon_pack.xml entries.
func DiffIcons(c *gin.Context) {
	iconDir := c.PostForm("icon_dir")
	if iconDir == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing icon_dir field"})
		return
	}

	iconPackPath := c.PostForm("icon_pack")
	if iconPackPath == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing icon_pack field"})
		return
	}

	missing, err := operation.FindMissingIcons(iconDir, iconPackPath)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "diff failed: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"missing_icons": missing,
		"count":         len(missing),
	})
}
