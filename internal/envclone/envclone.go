// Package envclone provides utilities for cloning and forking .env files
// between environments, with optional key filtering and renaming.
package envclone

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Options controls the behaviour of a Clone operation.
type Options struct {
	// Only copy keys present in this list. Empty means copy all.
	IncludeKeys []string
	// Keys to exclude from the clone.
	ExcludeKeys []string
	// Rename maps old key names to new key names.
	Rename map[string]string
	// Overwrite allows the destination file to be replaced if it already exists.
	Overwrite bool
}

// Clone reads src, applies opts, and writes the result to dst.
func Clone(src, dst string, opts Options) error {
	if src == "" || dst == "" {
		return fmt.Errorf("envclone: src and dst paths must not be empty")
	}

	entries, err := parseFile(src)
	if err != nil {
		return fmt.Errorf("envclone: read source: %w", err)
	}

	if !opts.Overwrite {
		if _, err := os.Stat(dst); err == nil {
			return fmt.Errorf("envclone: destination already exists: %s", dst)
		}
	}

	include := toSet(opts.IncludeKeys)
	exclude := toSet(opts.ExcludeKeys)

	f, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("envclone: create destination: %w", err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, line := range entries {
		if line.comment || line.blank {
			fmt.Fprintln(w, line.raw)
			continue
		}
		key := line.key
		if len(include) > 0 && !include[key] {
			continue
		}
		if exclude[key] {
			continue
		}
		if newKey, ok := opts.Rename[key]; ok {
			key = newKey
		}
		fmt.Fprintf(w, "%s=%s\n", key, line.value)
	}
	return w.Flush()
}

type entry struct {
	raw     string
	key     string
	value   string
	comment bool
	blank   bool
}

func parseFile(path string) ([]entry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var entries []entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			entries = append(entries, entry{raw: line, blank: true})
			continue
		}
		if strings.HasPrefix(trimmed, "#") {
			entries = append(entries, entry{raw: line, comment: true})
			continue
		}
		idx := strings.IndexByte(trimmed, '=')
		if idx < 0 {
			entries = append(entries, entry{raw: line, comment: true})
			continue
		}
		entries = append(entries, entry{
			raw:   line,
			key:   trimmed[:idx],
			value: trimmed[idx+1:],
		})
	}
	return entries, scanner.Err()
}

func toSet(keys []string) map[string]bool {
	s := make(map[string]bool, len(keys))
	for _, k := range keys {
		s[k] = true
	}
	return s
}
