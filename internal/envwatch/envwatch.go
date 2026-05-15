// Package envwatch provides file change detection for .env files,
// comparing the current state against a previously recorded checksum.
package envwatch

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

// ErrNoBaseline is returned when no baseline checksum has been recorded.
var ErrNoBaseline = errors.New("envwatch: no baseline checksum recorded")

// Status describes the change state of a watched file.
type Status struct {
	Path      string
	Changed   bool
	Baseline  string
	Current   string
	CheckedAt time.Time
}

// String returns a human-readable summary of the status.
func (s Status) String() string {
	if s.Changed {
		return fmt.Sprintf("%s: CHANGED (baseline=%s, current=%s)", s.Path, s.Baseline[:8], s.Current[:8])
	}
	return fmt.Sprintf("%s: unchanged", s.Path)
}

// Checksum computes the SHA-256 hex digest of the file at path.
func Checksum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", fmt.Errorf("envwatch: open %s: %w", path, err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("envwatch: hash %s: %w", path, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Watch compares the current checksum of path against baseline.
// Returns ErrNoBaseline if baseline is empty.
func Watch(path, baseline string) (Status, error) {
	if baseline == "" {
		return Status{}, ErrNoBaseline
	}

	current, err := Checksum(path)
	if err != nil {
		return Status{}, err
	}

	return Status{
		Path:      path,
		Changed:   current != baseline,
		Baseline:  baseline,
		Current:   current,
		CheckedAt: time.Now().UTC(),
	}, nil
}
