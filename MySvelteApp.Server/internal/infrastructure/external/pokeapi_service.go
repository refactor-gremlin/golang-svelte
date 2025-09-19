package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"mysvelteapp/server/internal/application/pokemon"
)

const (
	pokemonAPIBaseURL = "https://pokeapi.co/api/v2/pokemon/"
	pokemonCountURL   = "https://pokeapi.co/api/v2/pokemon-species/?limit=0"
)

// PokeApiRandomPokemonService implements the IRandomPokemonService interface.
// This mirrors the C# PokeApiRandomPokemonService exactly.
type PokeApiRandomPokemonService struct {
	httpClient *http.Client
}

// NewPokeApiRandomPokemonService creates a new instance of PokeApiRandomPokemonService.
func NewPokeApiRandomPokemonService(httpClient *http.Client) *PokeApiRandomPokemonService {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &PokeApiRandomPokemonService{
		httpClient: httpClient,
	}
}

// GetRandomPokemonAsync retrieves a random Pokemon from the PokeAPI.
// This mirrors the C# implementation exactly.
func (p *PokeApiRandomPokemonService) GetRandomPokemonAsync(ctx context.Context) (*pokemon.RandomPokemonDto, error) {
	count, err := p.getPokemonCountAsync(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Pokemon count: %w", err)
	}

	randomPokemon := rand.Intn(count) + 1
	pokemonURL := fmt.Sprintf("%s%d", pokemonAPIBaseURL, randomPokemon)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pokemonURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get Pokemon data: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Pokemon API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var pokeAPI pokeAPIResponse
	if err := json.Unmarshal(body, &pokeAPI); err != nil {
		return nil, fmt.Errorf("failed to deserialize Pokemon data: %w", err)
	}

	// Join types like in C# implementation
	var types []string
	for _, t := range pokeAPI.Types {
		types = append(types, t.Type.Name)
	}
	typeStr := strings.Join(types, ", ")

	return &pokemon.RandomPokemonDto{
		Name:  &pokeAPI.Name,
		Type:  &typeStr,
		Image: pokeAPI.Sprites.FrontDefault,
	}, nil
}

// getPokemonCountAsync gets the total count of Pokemon from the API.
// This mirrors the C# implementation exactly.
func (p *PokeApiRandomPokemonService) getPokemonCountAsync(ctx context.Context) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pokemonCountURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create count request: %w", err)
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed to get Pokemon count: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("Pokemon count API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed to read count response body: %w", err)
	}

	var countResponse struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal(body, &countResponse); err != nil {
		return 0, fmt.Errorf("failed to deserialize count data: %w", err)
	}

	return countResponse.Count, nil
}

// pokeAPIResponse represents the Pokemon API response structure.
// This mirrors the C# PokeApiResponse class exactly.
type pokeAPIResponse struct {
	Name    string         `json:"name"`
	Types   []pokeAPIType  `json:"types"`
	Sprites pokeAPISprites `json:"sprites"`
}

// pokeAPIType represents a Pokemon type in the API response.
// This mirrors the C# PokeApiType class exactly.
type pokeAPIType struct {
	Type typeInfo `json:"type"`
}

// typeInfo represents type information in the API response.
// This mirrors the C# TypeInfo class exactly.
type typeInfo struct {
	Name string `json:"name"`
}

// pokeAPISprites represents the sprites in the Pokemon API response.
// This mirrors the C# PokeApiSprites class exactly.
type pokeAPISprites struct {
	FrontDefault *string `json:"front_default"`
}
