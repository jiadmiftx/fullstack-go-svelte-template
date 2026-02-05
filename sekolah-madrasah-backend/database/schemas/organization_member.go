package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationMember struct {
	Id             uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserId         uuid.UUID      `gorm:"type:uuid;not null;index"`
	OrganizationId uuid.UUID      `gorm:"type:uuid;not null;index"`
	RoleId         uuid.UUID      `gorm:"type:uuid;not null;index"`
	IsActive       bool           `gorm:"default:true"`
	JoinedAt       time.Time
	InvitedBy      *uuid.UUID     `gorm:"type:uuid"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	User         *User         `gorm:"foreignKey:UserId"`
	Organization *Organization `gorm:"foreignKey:OrganizationId"`
	Role         *Role         `gorm:"foreignKey:RoleId"`
	Inviter      *User         `gorm:"foreignKey:InvitedBy"`
}

func (OrganizationMember) TableName() string { return "organization_members" }

func (om *OrganizationMember) BeforeCreate(tx *gorm.DB) (err error) {
	if om.Id == uuid.Nil {
		om.Id = uuid.New()
	}
	om.JoinedAt = time.Now()
	om.CreatedAt = time.Now()
	om.UpdatedAt = time.Now()
	return
}

func (om *OrganizationMember) BeforeUpdate(tx *gorm.DB) (err error) {
	om.UpdatedAt = time.Now()
	return
}
