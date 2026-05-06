package teamkeys

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"
)

const DefaultTeamKeysFile = ".envseal-team.json"

// TeamKeys holds a map of named public keys for team members.
type TeamKeys struct {
	Keys map[string]string `json:"keys"`
}

// New returns an empty TeamKeys instance.
func New() *TeamKeys {
	return &TeamKeys{Keys: make(map[string]string)}
}

// Add adds or updates a named recipient public key.
func (t *TeamKeys) Add(name, publicKey string) error {
	if name == "" {
		return errors.New("name must not be empty")
	}
	if _, err := age.ParseX25519Recipient(publicKey); err != nil {
		return fmt.Errorf("invalid public key for %q: %w", name, err)
	}
	t.Keys[name] = publicKey
	return nil
}

// Remove removes a named recipient by name.
func (t *TeamKeys) Remove(name string) error {
	if _, ok := t.Keys[name]; !ok {
		return fmt.Errorf("key %q not found", name)
	}
	delete(t.Keys, name)
	return nil
}

// Recipients returns a list of age.Recipient for all team members.
func (t *TeamKeys) Recipients() ([]age.Recipient, error) {
	var recipients []age.Recipient
	for name, pub := range t.Keys {
		r, err := age.ParseX25519Recipient(pub)
		if err != nil {
			return nil, fmt.Errorf("invalid key for %q: %w", name, err)
		}
		recipients = append(recipients, r)
	}
	return recipients, nil
}

// Save writes the TeamKeys to a JSON file at the given path.
func Save(path string, t *TeamKeys) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("creating directory: %w", err)
	}
	data, err := json.MarshalIndent(t, "", "  ")
	if err != nil {
		return fmt.Errorf("marshalling team keys: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

// Load reads TeamKeys from a JSON file at the given path.
func Load(path string) (*TeamKeys, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading team keys file: %w", err)
	}
	var t TeamKeys
	if err := json.Unmarshal(data, &t); err != nil {
		return nil, fmt.Errorf("parsing team keys file: %w", err)
	}
	if t.Keys == nil {
		t.Keys = make(map[string]string)
	}
	return &t, nil
}
