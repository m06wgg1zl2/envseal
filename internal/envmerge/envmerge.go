// Package envmerge provides utilities for merging multiple .env files,
// with later files taking precedence over earlier ones.
package envmerge

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Result holds the merged key-value pairs and metadata about the merge.
type Result struct {
	// Pairs is the final merged map of key to value.
	Pairs map[string]string
	// Order preserves insertion order of keys for deterministic output.
	Order []string
	// Sources maps each key to the file it was last defined in.
	Sources map[string]string
}

// Merge reads the given env files in order and merges them into a single
// Result. Later files override keys from earlier files.
func Merge(paths []string) (*Result, error) {
	if len(paths) == 0 {
		return nil, fmt.Errorf("envmerge: no files provided")
	}

	r := &Result{
		Pairs:   make(map[string]string),
		Order:   []string{},
		Sources: make(map[string]string),
	}

	for _, path := range paths {
		pairs, err := parseFile(path)
		if err != nil {
			return nil, fmt.Errorf("envmerge: reading %s: %w", path, err)
		}
		for _, kv := range pairs {
			if _, exists := r.Pairs[kv[0]]; !exists {
				r.Order = append(r.Order, kv[0])
			}
			r.Pairs[kv[0]] = kv[1]
			r.Sources[kv[0]] = path
		}
	}

	return r, nil
}

// Write serialises the merged result as a .env file to the given path.
func Write(r *Result, dest string) error {
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("envmerge: creating %s: %w", dest, err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)
	for _, k := range r.Order {
		fmt.Fprintf(w, "%s=%s\n", k, r.Pairs[k])
	}
	return w.Flush()
}

// parseFile reads a .env file and returns ordered key-value pairs.
func parseFile(path string) ([][2]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var pairs [][2]string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 1 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		val := strings.TrimSpace(line[idx+1:])
		pairs = append(pairs, [2]string{key, val})
	}
	return pairs, scanner.Err()
}
