package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Unit represents a School or Educational Unit within an Organization.
// Example: "SMA Negeri 1", "SD Islam Al-Azhar" under organization "Yayasan Al-Azhar"
type Unit struct {
	Id             uuid.UUID `gorm:"type:uuid;primaryKey"`
	OrganizationId uuid.UUID `gorm:"type:uuid;not null;index"`
	Name           string    `gorm:"type:varchar(255);not null"`
	Code           string    `gorm:"type:varchar(50);uniqueIndex;not null"` // Unique unit code (e.g. NPSN)
	Slug           string    `gorm:"type:varchar(100);uniqueIndex"`         // URL slug for public registration
	Type           string    `gorm:"type:varchar(50);not null"`             // SD/SMP/SMA/TK/Madrasah
	Address        string    `gorm:"type:text"`
	Phone          string    `gorm:"type:varchar(20)"`
	Email          string    `gorm:"type:varchar(255)"`
	Logo           string    `gorm:"type:varchar(500)"`
	IsActive       bool      `gorm:"default:true"`
	Settings       string    `gorm:"type:jsonb;default:'{}'"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	DeletedAt      gorm.DeletedAt `gorm:"index"`

	Organization *Organization `gorm:"foreignKey:OrganizationId"`
}

func (Unit) TableName() string { return "units" }

func (s *Unit) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Id == uuid.Nil {
		s.Id = uuid.New()
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return
}

func (s *Unit) BeforeUpdate(tx *gorm.DB) (err error) {
	s.UpdatedAt = time.Now()
	return
}
