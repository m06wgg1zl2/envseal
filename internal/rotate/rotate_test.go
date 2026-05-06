package rotate_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envseal/internal/keystore"
	"github.com/user/envseal/internal/rotate"
	"github.com/user/envseal/internal/seal"
	"github.com/user/envseal/internal/teamkeys"
)

func setupRotateDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func TestRotate_RoundTrip(t *testing.T) {
	dir := setupRotateDir(t)

	// Generate an identity for sealing and rotation.
	identityPath := filepath.Join(dir, "key.age")
	identity, err := keystore.GenerateIdentity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}
	if err := keystore.SaveIdentity(identity, identityPath); err != nil {
		t.Fatalf("save identity: %v", err)
	}

	// Create a team keys file with the identity's public key.
	tk := teamkeys.New()
	if err := tk.Add("alice", identity.Recipient().String()); err != nil {
		t.Fatalf("add team key: %v", err)
	}
	teamKeysPath := filepath.Join(dir, "team.json")
	if err := tk.Save(teamKeysPath); err != nil {
		t.Fatalf("save team keys: %v", err)
	}

	// Write a sample .env file and seal it.
	envPath := filepath.Join(dir, ".env")
	if err := os.WriteFile(envPath, []byte("SECRET=hello\nTOKEN=world\n"), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}
	sealedPath := filepath.Join(dir, ".env.sealed.json")
	if err := seal.Seal(envPath, sealedPath, teamKeysPath, identityPath); err != nil {
		t.Fatalf("seal: %v", err)
	}

	// Rotate with the same team keys (no-op key change, but exercises the path).
	err = rotate.Rotate(rotate.Options{
		SealedFile:   sealedPath,
		TeamKeysFile: teamKeysPath,
		IdentityFile: identityPath,
	})
	if err != nil {
		t.Fatalf("rotate: %v", err)
	}

	// Unseal after rotation and verify plaintext.
	outPath := filepath.Join(dir, ".env.out")
	if err := seal.Unseal(sealedPath, outPath, identityPath); err != nil {
		t.Fatalf("unseal after rotate: %v", err)
	}
	got, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if string(got) != "SECRET=hello\nTOKEN=world\n" {
		t.Errorf("unexpected plaintext after rotate: %q", string(got))
	}
}

func TestRotate_MissingSealedFile(t *testing.T) {
	dir := setupRotateDir(t)
	err := rotate.Rotate(rotate.Options{
		SealedFile:   filepath.Join(dir, "nonexistent.json"),
		TeamKeysFile: filepath.Join(dir, "team.json"),
		IdentityFile: filepath.Join(dir, "key.age"),
	})
	if err == nil {
		t.Fatal("expected error for missing sealed file, got nil")
	}
}

func TestRotate_EmptySealedFilePath(t *testing.T) {
	err := rotate.Rotate(rotate.Options{})
	if err == nil {
		t.Fatal("expected error for empty sealed file path, got nil")
	}
}
