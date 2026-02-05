package permission_repository

import (
	"context"

	"github.com/google/uuid"
	"sekolah-madrasah/pkg/paginate_utils"
)

type PermissionRepository interface {
	GetPermission(ctx context.Context, filter PermissionFilter) (Permission, int, error)
	GetPermissions(ctx context.Context, filter PermissionFilter, paginate *paginate_utils.PaginateData) ([]Permission, int, error)
	CreatePermission(ctx context.Context, permission Permission) (Permission, int, error)
	DeletePermission(ctx context.Context, filter PermissionFilter) (int, error)

	GetRolePermissions(ctx context.Context, roleId uuid.UUID) ([]Permission, int, error)
	AssignPermissionToRole(ctx context.Context, roleId uuid.UUID, permissionId uuid.UUID) (int, error)
	RemovePermissionFromRole(ctx context.Context, roleId uuid.UUID, permissionId uuid.UUID) (int, error)
	SetRolePermissions(ctx context.Context, roleId uuid.UUID, permissionIds []uuid.UUID) (int, error)
}
