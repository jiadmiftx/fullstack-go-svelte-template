package role_controller

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	Id             uuid.UUID    `json:"id"`
	OrganizationId *uuid.UUID   `json:"organization_id,omitempty"`
	Name           string       `json:"name"`
	DisplayName    string       `json:"display_name"`
	Type           string       `json:"type"`
	Level          int          `json:"level"`
	Description    string       `json:"description,omitempty"`
	IsDefault      bool         `json:"is_default"`
	Permissions    []Permission `json:"permissions,omitempty"`
	CreatedAt      time.Time    `json:"created_at"`
	UpdatedAt      time.Time    `json:"updated_at"`
}

type Permission struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Resource    string    `json:"resource"`
	Action      string    `json:"action"`
	Description string    `json:"description,omitempty"`
}

type CreateRoleRequest struct {
	OrganizationId *uuid.UUID  `json:"organization_id"`
	Name           string      `json:"name"`
	DisplayName    string      `json:"display_name"`
	Type           string      `json:"type"`
	Level          int         `json:"level"`
	Description    string      `json:"description"`
	IsDefault      bool        `json:"is_default"`
	PermissionIds  []uuid.UUID `json:"permission_ids"`
}

type UpdateRoleRequest struct {
	Name          string      `json:"name"`
	DisplayName   string      `json:"display_name"`
	Description   string      `json:"description"`
	Level         int         `json:"level"`
	PermissionIds []uuid.UUID `json:"permission_ids"`
}
