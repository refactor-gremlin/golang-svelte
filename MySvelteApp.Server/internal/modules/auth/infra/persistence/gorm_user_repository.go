package persistence

import (
	"context"
	"errors"
	"fmt"
	"strings"

	sqlite "github.com/mattn/go-sqlite3"
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

	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		if isUniqueConstraintError(err) {
			if conflict := r.detectConflict(ctx, user); conflict != nil {
				return conflict
			}
			return authapp.ConflictError{Message: "A user with the provided credentials already exists."}
		}
		return err
	}

	return nil
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
	normalized := strings.ToLower(strings.TrimSpace(email))
	if normalized == "" {
		return false, fmt.Errorf("email cannot be blank")
	}

	var count int64
	if err := r.db.WithContext(ctx).
		Model(&authdomain.User{}).
		Where("email = ?", normalized).
		Count(&count).
		Error; err != nil {
		return false, err
	}

	return count > 0, nil
}

func (r *GormUserRepository) detectConflict(ctx context.Context, user *authdomain.User) error {
	if exists, err := r.UsernameExists(ctx, user.Username); err == nil && exists {
		return authapp.ConflictError{Message: "This username is already taken. Please choose a different one."}
	}
	if exists, err := r.EmailExists(ctx, user.Email); err == nil && exists {
		return authapp.ConflictError{Message: "This email is already registered. Please use a different email address."}
	}
	return nil
}

func isUniqueConstraintError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	var sqliteErr sqlite.Error
	if errors.As(err, &sqliteErr) {
		if sqliteErr.Code == sqlite.ErrConstraint {
			switch sqliteErr.ExtendedCode {
			case sqlite.ErrConstraintUnique, sqlite.ErrConstraintPrimaryKey:
				return true
			}
		}
	}
	return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
