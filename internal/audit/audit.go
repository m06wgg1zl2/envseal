// Package audit provides a simple audit log for tracking seal/unseal operations.
package audit

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

const DefaultAuditLogPath = ".envseal_audit.log"

// EventType represents the type of audit event.
type EventType string

const (
	EventSeal   EventType = "seal"
	EventUnseal EventType = "unseal"
	EventRotate EventType = "rotate"
	EventAddKey EventType = "add_key"
)

// Entry represents a single audit log entry.
type Entry struct {
	Timestamp time.Time `json:"timestamp"`
	Event     EventType `json:"event"`
	User      string    `json:"user,omitempty"`
	File      string    `json:"file,omitempty"`
	Version   int       `json:"version,omitempty"`
	Note      string    `json:"note,omitempty"`
}

// Log appends an audit entry to the given log file path.
func Log(logPath string, entry Entry) error {
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now().UTC()
	}

	data, err := json.Marshal(entry)
	if err != nil {
		return fmt.Errorf("audit: marshal entry: %w", err)
	}
	data = append(data, '\n')

	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		return fmt.Errorf("audit: open log file: %w", err)
	}
	defer f.Close()

	if _, err := f.Write(data); err != nil {
		return fmt.Errorf("audit: write entry: %w", err)
	}
	return nil
}

// ReadAll reads all audit entries from the given log file path.
func ReadAll(logPath string) ([]Entry, error) {
	data, err := os.ReadFile(logPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("audit: read log file: %w", err)
	}

	var entries []Entry
	decoder := json.NewDecoder(bytes.NewReader(data))
	for decoder.More() {
		var e Entry
		if err := decoder.Decode(&e); err != nil {
			return nil, fmt.Errorf("audit: decode entry: %w", err)
		}
		entries = append(entries, e)
	}
	return entries, nil
}

// ensure bytesReader-replacement compiles away; keep io.Reader interface satisfied.
var _ io.Reader = (*bytes.Reader)(nil)
