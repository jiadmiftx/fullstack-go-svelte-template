package unit_controller

import (
	"sekolah-madrasah/app/use_case/unit_use_case"

	"github.com/google/uuid"
)

type UnitController interface {
	GetUnit(c interface{})
	GetUnits(c interface{})
	CreateUnit(c interface{})
	UpdateUnit(c interface{})
	DeleteUnit(c interface{})
}

// DTOs for HTTP layer
type Unit struct {
	Id             uuid.UUID `json:"id"`
	OrganizationId uuid.UUID `json:"organization_id"`
	Name           string    `json:"name"`
	Code           string    `json:"code"`
	Type           string    `json:"type"`
	Address        string    `json:"address"`
	Phone          string    `json:"phone"`
	Email          string    `json:"email"`
	Logo           string    `json:"logo"`
	IsActive       bool      `json:"is_active"`
}

type CreateUnitRequest struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Type    string `json:"type,omitempty"`
	Address string `json:"address,omitempty"`
	Phone   string `json:"phone,omitempty"`
	Email   string `json:"email,omitempty"`
	Logo    string `json:"logo,omitempty"`
}

type UpdateUnitRequest struct {
	Name     string `json:"name,omitempty"`
	Address  string `json:"address,omitempty"`
	Phone    string `json:"phone,omitempty"`
	Email    string `json:"email,omitempty"`
	Logo     string `json:"logo,omitempty"`
	Settings string `json:"settings,omitempty"`
}

func toUnitResponse(s unit_use_case.Unit) Unit {
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
	}
}
