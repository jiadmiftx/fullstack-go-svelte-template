package org_member_repository

import (
	"context"

	"sekolah-madrasah/pkg/paginate_utils"
)

type OrgMemberRepository interface {
	GetMember(ctx context.Context, filter OrgMemberFilter) (OrganizationMember, int, error)
	GetMembers(ctx context.Context, filter OrgMemberFilter, paginate *paginate_utils.PaginateData) ([]OrganizationMember, int, error)
	AddMember(ctx context.Context, member OrganizationMember) (OrganizationMember, int, error)
	UpdateMember(ctx context.Context, filter OrgMemberFilter, member OrganizationMember) (int, error)
	RemoveMember(ctx context.Context, filter OrgMemberFilter) (int, error)
}
