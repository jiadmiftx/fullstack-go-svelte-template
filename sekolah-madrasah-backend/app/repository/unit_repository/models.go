package unit_repository

import (
	"time"

	"github.com/google/uuid"
)

type Unit struct {
	Id             uuid.UUID
	OrganizationId uuid.UUID
	Name           string
	Code           string
	Type           string
	Address        string
	Phone          string
	Email          string
	Logo           string
	IsActive       bool
	Settings       string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}
