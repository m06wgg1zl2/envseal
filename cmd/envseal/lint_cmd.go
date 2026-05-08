package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/yourorg/envseal/internal/lint"
)

// runLint reads a .env file and reports any lint issues.
// Exits with a non-zero status if errors are found.
func runLint(args []string) error {
	fs := flag.NewFlagSet("lint", flag.ContinueOnError)
	envFile := fs.String("env", ".env", "path to the .env file to lint")
	strict := fs.Bool("strict", false, "exit non-zero on warnings as well as errors")

	if err := fs.Parse(args); err != nil {
		return err
	}

	data, err := os.ReadFile(*envFile)
	if err != nil {
		return fmt.Errorf("reading %s: %w", *envFile, err)
	}

	issues := lint.Check(string(data))

	if len(issues) == 0 {
		fmt.Printf("✔ %s looks good — no issues found\n", *envFile)
		return nil
	}

	for _, issue := range issues {
		fmt.Println(issue)
	}

	hasErr := lint.HasErrors(issues)
	if hasErr || *strict {
		return fmt.Errorf("lint failed: %d issue(s) found in %s", len(issues), *envFile)
	}

	// Warnings only, non-strict mode: print summary but succeed.
	fmt.Printf("⚠ %d warning(s) in %s\n", len(issues), *envFile)
	return nil
}
