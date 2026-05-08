package lint_test

import (
	"testing"

	"github.com/yourorg/envseal/internal/lint"
)

func TestCheck_ValidEnv(t *testing.T) {
	content := "DB_HOST=localhost\nDB_PORT=5432\nAPI_KEY=secret\n"
	issues := lint.Check(content)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d: %v", len(issues), issues)
	}
}

func TestCheck_EmptyValue(t *testing.T) {
	content := "DB_HOST=\n"
	issues := lint.Check(content)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "warning" {
		t.Errorf("expected warning, got %s", issues[0].Severity)
	}
	if issues[0].Key != "DB_HOST" {
		t.Errorf("expected key DB_HOST, got %s", issues[0].Key)
	}
}

func TestCheck_InvalidLine(t *testing.T) {
	content := "NOTAVALIDLINE\n"
	issues := lint.Check(content)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "error" {
		t.Errorf("expected error severity, got %s", issues[0].Severity)
	}
}

func TestCheck_DuplicateKey(t *testing.T) {
	content := "FOO=bar\nFOO=baz\n"
	issues := lint.Check(content)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue for duplicate, got %d", len(issues))
	}
	if issues[0].Severity != "warning" {
		t.Errorf("expected warning for duplicate, got %s", issues[0].Severity)
	}
}

func TestCheck_IgnoresComments(t *testing.T) {
	content := "# This is a comment\n\nDB=val\n"
	issues := lint.Check(content)
	if len(issues) != 0 {
		t.Errorf("expected no issues, got %d", len(issues))
	}
}

func TestCheck_EmptyKey(t *testing.T) {
	content := "=value\n"
	issues := lint.Check(content)
	if len(issues) != 1 {
		t.Fatalf("expected 1 issue, got %d", len(issues))
	}
	if issues[0].Severity != "error" {
		t.Errorf("expected error for empty key, got %s", issues[0].Severity)
	}
}

func TestHasErrors_True(t *testing.T) {
	issues := []lint.Issue{{Severity: "error"}, {Severity: "warning"}}
	if !lint.HasErrors(issues) {
		t.Error("expected HasErrors to return true")
	}
}

func TestHasErrors_False(t *testing.T) {
	issues := []lint.Issue{{Severity: "warning"}}
	if lint.HasErrors(issues) {
		t.Error("expected HasErrors to return false")
	}
}

func TestIssue_String(t *testing.T) {
	i := lint.Issue{Line: 3, Key: "FOO", Message: "value is empty", Severity: "warning"}
	s := i.String()
	if s == "" {
		t.Error("expected non-empty string from Issue.String()")
	}
}
