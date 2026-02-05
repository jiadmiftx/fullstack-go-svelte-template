package auth_use_case

import "context"

type AuthUseCase interface {
	Login(ctx context.Context, req LoginRequest) (LoginResponse, int, error)
	Register(ctx context.Context, req RegisterRequest) (UserInfo, int, error)
	RefreshToken(ctx context.Context, req RefreshTokenRequest) (LoginResponse, int, error)
}
