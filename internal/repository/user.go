package repository

import (
	"errors"
	"log"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/radiatus-ai/auth-service/internal/model"
)

var (
	ErrUserNotFound = errors.New("user not found")
)

type UserRepository interface {
	Create(user *model.User) error
	GetByEmail(email string) (*model.User, error)
	GetByID(id uuid.UUID) (*model.User, error)
	GetByGoogleID(googleID string) (*model.User, error)
	ExistsByEmail(email string) (bool, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *model.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		log.Printf("SQL: %v\n", result.Statement.SQL.String())
		log.Printf("Vars: %v\n", result.Statement.Vars)
		return result.Error
	}
	return result.Error
}

func (r *userRepository) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByID(id uuid.UUID) (*model.User, error) {
	var user model.User
	if err := r.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) GetByGoogleID(googleID string) (*model.User, error) {
	var user model.User
	if err := r.db.Where("google_id = ?", googleID).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) ExistsByEmail(email string) (bool, error) {
	var count int64
	if err := r.db.Model(&model.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}
