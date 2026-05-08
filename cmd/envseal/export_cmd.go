package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envseal/internal/export"
	"github.com/yourorg/envseal/internal/keystore"
	"github.com/yourorg/envseal/internal/seal"
)

// runExport decrypts a sealed file and writes its contents to stdout in the
// requested format. Usage:
//
//	envseal export [--format shell|dotenv|json] [--key KEY,...] <sealed-file>
func runExport(args []string) error {
	fs := flag.NewFlagSet("export", flag.ContinueOnError)
	formatFlag := fs.String("format", "dotenv", "output format: shell, dotenv, json")
	keysFlag := fs.String("key", "", "comma-separated list of keys to export (default: all)")
	keyPath := fs.String("identity", keystore.DefaultKeyPath(), "path to age identity file")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() < 1 {
		return fmt.Errorf("export: sealed file path required")
	}
	sealedPath := fs.Arg(0)

	identity, err := keystore.LoadIdentity(*keyPath)
	if err != nil {
		return fmt.Errorf("export: load identity: %w", err)
	}

	envContent, err := seal.Unseal(sealedPath, identity)
	if err != nil {
		return fmt.Errorf("export: unseal: %w", err)
	}

	var filterKeys []string
	if *keysFlag != "" {
		for _, k := range splitComma(*keysFlag) {
			if k != "" {
				filterKeys = append(filterKeys, k)
			}
		}
	}

	opts := export.Options{
		Format: export.Format(*formatFlag),
		Keys:   filterKeys,
	}

	if err := export.Write(os.Stdout, envContent, opts); err != nil {
		return fmt.Errorf("export: write: %w", err)
	}
	return nil
}

// splitComma splits s on commas and trims whitespace from each element.
func splitComma(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			part := strings.TrimSpace(s[start:i])
			out = append(out, part)
			start = i + 1
		}
	}
	return out
}
