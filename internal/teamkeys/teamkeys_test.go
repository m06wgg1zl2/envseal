package teamkeys_test

import (
	"os"
	"path/filepath"
	"testing"

	"filippo.io/age"

	"github.com/yourusername/envseal/internal/teamkeys"
)

func generatePublicKey(t *testing.T) string {
	t.Helper()
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("generating identity: %v", err)
	}
	return identity.Recipient().String()
}

func TestNew(t *testing.T) {
	tk := teamkeys.New()
	if tk.Keys == nil {
		t.Fatal("expected non-nil Keys map")
	}
	if len(tk.Keys) != 0 {
		t.Fatalf("expected empty Keys map, got %d entries", len(tk.Keys))
	}
}

func TestAdd_Valid(t *testing.T) {
	tk := teamkeys.New()
	pub := generatePublicKey(t)
	if err := tk.Add("alice", pub); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tk.Keys["alice"] != pub {
		t.Fatalf("expected key to be stored")
	}
}

func TestAdd_EmptyName(t *testing.T) {
	tk := teamkeys.New()
	pub := generatePublicKey(t)
	if err := tk.Add("", pub); err == nil {
		t.Fatal("expected error for empty name")
	}
}

func TestAdd_InvalidKey(t *testing.T) {
	tk := teamkeys.New()
	if err := tk.Add("bob", "not-a-valid-key"); err == nil {
		t.Fatal("expected error for invalid public key")
	}
}

func TestRemove(t *testing.T) {
	tk := teamkeys.New()
	pub := generatePublicKey(t)
	_ = tk.Add("alice", pub)
	if err := tk.Remove("alice"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := tk.Keys["alice"]; ok {
		t.Fatal("expected key to be removed")
	}
}

func TestRemove_NotFound(t *testing.T) {
	tk := teamkeys.New()
	if err := tk.Remove("nobody"); err == nil {
		t.Fatal("expected error for missing key")
	}
}

func TestRecipients(t *testing.T) {
	tk := teamkeys.New()
	_ = tk.Add("alice", generatePublicKey(t))
	_ = tk.Add("bob", generatePublicKey(t))
	recipients, err := tk.Recipients()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(recipients) != 2 {
		t.Fatalf("expected 2 recipients, got %d", len(recipients))
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, teamkeys.DefaultTeamKeysFile)

	tk := teamkeys.New()
	_ = tk.Add("alice", generatePublicKey(t))

	if err := teamkeys.Save(path, tk); err != nil {
		t.Fatalf("save error: %v", err)
	}

	loaded, err := teamkeys.Load(path)
	if err != nil {
		t.Fatalf("load error: %v", err)
	}
	if loaded.Keys["alice"] != tk.Keys["alice"] {
		t.Fatalf("loaded key mismatch")
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	_, err := teamkeys.Load("/nonexistent/path/.envseal-team.json")
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "bad.json")
	_ = os.WriteFile(path, []byte("not json"), 0o600)
	_, err := teamkeys.Load(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}
