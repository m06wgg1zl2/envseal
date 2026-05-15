// Package envclone provides utilities for cloning .env files between
// environments or configurations.
//
// It supports selective key inclusion and exclusion, key renaming, and
// controlled overwrite behaviour. Comments and blank lines are preserved
// when no filtering would remove their associated context.
//
// Typical usage:
//
//	err := envclone.Clone(".env", ".env.staging", envclone.Options{
//		ExcludeKeys: []string{"DATABASE_URL"},
//		Rename:      map[string]string{"APP_PORT": "PORT"},
//		Overwrite:   false,
//	})
package envclone
