package controllers

import (
	"encoding/json"
	"net/http"

	"mysvelteapp/server/internal/application/pokemon"
)

// RandomPokemonController exposes HTTP handlers for Pokemon routes.
// This mirrors the C# RandomPokemonController exactly.
type RandomPokemonController struct {
	pokemonService pokemon.IRandomPokemonService
}

// NewRandomPokemonController wires the controller to the Pokemon service.
func NewRandomPokemonController(pokemonService pokemon.IRandomPokemonService) *RandomPokemonController {
	return &RandomPokemonController{pokemonService: pokemonService}
}

// Get godoc
// @Summary Get a random Pokemon
// @Description Retrieves a random Pokemon from the PokeAPI
// @Tags pokemon
// @Accept json
// @Produce json
// @Success 200 {object} pokemon.RandomPokemonDto
// @Router /RandomPokemon [get]
func (c *RandomPokemonController) Get(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	pokemon, err := c.pokemonService.GetRandomPokemonAsync(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "Failed to get random Pokemon")
		return
	}

	writeJSONSuccess(w, http.StatusOK, pokemon)
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	response := map[string]string{"error": message}
	_ = json.NewEncoder(w).Encode(response)
}

func writeJSONSuccess(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
