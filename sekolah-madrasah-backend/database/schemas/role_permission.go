package schemas

import (
	"time"

	"github.com/google/uuid"
)

type RolePermission struct {
	Id           uuid.UUID `gorm:"type:uuid;primaryKey"`
	RoleId       uuid.UUID `gorm:"type:uuid;not null;index"`
	PermissionId uuid.UUID `gorm:"type:uuid;not null;index"`
	CreatedAt    time.Time

	Role       *Role       `gorm:"foreignKey:RoleId"`
	Permission *Permission `gorm:"foreignKey:PermissionId"`
}

func (RolePermission) TableName() string { return "role_permissions" }

func (rp *RolePermission) BeforeCreate(tx interface{}) (err error) {
	if rp.Id == uuid.Nil {
		rp.Id = uuid.New()
	}
	rp.CreatedAt = time.Now()
	return
}
