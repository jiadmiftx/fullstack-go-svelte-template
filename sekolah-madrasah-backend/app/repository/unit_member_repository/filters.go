package unit_member_repository

import (
	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
)

type UnitMemberFilter struct {
	Id       *uuid.UUID
	UserId   *uuid.UUID
	UnitId *uuid.UUID
	Role     *schemas.UnitMemberRole
	IsActive *bool
}
