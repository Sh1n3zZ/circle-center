package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

// HashSHA256Hex returns the lowercase hex-encoded SHA-256 hash of the input string.
// Input will be trimmed before hashing.
func HashSHA256Hex(input string) string {
	trimmed := strings.TrimSpace(input)
	sum := sha256.Sum256([]byte(trimmed))
	return hex.EncodeToString(sum[:])
}

// GenerateSecureToken returns a hex-encoded cryptographically secure random token.
// byteLength controls entropy; e.g., 32 -> 64 hex chars.
func GenerateSecureToken(byteLength int) (string, error) {
	if byteLength <= 0 {
		byteLength = 32
	}
	buf := make([]byte, byteLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf), nil
}
