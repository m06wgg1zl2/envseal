package envnormalize

import (
	"strings"
	"testing"
)

func TestNormalize_UppercaseKeys(t *testing.T) {
	env := map[string]string{"db_host": "localhost", "api_key": "secret"}
	out, err := Normalize(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["DB_HOST"]; !ok {
		t.Error("expected DB_HOST to be present")
	}
	if _, ok := out["API_KEY"]; !ok {
		t.Error("expected API_KEY to be present")
	}
}

func TestNormalize_TrimValues(t *testing.T) {
	env := map[string]string{"KEY": "  value  "}
	out, err := Normalize(env, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if out["KEY"] != "value" {
		t.Errorf("expected trimmed value, got %q", out["KEY"])
	}
}

func TestNormalize_NilEnvReturnsError(t *testing.T) {
	_, err := Normalize(nil, DefaultOptions())
	if err == nil {
		t.Fatal("expected error for nil env")
	}
}

func TestNormalize_NoUppercase(t *testing.T) {
	opts := DefaultOptions()
	opts.UppercaseKeys = false
	env := map[string]string{"my_key": "val"}
	out, err := Normalize(env, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if _, ok := out["my_key"]; !ok {
		t.Error("expected original key casing to be preserved")
	}
}

func TestNormalizeLines_PreservesComments(t *testing.T) {
	src := "# comment\nKEY=value\n"
	lines, err := NormalizeLines(src, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lines[0] != "# comment" {
		t.Errorf("expected comment preserved, got %q", lines[0])
	}
}

func TestNormalizeLines_DeduplicatesKeys(t *testing.T) {
	src := "KEY=first\nKEY=second\n"
	lines, err := NormalizeLines(src, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 1 {
		t.Fatalf("expected 1 line after dedup, got %d", len(lines))
	}
	if lines[0] != "KEY=second" {
		t.Errorf("expected last value to win, got %q", lines[0])
	}
}

func TestNormalizeLines_UppercasesKeys(t *testing.T) {
	src := "db_host=localhost\n"
	lines, err := NormalizeLines(src, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.HasPrefix(lines[0], "DB_HOST=") {
		t.Errorf("expected DB_HOST= prefix, got %q", lines[0])
	}
}

func TestNormalizeLines_TrimsValues(t *testing.T) {
	src := "KEY=  spaced  \n"
	lines, err := NormalizeLines(src, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if lines[0] != "KEY=spaced" {
		t.Errorf("expected trimmed value, got %q", lines[0])
	}
}

func TestNormalizeLines_PreservesBlankLines(t *testing.T) {
	src := "KEY=val\n\nOTHER=x\n"
	lines, err := NormalizeLines(src, DefaultOptions())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(lines) != 3 {
		t.Errorf("expected 3 lines, got %d", len(lines))
	}
}
