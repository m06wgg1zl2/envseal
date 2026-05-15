// Package envprofile supports named environment profiles (e.g. dev, staging, prod)
// that can be switched, merged, and tracked within envseal.
package envprofile

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const DefaultProfilesFile = ".envseal-profiles.json"

// Profile represents a named set of environment variable overrides.
type Profile struct {
	Name      string            `json:"name"`
	Vars      map[string]string `json:"vars"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

// Store holds all profiles for a project.
type Store struct {
	Active   string             `json:"active"`
	Profiles map[string]Profile `json:"profiles"`
}

// Load reads the profiles store from disk. Returns an empty store if not found.
func Load(path string) (*Store, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return &Store{Profiles: make(map[string]Profile)}, nil
	}
	if err != nil {
		return nil, fmt.Errorf("envprofile: read %s: %w", path, err)
	}
	var s Store
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, fmt.Errorf("envprofile: parse %s: %w", path, err)
	}
	if s.Profiles == nil {
		s.Profiles = make(map[string]Profile)
	}
	return &s, nil
}

// Save writes the store to disk.
func Save(path string, s *Store) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return fmt.Errorf("envprofile: mkdir: %w", err)
	}
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return fmt.Errorf("envprofile: marshal: %w", err)
	}
	return os.WriteFile(path, data, 0o600)
}

// Set adds or updates a profile in the store.
func Set(s *Store, name string, vars map[string]string) error {
	if name == "" {
		return errors.New("envprofile: profile name must not be empty")
	}
	now := time.Now().UTC()
	existing, ok := s.Profiles[name]
	if !ok {
		existing = Profile{Name: name, CreatedAt: now}
	}
	existing.Vars = vars
	existing.UpdatedAt = now
	s.Profiles[name] = existing
	return nil
}

// Switch sets the active profile by name.
func Switch(s *Store, name string) error {
	if _, ok := s.Profiles[name]; !ok {
		return fmt.Errorf("envprofile: profile %q not found", name)
	}
	s.Active = name
	return nil
}

// ActiveVars returns the vars for the currently active profile, or an empty map.
func ActiveVars(s *Store) map[string]string {
	if s.Active == "" {
		return map[string]string{}
	}
	p, ok := s.Profiles[s.Active]
	if !ok {
		return map[string]string{}
	}
	out := make(map[string]string, len(p.Vars))
	for k, v := range p.Vars {
		out[k] = v
	}
	return out
}
