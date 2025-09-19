package app

import (
	"context"

	pokemondomain "mysvelteapp/server_new/internal/modules/pokemon/domain"
)

// RandomPokemonPort defines the contract for retrieving random Pokemon data.
type RandomPokemonPort interface {
	GetRandomPokemon(ctx context.Context) (*pokemondomain.RandomPokemon, error)
}
