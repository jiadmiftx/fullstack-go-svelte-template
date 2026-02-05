package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Organization represents a Yayasan or Educational Institution.
// Conceptually: Organization = Yayasan (Foundation)
// Example: "Yayasan Pendidikan Islam" (code: YPI-001)
// Contains multiple Schools/Units (stored in units table)
type Organization struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey"`
	OwnerId     uuid.UUID `gorm:"type:uuid;not null;index"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Code        string    `gorm:"type:varchar(50);uniqueIndex;not null"`
	Type        string    `gorm:"type:varchar(50);not null"`
	Description string    `gorm:"type:text"`
	Address     string    `gorm:"type:text"`
	Logo        string    `gorm:"type:varchar(500)"`
	IsActive    bool      `gorm:"default:true"`
	Settings    string    `gorm:"type:jsonb;default:'{}'"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`

	Owner *User `gorm:"foreignKey:OwnerId"`
}

func (Organization) TableName() string { return "organizations" }

func (o *Organization) BeforeCreate(tx *gorm.DB) (err error) {
	if o.Id == uuid.Nil {
		o.Id = uuid.New()
	}
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return
}

func (o *Organization) BeforeUpdate(tx *gorm.DB) (err error) {
	o.UpdatedAt = time.Now()
	return
}
