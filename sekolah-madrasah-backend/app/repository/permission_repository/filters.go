package permission_repository

import "github.com/google/uuid"

type PermissionFilter struct {
	Id       *uuid.UUID
	Name     *string
	Resource *string
	Action   *string
}

type RolePermissionFilter struct {
	Id           *uuid.UUID
	RoleId       *uuid.UUID
	PermissionId *uuid.UUID
}
