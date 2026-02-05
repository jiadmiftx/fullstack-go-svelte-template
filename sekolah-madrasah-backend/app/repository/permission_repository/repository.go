package permission_repository

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"sekolah-madrasah/app/repository/common"
	"sekolah-madrasah/database/schemas"
	"sekolah-madrasah/pkg/paginate_utils"
	"gorm.io/gorm"
)

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) applyFilter(query *gorm.DB, filter PermissionFilter) *gorm.DB {
	if filter.Id != nil {
		query = query.Where("id = ?", *filter.Id)
	}
	if filter.Name != nil {
		query = query.Where("name = ?", *filter.Name)
	}
	if filter.Resource != nil {
		query = query.Where("resource = ?", *filter.Resource)
	}
	if filter.Action != nil {
		query = query.Where("action = ?", *filter.Action)
	}
	return query
}

func (r *permissionRepository) toModel(schema schemas.Permission) Permission {
	return Permission{
		Id:          schema.Id,
		Name:        schema.Name,
		Resource:    schema.Resource,
		Action:      schema.Action,
		Description: schema.Description,
		CreatedAt:   schema.CreatedAt,
	}
}

func (r *permissionRepository) GetPermission(ctx context.Context, filter PermissionFilter) (Permission, int, error) {
	var schema schemas.Permission
	query := r.db.WithContext(ctx).Model(&schemas.Permission{})
	query = r.applyFilter(query, filter)

	if err := query.First(&schema).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return Permission{}, code, err
	}

	return r.toModel(schema), http.StatusOK, nil
}

func (r *permissionRepository) GetPermissions(ctx context.Context, filter PermissionFilter, paginate *paginate_utils.PaginateData) ([]Permission, int, error) {
	var schemaList []schemas.Permission
	query := r.db.WithContext(ctx).Model(&schemas.Permission{})
	query = r.applyFilter(query, filter)

	if err := common.CountTotal(query, paginate); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	query = common.ApplyPagination(query, paginate)
	query = common.ApplyOrderBy(query, "resource ASC, action ASC")

	if err := query.Find(&schemaList).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return nil, code, err
	}

	permissions := make([]Permission, len(schemaList))
	for i, s := range schemaList {
		permissions[i] = r.toModel(s)
	}

	return permissions, http.StatusOK, nil
}

func (r *permissionRepository) CreatePermission(ctx context.Context, permission Permission) (Permission, int, error) {
	schema := schemas.Permission{
		Id:          permission.Id,
		Name:        permission.Name,
		Resource:    permission.Resource,
		Action:      permission.Action,
		Description: permission.Description,
	}

	if schema.Id == uuid.Nil {
		schema.Id = uuid.New()
	}
	schema.CreatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(&schema).Error; err != nil {
		return Permission{}, http.StatusInternalServerError, err
	}

	return r.toModel(schema), http.StatusCreated, nil
}

func (r *permissionRepository) DeletePermission(ctx context.Context, filter PermissionFilter) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.Permission{})
	query = r.applyFilter(query, filter)

	result := query.Delete(&schemas.Permission{})
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	if result.RowsAffected == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return http.StatusOK, nil
}

func (r *permissionRepository) GetRolePermissions(ctx context.Context, roleId uuid.UUID) ([]Permission, int, error) {
	var permissions []schemas.Permission

	err := r.db.WithContext(ctx).
		Table("permissions").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Where("role_permissions.role_id = ?", roleId).
		Find(&permissions).Error

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	result := make([]Permission, len(permissions))
	for i, p := range permissions {
		result[i] = r.toModel(p)
	}

	return result, http.StatusOK, nil
}

func (r *permissionRepository) AssignPermissionToRole(ctx context.Context, roleId uuid.UUID, permissionId uuid.UUID) (int, error) {
	var existing schemas.RolePermission
	err := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleId, permissionId).
		First(&existing).Error

	if err == nil {
		return http.StatusOK, nil
	}

	rolePermission := schemas.RolePermission{
		Id:           uuid.New(),
		RoleId:       roleId,
		PermissionId: permissionId,
		CreatedAt:    time.Now(),
	}

	if err := r.db.WithContext(ctx).Create(&rolePermission).Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusCreated, nil
}

func (r *permissionRepository) RemovePermissionFromRole(ctx context.Context, roleId uuid.UUID, permissionId uuid.UUID) (int, error) {
	result := r.db.WithContext(ctx).
		Where("role_id = ? AND permission_id = ?", roleId, permissionId).
		Delete(&schemas.RolePermission{})

	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}

	return http.StatusOK, nil
}

func (r *permissionRepository) SetRolePermissions(ctx context.Context, roleId uuid.UUID, permissionIds []uuid.UUID) (int, error) {
	tx := r.db.WithContext(ctx).Begin()

	if err := tx.Where("role_id = ?", roleId).Delete(&schemas.RolePermission{}).Error; err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, err
	}

	for _, permId := range permissionIds {
		rp := schemas.RolePermission{
			Id:           uuid.New(),
			RoleId:       roleId,
			PermissionId: permId,
			CreatedAt:    time.Now(),
		}
		if err := tx.Create(&rp).Error; err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, err
		}
	}

	if err := tx.Commit().Error; err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}
