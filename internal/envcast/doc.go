// Package envcast provides typed conversion helpers for environment variable maps.
//
// It is designed to work alongside the envseal unsealed environment data,
// allowing callers to safely parse string values into Go primitives such as
// int, bool, float64, and []string — each with a configurable default fallback.
//
// Example usage:
//
//	port, err := envcast.ToInt(env, "PORT", 8080)
//	debug, err := envcast.ToBool(env, "DEBUG", false)
//	hosts := envcast.ToStringSlice(env, "ALLOWED_HOSTS", ",")
//
// Validate can be used to collect multiple conversion errors at once:
//
//	results := []envcast.Result{
//	    {Key: "PORT", Error: portErr},
//	    {Key: "DEBUG", Error: debugErr},
//	}
//	if errs := envcast.Validate(results); len(errs) > 0 { ... }
package envcast
