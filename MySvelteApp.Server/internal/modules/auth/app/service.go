package app

import (
	"context"
	"regexp"
	"strings"
	"unicode"

	authdomain "mysvelteapp/server_new/internal/modules/auth/domain"
)

const (
	minUsernameLength = 3
	minPasswordLength = 8
	maxPasswordLength = 512
)

var (
	usernameRegex = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	emailRegex    = regexp.MustCompile(`^[^\s@]+@[^\s@.]+\.[^\s@.]+$`)
)

// Service exposes the authentication use-cases.
type Service struct {
	users  UserRepository
	hasher PasswordHasher
	tokens TokenGenerator
}

// NewService wires the service dependencies.
func NewService(users UserRepository, hasher PasswordHasher, tokens TokenGenerator) *Service {
	return &Service{
		users:  users,
		hasher: hasher,
		tokens: tokens,
	}
}

// Register creates a new user account when the command is valid.
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

// Login authenticates an existing user with the provided credentials.
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
	case len(username) < minUsernameLength:
		return ValidationError{Message: "Username must be at least 3 characters long."}
	case len(username) > authdomain.MaxUsernameLength:
		return ValidationError{Message: "Username must not exceed 64 characters."}
	case !usernameRegex.MatchString(username):
		return ValidationError{Message: "Username can only contain letters, numbers, and underscores."}
	}

	email := strings.TrimSpace(cmd.Email)
	switch {
	case email == "":
		return ValidationError{Message: "Email is required."}
	case len(email) > authdomain.MaxEmailLength:
		return ValidationError{Message: "Email must not exceed 320 characters."}
	case strings.Contains(email, ".."):
		return ValidationError{Message: "Please enter a valid email address."}
	case !emailRegex.MatchString(email):
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
