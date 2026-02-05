package permission_use_case

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

type CreatePermissionRequest struct {
	Name        string
	Resource    string
	Action      string
	Description string
}

type PermissionFilter struct {
	Resource *string
	Action   *string
}
