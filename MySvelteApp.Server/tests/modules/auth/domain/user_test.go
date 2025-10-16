package domain_test

import (
	"strings"
	"testing"

	authdomain "mysvelteapp/server_new/internal/modules/auth/domain"
)

// TestNewUserValid confirms canonical inputs construct a valid user.
// Arrange: prepare a typical set of field values.
// Act: call NewUser with those values.
// Assert: expect trimmed username and normalised email.
func TestNewUserValid(t *testing.T) {
	// Arrange
	username := "test_user"
	email := "user@example.com"

	// Act
	user, err := authdomain.NewUser(username, email, "hash", "salt")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Username != username {
		t.Fatalf("expected username to be preserved, got %s", user.Username)
	}
	if user.Email != email {
		t.Fatalf("expected email to be normalised, got %s", user.Email)
	}
}

// TestNewUserTrimsAndNormalizesFields checks trimming and lowercasing logic.
// Arrange: provide padded username and mixed-case email.
// Act: call NewUser with the noisy inputs.
// Assert: expect trimmed username and lowercase email in the result.
func TestNewUserTrimsAndNormalizesFields(t *testing.T) {
	// Arrange
	username := "  spaced  "
	email := " MixedCase@Example.COM "

	// Act
	user, err := authdomain.NewUser(username, email, "hash", "salt")

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Username != "spaced" {
		t.Fatalf("expected username to be trimmed, got %q", user.Username)
	}
	if user.Email != "mixedcase@example.com" {
		t.Fatalf("expected email to be lowercased, got %q", user.Email)
	}
}

// TestNewUserEmptyUsername enforces username presence.
// Arrange: use whitespace-only username.
// Act: attempt to build a user.
// Assert: expect an error.
func TestNewUserEmptyUsername(t *testing.T) {
	// Arrange
	username := " "
	email := "user@example.com"

	// Act
	_, err := authdomain.NewUser(username, email, "hash", "salt")

	// Assert
	if err == nil {
		t.Fatalf("expected error for empty username")
	}
}

// TestNewUserEmptyEmail enforces email presence.
// Arrange: set the email to whitespace.
// Act: invoke NewUser.
// Assert: expect an error response.
func TestNewUserEmptyEmail(t *testing.T) {
	// Arrange
	username := "test_user"
	email := " "

	// Act
	_, err := authdomain.NewUser(username, email, "hash", "salt")

	// Assert
	if err == nil {
		t.Fatalf("expected error for empty email")
	}
}

// TestNewUserEmptyHash prevents empty password hashes.
// Arrange: pass an empty hash.
// Act: call NewUser.
// Assert: expect an error.
func TestNewUserEmptyHash(t *testing.T) {
	// Arrange
	username := "test_user"
	email := "user@example.com"

	// Act
	_, err := authdomain.NewUser(username, email, "", "salt")

	// Assert
	if err == nil {
		t.Fatalf("expected error for empty hash")
	}
}

// TestNewUserEmptySalt prevents empty password salts.
// Arrange: set the salt to empty.
// Act: call NewUser.
// Assert: expect an error.
func TestNewUserEmptySalt(t *testing.T) {
	// Arrange
	username := "test_user"
	email := "user@example.com"

	// Act
	_, err := authdomain.NewUser(username, email, "hash", "")

	// Assert
	if err == nil {
		t.Fatalf("expected error for empty salt")
	}
}

// TestNewUserUsernameTooLong guards the maximum username length.
// Arrange: craft a username beyond MaxUsernameLength.
// Act: call NewUser with the long username.
// Assert: expect an error describing the limit breach.
func TestNewUserUsernameTooLong(t *testing.T) {
	// Arrange
	username := strings.Repeat("a", authdomain.MaxUsernameLength+1)
	email := "user@example.com"

	// Act
	_, err := authdomain.NewUser(username, email, "hash", "salt")

	// Assert
	if err == nil {
		t.Fatalf("expected error for long username")
	}
}

// TestNewUserEmailTooLong guards the maximum email length.
// Arrange: compose an email longer than MaxEmailLength.
// Act: attempt to create the user.
// Assert: expect an error.
func TestNewUserEmailTooLong(t *testing.T) {
	// Arrange
	tooLongEmail := strings.Repeat("a", authdomain.MaxEmailLength-10) + "@example.com"
	if len(tooLongEmail) <= authdomain.MaxEmailLength {
		t.Fatalf("constructed email is not longer than max")
	}

	// Act
	_, err := authdomain.NewUser("test_user", tooLongEmail, "hash", "salt")

	// Assert
	if err == nil {
		t.Fatalf("expected error for long email")
	}
}
