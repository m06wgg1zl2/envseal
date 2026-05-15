package envprofile_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/envseal/internal/envprofile"
)

func tmpPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "profiles.json")
}

func TestLoad_EmptyIfNotExist(t *testing.T) {
	s, err := envprofile.Load("/nonexistent/path/profiles.json")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(s.Profiles) != 0 {
		t.Errorf("expected empty profiles, got %d", len(s.Profiles))
	}
}

func TestSet_AddsProfile(t *testing.T) {
	s, _ := envprofile.Load("/nonexistent")
	err := envprofile.Set(s, "dev", map[string]string{"DEBUG": "true"})
	if err != nil {
		t.Fatalf("Set failed: %v", err)
	}
	p, ok := s.Profiles["dev"]
	if !ok {
		t.Fatal("expected profile 'dev' to exist")
	}
	if p.Vars["DEBUG"] != "true" {
		t.Errorf("expected DEBUG=true, got %s", p.Vars["DEBUG"])
	}
}

func TestSet_EmptyNameReturnsError(t *testing.T) {
	s := &envprofile.Store{Profiles: make(map[string]envprofile.Profile)}
	if err := envprofile.Set(s, "", map[string]string{}); err == nil {
		t.Error("expected error for empty name")
	}
}

func TestSwitch_SetsActive(t *testing.T) {
	s, _ := envprofile.Load("/nonexistent")
	_ = envprofile.Set(s, "staging", map[string]string{"ENV": "staging"})
	if err := envprofile.Switch(s, "staging"); err != nil {
		t.Fatalf("Switch failed: %v", err)
	}
	if s.Active != "staging" {
		t.Errorf("expected active=staging, got %s", s.Active)
	}
}

func TestSwitch_UnknownProfileReturnsError(t *testing.T) {
	s := &envprofile.Store{Profiles: make(map[string]envprofile.Profile)}
	if err := envprofile.Switch(s, "unknown"); err == nil {
		t.Error("expected error for unknown profile")
	}
}

func TestSaveAndLoad_RoundTrip(t *testing.T) {
	path := tmpPath(t)
	s, _ := envprofile.Load(path)
	_ = envprofile.Set(s, "prod", map[string]string{"LOG_LEVEL": "warn"})
	_ = envprofile.Switch(s, "prod")
	if err := envprofile.Save(path, s); err != nil {
		t.Fatalf("Save failed: %v", err)
	}
	s2, err := envprofile.Load(path)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if s2.Active != "prod" {
		t.Errorf("expected active=prod, got %s", s2.Active)
	}
	if s2.Profiles["prod"].Vars["LOG_LEVEL"] != "warn" {
		t.Errorf("expected LOG_LEVEL=warn")
	}
}

func TestActiveVars_ReturnsVarsForActiveProfile(t *testing.T) {
	s, _ := envprofile.Load("/nonexistent")
	_ = envprofile.Set(s, "dev", map[string]string{"A": "1", "B": "2"})
	_ = envprofile.Switch(s, "dev")
	vars := envprofile.ActiveVars(s)
	if vars["A"] != "1" || vars["B"] != "2" {
		t.Errorf("unexpected vars: %v", vars)
	}
}

func TestActiveVars_EmptyWhenNoActive(t *testing.T) {
	s := &envprofile.Store{Profiles: make(map[string]envprofile.Profile)}
	if len(envprofile.ActiveVars(s)) != 0 {
		t.Error("expected empty vars when no active profile")
	}
}

func TestLoad_InvalidJSON(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	_ = os.WriteFile(path, []byte("{invalid"), 0o600)
	if _, err := envprofile.Load(path); err == nil {
		t.Error("expected error for invalid JSON")
	}
}
