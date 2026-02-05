package schemas

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	Id          uuid.UUID `gorm:"type:uuid;primaryKey"`
	Name        string    `gorm:"type:varchar(100);uniqueIndex;not null"`
	Resource    string    `gorm:"type:varchar(50);not null;index"`
	Action      string    `gorm:"type:varchar(20);not null"`
	Description string    `gorm:"type:text"`
	CreatedAt   time.Time
}

func (Permission) TableName() string { return "permissions" }

func (p *Permission) BeforeCreate(tx interface{}) (err error) {
	if p.Id == uuid.Nil {
		p.Id = uuid.New()
	}
	p.CreatedAt = time.Now()
	return
}
