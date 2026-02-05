package membership_service

import (
	"context"
	"net/http"

	"sekolah-madrasah/database/schemas"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MembershipService interface {
	GetUserMemberships(ctx context.Context, userId uuid.UUID) (UserMemberships, int, error)
}

type membershipService struct {
	db *gorm.DB
}

func NewMembershipService(db *gorm.DB) MembershipService {
	return &membershipService{db: db}
}

// Response DTOs
type OrganizationMembership struct {
	OrgId    uuid.UUID `json:"org_id"`
	OrgName  string    `json:"org_name"`
	RoleId   uuid.UUID `json:"role_id"`
	RoleName string    `json:"role_name"`
}

type UnitMembership struct {
	UnitMemberId uuid.UUID `json:"unit_member_id"`
	UnitId       uuid.UUID `json:"unit_id"`
	UnitName     string    `json:"perumahan_name"`
	OrgId        uuid.UUID `json:"org_id"`
	OrgName      string    `json:"org_name"`
	Role         string    `json:"role"` // pengurus, warga, admin, staff, parent
	IsActive     bool      `json:"is_active"`
}

type UserMemberships struct {
	UserId                  uuid.UUID                `json:"user_id"`
	IsSuperAdmin            bool                     `json:"is_super_admin"`
	OrganizationMemberships []OrganizationMembership `json:"organization_memberships"`
	UnitMemberships         []UnitMembership         `json:"unit_memberships"`
}

func (s *membershipService) GetUserMemberships(ctx context.Context, userId uuid.UUID) (UserMemberships, int, error) {
	result := UserMemberships{
		UserId:                  userId,
		OrganizationMemberships: []OrganizationMembership{},
		UnitMemberships:         []UnitMembership{},
	}

	// Get user to check is_super_admin
	var user schemas.User
	if err := s.db.WithContext(ctx).Where("id = ?", userId).First(&user).Error; err != nil {
		return result, http.StatusNotFound, err
	}
	result.IsSuperAdmin = user.IsSuperAdmin

	// Get organization memberships
	var orgMembers []schemas.OrganizationMember
	if err := s.db.WithContext(ctx).
		Preload("Organization").
		Preload("Role").
		Where("user_id = ?", userId).
		Find(&orgMembers).Error; err != nil {
		return result, http.StatusInternalServerError, err
	}

	for _, om := range orgMembers {
		membership := OrganizationMembership{
			OrgId:  om.OrganizationId,
			RoleId: om.RoleId,
		}
		if om.Organization != nil {
			membership.OrgName = om.Organization.Name
		}
		if om.Role != nil {
			membership.RoleName = om.Role.Name
		}
		result.OrganizationMemberships = append(result.OrganizationMemberships, membership)
	}

	// Get perumahan memberships
	var perumahanMembers []schemas.UnitMember
	if err := s.db.WithContext(ctx).
		Preload("Unit").
		Preload("Unit.Organization").
		Where("user_id = ?", userId).
		Where("is_active = ?", true).
		Find(&perumahanMembers).Error; err != nil {
		return result, http.StatusInternalServerError, err
	}

	for _, sm := range perumahanMembers {
		membership := UnitMembership{
			UnitMemberId: sm.Id,
			UnitId:       sm.UnitId,
			Role:         string(sm.Role),
			IsActive:     sm.IsActive,
		}
		if sm.Unit != nil {
			membership.UnitName = sm.Unit.Name
			membership.OrgId = sm.Unit.OrganizationId
			if sm.Unit.Organization != nil {
				membership.OrgName = sm.Unit.Organization.Name
			}
		}
		result.UnitMemberships = append(result.UnitMemberships, membership)
	}

	return result, http.StatusOK, nil
}
