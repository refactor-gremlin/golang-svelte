package api

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	pokemonapp "mysvelteapp/server_new/internal/modules/pokemon/app"
)

// GetRandomPokemon handles requests for a random Pokemon.
//
// Simplified Data Flow:
// HTTP GET /RandomPokemon
// → Context propagation for cancellation support
// → External API integration (PokeAPI)
// → RandomPokemonResponse JSON serialization
//
// Error Handling:
// - 503: Upstream Pokemon API outage or failure
// - 500: Unexpected internal errors
// - Context cancellation is properly respected
//
// External Dependencies:
// This endpoint relies on the external PokeAPI service being available.
// Network failures or API outages will result in 503 status codes.
func GetRandomPokemon(service *pokemonapp.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		pokemon, err := service.GetRandomPokemon(c.Request.Context())
		if err != nil {
			c.Error(err) // attach to gin context for observability

			var upstreamErr pokemonapp.ExternalAPIError
			switch {
			case errors.As(err, &upstreamErr):
				c.JSON(http.StatusServiceUnavailable, gin.H{"error": "Pokemon service temporarily unavailable"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get random Pokemon"})
			}
			return
		}

		c.JSON(http.StatusOK, pokemon)
	}
}
