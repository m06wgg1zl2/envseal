package envdiff

import (
	"strings"
	"testing"
)

func TestCompare_NoChanges(t *testing.T) {
	old := map[string]string{"FOO": "bar", "BAZ": "qux"}
	new := map[string]string{"FOO": "bar", "BAZ": "qux"}
	r := Compare(old, new, false)
	if r.HasChanges() {
		t.Error("expected no changes")
	}
	if len(r.Changes) != 0 {
		t.Errorf("expected 0 changes, got %d", len(r.Changes))
	}
}

func TestCompare_Added(t *testing.T) {
	old := map[string]string{}
	new := map[string]string{"NEW_KEY": "value"}
	r := Compare(old, new, false)
	if !r.HasChanges() {
		t.Fatal("expected changes")
	}
	if len(r.Changes) != 1 || r.Changes[0].Kind != Added {
		t.Errorf("expected one Added change, got %+v", r.Changes)
	}
	if r.Changes[0].Key != "NEW_KEY" || r.Changes[0].NewValue != "value" {
		t.Errorf("unexpected change: %+v", r.Changes[0])
	}
}

func TestCompare_Removed(t *testing.T) {
	old := map[string]string{"OLD_KEY": "gone"}
	new := map[string]string{}
	r := Compare(old, new, false)
	if len(r.Changes) != 1 || r.Changes[0].Kind != Removed {
		t.Errorf("expected one Removed change, got %+v", r.Changes)
	}
	if r.Changes[0].OldValue != "gone" {
		t.Errorf("unexpected OldValue: %s", r.Changes[0].OldValue)
	}
}

func TestCompare_Modified(t *testing.T) {
	old := map[string]string{"KEY": "old"}
	new := map[string]string{"KEY": "new"}
	r := Compare(old, new, false)
	if len(r.Changes) != 1 || r.Changes[0].Kind != Modified {
		t.Errorf("expected one Modified change, got %+v", r.Changes)
	}
	c := r.Changes[0]
	if c.OldValue != "old" || c.NewValue != "new" {
		t.Errorf("unexpected values: old=%s new=%s", c.OldValue, c.NewValue)
	}
}

func TestCompare_IncludeUnchanged(t *testing.T) {
	old := map[string]string{"A": "1", "B": "2"}
	new := map[string]string{"A": "1", "B": "2"}
	r := Compare(old, new, true)
	if len(r.Changes) != 2 {
		t.Errorf("expected 2 unchanged entries, got %d", len(r.Changes))
	}
	for _, c := range r.Changes {
		if c.Kind != Unchanged {
			t.Errorf("expected Unchanged, got %s", c.Kind)
		}
	}
}

func TestCompare_SortedByKey(t *testing.T) {
	old := map[string]string{"Z": "1", "A": "2"}
	new := map[string]string{}
	r := Compare(old, new, false)
	if len(r.Changes) != 2 {
		t.Fatalf("expected 2 changes, got %d", len(r.Changes))
	}
	if r.Changes[0].Key != "A" || r.Changes[1].Key != "Z" {
		t.Errorf("expected sorted keys, got %s, %s", r.Changes[0].Key, r.Changes[1].Key)
	}
}

func TestFormat_ContainsSymbols(t *testing.T) {
	r := Result{Changes: []Change{
		{Key: "A", Kind: Added, NewValue: "1"},
		{Key: "B", Kind: Removed, OldValue: "2"},
		{Key: "C", Kind: Modified, OldValue: "old", NewValue: "new"},
	}}
	out := Format(r)
	if !strings.Contains(out, "+ A=1") {
		t.Errorf("missing added line in:\n%s", out)
	}
	if !strings.Contains(out, "- B=2") {
		t.Errorf("missing removed line in:\n%s", out)
	}
	if !strings.Contains(out, "~ C: old -> new") {
		t.Errorf("missing modified line in:\n%s", out)
	}
}
