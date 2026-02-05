package permission_use_case

import (
	"context"

	"github.com/google/uuid"
	"sekolah-madrasah/pkg/paginate_utils"
)

type PermissionUseCase interface {
	GetPermission(ctx context.Context, id uuid.UUID) (Permission, int, error)
	GetPermissions(ctx context.Context, filter PermissionFilter, paginate *paginate_utils.PaginateData) ([]Permission, int, error)
	CreatePermission(ctx context.Context, req CreatePermissionRequest) (Permission, int, error)
	DeletePermission(ctx context.Context, id uuid.UUID) (int, error)
}
