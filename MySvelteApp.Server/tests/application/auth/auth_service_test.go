package auth_test

import (
	"context"
	"testing"

	appauth "mysvelteapp/server/internal/application/auth"
	"mysvelteapp/server/internal/domain/entities"
	"mysvelteapp/server/internal/infrastructure/security"
)

type memoryUserRepository struct {
	usersByUsername map[string]*entities.User
	usersByEmail    map[string]*entities.User
	nextID          uint
}

func newMemoryUserRepository() *memoryUserRepository {
	return &memoryUserRepository{
		usersByUsername: make(map[string]*entities.User),
		usersByEmail:    make(map[string]*entities.User),
		nextID:          1,
	}
}

func (m *memoryUserRepository) Add(_ context.Context, user *entities.User) error {
	clone := *user
	clone.ID = m.nextID
	m.nextID++

	m.usersByUsername[clone.Username] = &clone
	m.usersByEmail[clone.Email] = &clone

	user.ID = clone.ID
	return nil
}

func (m *memoryUserRepository) GetByUsername(_ context.Context, username string) (*entities.User, error) {
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
	_, ok := m.usersByEmail[email]
	return ok, nil
}

type stubTokenGenerator struct{}

func (stubTokenGenerator) GenerateToken(_ *entities.User) (string, error) {
	return "token-123", nil
}

func TestRegisterSuccess(t *testing.T) {
	repo := newMemoryUserRepository()
	hasher := security.NewHMACPasswordHasher()
	service := appauth.NewService(repo, hasher, stubTokenGenerator{})

	result, err := service.Register(context.Background(), appauth.RegisterRequest{
		Username: "new_user",
		Email:    "user@example.com",
		Password: "Password123",
	})

	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.Success {
		t.Fatalf("expected registration to succeed, got error %s", result.ErrorMessage)
	}
	if result.Token == "" {
		t.Fatalf("expected token to be returned")
	}
	if result.UserID == 0 {
		t.Fatalf("expected user ID to be assigned")
	}
}

func TestRegisterDuplicateUsername(t *testing.T) {
	repo := newMemoryUserRepository()
	hasher := security.NewHMACPasswordHasher()
	service := appauth.NewService(repo, hasher, stubTokenGenerator{})

	_, err := service.Register(context.Background(), appauth.RegisterRequest{
		Username: "duplicate",
		Email:    "first@example.com",
		Password: "Password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	result, err := service.Register(context.Background(), appauth.RegisterRequest{
		Username: "duplicate",
		Email:    "second@example.com",
		Password: "Password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Success {
		t.Fatalf("expected registration to fail for duplicate username")
	}
	if result.ErrorType != appauth.AuthErrorTypeConflict {
		t.Fatalf("expected conflict error, got %s", result.ErrorType)
	}
}

func TestLoginSuccess(t *testing.T) {
	repo := newMemoryUserRepository()
	hasher := security.NewHMACPasswordHasher()
	service := appauth.NewService(repo, hasher, stubTokenGenerator{})

	registerResult, err := service.Register(context.Background(), appauth.RegisterRequest{
		Username: "login_user",
		Email:    "login@example.com",
		Password: "Password123",
	})
	if err != nil || !registerResult.Success {
		t.Fatalf("registration failed: %v %v", err, registerResult)
	}

	result, err := service.Login(context.Background(), appauth.LoginRequest{
		Username: "login_user",
		Password: "Password123",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if !result.Success {
		t.Fatalf("expected login to succeed")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	repo := newMemoryUserRepository()
	hasher := security.NewHMACPasswordHasher()
	service := appauth.NewService(repo, hasher, stubTokenGenerator{})

	registerResult, err := service.Register(context.Background(), appauth.RegisterRequest{
		Username: "login_user",
		Email:    "login@example.com",
		Password: "Password123",
	})
	if err != nil || !registerResult.Success {
		t.Fatalf("registration failed: %v %v", err, registerResult)
	}

	result, err := service.Login(context.Background(), appauth.LoginRequest{
		Username: "login_user",
		Password: "WrongPassword",
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result.Success {
		t.Fatalf("expected login to fail with incorrect password")
	}
	if result.ErrorType != appauth.AuthErrorTypeUnauthorized {
		t.Fatalf("expected unauthorized error, got %s", result.ErrorType)
	}
}
