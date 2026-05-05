package keystore_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envseal/internal/keystore"
)

func TestGenerateIdentity(t *testing.T) {
	id, err := keystore.GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity() error: %v", err)
	}
	if id.PublicKey == "" {
		t.Error("expected non-empty public key")
	}
	if id.Private == nil {
		t.Error("expected non-nil private identity")
	}
	if id.Recipient == nil {
		t.Error("expected non-nil recipient")
	}
}

func TestSaveAndLoadIdentity(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "identity.age")

	original, err := keystore.GenerateIdentity()
	if err != nil {
		t.Fatalf("GenerateIdentity() error: %v", err)
	}

	if err := keystore.SaveIdentity(original, path); err != nil {
		t.Fatalf("SaveIdentity() error: %v", err)
	}

	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat key file: %v", err)
	}
	if info.Mode().Perm() != 0600 {
		t.Errorf("expected file mode 0600, got %v", info.Mode().Perm())
	}

	loaded, err := keystore.LoadIdentity(path)
	if err != nil {
		t.Fatalf("LoadIdentity() error: %v", err)
	}

	if loaded.PublicKey != original.PublicKey {
		t.Errorf("public key mismatch: got %q, want %q", loaded.PublicKey, original.PublicKey)
	}
}

func TestSaveIdentity_FailsIfExists(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "identity.age")

	id, _ := keystore.GenerateIdentity()
	_ = keystore.SaveIdentity(id, path)

	err := keystore.SaveIdentity(id, path)
	if err == nil {
		t.Error("expected error when saving to existing file, got nil")
	}
}
