package role_repository

import "github.com/google/uuid"

type RoleFilter struct {
	Id             *uuid.UUID
	OrganizationId *uuid.UUID
	Name           *string
	Type           *string
	IsDefault      *bool
	IsGlobal       *bool
}
