package role_use_case

import (
	"context"
	"errors"
	"net/http"

	"sekolah-madrasah/app/repository/permission_repository"
	"sekolah-madrasah/app/repository/role_repository"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
)

type roleUseCase struct {
	roleRepo       role_repository.RoleRepository
	permissionRepo permission_repository.PermissionRepository
}

func NewRoleUseCase(roleRepo role_repository.RoleRepository, permissionRepo permission_repository.PermissionRepository) RoleUseCase {
	return &roleUseCase{
		roleRepo:       roleRepo,
		permissionRepo: permissionRepo,
	}
}

func (u *roleUseCase) toRole(r role_repository.Role, permissions []permission_repository.Permission) Role {
	perms := make([]Permission, len(permissions))
	for i, p := range permissions {
		perms[i] = Permission{
			Id:          p.Id,
			Name:        p.Name,
			Resource:    p.Resource,
			Action:      p.Action,
			Description: p.Description,
		}
	}

	return Role{
		Id:             r.Id,
		OrganizationId: r.OrganizationId,
		Name:           r.Name,
		DisplayName:    r.DisplayName,
		Type:           r.Type,
		Level:          r.Level,
		Description:    r.Description,
		IsDefault:      r.IsDefault,
		Permissions:    perms,
		CreatedAt:      r.CreatedAt,
		UpdatedAt:      r.UpdatedAt,
	}
}

func (u *roleUseCase) GetRole(ctx context.Context, id uuid.UUID) (Role, int, error) {
	role, code, err := u.roleRepo.GetRole(ctx, role_repository.RoleFilter{Id: &id})
	if err != nil {
		return Role{}, code, err
	}

	permissions, _, _ := u.permissionRepo.GetRolePermissions(ctx, id)

	return u.toRole(role, permissions), http.StatusOK, nil
}

func (u *roleUseCase) GetRoles(ctx context.Context, filter RoleFilter, paginate *paginate_utils.PaginateData) ([]Role, int, error) {
	repoFilter := role_repository.RoleFilter{
		Id:             filter.Id,
		OrganizationId: filter.OrganizationId,
		Name:           filter.Name,
		Type:           filter.Type,
		IsGlobal:       filter.IsGlobal,
	}

	roles, code, err := u.roleRepo.GetRoles(ctx, repoFilter, paginate)
	if err != nil {
		return nil, code, err
	}

	result := make([]Role, len(roles))
	for i, r := range roles {
		permissions, _, _ := u.permissionRepo.GetRolePermissions(ctx, r.Id)
		result[i] = u.toRole(r, permissions)
	}

	return result, http.StatusOK, nil
}

func (u *roleUseCase) CreateRole(ctx context.Context, req CreateRoleRequest) (Role, int, error) {
	existingRole, code, _ := u.roleRepo.GetRole(ctx, role_repository.RoleFilter{
		Name:           &req.Name,
		OrganizationId: req.OrganizationId,
	})
	if code == http.StatusOK && existingRole.Id != uuid.Nil {
		return Role{}, http.StatusConflict, errors.New("role with this name already exists")
	}

	roleType := req.Type
	if roleType == "" {
		roleType = "custom"
	}

	newRole := role_repository.Role{
		OrganizationId: req.OrganizationId,
		Name:           req.Name,
		DisplayName:    req.DisplayName,
		Type:           roleType,
		Level:          req.Level,
		Description:    req.Description,
		IsDefault:      req.IsDefault,
	}

	createdRole, code, err := u.roleRepo.CreateRole(ctx, newRole)
	if err != nil {
		return Role{}, code, err
	}

	// Always set permissions (even if empty)
	u.permissionRepo.SetRolePermissions(ctx, createdRole.Id, req.PermissionIds)

	permissions, _, _ := u.permissionRepo.GetRolePermissions(ctx, createdRole.Id)

	return u.toRole(createdRole, permissions), http.StatusCreated, nil
}

func (u *roleUseCase) UpdateRole(ctx context.Context, id uuid.UUID, req UpdateRoleRequest) (Role, int, error) {
	existingRole, code, err := u.roleRepo.GetRole(ctx, role_repository.RoleFilter{Id: &id})
	if err != nil {
		return Role{}, code, err
	}

	if existingRole.Type == "system" {
		return Role{}, http.StatusForbidden, errors.New("cannot modify system role")
	}

	updateData := role_repository.Role{
		Name:        req.Name,
		DisplayName: req.DisplayName,
		Description: req.Description,
		Level:       req.Level,
	}

	code, err = u.roleRepo.UpdateRole(ctx, role_repository.RoleFilter{Id: &id}, updateData)
	if err != nil {
		return Role{}, code, err
	}

	// Always update permissions when provided (even if empty to clear all)
	u.permissionRepo.SetRolePermissions(ctx, id, req.PermissionIds)

	updatedRole, code, err := u.roleRepo.GetRole(ctx, role_repository.RoleFilter{Id: &id})
	if err != nil {
		return Role{}, code, err
	}

	permissions, _, _ := u.permissionRepo.GetRolePermissions(ctx, id)

	return u.toRole(updatedRole, permissions), http.StatusOK, nil
}

func (u *roleUseCase) DeleteRole(ctx context.Context, id uuid.UUID) (int, error) {
	role, code, err := u.roleRepo.GetRole(ctx, role_repository.RoleFilter{Id: &id})
	if err != nil {
		return code, err
	}

	if role.Type == "system" {
		return http.StatusForbidden, errors.New("cannot delete system role")
	}

	u.permissionRepo.SetRolePermissions(ctx, id, []uuid.UUID{})

	return u.roleRepo.DeleteRole(ctx, role_repository.RoleFilter{Id: &id})
}
