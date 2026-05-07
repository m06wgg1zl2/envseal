package diff_test

import (
	"strings"
	"testing"

	"github.com/yourorg/envseal/internal/diff"
)

func TestCompare_NoChanges(t *testing.T) {
	content := "FOO=bar\nBAZ=qux\n"
	changes := diff.Compare(content, content)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestCompare_Added(t *testing.T) {
	old := "FOO=bar\n"
	new := "FOO=bar\nNEW_KEY=secret\n"
	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != diff.Added || changes[0].Key != "NEW_KEY" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompare_Removed(t *testing.T) {
	old := "FOO=bar\nOLD_KEY=gone\n"
	new := "FOO=bar\n"
	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != diff.Removed || changes[0].Key != "OLD_KEY" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompare_Modified(t *testing.T) {
	old := "FOO=bar\n"
	new := "FOO=newval\n"
	changes := diff.Compare(old, new)
	if len(changes) != 1 {
		t.Fatalf("expected 1 change, got %d", len(changes))
	}
	if changes[0].Type != diff.Modified || changes[0].Key != "FOO" {
		t.Errorf("unexpected change: %+v", changes[0])
	}
}

func TestCompare_IgnoresComments(t *testing.T) {
	old := "# comment\nFOO=bar\n"
	new := "FOO=bar\n"
	changes := diff.Compare(old, new)
	if len(changes) != 0 {
		t.Fatalf("expected 0 changes, got %d", len(changes))
	}
}

func TestCompare_SortedOutput(t *testing.T) {
	old := ""
	new := "ZZZ=1\nAAA=2\nMMM=3\n"
	changes := diff.Compare(old, new)
	if changes[0].Key != "AAA" || changes[1].Key != "MMM" || changes[2].Key != "ZZZ" {
		t.Errorf("changes not sorted: %v", changes)
	}
}

func TestSummary_NoChanges(t *testing.T) {
	s := diff.Summary(nil)
	if s != "no changes" {
		t.Errorf("expected 'no changes', got %q", s)
	}
}

func TestSummary_ContainsKeys(t *testing.T) {
	old := "FOO=bar\n"
	new := "FOO=baz\nNEW=val\n"
	changes := diff.Compare(old, new)
	s := diff.Summary(changes)
	if !strings.Contains(s, "FOO") || !strings.Contains(s, "NEW") {
		t.Errorf("summary missing expected keys: %s", s)
	}
}

func TestChange_String_MasksValues(t *testing.T) {
	c := diff.Change{Key: "SECRET", Type: diff.Modified, OldVal: "oldpass", NewVal: "newpass"}
	s := c.String()
	if strings.Contains(s, "oldpass") || strings.Contains(s, "newpass") {
		t.Errorf("Change.String() should not expose values, got: %s", s)
	}
}
