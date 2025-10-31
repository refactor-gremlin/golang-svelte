package app

// Package app provides the simplified service layer for the Pokemon module.
// This module demonstrates a streamlined approach where business logic is
// contained within a single service layer without unnecessary abstractions.
//
// Simplified Data Flow:
// HTTP Handler → Service → External API (PokeAPI)
// HTTP Response ← Service ← External API Response
//
// Architecture Decision:
// Unlike the auth module, the pokemon module has simple data retrieval
// requirements without complex business rules, so we've eliminated:
// - Repository pattern (no database persistence)
// - Port/Adapter pattern (direct HTTP client usage)
// - Domain layer (simple data transfer objects)
//
// This results in a cleaner, more maintainable codebase for simple use cases.

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	pokemonAPIBaseURL    = "https://pokeapi.co/api/v2/pokemon/"
	pokemonCountURL      = "https://pokeapi.co/api/v2/pokemon-species/?limit=0"
	pokemonCountCacheTTL = 10 * time.Minute
)

// ExternalAPIError indicates the upstream Pokemon API could not satisfy the request.
type ExternalAPIError struct {
	Err error
}

func (e ExternalAPIError) Error() string {
	return "external pokemon service unavailable"
}

// Unwrap returns the underlying error for logging/inspection.
func (e ExternalAPIError) Unwrap() error {
	return e.Err
}

type countCache struct {
	value     int
	expiresAt time.Time
}

// RandomPokemonResponse represents the HTTP response model for a random Pokemon.
// This model is designed specifically for JSON serialization and follows the
// frontend's expected contract with camelCase fields and proper null handling.
//
// Response Design:
// - Uses pointers to properly handle null values from external API
// - omitempty tags ensure clean JSON output
// - Focuses on essential fields required by the frontend
type RandomPokemonResponse struct {
	Name  *string `json:"name,omitempty"`  // Pokemon name (nullable)
	Type  *string `json:"type,omitempty"`  // Pokemon types as string (nullable)
	Image *string `json:"image,omitempty"` // Pokemon sprite URL (nullable)
}

// Service handles Pokemon operations by directly integrating with the external PokeAPI.
// This service demonstrates a simplified architecture where external integration
// is handled directly without additional abstraction layers.
//
// Design Philosophy:
// - Direct HTTP client usage for simple external API integration
// - No repository pattern since we're not persisting data
// - No domain layer since we're just data transformation
// - Focused on the specific use case (random Pokemon retrieval)
type Service struct {
	httpClient *http.Client // HTTP client for external API calls
	random     *rand.Rand   // pseudo-random generator, seeded once

	pokemonBaseURL string
	countURL       string

	randMu sync.Mutex
	mu     sync.RWMutex
	cache  countCache
	now    func() time.Time
}

// NewService creates a new Pokemon service with configurable HTTP client.
// If no HTTP client is provided, a default one with 30-second timeout is created.
// This design allows for easy testing with mock HTTP clients.
func NewService(httpClient *http.Client) *Service {
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 30 * time.Second}
	}
	return &Service{
		httpClient:     httpClient,
		random:         rand.New(rand.NewSource(time.Now().UnixNano())),
		pokemonBaseURL: pokemonAPIBaseURL,
		countURL:       pokemonCountURL,
		now:            time.Now,
	}
}

// GetRandomPokemon retrieves a random Pokemon from the external PokeAPI.
//
// Data Flow Process:
// 1. Fetch total Pokemon count from PokeAPI species endpoint
// 2. Generate random ID within the valid range
// 3. Retrieve specific Pokemon data using the random ID
// 4. Transform the complex API response into our simplified RandomPokemon DTO
// 5. Return the simplified data structure
//
// Error Handling:
// - Context cancellation is respected throughout the process
// - Network errors are wrapped with descriptive messages
// - API response validation ensures data integrity
// - Timeout handling prevents hanging requests
//
// External Dependencies:
// - PokeAPI (https://pokeapi.co/api/v2/)
// - Requires internet connectivity
// - Rate limiting considerations (PokeAPI is generous but has limits)
func (s *Service) GetRandomPokemon(ctx context.Context) (*RandomPokemonResponse, error) {
	count, err := s.getPokemonCount(ctx)
	if err != nil {
		return nil, err
	}

	if count <= 0 {
		return nil, ExternalAPIError{Err: fmt.Errorf("invalid pokemon count: %d", count)}
	}

	s.randMu.Lock()
	randomPokemon := s.random.Intn(count) + 1
	s.randMu.Unlock()
	pokemonURL := fmt.Sprintf("%s%d", s.pokemonBaseURL, randomPokemon)

	apiResp, err := s.fetchPokemon(ctx, pokemonURL)
	if err != nil {
		return nil, err
	}

	var types []string
	for _, t := range apiResp.Types {
		types = append(types, t.Type.Name)
	}
	typeStr := strings.Join(types, ", ")

	return &RandomPokemonResponse{
		Name:  &apiResp.Name,
		Type:  &typeStr,
		Image: apiResp.Sprites.FrontDefault,
	}, nil
}

func (s *Service) getPokemonCount(ctx context.Context) (int, error) {
	if count, ok := s.cachedPokemonCount(); ok {
		return count, nil
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, s.countURL, nil)
	if err != nil {
		return 0, ExternalAPIError{Err: fmt.Errorf("create count request: %w", err)}
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return 0, ctxErr
		}
		return 0, ExternalAPIError{Err: fmt.Errorf("retrieve count: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, ExternalAPIError{Err: fmt.Errorf("count request returned status %d", resp.StatusCode)}
	}

	var countResp struct {
		Count int `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&countResp); err != nil {
		return 0, ExternalAPIError{Err: fmt.Errorf("decode count response: %w", err)}
	}

	if countResp.Count <= 0 {
		return 0, ExternalAPIError{Err: fmt.Errorf("received non-positive count %d", countResp.Count)}
	}

	s.storePokemonCount(countResp.Count)

	return countResp.Count, nil
}

func (s *Service) cachedPokemonCount() (int, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.cache.value > 0 && s.cache.expiresAt.After(s.now()) {
		return s.cache.value, true
	}
	return 0, false
}

func (s *Service) storePokemonCount(count int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.cache = countCache{
		value:     count,
		expiresAt: s.now().Add(pokemonCountCacheTTL),
	}
}

func (s *Service) fetchPokemon(ctx context.Context, url string) (*pokemonAPIResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, ExternalAPIError{Err: fmt.Errorf("create pokemon request: %w", err)}
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return nil, ctxErr
		}
		return nil, ExternalAPIError{Err: fmt.Errorf("retrieve pokemon: %w", err)}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, ExternalAPIError{Err: fmt.Errorf("pokemon request returned status %d", resp.StatusCode)}
	}

	var apiResp pokemonAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&apiResp); err != nil {
		return nil, ExternalAPIError{Err: fmt.Errorf("decode pokemon response: %w", err)}
	}

	return &apiResp, nil
}

type pokemonAPIResponse struct {
	Name  string `json:"name"`
	Types []struct {
		Type struct {
			Name string `json:"name"`
		} `json:"type"`
	} `json:"types"`
	Sprites struct {
		FrontDefault *string `json:"front_default"`
	} `json:"sprites"`
}
