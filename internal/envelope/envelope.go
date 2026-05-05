package envelope

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

// Envelope wraps encrypted .env content with metadata.
type Envelope struct {
	Version   int       `json:"version"`
	CreatedAt time.Time `json:"created_at"`
	Recipients []string `json:"recipients"`
	Ciphertext string   `json:"ciphertext"` // base64-encoded age ciphertext
}

// New creates a new Envelope with the given recipients and ciphertext.
func New(recipients []string, ciphertext string) *Envelope {
	return &Envelope{
		Version:    1,
		CreatedAt:  time.Now().UTC(),
		Recipients: recipients,
		Ciphertext: ciphertext,
	}
}

// Save writes the envelope as JSON to the given file path.
func Save(path string, env *Envelope) error {
	data, err := json.MarshalIndent(env, "", "  ")
	if err != nil {
		return fmt.Errorf("envelope: marshal failed: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("envelope: write failed: %w", err)
	}
	return nil
}

// Load reads and parses an envelope from the given file path.
func Load(path string) (*Envelope, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("envelope: read failed: %w", err)
	}
	var env Envelope
	if err := json.Unmarshal(data, &env); err != nil {
		return nil, fmt.Errorf("envelope: parse failed: %w", err)
	}
	if env.Version != 1 {
		return nil, fmt.Errorf("envelope: unsupported version %d", env.Version)
	}
	return &env, nil
}
