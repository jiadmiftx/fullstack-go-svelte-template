package user_repository

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

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) applyFilter(query *gorm.DB, filter UserFilter) *gorm.DB {
	if filter.Id != nil {
		query = query.Where("id = ?", *filter.Id)
	}
	if filter.Email != nil {
		query = query.Where("email = ?", *filter.Email)
	}
	if filter.IsSuperAdmin != nil {
		query = query.Where("is_super_admin = ?", *filter.IsSuperAdmin)
	}
	if filter.IsActive != nil {
		query = query.Where("is_active = ?", *filter.IsActive)
	}
	// Platform-only filter: exclude users who are in organization_members or unit_members
	if filter.PlatformOnly != nil && *filter.PlatformOnly {
		query = query.Where(`
			(is_super_admin = true) OR 
			(id NOT IN (SELECT user_id FROM organization_members WHERE deleted_at IS NULL) 
			 AND id NOT IN (SELECT user_id FROM unit_members WHERE deleted_at IS NULL))
		`)
	}
	return query
}

func (r *userRepository) toModel(schema schemas.User) User {
	return User{
		Id:              schema.Id,
		Email:           schema.Email,
		Password:        schema.Password,
		FullName:        schema.FullName,
		Phone:           schema.Phone,
		Avatar:          schema.Avatar,
		IsSuperAdmin:    schema.IsSuperAdmin,
		IsActive:        schema.IsActive,
		EmailVerifiedAt: schema.EmailVerifiedAt,
		LastLoginAt:     schema.LastLoginAt,
		CreatedAt:       schema.CreatedAt,
		UpdatedAt:       schema.UpdatedAt,
	}
}

func (r *userRepository) toSchema(model User) schemas.User {
	return schemas.User{
		Id:              model.Id,
		Email:           model.Email,
		Password:        model.Password,
		FullName:        model.FullName,
		Phone:           model.Phone,
		Avatar:          model.Avatar,
		IsSuperAdmin:    model.IsSuperAdmin,
		IsActive:        model.IsActive,
		EmailVerifiedAt: model.EmailVerifiedAt,
		LastLoginAt:     model.LastLoginAt,
	}
}

func (r *userRepository) GetUser(ctx context.Context, filter UserFilter) (User, int, error) {
	var schema schemas.User
	query := r.db.WithContext(ctx).Model(&schemas.User{})
	query = r.applyFilter(query, filter)

	if err := query.First(&schema).Error; err != nil {
		code, err := common.HandleGORMError(err)
		return User{}, code, err
	}

	return r.toModel(schema), http.StatusOK, nil
}

func (r *userRepository) GetUsers(ctx context.Context, filter UserFilter, paginate *paginate_utils.PaginateData) ([]User, int, error) {
	var schemaList []schemas.User
	query := r.db.WithContext(ctx).Model(&schemas.User{})
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

	users := make([]User, len(schemaList))
	for i, s := range schemaList {
		users[i] = r.toModel(s)
	}

	return users, http.StatusOK, nil
}

func (r *userRepository) CreateUser(ctx context.Context, user User) (User, int, error) {
	schema := r.toSchema(user)
	if schema.Id == uuid.Nil {
		schema.Id = uuid.New()
	}
	schema.CreatedAt = time.Now()
	schema.UpdatedAt = time.Now()

	if err := r.db.WithContext(ctx).Create(&schema).Error; err != nil {
		return User{}, http.StatusInternalServerError, err
	}

	return r.toModel(schema), http.StatusCreated, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, filter UserFilter, user User) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.User{})
	query = r.applyFilter(query, filter)

	updates := map[string]interface{}{
		"updated_at": time.Now(),
	}

	if user.Email != "" {
		updates["email"] = user.Email
	}
	if user.Password != "" {
		updates["password"] = user.Password
	}
	if user.FullName != "" {
		updates["full_name"] = user.FullName
	}
	if user.Phone != "" {
		updates["phone"] = user.Phone
	}
	if user.Avatar != "" {
		updates["avatar"] = user.Avatar
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

func (r *userRepository) DeleteUser(ctx context.Context, filter UserFilter) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.User{})
	query = r.applyFilter(query, filter)

	result := query.Delete(&schemas.User{})
	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}
	if result.RowsAffected == 0 {
		return http.StatusNotFound, gorm.ErrRecordNotFound
	}

	return http.StatusOK, nil
}

func (r *userRepository) UpdateLastLogin(ctx context.Context, filter UserFilter) (int, error) {
	query := r.db.WithContext(ctx).Model(&schemas.User{})
	query = r.applyFilter(query, filter)

	now := time.Now()
	result := query.Updates(map[string]interface{}{
		"last_login_at": now,
		"updated_at":    now,
	})

	if result.Error != nil {
		return http.StatusInternalServerError, result.Error
	}

	return http.StatusOK, nil
}

// HasPendingApproval checks if user has any warga profile with pending approval status
func (r *userRepository) HasPendingApproval(ctx context.Context, userId uuid.UUID) (bool, error) {
	status, err := r.GetApprovalStatus(ctx, userId)
	if err != nil {
		return false, err
	}
	return status == "pending", nil
}

// GetApprovalStatus returns the approval status of a user's warga profile
// Returns: "pending", "rejected", "approved", or "" if no warga profile
func (r *userRepository) GetApprovalStatus(ctx context.Context, userId uuid.UUID) (string, error) {
	var status string
	err := r.db.WithContext(ctx).
		Table("warga_profiles").
		Select("warga_profiles.approval_status").
		Joins("JOIN unit_members ON unit_members.id = warga_profiles.unit_member_id").
		Where("unit_members.user_id = ?", userId).
		Where("warga_profiles.deleted_at IS NULL").
		Limit(1).
		Scan(&status).Error
	if err != nil {
		return "", err
	}
	return status, nil
}
