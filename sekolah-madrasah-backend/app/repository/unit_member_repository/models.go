package unit_member_repository

import (
	"time"

	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
)

type UnitMember struct {
	Id        uuid.UUID
	UserId    uuid.UUID
	UnitId    uuid.UUID
	Role      schemas.UnitMemberRole
	IsActive  bool
	JoinedAt  time.Time
	InvitedBy *uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time

	// Preloaded relations
	User *UserInfo
	Unit *UnitInfo
}

type UserInfo struct {
	Id       uuid.UUID
	FullName string
	Email    string
}

type UnitInfo struct {
	Id   uuid.UUID
	Name string
	Code string
	Type string
}
