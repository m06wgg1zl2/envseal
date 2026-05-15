package envresolve_test

import (
	"os"
	"testing"

	"github.com/yourusername/envseal/internal/envresolve"
)

func TestResolve_NoReferences(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	got, err := envresolve.Resolve(env, envresolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["FOO"] != "bar" || got["BAZ"] != "qux" {
		t.Errorf("expected unchanged values, got %v", got)
	}
}

func TestResolve_BraceStyle(t *testing.T) {
	env := map[string]string{
		"BASE_URL": "https://example.com",
		"API_URL":  "${BASE_URL}/api",
	}
	got, err := envresolve.Resolve(env, envresolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "https://example.com/api"
	if got["API_URL"] != want {
		t.Errorf("API_URL = %q, want %q", got["API_URL"], want)
	}
}

func TestResolve_DollarStyle(t *testing.T) {
	env := map[string]string{
		"HOST": "localhost",
		"DSN":  "postgres://$HOST/db",
	}
	got, err := envresolve.Resolve(env, envresolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := "postgres://localhost/db"
	if got["DSN"] != want {
		t.Errorf("DSN = %q, want %q", got["DSN"], want)
	}
}

func TestResolve_FallbackToOS(t *testing.T) {
	os.Setenv("_ENVSEAL_TEST_OS_VAR", "from-os")
	defer os.Unsetenv("_ENVSEAL_TEST_OS_VAR")

	env := map[string]string{"GREETING": "hello ${_ENVSEAL_TEST_OS_VAR}"}
	opts := envresolve.DefaultOptions()
	opts.FallbackToOS = true

	got, err := envresolve.Resolve(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["GREETING"] != "hello from-os" {
		t.Errorf("GREETING = %q, want %q", got["GREETING"], "hello from-os")
	}
}

func TestResolve_MissingVar_NoFail(t *testing.T) {
	env := map[string]string{"VAL": "${UNDEFINED_XYZ}"}
	opts := envresolve.DefaultOptions()
	opts.FallbackToOS = false
	opts.FailOnMissing = false

	got, err := envresolve.Resolve(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["VAL"] != "" {
		t.Errorf("VAL = %q, want empty string", got["VAL"])
	}
}

func TestResolve_MissingVar_FailOnMissing(t *testing.T) {
	env := map[string]string{"VAL": "${UNDEFINED_XYZ}"}
	opts := envresolve.DefaultOptions()
	opts.FallbackToOS = false
	opts.FailOnMissing = true

	_, err := envresolve.Resolve(env, opts)
	if err == nil {
		t.Fatal("expected error for missing variable, got nil")
	}
}

func TestResolve_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"A": "1", "B": "${A}-2"}
	original := map[string]string{"A": "1", "B": "${A}-2"}

	_, err := envresolve.Resolve(env, envresolve.DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	for k, v := range original {
		if env[k] != v {
			t.Errorf("input mutated: env[%q] = %q, want %q", k, env[k], v)
		}
	}
}
