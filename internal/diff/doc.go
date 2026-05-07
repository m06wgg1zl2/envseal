// Package diff compares two plaintext .env file snapshots and produces
// a structured list of key-level changes (added, removed, modified).
//
// Values are intentionally masked in all human-readable output to prevent
// secrets from appearing in logs or terminal output. Only key names and
// change types are surfaced.
//
// Typical usage:
//
//	changes := diff.Compare(oldEnvContent, newEnvContent)
//	fmt.Println(diff.Summary(changes))
package diff
