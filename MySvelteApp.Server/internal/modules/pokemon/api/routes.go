package api

import (
	"github.com/gin-gonic/gin"
	pokemonapp "mysvelteapp/server_new/internal/modules/pokemon/app"
)

// RegisterRoutes mounts the pokemon routes beneath the provided router group.
// This function sets up the REST API endpoints for Pokemon-related operations.
//
// Routes Defined:
// GET /RandomPokemon - Retrieves a random Pokemon from external API
//
// Design Notes:
// - No route grouping needed as there's only one endpoint
// - Endpoint name follows camelCase convention for API consistency
// - Direct route mapping to handler for simplicity
//
// Data Flow:
// HTTP Request → Gin Router → Handler → Service → External API (PokeAPI)
// HTTP Response ← Handler ← Service ← External API Response ← JSON Transformation
func RegisterRoutes(router gin.IRouter, service *pokemonapp.Service) {
	router.GET("/RandomPokemon", GetRandomPokemon(service))
}
