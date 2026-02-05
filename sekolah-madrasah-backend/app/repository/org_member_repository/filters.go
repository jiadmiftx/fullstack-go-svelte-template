package org_member_repository

import "github.com/google/uuid"

type OrgMemberFilter struct {
	Id             *uuid.UUID
	UserId         *uuid.UUID
	OrganizationId *uuid.UUID
	RoleId         *uuid.UUID
	IsActive       *bool
}
