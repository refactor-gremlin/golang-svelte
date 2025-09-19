package api

// AuthSuccessResponse matches the JSON contract expected by the frontend generator.
// @name AuthSuccessResponse
type AuthSuccessResponse struct {
	Token    string `json:"token"`
	UserID   uint   `json:"userId"`
	Username string `json:"username"`
}

// AuthErrorResponse wraps error messages in a serialisable structure.
// @name AuthErrorResponse
type AuthErrorResponse struct {
	Message string `json:"message"`
}

// RegisterRequest represents the registration payload.
// @name RegisterRequest
type RegisterRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Username string `json:"username"`
}

// LoginRequest represents the login payload.
// @name LoginRequest
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
