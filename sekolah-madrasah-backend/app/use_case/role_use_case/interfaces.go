package role_use_case

import (
	"context"

	"github.com/google/uuid"
	"sekolah-madrasah/pkg/paginate_utils"
)

type RoleUseCase interface {
	GetRole(ctx context.Context, id uuid.UUID) (Role, int, error)
	GetRoles(ctx context.Context, filter RoleFilter, paginate *paginate_utils.PaginateData) ([]Role, int, error)
	CreateRole(ctx context.Context, req CreateRoleRequest) (Role, int, error)
	UpdateRole(ctx context.Context, id uuid.UUID, req UpdateRoleRequest) (Role, int, error)
	DeleteRole(ctx context.Context, id uuid.UUID) (int, error)
}
