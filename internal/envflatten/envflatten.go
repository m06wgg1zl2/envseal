// Package envflatten provides utilities for flattening nested key structures
// in environment variable maps. It handles dot-notation and double-underscore
// delimited keys, converting them into a flat key=value representation.
package envflatten

import (
	"errors"
	"fmt"
	"strings"
)

// Options controls how flattening is performed.
type Options struct {
	// Separator is the delimiter used to join nested key segments.
	// Defaults to "__" (double underscore).
	Separator string
	// UppercaseKeys converts all resulting keys to uppercase.
	UppercaseKeys bool
}

// DefaultOptions returns a sensible default configuration.
func DefaultOptions() Options {
	return Options{
		Separator:     "__",
		UppercaseKeys: true,
	}
}

// Flatten takes a nested map (map keys may contain dots or the configured
// separator) and returns a flat map with joined keys. Input values that are
// empty strings are preserved as-is.
func Flatten(env map[string]string, opts Options) (map[string]string, error) {
	if env == nil {
		return nil, errors.New("envflatten: env map must not be nil")
	}
	if opts.Separator == "" {
		return nil, errors.New("envflatten: separator must not be empty")
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		if k == "" {
			return nil, errors.New("envflatten: key must not be empty")
		}
		// Normalise dot-notation to the configured separator.
		normKey := strings.ReplaceAll(k, ".", opts.Separator)
		if opts.UppercaseKeys {
			normKey = strings.ToUpper(normKey)
		}
		if _, exists := result[normKey]; exists {
			return nil, fmt.Errorf("envflatten: duplicate key after normalisation: %q", normKey)
		}
		result[normKey] = v
	}
	return result, nil
}

// Prefix prepends a prefix (followed by the separator) to every key in env.
func Prefix(env map[string]string, prefix string, opts Options) (map[string]string, error) {
	if env == nil {
		return nil, errors.New("envflatten: env map must not be nil")
	}
	if prefix == "" {
		return nil, errors.New("envflatten: prefix must not be empty")
	}
	if opts.Separator == "" {
		return nil, errors.New("envflatten: separator must not be empty")
	}

	pfx := prefix
	if opts.UppercaseKeys {
		pfx = strings.ToUpper(pfx)
	}

	result := make(map[string]string, len(env))
	for k, v := range env {
		normKey := pfx + opts.Separator + k
		if opts.UppercaseKeys {
			normKey = strings.ToUpper(normKey)
		}
		result[normKey] = v
	}
	return result, nil
}
