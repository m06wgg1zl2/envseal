// Package envprofile provides named environment profile management for envseal.
//
// Profiles allow teams to define named sets of environment variable overrides
// (e.g. "dev", "staging", "prod") and switch between them quickly. Profiles
// are stored in a JSON file alongside sealed env files and are safe to commit
// since they contain only non-secret overrides.
//
// Typical usage:
//
//	s, err := envprofile.Load(".envseal-profiles.json")
//	_ = envprofile.Set(s, "dev", map[string]string{"DEBUG": "true"})
//	_ = envprofile.Switch(s, "dev")
//	vars := envprofile.ActiveVars(s)
//	_ = envprofile.Save(".envseal-profiles.json", s)
package envprofile
