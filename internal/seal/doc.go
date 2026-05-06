// Package seal provides the high-level Seal and Unseal operations for envseal.
//
// Seal reads a plaintext .env file, encrypts its contents using age encryption
// for all configured team recipients, and writes a versioned envelope to disk.
// The resulting .sealed file is safe to commit to version control.
//
// Unseal performs the reverse: it reads a .sealed envelope, decrypts it using
// the caller's local age identity, and writes the plaintext .env file.
//
// Typical usage:
//
//	err := seal.Seal(seal.SealOptions{
//		EnvFile:  ".env",
//		TeamFile: ".envseal/team.json",
//	})
//
//	err := seal.Unseal(seal.UnsealOptions{
//		SealedFile: ".env.sealed",
//		OutputFile: ".env",
//	})
package seal
