package domain_test

import (
	"testing"

	"mysvelteapp/server/internal/domain/entities"
)

func TestNewUserValid(t *testing.T) {
	user, err := entities.NewUser("test_user", "user@example.com", "hash", "salt")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if user.Username != "test_user" {
		t.Fatalf("expected username to be preserved, got %s", user.Username)
	}
	if user.Email != "user@example.com" {
		t.Fatalf("expected email to be normalised, got %s", user.Email)
	}
}

func TestNewUserEmptyUsername(t *testing.T) {
	_, err := entities.NewUser(" ", "user@example.com", "hash", "salt")
	if err == nil {
		t.Fatalf("expected error for empty username")
	}
}

func TestNewUserEmptyEmail(t *testing.T) {
	_, err := entities.NewUser("test_user", " ", "hash", "salt")
	if err == nil {
		t.Fatalf("expected error for empty email")
	}
}

func TestNewUserEmptyHash(t *testing.T) {
	_, err := entities.NewUser("test_user", "user@example.com", "", "salt")
	if err == nil {
		t.Fatalf("expected error for empty hash")
	}
}

func TestNewUserEmptySalt(t *testing.T) {
	_, err := entities.NewUser("test_user", "user@example.com", "hash", "")
	if err == nil {
		t.Fatalf("expected error for empty salt")
	}
}
