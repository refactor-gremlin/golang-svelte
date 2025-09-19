package auth

// RegisterRequest represents the payload needed to create a new user.
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

// AuthErrorType classifies why an authentication operation failed.
type AuthErrorType string

const (
	AuthErrorTypeValidation   AuthErrorType = "validation"
	AuthErrorTypeUnauthorized AuthErrorType = "unauthorized"
	AuthErrorTypeConflict     AuthErrorType = "conflict"
	AuthErrorTypeUnknown      AuthErrorType = "unknown"
)

// AuthResult mirrors the previous .NET contract so presentation logic can stay similar.
type AuthResult struct {
	Success      bool
	Token        string
	UserID       uint
	Username     string
	ErrorMessage string
	ErrorType    AuthErrorType
}
