package svc

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"circle-center/processor/operation"
)

// MergeRequest represents the request body for merge operations
type MergeRequest struct {
	// Selected components to merge (empty means merge all)
	Components []string `json:"components"`
	// If true, merge into first file; if false, merge into second file
	MergeIntoFirst bool `json:"merge_into_first"`
}

// MergeAppFilters handles POST /mergeappfilters which accepts two appfilter.xml
// files and merges them according to the specified options.
func MergeAppFilters(c *gin.Context) {
	// Get files from form data
	file1, err := c.FormFile("file1")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file1"})
		return
	}

	file2, err := c.FormFile("file2")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing file2"})
		return
	}

	// Parse merge options from form
	var req MergeRequest
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	// Create temporary directory for processing
	tmpDir, err := os.MkdirTemp("", "merge_*")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "create temp dir failed"})
		return
	}
	defer os.RemoveAll(tmpDir)

	// Save uploaded files
	file1Path := filepath.Join(tmpDir, "file1.xml")
	if err := c.SaveUploadedFile(file1, file1Path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save file1 failed"})
		return
	}

	file2Path := filepath.Join(tmpDir, "file2.xml")
	if err := c.SaveUploadedFile(file2, file2Path); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "save file2 failed"})
		return
	}

	// Prepare merge request
	mergeReq := operation.MergeRequest{
		FirstFile:      file1Path,
		SecondFile:     file2Path,
		OutputFile:     filepath.Join(tmpDir, "output.xml"),
		Mode:           operation.MergeSelected,
		MergeIntoFirst: req.MergeIntoFirst,
	}

	// If no components specified, merge all
	if len(req.Components) == 0 {
		mergeReq.Mode = operation.MergeAll
	} else {
		mergeReq.SelectedComponents = req.Components
	}

	// Perform merge in memory
	result, content, err := operation.MergeAppFiltersInMemory(mergeReq)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Return result
	c.JSON(http.StatusOK, gin.H{
		"result":  result,
		"content": content,
	})
}
