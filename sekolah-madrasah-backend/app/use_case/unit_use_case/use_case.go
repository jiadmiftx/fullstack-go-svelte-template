package unit_use_case

import (
	"context"
	"errors"
	"net/http"

	"sekolah-madrasah/app/repository/unit_repository"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
)

type unitUseCase struct {
	unitRepo unit_repository.UnitRepository
}

func NewUnitUseCase(unitRepo unit_repository.UnitRepository) UnitUseCase {
	return &unitUseCase{unitRepo: unitRepo}
}

func (u *unitUseCase) toUnit(s unit_repository.Unit) Unit {
	return Unit{
		Id:             s.Id,
		OrganizationId: s.OrganizationId,
		Name:           s.Name,
		Code:           s.Code,
		Type:           s.Type,
		Address:        s.Address,
		Phone:          s.Phone,
		Email:          s.Email,
		Logo:           s.Logo,
		IsActive:       s.IsActive,
		Settings:       s.Settings,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

func (u *unitUseCase) GetUnit(ctx context.Context, id uuid.UUID) (Unit, int, error) {
	perumahan, code, err := u.unitRepo.GetUnit(ctx, unit_repository.UnitFilter{Id: &id})
	if err != nil {
		return Unit{}, code, err
	}
	return u.toUnit(perumahan), http.StatusOK, nil
}

func (u *unitUseCase) GetUnits(ctx context.Context, filter UnitFilter, paginate *paginate_utils.PaginateData) ([]Unit, int, error) {
	repoFilter := unit_repository.UnitFilter{
		OrganizationId: filter.OrganizationId,
		Type:           filter.Type,
		IsActive:       filter.IsActive,
	}

	perumahans, code, err := u.unitRepo.GetUnits(ctx, repoFilter, paginate)
	if err != nil {
		return nil, code, err
	}

	result := make([]Unit, len(perumahans))
	for i, s := range perumahans {
		result[i] = u.toUnit(s)
	}

	return result, http.StatusOK, nil
}

func (u *unitUseCase) CreateUnit(ctx context.Context, organizationId uuid.UUID, req CreateUnitRequest) (Unit, int, error) {
	// Check if code already exists
	existingUnit, code, _ := u.unitRepo.GetUnit(ctx, unit_repository.UnitFilter{Code: &req.Code})
	if code == http.StatusOK && existingUnit.Id != uuid.Nil {
		return Unit{}, http.StatusConflict, errors.New("perumahan code already exists")
	}

	perumahanType := req.Type
	if perumahanType == "" {
		perumahanType = "SMP"
	}

	newUnit := unit_repository.Unit{
		OrganizationId: organizationId,
		Name:           req.Name,
		Code:           req.Code,
		Type:           perumahanType,
		Address:        req.Address,
		Phone:          req.Phone,
		Email:          req.Email,
		Logo:           req.Logo,
		IsActive:       true,
	}

	createdUnit, code, err := u.unitRepo.CreateUnit(ctx, newUnit)
	if err != nil {
		return Unit{}, code, err
	}

	return u.toUnit(createdUnit), http.StatusCreated, nil
}

func (u *unitUseCase) UpdateUnit(ctx context.Context, id uuid.UUID, req UpdateUnitRequest) (Unit, int, error) {
	_, code, err := u.unitRepo.GetUnit(ctx, unit_repository.UnitFilter{Id: &id})
	if err != nil {
		return Unit{}, code, err
	}

	updateData := unit_repository.Unit{
		Name:     req.Name,
		Address:  req.Address,
		Phone:    req.Phone,
		Email:    req.Email,
		Logo:     req.Logo,
		Settings: req.Settings,
	}

	code, err = u.unitRepo.UpdateUnit(ctx, unit_repository.UnitFilter{Id: &id}, updateData)
	if err != nil {
		return Unit{}, code, err
	}

	updatedUnit, code, err := u.unitRepo.GetUnit(ctx, unit_repository.UnitFilter{Id: &id})
	if err != nil {
		return Unit{}, code, err
	}

	return u.toUnit(updatedUnit), http.StatusOK, nil
}

func (u *unitUseCase) DeleteUnit(ctx context.Context, id uuid.UUID) (int, error) {
	_, code, err := u.unitRepo.GetUnit(ctx, unit_repository.UnitFilter{Id: &id})
	if err != nil {
		return code, err
	}

	return u.unitRepo.DeleteUnit(ctx, unit_repository.UnitFilter{Id: &id})
}
