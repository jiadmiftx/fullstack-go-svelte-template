package auth_use_case

import (
	"context"
	"errors"
	"net/http"
	"time"

	"sekolah-madrasah/app/repository/user_repository"
	"sekolah-madrasah/pkg/auth_utils"

	"golang.org/x/crypto/bcrypt"
)

const (
	AccessTokenDuration  = 24 * time.Hour
	RefreshTokenDuration = 7 * 24 * time.Hour
)

type authUseCase struct {
	userRepo user_repository.UserRepository
}

func NewAuthUseCase(userRepo user_repository.UserRepository) AuthUseCase {
	return &authUseCase{userRepo: userRepo}
}

func (u *authUseCase) Login(ctx context.Context, req LoginRequest) (LoginResponse, int, error) {
	user, code, err := u.userRepo.GetUser(ctx, user_repository.UserFilter{
		Email: &req.Email,
	})
	if err != nil {
		if code == http.StatusNotFound {
			return LoginResponse{}, http.StatusUnauthorized, errors.New("invalid email or password")
		}
		return LoginResponse{}, code, err
	}

	if !user.IsActive {
		return LoginResponse{}, http.StatusForbidden, errors.New("account is not active")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return LoginResponse{}, http.StatusUnauthorized, errors.New("invalid email or password")
	}

	// Check if user is a warga with pending or rejected approval
	approvalStatus, err := u.userRepo.GetApprovalStatus(ctx, user.Id)
	if err == nil {
		if approvalStatus == "pending" {
			return LoginResponse{}, http.StatusForbidden, errors.New("pendaftaran Anda masih menunggu verifikasi dari pengurus")
		}
		if approvalStatus == "rejected" {
			return LoginResponse{}, http.StatusForbidden, errors.New("pendaftaran Anda telah ditolak, silakan hubungi pengurus RT")
		}
	}

	accessToken, err := auth_utils.GenerateToken(auth_utils.TokenParams{
		UserID: user.Id,
	}, AccessTokenDuration)
	if err != nil {
		return LoginResponse{}, http.StatusInternalServerError, err
	}

	refreshToken, err := auth_utils.GenerateToken(auth_utils.TokenParams{
		UserID: user.Id,
	}, RefreshTokenDuration)
	if err != nil {
		return LoginResponse{}, http.StatusInternalServerError, err
	}

	u.userRepo.UpdateLastLogin(ctx, user_repository.UserFilter{Id: &user.Id})

	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(AccessTokenDuration).Unix(),
		User: UserInfo{
			Id:           user.Id,
			Email:        user.Email,
			FullName:     user.FullName,
			IsSuperAdmin: user.IsSuperAdmin,
			IsActive:     user.IsActive,
			LastLoginAt:  user.LastLoginAt,
		},
	}, http.StatusOK, nil
}

func (u *authUseCase) Register(ctx context.Context, req RegisterRequest) (UserInfo, int, error) {
	existingUser, code, _ := u.userRepo.GetUser(ctx, user_repository.UserFilter{
		Email: &req.Email,
	})
	if code == http.StatusOK && existingUser.Id.String() != "" {
		return UserInfo{}, http.StatusConflict, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return UserInfo{}, http.StatusInternalServerError, err
	}

	newUser := user_repository.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		IsActive: true,
	}

	createdUser, code, err := u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return UserInfo{}, code, err
	}

	return UserInfo{
		Id:           createdUser.Id,
		Email:        createdUser.Email,
		FullName:     createdUser.FullName,
		IsSuperAdmin: createdUser.IsSuperAdmin,
		IsActive:     createdUser.IsActive,
	}, http.StatusCreated, nil
}

func (u *authUseCase) RefreshToken(ctx context.Context, req RefreshTokenRequest) (LoginResponse, int, error) {
	claims, err := auth_utils.ValidateToken(req.RefreshToken)
	if err != nil {
		return LoginResponse{}, http.StatusUnauthorized, errors.New("invalid refresh token")
	}

	user, code, err := u.userRepo.GetUser(ctx, user_repository.UserFilter{
		Id: &claims.UserID,
	})
	if err != nil {
		return LoginResponse{}, code, err
	}

	if !user.IsActive {
		return LoginResponse{}, http.StatusForbidden, errors.New("account is not active")
	}

	accessToken, err := auth_utils.GenerateToken(auth_utils.TokenParams{
		UserID: user.Id,
	}, AccessTokenDuration)
	if err != nil {
		return LoginResponse{}, http.StatusInternalServerError, err
	}

	refreshToken, err := auth_utils.GenerateToken(auth_utils.TokenParams{
		UserID: user.Id,
	}, RefreshTokenDuration)
	if err != nil {
		return LoginResponse{}, http.StatusInternalServerError, err
	}

	return LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(AccessTokenDuration).Unix(),
		User: UserInfo{
			Id:           user.Id,
			Email:        user.Email,
			FullName:     user.FullName,
			IsSuperAdmin: user.IsSuperAdmin,
			IsActive:     user.IsActive,
			LastLoginAt:  user.LastLoginAt,
		},
	}, http.StatusOK, nil
}
