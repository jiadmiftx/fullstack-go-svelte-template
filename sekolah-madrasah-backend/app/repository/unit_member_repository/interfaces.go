package unit_member_repository

import (
	"context"

	"sekolah-madrasah/pkg/paginate_utils"
)

type UnitMemberRepository interface {
	GetMember(ctx context.Context, filter UnitMemberFilter) (UnitMember, int, error)
	GetMembers(ctx context.Context, filter UnitMemberFilter, paginate *paginate_utils.PaginateData) ([]UnitMember, int, error)
	AddMember(ctx context.Context, member UnitMember) (UnitMember, int, error)
	UpdateMember(ctx context.Context, filter UnitMemberFilter, member UnitMember) (int, error)
	RemoveMember(ctx context.Context, filter UnitMemberFilter) (int, error)
}
