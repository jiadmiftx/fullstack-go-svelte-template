package role_repository

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	Id             uuid.UUID
	OrganizationId *uuid.UUID
	Name           string
	DisplayName    string
	Type           string
	Level          int
	Description    string
	IsDefault      bool
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
