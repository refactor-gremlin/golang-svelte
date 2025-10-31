package app

// Package ports defines the dependency interfaces for the authentication service.
// These interfaces enable dependency injection and loose coupling between the
// service layer and infrastructure implementations.
//
// Interface Design Principles:
// - Define interfaces in the consuming layer (app layer)
// - Keep interfaces focused and small
// - Return domain entities, not infrastructure types
// - Use context for cancellation and timeout handling
// - Handle errors explicitly, don't panic

import (
	"context"

	authdomain "mysvelteapp/server_new/internal/modules/auth/domain"
)

// UserRepository exposes persistence operations required by the auth use-cases.
// This interface defines the contract for user data storage and retrieval,
// abstracting away the specific persistence mechanism (database, file, etc.).
type UserRepository interface {
	// Add persists a new user to the storage system.
	// Returns an error if the user cannot be stored (duplicate key, constraint violation, etc.)
	Add(ctx context.Context, user *authdomain.User) error

	// GetByUsername retrieves a user by their unique username.
	// Returns nil if the user is not found (no error should be returned for not found)
	GetByUsername(ctx context.Context, username string) (*authdomain.User, error)

	// UsernameExists checks if a username is already taken.
	// Used for validation during registration to prevent duplicate usernames.
	UsernameExists(ctx context.Context, username string) (bool, error)

	// EmailExists checks if an email address is already registered.
	// Used for validation during registration to prevent duplicate accounts.
	EmailExists(ctx context.Context, email string) (bool, error)
}

// PasswordHasher handles secure password hashing and verification.
// This interface abstracts the specific hashing algorithm and parameters,
// allowing for easy upgrades to stronger algorithms in the future.
type PasswordHasher interface {
	// HashPassword creates a secure hash and salt for the given password.
	// Returns both hash and salt for storage; never returns empty values for valid passwords.
	HashPassword(password string) (hash string, salt string, err error)

	// VerifyPassword checks if the provided password matches the stored hash and salt.
	// Uses constant-time comparison to prevent timing attacks.
	// Returns true only if the password is correct; errors indicate system failures.
	VerifyPassword(password, hash, salt string) (bool, error)
}

// TokenGenerator creates authentication tokens for validated users.
// This interface abstracts the token format and signing mechanism,
// supporting JWT, opaque tokens, or future token formats.
type TokenGenerator interface {
	// GenerateToken creates a new authentication token for the given user.
	// The token should contain sufficient information for subsequent authentication.
	// Returns an error if token generation fails (signing key issues, configuration errors, etc.)
	GenerateToken(user *authdomain.User) (string, error)
}
