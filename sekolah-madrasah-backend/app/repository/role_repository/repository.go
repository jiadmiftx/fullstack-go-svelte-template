package role_repository

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

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) applyFilter(query *gorm.DB, filter RoleFilter) *gorm.DB {
	if filter.Id != nil {
		query = query.Where("id = ?", *filter.Id)
	}
	if filter.OrganizationId != nil {
		query = query.Where("organization_id = ?", *filter.OrganizationId)
	}
	if filter.Name != nil {
		query = query.Where("name = ?", *filter.Name)
	}
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.IsDefault != nil {
		query = query.Where("is_default = ?", *filter.IsDefault)
	}
	if filter.IsGlobal != nil && *filter.IsGlobal {
		query = query.Where("organization_id IS NULL")
	}
	return query
}

func (r *roleRepository) toModel(schema schemas.Role) Role {
	return Role{
		Id:             schema.Id,
		OrganizationId: schema.OrganizationId,
		Name:           schema.Name,
		DisplayName:    schema.DisplayName,
		Type:           schema.Type,
		Level:          schema.Level,
		Description:    schema.Description,
		IsDefault:      schema.IsDefault,
		CreatedAt:      schema.CreatedAt,
		UpdatedAt:      schema.UpdatedAt,
	}
}

func (r *roleRepository) toSchema(model Role) schemas.Role {
	return schemas.Role{
		Id:             model.Id,
		OrganizationId: model.OrganizationId,
		Name:           model.Name,
		DisplayName:    model.DisplayName,
		Type:           model.Type,
		Level:          model.Level,
		Description:    model.Description,
		IsDefault:      model.IsDefault,
	}
}

func (r *roleRepository) GetRole(ctx context.Context, filter RoleFilter) (Role, int, error) {
	var schema schemas.Role
	query := r.db.WithContext(ctx).Model(&schemas.Role{})
	query = r.applyFilter(query, filter)

	if err := query.First(&schema).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return Role{}, code, err
	}

	return r.toModel(schema), http.StatusOK, nil
}

func (r *roleRepository) GetRoles(ctx context.Context, filter RoleFilter, paginate *paginate_utils.PaginateData) ([]Role, int, error) {
	var schemaList []schemas.Role
	query := r.db.WithContext(ctx).Model(&schemas.Role{})
	query = r.applyFilter(query, filter)

	if err := common.CountTotal(query, paginate); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	query = common.ApplyPagination(query, paginate)
	query = common.ApplyOrderBy(query, "level DESC, created_at DESC")

	if err := query.Find(&schemaList).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return nil, code, err
	}

	roles := make([]Role, len(schemaList))
	for i, s := range schemaList {
		roles[i] = r.toModel(s)
	}

	return roles, http.StatusOK, nil
}

func (r *roleRepository) CreateRole(ctx context.Context, role Role) (Role, int, error) {
	schema := r.toSchema(role)
	if schema.Id == uuid.Nil {
		schema.Id = uuid.New()
	}
	schema.CreatedAt = time.Now()
	schema.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(&schema).Error; err != nil {
		return Role{}, http.StatusInternalServerError, err
	}

	return r.toModel(schema), http.StatusCreated, nil
}

func (r *roleRepository) UpdateRole(ctx context.Context, filter RoleFilter, role Role) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.Role{})
	query = r.applyFilter(query, filter)

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if role.Name != "" {
		updates["name"] = role.Name
	}
	if role.DisplayName != "" {
		updates["display_name"] = role.DisplayName
	}
	if role.Description != "" {
		updates["description"] = role.Description
	}
	if role.Level != 0 {
		updates["level"] = role.Level
	}

	result := query.Updates(updates)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	if result.RowsAffected == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return http.StatusOK, nil
}

func (r *roleRepository) DeleteRole(ctx context.Context, filter RoleFilter) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.Role{})
	query = r.applyFilter(query, filter)

	result := query.Delete(&schemas.Role{})
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	if result.RowsAffected == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return http.StatusOK, nil
}
