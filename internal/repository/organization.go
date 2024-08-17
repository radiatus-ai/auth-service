// internal/repository/organization.go
package repository

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/radiatus-ai/auth-service/internal/model"
)

var (
	ErrOrganizationNotFound = errors.New("organization not found")
)

type OrganizationRepository interface {
	Create(org *model.Organization) error
	GetByID(id uuid.UUID) (*model.Organization, error)
	Update(org *model.Organization) error
	Delete(id uuid.UUID) error
	List() ([]model.Organization, error)
	AddUser(orgID, userID uuid.UUID) error
	RemoveUser(orgID, userID uuid.UUID) error
	GetUserOrganizations(userID uuid.UUID) ([]model.Organization, error)
	GetUserOrganization(userID uuid.UUID) (*model.Organization, error)
}

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) Create(org *model.Organization) error {
	return r.db.Create(org).Error
}

func (r *organizationRepository) GetByID(id uuid.UUID) (*model.Organization, error) {
	var org model.Organization
	if err := r.db.First(&org, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrganizationNotFound
		}
		return nil, err
	}
	return &org, nil
}

func (r *organizationRepository) Update(org *model.Organization) error {
	return r.db.Save(org).Error
}

func (r *organizationRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&model.Organization{}, id).Error
}

func (r *organizationRepository) List() ([]model.Organization, error) {
	var orgs []model.Organization
	if err := r.db.Find(&orgs).Error; err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *organizationRepository) AddUser(orgID, userID uuid.UUID) error {
	return r.db.Exec("INSERT INTO user_organizations (user_id, organization_id) VALUES (?, ?)", userID, orgID).Error
}

func (r *organizationRepository) RemoveUser(orgID, userID uuid.UUID) error {
	return r.db.Exec("DELETE FROM user_organizations WHERE user_id = ? AND organization_id = ?", userID, orgID).Error
}

func (r *organizationRepository) GetUserOrganizations(userID uuid.UUID) ([]model.Organization, error) {
	var orgs []model.Organization
	err := r.db.Joins("JOIN user_organizations ON user_organizations.organization_id = organizations.id").
		Where("user_organizations.user_id = ?", userID).
		Find(&orgs).Error
	if err != nil {
		return nil, err
	}
	return orgs, nil
}

func (r *organizationRepository) GetUserOrganization(userID uuid.UUID) (*model.Organization, error) {
	var org model.Organization
	err := r.db.
		Joins("JOIN user_organizations ON user_organizations.organization_id = organizations.id").
		Where("user_organizations.user_id = ?", userID).
		First(&org).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrOrganizationNotFound
		}
		return nil, err
	}
	return &org, nil
}
