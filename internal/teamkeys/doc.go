// Package teamkeys manages the set of age public keys for team members
// who are authorized to decrypt .env files sealed with envseal.
//
// A team keys file (.envseal-team.json) is a JSON document that maps
// human-readable names to age X25519 public keys. It is intended to be
// committed to version control alongside sealed .env files so that any
// authorized team member can be identified and re-encryption can be
// performed when the team changes.
//
// Example usage:
//
//	tk, err := teamkeys.Load(".envseal-team.json")
//	if err != nil {
//		// handle error
//	}
//	_ = tk.Add("alice", "age1...")
//	recipients, _ := tk.Recipients()
//	// pass recipients to crypto.Encrypt
package teamkeys
