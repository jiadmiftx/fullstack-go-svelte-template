package user_repository

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUserModel(t *testing.T) {
	now := time.Now()
	user := User{
		Id:        uuid.New(),
		Email:     "test@example.com",
		Password:  "hashedpassword",
		FullName:  "John Doe",
		IsActive:  true,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if user.Email != "test@example.com" {
		t.Errorf("Expected email test@example.com, got %s", user.Email)
	}

	if !user.IsActive {
		t.Error("Expected user to be active")
	}

	if user.Id == uuid.Nil {
		t.Error("Expected user ID to be set")
	}

	if user.FullName != "John Doe" {
		t.Errorf("Expected FullName 'John Doe', got %s", user.FullName)
	}
}

func TestUserFilter(t *testing.T) {
	email := "test@example.com"
	isActive := true
	
	filter := UserFilter{
		Email:    &email,
		IsActive: &isActive,
	}

	if *filter.Email != email {
		t.Errorf("Expected email %s, got %s", email, *filter.Email)
	}

	if *filter.IsActive != isActive {
		t.Errorf("Expected IsActive %v, got %v", isActive, *filter.IsActive)
	}
}

func TestUserFilterNilFields(t *testing.T) {
	filter := UserFilter{}

	if filter.Id != nil {
		t.Error("Expected Id to be nil")
	}

	if filter.Email != nil {
		t.Error("Expected Email to be nil")
	}

	if filter.IsActive != nil {
		t.Error("Expected IsActive to be nil")
	}
}

type MockUserRepository struct {
	users []User
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: []User{},
	}
}

func (m *MockUserRepository) CreateUser(ctx context.Context, user User) (User, int, error) {
	if user.Id == uuid.Nil {
		user.Id = uuid.New()
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	m.users = append(m.users, user)
	return user, 201, nil
}

func (m *MockUserRepository) GetUser(ctx context.Context, filter UserFilter) (User, int, error) {
	for _, u := range m.users {
		if filter.Id != nil && u.Id == *filter.Id {
			return u, 200, nil
		}
		if filter.Email != nil && u.Email == *filter.Email {
			return u, 200, nil
		}
	}
	return User{}, 404, nil
}

func TestMockUserRepository_CreateAndGet(t *testing.T) {
	repo := NewMockUserRepository()
	ctx := context.Background()

	newUser := User{
		Email:    "test@example.com",
		Password: "hashed",
		FullName: "Test User",
		IsActive: true,
	}

	created, code, err := repo.CreateUser(ctx, newUser)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if code != 201 {
		t.Errorf("Expected status code 201, got %d", code)
	}

	if created.Id == uuid.Nil {
		t.Error("Expected ID to be generated")
	}

	email := created.Email
	found, code, _ := repo.GetUser(ctx, UserFilter{Email: &email})
	if code != 200 {
		t.Errorf("Expected status 200, got %d", code)
	}

	if found.Id != created.Id {
		t.Error("Expected to find the same user")
	}
}
