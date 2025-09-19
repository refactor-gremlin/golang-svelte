package app

import (
	"context"

	authdomain "mysvelteapp/server_new/internal/modules/auth/domain"
)

// UserRepository exposes persistence operations required by the auth use-cases.
type UserRepository interface {
	Add(ctx context.Context, user *authdomain.User) error
	GetByUsername(ctx context.Context, username string) (*authdomain.User, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}

// PasswordHasher hashes and verifies passwords.
type PasswordHasher interface {
	HashPassword(password string) (hash string, salt string, err error)
	VerifyPassword(password, hash, salt string) (bool, error)
}

// TokenGenerator issues access tokens for authenticated users.
type TokenGenerator interface {
	GenerateToken(user *authdomain.User) (string, error)
}
