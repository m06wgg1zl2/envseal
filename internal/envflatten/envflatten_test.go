package envflatten_test

import (
	"testing"

	"github.com/nicholasgasior/envseal/internal/envflatten"
)

func TestFlatten_DotsToSeparator(t *testing.T) {
	env := map[string]string{
		"db.host": "localhost",
		"db.port": "5432",
	}
	opts := envflatten.DefaultOptions()
	got, err := envflatten.Flatten(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB__HOST"] != "localhost" {
		t.Errorf("expected DB__HOST=localhost, got %q", got["DB__HOST"])
	}
	if got["DB__PORT"] != "5432" {
		t.Errorf("expected DB__PORT=5432, got %q", got["DB__PORT"])
	}
}

func TestFlatten_UppercaseKeys(t *testing.T) {
	env := map[string]string{"app_name": "envseal"}
	opts := envflatten.DefaultOptions()
	got, err := envflatten.Flatten(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := got["APP_NAME"]; !ok {
		t.Errorf("expected key APP_NAME to exist")
	}
}

func TestFlatten_NoUppercase(t *testing.T) {
	env := map[string]string{"App_Name": "envseal"}
	opts := envflatten.Options{Separator: "__", UppercaseKeys: false}
	got, err := envflatten.Flatten(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["App_Name"] != "envseal" {
		t.Errorf("expected key to remain unchanged, got %v", got)
	}
}

func TestFlatten_NilEnvReturnsError(t *testing.T) {
	_, err := envflatten.Flatten(nil, envflatten.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestFlatten_EmptyKeyReturnsError(t *testing.T) {
	env := map[string]string{"": "value"}
	_, err := envflatten.Flatten(env, envflatten.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty key")
	}
}

func TestFlatten_DuplicateAfterNormalisationReturnsError(t *testing.T) {
	env := map[string]string{
		"db.host": "localhost",
		"db__host": "remotehost",
	}
	_, err := envflatten.Flatten(env, envflatten.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for duplicate normalised key")
	}
}

func TestPrefix_AddsPrefix(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	opts := envflatten.DefaultOptions()
	got, err := envflatten.Prefix(env, "DB", opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got["DB__HOST"] != "localhost" {
		t.Errorf("expected DB__HOST=localhost, got %q", got["DB__HOST"])
	}
	if got["DB__PORT"] != "5432" {
		t.Errorf("expected DB__PORT=5432, got %q", got["DB__PORT"])
	}
}

func TestPrefix_EmptyPrefixReturnsError(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	_, err := envflatten.Prefix(env, "", envflatten.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for empty prefix")
	}
}

func TestPrefix_NilEnvReturnsError(t *testing.T) {
	_, err := envflatten.Prefix(nil, "APP", envflatten.DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}
