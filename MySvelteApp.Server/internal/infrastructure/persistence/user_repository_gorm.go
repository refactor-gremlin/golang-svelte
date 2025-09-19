package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"mysvelteapp/server/internal/application/common/interfaces"
	"mysvelteapp/server/internal/domain/entities"
)

var _ interfaces.UserRepository = (*UserRepository)(nil)

// UserRepository persists users using GORM.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository wires a GORM-backed implementation of the auth repository interface.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Add inserts the provided user into the database.
func (r *UserRepository) Add(ctx context.Context, user *entities.User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByUsername fetches a user by username; returns nil when not found.
func (r *UserRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	trimmed := strings.TrimSpace(username)
	if trimmed == "" {
		return nil, fmt.Errorf("username cannot be blank")
	}

	var user entities.User
	err := r.db.WithContext(ctx).
		Where("username = ?", trimmed).
		Take(&user).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// UsernameExists checks whether a username is already stored.
func (r *UserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	trimmed := strings.TrimSpace(username)
	if trimmed == "" {
		return false, fmt.Errorf("username cannot be blank")
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("username = ?", trimmed).
		Count(&count).
		Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// EmailExists checks whether an email address is already stored.
func (r *UserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	trimmed := strings.TrimSpace(email)
	if trimmed == "" {
		return false, fmt.Errorf("email cannot be blank")
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&entities.User{}).
		Where("email = ?", trimmed).
		Count(&count).
		Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
