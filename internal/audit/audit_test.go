package audit_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/envseal/internal/audit"
)

func TestLog_CreatesFileAndAppends(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	entry1 := audit.Entry{
		Event:   audit.EventSeal,
		File:    ".env.sealed",
		Version: 1,
		Note:    "initial seal",
	}
	if err := audit.Log(logPath, entry1); err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	entry2 := audit.Entry{
		Event:   audit.EventUnseal,
		File:    ".env.sealed",
		Version: 1,
	}
	if err := audit.Log(logPath, entry2); err != nil {
		t.Fatalf("Log() second entry error: %v", err)
	}

	entries, err := audit.ReadAll(logPath)
	if err != nil {
		t.Fatalf("ReadAll() error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].Event != audit.EventSeal {
		t.Errorf("expected first event %q, got %q", audit.EventSeal, entries[0].Event)
	}
	if entries[1].Event != audit.EventUnseal {
		t.Errorf("expected second event %q, got %q", audit.EventUnseal, entries[1].Event)
	}
}

func TestLog_SetsTimestampIfZero(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	before := time.Now().UTC()
	entry := audit.Entry{Event: audit.EventRotate}
	if err := audit.Log(logPath, entry); err != nil {
		t.Fatalf("Log() error: %v", err)
	}
	after := time.Now().UTC()

	entries, err := audit.ReadAll(logPath)
	if err != nil {
		t.Fatalf("ReadAll() error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	ts := entries[0].Timestamp
	if ts.Before(before) || ts.After(after) {
		t.Errorf("timestamp %v not in expected range [%v, %v]", ts, before, after)
	}
}

func TestReadAll_EmptyIfFileNotExist(t *testing.T) {
	entries, err := audit.ReadAll("/nonexistent/path/audit.log")
	if err != nil {
		t.Fatalf("ReadAll() unexpected error: %v", err)
	}
	if entries != nil {
		t.Errorf("expected nil entries for missing file, got %v", entries)
	}
}

func TestLog_FilePermissions(t *testing.T) {
	dir := t.TempDir()
	logPath := filepath.Join(dir, "audit.log")

	entry := audit.Entry{Event: audit.EventAddKey, User: "alice"}
	if err := audit.Log(logPath, entry); err != nil {
		t.Fatalf("Log() error: %v", err)
	}

	info, err := os.Stat(logPath)
	if err != nil {
		t.Fatalf("Stat() error: %v", err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Errorf("expected file permissions 0600, got %04o", perm)
	}
}
