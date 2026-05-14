// Package envtemplate provides functionality to generate a .env.template file
// from a sealed .env file, stripping values but preserving keys and comments.
package envtemplate

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Entry represents a single line in the template output.
type Entry struct {
	Key     string
	Comment string
	Blank   bool
}

// Generate reads the source env file and writes a template to dest,
// replacing all values with empty strings while preserving keys and comments.
func Generate(src, dest string) error {
	in, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("envtemplate: open source: %w", err)
	}
	defer in.Close()

	entries, err := parse(in)
	if err != nil {
		return fmt.Errorf("envtemplate: parse: %w", err)
	}

	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("envtemplate: open dest: %w", err)
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	for _, e := range entries {
		switch {
		case e.Blank:
			fmt.Fprintln(w)
		case e.Comment != "":
			fmt.Fprintln(w, e.Comment)
		default:
			fmt.Fprintf(w, "%s=\n", e.Key)
		}
	}

	return w.Flush()
}

func parse(f *os.File) ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		switch {
		case trimmed == "":
			entries = append(entries, Entry{Blank: true})
		case strings.HasPrefix(trimmed, "#"):
			entries = append(entries, Entry{Comment: line})
		default:
			parts := strings.SplitN(trimmed, "=", 2)
			key := strings.TrimSpace(parts[0])
			if key == "" {
				continue
			}
			entries = append(entries, Entry{Key: key})
		}
	}
	return entries, scanner.Err()
}
