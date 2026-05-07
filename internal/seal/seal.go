// Package seal provides Seal and Unseal operations for .env files,
// encrypting them into a versioned envelope using age recipients.
package seal

import (
	"fmt"
	"os"
	"time"

	"filippo.io/age"

	"github.com/yourorg/envseal/internal/audit"
	"github.com/yourorg/envseal/internal/crypto"
	"github.com/yourorg/envseal/internal/diff"
	"github.com/yourorg/envseal/internal/envelope"
	"github.com/yourorg/envseal/internal/version"
)

// Seal reads the plaintext .env file at envPath, encrypts it for all
// recipients, and writes a versioned sealed envelope to sealedPath.
// If a previous sealed file exists its version is bumped; otherwise
// a new envelope is created at version 1.
func Seal(envPath, sealedPath string, recipients []age.Recipient, identity age.Identity) error {
	plaintext, err := os.ReadFile(envPath)
	if err != nil {
		return fmt.Errorf("seal: read env file: %w", err)
	}

	ciphertext, err := crypto.Encrypt(plaintext, recipients)
	if err != nil {
		return fmt.Errorf("seal: encrypt: %w", err)
	}

	var env *envelope.Envelope

	existing, loadErr := envelope.Load(sealedPath)
	if loadErr == nil {
		// Decrypt existing to diff against new plaintext.
		oldPlain, decErr := crypto.Decrypt(existing.Ciphertext, []age.Identity{identity})
		if decErr == nil {
			changes := diff.Compare(string(oldPlain), string(plaintext))
			if len(changes) == 0 {
				return nil // nothing changed — skip re-seal
			}
			_ = audit.Log(sealedPath+".audit", audit.Entry{
				Action:  "seal",
				Message: diff.Summary(changes),
			})
		}
		v := version.Bump(existing.Version)
		env = &envelope.Envelope{
			Version:    v,
			Ciphertext: ciphertext,
		}
	} else {
		v := version.New(version.Checksum(plaintext), time.Now())
		env = &envelope.Envelope{
			Version:    v,
			Ciphertext: ciphertext,
		}
		_ = audit.Log(sealedPath+".audit", audit.Entry{
			Action:  "seal",
			Message: "initial seal",
		})
	}

	if err := envelope.Save(sealedPath, env); err != nil {
		return fmt.Errorf("seal: save envelope: %w", err)
	}
	return nil
}

// Unseal reads the sealed envelope at sealedPath, decrypts it using the
// provided identity, and writes the plaintext to envPath.
func Unseal(sealedPath, envPath string, identity age.Identity) error {
	env, err := envelope.Load(sealedPath)
	if err != nil {
		return fmt.Errorf("unseal: load envelope: %w", err)
	}

	plaintext, err := crypto.Decrypt(env.Ciphertext, []age.Identity{identity})
	if err != nil {
		return fmt.Errorf("unseal: decrypt: %w", err)
	}

	if err := os.WriteFile(envPath, plaintext, 0600); err != nil {
		return fmt.Errorf("unseal: write env file: %w", err)
	}

	_ = audit.Log(sealedPath+".audit", audit.Entry{
		Action:  "unseal",
		Message: fmt.Sprintf("unsealed version %d", env.Version.Number),
	})

	return nil
}
