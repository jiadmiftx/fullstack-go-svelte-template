package user_use_case

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id              uuid.UUID
	Email           string
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

type CreateUserRequest struct {
	Email    string
	Password string
	FullName string
	Phone    string
}

type UpdateUserRequest struct {
	Email    string
	FullName string
	Phone    string
	Avatar   string
	IsActive *bool
}

type UserFilter struct {
	Id           *uuid.UUID
	Email        *string
	IsSuperAdmin *bool
	IsActive     *bool
	PlatformOnly *bool
}
