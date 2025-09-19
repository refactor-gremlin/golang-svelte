package pokemon

// RandomPokemonDto represents the response model for a random Pokemon.
// This mirrors the C# RandomPokemonDto exactly.
type RandomPokemonDto struct {
	Name  *string `json:"name,omitempty"`
	Type  *string `json:"type,omitempty"`
	Image *string `json:"image,omitempty"`
}
