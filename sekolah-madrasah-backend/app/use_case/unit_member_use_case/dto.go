package unit_member_use_case

import (
	"time"

	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
)

type UnitMember struct {
	Id        uuid.UUID              `json:"id"`
	UserId    uuid.UUID              `json:"user_id"`
	UnitId    uuid.UUID              `json:"unit_id"`
	Role      schemas.UnitMemberRole `json:"role"`
	IsActive  bool                   `json:"is_active"`
	JoinedAt  time.Time              `json:"joined_at"`
	InvitedBy *uuid.UUID             `json:"invited_by,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`

	User *UserInfo `json:"user,omitempty"`
	Unit *UnitInfo `json:"unit,omitempty"`
}

type UserInfo struct {
	Id       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Email    string    `json:"email"`
}

type UnitInfo struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Code string    `json:"code"`
	Type string    `json:"type"`
}

type AddMemberRequest struct {
	UserId uuid.UUID              `json:"user_id" binding:"required"`
	Role   schemas.UnitMemberRole `json:"role" binding:"required"`
}

type UpdateMemberRequest struct {
	Role     schemas.UnitMemberRole `json:"role,omitempty"`
	IsActive *bool                  `json:"is_active,omitempty"`
}

type UnitMemberFilter struct {
	UnitId   *uuid.UUID
	UserId   *uuid.UUID
	Role     *schemas.UnitMemberRole
	IsActive *bool
}
