// Package envdiff provides utilities for computing structured diffs
// between two sets of environment variables and rendering the result
// in human-readable or machine-readable formats.
package envdiff

import (
	"fmt"
	"sort"
	"strings"
)

// ChangeKind describes the type of change for a single key.
type ChangeKind string

const (
	Added    ChangeKind = "added"
	Removed  ChangeKind = "removed"
	Modified ChangeKind = "modified"
	Unchanged ChangeKind = "unchanged"
)

// Change represents a single key-level difference.
type Change struct {
	Key      string
	Kind     ChangeKind
	OldValue string
	NewValue string
}

// Result holds all changes between two env maps.
type Result struct {
	Changes []Change
}

// HasChanges reports whether any meaningful changes exist.
func (r Result) HasChanges() bool {
	for _, c := range r.Changes {
		if c.Kind != Unchanged {
			return true
		}
	}
	return false
}

// Compare computes the diff between oldEnv and newEnv.
// Only keys that differ are included unless includeUnchanged is true.
func Compare(oldEnv, newEnv map[string]string, includeUnchanged bool) Result {
	seen := make(map[string]bool)
	var changes []Change

	for k, oldVal := range oldEnv {
		seen[k] = true
		newVal, exists := newEnv[k]
		switch {
		case !exists:
			changes = append(changes, Change{Key: k, Kind: Removed, OldValue: oldVal})
		case oldVal != newVal:
			changes = append(changes, Change{Key: k, Kind: Modified, OldValue: oldVal, NewValue: newVal})
		default:
			if includeUnchanged {
				changes = append(changes, Change{Key: k, Kind: Unchanged, OldValue: oldVal, NewValue: newVal})
			}
		}
	}

	for k, newVal := range newEnv {
		if !seen[k] {
			changes = append(changes, Change{Key: k, Kind: Added, NewValue: newVal})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return Result{Changes: changes}
}

// Format renders the diff result as a human-readable string.
func Format(r Result) string {
	var sb strings.Builder
	for _, c := range r.Changes {
		switch c.Kind {
		case Added:
			fmt.Fprintf(&sb, "+ %s=%s\n", c.Key, c.NewValue)
		case Removed:
			fmt.Fprintf(&sb, "- %s=%s\n", c.Key, c.OldValue)
		case Modified:
			fmt.Fprintf(&sb, "~ %s: %s -> %s\n", c.Key, c.OldValue, c.NewValue)
		case Unchanged:
			fmt.Fprintf(&sb, "  %s=%s\n", c.Key, c.NewValue)
		}
	}
	return sb.String()
}
