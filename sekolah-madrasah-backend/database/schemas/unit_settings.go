package schemas

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UnitSettings stores unit settings like period duration
type UnitSettings struct {
	Id               uuid.UUID `gorm:"type:uuid;primaryKey"`
	UnitId           uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	PeriodDuration   int       `gorm:"type:int;default:40"`              // Duration per period in minutes (default 40)
	StartTime        string    `gorm:"type:varchar(10);default:'07:00'"` // Start time of day
	TotalPeriods     int       `gorm:"type:int;default:9"`               // Total periods per day
	BreakAfterPeriod int       `gorm:"type:int;default:3"`               // Break after period n
	BreakDuration    int       `gorm:"type:int;default:15"`              // Break duration (minutes)

	// Semester Settings
	AcademicYear    string     `gorm:"type:varchar(20)"`   // "2025/2026"
	CurrentSemester int        `gorm:"type:int;default:1"` // 1 or 2
	Semester1Start  *time.Time `gorm:"type:date"`          // Semester 1 start
	Semester1End    *time.Time `gorm:"type:date"`          // Semester 1 end
	Semester2Start  *time.Time `gorm:"type:date"`          // Semester 2 start
	Semester2End    *time.Time `gorm:"type:date"`          // Semester 2 end

	CreatedAt time.Time
	UpdatedAt time.Time

	Unit *Unit `gorm:"foreignKey:UnitId"`
}

func (UnitSettings) TableName() string { return "unit_settings" }

func (s *UnitSettings) BeforeCreate(tx *gorm.DB) (err error) {
	if s.Id == uuid.Nil {
		s.Id = uuid.New()
	}
	s.CreatedAt = time.Now()
	s.UpdatedAt = time.Now()
	return
}

func (s *UnitSettings) BeforeUpdate(tx *gorm.DB) (err error) {
	s.UpdatedAt = time.Now()
	return
}
