package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/h2non/filetype"
)

// ProcessImage decodes image bytes, optionally resizes to a square of targetSize,
// and re-encodes using a format derived from the original data. For JPEG, quality
// (1-100) is applied; for other formats, default encoding is used.
func ProcessImage(imageBytes []byte, targetSize int, quality int) ([]byte, error) {
	if len(imageBytes) == 0 {
		return nil, fmt.Errorf("empty image data")
	}

	kind, err := filetype.Match(imageBytes)
	if err != nil {
		return nil, fmt.Errorf("detect type: %w", err)
	}

	format := imaging.JPEG
	ext := strings.ToLower(kind.Extension)
	switch ext {
	case "jpg", "jpeg":
		format = imaging.JPEG
	case "png":
		format = imaging.PNG
	case "gif":
		format = imaging.GIF
	case "tiff", "tif":
		format = imaging.TIFF
	case "bmp":
		format = imaging.BMP
	default:
		format = imaging.JPEG
	}

	src, err := imaging.Decode(bytes.NewReader(imageBytes))
	if err != nil {
		return nil, fmt.Errorf("decode: %w", err)
	}

	dst := src
	if targetSize > 0 {
		dst = imaging.Fill(src, targetSize, targetSize, imaging.Center, imaging.Lanczos)
	}

	var buf bytes.Buffer
	var w io.Writer = &buf
	switch format {
	case imaging.JPEG:
		if err := imaging.Encode(w, dst, imaging.JPEG, imaging.JPEGQuality(quality)); err != nil {
			return nil, fmt.Errorf("encode jpeg: %w", err)
		}
	case imaging.PNG:
		if err := imaging.Encode(w, dst, imaging.PNG); err != nil {
			return nil, fmt.Errorf("encode png: %w", err)
		}
	case imaging.GIF:
		if err := imaging.Encode(w, dst, imaging.GIF); err != nil {
			return nil, fmt.Errorf("encode gif: %w", err)
		}
	case imaging.TIFF:
		if err := imaging.Encode(w, dst, imaging.TIFF); err != nil {
			return nil, fmt.Errorf("encode tiff: %w", err)
		}
	case imaging.BMP:
		if err := imaging.Encode(w, dst, imaging.BMP); err != nil {
			return nil, fmt.Errorf("encode bmp: %w", err)
		}
	default:
		if err := imaging.Encode(w, dst, imaging.JPEG, imaging.JPEGQuality(quality)); err != nil {
			return nil, fmt.Errorf("encode default jpeg: %w", err)
		}
	}

	return buf.Bytes(), nil
}

// ReadFileBytes reads the whole file content into memory.
func ReadFileBytes(abs string) ([]byte, error) {
	// Keep as a tiny helper to isolate IO from handlers
	return os.ReadFile(abs)
}
