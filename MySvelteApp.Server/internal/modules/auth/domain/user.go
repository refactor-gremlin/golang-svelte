package domain

// Package domain contains the core business entities for the authentication module.
// This layer defines the fundamental business concepts and invariants that must
// always be preserved, independent of any external concerns like persistence or API.
//
// Domain Design Principles:
// - Rich domain models with business logic
// - Enforced invariants through constructors
// - No external dependencies (frameworks, databases, etc.)
// - Pure Go with business rule validation

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"
)

const (
	// MinUsernameLength enforces the minimum username length allowed by the domain.
	MinUsernameLength = 3
	// MaxUsernameLength mirrors the legacy constraints.
	MaxUsernameLength = 64
	// MaxEmailLength mirrors the legacy constraints.
	MaxEmailLength = 320
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	emailRegex    = regexp.MustCompile(`^[^\s@]+@[^\s@.]+\.[^\s@.]+$`)
)

// ValidUsername reports whether the supplied username satisfies domain rules.
func ValidUsername(username string) bool {
	username = strings.TrimSpace(username)
	if len(username) < MinUsernameLength || len(username) > MaxUsernameLength {
		return false
	}
	return usernameRegex.MatchString(username)
}

// ValidEmail reports whether the supplied email satisfies domain rules.
func ValidEmail(email string) bool {
	email = strings.TrimSpace(email)
	if len(email) == 0 || len(email) > MaxEmailLength {
		return false
	}
	if strings.Contains(email, "..") {
		return false
	}
	return emailRegex.MatchString(email)
}

// User represents an authenticated user entity in the system.
// This aggregate root encapsulates the core user data and business invariants.
//
// Business Invariants:
// - Username and email must be unique within the system
// - Username must be 3-64 characters, alphanumeric + underscores
// - Email must be valid format and <= 320 characters
// - Password hash and salt must not be empty
// - All string fields are automatically trimmed and normalized
//
// GORM tags are used for persistence mapping but the domain
// logic is independent of any specific storage mechanism.
type User struct {
	ID           uint      `gorm:"primaryKey"`                    // Unique identifier
	Username     string    `gorm:"size:64;uniqueIndex;not null"`  // Unique username
	Email        string    `gorm:"size:320;uniqueIndex;not null"` // Normalized email
	PasswordHash string    `gorm:"size:512;not null"`             // Hashed password
	PasswordSalt string    `gorm:"size:256;not null"`             // Password salt
	CreatedAt    time.Time `gorm:"autoCreateTime"`                // Creation timestamp
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`                // Last update timestamp
}

// NewUser creates a new User aggregate after enforcing all business invariants.
// This constructor ensures that only valid User entities can be created,
// protecting the domain from invalid state.
//
// Validation Rules Applied:
// - Username: trimmed, non-empty, length check, alphanumeric validation
// - Email: trimmed, normalized to lowercase, format validation, length check
// - Password: non-empty hash and salt (actual validation happens in service layer)
//
// Returns an error if any invariant is violated, never returns an invalid User.
func NewUser(username, email, passwordHash, passwordSalt string) (*User, error) {
	username = strings.TrimSpace(username)
	if len(username) == 0 {
		return nil, errors.New("username cannot be empty")
	}
	if len(username) < MinUsernameLength {
		return nil, fmt.Errorf("username must be at least %d characters", MinUsernameLength)
	}
	if len(username) > MaxUsernameLength {
		return nil, fmt.Errorf("username must not exceed %d characters", MaxUsernameLength)
	}
	if !usernameRegex.MatchString(username) {
		return nil, errors.New("username may only contain letters, numbers, and underscores")
	}

	if len(passwordHash) == 0 {
		return nil, errors.New("password hash cannot be empty")
	}
	if len(passwordSalt) == 0 {
		return nil, errors.New("password salt cannot be empty")
	}

	trimmedEmail := strings.TrimSpace(email)
	if len(trimmedEmail) == 0 {
		return nil, errors.New("email cannot be empty")
	}
	normalizedEmail := strings.ToLower(trimmedEmail)
	if len(normalizedEmail) > MaxEmailLength {
		return nil, fmt.Errorf("email must not exceed %d characters", MaxEmailLength)
	}
	if !ValidEmail(normalizedEmail) {
		return nil, errors.New("email format is invalid")
	}

	return &User{
		Username:     username,
		Email:        normalizedEmail,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
	}, nil
}
