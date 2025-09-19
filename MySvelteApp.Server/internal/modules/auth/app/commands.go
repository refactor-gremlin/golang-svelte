package app

// RegisterRequest represents the payload required to create a new user account.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// LoginRequest represents the credentials submitted by an existing user.
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthSuccess encapsulates the data returned on successful authentication.
type AuthSuccess struct {
	Token    string
	UserID   uint
	Username string
}
