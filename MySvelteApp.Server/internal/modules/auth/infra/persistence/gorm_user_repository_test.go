package persistence

import (
	"context"
	"errors"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	authapp "mysvelteapp/server_new/internal/modules/auth/app"
	authdomain "mysvelteapp/server_new/internal/modules/auth/domain"
)

func newTestRepository(t *testing.T) *GormUserRepository {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	if err := db.AutoMigrate(&authdomain.User{}); err != nil {
		t.Fatalf("migrate schema: %v", err)
	}

	return NewGormUserRepository(db)
}

func TestAddReturnsConflictForDuplicateUsername(t *testing.T) {
	t.Parallel()

	repo := newTestRepository(t)
	ctx := context.Background()

	first, err := authdomain.NewUser("duplicate", "first@example.com", "hash1", "salt1")
	if err != nil {
		t.Fatalf("NewUser: %v", err)
	}
	if err := repo.Add(ctx, first); err != nil {
		t.Fatalf("Add first: %v", err)
	}

	second, err := authdomain.NewUser("duplicate", "second@example.com", "hash2", "salt2")
	if err != nil {
		t.Fatalf("NewUser second: %v", err)
	}
	if err := repo.Add(ctx, second); err == nil {
		t.Fatalf("expected conflict error")
	} else {
		if !authapp.IsConflictError(err) {
			t.Fatalf("expected conflict error, got %v", err)
		}
		var conflict authapp.ConflictError
		if !errors.As(err, &conflict) {
			t.Fatalf("expected ConflictError type")
		}
		expected := "This username is already taken. Please choose a different one."
		if conflict.Message != expected {
			t.Fatalf("expected message %q, got %q", expected, conflict.Message)
		}
	}
}

func TestAddReturnsConflictForDuplicateEmail(t *testing.T) {
	t.Parallel()

	repo := newTestRepository(t)
	ctx := context.Background()

	first, err := authdomain.NewUser("first_user", "user@example.com", "hash1", "salt1")
	if err != nil {
		t.Fatalf("NewUser: %v", err)
	}
	if err := repo.Add(ctx, first); err != nil {
		t.Fatalf("Add first: %v", err)
	}

	second, err := authdomain.NewUser("second_user", "USER@example.com", "hash2", "salt2")
	if err != nil {
		t.Fatalf("NewUser second: %v", err)
	}
	if err := repo.Add(ctx, second); err == nil {
		t.Fatalf("expected conflict error")
	} else {
		if !authapp.IsConflictError(err) {
			t.Fatalf("expected conflict error, got %v", err)
		}
		var conflict authapp.ConflictError
		if !errors.As(err, &conflict) {
			t.Fatalf("expected ConflictError type")
		}
		expected := "This email is already registered. Please use a different email address."
		if conflict.Message != expected {
			t.Fatalf("expected message %q, got %q", expected, conflict.Message)
		}
	}
}
