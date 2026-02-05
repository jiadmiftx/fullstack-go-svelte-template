package organization_use_case

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	Id           uuid.UUID
	OwnerId      uuid.UUID
	Name         string
	Code         string
	Type         string
	Description  string
	Address      string
	Logo         string
	IsActive     bool
	Settings     string
	MemberCount  int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type OrganizationMember struct {
	Id             uuid.UUID
	UserId         uuid.UUID
	OrganizationId uuid.UUID
	RoleId         uuid.UUID
	IsActive       bool
	JoinedAt       time.Time
	InvitedBy      *uuid.UUID
}

type CreateOrganizationRequest struct {
	Name        string
	Code        string
	Type        string
	Description string
	Address     string
	Logo        string
}

type UpdateOrganizationRequest struct {
	Name        string
	Description string
	Address     string
	Logo        string
	Settings    string
}

type AddMemberRequest struct {
	UserId    uuid.UUID
	RoleId    uuid.UUID
	InvitedBy *uuid.UUID
}

type UpdateMemberRequest struct {
	RoleId   uuid.UUID
	IsActive bool
}

type OrganizationFilter struct {
	OwnerId  *uuid.UUID
	Type     *string
	IsActive *bool
}
