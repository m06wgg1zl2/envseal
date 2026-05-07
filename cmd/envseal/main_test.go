package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/user/envseal/internal/keystore"
	"github.com/user/envseal/internal/teamkeys"
)

func TestRunInit(t *testing.T) {
	dir := t.TempDir()
	keyFile := filepath.Join(dir, "identity.txt")

	// Should not panic or exit
	runInit([]string{"-key", keyFile})

	if _, err := os.Stat(keyFile); err != nil {
		t.Fatalf("identity file not created: %v", err)
	}

	identity, err := keystore.LoadIdentity(keyFile)
	if err != nil {
		t.Fatalf("load identity: %v", err)
	}
	if identity == nil {
		t.Error("expected non-nil identity")
	}
}

func TestRunInit_Idempotent(t *testing.T) {
	dir := t.TempDir()
	keyFile := filepath.Join(dir, "identity.txt")

	// Run init twice; the second call should not overwrite the existing key.
	runInit([]string{"-key", keyFile})
	identity1, err := keystore.LoadIdentity(keyFile)
	if err != nil {
		t.Fatalf("load identity after first init: %v", err)
	}

	runInit([]string{"-key", keyFile})
	identity2, err := keystore.LoadIdentity(keyFile)
	if err != nil {
		t.Fatalf("load identity after second init: %v", err)
	}

	if identity1.Recipient().String() != identity2.Recipient().String() {
		t.Error("second init overwrote existing identity")
	}
}

func TestRunSealUnseal_Integration(t *testing.T) {
	dir := t.TempDir()
	keyFile := filepath.Join(dir, "identity.txt")
	teamFile := filepath.Join(dir, "team.json")
	envFile := filepath.Join(dir, ".env")
	sealedFile := filepath.Join(dir, ".env.sealed")
	outFile := filepath.Join(dir, ".env.out")

	// Init identity
	runInit([]string{"-key", keyFile})

	identity, err := keystore.LoadIdentity(keyFile)
	if err != nil {
		t.Fatalf("load identity: %v", err)
	}

	// Create team with our key
	tm := teamkeys.New()
	if err := tm.Add("me", identity.Recipient().String()); err != nil {
		t.Fatalf("add key: %v", err)
	}
	if err := teamkeys.Save(teamFile, tm); err != nil {
		t.Fatalf("save team: %v", err)
	}

	// Write env file
	envContent := "API_KEY=abc123\nDEBUG=true\n"
	if err := os.WriteFile(envFile, []byte(envContent), 0600); err != nil {
		t.Fatalf("write env: %v", err)
	}

	// Seal
	runSeal([]string{"-env", envFile, "-out", sealedFile, "-team", teamFile, "-key", keyFile})

	if _, err := os.Stat(sealedFile); err != nil {
		t.Fatalf("sealed file missing: %v", err)
	}

	// Unseal
	runUnseal([]string{"-in", sealedFile, "-out", outFile, "-key", keyFile})

	got, err := os.ReadFile(outFile)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if string(got) != envContent {
		t.Errorf("content mismatch: got %q, want %q", got, envContent)
	}
}

func TestRunAddKey(t *testing.T) {
	dir := t.TempDir()
	teamFile := filepath.Join(dir, "team.json")

	identity, err := keystore.GenerateIdentity()
	if err != nil {
		t.Fatalf("generate identity: %v", err)
	}
	pubKey := identity.Recipient().String()

	runAddKey([]string{"-team", teamFile, "-name", "alice", "-pubkey", pubKey})

	tm, err := teamkeys.Load(teamFile)
	if err != nil {
		t.Fatalf("load team: %v", err)
	}
	if _, ok := tm.Keys()["alice"]; !ok {
		t.Error("expected alice in team keys")
	}
}
