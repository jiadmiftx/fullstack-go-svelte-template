package organization_use_case

import (
	"context"

	"github.com/google/uuid"
	"sekolah-madrasah/pkg/paginate_utils"
)

type OrganizationUseCase interface {
	GetOrganization(ctx context.Context, id uuid.UUID) (Organization, int, error)
	GetOrganizations(ctx context.Context, filter OrganizationFilter, paginate *paginate_utils.PaginateData) ([]Organization, int, error)
	CreateOrganization(ctx context.Context, ownerId uuid.UUID, req CreateOrganizationRequest) (Organization, int, error)
	UpdateOrganization(ctx context.Context, id uuid.UUID, req UpdateOrganizationRequest) (Organization, int, error)
	DeleteOrganization(ctx context.Context, id uuid.UUID) (int, error)

	GetMembers(ctx context.Context, orgId uuid.UUID, paginate *paginate_utils.PaginateData) ([]OrganizationMember, int, error)
	AddMember(ctx context.Context, orgId uuid.UUID, req AddMemberRequest) (OrganizationMember, int, error)
	UpdateMember(ctx context.Context, orgId uuid.UUID, userId uuid.UUID, req UpdateMemberRequest) (int, error)
	RemoveMember(ctx context.Context, orgId uuid.UUID, userId uuid.UUID) (int, error)
}
