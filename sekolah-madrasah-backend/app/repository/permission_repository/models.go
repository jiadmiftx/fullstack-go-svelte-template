package permission_repository

import (
	"time"

	"github.com/google/uuid"
)

type Permission struct {
	Id          uuid.UUID
	Name        string
	Resource    string
	Action      string
	Description string
	CreatedAt   time.Time
}

type RolePermission struct {
	Id           uuid.UUID
	RoleId       uuid.UUID
	PermissionId uuid.UUID
	CreatedAt    time.Time
}
