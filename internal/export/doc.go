// Package export renders decrypted environment variables into various output
// formats suitable for consumption by shells and other tooling.
//
// Supported formats:
//
//   - shell  — emits `export KEY="value"` lines, suitable for eval in bash/zsh.
//   - dotenv — emits `KEY=value` lines in standard dotenv format.
//   - json   — emits a JSON object mapping keys to values.
//
// Callers may optionally restrict output to a subset of keys via Options.Keys.
// Comments and blank lines in the source env content are always stripped.
package export
