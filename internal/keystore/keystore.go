package keystore

import (
	"errors"
	"os"
	"path/filepath"

	"filippo.io/age"
)

const (
	DefaultKeyDir  = ".envseal"
	DefaultKeyFile = "identity.age"
)

// Identity wraps an age identity with metadata.
type Identity struct {
	Recipient age.Recipient
	Private   age.Identity
	PublicKey string
}

// GenerateIdentity creates a new age X25519 key pair.
func GenerateIdentity() (*Identity, error) {
	identity, err := age.GenerateX25519Identity()
	if err != nil {
		return nil, err
	}
	return &Identity{
		Recipient: identity.Recipient(),
		Private:   identity,
		PublicKey: identity.Recipient().String(),
	}, nil
}

// SaveIdentity writes the private key to disk at the given path.
func SaveIdentity(id *Identity, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(id.Private.(*age.X25519Identity).String() + "\n")
	return err
}

// LoadIdentity reads an age private key from disk.
func LoadIdentity(path string) (*Identity, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	identities, err := age.ParseIdentities(strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}
	if len(identities) == 0 {
		return nil, errors.New("no identities found in key file")
	}
	x, ok := identities[0].(*age.X25519Identity)
	if !ok {
		return nil, errors.New("unsupported identity type")
	}
	return &Identity{
		Recipient: x.Recipient(),
		Private:   x,
		PublicKey: x.Recipient().String(),
	}, nil
}

// DefaultKeyPath returns the default path for the user's identity key.
func DefaultKeyPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, DefaultKeyDir, DefaultKeyFile)
}
