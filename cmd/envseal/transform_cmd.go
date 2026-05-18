package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envseal/internal/envtransform"
)

func runTransform(args []string) error {
	fs := flag.NewFlagSet("transform", flag.ContinueOnError)
	addPrefix := fs.String("add-prefix", "", "prefix to add to all keys")
	stripPrefix := fs.String("strip-prefix", "", "prefix to strip from all keys")
	filter := fs.String("filter", "", "regex pattern to filter keys")
	renameRaw := fs.String("rename", "", "comma-separated OLD=NEW rename pairs")
	inputFile := fs.String("input", "", "JSON file containing env map (required)")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *inputFile == "" {
		return fmt.Errorf("transform: --input is required")
	}

	data, err := os.ReadFile(*inputFile)
	if err != nil {
		return fmt.Errorf("transform: reading input: %w", err)
	}
	env := map[string]string{}
	if err := json.Unmarshal(data, &env); err != nil {
		return fmt.Errorf("transform: parsing input JSON: %w", err)
	}

	renames := map[string]string{}
	if *renameRaw != "" {
		for _, pair := range strings.Split(*renameRaw, ",") {
			parts := strings.SplitN(strings.TrimSpace(pair), "=", 2)
			if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
				return fmt.Errorf("transform: invalid rename pair %q", pair)
			}
			renames[parts[0]] = parts[1]
		}
	}

	out, err := envtransform.Transform(env, envtransform.Options{
		AddPrefix:     *addPrefix,
		StripPrefix:   *stripPrefix,
		FilterPattern: *filter,
		RenameKeys:    renames,
	})
	if err != nil {
		return fmt.Errorf("transform: %w", err)
	}

	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(out)
}
