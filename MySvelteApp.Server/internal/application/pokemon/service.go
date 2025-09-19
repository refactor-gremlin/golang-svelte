package pokemon

import "context"

// IRandomPokemonService defines the contract for retrieving random Pokemon data.
// This mirrors the C# IRandomPokemonService interface exactly.
type IRandomPokemonService interface {
	GetRandomPokemonAsync(ctx context.Context) (*RandomPokemonDto, error)
}
