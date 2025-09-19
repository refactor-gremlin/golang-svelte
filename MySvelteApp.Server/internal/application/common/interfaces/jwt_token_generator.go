package interfaces

import (
	"mysvelteapp/server/internal/domain/entities"
)

// JwtTokenGenerator abstracts the token creation so the application layer stays persistence-agnostic.
type JwtTokenGenerator interface {
	GenerateToken(user *entities.User) (string, error)
}
