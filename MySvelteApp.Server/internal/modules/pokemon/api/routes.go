package api

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the pokemon routes beneath the provided router group.
func RegisterRoutes(router gin.IRouter, handlers *Handlers) {
	router.GET("/RandomPokemon", handlers.GetRandomPokemon)
}
