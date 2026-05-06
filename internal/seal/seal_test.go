package seal_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envseal/internal/keystore"
	"github.com/user/envseal/internal/seal"
	"github.com/user/envseal/internal/teamkeys"
)

func setupTestDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	return dir
}

func TestSealUnseal_RoundTrip(t *testing.T) {
	dir := setupTestDir(t)

	// Generate identity
	identity, err := keystore.GenerateIdentity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}
	keyFile := filepath.Join(dir, "identity.txt")
	if err := keystore.SaveIdentity(keyFile, identity); err != nil {
		t.Fatalf("save identity: %v", err)
	}

	// Create team keys with our public key
	tm := teamkeys.New()
	if err := tm.Add("testuser", identity.Recipient().String()); err != nil {
		t.Fatalf("add team key: %v", err)
	}
	teamFile := filepath.Join(dir, "team.json")
	if err := teamkeys.Save(teamFile, tm); err != nil {
		t.Fatalf("save team: %v", err)
	}

	// Write plaintext env file
	envContent := []byte("DB_HOST=localhost\nDB_PORT=5432\nSECRET=supersecret\n")
	envFile := filepath.Join(dir, ".env")
	if err := os.WriteFile(envFile, envContent, 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	sealedFile := filepath.Join(dir, ".env.sealed")
	outFile := filepath.Join(dir, ".env.out")

	// Seal
	err = seal.Seal(seal.SealOptions{
		EnvFile:    envFile,
		OutputFile: sealedFile,
		TeamFile:   teamFile,
		KeyFile:    keyFile,
	})
	if err != nil {
		t.Fatalf("seal: %v", err)
	}

	// Verify sealed file exists
	if _, err := os.Stat(sealedFile); err != nil {
		t.Fatalf("sealed file not found: %v", err)
	}

	// Unseal
	err = seal.Unseal(seal.UnsealOptions{
		SealedFile: sealedFile,
		OutputFile: outFile,
		KeyFile:    keyFile,
	})
	if err != nil {
		t.Fatalf("unseal: %v", err)
	}

	got, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}

	if string(got) != string(envContent) {
		t.Errorf("content mismatch: got %q, want %q", got, envContent)
	}
}

func TestSeal_MissingEnvFile(t *testing.T) {
	dir := setupTestDir(t)
	tm := teamkeys.New()
	teamFile := filepath.Join(dir, "team.json")
	_ = json.NewEncoder(func() *os.File { f, _ := os.Create(teamFile); return f }()).Encode(tm)

	err := seal.Seal(seal.SealOptions{
		EnvFile:  filepath.Join(dir, "nonexistent.env"),
		TeamFile: teamFile,
	})
	if err == nil {
		t.Error("expected error for missing env file")
	}
}

func TestUnseal_MissingSealedFile(t *testing.T) {
	dir := setupTestDir(t)
	err := seal.Unseal(seal.UnsealOptions{
		SealedFile: filepath.Join(dir, "nonexistent.sealed"),
	})
	if err == nil {
		t.Error("expected error for missing sealed file")
	}
}
