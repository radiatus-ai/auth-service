package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	ID            uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	Email         string         `gorm:"unique;not null" json:"email"`
	GoogleID      string         `gorm:"unique" json:"google_id,omitempty"`
	Password      string         `json:"-"` // Excluded from JSON output
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	Organizations []Organization `gorm:"many2many:user_organizations;" json:"organizations,omitempty"`
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return nil
}
