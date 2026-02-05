package unit_repository

import (
	"context"
	"net/http"
	"time"

	"sekolah-madrasah/app/repository/common"
	"sekolah-madrasah/database/schemas"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type unitRepository struct {
	db *gorm.DB
}

func NewUnitRepository(db *gorm.DB) UnitRepository {
	return &unitRepository{db: db}
}

func (r *unitRepository) applyFilter(query *gorm.DB, filter UnitFilter) *gorm.DB {
	if filter.Id != nil {
		query = query.Where("id = ?", *filter.Id)
	}
	if filter.OrganizationId != nil {
		query = query.Where("organization_id = ?", *filter.OrganizationId)
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

func (r *unitRepository) toModel(schema schemas.Unit) Unit {
	return Unit{
		Id:             schema.Id,
		OrganizationId: schema.OrganizationId,
		Name:           schema.Name,
		Code:           schema.Code,
		Type:           schema.Type,
		Address:        schema.Address,
		Phone:          schema.Phone,
		Email:          schema.Email,
		Logo:           schema.Logo,
		IsActive:       schema.IsActive,
		Settings:       schema.Settings,
		CreatedAt:      schema.CreatedAt,
		UpdatedAt:      schema.UpdatedAt,
	}
}

func (r *unitRepository) GetUnit(ctx context.Context, filter UnitFilter) (Unit, int, error) {
	var schema schemas.Unit
	query := r.db.WithContext(ctx).Model(&schemas.Unit{})
	query = r.applyFilter(query, filter)

	if err := query.First(&schema).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return Unit{}, code, err
	}

	return r.toModel(schema), http.StatusOK, nil
}

func (r *unitRepository) GetUnits(ctx context.Context, filter UnitFilter, paginate *paginate_utils.PaginateData) ([]Unit, int, error) {
	var schemaList []schemas.Unit
	query := r.db.WithContext(ctx).Model(&schemas.Unit{})
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

	perumahans := make([]Unit, len(schemaList))
	for i, s := range schemaList {
		perumahans[i] = r.toModel(s)
	}

	return perumahans, http.StatusOK, nil
}

func (r *unitRepository) CreateUnit(ctx context.Context, perumahan Unit) (Unit, int, error) {
	schema := schemas.Unit{
		Id:             perumahan.Id,
		OrganizationId: perumahan.OrganizationId,
		Name:           perumahan.Name,
		Code:           perumahan.Code,
		Type:           perumahan.Type,
		Address:        perumahan.Address,
		Phone:          perumahan.Phone,
		Email:          perumahan.Email,
		Logo:           perumahan.Logo,
		IsActive:       perumahan.IsActive,
		Settings:       perumahan.Settings,
	}

	if schema.Id == uuid.Nil {
		schema.Id = uuid.New()
	}
	schema.CreatedAt = time.Now()
	schema.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(&schema).Error; err != nil {
		return Unit{}, http.StatusInternalServerError, err
	}

	return r.toModel(schema), http.StatusCreated, nil
}

func (r *unitRepository) UpdateUnit(ctx context.Context, filter UnitFilter, perumahan Unit) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.Unit{})
	query = r.applyFilter(query, filter)

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if perumahan.Name != "" {
		updates["name"] = perumahan.Name
	}
	if perumahan.Address != "" {
		updates["address"] = perumahan.Address
	}
	if perumahan.Phone != "" {
		updates["phone"] = perumahan.Phone
	}
	if perumahan.Email != "" {
		updates["email"] = perumahan.Email
	}
	if perumahan.Logo != "" {
		updates["logo"] = perumahan.Logo
	}
	if perumahan.Settings != "" {
		updates["settings"] = perumahan.Settings
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

func (r *unitRepository) DeleteUnit(ctx context.Context, filter UnitFilter) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.Unit{})
	query = r.applyFilter(query, filter)

	result := query.Delete(&schemas.Unit{})
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	if result.RowsAffected == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return http.StatusOK, nil
}
