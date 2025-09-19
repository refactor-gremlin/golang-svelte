package auth

import (
	"context"
	"regexp"
	"strings"
	"unicode"

	"mysvelteapp/server/internal/application/common/interfaces"
	"mysvelteapp/server/internal/domain/entities"
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

// Service contains the authentication use-cases.
type Service struct {
	userRepository interfaces.UserRepository
	passwordHasher interfaces.PasswordHasher
	tokenGenerator interfaces.JwtTokenGenerator
}

// NewService wires the dependencies of the authentication use-cases.
func NewService(
	userRepository interfaces.UserRepository,
	passwordHasher interfaces.PasswordHasher,
	tokenGenerator interfaces.JwtTokenGenerator,
) *Service {
	return &Service{
		userRepository: userRepository,
		passwordHasher: passwordHasher,
		tokenGenerator: tokenGenerator,
	}
}

// Register creates a new user account when the provided payload is valid.
func (s *Service) Register(ctx context.Context, request RegisterRequest) (AuthResult, error) {
	if err := ctx.Err(); err != nil {
		return AuthResult{}, err
	}

	if validationResult := validateRegisterRequest(request); validationResult != nil {
		return validationResult.toAuthResult(), nil
	}

	trimmedUsername := strings.TrimSpace(request.Username)
	normalizedEmail := strings.ToLower(strings.TrimSpace(request.Email))

	usernameExists, err := s.userRepository.UsernameExists(ctx, trimmedUsername)
	if err != nil {
		return AuthResult{}, err
	}
	if usernameExists {
		return AuthResult{
			Success:      false,
			ErrorMessage: "This username is already taken. Please choose a different one.",
			ErrorType:    AuthErrorTypeConflict,
		}, nil
	}

	emailExists, err := s.userRepository.EmailExists(ctx, normalizedEmail)
	if err != nil {
		return AuthResult{}, err
	}
	if emailExists {
		return AuthResult{
			Success:      false,
			ErrorMessage: "This email is already registered. Please use a different email address.",
			ErrorType:    AuthErrorTypeConflict,
		}, nil
	}

	hash, salt, err := s.passwordHasher.HashPassword(request.Password)
	if err != nil {
		return AuthResult{}, err
	}

	user, err := entities.NewUser(trimmedUsername, normalizedEmail, hash, salt)
	if err != nil {
		return AuthResult{}, err
	}

	if err := s.userRepository.Add(ctx, user); err != nil {
		return AuthResult{}, err
	}

	token, err := s.tokenGenerator.GenerateToken(user)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		Success:  true,
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}

// Login authenticates an existing user based on their credentials.
func (s *Service) Login(ctx context.Context, request LoginRequest) (AuthResult, error) {
	if err := ctx.Err(); err != nil {
		return AuthResult{}, err
	}

	if validationResult := validateLoginRequest(request); validationResult != nil {
		return validationResult.toAuthResult(), nil
	}

	trimmedUsername := strings.TrimSpace(request.Username)

	user, err := s.userRepository.GetByUsername(ctx, trimmedUsername)
	if err != nil {
		return AuthResult{}, err
	}
	if user == nil {
		return unauthorizedResult(), nil
	}

	valid, err := s.passwordHasher.VerifyPassword(request.Password, user.PasswordHash, user.PasswordSalt)
	if err != nil {
		return AuthResult{}, err
	}
	if !valid {
		return unauthorizedResult(), nil
	}

	token, err := s.tokenGenerator.GenerateToken(user)
	if err != nil {
		return AuthResult{}, err
	}

	return AuthResult{
		Success:  true,
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}

// validationError holds validation failure context before conversion to AuthResult.
type validationError struct {
	message string
	kind    AuthErrorType
}

func (v validationError) toAuthResult() AuthResult {
	return AuthResult{
		Success:      false,
		ErrorMessage: v.message,
		ErrorType:    v.kind,
	}
}

func validateRegisterRequest(request RegisterRequest) *validationError {
	username := strings.TrimSpace(request.Username)
	switch {
	case username == "":
		return &validationError{"Username is required.", AuthErrorTypeValidation}
	case len(username) < minUsernameLength:
		return &validationError{"Username must be at least 3 characters long.", AuthErrorTypeValidation}
	case len(username) > entities.MaxUsernameLength:
		return &validationError{
			message: "Username must not exceed 64 characters.",
			kind:    AuthErrorTypeValidation,
		}
	case !usernameRegex.MatchString(username):
		return &validationError{
			message: "Username can only contain letters, numbers, and underscores.",
			kind:    AuthErrorTypeValidation,
		}
	}

	email := strings.TrimSpace(request.Email)
	switch {
	case email == "":
		return &validationError{"Email is required.", AuthErrorTypeValidation}
	case len(email) > entities.MaxEmailLength:
		return &validationError{
			message: "Email must not exceed 320 characters.",
			kind:    AuthErrorTypeValidation,
		}
	case strings.Contains(email, ".."):
		return &validationError{
			message: "Please enter a valid email address.",
			kind:    AuthErrorTypeValidation,
		}
	case !emailRegex.MatchString(email):
		return &validationError{
			message: "Please enter a valid email address.",
			kind:    AuthErrorTypeValidation,
		}
	}

	switch {
	case strings.TrimSpace(request.Password) == "":
		return &validationError{"Password is required.", AuthErrorTypeValidation}
	case len(request.Password) < minPasswordLength:
		return &validationError{"Password must be at least 8 characters long.", AuthErrorTypeValidation}
	case len(request.Password) > maxPasswordLength:
		return &validationError{
			message: "Password must not exceed 512 characters.",
			kind:    AuthErrorTypeValidation,
		}
	case !passwordMeetsRequirements(request.Password):
		return &validationError{
			message: "Password must contain at least one uppercase letter, one lowercase letter, and one number.",
			kind:    AuthErrorTypeValidation,
		}
	}

	return nil
}

func validateLoginRequest(request LoginRequest) *validationError {
	username := strings.TrimSpace(request.Username)
	if username == "" {
		return &validationError{"Username is required.", AuthErrorTypeValidation}
	}
	if strings.TrimSpace(request.Password) == "" {
		return &validationError{"Password is required.", AuthErrorTypeValidation}
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

func unauthorizedResult() AuthResult {
	return AuthResult{
		Success:      false,
		ErrorMessage: "Invalid username or password. Please check your credentials and try again.",
		ErrorType:    AuthErrorTypeUnauthorized,
	}
}
