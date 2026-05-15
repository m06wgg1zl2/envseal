package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/yourorg/envseal/internal/envvalidate"
	"github.com/yourorg/envseal/internal/seal"
)

// runValidate decrypts the sealed file and validates its contents against a
// JSON rules file. Exits with a non-zero status if any violations are found.
//
// Usage: envseal validate --sealed <file> --rules <rules.json> --key <keyfile>
func runValidate(args []string) error {
	var sealedFile, rulesFile, keyFile string

	for i := 0; i < len(args)-1; i++ {
		switch args[i] {
		case "--sealed":
			sealedFile = args[i+1]
		case "--rules":
			rulesFile = args[i+1]
		case "--key":
			keyFile = args[i+1]
		}
	}

	if sealedFile == "" {
		return fmt.Errorf("--sealed is required")
	}
	if rulesFile == "" {
		return fmt.Errorf("--rules is required")
	}
	if keyFile == "" {
		return fmt.Errorf("--key is required")
	}

	// Load rules from JSON file.
	data, err := os.ReadFile(rulesFile)
	if err != nil {
		return fmt.Errorf("reading rules file: %w", err)
	}
	var rules []envvalidate.Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return fmt.Errorf("parsing rules file: %w", err)
	}

	// Unseal to a temp file, then read env.
	tmpFile, err := os.CreateTemp("", "envseal-validate-*.env")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if err := seal.Unseal(sealedFile, tmpFile.Name(), keyFile); err != nil {
		return fmt.Errorf("unsealing: %w", err)
	}

	env, err := parseEnvFile(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("parsing env: %w", err)
	}

	violations := envvalidate.Check(env, rules)
	if envvalidate.HasViolations(violations) {
		fmt.Fprintf(os.Stderr, "Validation failed with %d violation(s):\n", len(violations))
		for _, v := range violations {
			fmt.Fprintf(os.Stderr, "  - %s\n", v.Error())
		}
		os.Exit(1)
	}

	fmt.Println("Validation passed.")
	return nil
}

// parseEnvFile reads a KEY=VALUE env file into a map.
func parseEnvFile(path string) (map[string]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	env := make(map[string]string)
	for _, line := range splitLines(string(data)) {
		if line == "" || line[0] == '#' {
			continue
		}
		if idx := indexByte(line, '='); idx >= 0 {
			env[line[:idx]] = line[idx+1:]
		}
	}
	return env, nil
}

func splitLines(s string) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func indexByte(s string, b byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == b {
			return i
		}
	}
	return -1
}
