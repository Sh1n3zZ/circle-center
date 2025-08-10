package storage

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// GenerateUniqueFilename creates a unique filename using UUID + timestamp + checksum.
// The ext parameter should be provided without leading dot, e.g., "jpg", "png".
func GenerateUniqueFilename(ext string) (string, error) {
	u := uuid.New()
	ts := time.Now().UnixNano()

	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}

	// Compute checksum over UUID + timestamp + random bytes
	hasher := sha256.New()
	hasher.Write([]byte(u.String()))
	hasher.Write([]byte(fmt.Sprintf("%d", ts)))
	hasher.Write(randomBytes)
	checksum := hex.EncodeToString(hasher.Sum(nil))

	shortChecksum := checksum[:12]
	cleanExt := strings.TrimPrefix(ext, ".")
	if cleanExt != "" {
		return fmt.Sprintf("%s_%d_%s.%s", u.String(), ts, shortChecksum, cleanExt), nil
	}
	return fmt.Sprintf("%s_%d_%s", u.String(), ts, shortChecksum), nil
}
