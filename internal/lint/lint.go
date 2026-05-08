// Package lint provides validation for .env file contents,
// checking for common issues such as missing values, suspicious keys,
// and formatting problems before sealing.
package lint

import (
	"bufio"
	"fmt"
	"strings"
)

// Issue represents a single lint finding.
type Issue struct {
	Line    int
	Key     string
	Message string
	Severity string // "error" or "warning"
}

func (i Issue) String() string {
	return fmt.Sprintf("%s (line %d): %s — %s", i.Severity, i.Line, i.Key, i.Message)
}

// Check validates the contents of a .env file and returns any issues found.
func Check(content string) []Issue {
	var issues []Issue
	scanner := bufio.NewScanner(strings.NewReader(content))
	lineNum := 0
	seenKeys := make(map[string]int)

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		// Skip blank lines and comments
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Must contain '='
		eqIdx := strings.Index(trimmed, "=")
		if eqIdx < 0 {
			issues = append(issues, Issue{
				Line:     lineNum,
				Key:      trimmed,
				Message:  "line is not a valid KEY=VALUE pair",
				Severity: "error",
			})
			continue
		}

		key := strings.TrimSpace(trimmed[:eqIdx])
		value := strings.TrimSpace(trimmed[eqIdx+1:])

		// Empty key
		if key == "" {
			issues = append(issues, Issue{
				Line:     lineNum,
				Key:      "(empty)",
				Message:  "key must not be empty",
				Severity: "error",
			})
			continue
		}

		// Duplicate key
		if prev, ok := seenKeys[key]; ok {
			issues = append(issues, Issue{
				Line:     lineNum,
				Key:      key,
				Message:  fmt.Sprintf("duplicate key (first seen on line %d)", prev),
				Severity: "warning",
			})
		} else {
			seenKeys[key] = lineNum
		}

		// Empty value warning
		if value == "" {
			issues = append(issues, Issue{
				Line:     lineNum,
				Key:      key,
				Message:  "value is empty",
				Severity: "warning",
			})
		}
	}

	return issues
}

// HasErrors returns true if any issue has severity "error".
func HasErrors(issues []Issue) bool {
	for _, i := range issues {
		if i.Severity == "error" {
			return true
		}
	}
	return false
}
