package organization_repository

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

type organizationRepository struct {
	db *gorm.DB
}

func NewOrganizationRepository(db *gorm.DB) OrganizationRepository {
	return &organizationRepository{db: db}
}

func (r *organizationRepository) applyFilter(query *gorm.DB, filter OrganizationFilter) *gorm.DB {
	if filter.Id != nil {
		query = query.Where("id = ?", *filter.Id)
	}
	if filter.OwnerId != nil {
		query = query.Where("owner_id = ?", *filter.OwnerId)
	}
	if filter.Code != nil {
		query = query.Where("code = ?", *filter.Code)
	}
	if filter.Type != nil {
		query = query.Where("type = ?", *filter.Type)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	return query
}

func (r *organizationRepository) toModel(schema schemas.Organization) Organization {
	return Organization{
		Id:          schema.Id,
		OwnerId:     schema.OwnerId,
		Name:        schema.Name,
		Code:        schema.Code,
		Type:        schema.Type,
		Description: schema.Description,
		Address:     schema.Address,
		Logo:        schema.Logo,
		IsActive:    schema.IsActive,
		Settings:    schema.Settings,
		CreatedAt:   schema.CreatedAt,
		UpdatedAt:   schema.UpdatedAt,
	}
}

func (r *organizationRepository) GetOrganization(ctx context.Context, filter OrganizationFilter) (Organization, int, error) {
	var schema schemas.Organization
	query := r.db.WithContext(ctx).Model(&schemas.Organization{})
	query = r.applyFilter(query, filter)

	if err := query.First(&schema).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return Organization{}, code, err
	}

	return r.toModel(schema), http.StatusOK, nil
}

func (r *organizationRepository) GetOrganizations(ctx context.Context, filter OrganizationFilter, paginate *paginate_utils.PaginateData) ([]Organization, int, error) {
	var schemaList []schemas.Organization
	query := r.db.WithContext(ctx).Model(&schemas.Organization{})
	query = r.applyFilter(query, filter)

	if err := common.CountTotal(query, paginate); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	query = common.ApplyPagination(query, paginate)
	query = common.ApplyOrderBy(query, "created_at DESC")

	if err := query.Find(&schemaList).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return nil, code, err
	}

	orgs := make([]Organization, len(schemaList))
	for i, s := range schemaList {
		orgs[i] = r.toModel(s)
	}

	return orgs, http.StatusOK, nil
}

func (r *organizationRepository) CreateOrganization(ctx context.Context, org Organization) (Organization, int, error) {
	schema := schemas.Organization{
		Id:          org.Id,
		OwnerId:     org.OwnerId,
		Name:        org.Name,
		Code:        org.Code,
		Type:        org.Type,
		Description: org.Description,
		Address:     org.Address,
		Logo:        org.Logo,
		IsActive:    org.IsActive,
		Settings:    org.Settings,
	}

	if schema.Id == uuid.Nil {
		schema.Id = uuid.New()
	}
	schema.CreatedAt = time.Now()
	schema.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(&schema).Error; err != nil {
		return Organization{}, http.StatusInternalServerError, err
	}

	return r.toModel(schema), http.StatusCreated, nil
}

func (r *organizationRepository) UpdateOrganization(ctx context.Context, filter OrganizationFilter, org Organization) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.Organization{})
	query = r.applyFilter(query, filter)

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if org.Name != "" {
		updates["name"] = org.Name
	}
	if org.Description != "" {
		updates["description"] = org.Description
	}
	if org.Address != "" {
		updates["address"] = org.Address
	}
	if org.Logo != "" {
		updates["logo"] = org.Logo
	}
	if org.Settings != "" {
		updates["settings"] = org.Settings
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

func (r *organizationRepository) DeleteOrganization(ctx context.Context, filter OrganizationFilter) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.Organization{})
	query = r.applyFilter(query, filter)

	result := query.Delete(&schemas.Organization{})
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	if result.RowsAffected == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return http.StatusOK, nil
}
