// Package envresolve provides functionality to resolve variable references
// within .env files, expanding ${VAR} and $VAR style references.
package envresolve

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var refPattern = regexp.MustCompile(`\$\{([^}]+)\}|\$([A-Za-z_][A-Za-z0-9_]*)`)

// Options configures resolution behaviour.
type Options struct {
	// FallbackToOS allows falling back to OS environment variables
	// when a key is not found in the provided env map.
	FallbackToOS bool
	// FailOnMissing returns an error if a referenced variable cannot be resolved.
	FailOnMissing bool
}

// DefaultOptions returns sensible defaults.
func DefaultOptions() Options {
	return Options{
		FallbackToOS:  true,
		FailOnMissing: false,
	}
}

// Resolve expands variable references in the values of env using the keys
// already present in env (and optionally OS environment variables).
// It returns a new map with all resolvable references expanded.
func Resolve(env map[string]string, opts Options) (map[string]string, error) {
	resolved := make(map[string]string, len(env))
	for k, v := range env {
		resolved[k] = v
	}

	for k, v := range resolved {
		expanded, err := expand(v, resolved, opts)
		if err != nil {
			return nil, fmt.Errorf("resolving %q: %w", k, err)
		}
		resolved[k] = expanded
	}

	return resolved, nil
}

func expand(value string, env map[string]string, opts Options) (string, error) {
	var resolveErr error
	result := refPattern.ReplaceAllStringFunc(value, func(match string) string {
		if resolveErr != nil {
			return match
		}
		key := extractKey(match)
		if val, ok := env[key]; ok {
			return val
		}
		if opts.FallbackToOS {
			if val, ok := os.LookupEnv(key); ok {
				return val
			}
		}
		if opts.FailOnMissing {
			resolveErr = fmt.Errorf("undefined variable %q", key)
			return match
		}
		return ""
	})
	if resolveErr != nil {
		return "", resolveErr
	}
	return result, nil
}

func extractKey(match string) string {
	match = strings.TrimPrefix(match, "$")
	match = strings.TrimPrefix(match, "{")
	match = strings.TrimSuffix(match, "}")
	return match
}
