package api

import (
	"net/http"

	"github.com/gin-gonic/gin"

	pokemonapp "mysvelteapp/server_new/internal/modules/pokemon/app"
)

// Handlers exposes HTTP endpoints for the pokemon module.
type Handlers struct {
	service *pokemonapp.Service
}

// NewHandlers wires the pokemon service into HTTP handlers.
func NewHandlers(service *pokemonapp.Service) *Handlers {
	return &Handlers{service: service}
}

// GetRandomPokemon godoc
// @Summary Get a random Pokemon
// @Description Retrieves a random Pokemon from the PokeAPI
// @Tags pokemon
// @Accept json
// @Produce json
// @Success 200 {object} RandomPokemonResponse
// @Failure 500 {object} map[string]string
// @Router /RandomPokemon [get]
func (h *Handlers) GetRandomPokemon(c *gin.Context) {
	pokemon, err := h.service.GetRandomPokemon(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get random Pokemon"})
		return
	}

	c.JSON(http.StatusOK, RandomPokemonResponse{
		Name:  pokemon.Name,
		Type:  pokemon.Type,
		Image: pokemon.Image,
	})
}
