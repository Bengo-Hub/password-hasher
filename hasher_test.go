package passwordhasher

import (
	"crypto/rand"
	"strings"
	"testing"
)

func TestHasher_Hash(t *testing.T) {
	hasher := NewHasher()

	password := "ChangeMe123!"
	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	// Verify format: $argon2id$v=19$m=65536,t=3,p=2$<salt>$<hash>
	if !strings.HasPrefix(hash, "$argon2id$v=19$") {
		t.Errorf("Invalid hash format: %s", hash)
	}

	parts := strings.Split(hash, "$")
	if len(parts) != 6 {
		t.Errorf("Expected 6 parts, got %d", len(parts))
	}

	// Verify parameters
	if parts[3] != "m=65536,t=3,p=2" {
		t.Errorf("Invalid parameters: %s", parts[3])
	}
}

func TestHasher_Verify_Success(t *testing.T) {
	hasher := NewHasher()

	password := "ChangeMe123!"
	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	// Verify correct password
	err = hasher.Verify(password, hash)
	if err != nil {
		t.Errorf("Verify failed for correct password: %v", err)
	}
}

func TestHasher_Verify_Failure(t *testing.T) {
	hasher := NewHasher()

	password := "ChangeMe123!"
	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	// Verify wrong password
	err = hasher.Verify("WrongPassword", hash)
	if err != ErrPasswordMismatch {
		t.Errorf("Expected ErrPasswordMismatch, got: %v", err)
	}
}

func TestHasher_Verify_InvalidFormat(t *testing.T) {
	hasher := NewHasher()

	tests := []struct {
		name string
		hash string
	}{
		{"empty", ""},
		{"invalid format", "invalid-hash"},
		{"wrong algorithm", "$bcrypt$v=19$m=65536,t=3,p=2$salt$hash"},
		{"wrong version", "$argon2id$v=18$m=65536,t=3,p=2$salt$hash"},
		{"missing parts", "$argon2id$v=19$m=65536,t=3,p=2$salt"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := hasher.Verify("password", tt.hash)
			if err == nil {
				t.Error("Expected error for invalid hash")
			}
		})
	}
}

func TestHasher_HashWithSalt(t *testing.T) {
	hasher := NewHasher()

	password := "ChangeMe123!"
	salt := make([]byte, DefaultSaltLength)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("Failed to generate salt: %v", err)
	}

	hash1, err := hasher.HashWithSalt(password, salt)
	if err != nil {
		t.Fatalf("HashWithSalt failed: %v", err)
	}

	hash2, err := hasher.HashWithSalt(password, salt)
	if err != nil {
		t.Fatalf("HashWithSalt failed: %v", err)
	}

	// Same salt should produce same hash
	if hash1 != hash2 {
		t.Error("Same salt should produce identical hashes")
	}

	// Verify the hash
	err = hasher.Verify(password, hash1)
	if err != nil {
		t.Errorf("Verify failed: %v", err)
	}
}

func TestHasher_CustomParameters(t *testing.T) {
	hasher := NewCustomHasher(32768, 2, 1, 32)

	password := "ChangeMe123!"
	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	// Verify format includes custom parameters
	if !strings.Contains(hash, "m=32768,t=2,p=1") {
		t.Errorf("Hash doesn't contain custom parameters: %s", hash)
	}

	// Verify password
	err = hasher.Verify(password, hash)
	if err != nil {
		t.Errorf("Verify failed: %v", err)
	}
}

func TestHasher_EmptyPassword(t *testing.T) {
	hasher := NewHasher()

	_, err := hasher.Hash("")
	if err == nil {
		t.Error("Expected error for empty password")
	}

	err = hasher.Verify("", "some-hash")
	if err == nil {
		t.Error("Expected error for empty password in Verify")
	}
}

// TestAuthServiceCompatibility verifies compatibility with auth-service Go implementation
func TestAuthServiceCompatibility(t *testing.T) {
	hasher := NewHasher()

	// Test with the same password format used in auth-service
	password := "ChangeMe123!"

	// Generate hash
	hash, err := hasher.Hash(password)
	if err != nil {
		t.Fatalf("Hash failed: %v", err)
	}

	// Verify format matches auth-service
	expectedPrefix := "$argon2id$v=19$m=65536,t=3,p=2$"
	if !strings.HasPrefix(hash, expectedPrefix) {
		t.Errorf("Hash format doesn't match auth-service. Got: %s", hash)
	}

	// Verify we can validate the hash
	err = hasher.Verify(password, hash)
	if err != nil {
		t.Errorf("Failed to verify hash: %v", err)
	}

	// Verify wrong password fails
	err = hasher.Verify("WrongPassword", hash)
	if err != ErrPasswordMismatch {
		t.Errorf("Expected ErrPasswordMismatch for wrong password, got: %v", err)
	}
}

func BenchmarkHash(b *testing.B) {
	hasher := NewHasher()
	password := "ChangeMe123!"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = hasher.Hash(password)
	}
}

func BenchmarkVerify(b *testing.B) {
	hasher := NewHasher()
	password := "ChangeMe123!"
	hash, _ := hasher.Hash(password)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = hasher.Verify(password, hash)
	}
}
