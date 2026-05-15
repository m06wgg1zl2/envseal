// Package envnormalize provides utilities for normalizing .env file contents,
// including trimming whitespace, uppercasing keys, and removing duplicate entries.
package envnormalize

import (
	"bufio"
	"fmt"
	"strings"
)

// Options controls normalization behaviour.
type Options struct {
	// UppercaseKeys converts all keys to UPPER_CASE.
	UppercaseKeys bool
	// TrimValues strips leading/trailing whitespace from values.
	TrimValues bool
	// DeduplicateKeys keeps only the last occurrence of a duplicate key.
	DeduplicateKeys bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		UppercaseKeys:   true,
		TrimValues:      true,
		DeduplicateKeys: true,
	}
}

// Normalize applies the given options to env, returning a new map with
// normalized key/value pairs. Comments and blank lines are not represented
// in the output map; use Write to produce a normalized file.
func Normalize(env map[string]string, opts Options) (map[string]string, error) {
	if env == nil {
		return nil, fmt.Errorf("envnormalize: env map must not be nil")
	}
	out := make(map[string]string, len(env))
	for k, v := range env {
		key := k
		if opts.UppercaseKeys {
			key = strings.ToUpper(k)
		}
		val := v
		if opts.TrimValues {
			val = strings.TrimSpace(v)
		}
		out[key] = val
	}
	return out, nil
}

// NormalizeLines reads KEY=VALUE lines from src, applies opts, and returns
// normalized lines. Blank lines and comments are preserved in place.
// When DeduplicateKeys is true, only the last assignment for a key is kept.
func NormalizeLines(src string, opts Options) ([]string, error) {
	type entry struct {
		line    string
		key     string
		isKV    bool
	}

	var entries []entry
	scanner := bufio.NewScanner(strings.NewReader(src))
	for scanner.Scan() {
		raw := scanner.Text()
		trimmed := strings.TrimSpace(raw)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			entries = append(entries, entry{line: raw})
			continue
		}
		idx := strings.IndexByte(trimmed, '=')
		if idx <= 0 {
			entries = append(entries, entry{line: raw})
			continue
		}
		k := trimmed[:idx]
		v := trimmed[idx+1:]
		if opts.UppercaseKeys {
			k = strings.ToUpper(k)
		}
		if opts.TrimValues {
			v = strings.TrimSpace(v)
		}
		normLine := k + "=" + v
		entries = append(entries, entry{line: normLine, key: k, isKV: true})
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("envnormalize: scan error: %w", err)
	}

	if opts.DeduplicateKeys {
		// Mark earlier duplicate KV entries as removed.
		lastIdx := make(map[string]int)
		for i, e := range entries {
			if e.isKV {
				lastIdx[e.key] = i
			}
		}
		var deduped []string
		for i, e := range entries {
			if e.isKV && lastIdx[e.key] != i {
				continue
			}
			deduped = append(deduped, e.line)
		}
		return deduped, nil
	}

	out := make([]string, len(entries))
	for i, e := range entries {
		out[i] = e.line
	}
	return out, nil
}
