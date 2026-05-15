// Package envcast provides typed conversion utilities for environment variable values.
// It supports parsing strings into common Go types with optional default fallbacks.
package envcast

import (
	"fmt"
	"strconv"
	"strings"
)

// Result holds a parsed value and any conversion error.
type Result struct {
	Key   string
	Raw   string
	Error error
}

// ToString returns the raw string value for the given key.
func ToString(env map[string]string, key, defaultVal string) string {
	if v, ok := env[key]; ok && v != "" {
		return v
	}
	return defaultVal
}

// ToInt parses the value as an integer, returning defaultVal on failure.
func ToInt(env map[string]string, key string, defaultVal int) (int, error) {
	v, ok := env[key]
	if !ok || v == "" {
		return defaultVal, nil
	}
	n, err := strconv.Atoi(strings.TrimSpace(v))
	if err != nil {
		return defaultVal, fmt.Errorf("envcast: key %q: cannot parse %q as int: %w", key, v, err)
	}
	return n, nil
}

// ToBool parses the value as a boolean (true/false/1/0/yes/no), returning defaultVal on failure.
func ToBool(env map[string]string, key string, defaultVal bool) (bool, error) {
	v, ok := env[key]
	if !ok || v == "" {
		return defaultVal, nil
	}
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "true", "1", "yes", "on":
		return true, nil
	case "false", "0", "no", "off":
		return false, nil
	default:
		return defaultVal, fmt.Errorf("envcast: key %q: cannot parse %q as bool", key, v)
	}
}

// ToFloat64 parses the value as a float64, returning defaultVal on failure.
func ToFloat64(env map[string]string, key string, defaultVal float64) (float64, error) {
	v, ok := env[key]
	if !ok || v == "" {
		return defaultVal, nil
	}
	f, err := strconv.ParseFloat(strings.TrimSpace(v), 64)
	if err != nil {
		return defaultVal, fmt.Errorf("envcast: key %q: cannot parse %q as float64: %w", key, v, err)
	}
	return f, nil
}

// ToStringSlice splits the value by sep and trims whitespace from each element.
func ToStringSlice(env map[string]string, key, sep string) []string {
	v, ok := env[key]
	if !ok || v == "" {
		return nil
	}
	parts := strings.Split(v, sep)
	result := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			result = append(result, t)
		}
	}
	return result
}

// Validate runs all conversions and returns any errors encountered.
func Validate(results []Result) []error {
	var errs []error
	for _, r := range results {
		if r.Error != nil {
			errs = append(errs, r.Error)
		}
	}
	return errs
}
