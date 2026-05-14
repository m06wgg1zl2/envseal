package envredact_test

import (
	"testing"

	"github.com/nicholasgasior/envseal/internal/envredact"
)

func TestRedact_SensitiveKeysAreRedacted(t *testing.T) {
	env := map[string]string{
		"DB_PASSWORD": "supersecret",
		"API_KEY":     "abc123",
		"APP_NAME":    "myapp",
		"PORT":        "8080",
	}

	opts := envredact.DefaultOptions()
	got, err := envredact.Redact(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if got["DB_PASSWORD"] != "[REDACTED]" {
		t.Errorf("DB_PASSWORD: want [REDACTED], got %q", got["DB_PASSWORD"])
	}
	if got["API_KEY"] != "[REDACTED]" {
		t.Errorf("API_KEY: want [REDACTED], got %q", got["API_KEY"])
	}
	if got["APP_NAME"] != "myapp" {
		t.Errorf("APP_NAME: want myapp, got %q", got["APP_NAME"])
	}
	if got["PORT"] != "8080" {
		t.Errorf("PORT: want 8080, got %q", got["PORT"])
	}
}

func TestRedact_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"SECRET_KEY": "original"}
	opts := envredact.DefaultOptions()
	_, err := envredact.Redact(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if env["SECRET_KEY"] != "original" {
		t.Error("Redact mutated the input map")
	}
}

func TestRedact_EmptyEnv(t *testing.T) {
	got, err := envredact.Redact(map[string]string{}, envredact.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(got) != 0 {
		t.Errorf("expected empty map, got %v", got)
	}
}

func TestRedact_CustomPattern(t *testing.T) {
	env := map[string]string{
		"STRIPE_KEY": "sk_live_xxx",
		"HOSTNAME":   "localhost",
	}
	opts := envredact.Options{Patterns: []string{"(?i)stripe"}}
	got, err := envredact.Redact(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["STRIPE_KEY"] != "[REDACTED]" {
		t.Errorf("STRIPE_KEY: want [REDACTED], got %q", got["STRIPE_KEY"])
	}
	if got["HOSTNAME"] != "localhost" {
		t.Errorf("HOSTNAME: want localhost, got %q", got["HOSTNAME"])
	}
}

func TestRedact_InvalidPattern(t *testing.T) {
	env := map[string]string{"KEY": "value"}
	opts := envredact.Options{Patterns: []string{"[invalid"}}
	_, err := envredact.Redact(env, opts)
	if err == nil {
		t.Error("expected error for invalid regex pattern, got nil")
	}
}

func TestIsSensitiveKey(t *testing.T) {
	cases := []struct {
		key      string
		want     bool
	}{
		{"DB_PASSWORD", true},
		{"AUTH_TOKEN", true},
		{"APP_ENV", false},
		{"PRIVATE_KEY", true},
		{"DEBUG", false},
	}
	for _, tc := range cases {
		got, err := envredact.IsSensitiveKey(tc.key, envredact.DefaultSensitivePatterns)
		if err != nil {
			t.Fatalf("key %q: unexpected error: %v", tc.key, err)
		}
		if got != tc.want {
			t.Errorf("IsSensitiveKey(%q): want %v, got %v", tc.key, tc.want, got)
		}
	}
}
