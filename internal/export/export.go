// Package export provides functionality to export decrypted .env variables
// into various formats such as shell export statements or dotenv format.
package export

import (
	"fmt"
	"io"
	"strings"
)

// Format represents the output format for exported variables.
type Format string

const (
	// FormatShell outputs variables as shell export statements.
	FormatShell Format = "shell"
	// FormatDotenv outputs variables in standard dotenv format.
	FormatDotenv Format = "dotenv"
	// FormatJSON outputs variables as a JSON object.
	FormatJSON Format = "json"
)

// Options configures the export behaviour.
type Options struct {
	Format Format
	// Keys filters output to only these keys. If empty, all keys are exported.
	Keys []string
}

// Write parses the env content and writes it to w in the requested format.
func Write(w io.Writer, envContent string, opts Options) error {
	vars, err := parse(envContent)
	if err != nil {
		return fmt.Errorf("export: parse env: %w", err)
	}

	if len(opts.Keys) > 0 {
		vars = filter(vars, opts.Keys)
	}

	switch opts.Format {
	case FormatShell:
		return writeShell(w, vars)
	case FormatDotenv:
		return writeDotenv(w, vars)
	case FormatJSON:
		return writeJSON(w, vars)
	default:
		return fmt.Errorf("export: unknown format %q", opts.Format)
	}
}

func writeShell(w io.Writer, vars []entry) error {
	for _, e := range vars {
		if _, err := fmt.Fprintf(w, "export %s=%q\n", e.key, e.value); err != nil {
			return err
		}
	}
	return nil
}

func writeDotenv(w io.Writer, vars []entry) error {
	for _, e := range vars {
		if _, err := fmt.Fprintf(w, "%s=%s\n", e.key, e.value); err != nil {
			return err
		}
	}
	return nil
}

func writeJSON(w io.Writer, vars []entry) error {
	fmt.Fprint(w, "{\n")
	for i, e := range vars {
		comma := ","
		if i == len(vars)-1 {
			comma = ""
		}
		if _, err := fmt.Fprintf(w, "  %q: %q%s\n", e.key, e.value, comma); err != nil {
			return err
		}
	}
	fmt.Fprint(w, "}\n")
	return nil
}

type entry struct {
	key   string
	value string
}

func parse(content string) ([]entry, error) {
	var entries []entry
	for _, line := range strings.Split(content, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		idx := strings.IndexByte(line, '=')
		if idx < 0 {
			continue
		}
		key := strings.TrimSpace(line[:idx])
		value := strings.TrimSpace(line[idx+1:])
		value = strings.Trim(value, `"`)
		entries = append(entries, entry{key: key, value: value})
	}
	return entries, nil
}

func filter(vars []entry, keys []string) []entry {
	set := make(map[string]struct{}, len(keys))
	for _, k := range keys {
		set[k] = struct{}{}
	}
	var out []entry
	for _, e := range vars {
		if _, ok := set[e.key]; ok {
			out = append(out, e)
		}
	}
	return out
}
