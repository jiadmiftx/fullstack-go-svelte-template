package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	Id              uuid.UUID      `gorm:"type:uuid;primaryKey"`
	Email           string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Password        string         `gorm:"type:varchar(255);not null"`
	FullName        string         `gorm:"type:varchar(100)"`
	Phone           string         `gorm:"type:varchar(20)"`
	Avatar          string         `gorm:"type:varchar(500)"`
	IsSuperAdmin    bool           `gorm:"default:false"`
	IsActive        bool           `gorm:"default:true"`
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
	DeletedAt       gorm.DeletedAt `gorm:"index"`
}

func (User) TableName() string { return "users" }

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.Id == uuid.Nil {
		u.Id = uuid.New()
	}
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return
}

func (u *User) BeforeUpdate(tx *gorm.DB) (err error) {
	u.UpdatedAt = time.Now()
	return
}
