package auth

import (
	"testing"

	"github.com/radiatus-ai/auth-service/internal/model"
	"github.com/radiatus-ai/auth-service/internal/repository"
)

type mockUserRepository struct {
	users map[string]*model.User
}

func (m *mockUserRepository) Create(user *model.User) error {
	if _, exists := m.users[user.Email]; exists {
		return ErrUserAlreadyExists
	}
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepository) GetByEmail(email string) (*model.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, repository.ErrUserNotFound
	}
	return user, nil
}

func (m *mockUserRepository) ExistsByEmail(email string) (bool, error) {
	_, exists := m.users[email]
	return exists, nil
}

func TestRegister(t *testing.T) {
	repo := &mockUserRepository{users: make(map[string]*model.User)}
	service := NewService(repo, "test-secret")

	err := service.Register("test@example.com", "password")
	if err != nil {
		t.Errorf("Expected successful registration, got error: %v", err)
	}

	err = service.Register("test@example.com", "password")
	if err != ErrUserAlreadyExists {
		t.Errorf("Expected ErrUserAlreadyExists, got: %v", err)
	}
}

func TestLogin(t *testing.T) {
	repo := &mockUserRepository{users: make(map[string]*model.User)}
	service := NewService(repo, "test-secret")

	// Register a user first
	_ = service.Register("test@example.com", "password")

	token, err := service.Login("test@example.com", "password")
	if err != nil {
		t.Errorf("Expected successful login, got error: %v", err)
	}
	if token == "" {
		t.Error("Expected non-empty token")
	}

	_, err = service.Login("test@example.com", "wrong-password")
	if err != ErrInvalidCredentials {
		t.Errorf("Expected ErrInvalidCredentials, got: %v", err)
	}

	_, err = service.Login("nonexistent@example.com", "password")
	if err != repository.ErrUserNotFound {
		t.Errorf("Expected ErrUserNotFound, got: %v", err)
	}
}

func TestVerifyToken(t *testing.T) {
	repo := &mockUserRepository{users: make(map[string]*model.User)}
	service := NewService(repo, "test-secret")

	// Register and login a user to get a token
	_ = service.Register("test@example.com", "password")
	token, _ := service.Login("test@example.com", "password")

	userID, err := service.VerifyToken(token)
	if err != nil {
		t.Errorf("Expected successful token verification, got error: %v", err)
	}
	if userID == "" {
		t.Error("Expected non-empty user ID")
	}

	_, err = service.VerifyToken("invalid-token")
	if err == nil {
		t.Error("Expected error for invalid token, got nil")
	}
}
