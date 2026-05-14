package envtemplate_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/nicholasgasior/envseal/internal/envtemplate"
)

func writeEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0644); err != nil {
		t.Fatalf("writeEnv: %v", err)
	}
	return p
}

func readFile(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readFile: %v", err)
	}
	return string(b)
}

func TestGenerate_StripsValues(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "DB_HOST=localhost\nDB_PORT=5432\nSECRET=supersecret\n")
	dest := filepath.Join(dir, ".env.template")

	if err := envtemplate.Generate(src, dest); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	got := readFile(t, dest)
	for _, line := range []string{"DB_HOST=", "DB_PORT=", "SECRET="} {
		if !strings.Contains(got, line) {
			t.Errorf("expected %q in output, got:\n%s", line, got)
		}
	}
	if strings.Contains(got, "localhost") || strings.Contains(got, "supersecret") {
		t.Errorf("output should not contain original values, got:\n%s", got)
	}
}

func TestGenerate_PreservesComments(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "# Database config\nDB_HOST=localhost\n")
	dest := filepath.Join(dir, ".env.template")

	if err := envtemplate.Generate(src, dest); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	got := readFile(t, dest)
	if !strings.Contains(got, "# Database config") {
		t.Errorf("expected comment preserved, got:\n%s", got)
	}
}

func TestGenerate_PreservesBlankLines(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "A=1\n\nB=2\n")
	dest := filepath.Join(dir, ".env.template")

	if err := envtemplate.Generate(src, dest); err != nil {
		t.Fatalf("Generate: %v", err)
	}

	got := readFile(t, dest)
	lines := strings.Split(strings.TrimRight(got, "\n"), "\n")
	if len(lines) != 3 {
		t.Errorf("expected 3 lines (including blank), got %d:\n%s", len(lines), got)
	}
}

func TestGenerate_MissingSource(t *testing.T) {
	dir := t.TempDir()
	err := envtemplate.Generate(filepath.Join(dir, "nonexistent.env"), filepath.Join(dir, "out"))
	if err == nil {
		t.Error("expected error for missing source file")
	}
}
