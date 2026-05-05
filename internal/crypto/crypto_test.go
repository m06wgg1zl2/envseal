package crypto_test

import (
	"testing"

	"filippo.io/age"

	"github.com/yourusername/envseal/internal/crypto"
)

func generateTestIdentity(t *testing.T) *age.X25519Identity {
	t.Helper()
	id, err := age.GenerateX25519Identity()
	if err != nil {
		t.Fatalf("failed to generate identity: %v", err)
	}
	return id
}

func TestEncryptDecrypt_RoundTrip(t *testing.T) {
	id := generateTestIdentity(t)
	plaintext := []byte("SECRET=supersecret\nDB_PASS=hunter2")

	ciphertext, err := crypto.Encrypt(plaintext, []age.Recipient{id.Recipient()})
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	if len(ciphertext) == 0 {
		t.Fatal("expected non-empty ciphertext")
	}

	decrypted, err := crypto.Decrypt(ciphertext, []age.Identity{id})
	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	if string(decrypted) != string(plaintext) {
		t.Errorf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncrypt_NoRecipients(t *testing.T) {
	_, err := crypto.Encrypt([]byte("data"), nil)
	if err == nil {
		t.Fatal("expected error with no recipients, got nil")
	}
}

func TestDecrypt_NoIdentities(t *testing.T) {
	id := generateTestIdentity(t)
	ciphertext, err := crypto.Encrypt([]byte("data"), []age.Recipient{id.Recipient()})
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = crypto.Decrypt(ciphertext, nil)
	if err == nil {
		t.Fatal("expected error with no identities, got nil")
	}
}

func TestDecrypt_WrongIdentity(t *testing.T) {
	id1 := generateTestIdentity(t)
	id2 := generateTestIdentity(t)

	ciphertext, err := crypto.Encrypt([]byte("secret"), []age.Recipient{id1.Recipient()})
	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	_, err = crypto.Decrypt(ciphertext, []age.Identity{id2})
	if err == nil {
		t.Fatal("expected error decrypting with wrong identity, got nil")
	}
}

func TestParseRecipient_Valid(t *testing.T) {
	id := generateTestIdentity(t)
	pubkey := id.Recipient().String()

	recipient, err := crypto.ParseRecipient(pubkey)
	if err != nil {
		t.Fatalf("ParseRecipient failed: %v", err)
	}
	if recipient == nil {
		t.Fatal("expected non-nil recipient")
	}
}

func TestParseRecipient_Invalid(t *testing.T) {
	_, err := crypto.ParseRecipient("not-a-valid-key")
	if err == nil {
		t.Fatal("expected error for invalid key, got nil")
	}
}
