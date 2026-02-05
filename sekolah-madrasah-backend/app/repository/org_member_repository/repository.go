package org_member_repository

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

type orgMemberRepository struct {
	db *gorm.DB
}

func NewOrgMemberRepository(db *gorm.DB) OrgMemberRepository {
	return &orgMemberRepository{db: db}
}

func (r *orgMemberRepository) applyFilter(query *gorm.DB, filter OrgMemberFilter) *gorm.DB {
	if filter.Id != nil {
		query = query.Where("id = ?", *filter.Id)
	}
	if filter.UserId != nil {
		query = query.Where("user_id = ?", *filter.UserId)
	}
	if filter.OrganizationId != nil {
		query = query.Where("organization_id = ?", *filter.OrganizationId)
	}
	if filter.RoleId != nil {
		query = query.Where("role_id = ?", *filter.RoleId)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	return query
}

func (r *orgMemberRepository) toModel(schema schemas.OrganizationMember) OrganizationMember {
	return OrganizationMember{
		Id:             schema.Id,
		UserId:         schema.UserId,
		OrganizationId: schema.OrganizationId,
		RoleId:         schema.RoleId,
		IsActive:       schema.IsActive,
		JoinedAt:       schema.JoinedAt,
		InvitedBy:      schema.InvitedBy,
		CreatedAt:      schema.CreatedAt,
		UpdatedAt:      schema.UpdatedAt,
	}
}

func (r *orgMemberRepository) GetMember(ctx context.Context, filter OrgMemberFilter) (OrganizationMember, int, error) {
	var schema schemas.OrganizationMember
	query := r.db.WithContext(ctx).Model(&schemas.OrganizationMember{})
	query = r.applyFilter(query, filter)

	if err := query.First(&schema).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return OrganizationMember{}, code, err
	}

	return r.toModel(schema), http.StatusOK, nil
}

func (r *orgMemberRepository) GetMembers(ctx context.Context, filter OrgMemberFilter, paginate *paginate_utils.PaginateData) ([]OrganizationMember, int, error) {
	var schemaList []schemas.OrganizationMember
	query := r.db.WithContext(ctx).Model(&schemas.OrganizationMember{})
	query = r.applyFilter(query, filter)

	if err := common.CountTotal(query, paginate); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	query = common.ApplyPagination(query, paginate)
	query = common.ApplyOrderBy(query, "joined_at DESC")

	if err := query.Find(&schemaList).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return nil, code, err
	}

	members := make([]OrganizationMember, len(schemaList))
	for i, s := range schemaList {
		members[i] = r.toModel(s)
	}

	return members, http.StatusOK, nil
}

func (r *orgMemberRepository) AddMember(ctx context.Context, member OrganizationMember) (OrganizationMember, int, error) {
	schema := schemas.OrganizationMember{
		Id:             member.Id,
		UserId:         member.UserId,
		OrganizationId: member.OrganizationId,
		RoleId:         member.RoleId,
		IsActive:       member.IsActive,
		InvitedBy:      member.InvitedBy,
	}

	if schema.Id == uuid.Nil {
		schema.Id = uuid.New()
	}
	schema.JoinedAt = time.Now()
	schema.CreatedAt = time.Now()
	schema.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(&schema).Error; err != nil {
		return OrganizationMember{}, http.StatusInternalServerError, err
	}

	return r.toModel(schema), http.StatusCreated, nil
}

func (r *orgMemberRepository) UpdateMember(ctx context.Context, filter OrgMemberFilter, member OrganizationMember) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.OrganizationMember{})
	query = r.applyFilter(query, filter)

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if member.RoleId != uuid.Nil {
		updates["role_id"] = member.RoleId
	}

	updates["is_active"] = member.IsActive

	result := query.Updates(updates)
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	if result.RowsAffected == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return http.StatusOK, nil
}

func (r *orgMemberRepository) RemoveMember(ctx context.Context, filter OrgMemberFilter) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.OrganizationMember{})
	query = r.applyFilter(query, filter)

	result := query.Delete(&schemas.OrganizationMember{})
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	if result.RowsAffected == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return http.StatusOK, nil
}
