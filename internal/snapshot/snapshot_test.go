package snapshot_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourusername/envseal/internal/snapshot"
)

func TestRecord_CreatesFileAndAppends(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envseal_snapshots.json")

	entry1 := snapshot.Entry{Version: 1, Checksum: "abc123", SealedAt: time.Now().UTC()}
	if err := snapshot.Record(path, entry1); err != nil {
		t.Fatalf("Record: %v", err)
	}

	entry2 := snapshot.Entry{Version: 2, Checksum: "def456", SealedAt: time.Now().UTC()}
	if err := snapshot.Record(path, entry2); err != nil {
		t.Fatalf("Record second: %v", err)
	}

	h, err := snapshot.LoadHistory(path)
	if err != nil {
		t.Fatalf("LoadHistory: %v", err)
	}
	if len(h.Entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(h.Entries))
	}
	if h.Entries[0].Checksum != "abc123" {
		t.Errorf("expected abc123, got %s", h.Entries[0].Checksum)
	}
	if h.Entries[1].Version != 2 {
		t.Errorf("expected version 2, got %d", h.Entries[1].Version)
	}
}

func TestRecord_SetsTimestampIfZero(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envseal_snapshots.json")

	entry := snapshot.Entry{Version: 1, Checksum: "xyz"}
	if err := snapshot.Record(path, entry); err != nil {
		t.Fatalf("Record: %v", err)
	}

	h, _ := snapshot.LoadHistory(path)
	if h.Entries[0].SealedAt.IsZero() {
		t.Error("expected SealedAt to be set automatically")
	}
}

func TestRecord_InvalidEntry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envseal_snapshots.json")

	if err := snapshot.Record(path, snapshot.Entry{Version: 0, Checksum: "abc"}); err == nil {
		t.Error("expected error for version=0")
	}
	if err := snapshot.Record(path, snapshot.Entry{Version: 1, Checksum: ""}); err == nil {
		t.Error("expected error for empty checksum")
	}
}

func TestLoadHistory_EmptyIfFileNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.json")
	h, err := snapshot.LoadHistory(path)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(h.Entries) != 0 {
		t.Errorf("expected empty history, got %d entries", len(h.Entries))
	}
}

func TestLoadHistory_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not-json"), 0o644)

	_, err := snapshot.LoadHistory(path)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLatest_ReturnsLastEntry(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, ".envseal_snapshots.json")

	_ = snapshot.Record(path, snapshot.Entry{Version: 1, Checksum: "first"})
	_ = snapshot.Record(path, snapshot.Entry{Version: 2, Checksum: "second"})

	latest, err := snapshot.Latest(path)
	if err != nil {
		t.Fatalf("Latest: %v", err)
	}
	if latest.Checksum != "second" {
		t.Errorf("expected 'second', got %s", latest.Checksum)
	}
}

func TestLatest_ErrorIfNoEntries(t *testing.T) {
	path := filepath.Join(t.TempDir(), "empty.json")
	_, err := snapshot.Latest(path)
	if err == nil {
		t.Error("expected error when no entries exist")
	}
}

func TestDefaultSnapshotPath(t *testing.T) {
	p := snapshot.DefaultSnapshotPath("/some/dir/.env.sealed")
	expected := "/some/dir/.envseal_snapshots.json"
	if p != expected {
		t.Errorf("expected %s, got %s", expected, p)
	}
}
