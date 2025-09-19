package app_test

import (
	"context"
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
	m.usersByEmail[clone.Email] = &clone

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
	_, ok := m.usersByEmail[email]
	return ok, nil
}

type stubTokenGenerator struct{}

func (stubTokenGenerator) GenerateToken(_ *authdomain.User) (string, error) {
	return "token-123", nil
}

func TestRegisterSuccess(t *testing.T) {
	repo := newMemoryUserRepository()
	hasher := authsecurity.NewHMACPasswordHasher()
	service := authapp.NewService(repo, hasher, stubTokenGenerator{})

	result, err := service.Register(context.Background(), authapp.RegisterRequest{
		Username: "new_user",
		Email:    "user@example.com",
		Password: "Password123",
	})

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
}

func TestRegisterDuplicateUsername(t *testing.T) {
	repo := newMemoryUserRepository()
	hasher := authsecurity.NewHMACPasswordHasher()
	service := authapp.NewService(repo, hasher, stubTokenGenerator{})

	_, err := service.Register(context.Background(), authapp.RegisterRequest{
		Username: "duplicate",
		Email:    "first@example.com",
		Password: "Password123",
	})
	if err != nil {
		t.Fatalf("expected first registration to succeed, got %v", err)
	}

	result, err := service.Register(context.Background(), authapp.RegisterRequest{
		Username: "duplicate",
		Email:    "second@example.com",
		Password: "Password123",
	})
	if err == nil {
		t.Fatalf("expected conflict error, got result %+v", result)
	}
	if !authapp.IsConflictError(err) {
		t.Fatalf("expected conflict error, got %v", err)
	}
}

func TestLoginSuccess(t *testing.T) {
	repo := newMemoryUserRepository()
	hasher := authsecurity.NewHMACPasswordHasher()
	service := authapp.NewService(repo, hasher, stubTokenGenerator{})

	_, err := service.Register(context.Background(), authapp.RegisterRequest{
		Username: "login_user",
		Email:    "login@example.com",
		Password: "Password123",
	})
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	result, err := service.Login(context.Background(), authapp.LoginRequest{
		Username: "login_user",
		Password: "Password123",
	})
	if err != nil {
		t.Fatalf("expected login to succeed, got %v", err)
	}
	if result == nil || result.Token == "" {
		t.Fatalf("expected login to return token")
	}
}

func TestLoginInvalidPassword(t *testing.T) {
	repo := newMemoryUserRepository()
	hasher := authsecurity.NewHMACPasswordHasher()
	service := authapp.NewService(repo, hasher, stubTokenGenerator{})

	_, err := service.Register(context.Background(), authapp.RegisterRequest{
		Username: "login_user",
		Email:    "login@example.com",
		Password: "Password123",
	})
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	result, err := service.Login(context.Background(), authapp.LoginRequest{
		Username: "login_user",
		Password: "WrongPassword",
	})
	if err == nil {
		t.Fatalf("expected unauthorized error, got result %+v", result)
	}
	if !authapp.IsUnauthorizedError(err) {
		t.Fatalf("expected unauthorized error, got %v", err)
	}
}
