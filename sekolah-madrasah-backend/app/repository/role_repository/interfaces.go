package role_repository

import (
	"context"

	"sekolah-madrasah/pkg/paginate_utils"
)

type RoleRepository interface {
	GetRole(ctx context.Context, filter RoleFilter) (Role, int, error)
	GetRoles(ctx context.Context, filter RoleFilter, paginate *paginate_utils.PaginateData) ([]Role, int, error)
	CreateRole(ctx context.Context, role Role) (Role, int, error)
	UpdateRole(ctx context.Context, filter RoleFilter, role Role) (int, error)
	DeleteRole(ctx context.Context, filter RoleFilter) (int, error)
}
