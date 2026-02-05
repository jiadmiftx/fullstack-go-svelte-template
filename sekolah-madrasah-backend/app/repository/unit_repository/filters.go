package unit_repository

import "github.com/google/uuid"

type UnitFilter struct {
	Id             *uuid.UUID
	OrganizationId *uuid.UUID
	Code           *string
	Type           *string
	IsActive       *bool
}
