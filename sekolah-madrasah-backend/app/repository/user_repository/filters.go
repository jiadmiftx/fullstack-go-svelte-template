package user_repository

import "github.com/google/uuid"

type UserFilter struct {
	Id           *uuid.UUID
	Email        *string
	IsSuperAdmin *bool
	IsActive     *bool
	PlatformOnly *bool // Filter to show only platform-level users (super admin or unassigned)
}
