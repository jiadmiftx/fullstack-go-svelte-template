package unit_member_use_case

import (
	"context"

	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
)

type UnitMemberUseCase interface {
	GetMember(ctx context.Context, perumahanId, memberId uuid.UUID) (UnitMember, int, error)
	GetMembers(ctx context.Context, perumahanId uuid.UUID, filter UnitMemberFilter, paginate *paginate_utils.PaginateData) ([]UnitMember, int, error)
	AddMember(ctx context.Context, perumahanId uuid.UUID, req AddMemberRequest, invitedBy *uuid.UUID) (UnitMember, int, error)
	UpdateMember(ctx context.Context, perumahanId, memberId uuid.UUID, req UpdateMemberRequest) (UnitMember, int, error)
	RemoveMember(ctx context.Context, perumahanId, memberId uuid.UUID) (int, error)
}
