package pokeapi

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"

	pokemonapp "mysvelteapp/server_new/internal/modules/pokemon/app"
	pokemondomain "mysvelteapp/server_new/internal/modules/pokemon/domain"
)

const (
	pokemonAPIBaseURL = "https://pokeapi.co/api/v2/pokemon/"
	pokemonCountURL   = "https://pokeapi.co/api/v2/pokemon-species/?limit=0"
)

var _ pokemonapp.RandomPokemonPort = (*Adapter)(nil)

// Adapter integrates with the external PokeAPI.
type Adapter struct {
	httpClient *http.Client
}

// NewAdapter creates a new Adapter instance.
func NewAdapter(httpClient *http.Client) *Adapter {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Adapter{httpClient: httpClient}
}

// GetRandomPokemon retrieves a random Pokemon from the PokeAPI.
func (a *Adapter) GetRandomPokemon(ctx context.Context) (*pokemondomain.RandomPokemon, error) {
	count, err := a.getPokemonCount(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Pokemon count: %w", err)
	}

	randomPokemon := rand.Intn(count) + 1
	pokemonURL := fmt.Sprintf("%s%d", pokemonAPIBaseURL, randomPokemon)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pokemonURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := a.httpClient.Do(req)
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

	var apiResp pokeAPIResponse
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to deserialize Pokemon data: %w", err)
	}

	var types []string
	for _, t := range apiResp.Types {
		types = append(types, t.Type.Name)
	}
	typeStr := strings.Join(types, ", ")

	return &pokemondomain.RandomPokemon{
		Name:  &apiResp.Name,
		Type:  &typeStr,
		Image: apiResp.Sprites.FrontDefault,
	}, nil
}

func (a *Adapter) getPokemonCount(ctx context.Context) (int, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, pokemonCountURL, nil)
	if err != nil {
		return 0, fmt.Errorf("failed to create count request: %w", err)
	}

	resp, err := a.httpClient.Do(req)
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

	var countResp struct {
		Count int `json:"count"`
	}
	if err := json.Unmarshal(body, &countResp); err != nil {
		return 0, fmt.Errorf("failed to deserialize count data: %w", err)
	}

	return countResp.Count, nil
}

type pokeAPIResponse struct {
	Name    string         `json:"name"`
	Types   []pokeAPIType  `json:"types"`
	Sprites pokeAPISprites `json:"sprites"`
}

type pokeAPIType struct {
	Type typeInfo `json:"type"`
}

type typeInfo struct {
	Name string `json:"name"`
}

type pokeAPISprites struct {
	FrontDefault *string `json:"front_default"`
}
