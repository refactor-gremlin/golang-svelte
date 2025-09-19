package app

import (
	"context"

	pokemondomain "mysvelteapp/server_new/internal/modules/pokemon/domain"
)

// Service orchestrates Pokemon use-cases.
type Service struct {
	port RandomPokemonPort
}

// NewService wires the port into the service.
func NewService(port RandomPokemonPort) *Service {
	return &Service{port: port}
}

// GetRandomPokemon fetches a random Pokemon using the configured port.
func (s *Service) GetRandomPokemon(ctx context.Context) (*pokemondomain.RandomPokemon, error) {
	return s.port.GetRandomPokemon(ctx)
}
