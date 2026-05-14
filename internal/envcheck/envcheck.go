// Package envcheck provides functionality to compare a decrypted .env file
// against a .env.template to detect missing or extra keys.
package envcheck

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Result holds the outcome of a template compliance check.
type Result struct {
	MissingKeys []string // keys present in template but absent in env
	ExtraKeys   []string // keys present in env but absent in template
}

// HasIssues returns true if any missing or extra keys were found.
func (r Result) HasIssues() bool {
	return len(r.MissingKeys) > 0 || len(r.ExtraKeys) > 0
}

// Format returns a human-readable summary of the result.
func (r Result) Format() string {
	var sb strings.Builder
	for _, k := range r.MissingKeys {
		fmt.Fprintf(&sb, "  missing : %s\n", k)
	}
	for _, k := range r.ExtraKeys {
		fmt.Fprintf(&sb, "  extra   : %s\n", k)
	}
	return sb.String()
}

// Check compares the keys in envPath against those declared in templatePath.
// Comment lines and blank lines are ignored in both files.
func Check(envPath, templatePath string) (Result, error) {
	envKeys, err := parseKeys(envPath)
	if err != nil {
		return Result{}, fmt.Errorf("reading env file: %w", err)
	}
	tplKeys, err := parseKeys(templatePath)
	if err != nil {
		return Result{}, fmt.Errorf("reading template file: %w", err)
	}

	var res Result
	for k := range tplKeys {
		if !envKeys[k] {
			res.MissingKeys = append(res.MissingKeys, k)
		}
	}
	for k := range envKeys {
		if !tplKeys[k] {
			res.ExtraKeys = append(res.ExtraKeys, k)
		}
	}
	return res, nil
}

// parseKeys reads a file and returns a set of key names found on KEY=... lines.
func parseKeys(path string) (map[string]bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	keys := make(map[string]bool)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		if key != "" {
			keys[key] = true
		}
	}
	return keys, scanner.Err()
}
