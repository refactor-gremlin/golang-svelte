package app_test

import (
	"context"
	"errors"
	"strings"
	"testing"

	authapp "mysvelteapp/server_new/internal/modules/auth/app"
	authdomain "mysvelteapp/server_new/internal/modules/auth/domain"
	authsecurity "mysvelteapp/server_new/internal/modules/auth/infra/security"
)

type memoryUserRepository struct {
	usersByUsername map[string]*authdomain.User
	usersByEmail    map[string]*authdomain.User
	nextID          uint
}

func newMemoryUserRepository() *memoryUserRepository {
	return &memoryUserRepository{
		usersByUsername: make(map[string]*authdomain.User),
		usersByEmail:    make(map[string]*authdomain.User),
		nextID:          1,
	}
}

func (m *memoryUserRepository) Add(_ context.Context, user *authdomain.User) error {
	clone := *user
	clone.ID = m.nextID
	m.nextID++

	m.usersByUsername[clone.Username] = &clone
	m.usersByEmail[strings.ToLower(clone.Email)] = &clone

	user.ID = clone.ID
	return nil
}

func (m *memoryUserRepository) GetByUsername(_ context.Context, username string) (*authdomain.User, error) {
	if user, ok := m.usersByUsername[username]; ok {
		clone := *user
		return &clone, nil
	}
	return nil, nil
}

func (m *memoryUserRepository) UsernameExists(_ context.Context, username string) (bool, error) {
	_, ok := m.usersByUsername[username]
	return ok, nil
}

func (m *memoryUserRepository) EmailExists(_ context.Context, email string) (bool, error) {
	_, ok := m.usersByEmail[strings.ToLower(email)]
	return ok, nil
}

type stubTokenGenerator struct{}

func (stubTokenGenerator) GenerateToken(_ *authdomain.User) (string, error) {
	return "token-123", nil
}

func newAuthService(repo *memoryUserRepository) *authapp.Service {
	hasher := authsecurity.NewHMACPasswordHasher()
	return authapp.NewService(repo, hasher, stubTokenGenerator{})
}

// TestRegisterSuccess validates the happy-path registration flow.
// Arrange: configure in-memory dependencies with a fresh auth service.
// Act: call Register with valid, mixed-case input.
// Assert: expect a token, persisted user ID, and normalised stored fields.
func TestRegisterSuccess(t *testing.T) {
	// Arrange
	repo := newMemoryUserRepository()
	service := newAuthService(repo)

	// Act
	result, err := service.Register(context.Background(), authapp.RegisterRequest{
		Username: "new_user",
		Email:    "NEW_user@example.COM ",
		Password: "Password123",
	})

	// Assert
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result == nil {
		t.Fatalf("expected result to be returned")
	}
	if result.Token == "" {
		t.Fatalf("expected token to be returned")
	}
	if result.UserID == 0 {
		t.Fatalf("expected user ID to be assigned")
	}
	if result.Username != "new_user" {
		t.Fatalf("expected username to be preserved, got %q", result.Username)
	}

	stored := repo.usersByUsername["new_user"]
	if stored == nil {
		t.Fatalf("expected user to be stored in repository")
	}
	if stored.Email != "new_user@example.com" {
		t.Fatalf("expected email to be normalised, got %q", stored.Email)
	}
	if stored.PasswordHash == "Password123" {
		t.Fatalf("expected password hash to differ from plain text")
	}
}

// TestRegisterDuplicateUsername ensures duplicate usernames are rejected.
// Arrange: seed a user in the repository.
// Act: attempt a second registration with the same username.
// Assert: expect a typed conflict error with the friendly message.
func TestRegisterDuplicateUsername(t *testing.T) {
	// Arrange
	repo := newMemoryUserRepository()
	service := newAuthService(repo)

	// Arrange (seed existing user)
	_, err := service.Register(context.Background(), authapp.RegisterRequest{
		Username: "duplicate",
		Email:    "first@example.com",
		Password: "Password123",
	})
	if err != nil {
		t.Fatalf("expected first registration to succeed, got %v", err)
	}

	// Act
	result, err := service.Register(context.Background(), authapp.RegisterRequest{
		Username: "duplicate",
		Email:    "second@example.com",
		Password: "Password123",
	})
	// Assert
	if err == nil {
		t.Fatalf("expected conflict error, got result %+v", result)
	}
	if !authapp.IsConflictError(err) {
		t.Fatalf("expected conflict error, got %v", err)
	}

	var conflict authapp.ConflictError
	if !errors.As(err, &conflict) {
		t.Fatalf("expected ConflictError type")
	}
	expected := "This username is already taken. Please choose a different one."
	if conflict.Message != expected {
		t.Fatalf("expected conflict message %q, got %q", expected, conflict.Message)
	}
}

// TestRegisterDuplicateEmail ensures email uniqueness is enforced case-insensitively.
// Arrange: seed a user whose email differs only by case.
// Act: register another user with the same email.
// Assert: expect a conflict error describing the duplicate email.
func TestRegisterDuplicateEmail(t *testing.T) {
	// Arrange
	repo := newMemoryUserRepository()
	service := newAuthService(repo)

	// Arrange (seed existing email)
	_, err := service.Register(context.Background(), authapp.RegisterRequest{
		Username: "first_user",
		Email:    "user@example.com",
		Password: "Password123",
	})
	if err != nil {
		t.Fatalf("expected first registration to succeed, got %v", err)
	}

	// Act
	result, err := service.Register(context.Background(), authapp.RegisterRequest{
		Username: "second_user",
		Email:    "USER@example.com",
		Password: "Password123",
	})
	// Assert
	if err == nil {
		t.Fatalf("expected conflict error, got result %+v", result)
	}
	if !authapp.IsConflictError(err) {
		t.Fatalf("expected conflict error, got %v", err)
	}

	var conflict authapp.ConflictError
	if !errors.As(err, &conflict) {
		t.Fatalf("expected ConflictError type")
	}
	expected := "This email is already registered. Please use a different email address."
	if conflict.Message != expected {
		t.Fatalf("expected conflict message %q, got %q", expected, conflict.Message)
	}
}

// TestRegisterValidationErrors covers validation failures for the register command.
// Arrange: table-drive invalid payloads.
// Act: invoke Register for each case.
// Assert: expect ValidationError with the matching message.
func TestRegisterValidationErrors(t *testing.T) {
	// Arrange
	repo := newMemoryUserRepository()
	service := newAuthService(repo)

	testCases := []struct {
		name    string
		payload authapp.RegisterRequest
		message string
	}{
		{
			name: "empty username",
			payload: authapp.RegisterRequest{
				Username: "   ",
				Email:    "user@example.com",
				Password: "Password123",
			},
			message: "Username is required.",
		},
		{
			name: "short username",
			payload: authapp.RegisterRequest{
				Username: "ab",
				Email:    "user@example.com",
				Password: "Password123",
			},
			message: "Username must be at least 3 characters long.",
		},
		{
			name: "long username",
			payload: authapp.RegisterRequest{
				Username: strings.Repeat("a", authdomain.MaxUsernameLength+1),
				Email:    "user@example.com",
				Password: "Password123",
			},
			message: "Username must not exceed 64 characters.",
		},
		{
			name: "invalid username characters",
			payload: authapp.RegisterRequest{
				Username: "bad-user",
				Email:    "user@example.com",
				Password: "Password123",
			},
			message: "Username can only contain letters, numbers, and underscores.",
		},
		{
			name: "empty email",
			payload: authapp.RegisterRequest{
				Username: "valid_user",
				Email:    "   ",
				Password: "Password123",
			},
			message: "Email is required.",
		},
		{
			name: "invalid email",
			payload: authapp.RegisterRequest{
				Username: "valid_user",
				Email:    "user@@example.com",
				Password: "Password123",
			},
			message: "Please enter a valid email address.",
		},
		{
			name: "empty password",
			payload: authapp.RegisterRequest{
				Username: "valid_user",
				Email:    "user@example.com",
				Password: "   ",
			},
			message: "Password is required.",
		},
		{
			name: "short password",
			payload: authapp.RegisterRequest{
				Username: "valid_user",
				Email:    "user@example.com",
				Password: "Abc123",
			},
			message: "Password must be at least 8 characters long.",
		},
		{
			name: "password missing complexity",
			payload: authapp.RegisterRequest{
				Username: "valid_user",
				Email:    "user@example.com",
				Password: "alllowercase",
			},
			message: "Password must contain at least one uppercase letter, one lowercase letter, and one number.",
		},
		{
			name: "password too long",
			payload: authapp.RegisterRequest{
				Username: "valid_user",
				Email:    "user@example.com",
				Password: strings.Repeat("A", 513),
			},
			message: "Password must not exceed 512 characters.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			payload := tc.payload

			// Act
			result, err := service.Register(context.Background(), payload)
			if err == nil {
				t.Fatalf("expected validation error, got result %+v", result)
			}

			// Assert
			if !authapp.IsValidationError(err) {
				t.Fatalf("expected validation error, got %v", err)
			}

			var validation authapp.ValidationError
			if !errors.As(err, &validation) {
				t.Fatalf("expected ValidationError type")
			}
			if validation.Message != tc.message {
				t.Fatalf("expected message %q, got %q", tc.message, validation.Message)
			}
		})
	}
}

// TestLoginSuccess proves valid credentials produce an auth token.
// Arrange: seed a registered user.
// Act: call Login with trimmed credentials.
// Assert: expect a token-bearing response and trimmed username.
func TestLoginSuccess(t *testing.T) {
	// Arrange
	repo := newMemoryUserRepository()
	service := newAuthService(repo)

	// Arrange (seed credentials)
	_, err := service.Register(context.Background(), authapp.RegisterRequest{
		Username: "login_user",
		Email:    "login@example.com",
		Password: "Password123",
	})
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	// Act
	result, err := service.Login(context.Background(), authapp.LoginRequest{
		Username: " login_user ",
		Password: "Password123",
	})

	// Assert
	if err != nil {
		t.Fatalf("expected login to succeed, got %v", err)
	}
	if result == nil || result.Token == "" {
		t.Fatalf("expected login to return token")
	}
	if result.Username != "login_user" {
		t.Fatalf("expected username to be trimmed, got %q", result.Username)
	}
}

// TestLoginInvalidPassword ensures incorrect passwords fail authentication.
// Arrange: register a known user.
// Act: Login with the wrong password.
// Assert: expect an UnauthorizedError with the friendly message.
func TestLoginInvalidPassword(t *testing.T) {
	// Arrange
	repo := newMemoryUserRepository()
	service := newAuthService(repo)

	// Arrange (seed user)
	_, err := service.Register(context.Background(), authapp.RegisterRequest{
		Username: "login_user",
		Email:    "login@example.com",
		Password: "Password123",
	})
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	// Act
	result, err := service.Login(context.Background(), authapp.LoginRequest{
		Username: "login_user",
		Password: "WrongPassword",
	})

	// Assert
	if err == nil {
		t.Fatalf("expected unauthorized error, got result %+v", result)
	}
	if !authapp.IsUnauthorizedError(err) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}

	var unauthorized authapp.UnauthorizedError
	if !errors.As(err, &unauthorized) {
		t.Fatalf("expected UnauthorizedError type")
	}
	expected := "Invalid username or password. Please check your credentials and try again."
	if unauthorized.Message != expected {
		t.Fatalf("expected unauthorized message %q, got %q", expected, unauthorized.Message)
	}
}

// TestLoginUnknownUser verifies missing accounts are treated as unauthorized.
// Arrange: create an empty repository.
// Act: attempt to log in with an unknown username.
// Assert: expect UnauthorizedError to be returned.
func TestLoginUnknownUser(t *testing.T) {
	// Arrange
	repo := newMemoryUserRepository()
	service := newAuthService(repo)

	// Act
	result, err := service.Login(context.Background(), authapp.LoginRequest{
		Username: "missing",
		Password: "Password123",
	})

	// Assert
	if err == nil {
		t.Fatalf("expected unauthorized error, got result %+v", result)
	}
	if !authapp.IsUnauthorizedError(err) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}

// TestLoginValidationErrors ensures login input validation mirrors production rules.
// Arrange: define invalid login payloads.
// Act: call Login for each case.
// Assert: expect ValidationError with the exact message.
func TestLoginValidationErrors(t *testing.T) {
	// Arrange
	repo := newMemoryUserRepository()
	service := newAuthService(repo)

	testCases := []struct {
		name    string
		payload authapp.LoginRequest
		message string
	}{
		{
			name: "empty username",
			payload: authapp.LoginRequest{
				Username: "   ",
				Password: "Password123",
			},
			message: "Username is required.",
		},
		{
			name: "empty password",
			payload: authapp.LoginRequest{
				Username: "valid",
				Password: "   ",
			},
			message: "Password is required.",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Arrange
			payload := tc.payload

			// Act
			result, err := service.Login(context.Background(), payload)
			if err == nil {
				t.Fatalf("expected validation error, got result %+v", result)
			}

			// Assert
			if !authapp.IsValidationError(err) {
				t.Fatalf("expected validation error, got %v", err)
			}

			var validation authapp.ValidationError
			if !errors.As(err, &validation) {
				t.Fatalf("expected ValidationError type")
			}
			if validation.Message != tc.message {
				t.Fatalf("expected message %q, got %q", tc.message, validation.Message)
			}
		})
	}
}
