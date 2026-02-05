package organization_repository

import "github.com/google/uuid"

type OrganizationFilter struct {
	Id       *uuid.UUID
	OwnerId  *uuid.UUID
	Code     *string
	Type     *string
	IsActive *bool
}
