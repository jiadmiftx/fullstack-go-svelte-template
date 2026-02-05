package auth_use_case

import (
	"time"

	"github.com/google/uuid"
)

type LoginRequest struct {
	Email    string
	Password string
}

type LoginResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresAt    int64
	User         UserInfo
}

type RegisterRequest struct {
	Email    string
	Password string
	FullName string
}

type RefreshTokenRequest struct {
	RefreshToken string
}

type UserInfo struct {
	Id           uuid.UUID
	Email        string
	FullName     string
	IsSuperAdmin bool
	IsActive     bool
	LastLoginAt  *time.Time
}
