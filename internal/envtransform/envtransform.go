// Package envtransform provides utilities for applying key/value transformations
// to environment variable maps, such as prefixing keys, renaming keys, and
// filtering by pattern.
package envtransform

import (
	"errors"
	"regexp"
	"strings"
)

// Options controls how transformations are applied.
type Options struct {
	// AddPrefix prepends a string to every key.
	AddPrefix string
	// StripPrefix removes a leading prefix from every key (applied after AddPrefix).
	StripPrefix string
	// RenameKeys maps old key names to new key names.
	RenameKeys map[string]string
	// FilterPattern, if non-empty, keeps only keys matching the pattern.
	FilterPattern string
}

// Transform applies the given Options to env and returns a new map.
// The original map is never mutated.
func Transform(env map[string]string, opts Options) (map[string]string, error) {
	if env == nil {
		return nil, errors.New("envtransform: env must not be nil")
	}

	var filter *regexp.Regexp
	if opts.FilterPattern != "" {
		var err error
		filter, err = regexp.Compile(opts.FilterPattern)
		if err != nil {
			return nil, fmt.Errorf("envtransform: invalid filter pattern: %w", err)
		}
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if filter != nil && !filter.MatchString(k) {
			continue
		}
		newKey := k
		if opts.AddPrefix != "" {
			newKey = opts.AddPrefix + newKey
		}
		if opts.StripPrefix != "" {
			newKey = strings.TrimPrefix(newKey, opts.StripPrefix)
		}
		if renamed, ok := opts.RenameKeys[newKey]; ok {
			newKey = renamed
		}
		if newKey == "" {
			continue
		}
		out[newKey] = v
	}
	return out, nil
}
