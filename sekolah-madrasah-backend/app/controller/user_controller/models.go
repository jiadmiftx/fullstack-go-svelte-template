package user_controller

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	Id              uuid.UUID  `json:"id"`
	Email           string     `json:"email"`
	FullName        string     `json:"full_name"`
	Phone           string     `json:"phone,omitempty"`
	Avatar          string     `json:"avatar,omitempty"`
	IsSuperAdmin    bool       `json:"is_super_admin"`
	IsActive        bool       `json:"is_active"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	LastLoginAt     *time.Time `json:"last_login_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
}

type UpdateUserRequest struct {
	Email    string `json:"email"`
	FullName string `json:"full_name"`
	Phone    string `json:"phone"`
	Avatar   string `json:"avatar"`
	IsActive *bool  `json:"is_active"`
}
