package api

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the auth routes beneath the provided router group.
func RegisterRoutes(router gin.IRouter, handlers *Handlers) {
	auth := router.Group("/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
}
