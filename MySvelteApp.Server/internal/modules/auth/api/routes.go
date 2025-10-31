package api

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the authentication routes beneath the provided router group.
// This function sets up the REST API endpoints for user authentication operations.
//
// Routes Defined:
// POST /auth/register - Creates a new user account
// POST /auth/login    - Authenticates existing users
//
// Route Grouping:
// All auth routes are grouped under "/auth" prefix for organization
// and to enable future auth-specific middleware (rate limiting, etc.)
//
// Data Flow:
// HTTP Request → Gin Router → Handler → Service → Repository/Infrastructure
// HTTP Response ← Handler ← Service ← Repository/Infrastructure ← Database/External API
func RegisterRoutes(router gin.IRouter, handlers *Handlers) {
	auth := router.Group("/auth")
	auth.POST("/register", handlers.Register)
	auth.POST("/login", handlers.Login)
}
