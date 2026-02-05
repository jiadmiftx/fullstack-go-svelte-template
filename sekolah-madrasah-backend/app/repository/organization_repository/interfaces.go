package organization_repository

import (
	"context"

	"sekolah-madrasah/pkg/paginate_utils"
)

type OrganizationRepository interface {
	GetOrganization(ctx context.Context, filter OrganizationFilter) (Organization, int, error)
	GetOrganizations(ctx context.Context, filter OrganizationFilter, paginate *paginate_utils.PaginateData) ([]Organization, int, error)
	CreateOrganization(ctx context.Context, org Organization) (Organization, int, error)
	UpdateOrganization(ctx context.Context, filter OrganizationFilter, org Organization) (int, error)
	DeleteOrganization(ctx context.Context, filter OrganizationFilter) (int, error)
}
