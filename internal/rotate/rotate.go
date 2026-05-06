// Package rotate provides functionality to re-encrypt sealed env files
// with a new or updated set of recipient keys.
package rotate

import (
	"fmt"

	"github.com/user/envseal/internal/crypto"
	"github.com/user/envseal/internal/envelope"
	"github.com/user/envseal/internal/keystore"
	"github.com/user/envseal/internal/seal"
	"github.com/user/envseal/internal/teamkeys"
)

// Options configures a key rotation operation.
type Options struct {
	// SealedFile is the path to the existing .sealed.json file.
	SealedFile string
	// TeamKeysFile is the path to the team keys JSON file.
	TeamKeysFile string
	// IdentityFile is the path to the caller's age private key.
	IdentityFile string
	// OutputFile is the path to write the rotated sealed file.
	// If empty, SealedFile is overwritten.
	OutputFile string
}

// Rotate decrypts the sealed env file using the provided identity, then
// re-encrypts it with the current set of team recipients.
func Rotate(opts Options) error {
	if opts.SealedFile == "" {
		return fmt.Errorf("rotate: sealed file path is required")
	}
	if opts.IdentityFile == "" {
		opts.IdentityFile = keystore.DefaultKeyPath()
	}
	if opts.OutputFile == "" {
		opts.OutputFile = opts.SealedFile
	}

	// Load existing envelope to get current version metadata.
	env, err := envelope.Load(opts.SealedFile)
	if err != nil {
		return fmt.Errorf("rotate: loading sealed file: %w", err)
	}

	// Load caller identity for decryption.
	identity, err := keystore.LoadIdentity(opts.IdentityFile)
	if err != nil {
		return fmt.Errorf("rotate: loading identity: %w", err)
	}

	// Decrypt the current ciphertext.
	plaintext, err := crypto.Decrypt(env.Ciphertext, identity)
	if err != nil {
		return fmt.Errorf("rotate: decrypting envelope: %w", err)
	}

	// Load team keys to build new recipient list.
	tk, err := teamkeys.Load(opts.TeamKeysFile)
	if err != nil {
		return fmt.Errorf("rotate: loading team keys: %w", err)
	}

	recipients, err := tk.Recipients()
	if err != nil {
		return fmt.Errorf("rotate: resolving recipients: %w", err)
	}
	if len(recipients) == 0 {
		return fmt.Errorf("rotate: no recipients found in team keys file")
	}

	// Re-encrypt with updated recipients.
	newCiphertext, err := crypto.Encrypt(plaintext, recipients)
	if err != nil {
		return fmt.Errorf("rotate: re-encrypting: %w", err)
	}

	// Preserve version metadata, update ciphertext.
	env.Ciphertext = newCiphertext

	// Write rotated envelope.
	if err := seal.WriteEnvelope(opts.OutputFile, env); err != nil {
		return fmt.Errorf("rotate: writing rotated file: %w", err)
	}

	return nil
}
