package envmerge_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourusername/envseal/internal/envmerge"
)

func writeEnv(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0600); err != nil {
		t.Fatalf("writeEnv: %v", err)
	}
	return p
}

func TestMerge_BasicPrecedence(t *testing.T) {
	dir := t.TempDir()
	a := writeEnv(t, dir, "a.env", "FOO=from_a\nBAR=bar_a\n")
	b := writeEnv(t, dir, "b.env", "FOO=from_b\nBAZ=baz_b\n")

	r, err := envmerge.Merge([]string{a, b})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if r.Pairs["FOO"] != "from_b" {
		t.Errorf("FOO: want from_b, got %s", r.Pairs["FOO"])
	}
	if r.Pairs["BAR"] != "bar_a" {
		t.Errorf("BAR: want bar_a, got %s", r.Pairs["BAR"])
	}
	if r.Pairs["BAZ"] != "baz_b" {
		t.Errorf("BAZ: want baz_b, got %s", r.Pairs["BAZ"])
	}
	if r.Sources["FOO"] != b {
		t.Errorf("FOO source: want %s, got %s", b, r.Sources["FOO"])
	}
}

func TestMerge_IgnoresCommentsAndBlanks(t *testing.T) {
	dir := t.TempDir()
	p := writeEnv(t, dir, "c.env", "# comment\n\nKEY=value\n")

	r, err := envmerge.Merge([]string{p})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(r.Pairs) != 1 {
		t.Errorf("expected 1 key, got %d", len(r.Pairs))
	}
	if r.Pairs["KEY"] != "value" {
		t.Errorf("KEY: want value, got %s", r.Pairs["KEY"])
	}
}

func TestMerge_NoFiles(t *testing.T) {
	_, err := envmerge.Merge([]string{})
	if err == nil {
		t.Fatal("expected error for empty file list")
	}
}

func TestMerge_MissingFile(t *testing.T) {
	_, err := envmerge.Merge([]string{"/nonexistent/path.env"})
	if err == nil {
		t.Fatal("expected error for missing file")
	}
}

func TestWrite_ProducesValidEnv(t *testing.T) {
	dir := t.TempDir()
	a := writeEnv(t, dir, "a.env", "ALPHA=1\nBETA=2\n")
	b := writeEnv(t, dir, "b.env", "BETA=overridden\nGAMMA=3\n")

	r, err := envmerge.Merge([]string{a, b})
	if err != nil {
		t.Fatalf("merge: %v", err)
	}

	out := filepath.Join(dir, "merged.env")
	if err := envmerge.Write(r, out); err != nil {
		t.Fatalf("write: %v", err)
	}

	data, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("read: %v", err)
	}
	content := string(data)
	for _, want := range []string{"ALPHA=1", "BETA=overridden", "GAMMA=3"} {
		if !contains(content, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(s) > 0 && containsStr(s, sub))
}

func containsStr(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
