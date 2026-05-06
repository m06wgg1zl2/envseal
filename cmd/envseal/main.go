package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/user/envseal/internal/keystore"
	"github.com/user/envseal/internal/seal"
	"github.com/user/envseal/internal/teamkeys"
)

const usage = `envseal — encrypt and version-control .env files

Usage:
  envseal <command> [flags]

Commands:
  init      Generate a new age identity
  seal      Encrypt a .env file
  unseal    Decrypt a .env.sealed file
  add-key   Add a team member's public key
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		runInit(os.Args[2:])
	case "seal":
		runSeal(os.Args[2:])
	case "unseal":
		runUnseal(os.Args[2:])
	case "add-key":
		runAddKey(os.Args[2:])
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		fmt.Fprint(os.Stderr, usage)
		os.Exit(1)
	}
}

func runInit(args []string) {
	fs := flag.NewFlagSet("init", flag.ExitOnError)
	keyFile := fs.String("key", keystore.DefaultKeyPath(), "path to write identity")
	_ = fs.Parse(args)

	identity, err := keystore.GenerateIdentity()
	dieIf(err)
	dieIf(keystore.SaveIdentity(*keyFile, identity))
	fmt.Printf("Identity generated: %s\nPublic key: %s\n", *keyFile, identity.Recipient())
}

func runSeal(args []string) {
	fs := flag.NewFlagSet("seal", flag.ExitOnError)
	envFile := fs.String("env", ".env", "plaintext .env file to seal")
	outFile := fs.String("out", "", "output sealed file (default: <env>.sealed)")
	teamFile := fs.String("team", ".envseal/team.json", "team keys file")
	keyFile := fs.String("key", keystore.DefaultKeyPath(), "identity key file")
	_ = fs.Parse(args)

	dieIf(seal.Seal(seal.SealOptions{
		EnvFile:    *envFile,
		OutputFile: *outFile,
		TeamFile:   *teamFile,
		KeyFile:    *keyFile,
	}))
	fmt.Println("Sealed successfully.")
}

func runUnseal(args []string) {
	fs := flag.NewFlagSet("unseal", flag.ExitOnError)
	sealedFile := fs.String("in", ".env.sealed", "sealed file to decrypt")
	outFile := fs.String("out", ".env", "output plaintext file")
	keyFile := fs.String("key", keystore.DefaultKeyPath(), "identity key file")
	_ = fs.Parse(args)

	dieIf(seal.Unseal(seal.UnsealOptions{
		SealedFile: *sealedFile,
		OutputFile: *outFile,
		KeyFile:    *keyFile,
	}))
	fmt.Println("Unsealed successfully.")
}

func runAddKey(args []string) {
	fs := flag.NewFlagSet("add-key", flag.ExitOnError)
	teamFile := fs.String("team", ".envseal/team.json", "team keys file")
	name := fs.String("name", "", "team member name (required)")
	pubKey := fs.String("pubkey", "", "age public key (required)")
	_ = fs.Parse(args)

	if *name == "" || *pubKey == "" {
		fmt.Fprintln(os.Stderr, "--name and --pubkey are required")
		os.Exit(1)
	}

	tm, err := teamkeys.Load(*teamFile)
	if os.IsNotExist(err) {
		tm = teamkeys.New()
	} else {
		dieIf(err)
	}

	dieIf(tm.Add(*name, *pubKey))
	dieIf(teamkeys.Save(*teamFile, tm))
	fmt.Printf("Added key for %s\n", *name)
}

func dieIf(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
