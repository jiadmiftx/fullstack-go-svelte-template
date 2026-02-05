package user_use_case

import (
	"context"
	"errors"
	"net/http"

	"sekolah-madrasah/app/repository/user_repository"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userUseCase struct {
	userRepo user_repository.UserRepository
}

func NewUserUseCase(userRepo user_repository.UserRepository) UserUseCase {
	return &userUseCase{userRepo: userRepo}
}

func (u *userUseCase) toUser(repoUser user_repository.User) User {
	return User{
		Id:              repoUser.Id,
		Email:           repoUser.Email,
		FullName:        repoUser.FullName,
		Phone:           repoUser.Phone,
		Avatar:          repoUser.Avatar,
		IsSuperAdmin:    repoUser.IsSuperAdmin,
		IsActive:        repoUser.IsActive,
		EmailVerifiedAt: repoUser.EmailVerifiedAt,
		LastLoginAt:     repoUser.LastLoginAt,
		CreatedAt:       repoUser.CreatedAt,
		UpdatedAt:       repoUser.UpdatedAt,
	}
}

func (u *userUseCase) GetUser(ctx context.Context, id uuid.UUID) (User, int, error) {
	user, code, err := u.userRepo.GetUser(ctx, user_repository.UserFilter{
		Id: &id,
	})
	if err != nil {
		return User{}, code, err
	}

	return u.toUser(user), http.StatusOK, nil
}

func (u *userUseCase) GetUsers(ctx context.Context, filter UserFilter, paginate *paginate_utils.PaginateData) ([]User, int, error) {
	repoFilter := user_repository.UserFilter{
		Id:           filter.Id,
		Email:        filter.Email,
		IsSuperAdmin: filter.IsSuperAdmin,
		IsActive:     filter.IsActive,
		PlatformOnly: filter.PlatformOnly,
	}

	users, code, err := u.userRepo.GetUsers(ctx, repoFilter, paginate)
	if err != nil {
		return nil, code, err
	}

	result := make([]User, len(users))
	for i, user := range users {
		result[i] = u.toUser(user)
	}

	return result, http.StatusOK, nil
}

func (u *userUseCase) CreateUser(ctx context.Context, req CreateUserRequest) (User, int, error) {
	existingUser, code, _ := u.userRepo.GetUser(ctx, user_repository.UserFilter{
		Email: &req.Email,
	})
	if code == http.StatusOK && existingUser.Id != uuid.Nil {
		return User{}, http.StatusConflict, errors.New("email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, http.StatusInternalServerError, err
	}

	newUser := user_repository.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		FullName: req.FullName,
		Phone:    req.Phone,
		IsActive: true,
	}

	createdUser, code, err := u.userRepo.CreateUser(ctx, newUser)
	if err != nil {
		return User{}, code, err
	}

	return u.toUser(createdUser), http.StatusCreated, nil
}

func (u *userUseCase) UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (User, int, error) {
	existingUser, code, err := u.userRepo.GetUser(ctx, user_repository.UserFilter{
		Id: &id,
	})
	if err != nil {
		return User{}, code, err
	}

	if req.Email != "" && req.Email != existingUser.Email {
		checkUser, checkCode, _ := u.userRepo.GetUser(ctx, user_repository.UserFilter{
			Email: &req.Email,
		})
		if checkCode == http.StatusOK && checkUser.Id != uuid.Nil {
			return User{}, http.StatusConflict, errors.New("email already registered")
		}
	}

	updateData := user_repository.User{
		Email:    req.Email,
		FullName: req.FullName,
		Phone:    req.Phone,
		Avatar:   req.Avatar,
	}

	code, err = u.userRepo.UpdateUser(ctx, user_repository.UserFilter{Id: &id}, updateData)
	if err != nil {
		return User{}, code, err
	}

	updatedUser, code, err := u.userRepo.GetUser(ctx, user_repository.UserFilter{Id: &id})
	if err != nil {
		return User{}, code, err
	}

	return u.toUser(updatedUser), http.StatusOK, nil
}

func (u *userUseCase) DeleteUser(ctx context.Context, id uuid.UUID) (int, error) {
	user, code, err := u.userRepo.GetUser(ctx, user_repository.UserFilter{
		Id: &id,
	})
	if err != nil {
		return code, err
	}

	// Prevent deletion of SuperAdmin users
	if user.IsSuperAdmin {
		return http.StatusForbidden, errors.New("cannot delete super admin users")
	}

	return u.userRepo.DeleteUser(ctx, user_repository.UserFilter{Id: &id})
}
