// Package diff provides utilities for comparing plaintext .env files
// and producing human-readable change summaries between versions.
package diff

import (
	"fmt"
	"sort"
	"strings"
)

// Change represents a single key-level change between two env snapshots.
type Change struct {
	Key    string
	Type   ChangeType
	OldVal string
	NewVal string
}

// ChangeType describes the kind of diff operation.
type ChangeType string

const (
	Added    ChangeType = "added"
	Removed  ChangeType = "removed"
	Modified ChangeType = "modified"
)

// String returns a human-readable description of the change.
// Values are masked to avoid leaking secrets in logs.
func (c Change) String() string {
	switch c.Type {
	case Added:
		return fmt.Sprintf("+ %s (added)", c.Key)
	case Removed:
		return fmt.Sprintf("- %s (removed)", c.Key)
	case Modified:
		return fmt.Sprintf("~ %s (modified)", c.Key)
	}
	return ""
}

// Compare parses two env file contents and returns the list of changes.
// Both inputs are expected to be in KEY=VALUE format, one per line.
func Compare(oldContent, newContent string) []Change {
	oldMap := parseEnv(oldContent)
	newMap := parseEnv(newContent)

	var changes []Change

	for k, oldVal := range oldMap {
		if newVal, ok := newMap[k]; !ok {
			changes = append(changes, Change{Key: k, Type: Removed, OldVal: oldVal})
		} else if oldVal != newVal {
			changes = append(changes, Change{Key: k, Type: Modified, OldVal: oldVal, NewVal: newVal})
		}
	}

	for k, newVal := range newMap {
		if _, ok := oldMap[k]; !ok {
			changes = append(changes, Change{Key: k, Type: Added, NewVal: newVal})
		}
	}

	sort.Slice(changes, func(i, j int) bool {
		return changes[i].Key < changes[j].Key
	})

	return changes
}

// Summary returns a compact multi-line string summarising all changes.
func Summary(changes []Change) string {
	if len(changes) == 0 {
		return "no changes"
	}
	lines := make([]string, 0, len(changes))
	for _, c := range changes {
		lines = append(lines, c.String())
	}
	return strings.Join(lines, "\n")
}

// parseEnv converts KEY=VALUE text into a map, ignoring blank lines and comments.
func parseEnv(content string) map[string]string {
	m := make(map[string]string)
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			m[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
		}
	}
	return m
}
