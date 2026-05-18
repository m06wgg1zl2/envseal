package envtransform_test

import (
	"testing"

	"github.com/yourorg/envseal/internal/envtransform"
)

func TestTransform_NoOp(t *testing.T) {
	env := map[string]string{"FOO": "bar", "BAZ": "qux"}
	out, err := envtransform.Transform(env, envtransform.Options{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 || out["FOO"] != "bar" || out["BAZ"] != "qux" {
		t.Errorf("expected no-op transform, got %v", out)
	}
}

func TestTransform_AddPrefix(t *testing.T) {
	env := map[string]string{"HOST": "localhost", "PORT": "5432"}
	out, err := envtransform.Transform(env, envtransform.Options{AddPrefix: "DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["DB_HOST"] != "localhost" || out["DB_PORT"] != "5432" {
		t.Errorf("unexpected output: %v", out)
	}
	if _, ok := out["HOST"]; ok {
		t.Error("original key should not be present")
	}
}

func TestTransform_StripPrefix(t *testing.T) {
	env := map[string]string{"APP_NAME": "envseal", "APP_ENV": "prod"}
	out, err := envtransform.Transform(env, envtransform.Options{StripPrefix: "APP_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NAME"] != "envseal" || out["ENV"] != "prod" {
		t.Errorf("unexpected output: %v", out)
	}
}

func TestTransform_RenameKeys(t *testing.T) {
	env := map[string]string{"OLD_KEY": "value"}
	out, err := envtransform.Transform(env, envtransform.Options{
		RenameKeys: map[string]string{"OLD_KEY": "NEW_KEY"},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["NEW_KEY"] != "value" {
		t.Errorf("expected NEW_KEY=value, got %v", out)
	}
	if _, ok := out["OLD_KEY"]; ok {
		t.Error("old key should not be present after rename")
	}
}

func TestTransform_FilterPattern(t *testing.T) {
	env := map[string]string{"DB_HOST": "localhost", "APP_NAME": "envseal", "DB_PORT": "5432"}
	out, err := envtransform.Transform(env, envtransform.Options{FilterPattern: "^DB_"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 2 {
		t.Errorf("expected 2 keys, got %d: %v", len(out), out)
	}
	if _, ok := out["APP_NAME"]; ok {
		t.Error("APP_NAME should have been filtered out")
	}
}

func TestTransform_InvalidFilterPattern(t *testing.T) {
	env := map[string]string{"FOO": "bar"}
	_, err := envtransform.Transform(env, envtransform.Options{FilterPattern: "[invalid"})
	if err == nil {
		t.Error("expected error for invalid regex, got nil")
	}
}

func TestTransform_NilEnvReturnsError(t *testing.T) {
	_, err := envtransform.Transform(nil, envtransform.Options{})
	if err == nil {
		t.Error("expected error for nil env, got nil")
	}
}

func TestTransform_StripToEmptyKeyDropped(t *testing.T) {
	env := map[string]string{"PREFIX": "value"}
	out, err := envtransform.Transform(env, envtransform.Options{StripPrefix: "PREFIX"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(out) != 0 {
		t.Errorf("expected empty output when key becomes empty, got %v", out)
	}
}
