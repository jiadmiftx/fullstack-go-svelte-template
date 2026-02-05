package user_repository

import (
	"context"

	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
)

type UserRepository interface {
	GetUser(ctx context.Context, filter UserFilter) (User, int, error)
	GetUsers(ctx context.Context, filter UserFilter, paginate *paginate_utils.PaginateData) ([]User, int, error)
	CreateUser(ctx context.Context, user User) (User, int, error)
	UpdateUser(ctx context.Context, filter UserFilter, user User) (int, error)
	DeleteUser(ctx context.Context, filter UserFilter) (int, error)
	UpdateLastLogin(ctx context.Context, filter UserFilter) (int, error)
	HasPendingApproval(ctx context.Context, userId uuid.UUID) (bool, error)
	GetApprovalStatus(ctx context.Context, userId uuid.UUID) (string, error)
}
