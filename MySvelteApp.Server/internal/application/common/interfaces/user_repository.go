package interfaces

import (
	"context"

	"mysvelteapp/server/internal/domain/entities"
)

// UserRepository exposes the persistence operations required by the auth use-cases.
type UserRepository interface {
	Add(ctx context.Context, user *entities.User) error
	GetByUsername(ctx context.Context, username string) (*entities.User, error)
	UsernameExists(ctx context.Context, username string) (bool, error)
	EmailExists(ctx context.Context, email string) (bool, error)
}
