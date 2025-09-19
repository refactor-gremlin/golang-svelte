package entities

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	// MaxUsernameLength mirrors the original .NET constraints
	MaxUsernameLength = 64
	// MaxEmailLength mirrors the original .NET constraints
	MaxEmailLength = 320
)

// User represents an authenticated user in the system.
// The struct tags map the entity to GORM so we can persist it with the same constraints.
type User struct {
	ID           uint      `gorm:"primaryKey"`
	Username     string    `gorm:"size:64;uniqueIndex;not null"`
	Email        string    `gorm:"size:320;uniqueIndex;not null"`
	PasswordHash string    `gorm:"size:512;not null"`
	PasswordSalt string    `gorm:"size:256;not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`
}

// NewUser applies the invariants that previously lived in the .NET domain model.
func NewUser(username, email, passwordHash, passwordSalt string) (*User, error) {
	username = strings.TrimSpace(username)
	if len(username) == 0 {
		return nil, errors.New("username cannot be empty")
	}
	if len(username) > MaxUsernameLength {
		return nil, fmt.Errorf("username must not exceed %d characters", MaxUsernameLength)
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

	return &User{
		Username:     username,
		Email:        normalizedEmail,
		PasswordHash: passwordHash,
		PasswordSalt: passwordSalt,
	}, nil
}
