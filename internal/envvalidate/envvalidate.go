// Package envvalidate provides type-aware validation of environment variable values.
// It supports rules such as required, type constraints, min/max length, and regex patterns.
package envvalidate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Rule defines validation constraints for a single environment variable.
type Rule struct {
	Key      string
	Required bool
	Type     string // "string", "int", "bool", "float"
	MinLen   int
	MaxLen   int    // 0 means no limit
	Pattern  string // optional regex
}

// Violation describes a single validation failure.
type Violation struct {
	Key     string
	Message string
}

func (v Violation) Error() string {
	return fmt.Sprintf("%s: %s", v.Key, v.Message)
}

// Check validates the provided env map against the given rules.
// It returns a slice of Violations (empty if all pass).
func Check(env map[string]string, rules []Rule) []Violation {
	var violations []Violation

	for _, r := range rules {
		val, exists := env[r.Key]

		if r.Required && (!exists || strings.TrimSpace(val) == "") {
			violations = append(violations, Violation{Key: r.Key, Message: "required but missing or empty"})
			continue
		}

		if !exists {
			continue
		}

		if r.Type != "" {
			if err := checkType(r.Key, val, r.Type); err != nil {
				violations = append(violations, *err)
			}
		}

		if r.MinLen > 0 && len(val) < r.MinLen {
			violations = append(violations, Violation{Key: r.Key, Message: fmt.Sprintf("value too short (min %d)", r.MinLen)})
		}

		if r.MaxLen > 0 && len(val) > r.MaxLen {
			violations = append(violations, Violation{Key: r.Key, Message: fmt.Sprintf("value too long (max %d)", r.MaxLen)})
		}

		if r.Pattern != "" {
			re, err := regexp.Compile(r.Pattern)
			if err != nil {
				violations = append(violations, Violation{Key: r.Key, Message: fmt.Sprintf("invalid pattern: %v", err)})
			} else if !re.MatchString(val) {
				violations = append(violations, Violation{Key: r.Key, Message: fmt.Sprintf("does not match pattern %q", r.Pattern)})
			}
		}
	}

	return violations
}

// HasViolations returns true if any violations were found.
func HasViolations(violations []Violation) bool {
	return len(violations) > 0
}

func checkType(key, val, typ string) *Violation {
	switch typ {
	case "int":
		if _, err := strconv.Atoi(val); err != nil {
			return &Violation{Key: key, Message: "expected integer value"}
		}
	case "bool":
		lower := strings.ToLower(val)
		if lower != "true" && lower != "false" && lower != "1" && lower != "0" {
			return &Violation{Key: key, Message: "expected boolean value (true/false/1/0)"}
		}
	case "float":
		if _, err := strconv.ParseFloat(val, 64); err != nil {
			return &Violation{Key: key, Message: "expected float value"}
		}
	}
	return nil
}
