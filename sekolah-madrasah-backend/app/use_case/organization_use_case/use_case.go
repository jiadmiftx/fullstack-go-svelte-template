package organization_use_case

import (
	"context"
	"errors"
	"net/http"

	"github.com/google/uuid"
	"sekolah-madrasah/app/repository/org_member_repository"
	"sekolah-madrasah/app/repository/organization_repository"
	"sekolah-madrasah/pkg/paginate_utils"
)

type organizationUseCase struct {
	orgRepo       organization_repository.OrganizationRepository
	memberRepo    org_member_repository.OrgMemberRepository
}

func NewOrganizationUseCase(
	orgRepo organization_repository.OrganizationRepository,
	memberRepo org_member_repository.OrgMemberRepository,
) OrganizationUseCase {
	return &organizationUseCase{
		orgRepo:    orgRepo,
		memberRepo: memberRepo,
	}
}

func (u *organizationUseCase) toOrganization(o organization_repository.Organization) Organization {
	return Organization{
		Id:          o.Id,
		OwnerId:     o.OwnerId,
		Name:        o.Name,
		Code:        o.Code,
		Type:        o.Type,
		Description: o.Description,
		Address:     o.Address,
		Logo:        o.Logo,
		IsActive:    o.IsActive,
		Settings:    o.Settings,
		CreatedAt:   o.CreatedAt,
		UpdatedAt:   o.UpdatedAt,
	}
}

func (u *organizationUseCase) toMember(m org_member_repository.OrganizationMember) OrganizationMember {
	return OrganizationMember{
		Id:             m.Id,
		UserId:         m.UserId,
		OrganizationId: m.OrganizationId,
		RoleId:         m.RoleId,
		IsActive:       m.IsActive,
		JoinedAt:       m.JoinedAt,
		InvitedBy:      m.InvitedBy,
	}
}

func (u *organizationUseCase) GetOrganization(ctx context.Context, id uuid.UUID) (Organization, int, error) {
	org, code, err := u.orgRepo.GetOrganization(ctx, organization_repository.OrganizationFilter{Id: &id})
	if err != nil {
		return Organization{}, code, err
	}

	result := u.toOrganization(org)

	memberPaginate := &paginate_utils.PaginateData{Page: 1, Limit: 1}
	_, _, _ = u.memberRepo.GetMembers(ctx, org_member_repository.OrgMemberFilter{OrganizationId: &id}, memberPaginate)
	result.MemberCount = int(memberPaginate.TotalData)

	return result, http.StatusOK, nil
}

func (u *organizationUseCase) GetOrganizations(ctx context.Context, filter OrganizationFilter, paginate *paginate_utils.PaginateData) ([]Organization, int, error) {
	repoFilter := organization_repository.OrganizationFilter{
		OwnerId:  filter.OwnerId,
		Type:     filter.Type,
		IsActive: filter.IsActive,
	}

	orgs, code, err := u.orgRepo.GetOrganizations(ctx, repoFilter, paginate)
	if err != nil {
		return nil, code, err
	}

	result := make([]Organization, len(orgs))
	for i, o := range orgs {
		result[i] = u.toOrganization(o)
	}

	return result, http.StatusOK, nil
}

func (u *organizationUseCase) CreateOrganization(ctx context.Context, ownerId uuid.UUID, req CreateOrganizationRequest) (Organization, int, error) {
	existingOrg, code, _ := u.orgRepo.GetOrganization(ctx, organization_repository.OrganizationFilter{
		Code: &req.Code,
	})
	if code == http.StatusOK && existingOrg.Id != uuid.Nil {
		return Organization{}, http.StatusConflict, errors.New("organization code already exists")
	}

	orgType := req.Type
	if orgType == "" {
		orgType = "general"
	}

	newOrg := organization_repository.Organization{
		OwnerId:     ownerId,
		Name:        req.Name,
		Code:        req.Code,
		Type:        orgType,
		Description: req.Description,
		Address:     req.Address,
		Logo:        req.Logo,
		IsActive:    true,
	}

	createdOrg, code, err := u.orgRepo.CreateOrganization(ctx, newOrg)
	if err != nil {
		return Organization{}, code, err
	}

	return u.toOrganization(createdOrg), http.StatusCreated, nil
}

func (u *organizationUseCase) UpdateOrganization(ctx context.Context, id uuid.UUID, req UpdateOrganizationRequest) (Organization, int, error) {
	_, code, err := u.orgRepo.GetOrganization(ctx, organization_repository.OrganizationFilter{Id: &id})
	if err != nil {
		return Organization{}, code, err
	}

	updateData := organization_repository.Organization{
		Name:        req.Name,
		Description: req.Description,
		Address:     req.Address,
		Logo:        req.Logo,
		Settings:    req.Settings,
	}

	code, err = u.orgRepo.UpdateOrganization(ctx, organization_repository.OrganizationFilter{Id: &id}, updateData)
	if err != nil {
		return Organization{}, code, err
	}

	updatedOrg, code, err := u.orgRepo.GetOrganization(ctx, organization_repository.OrganizationFilter{Id: &id})
	if err != nil {
		return Organization{}, code, err
	}

	return u.toOrganization(updatedOrg), http.StatusOK, nil
}

func (u *organizationUseCase) DeleteOrganization(ctx context.Context, id uuid.UUID) (int, error) {
	_, code, err := u.orgRepo.GetOrganization(ctx, organization_repository.OrganizationFilter{Id: &id})
	if err != nil {
		return code, err
	}

	return u.orgRepo.DeleteOrganization(ctx, organization_repository.OrganizationFilter{Id: &id})
}

func (u *organizationUseCase) GetMembers(ctx context.Context, orgId uuid.UUID, paginate *paginate_utils.PaginateData) ([]OrganizationMember, int, error) {
	members, code, err := u.memberRepo.GetMembers(ctx, org_member_repository.OrgMemberFilter{
		OrganizationId: &orgId,
	}, paginate)
	if err != nil {
		return nil, code, err
	}

	result := make([]OrganizationMember, len(members))
	for i, m := range members {
		result[i] = u.toMember(m)
	}

	return result, http.StatusOK, nil
}

func (u *organizationUseCase) AddMember(ctx context.Context, orgId uuid.UUID, req AddMemberRequest) (OrganizationMember, int, error) {
	existingMember, code, _ := u.memberRepo.GetMember(ctx, org_member_repository.OrgMemberFilter{
		UserId:         &req.UserId,
		OrganizationId: &orgId,
	})
	if code == http.StatusOK && existingMember.Id != uuid.Nil {
		return OrganizationMember{}, http.StatusConflict, errors.New("user is already a member of this organization")
	}

	newMember := org_member_repository.OrganizationMember{
		UserId:         req.UserId,
		OrganizationId: orgId,
		RoleId:         req.RoleId,
		IsActive:       true,
		InvitedBy:      req.InvitedBy,
	}

	createdMember, code, err := u.memberRepo.AddMember(ctx, newMember)
	if err != nil {
		return OrganizationMember{}, code, err
	}

	return u.toMember(createdMember), http.StatusCreated, nil
}

func (u *organizationUseCase) UpdateMember(ctx context.Context, orgId uuid.UUID, userId uuid.UUID, req UpdateMemberRequest) (int, error) {
	_, code, err := u.memberRepo.GetMember(ctx, org_member_repository.OrgMemberFilter{
		UserId:         &userId,
		OrganizationId: &orgId,
	})
	if err != nil {
		return code, err
	}

	updateData := org_member_repository.OrganizationMember{
		RoleId:   req.RoleId,
		IsActive: req.IsActive,
	}

	return u.memberRepo.UpdateMember(ctx, org_member_repository.OrgMemberFilter{
		UserId:         &userId,
		OrganizationId: &orgId,
	}, updateData)
}

func (u *organizationUseCase) RemoveMember(ctx context.Context, orgId uuid.UUID, userId uuid.UUID) (int, error) {
	_, code, err := u.memberRepo.GetMember(ctx, org_member_repository.OrgMemberFilter{
		UserId:         &userId,
		OrganizationId: &orgId,
	})
	if err != nil {
		return code, err
	}

	return u.memberRepo.RemoveMember(ctx, org_member_repository.OrgMemberFilter{
		UserId:         &userId,
		OrganizationId: &orgId,
	})
}
