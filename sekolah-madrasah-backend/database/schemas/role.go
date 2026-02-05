package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Role struct {
	Id             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	OrganizationId *uuid.UUID     `gorm:"type:uuid;index"`
	Name           string         `gorm:"type:varchar(50);not null"`
	DisplayName    string         `gorm:"type:varchar(100);not null"`
	Type           string         `gorm:"type:varchar(20);not null;default:'custom'"`
	Level          int            `gorm:"not null;default:0"`
	Description    string         `gorm:"type:text"`
	IsDefault      bool           `gorm:"default:false"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	Organization *Organization `gorm:"foreignKey:OrganizationId"`
}

func (Role) TableName() string { return "roles" }

func (r *Role) BeforeCreate(tx *gorm.DB) (err error) {
	if r.Id == uuid.Nil {
		r.Id = uuid.New()
	}
	r.CreatedAt = time.Now()
	r.UpdatedAt = time.Now()
	return
}

func (r *Role) BeforeUpdate(tx *gorm.DB) (err error) {
	r.UpdatedAt = time.Now()
	return
}
