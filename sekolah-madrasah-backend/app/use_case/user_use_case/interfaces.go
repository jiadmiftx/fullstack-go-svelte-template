package user_use_case

import (
	"context"

	"github.com/google/uuid"
	"sekolah-madrasah/pkg/paginate_utils"
)

type UserUseCase interface {
	GetUser(ctx context.Context, id uuid.UUID) (User, int, error)
	GetUsers(ctx context.Context, filter UserFilter, paginate *paginate_utils.PaginateData) ([]User, int, error)
	CreateUser(ctx context.Context, req CreateUserRequest) (User, int, error)
	UpdateUser(ctx context.Context, id uuid.UUID, req UpdateUserRequest) (User, int, error)
	DeleteUser(ctx context.Context, id uuid.UUID) (int, error)
}
