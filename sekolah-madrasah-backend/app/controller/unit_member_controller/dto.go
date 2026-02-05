package unit_member_controller

import (
	"time"

	"sekolah-madrasah/app/use_case/unit_member_use_case"
	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
)

// DTOs for HTTP layer
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
	User      *UserInfo              `json:"user,omitempty"`
	Unit      *UnitInfo              `json:"unit,omitempty"`
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

func toMemberResponse(m unit_member_use_case.UnitMember) UnitMember {
	member := UnitMember{
		Id:        m.Id,
		UserId:    m.UserId,
		UnitId:    m.UnitId,
		Role:      m.Role,
		IsActive:  m.IsActive,
		JoinedAt:  m.JoinedAt,
		InvitedBy: m.InvitedBy,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	if m.User != nil {
		member.User = &UserInfo{
			Id:       m.User.Id,
			FullName: m.User.FullName,
			Email:    m.User.Email,
		}
	}

	if m.Unit != nil {
		member.Unit = &UnitInfo{
			Id:   m.Unit.Id,
			Name: m.Unit.Name,
			Code: m.Unit.Code,
			Type: m.Unit.Type,
		}
	}

	return member
}
