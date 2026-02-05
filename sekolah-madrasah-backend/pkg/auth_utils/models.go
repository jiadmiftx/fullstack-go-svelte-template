package auth_utils

import (
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"sekolah-madrasah/pkg/rbac"
)

type AuthClaim struct {
	UserID uuid.UUID `json:"user_id" default:"" order:"1"`
	Role   string    `json:"role" default:"user"`
	Exp    int64     `json:"exp" default:"1h" order:"10"`
	Level  rbac.Role

	jwt.StandardClaims
}

