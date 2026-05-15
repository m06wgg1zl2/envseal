package envclone_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/your-org/envseal/internal/envclone"
)

func writeEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeEnv: %v", err)
	}
	return p
}

func readEnv(t *testing.T, path string) string {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("readEnv: %v", err)
	}
	return string(b)
}

func TestClone_CopiesAllKeys(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "FOO=bar\nBAZ=qux\n")
	dst := filepath.Join(dir, ".env.clone")

	if err := envclone.Clone(src, dst, envclone.Options{}); err != nil {
		t.Fatalf("Clone: %v", err)
	}
	out := readEnv(t, dst)
	if !strings.Contains(out, "FOO=bar") || !strings.Contains(out, "BAZ=qux") {
		t.Errorf("expected all keys in output, got:\n%s", out)
	}
}

func TestClone_IncludeKeys(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "FOO=bar\nBAZ=qux\nSECRET=hidden\n")
	dst := filepath.Join(dir, ".env.out")

	opts := envclone.Options{IncludeKeys: []string{"FOO", "BAZ"}}
	if err := envclone.Clone(src, dst, opts); err != nil {
		t.Fatalf("Clone: %v", err)
	}
	out := readEnv(t, dst)
	if strings.Contains(out, "SECRET") {
		t.Errorf("SECRET should have been excluded, got:\n%s", out)
	}
	if !strings.Contains(out, "FOO=bar") {
		t.Errorf("FOO should be present, got:\n%s", out)
	}
}

func TestClone_ExcludeKeys(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "FOO=bar\nSECRET=hidden\n")
	dst := filepath.Join(dir, ".env.out")

	opts := envclone.Options{ExcludeKeys: []string{"SECRET"}}
	if err := envclone.Clone(src, dst, opts); err != nil {
		t.Fatalf("Clone: %v", err)
	}
	out := readEnv(t, dst)
	if strings.Contains(out, "SECRET") {
		t.Errorf("SECRET should be excluded, got:\n%s", out)
	}
}

func TestClone_RenameKeys(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "OLD_KEY=value\n")
	dst := filepath.Join(dir, ".env.out")

	opts := envclone.Options{Rename: map[string]string{"OLD_KEY": "NEW_KEY"}}
	if err := envclone.Clone(src, dst, opts); err != nil {
		t.Fatalf("Clone: %v", err)
	}
	out := readEnv(t, dst)
	if !strings.Contains(out, "NEW_KEY=value") {
		t.Errorf("expected NEW_KEY=value, got:\n%s", out)
	}
	if strings.Contains(out, "OLD_KEY") {
		t.Errorf("OLD_KEY should be renamed, got:\n%s", out)
	}
}

func TestClone_NoOverwriteByDefault(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "FOO=bar\n")
	dst := writeEnv(t, dir, ".env.out", "existing\n")

	err := envclone.Clone(src, dst, envclone.Options{})
	if err == nil {
		t.Fatal("expected error when dst exists and Overwrite is false")
	}
}

func TestClone_OverwriteFlag(t *testing.T) {
	dir := t.TempDir()
	src := writeEnv(t, dir, ".env", "FOO=newval\n")
	dst := writeEnv(t, dir, ".env.out", "FOO=oldval\n")

	if err := envclone.Clone(src, dst, envclone.Options{Overwrite: true}); err != nil {
		t.Fatalf("Clone with overwrite: %v", err)
	}
	out := readEnv(t, dst)
	if !strings.Contains(out, "FOO=newval") {
		t.Errorf("expected overwritten value, got:\n%s", out)
	}
}

func TestClone_EmptyPaths(t *testing.T) {
	if err := envclone.Clone("", "dst", envclone.Options{}); err == nil {
		t.Error("expected error for empty src")
	}
	if err := envclone.Clone("src", "", envclone.Options{}); err == nil {
		t.Error("expected error for empty dst")
	}
}
