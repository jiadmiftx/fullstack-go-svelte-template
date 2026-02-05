package permission_use_case

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"sekolah-madrasah/app/repository/permission_repository"
	"sekolah-madrasah/pkg/paginate_utils"
)

type permissionUseCase struct {
	permissionRepo permission_repository.PermissionRepository
}

func NewPermissionUseCase(permissionRepo permission_repository.PermissionRepository) PermissionUseCase {
	return &permissionUseCase{permissionRepo: permissionRepo}
}

func (u *permissionUseCase) toPermission(p permission_repository.Permission) Permission {
	return Permission{
		Id:          p.Id,
		Name:        p.Name,
		Resource:    p.Resource,
		Action:      p.Action,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
	}
}

func (u *permissionUseCase) GetPermission(ctx context.Context, id uuid.UUID) (Permission, int, error) {
	permission, code, err := u.permissionRepo.GetPermission(ctx, permission_repository.PermissionFilter{Id: &id})
	if err != nil {
		return Permission{}, code, err
	}

	return u.toPermission(permission), http.StatusOK, nil
}

func (u *permissionUseCase) GetPermissions(ctx context.Context, filter PermissionFilter, paginate *paginate_utils.PaginateData) ([]Permission, int, error) {
	repoFilter := permission_repository.PermissionFilter{
		Resource: filter.Resource,
		Action:   filter.Action,
	}

	permissions, code, err := u.permissionRepo.GetPermissions(ctx, repoFilter, paginate)
	if err != nil {
		return nil, code, err
	}

	result := make([]Permission, len(permissions))
	for i, p := range permissions {
		result[i] = u.toPermission(p)
	}

	return result, http.StatusOK, nil
}

func (u *permissionUseCase) CreatePermission(ctx context.Context, req CreatePermissionRequest) (Permission, int, error) {
	permName := fmt.Sprintf("%s.%s", req.Resource, req.Action)

	existingPerm, code, _ := u.permissionRepo.GetPermission(ctx, permission_repository.PermissionFilter{
		Name: &permName,
	})
	if code == http.StatusOK && existingPerm.Id != uuid.Nil {
		return Permission{}, http.StatusConflict, errors.New("permission already exists")
	}

	newPerm := permission_repository.Permission{
		Name:        permName,
		Resource:    req.Resource,
		Action:      req.Action,
		Description: req.Description,
	}

	createdPerm, code, err := u.permissionRepo.CreatePermission(ctx, newPerm)
	if err != nil {
		return Permission{}, code, err
	}

	return u.toPermission(createdPerm), http.StatusCreated, nil
}

func (u *permissionUseCase) DeletePermission(ctx context.Context, id uuid.UUID) (int, error) {
	_, code, err := u.permissionRepo.GetPermission(ctx, permission_repository.PermissionFilter{Id: &id})
	if err != nil {
		return code, err
	}

	return u.permissionRepo.DeletePermission(ctx, permission_repository.PermissionFilter{Id: &id})
}
