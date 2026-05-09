package policy_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envseal/internal/policy"
)

func TestCheck_NoViolations(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"APP_PORT": "8080",
	}
	rule := policy.Rule{
		Required: []string{"APP_NAME", "APP_PORT"},
	}
	v := policy.Check(env, rule)
	if len(v) != 0 {
		t.Fatalf("expected no violations, got %d", len(v))
	}
}

func TestCheck_MissingRequired(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
	}
	rule := policy.Rule{
		Required: []string{"APP_NAME", "DATABASE_URL"},
	}
	v := policy.Check(env, rule)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if v[0].Key != "DATABASE_URL" {
		t.Errorf("expected violation for DATABASE_URL, got %q", v[0].Key)
	}
}

func TestCheck_ForbiddenPresent(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"DEBUG":    "true",
	}
	rule := policy.Rule{
		Forbidden: []string{"DEBUG"},
	}
	v := policy.Check(env, rule)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(v))
	}
	if !strings.Contains(v[0].Message, "forbidden") {
		t.Errorf("expected forbidden message, got %q", v[0].Message)
	}
}

func TestCheck_PatternEnforced(t *testing.T) {
	env := map[string]string{
		"APP_NAME": "myapp",
		"lowercase": "bad",
	}
	rule := policy.Rule{
		Pattern: `^[A-Z][A-Z0-9_]+$`,
	}
	v := policy.Check(env, rule)
	if len(v) != 1 {
		t.Fatalf("expected 1 pattern violation, got %d", len(v))
	}
	if v[0].Key != "lowercase" {
		t.Errorf("expected violation for 'lowercase', got %q", v[0].Key)
	}
}

func TestCheck_InvalidPattern(t *testing.T) {
	env := map[string]string{"KEY": "val"}
	rule := policy.Rule{Pattern: `[invalid`}
	v := policy.Check(env, rule)
	if len(v) != 1 || v[0].Key != "__pattern__" {
		t.Fatalf("expected pattern compile error violation, got %+v", v)
	}
}

func TestHasViolations(t *testing.T) {
	if policy.HasViolations(nil) {
		t.Error("expected false for nil violations")
	}
	v := []policy.Violation{{Key: "X", Message: "missing"}}
	if !policy.HasViolations(v) {
		t.Error("expected true for non-empty violations")
	}
}

func TestFormat_NoViolations(t *testing.T) {
	out := policy.Format(nil)
	if !strings.Contains(out, "no policy violations") {
		t.Errorf("unexpected format output: %q", out)
	}
}

func TestFormat_WithViolations(t *testing.T) {
	v := []policy.Violation{
		{Key: "DB_URL", Message: "required key is missing"},
	}
	out := policy.Format(v)
	if !strings.Contains(out, "1 policy violation") {
		t.Errorf("expected count in output, got: %q", out)
	}
	if !strings.Contains(out, "DB_URL") {
		t.Errorf("expected key name in output, got: %q", out)
	}
}
