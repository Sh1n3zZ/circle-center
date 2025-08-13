package utils

import (
	"database/sql"
	"errors"
	"regexp"
	"strings"
)

// Slugify converts an input string to a URL-friendly slug
// - Lowercase all characters
// - Replace non-alphanumeric characters with '-'
// - Collapse multiple '-' and trim leading/trailing '-'
func Slugify(input string) string {
	str := strings.ToLower(strings.TrimSpace(input))
	// replace non-alphanumeric with hyphen
	re := regexp.MustCompile(`[^a-z0-9]+`)
	str = re.ReplaceAllString(str, "-")
	str = strings.Trim(str, "-")
	// collapse multiple hyphens
	re2 := regexp.MustCompile(`-+`)
	str = re2.ReplaceAllString(str, "-")
	return str
}

// AsBool attempts to interpret driver-returned types as boolean
// Supports common representations: bool, numeric, byte slice, and string
func AsBool(v interface{}) (bool, error) {
	switch val := v.(type) {
	case bool:
		return val, nil
	case int64:
		return val != 0, nil
	case int32:
		return val != 0, nil
	case int:
		return val != 0, nil
	case uint8:
		return val != 0, nil
	case []byte:
		s := strings.ToLower(string(val))
		return s == "1" || s == "t" || s == "true", nil
	case string:
		s := strings.ToLower(val)
		return s == "1" || s == "t" || s == "true", nil
	default:
		return false, errors.New("unsupported boolean type")
	}
}

// NullString returns empty string if the sql.NullString is invalid
func NullString(ns sql.NullString) string {
	if ns.Valid {
		return ns.String
	}
	return ""
}
