// Package envschema provides functionality for generating and validating
// JSON Schema documents from .env files and schema definitions.
//
// A schema describes the expected keys, types, constraints, and documentation
// for an environment configuration. It can be used to validate .env files,
// generate documentation, or onboard new team members.
package envschema

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// FieldType represents the expected type of an environment variable.
type FieldType string

const (
	TypeString  FieldType = "string"
	TypeInt     FieldType = "int"
	TypeBool    FieldType = "bool"
	TypeFloat   FieldType = "float"
	TypeURL     FieldType = "url"
	TypeEmail   FieldType = "email"
)

// Field describes a single environment variable in the schema.
type Field struct {
	Key         string    `json:"key"`
	Type        FieldType `json:"type"`
	Required    bool      `json:"required"`
	Description string    `json:"description,omitempty"`
	Default     string    `json:"default,omitempty"`
	Pattern     string    `json:"pattern,omitempty"`
}

// Schema represents the full schema for an environment configuration.
type Schema struct {
	Version string           `json:"version"`
	Fields  map[string]Field `json:"fields"`
}

// ValidationError describes a single schema violation.
type ValidationError struct {
	Key     string
	Message string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", e.Key, e.Message)
}

// Generate creates a Schema by inspecting the keys in a .env file.
// All fields are inferred as optional strings. The caller can enrich
// the schema after generation.
func Generate(envPath string) (*Schema, error) {
	keys, err := parseKeys(envPath)
	if err != nil {
		return nil, fmt.Errorf("envschema: read %s: %w", envPath, err)
	}

	fields := make(map[string]Field, len(keys))
	for _, k := range keys {
		fields[k] = Field{
			Key:  k,
			Type: TypeString,
		}
	}

	return &Schema{
		Version: "1",
		Fields:  fields,
	}, nil
}

// Save writes the schema to a JSON file at the given path.
func Save(s *Schema, path string) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("envschema: create %s: %w", path, err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	if err := enc.Encode(s); err != nil {
		return fmt.Errorf("envschema: encode: %w", err)
	}
	return nil
}

// Load reads a schema from a JSON file.
func Load(path string) (*Schema, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("envschema: open %s: %w", path, err)
	}
	defer f.Close()

	var s Schema
	if err := json.NewDecoder(f).Decode(&s); err != nil {
		return nil, fmt.Errorf("envschema: decode: %w", err)
	}
	return &s, nil
}

// Validate checks an env map against the schema and returns any violations.
func Validate(s *Schema, env map[string]string) []ValidationError {
	var errs []ValidationError

	for key, field := range s.Fields {
		val, present := env[key]
		if !present || val == "" {
			if field.Required {
				errs = append(errs, ValidationError{Key: key, Message: "required key is missing or empty"})
			}
			continue
		}

		if field.Pattern != "" {
			if !matchesPattern(field.Pattern, val) {
				errs = append(errs, ValidationError{Key: key, Message: fmt.Sprintf("value does not match pattern %q", field.Pattern)})
			}
		}
	}

	return errs
}

// HasViolations returns true if Validate produced any errors.
func HasViolations(errs []ValidationError) bool {
	return len(errs) > 0
}

// parseKeys reads the keys (not values) from a .env file.
func parseKeys(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var keys []string
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
			keys = append(keys, key)
		}
	}
	return keys, scanner.Err()
}

// matchesPattern performs a simple glob-style prefix/suffix/contains match.
// For full regex support callers should use internal/policy.
func matchesPattern(pattern, value string) bool {
	if strings.HasPrefix(pattern, "*") && strings.HasSuffix(pattern, "*") {
		return strings.Contains(value, strings.Trim(pattern, "*"))
	}
	if strings.HasPrefix(pattern, "*") {
		return strings.HasSuffix(value, strings.TrimPrefix(pattern, "*"))
	}
	if strings.HasSuffix(pattern, "*") {
		return strings.HasPrefix(value, strings.TrimSuffix(pattern, "*"))
	}
	return value == pattern
}
