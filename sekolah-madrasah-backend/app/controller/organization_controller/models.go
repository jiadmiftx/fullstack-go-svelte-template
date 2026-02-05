package organization_controller

import (
	"time"

	"github.com/google/uuid"
)

type Organization struct {
	Id          uuid.UUID `json:"id"`
	OwnerId     uuid.UUID `json:"owner_id"`
	Name        string    `json:"name"`
	Code        string    `json:"code"`
	Type        string    `json:"type"`
	Description string    `json:"description,omitempty"`
	Address     string    `json:"address,omitempty"`
	Logo        string    `json:"logo,omitempty"`
	IsActive    bool      `json:"is_active"`
	MemberCount int       `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type OrganizationMember struct {
	Id             uuid.UUID  `json:"id"`
	UserId         uuid.UUID  `json:"user_id"`
	OrganizationId uuid.UUID  `json:"organization_id"`
	RoleId         uuid.UUID  `json:"role_id"`
	IsActive       bool       `json:"is_active"`
	JoinedAt       time.Time  `json:"joined_at"`
	InvitedBy      *uuid.UUID `json:"invited_by,omitempty"`
}

type CreateOrganizationRequest struct {
	Name        string `json:"name"`
	Code        string `json:"code"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Logo        string `json:"logo"`
}

type UpdateOrganizationRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Address     string `json:"address"`
	Logo        string `json:"logo"`
	Settings    string `json:"settings"`
}

type AddMemberRequest struct {
	UserId uuid.UUID `json:"user_id"`
	RoleId uuid.UUID `json:"role_id"`
}

type UpdateMemberRequest struct {
	RoleId   uuid.UUID `json:"role_id"`
	IsActive bool      `json:"is_active"`
}
