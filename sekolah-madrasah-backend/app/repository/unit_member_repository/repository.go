package unit_member_repository

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

type unitMemberRepository struct {
	db *gorm.DB
}

func NewUnitMemberRepository(db *gorm.DB) UnitMemberRepository {
	return &unitMemberRepository{db: db}
}

func (r *unitMemberRepository) applyFilter(query *gorm.DB, filter UnitMemberFilter) *gorm.DB {
	if filter.Id != nil {
		query = query.Where("unit_members.id = ?", *filter.Id)
	}
	if filter.UserId != nil {
		query = query.Where("unit_members.user_id = ?", *filter.UserId)
	}
	if filter.UnitId != nil {
		query = query.Where("unit_members.unit_id = ?", *filter.UnitId)
	}
	if filter.Role != nil {
		query = query.Where("unit_members.role = ?", *filter.Role)
	}
	if filter.IsActive != nil {
		query = query.Where("unit_members.is_active = ?", *filter.IsActive)
	}
	return query
}

func (r *unitMemberRepository) toModel(schema schemas.UnitMember) UnitMember {
	member := UnitMember{
		Id:        schema.Id,
		UserId:    schema.UserId,
		UnitId:    schema.UnitId,
		Role:      schema.Role,
		IsActive:  schema.IsActive,
		JoinedAt:  schema.JoinedAt,
		InvitedBy: schema.InvitedBy,
		CreatedAt: schema.CreatedAt,
		UpdatedAt: schema.UpdatedAt,
	}

	if schema.User != nil {
		member.User = &UserInfo{
			Id:       schema.User.Id,
			FullName: schema.User.FullName,
			Email:    schema.User.Email,
		}
	}

	if schema.Unit != nil {
		member.Unit = &UnitInfo{
			Id:   schema.Unit.Id,
			Name: schema.Unit.Name,
			Code: schema.Unit.Code,
			Type: schema.Unit.Type,
		}
	}

	return member
}

func (r *unitMemberRepository) GetMember(ctx context.Context, filter UnitMemberFilter) (UnitMember, int, error) {
	var schema schemas.UnitMember
	query := r.db.WithContext(ctx).Model(&schemas.UnitMember{})
	query = query.Preload("User").Preload("Unit")
	query = r.applyFilter(query, filter)

	if err := query.First(&schema).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return UnitMember{}, code, err
	}

	return r.toModel(schema), http.StatusOK, nil
}

func (r *unitMemberRepository) GetMembers(ctx context.Context, filter UnitMemberFilter, paginate *paginate_utils.PaginateData) ([]UnitMember, int, error) {
	var schemaList []schemas.UnitMember
	query := r.db.WithContext(ctx).Model(&schemas.UnitMember{})
	query = r.applyFilter(query, filter)

	if err := common.CountTotal(query, paginate); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	query = query.Preload("User").Preload("Unit")
	query = common.ApplyPagination(query, paginate)
	query = common.ApplyOrderBy(query, "joined_at DESC")

	if err := query.Find(&schemaList).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return nil, code, err
	}

	members := make([]UnitMember, len(schemaList))
	for i, s := range schemaList {
		members[i] = r.toModel(s)
	}

	return members, http.StatusOK, nil
}

func (r *unitMemberRepository) AddMember(ctx context.Context, member UnitMember) (UnitMember, int, error) {
	schema := schemas.UnitMember{
		Id:        member.Id,
		UserId:    member.UserId,
		UnitId:    member.UnitId,
		Role:      member.Role,
		IsActive:  member.IsActive,
		InvitedBy: member.InvitedBy,
	}

	if schema.Id == uuid.Nil {
		schema.Id = uuid.New()
	}
	schema.JoinedAt = time.Now()
	schema.CreatedAt = time.Now()
	schema.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(&schema).Error; err != nil {
		return UnitMember{}, http.StatusInternalServerError, err
	}

	return r.toModel(schema), http.StatusCreated, nil
}

func (r *unitMemberRepository) UpdateMember(ctx context.Context, filter UnitMemberFilter, member UnitMember) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.UnitMember{})
	query = r.applyFilter(query, filter)

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if member.Role != "" {
		updates["role"] = member.Role
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

func (r *unitMemberRepository) RemoveMember(ctx context.Context, filter UnitMemberFilter) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.UnitMember{})
	query = r.applyFilter(query, filter)

	result := query.Delete(&schemas.UnitMember{})
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	if result.RowsAffected == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return http.StatusOK, nil
}
