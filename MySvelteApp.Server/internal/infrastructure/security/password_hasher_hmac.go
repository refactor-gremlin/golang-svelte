package security

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"errors"
	"fmt"
)

const defaultSaltSize = 64

// HMACPasswordHasher reproduces the HMACSHA512 approach used in the previous .NET implementation.
type HMACPasswordHasher struct {
	saltSize int
}

// NewHMACPasswordHasher constructs a hasher with the default salt size.
func NewHMACPasswordHasher() *HMACPasswordHasher {
	return &HMACPasswordHasher{saltSize: defaultSaltSize}
}

// HashPassword generates a base64-encoded hash and salt.
func (h *HMACPasswordHasher) HashPassword(password string) (string, string, error) {
	if password == "" {
		return "", "", errors.New("password cannot be empty")
	}

	salt := make([]byte, h.saltSize)
	if _, err := rand.Read(salt); err != nil {
		return "", "", fmt.Errorf("generate salt: %w", err)
	}

	mac := hmac.New(sha512.New, salt)
	if _, err := mac.Write([]byte(password)); err != nil {
		return "", "", fmt.Errorf("compute hash: %w", err)
	}

	hash := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(hash), base64.StdEncoding.EncodeToString(salt), nil
}

// VerifyPassword recomputes the hash using the stored salt and compares it to the stored hash.
func (h *HMACPasswordHasher) VerifyPassword(password, storedHash, storedSalt string) (bool, error) {
	if password == "" {
		return false, errors.New("password cannot be empty")
	}
	if storedHash == "" || storedSalt == "" {
		return false, errors.New("stored hash and salt must be provided")
	}

	decodedSalt, err := base64.StdEncoding.DecodeString(storedSalt)
	if err != nil {
		return false, fmt.Errorf("decode salt: %w", err)
	}

	decodedHash, err := base64.StdEncoding.DecodeString(storedHash)
	if err != nil {
		return false, fmt.Errorf("decode hash: %w", err)
	}

	mac := hmac.New(sha512.New, decodedSalt)
	if _, err := mac.Write([]byte(password)); err != nil {
		return false, fmt.Errorf("compute hash: %w", err)
	}

	computed := mac.Sum(nil)

	return hmac.Equal(computed, decodedHash), nil
}
