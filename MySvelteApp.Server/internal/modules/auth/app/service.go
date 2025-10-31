package app

// Package app contains the application service layer for the authentication module.
// This layer orchestrates business logic and coordinates between different infrastructure components.
//
// Data Flow Architecture:
// API Layer (HTTP) → Service Layer (Business Logic) → Infrastructure Layer (Data/External)
//
// The service layer is responsible for:
// - Business rule enforcement and validation
// - Coordinating between repositories and external services
// - Transaction management and data consistency
// - Domain events and use case orchestration
//
// Dependencies are injected via interfaces to maintain loose coupling and testability.

import (
	"context"
	"fmt"
	"strings"
	"unicode"

	authdomain "mysvelteapp/server_new/internal/modules/auth/domain"
)

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

const (
	minPasswordLength = 8
	maxPasswordLength = 512
)

// Service exposes the authentication use-cases and business logic.
// This service orchestrates the complete authentication flow by coordinating
// between user persistence, password security, and token generation.
//
// Data Flow Dependencies:
// - UserRepository: Handles user data persistence and retrieval
// - PasswordHasher: Manages secure password hashing and verification
// - TokenGenerator: Creates JWT tokens for authenticated sessions
//
// All dependencies are injected as interfaces to enable testing and
// maintain flexibility in implementation choices.
type Service struct {
	users  UserRepository // User persistence operations
	hasher PasswordHasher // Password security operations
	tokens TokenGenerator // JWT token generation
}

// NewService creates a new Service instance with injected dependencies.
// This follows the dependency injection pattern where all external
// dependencies are provided from the outside, making the service
// easily testable with mock implementations.
func NewService(users UserRepository, hasher PasswordHasher, tokens TokenGenerator) *Service {
	return &Service{
		users:  users,
		hasher: hasher,
		tokens: tokens,
	}
}

// Register creates a new user account after comprehensive validation.
//
// Data Flow:
// 1. Input validation (username, email, password format)
// 2. Business rule validation (username/email uniqueness)
// 3. Password hashing with salt generation
// 4. Domain entity creation with invariants
// 5. Database persistence
// 6. JWT token generation for immediate login
//
// Error Conditions:
// - ValidationError: Invalid input format or business rules
// - ConflictError: Username or email already exists
// - Repository errors: Database connectivity/operation failures
// - Token generation errors: JWT service failures
func (s *Service) Register(ctx context.Context, cmd RegisterRequest) (*AuthSuccess, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := validateRegister(cmd); err != nil {
		return nil, err
	}

	trimmedUsername := strings.TrimSpace(cmd.Username)
	normalizedEmail := strings.ToLower(strings.TrimSpace(cmd.Email))

	exists, err := s.users.UsernameExists(ctx, trimmedUsername)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ConflictError{Message: "This username is already taken. Please choose a different one."}
	}

	emailExists, err := s.users.EmailExists(ctx, normalizedEmail)
	if err != nil {
		return nil, err
	}
	if emailExists {
		return nil, ConflictError{Message: "This email is already registered. Please use a different email address."}
	}

	hash, salt, err := s.hasher.HashPassword(cmd.Password)
	if err != nil {
		return nil, err
	}

	user, err := authdomain.NewUser(trimmedUsername, normalizedEmail, hash, salt)
	if err != nil {
		return nil, err
	}

	if err := s.users.Add(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.tokens.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthSuccess{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}

// Login authenticates a user with username and password credentials.
//
// Data Flow:
// 1. Input validation (non-empty username/password)
// 2. Repository lookup by username
// 3. Password hash verification against stored hash and salt
// 4. JWT token generation for authenticated session
//
// Security Considerations:
// - Uses constant-time comparison for password verification
// - Returns generic error messages to prevent username enumeration
// - Generates fresh JWT token for each successful login
//
// Error Conditions:
// - ValidationError: Missing username or password
// - UnauthorizedError: Invalid credentials (generic message)
// - Repository errors: Database connectivity/operation failures
// - Token generation errors: JWT service failures
func (s *Service) Login(ctx context.Context, cmd LoginRequest) (*AuthSuccess, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}

	if err := validateLogin(cmd); err != nil {
		return nil, err
	}

	trimmedUsername := strings.TrimSpace(cmd.Username)

	user, err := s.users.GetByUsername(ctx, trimmedUsername)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, unauthorizedError()
	}

	valid, err := s.hasher.VerifyPassword(cmd.Password, user.PasswordHash, user.PasswordSalt)
	if err != nil {
		return nil, err
	}
	if !valid {
		return nil, unauthorizedError()
	}

	token, err := s.tokens.GenerateToken(user)
	if err != nil {
		return nil, err
	}

	return &AuthSuccess{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}

func validateRegister(cmd RegisterRequest) error {
	username := strings.TrimSpace(cmd.Username)
	switch {
	case username == "":
		return ValidationError{Message: "Username is required."}
	case len(username) < authdomain.MinUsernameLength:
		return ValidationError{Message: fmt.Sprintf("Username must be at least %d characters long.", authdomain.MinUsernameLength)}
	case len(username) > authdomain.MaxUsernameLength:
		return ValidationError{Message: fmt.Sprintf("Username must not exceed %d characters.", authdomain.MaxUsernameLength)}
	case !authdomain.ValidUsername(username):
		return ValidationError{Message: "Username can only contain letters, numbers, and underscores."}
	}

	email := strings.TrimSpace(cmd.Email)
	switch {
	case email == "":
		return ValidationError{Message: "Email is required."}
	case len(email) > authdomain.MaxEmailLength:
		return ValidationError{Message: fmt.Sprintf("Email must not exceed %d characters.", authdomain.MaxEmailLength)}
	case !authdomain.ValidEmail(email):
		return ValidationError{Message: "Please enter a valid email address."}
	}

	switch {
	case strings.TrimSpace(cmd.Password) == "":
		return ValidationError{Message: "Password is required."}
	case len(cmd.Password) < minPasswordLength:
		return ValidationError{Message: "Password must be at least 8 characters long."}
	case len(cmd.Password) > maxPasswordLength:
		return ValidationError{Message: "Password must not exceed 512 characters."}
	case !passwordMeetsRequirements(cmd.Password):
		return ValidationError{Message: "Password must contain at least one uppercase letter, one lowercase letter, and one number."}
	}

	return nil
}

func validateLogin(cmd LoginRequest) error {
	if strings.TrimSpace(cmd.Username) == "" {
		return ValidationError{Message: "Username is required."}
	}
	if strings.TrimSpace(cmd.Password) == "" {
		return ValidationError{Message: "Password is required."}
	}
	return nil
}

func passwordMeetsRequirements(password string) bool {
	var hasUpper, hasLower, hasDigit bool
	for _, r := range password {
		switch {
		case unicode.IsUpper(r):
			hasUpper = true
		case unicode.IsLower(r):
			hasLower = true
		case unicode.IsDigit(r):
			hasDigit = true
		}
	}
	return hasUpper && hasLower && hasDigit
}

func unauthorizedError() error {
	return UnauthorizedError{Message: "Invalid username or password. Please check your credentials and try again."}
}
