package envelope_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/envseal/internal/envelope"
)

func TestNew(t *testing.T) {
	recipients := []string{"age1abc123", "age1def456"}
	ciphertext := "base64encodeddata=="

	env := envelope.New(recipients, ciphertext)

	if env.Version != 1 {
		t.Errorf("expected version 1, got %d", env.Version)
	}
	if env.Ciphertext != ciphertext {
		t.Errorf("expected ciphertext %q, got %q", ciphertext, env.Ciphertext)
	}
	if len(env.Recipients) != 2 {
		t.Errorf("expected 2 recipients, got %d", len(env.Recipients))
	}
	if env.CreatedAt.After(time.Now().Add(time.Second)) {
		t.Error("created_at should not be in the future")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "test.env.sealed")

	original := envelope.New([]string{"age1xyz789"}, "encryptedpayload==")

	if err := envelope.Save(path, original); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	loaded, err := envelope.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if loaded.Version != original.Version {
		t.Errorf("version mismatch: got %d, want %d", loaded.Version, original.Version)
	}
	if loaded.Ciphertext != original.Ciphertext {
		t.Errorf("ciphertext mismatch: got %q, want %q", loaded.Ciphertext, original.Ciphertext)
	}
	if len(loaded.Recipients) != len(original.Recipients) {
		t.Errorf("recipients length mismatch: got %d, want %d", len(loaded.Recipients), len(original.Recipients))
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := envelope.Load("/nonexistent/path/file.sealed")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.sealed")

	if err := os.WriteFile(path, []byte("not json"), 0644); err != nil {
		t.Fatalf("setup failed: %v", err)
	}

	_, err := envelope.Load(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}
