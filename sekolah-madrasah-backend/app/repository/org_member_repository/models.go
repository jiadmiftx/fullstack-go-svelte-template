package org_member_repository

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationMember struct {
	Id             uuid.UUID
	UserId         uuid.UUID
	OrganizationId uuid.UUID
	RoleId         uuid.UUID
	IsActive       bool
	JoinedAt       time.Time
	InvitedBy      *uuid.UUID
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
