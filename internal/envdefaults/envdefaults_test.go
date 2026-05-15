package envdefaults_test

import (
	"testing"

	"github.com/yourusername/envseal/internal/envdefaults"
)

func TestApply_SetsAbsentKeys(t *testing.T) {
	env := map[string]string{"EXISTING": "yes"}
	defaults := []envdefaults.Default{
		{Key: "EXISTING", Value: "no"},
		{Key: "NEW_KEY", Value: "default_value"},
	}
	out, result, err := envdefaults.Apply(env, defaults)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["EXISTING"] != "yes" {
		t.Errorf("expected EXISTING=yes, got %q", out["EXISTING"])
	}
	if out["NEW_KEY"] != "default_value" {
		t.Errorf("expected NEW_KEY=default_value, got %q", out["NEW_KEY"])
	}
	if len(result.Applied) != 1 || result.Applied[0] != "NEW_KEY" {
		t.Errorf("unexpected Applied: %v", result.Applied)
	}
	if len(result.Skipped) != 1 || result.Skipped[0] != "EXISTING" {
		t.Errorf("unexpected Skipped: %v", result.Skipped)
	}
}

func TestApply_OverwriteFlag(t *testing.T) {
	env := map[string]string{"HOST": "localhost"}
	defaults := []envdefaults.Default{
		{Key: "HOST", Value: "0.0.0.0", Overwrite: true},
	}
	out, result, err := envdefaults.Apply(env, defaults)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["HOST"] != "0.0.0.0" {
		t.Errorf("expected HOST=0.0.0.0, got %q", out["HOST"])
	}
	if len(result.Overwritten) != 1 {
		t.Errorf("expected 1 overwritten key, got %v", result.Overwritten)
	}
}

func TestApply_EmptyValueTreatedAsAbsent(t *testing.T) {
	env := map[string]string{"PORT": ""}
	defaults := []envdefaults.Default{
		{Key: "PORT", Value: "8080"},
	}
	out, result, err := envdefaults.Apply(env, defaults)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["PORT"] != "8080" {
		t.Errorf("expected PORT=8080, got %q", out["PORT"])
	}
	if len(result.Applied) != 1 {
		t.Errorf("expected 1 applied key, got %v", result.Applied)
	}
}

func TestApply_NilEnvReturnsError(t *testing.T) {
	_, _, err := envdefaults.Apply(nil, nil)
	if err == nil {
		t.Fatal("expected error for nil env, got nil")
	}
}

func TestApply_EmptyKeyReturnsError(t *testing.T) {
	env := map[string]string{}
	defaults := []envdefaults.Default{{Key: "", Value: "v"}}
	_, _, err := envdefaults.Apply(env, defaults)
	if err == nil {
		t.Fatal("expected error for empty key, got nil")
	}
}

func TestApply_DoesNotMutateInput(t *testing.T) {
	env := map[string]string{"A": "1"}
	defaults := []envdefaults.Default{{Key: "B", Value: "2"}}
	_, _, _ = envdefaults.Apply(env, defaults)
	if _, ok := env["B"]; ok {
		t.Error("Apply mutated the input map")
	}
}

func TestFromMap_BuildsDefaults(t *testing.T) {
	m := map[string]string{"X": "1", "Y": "2"}
	defaults := envdefaults.FromMap(m)
	if len(defaults) != 2 {
		t.Errorf("expected 2 defaults, got %d", len(defaults))
	}
	for _, d := range defaults {
		if d.Key == "" {
			t.Error("FromMap produced entry with empty key")
		}
		if d.Overwrite {
			t.Error("FromMap should not set Overwrite=true")
		}
	}
}
