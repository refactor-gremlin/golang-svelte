package api

// Package api provides HTTP handlers and models for the authentication module.
// This layer handles HTTP request/response transformation and exposes the authentication endpoints.
//
// Data Flow:
// HTTP Request → Handler → Service → Repository → Database
// HTTP Response ← Handler ← Service ← Repository ← Database
//
// The API layer is responsible for:
// - HTTP request validation and binding
// - Response formatting and status codes
// - Error handling appropriate for HTTP clients
// - Swagger documentation generation

// AuthErrorResponse provides a standardized error response structure.
// All API errors are returned in this format to ensure consistent
// error handling on the client side.
// @name AuthErrorResponse
type AuthErrorResponse struct {
	Message string `json:"message"` // Human-readable error message
}

// AuthSuccessResponse provides the standard success payload for authentication endpoints.
// @name AuthSuccessResponse
type AuthSuccessResponse struct {
	Token    string `json:"token"`
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
}
