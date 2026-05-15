package envwatch_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envseal/internal/envwatch"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestChecksum_Deterministic(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, ".env", "KEY=value\n")

	a, err := envwatch.Checksum(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	b, err := envwatch.Checksum(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if a != b {
		t.Errorf("expected identical checksums, got %s vs %s", a, b)
	}
}

func TestChecksum_DiffersOnChange(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, ".env", "KEY=value\n")

	before, _ := envwatch.Checksum(p)
	if err := os.WriteFile(p, []byte("KEY=changed\n"), 0600); err != nil {
		t.Fatal(err)
	}
	after, _ := envwatch.Checksum(p)

	if before == after {
		t.Error("expected checksums to differ after file change")
	}
}

func TestWatch_Unchanged(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, ".env", "KEY=value\n")

	baseline, err := envwatch.Checksum(p)
	if err != nil {
		t.Fatal(err)
	}

	status, err := envwatch.Watch(p, baseline)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if status.Changed {
		t.Error("expected file to be unchanged")
	}
}

func TestWatch_Changed(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, ".env", "KEY=value\n")

	baseline, _ := envwatch.Checksum(p)
	_ = os.WriteFile(p, []byte("KEY=new\n"), 0600)

	status, err := envwatch.Watch(p, baseline)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !status.Changed {
		t.Error("expected file to be marked as changed")
	}
}

func TestWatch_NoBaseline(t *testing.T) {
	dir := t.TempDir()
	p := writeFile(t, dir, ".env", "KEY=value\n")

	_, err := envwatch.Watch(p, "")
	if err != envwatch.ErrNoBaseline {
		t.Errorf("expected ErrNoBaseline, got %v", err)
	}
}

func TestWatch_MissingFile(t *testing.T) {
	_, err := envwatch.Watch("/nonexistent/.env", "abc123")
	if err == nil {
		t.Error("expected error for missing file")
	}
}
