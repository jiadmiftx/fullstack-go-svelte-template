package user_repository

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id              uuid.UUID
	Email           string
	Password        string
	FullName        string
	Phone           string
	Avatar          string
	IsSuperAdmin    bool
	IsActive        bool
	EmailVerifiedAt *time.Time
	LastLoginAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}
