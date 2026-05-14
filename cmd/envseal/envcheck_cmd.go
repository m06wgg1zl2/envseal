package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourusername/envseal/internal/envcheck"
)

// runEnvCheck compares a (decrypted) .env file against a .env.template and
// reports any keys that are missing from or extra in the env file.
//
// Usage:
//
//	envseal check [-env <path>] [-template <path>]
func runEnvCheck(args []string) error {
	fs := flag.NewFlagSet("check", flag.ContinueOnError)
	envPath := fs.String("env", ".env", "path to the .env file")
	tplPath := fs.String("template", ".env.template", "path to the .env.template file")

	if err := fs.Parse(args); err != nil {
		return err
	}

	res, err := envcheck.Check(*envPath, *tplPath)
	if err != nil {
		return fmt.Errorf("envcheck: %w", err)
	}

	if !res.HasIssues() {
		fmt.Println("✓ .env matches template")
		return nil
	}

	fmt.Fprintf(os.Stderr, "✗ .env does not match template:\n%s", res.Format())
	return fmt.Errorf("compliance check failed: %d missing, %d extra",
		len(res.MissingKeys), len(res.ExtraKeys))
}
