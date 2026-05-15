package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/your-org/envseal/internal/envclone"
)

func runClone(args []string) error {
	fs := flag.NewFlagSet("clone", flag.ContinueOnError)
	include := fs.String("include", "", "comma-separated list of keys to include")
	exclude := fs.String("exclude", "", "comma-separated list of keys to exclude")
	rename := fs.String("rename", "", "comma-separated OLD=NEW pairs to rename keys")
	overwrite := fs.Bool("overwrite", false, "overwrite destination if it already exists")

	if err := fs.Parse(args); err != nil {
		return err
	}

	positional := fs.Args()
	if len(positional) != 2 {
		return fmt.Errorf("usage: envseal clone [flags] <src> <dst>")
	}
	src, dst := positional[0], positional[1]

	opts := envclone.Options{
		Overwrite: *overwrite,
	}

	if *include != "" {
		opts.IncludeKeys = splitComma(*include)
	}
	if *exclude != "" {
		opts.ExcludeKeys = splitComma(*exclude)
	}
	if *rename != "" {
		m, err := parseRenameMap(*rename)
		if err != nil {
			return err
		}
		opts.Rename = m
	}

	if err := envclone.Clone(src, dst, opts); err != nil {
		return err
	}

	fmt.Fprintf(os.Stdout, "cloned %s → %s\n", src, dst)
	return nil
}

func parseRenameMap(s string) (map[string]string, error) {
	m := make(map[string]string)
	for _, pair := range strings.Split(s, ",") {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("clone: invalid rename pair %q (expected OLD=NEW)", pair)
		}
		m[parts[0]] = parts[1]
	}
	return m, nil
}
