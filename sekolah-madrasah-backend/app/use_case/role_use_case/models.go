package role_use_case

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
	Permissions    []Permission
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type Permission struct {
	Id          uuid.UUID
	Name        string
	Resource    string
	Action      string
	Description string
}

type CreateRoleRequest struct {
	OrganizationId *uuid.UUID
	Name           string
	DisplayName    string
	Type           string
	Level          int
	Description    string
	IsDefault      bool
	PermissionIds  []uuid.UUID
}

type UpdateRoleRequest struct {
	Name          string
	DisplayName   string
	Description   string
	Level         int
	PermissionIds []uuid.UUID
}

type RoleFilter struct {
	Id             *uuid.UUID
	OrganizationId *uuid.UUID
	Name           *string
	Type           *string
	IsGlobal       *bool
}
