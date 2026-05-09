package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envseal/internal/policy"
	"github.com/yourorg/envseal/internal/seal"
)

// runPolicy checks a sealed env file against a policy rule defined via flags.
//
// Usage:
//
//	envseal policy --sealed=.env.sealed --required=KEY1,KEY2 --forbidden=DEBUG --pattern='^[A-Z_]+$'
func runPolicy(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: envseal policy --sealed=<file> [--required=K1,K2] [--forbidden=K] [--pattern=REGEX]")
	}

	var sealedFile, requiredRaw, forbiddenRaw, pattern string
	for _, arg := range args {
		switch {
		case strings.HasPrefix(arg, "--sealed="):
			sealedFile = strings.TrimPrefix(arg, "--sealed=")
		case strings.HasPrefix(arg, "--required="):
			requiredRaw = strings.TrimPrefix(arg, "--required=")
		case strings.HasPrefix(arg, "--forbidden="):
			forbiddenRaw = strings.TrimPrefix(arg, "--forbidden=")
		case strings.HasPrefix(arg, "--pattern="):
			pattern = strings.TrimPrefix(arg, "--pattern=")
		}
	}

	if sealedFile == "" {
		return fmt.Errorf("--sealed flag is required")
	}

	keyPath := os.Getenv("ENVSEAL_KEY")
	if keyPath == "" {
		keyPath = defaultKeyPath()
	}

	plaintext, err := seal.Unseal(sealedFile, keyPath)
	if err != nil {
		return fmt.Errorf("unseal failed: %w", err)
	}

	env, err := parseEnvJSON(plaintext)
	if err != nil {
		return fmt.Errorf("failed to parse env: %w", err)
	}

	rule := policy.Rule{
		Pattern: pattern,
	}
	if requiredRaw != "" {
		rule.Required = splitComma(requiredRaw)
	}
	if forbiddenRaw != "" {
		rule.Forbidden = splitComma(forbiddenRaw)
	}

	violations := policy.Check(env, rule)
	fmt.Println(policy.Format(violations))

	if policy.HasViolations(violations) {
		return fmt.Errorf("policy check failed")
	}
	return nil
}

// parseEnvJSON decodes a JSON-encoded map[string]string from plaintext bytes.
func parseEnvJSON(data []byte) (map[string]string, error) {
	var env map[string]string
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, err
	}
	return env, nil
}
