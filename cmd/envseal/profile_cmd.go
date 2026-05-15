package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/yourorg/envseal/internal/envprofile"
)

func runProfile(args []string) error {
	fs := flag.NewFlagSet("profile", flag.ContinueOnError)
	file := fs.String("file", envprofile.DefaultProfilesFile, "profiles file path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	sub := fs.Arg(0)
	rest := fs.Args()
	if len(rest) > 0 {
		rest = rest[1:]
	}

	s, err := envprofile.Load(*file)
	if err != nil {
		return err
	}

	switch sub {
	case "set":
		return profileSet(s, *file, rest)
	case "switch":
		return profileSwitch(s, *file, rest)
	case "list":
		return profileList(s)
	case "active":
		return profileActive(s)
	default:
		return fmt.Errorf("profile: unknown subcommand %q (set|switch|list|active)", sub)
	}
}

func profileSet(s *envprofile.Store, file string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("profile set: usage: profile set <name> KEY=VALUE ...")
	}
	name := args[0]
	vars := make(map[string]string)
	for _, kv := range args[1:] {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			return fmt.Errorf("profile set: invalid key=value pair: %q", kv)
		}
		vars[parts[0]] = parts[1]
	}
	if err := envprofile.Set(s, name, vars); err != nil {
		return err
	}
	if err := envprofile.Save(file, s); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "profile %q saved with %d var(s)\n", name, len(vars))
	return nil
}

func profileSwitch(s *envprofile.Store, file string, args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("profile switch: usage: profile switch <name>")
	}
	if err := envprofile.Switch(s, args[0]); err != nil {
		return err
	}
	if err := envprofile.Save(file, s); err != nil {
		return err
	}
	fmt.Fprintf(os.Stdout, "active profile set to %q\n", args[0])
	return nil
}

func profileList(s *envprofile.Store) error {
	for name, p := range s.Profiles {
		marker := " "
		if name == s.Active {
			marker = "*"
		}
		fmt.Fprintf(os.Stdout, "%s %s (%d vars)\n", marker, name, len(p.Vars))
	}
	return nil
}

func profileActive(s *envprofile.Store) error {
	vars := envprofile.ActiveVars(s)
	if s.Active == "" {
		fmt.Fprintln(os.Stdout, "no active profile")
		return nil
	}
	fmt.Fprintf(os.Stdout, "active: %s\n", s.Active)
	out, _ := json.MarshalIndent(vars, "", "  ")
	fmt.Fprintln(os.Stdout, string(out))
	return nil
}
