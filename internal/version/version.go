// Package version provides utilities for managing and comparing
// envelope versions to support safe re-sealing and change detection.
package version

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"
)

// Info holds metadata about a sealed envelope version.
type Info struct {
	Version   int       `json:"version"`
	SealedAt  time.Time `json:"sealed_at"`
	Checksum  string    `json:"checksum"`
	SealedBy  string    `json:"sealed_by,omitempty"`
}

// New creates a new Info for the given plaintext content.
func New(plaintext []byte, sealedBy string) Info {
	return Info{
		Version:  1,
		SealedAt: time.Now().UTC(),
		Checksum: Checksum(plaintext),
		SealedBy: sealedBy,
	}
}

// Bump returns a new Info with the version incremented and checksum updated.
func Bump(prev Info, plaintext []byte) Info {
	return Info{
		Version:  prev.Version + 1,
		SealedAt: time.Now().UTC(),
		Checksum: Checksum(plaintext),
		SealedBy: prev.SealedBy,
	}
}

// Checksum computes a SHA-256 hex digest of the given data.
func Checksum(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:])
}

// VerifyFile reads the file at path and checks its checksum against expected.
func VerifyFile(path, expected string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("version: reading file %q: %w", path, err)
	}
	actual := Checksum(data)
	if actual != expected {
		return fmt.Errorf("version: checksum mismatch for %q: expected %s, got %s", path, expected, actual)
	}
	return nil
}

// HasChanged returns true if the checksum of plaintext differs from prev.Checksum.
func HasChanged(prev Info, plaintext []byte) bool {
	return Checksum(plaintext) != prev.Checksum
}
