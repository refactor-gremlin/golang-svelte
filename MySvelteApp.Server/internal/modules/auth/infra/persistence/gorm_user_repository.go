package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"

	authapp "mysvelteapp/server_new/internal/modules/auth/app"
	authdomain "mysvelteapp/server_new/internal/modules/auth/domain"
)

var _ authapp.UserRepository = (*GormUserRepository)(nil)

// GormUserRepository persists users using GORM.
type GormUserRepository struct {
	db *gorm.DB
}

// NewGormUserRepository constructs a repository backed by GORM.
func NewGormUserRepository(db *gorm.DB) *GormUserRepository {
	return &GormUserRepository{db: db}
}

// Add inserts the provided user into the database.
func (r *GormUserRepository) Add(ctx context.Context, user *authdomain.User) error {
	if user == nil {
		return fmt.Errorf("user cannot be nil")
	}
	return r.db.WithContext(ctx).Create(user).Error
}

// GetByUsername fetches a user by username; returns nil when not found.
func (r *GormUserRepository) GetByUsername(ctx context.Context, username string) (*authdomain.User, error) {
	trimmed := strings.TrimSpace(username)
	if trimmed == "" {
		return nil, fmt.Errorf("username cannot be blank")
	}

	var user authdomain.User
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
func (r *GormUserRepository) UsernameExists(ctx context.Context, username string) (bool, error) {
	trimmed := strings.TrimSpace(username)
	if trimmed == "" {
		return false, fmt.Errorf("username cannot be blank")
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&authdomain.User{}).
		Where("username = ?", trimmed).
		Count(&count).
		Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

// EmailExists checks whether an email address is already stored.
func (r *GormUserRepository) EmailExists(ctx context.Context, email string) (bool, error) {
	trimmed := strings.TrimSpace(email)
	if trimmed == "" {
		return false, fmt.Errorf("email cannot be blank")
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&authdomain.User{}).
		Where("email = ?", trimmed).
		Count(&count).
		Error; err != nil {
		return false, err
	}

	return count > 0, nil
}
