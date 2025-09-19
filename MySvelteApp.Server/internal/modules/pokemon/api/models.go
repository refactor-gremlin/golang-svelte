package api

// RandomPokemonResponse represents the response model for a random Pokemon.
// @name RandomPokemonResponse
type RandomPokemonResponse struct {
	Name  *string `json:"name,omitempty"`
	Type  *string `json:"type,omitempty"`
	Image *string `json:"image,omitempty"`
}
