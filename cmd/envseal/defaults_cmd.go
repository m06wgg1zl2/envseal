package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourusername/envseal/internal/envdefaults"
)

// runDefaults applies a set of default key=value pairs to a decrypted .env
// file, printing the merged result to stdout or writing it to a file.
//
// Usage:
//
//	envseal defaults -defaults KEY=VAL,KEY2=VAL2 [-overwrite] [-out FILE] [-json]
func runDefaults(args []string) error {
	fs := flag.NewFlagSet("defaults", flag.ContinueOnError)
	rawDefaults := fs.String("defaults", "", "comma-separated KEY=VALUE pairs to apply as defaults (required)")
	overwrite := fs.Bool("overwrite", false, "overwrite existing keys with default values")
	outFile := fs.String("out", "", "write result to file instead of stdout")
	jsonOut := fs.Bool("json", false, "output as JSON object")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *rawDefaults == "" {
		return fmt.Errorf("defaults: -defaults flag is required")
	}

	// Parse KEY=VALUE pairs from flag.
	pairs := strings.Split(*rawDefaults, ",")
	var defs []envdefaults.Default
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		idx := strings.IndexByte(pair, '=')
		if idx <= 0 {
			return fmt.Errorf("defaults: invalid pair %q, expected KEY=VALUE", pair)
		}
		defs = append(defs, envdefaults.Default{
			Key:       pair[:idx],
			Value:     pair[idx+1:],
			Overwrite: *overwrite,
		})
	}

	// Build env from remaining positional args (KEY=VALUE) or an empty map.
	env := make(map[string]string)
	for _, arg := range fs.Args() {
		idx := strings.IndexByte(arg, '=')
		if idx > 0 {
			env[arg[:idx]] = arg[idx+1:]
		}
	}

	out, result, err := envdefaults.Apply(env, defs)
	if err != nil {
		return err
	}

	// Format output.
	var sb strings.Builder
	if *jsonOut {
		b, err := json.MarshalIndent(out, "", "  ")
		if err != nil {
			return err
		}
		sb.Write(b)
		sb.WriteByte('\n')
	} else {
		for k, v := range out {
			fmt.Fprintf(&sb, "%s=%s\n", k, v)
		}
	}

	if *outFile != "" {
		if err := os.WriteFile(*outFile, []byte(sb.String()), 0o600); err != nil {
			return fmt.Errorf("defaults: writing output: %w", err)
		}
	} else {
		fmt.Print(sb.String())
	}

	fmt.Fprintf(os.Stderr, "defaults: applied=%d skipped=%d overwritten=%d\n",
		len(result.Applied), len(result.Skipped), len(result.Overwritten))
	return nil
}
