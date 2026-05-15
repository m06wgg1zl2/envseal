// Package envdefaults provides utilities for applying default values to
// environment variable maps when keys are missing or empty.
package envdefaults

import "fmt"

// Default represents a single default entry with an optional description.
type Default struct {
	Key         string
	Value       string
	Description string
	Overwrite   bool // if true, overwrite even when key already exists
}

// Result holds the outcome of applying defaults.
type Result struct {
	Applied  []string // keys that were set
	Skipped  []string // keys that already had a value
	Overwritten []string // keys that were overwritten
}

// Apply merges defaults into env. Keys present in env are skipped unless
// the Default has Overwrite set to true. The original map is not mutated;
// a new map is returned.
func Apply(env map[string]string, defaults []Default) (map[string]string, Result, error) {
	if env == nil {
		return nil, Result{}, fmt.Errorf("envdefaults: env map must not be nil")
	}

	out := make(map[string]string, len(env))
	for k, v := range env {
		out[k] = v
	}

	var result Result
	for _, d := range defaults {
		if d.Key == "" {
			return nil, Result{}, fmt.Errorf("envdefaults: default entry has empty key")
		}
		existing, exists := out[d.Key]
		switch {
		case !exists || existing == "":
			out[d.Key] = d.Value
			result.Applied = append(result.Applied, d.Key)
		case d.Overwrite:
			out[d.Key] = d.Value
			result.Overwritten = append(result.Overwritten, d.Key)
		default:
			result.Skipped = append(result.Skipped, d.Key)
		}
	}
	return out, result, nil
}

// FromMap builds a []Default slice from a plain key→value map using empty
// descriptions and Overwrite=false. Useful for quick one-liners.
func FromMap(m map[string]string) []Default {
	defaults := make([]Default, 0, len(m))
	for k, v := range m {
		defaults = append(defaults, Default{Key: k, Value: v})
	}
	return defaults
}
