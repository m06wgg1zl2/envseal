// Package snapshot provides functionality to capture and compare
// sealed environment file states over time, enabling rollback detection
// and history tracking.
package snapshot

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Entry represents a single snapshot of a sealed env file.
type Entry struct {
	Version  int       `json:"version"`
	Checksum string    `json:"checksum"`
	SealedAt time.Time `json:"sealed_at"`
	Note     string    `json:"note,omitempty"`
}

// History holds a list of snapshot entries for a sealed file.
type History struct {
	Entries []Entry `json:"entries"`
}

// DefaultSnapshotPath returns the default path for a snapshot file
// relative to the given sealed file path.
func DefaultSnapshotPath(sealedFilePath string) string {
	dir := filepath.Dir(sealedFilePath)
	return filepath.Join(dir, ".envseal_snapshots.json")
}

// Record appends a new snapshot entry to the history file.
// If the history file does not exist, it is created.
func Record(snapshotPath string, entry Entry) error {
	if entry.Version <= 0 {
		return errors.New("snapshot: version must be greater than zero")
	}
	if entry.Checksum == "" {
		return errors.New("snapshot: checksum must not be empty")
	}
	if entry.SealedAt.IsZero() {
		entry.SealedAt = time.Now().UTC()
	}

	h, err := load(snapshotPath)
	if err != nil {
		return fmt.Errorf("snapshot: load history: %w", err)
	}

	h.Entries = append(h.Entries, entry)

	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return fmt.Errorf("snapshot: marshal: %w", err)
	}

	if err := os.WriteFile(snapshotPath, data, 0o644); err != nil {
		return fmt.Errorf("snapshot: write file: %w", err)
	}
	return nil
}

// LoadHistory reads and returns all snapshot entries from the history file.
// Returns an empty History if the file does not exist.
func LoadHistory(snapshotPath string) (History, error) {
	return load(snapshotPath)
}

// Latest returns the most recent snapshot entry, or an error if history is empty.
func Latest(snapshotPath string) (Entry, error) {
	h, err := load(snapshotPath)
	if err != nil {
		return Entry{}, err
	}
	if len(h.Entries) == 0 {
		return Entry{}, errors.New("snapshot: no entries recorded")
	}
	return h.Entries[len(h.Entries)-1], nil
}

func load(snapshotPath string) (History, error) {
	data, err := os.ReadFile(snapshotPath)
	if errors.Is(err, os.ErrNotExist) {
		return History{}, nil
	}
	if err != nil {
		return History{}, fmt.Errorf("read snapshot file: %w", err)
	}
	var h History
	if err := json.Unmarshal(data, &h); err != nil {
		return History{}, fmt.Errorf("parse snapshot file: %w", err)
	}
	return h, nil
}
