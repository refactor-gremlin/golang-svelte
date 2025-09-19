package security_test

import (
	"testing"

	authsecurity "mysvelteapp/server_new/internal/modules/auth/infra/security"
)

func TestHashAndVerifyPassword(t *testing.T) {
	hasher := authsecurity.NewHMACPasswordHasher()

	hash, salt, err := hasher.HashPassword("Password123")
	if err != nil {
		t.Fatalf("expected no error hashing password, got %v", err)
	}
	if hash == "" || salt == "" {
		t.Fatalf("expected hash and salt to be populated")
	}

	verified, err := hasher.VerifyPassword("Password123", hash, salt)
	if err != nil {
		t.Fatalf("expected no error verifying password, got %v", err)
	}
	if !verified {
		t.Fatalf("expected password to verify correctly")
	}

	verified, err = hasher.VerifyPassword("WrongPassword", hash, salt)
	if err != nil {
		t.Fatalf("expected no error verifying password, got %v", err)
	}
	if verified {
		t.Fatalf("expected verification to fail for incorrect password")
	}
}
