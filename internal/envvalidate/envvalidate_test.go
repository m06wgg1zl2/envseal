package envvalidate_test

import (
	"testing"

	"github.com/yourorg/envseal/internal/envvalidate"
)

func TestCheck_NoViolations(t *testing.T) {
	env := map[string]string{"PORT": "8080", "DEBUG": "true"}
	rules := []envvalidate.Rule{
		{Key: "PORT", Required: true, Type: "int"},
		{Key: "DEBUG", Required: true, Type: "bool"},
	}
	violations := envvalidate.Check(env, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestCheck_RequiredMissing(t *testing.T) {
	env := map[string]string{}
	rules := []envvalidate.Rule{
		{Key: "API_KEY", Required: true},
	}
	violations := envvalidate.Check(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
	if violations[0].Key != "API_KEY" {
		t.Errorf("expected violation for API_KEY, got %s", violations[0].Key)
	}
}

func TestCheck_InvalidInt(t *testing.T) {
	env := map[string]string{"PORT": "not-a-number"}
	rules := []envvalidate.Rule{{Key: "PORT", Type: "int"}}
	violations := envvalidate.Check(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheck_InvalidBool(t *testing.T) {
	env := map[string]string{"DEBUG": "yes"}
	rules := []envvalidate.Rule{{Key: "DEBUG", Type: "bool"}}
	violations := envvalidate.Check(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheck_InvalidFloat(t *testing.T) {
	env := map[string]string{"RATIO": "abc"}
	rules := []envvalidate.Rule{{Key: "RATIO", Type: "float"}}
	violations := envvalidate.Check(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheck_MinLen(t *testing.T) {
	env := map[string]string{"TOKEN": "abc"}
	rules := []envvalidate.Rule{{Key: "TOKEN", MinLen: 10}}
	violations := envvalidate.Check(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheck_MaxLen(t *testing.T) {
	env := map[string]string{"SHORT": "toolongvalue"}
	rules := []envvalidate.Rule{{Key: "SHORT", MaxLen: 5}}
	violations := envvalidate.Check(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheck_PatternMatch(t *testing.T) {
	env := map[string]string{"EMAIL": "user@example.com"}
	rules := []envvalidate.Rule{{Key: "EMAIL", Pattern: `^[\w.]+@[\w.]+$`}}
	violations := envvalidate.Check(env, rules)
	if len(violations) != 0 {
		t.Fatalf("expected no violations, got %v", violations)
	}
}

func TestCheck_PatternNoMatch(t *testing.T) {
	env := map[string]string{"EMAIL": "not-an-email"}
	rules := []envvalidate.Rule{{Key: "EMAIL", Pattern: `^[\w.]+@[\w.]+$`}}
	violations := envvalidate.Check(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d", len(violations))
	}
}

func TestCheck_InvalidPattern(t *testing.T) {
	env := map[string]string{"X": "value"}
	rules := []envvalidate.Rule{{Key: "X", Pattern: `[invalid`}}
	violations := envvalidate.Check(env, rules)
	if len(violations) != 1 {
		t.Fatalf("expected 1 violation for invalid pattern, got %d", len(violations))
	}
}

func TestHasViolations(t *testing.T) {
	if envvalidate.HasViolations(nil) {
		t.Error("expected false for nil violations")
	}
	v := []envvalidate.Violation{{Key: "X", Message: "bad"}}
	if !envvalidate.HasViolations(v) {
		t.Error("expected true for non-empty violations")
	}
}
