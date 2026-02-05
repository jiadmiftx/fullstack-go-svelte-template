package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UnitMemberRole defines the role of a user within a unit
type UnitMemberRole string

const (
	UnitMemberRoleAdmin    UnitMemberRole = "admin"
	UnitMemberRolePengurus UnitMemberRole = "pengurus"
	UnitMemberRoleWarga    UnitMemberRole = "warga"
	UnitMemberRoleParent   UnitMemberRole = "parent"
	UnitMemberRoleStaff    UnitMemberRole = "staff"
)

type UnitMember struct {
	Id        uuid.UUID      `gorm:"type:uuid;primaryKey"`
	UserId    uuid.UUID      `gorm:"type:uuid;not null;index"`
	UnitId    uuid.UUID      `gorm:"type:uuid;not null;index"`
	Role      UnitMemberRole `gorm:"type:varchar(20);not null;default:'staff'"`
	IsActive  bool           `gorm:"default:true"`
	JoinedAt  time.Time
	InvitedBy *uuid.UUID `gorm:"type:uuid"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	User    *User `gorm:"foreignKey:UserId"`
	Unit    *Unit `gorm:"foreignKey:UnitId"`
	Inviter *User `gorm:"foreignKey:InvitedBy"`
}

func (UnitMember) TableName() string { return "unit_members" }

func (sm *UnitMember) BeforeCreate(tx *gorm.DB) (err error) {
	if sm.Id == uuid.Nil {
		sm.Id = uuid.New()
	}
	sm.JoinedAt = time.Now()
	sm.CreatedAt = time.Now()
	sm.UpdatedAt = time.Now()
	return
}

func (sm *UnitMember) BeforeUpdate(tx *gorm.DB) (err error) {
	sm.UpdatedAt = time.Now()
	return
}
