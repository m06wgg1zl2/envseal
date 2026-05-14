package envcheck_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envseal/internal/envcheck"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeFile: %v", err)
	}
	return p
}

func TestCheck_NoIssues(t *testing.T) {
	dir := t.TempDir()
	env := writeFile(t, dir, ".env", "FOO=bar\nBAZ=qux\n")
	tpl := writeFile(t, dir, ".env.template", "FOO=\nBAZ=\n")

	res, err := envcheck.Check(env, tpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HasIssues() {
		t.Errorf("expected no issues, got:\n%s", res.Format())
	}
}

func TestCheck_MissingKey(t *testing.T) {
	dir := t.TempDir()
	env := writeFile(t, dir, ".env", "FOO=bar\n")
	tpl := writeFile(t, dir, ".env.template", "FOO=\nBAZ=\n")

	res, err := envcheck.Check(env, tpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.MissingKeys) != 1 || res.MissingKeys[0] != "BAZ" {
		t.Errorf("expected MissingKeys=[BAZ], got %v", res.MissingKeys)
	}
	if len(res.ExtraKeys) != 0 {
		t.Errorf("expected no extra keys, got %v", res.ExtraKeys)
	}
}

func TestCheck_ExtraKey(t *testing.T) {
	dir := t.TempDir()
	env := writeFile(t, dir, ".env", "FOO=bar\nSECRET=x\n")
	tpl := writeFile(t, dir, ".env.template", "FOO=\n")

	res, err := envcheck.Check(env, tpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(res.ExtraKeys) != 1 || res.ExtraKeys[0] != "SECRET" {
		t.Errorf("expected ExtraKeys=[SECRET], got %v", res.ExtraKeys)
	}
}

func TestCheck_IgnoresComments(t *testing.T) {
	dir := t.TempDir()
	env := writeFile(t, dir, ".env", "# comment\nFOO=bar\n")
	tpl := writeFile(t, dir, ".env.template", "# another comment\nFOO=\n")

	res, err := envcheck.Check(env, tpl)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.HasIssues() {
		t.Errorf("comments should be ignored, got issues:\n%s", res.Format())
	}
}

func TestCheck_MissingEnvFile(t *testing.T) {
	dir := t.TempDir()
	tpl := writeFile(t, dir, ".env.template", "FOO=\n")

	_, err := envcheck.Check(filepath.Join(dir, "nonexistent.env"), tpl)
	if err == nil {
		t.Error("expected error for missing env file")
	}
}

func TestCheck_MissingTemplateFile(t *testing.T) {
	dir := t.TempDir()
	env := writeFile(t, dir, ".env", "FOO=bar\n")

	_, err := envcheck.Check(env, filepath.Join(dir, "nonexistent.template"))
	if err == nil {
		t.Error("expected error for missing template file")
	}
}
