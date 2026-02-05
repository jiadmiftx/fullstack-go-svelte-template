package organization_repository

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	Id          uuid.UUID
	OwnerId     uuid.UUID
	Name        string
	Code        string
	Type        string
	Description string
	Address     string
	Logo        string
	IsActive    bool
	Settings    string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
