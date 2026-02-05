package unit_use_case

import (
	"time"

	"github.com/google/uuid"
)

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
	Settings       string    `json:"settings,omitempty"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type CreateUnitRequest struct {
	Name    string `json:"name"`
	Code    string `json:"code"`
	Type    string `json:"type"`
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

type UnitFilter struct {
	OrganizationId *uuid.UUID
	Type           *string
	IsActive       *bool
}
