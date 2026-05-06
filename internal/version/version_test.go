package version_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nicholasgasior/envseal/internal/version"
)

func TestNew_SetsFieldsCorrectly(t *testing.T) {
	plaintext := []byte("FOO=bar\nBAZ=qux")
	info := version.New(plaintext, "alice")

	if info.Version != 1 {
		t.Errorf("expected version 1, got %d", info.Version)
	}
	if info.SealedBy != "alice" {
		t.Errorf("expected SealedBy 'alice', got %q", info.SealedBy)
	}
	if info.Checksum == "" {
		t.Error("expected non-empty checksum")
	}
	if info.SealedAt.IsZero() {
		t.Error("expected non-zero SealedAt")
	}
}

func TestBump_IncrementsVersion(t *testing.T) {
	plaintext := []byte("FOO=bar")
	prev := version.New(plaintext, "bob")

	newPlaintext := []byte("FOO=bar\nNEW=value")
	bumped := version.Bump(prev, newPlaintext)

	if bumped.Version != prev.Version+1 {
		t.Errorf("expected version %d, got %d", prev.Version+1, bumped.Version)
	}
	if bumped.Checksum == prev.Checksum {
		t.Error("expected checksum to change after content update")
	}
	if bumped.SealedBy != prev.SealedBy {
		t.Errorf("expected SealedBy to be preserved: got %q", bumped.SealedBy)
	}
}

func TestChecksum_Deterministic(t *testing.T) {
	data := []byte("hello world")
	a := version.Checksum(data)
	b := version.Checksum(data)
	if a != b {
		t.Errorf("checksum not deterministic: %q != %q", a, b)
	}
}

func TestVerifyFile_Valid(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.env")
	content := []byte("KEY=value")
	if err := os.WriteFile(path, content, 0600); err != nil {
		t.Fatal(err)
	}
	expected := version.Checksum(content)
	if err := version.VerifyFile(path, expected); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestVerifyFile_Mismatch(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.env")
	if err := os.WriteFile(path, []byte("KEY=value"), 0600); err != nil {
		t.Fatal(err)
	}
	err := version.VerifyFile(path, "deadbeef")
	if err == nil {
		t.Error("expected checksum mismatch error, got nil")
	}
}

func TestHasChanged(t *testing.T) {
	original := []byte("FOO=1")
	info := version.New(original, "")

	if version.HasChanged(info, original) {
		t.Error("expected no change for same content")
	}
	if !version.HasChanged(info, []byte("FOO=2")) {
		t.Error("expected change detected for different content")
	}
}
