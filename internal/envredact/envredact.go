// Package envredact provides utilities for redacting sensitive values
// from environment variable maps before display or logging.
package envredact

import (
	"regexp"
	"strings"
)

// DefaultSensitivePatterns is a list of key name patterns considered sensitive.
var DefaultSensitivePatterns = []string{
	"(?i)password",
	"(?i)secret",
	"(?i)token",
	"(?i)api_key",
	"(?i)private_key",
	"(?i)auth",
	"(?i)credential",
}

const redactedValue = "[REDACTED]"

// Options configures redaction behaviour.
type Options struct {
	// Patterns is a list of regex patterns matched against key names.
	// Keys matching any pattern will have their values redacted.
	Patterns []string
}

// DefaultOptions returns Options populated with DefaultSensitivePatterns.
func DefaultOptions() Options {
	return Options{Patterns: DefaultSensitivePatterns}
}

// Redact returns a copy of env with sensitive values replaced by [REDACTED].
// Keys are matched case-insensitively against the configured patterns.
func Redact(env map[string]string, opts Options) (map[string]string, error) {
	compiled, err := compilePatterns(opts.Patterns)
	if err != nil {
		return nil, err
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		if isSensitive(k, compiled) {
			out[k] = redactedValue
		} else {
			out[k] = v
		}
	}
	return out, nil
}

// IsSensitiveKey reports whether the given key matches any of the provided
// regex patterns.
func IsSensitiveKey(key string, patterns []string) (bool, error) {
	compiled, err := compilePatterns(patterns)
	if err != nil {
		return false, err
	}
	return isSensitive(key, compiled), nil
}

func isSensitive(key string, patterns []*regexp.Regexp) bool {
	upper := strings.ToUpper(key)
	for _, re := range patterns {
		if re.MatchString(upper) {
			return true
		}
	}
	return false
}

func compilePatterns(patterns []string) ([]*regexp.Regexp, error) {
	var compiled []*regexp.Regexp
	for _, p := range patterns {
		re, err := regexp.Compile(strings.ToUpper(p))
		if err != nil {
			return nil, err
		}
		compiled = append(compiled, re)
	}
	return compiled, nil
}
