package domain

import (
	"strings"
	"testing"
)

func TestValidUsername(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"empty", "", false},
		{"too short", "ab", false},
		{"valid", "user_123", true},
		{"invalid chars", "user-123", false},
		{"too long", strings.Repeat("a", MaxUsernameLength+1), false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := ValidUsername(tc.input); got != tc.valid {
				t.Fatalf("ValidUsername(%q) = %v, want %v", tc.input, got, tc.valid)
			}
		})
	}
}

func TestValidEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		input string
		valid bool
	}{
		{"empty", "", false},
		{"valid", "user@example.com", true},
		{"double dot", "user..test@example.com", false},
		{"missing domain", "user@", false},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := ValidEmail(tc.input); got != tc.valid {
				t.Fatalf("ValidEmail(%q) = %v, want %v", tc.input, got, tc.valid)
			}
		})
	}
}

func TestNewUserRejectsInvalidUsername(t *testing.T) {
	t.Parallel()

	if _, err := NewUser("ab", "test@example.com", "hash", "salt"); err == nil {
		t.Fatalf("expected error for short username")
	}

	if _, err := NewUser("invalid-username", "test@example.com", "hash", "salt"); err == nil {
		t.Fatalf("expected error for invalid characters")
	}
}

func TestNewUserRejectsInvalidEmail(t *testing.T) {
	t.Parallel()

	if _, err := NewUser("valid_user", "", "hash", "salt"); err == nil {
		t.Fatalf("expected error for empty email")
	}

	if _, err := NewUser("valid_user", "user..test@example.com", "hash", "salt"); err == nil {
		t.Fatalf("expected error for double dot email")
	}
}
