package unit_member_use_case

import (
	"context"
	"errors"
	"net/http"

	"sekolah-madrasah/app/repository/unit_member_repository"
	"sekolah-madrasah/database/schemas"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
)

type perumahanMemberUseCase struct {
	memberRepo unit_member_repository.UnitMemberRepository
}

func NewUnitMemberUseCase(memberRepo unit_member_repository.UnitMemberRepository) UnitMemberUseCase {
	return &perumahanMemberUseCase{memberRepo: memberRepo}
}

func (u *perumahanMemberUseCase) toMember(m unit_member_repository.UnitMember) UnitMember {
	member := UnitMember{
		Id:        m.Id,
		UserId:    m.UserId,
		UnitId:  m.UnitId,
		Role:      m.Role,
		IsActive:  m.IsActive,
		JoinedAt:  m.JoinedAt,
		InvitedBy: m.InvitedBy,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}

	if m.User != nil {
		member.User = &UserInfo{
			Id:       m.User.Id,
			FullName: m.User.FullName,
			Email:    m.User.Email,
		}
	}

	if m.Unit != nil {
		member.Unit = &UnitInfo{
			Id:   m.Unit.Id,
			Name: m.Unit.Name,
			Code: m.Unit.Code,
			Type: m.Unit.Type,
		}
	}

	return member
}

func (u *perumahanMemberUseCase) GetMember(ctx context.Context, perumahanId, memberId uuid.UUID) (UnitMember, int, error) {
	member, code, err := u.memberRepo.GetMember(ctx, unit_member_repository.UnitMemberFilter{
		Id:       &memberId,
		UnitId: &perumahanId,
	})
	if err != nil {
		return UnitMember{}, code, err
	}
	return u.toMember(member), http.StatusOK, nil
}

func (u *perumahanMemberUseCase) GetMembers(ctx context.Context, perumahanId uuid.UUID, filter UnitMemberFilter, paginate *paginate_utils.PaginateData) ([]UnitMember, int, error) {
	repoFilter := unit_member_repository.UnitMemberFilter{
		UnitId: &perumahanId,
		UserId:   filter.UserId,
		Role:     filter.Role,
		IsActive: filter.IsActive,
	}

	members, code, err := u.memberRepo.GetMembers(ctx, repoFilter, paginate)
	if err != nil {
		return nil, code, err
	}

	result := make([]UnitMember, len(members))
	for i, m := range members {
		result[i] = u.toMember(m)
	}

	return result, http.StatusOK, nil
}

func (u *perumahanMemberUseCase) AddMember(ctx context.Context, perumahanId uuid.UUID, req AddMemberRequest, invitedBy *uuid.UUID) (UnitMember, int, error) {
	// Check if user is already a member of this perumahan
	existingMember, code, _ := u.memberRepo.GetMember(ctx, unit_member_repository.UnitMemberFilter{
		UserId:   &req.UserId,
		UnitId: &perumahanId,
	})
	if code == http.StatusOK && existingMember.Id != uuid.Nil {
		return UnitMember{}, http.StatusConflict, errors.New("user is already a member of this perumahan")
	}

	// Validate role
	if req.Role == "" {
		req.Role = schemas.UnitMemberRoleStaff
	}

	newMember := unit_member_repository.UnitMember{
		UserId:    req.UserId,
		UnitId:  perumahanId,
		Role:      req.Role,
		IsActive:  true,
		InvitedBy: invitedBy,
	}

	createdMember, code, err := u.memberRepo.AddMember(ctx, newMember)
	if err != nil {
		return UnitMember{}, code, err
	}

	// Fetch with relations
	member, code, err := u.memberRepo.GetMember(ctx, unit_member_repository.UnitMemberFilter{Id: &createdMember.Id})
	if err != nil {
		return u.toMember(createdMember), http.StatusCreated, nil
	}

	return u.toMember(member), http.StatusCreated, nil
}

func (u *perumahanMemberUseCase) UpdateMember(ctx context.Context, perumahanId, memberId uuid.UUID, req UpdateMemberRequest) (UnitMember, int, error) {
	// Check if member exists
	_, code, err := u.memberRepo.GetMember(ctx, unit_member_repository.UnitMemberFilter{
		Id:       &memberId,
		UnitId: &perumahanId,
	})
	if err != nil {
		return UnitMember{}, code, err
	}

	updateData := unit_member_repository.UnitMember{
		Role: req.Role,
	}
	if req.IsActive != nil {
		updateData.IsActive = *req.IsActive
	}

	code, err = u.memberRepo.UpdateMember(ctx, unit_member_repository.UnitMemberFilter{Id: &memberId}, updateData)
	if err != nil {
		return UnitMember{}, code, err
	}

	// Fetch updated member
	updatedMember, code, err := u.memberRepo.GetMember(ctx, unit_member_repository.UnitMemberFilter{Id: &memberId})
	if err != nil {
		return UnitMember{}, code, err
	}

	return u.toMember(updatedMember), http.StatusOK, nil
}

func (u *perumahanMemberUseCase) RemoveMember(ctx context.Context, perumahanId, memberId uuid.UUID) (int, error) {
	// Check if member exists
	_, code, err := u.memberRepo.GetMember(ctx, unit_member_repository.UnitMemberFilter{
		Id:       &memberId,
		UnitId: &perumahanId,
	})
	if err != nil {
		return code, err
	}

	return u.memberRepo.RemoveMember(ctx, unit_member_repository.UnitMemberFilter{Id: &memberId})
}
