// Package policy enforces environment variable rules such as required keys,
// forbidden keys, and naming conventions across sealed env files.
package policy

import (
	"fmt"
	"regexp"
	"strings"
)

// Rule defines a single policy rule applied to env variables.
type Rule struct {
	Required  []string `json:"required,omitempty"`
	Forbidden []string `json:"forbidden,omitempty"`
	Pattern   string   `json:"pattern,omitempty"`
}

// Violation represents a single policy violation.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("policy violation for %q: %s", v.Key, v.Message)
}

// Check validates the given env map against the provided Rule.
// It returns a list of Violations (empty if all rules pass).
func Check(env map[string]string, rule Rule) []Violation {
	var violations []Violation

	for _, key := range rule.Required {
		if _, ok := env[key]; !ok {
			violations = append(violations, Violation{
				Key:     key,
				Message: "required key is missing",
			})
		}
	}

	for _, key := range rule.Forbidden {
		if _, ok := env[key]; ok {
			violations = append(violations, Violation{
				Key:     key,
				Message: "forbidden key is present",
			})
		}
	}

	if rule.Pattern != "" {
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			violations = append(violations, Violation{
				Key:     "__pattern__",
				Message: fmt.Sprintf("invalid pattern %q: %v", rule.Pattern, err),
			})
		} else {
			for key := range env {
				if !re.MatchString(key) {
					violations = append(violations, Violation{
						Key:     key,
						Message: fmt.Sprintf("key does not match required pattern %q", rule.Pattern),
					})
				}
			}
		}
	}

	return violations
}

// HasViolations returns true if any violations exist.
func HasViolations(violations []Violation) bool {
	return len(violations) > 0
}

// Format returns a human-readable summary of violations.
func Format(violations []Violation) string {
	if len(violations) == 0 {
		return "no policy violations found"
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d policy violation(s):\n", len(violations)))
	for _, v := range violations {
		sb.WriteString(fmt.Sprintf("  - %s\n", v.Error()))
	}
	return strings.TrimRight(sb.String(), "\n")
}
