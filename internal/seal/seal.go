// Package seal provides high-level operations for encrypting and decrypting
// .env files using age encryption, combining envelope metadata and team keys.
package seal

import (
	"fmt"
	"os"
	"time"

	"github.com/user/envseal/internal/crypto"
	"github.com/user/envseal/internal/envelope"
	"github.com/user/envseal/internal/keystore"
	"github.com/user/envseal/internal/teamkeys"
)

// SealOptions configures the Seal operation.
type SealOptions struct {
	EnvFile    string
	OutputFile string
	TeamFile   string
	KeyFile    string
}

// UnsealOptions configures the Unseal operation.
type UnsealOptions struct {
	SealedFile string
	OutputFile string
	KeyFile    string
}

// Seal reads a plaintext .env file, encrypts it for all team recipients,
// and writes the sealed envelope to disk.
func Seal(opts SealOptions) error {
	plaintext, err := os.ReadFile(opts.EnvFile)
	if err != nil {
		return fmt.Errorf("reading env file: %w", err)
	}

	tm, err := teamkeys.Load(opts.TeamFile)
	if err != nil {
		return fmt.Errorf("loading team keys: %w", err)
	}

	recipients, err := tm.Recipients()
	if err != nil {
		return fmt.Errorf("parsing recipients: %w", err)
	}

	ciphertext, err := crypto.Encrypt(plaintext, recipients)
	if err != nil {
		return fmt.Errorf("encrypting: %w", err)
	}

	env := envelope.New(ciphertext)
	env.SealedAt = time.Now().UTC()
	env.Source = opts.EnvFile

	outPath := opts.OutputFile
	if outPath == "" {
		outPath = opts.EnvFile + ".sealed"
	}

	if err := envelope.Save(outPath, env); err != nil {
		return fmt.Errorf("saving envelope: %w", err)
	}

	return nil
}

// Unseal reads a sealed envelope, decrypts it using the local identity,
// and writes the plaintext .env file to disk.
func Unseal(opts UnsealOptions) error {
	env, err := envelope.Load(opts.SealedFile)
	if err != nil {
		return fmt.Errorf("loading envelope: %w", err)
	}

	keyPath := opts.KeyFile
	if keyPath == "" {
		keyPath = keystore.DefaultKeyPath()
	}

	identity, err := keystore.LoadIdentity(keyPath)
	if err != nil {
		return fmt.Errorf("loading identity: %w", err)
	}

	plaintext, err := crypto.Decrypt(env.Ciphertext, []interface{}{identity})
	if err != nil {
		return fmt.Errorf("decrypting: %w", err)
	}

	outPath := opts.OutputFile
	if outPath == "" {
		outPath = opts.SealedFile + ".env"
	}

	if err := os.WriteFile(outPath, plaintext, 0600); err != nil {
		return fmt.Errorf("writing output file: %w", err)
	}

	return nil
}
