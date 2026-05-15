package main

import (
	"os"
	"path/filepath"
	"testing"
)

func profileFile(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "profiles.json")
}

func TestRunProfile_Set(t *testing.T) {
	f := profileFile(t)
	err := runProfile([]string{"-file", f, "set", "dev", "DEBUG=true", "LOG=info"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, err := os.Stat(f); err != nil {
		t.Errorf("expected profiles file to exist: %v", err)
	}
}

func TestRunProfile_Switch(t *testing.T) {
	f := profileFile(t)
	_ = runProfile([]string{"-file", f, "set", "staging", "ENV=staging"})
	err := runProfile([]string{"-file", f, "switch", "staging"})
	if err != nil {
		t.Fatalf("Switch failed: %v", err)
	}
}

func TestRunProfile_Switch_UnknownProfile(t *testing.T) {
	f := profileFile(t)
	err := runProfile([]string{"-file", f, "switch", "ghost"})
	if err == nil {
		t.Error("expected error switching to unknown profile")
	}
}

func TestRunProfile_List(t *testing.T) {
	f := profileFile(t)
	_ = runProfile([]string{"-file", f, "set", "dev", "A=1"})
	_ = runProfile([]string{"-file", f, "set", "prod", "A=0"})
	if err := runProfile([]string{"-file", f, "list"}); err != nil {
		t.Fatalf("list failed: %v", err)
	}
}

func TestRunProfile_Active_NoProfile(t *testing.T) {
	f := profileFile(t)
	if err := runProfile([]string{"-file", f, "active"}); err != nil {
		t.Fatalf("active failed: %v", err)
	}
}

func TestRunProfile_Active_WithProfile(t *testing.T) {
	f := profileFile(t)
	_ = runProfile([]string{"-file", f, "set", "dev", "X=1"})
	_ = runProfile([]string{"-file", f, "switch", "dev"})
	if err := runProfile([]string{"-file", f, "active"}); err != nil {
		t.Fatalf("active failed: %v", err)
	}
}

func TestRunProfile_Set_InvalidKV(t *testing.T) {
	f := profileFile(t)
	err := runProfile([]string{"-file", f, "set", "dev", "NODEQUALS"})
	if err == nil {
		t.Error("expected error for invalid key=value")
	}
}

func TestRunProfile_UnknownSubcommand(t *testing.T) {
	f := profileFile(t)
	if err := runProfile([]string{"-file", f, "frobnicate"}); err == nil {
		t.Error("expected error for unknown subcommand")
	}
}
