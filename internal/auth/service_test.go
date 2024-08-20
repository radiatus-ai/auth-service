package auth

import (
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
