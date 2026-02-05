package auth_use_case

import (
	"context"
	"errors"
	"testing"
	"time"

	"sekolah-madrasah/app/repository/user_repository"
	"sekolah-madrasah/pkg/paginate_utils"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	users []user_repository.User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{users: []user_repository.User{}}
}

func (m *MockUserRepository) GetUser(ctx context.Context, filter user_repository.UserFilter) (user_repository.User, int, error) {
	for _, u := range m.users {
		if filter.Id != nil && u.Id == *filter.Id {
			return u, 200, nil
		}
		if filter.Email != nil && u.Email == *filter.Email {
			return u, 200, nil
		}
	}
	return user_repository.User{}, 404, errors.New("user not found")
}

func (m *MockUserRepository) GetUsers(ctx context.Context, filter user_repository.UserFilter, paginate *paginate_utils.PaginateData) ([]user_repository.User, int, error) {
	return m.users, 200, nil
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user user_repository.User) (user_repository.User, int, error) {
	if user.Id == uuid.Nil {
		user.Id = uuid.New()
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users = append(m.users, user)
	return user, 201, nil
}

func (m *MockUserRepository) UpdateUser(ctx context.Context, filter user_repository.UserFilter, user user_repository.User) (int, error) {
	for i, u := range m.users {
		if filter.Id != nil && u.Id == *filter.Id {
			m.users[i].FullName = user.FullName
			m.users[i].UpdatedAt = time.Now()
			return 200, nil
		}
	}
	return 404, nil
}

func (m *MockUserRepository) DeleteUser(ctx context.Context, filter user_repository.UserFilter) (int, error) {
	for i, u := range m.users {
		if filter.Id != nil && u.Id == *filter.Id {
			m.users = append(m.users[:i], m.users[i+1:]...)
			return 200, nil
		}
	}
	return 404, nil
}

func (m *MockUserRepository) UpdateLastLogin(ctx context.Context, filter user_repository.UserFilter) (int, error) {
	return 200, nil
}

func (m *MockUserRepository) HasPendingApproval(ctx context.Context, userId uuid.UUID) (bool, error) {
	return false, nil // Mock always returns not pending
}

func (m *MockUserRepository) GetApprovalStatus(ctx context.Context, userId uuid.UUID) (string, error) {
	return "approved", nil // Mock always returns approved
}

func TestAuthUseCase_Register(t *testing.T) {
	mockRepo := NewMockUserRepository()
	useCase := NewAuthUseCase(mockRepo)
	ctx := context.Background()

	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	userInfo, code, err := useCase.Register(ctx, req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if code != 201 {
		t.Errorf("Expected status 201, got %d", code)
	}

	if userInfo.Email != req.Email {
		t.Errorf("Expected email %s, got %s", req.Email, userInfo.Email)
	}

	if userInfo.Id == uuid.Nil {
		t.Error("Expected user ID to be generated")
	}

	if !userInfo.IsActive {
		t.Error("Expected user to be active")
	}
}

func TestAuthUseCase_RegisterDuplicateEmail(t *testing.T) {
	mockRepo := NewMockUserRepository()
	useCase := NewAuthUseCase(mockRepo)
	ctx := context.Background()

	req := RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		FullName: "Test User",
	}

	_, _, _ = useCase.Register(ctx, req)
	_, code, err := useCase.Register(ctx, req)

	if code != 409 {
		t.Errorf("Expected status 409 for duplicate email, got %d", code)
	}

	if err == nil {
		t.Error("Expected error for duplicate email")
	}
}

func TestAuthUseCase_Login(t *testing.T) {
	mockRepo := NewMockUserRepository()
	useCase := NewAuthUseCase(mockRepo)
	ctx := context.Background()

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockRepo.users = append(mockRepo.users, user_repository.User{
		Id:       uuid.New(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
		FullName: "Test User",
		IsActive: true,
	})

	req := LoginRequest{
		Email:    "test@example.com",
		Password: password,
	}

	resp, code, err := useCase.Login(ctx, req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if code != 200 {
		t.Errorf("Expected status 200, got %d", code)
	}

	if resp.AccessToken == "" {
		t.Error("Expected access token to be generated")
	}

	if resp.RefreshToken == "" {
		t.Error("Expected refresh token to be generated")
	}

	if resp.User.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", resp.User.Email)
	}
}

func TestAuthUseCase_LoginWrongPassword(t *testing.T) {
	mockRepo := NewMockUserRepository()
	useCase := NewAuthUseCase(mockRepo)
	ctx := context.Background()

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	mockRepo.users = append(mockRepo.users, user_repository.User{
		Id:       uuid.New(),
		Email:    "test@example.com",
		Password: string(hashedPassword),
		FullName: "Test User",
		IsActive: true,
	})

	req := LoginRequest{
		Email:    "test@example.com",
		Password: "wrongpassword",
	}

	_, code, err := useCase.Login(ctx, req)

	if code != 401 {
		t.Errorf("Expected status 401 for wrong password, got %d", code)
	}

	if err == nil {
		t.Error("Expected error for wrong password")
	}
}

func TestAuthUseCase_LoginUserNotFound(t *testing.T) {
	mockRepo := NewMockUserRepository()
	useCase := NewAuthUseCase(mockRepo)
	ctx := context.Background()

	req := LoginRequest{
		Email:    "nonexistent@example.com",
		Password: "password123",
	}

	_, code, err := useCase.Login(ctx, req)

	// Login returns 401 for invalid credentials (user not found or wrong password)
	if code != 401 {
		t.Errorf("Expected status 401 for user not found, got %d", code)
	}

	if err == nil {
		t.Error("Expected error for user not found")
	}
}

func TestAuthUseCase_LoginInactiveUser(t *testing.T) {
	mockRepo := NewMockUserRepository()
	useCase := NewAuthUseCase(mockRepo)
	ctx := context.Background()

	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockRepo.users = append(mockRepo.users, user_repository.User{
		Id:       uuid.New(),
		Email:    "inactive@example.com",
		Password: string(hashedPassword),
		FullName: "Inactive User",
		IsActive: false,
	})

	req := LoginRequest{
		Email:    "inactive@example.com",
		Password: password,
	}

	_, code, err := useCase.Login(ctx, req)

	if code != 403 {
		t.Errorf("Expected status 403 for inactive user, got %d", code)
	}

	if err == nil {
		t.Error("Expected error for inactive user")
	}
}
