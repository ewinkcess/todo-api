package service

import (
	"errors"
	"os"
	"testing"

	"todo-api/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByEmail(email string) (*domain.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(id uint) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

// ============================================
// Test Register
// ============================================

func TestRegister_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("FindByEmail", "ewink@gmail.com").Return(nil, errors.New("record not found"))
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	authService := NewAuthService(mockRepo)

	user, err := authService.Register("Ewink", "ewink@gmail.com", "password123")

	assert.NoError(t, err)
	assert.Equal(t, "Ewink", user.Name)
	assert.Equal(t, "ewink@gmail.com", user.Email)
	assert.NotEmpty(t, user.Password)
	mockRepo.AssertExpectations(t)
}

func TestRegister_EmailAlreadyExists(t *testing.T) {
	mockRepo := new(MockUserRepository)

	existingUser := &domain.User{
		ID:    1,
		Email: "ewink@gmail.com",
	}

	mockRepo.On("FindByEmail", "ewink@gmail.com").Return(existingUser, nil)

	authService := NewAuthService(mockRepo)

	user, err := authService.Register("Ewink", "ewink@gmail.com", "password123")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Equal(t, "email sudah terdaftar", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============================================
// Test Login
// ============================================

func TestLogin_Success(t *testing.T) {
	os.Setenv("JWT_SECRET", "rahasia_testing")

	mockRepo := new(MockUserRepository)

	existingUser := &domain.User{
		ID:       1,
		Email:    "ewink@gmail.com",
		Password: "$2a$14$g0cFkORF9qaksSgh.BXnlulsNPVf3jLUG3ON7UnpJFYCrOWtKBmM6",
	}

	mockRepo.On("FindByEmail", "ewink@gmail.com").Return(existingUser, nil)

	authService := NewAuthService(mockRepo)

	token, err := authService.Login("ewink@gmail.com", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	mockRepo.AssertExpectations(t)
}

func TestLogin_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)

	existingUser := &domain.User{
		ID:       1,
		Email:    "ewink@gmail.com",
		Password: "$2a$14$g0cFkORF9qaksSgh.BXnlulsNPVf3jLUG3ON7UnpJFYCrOWtKBmM6",
	}

	mockRepo.On("FindByEmail", "ewink@gmail.com").Return(existingUser, nil)

	authService := NewAuthService(mockRepo)

	token, err := authService.Login("ewink@gmail.com", "passwordsalah")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "email atau password salah", err.Error())
	mockRepo.AssertExpectations(t)
}

func TestLogin_EmailNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("FindByEmail", "tidakada@gmail.com").Return(nil, errors.New("record not found"))

	authService := NewAuthService(mockRepo)

	token, err := authService.Login("tidakada@gmail.com", "password123")

	assert.Error(t, err)
	assert.Empty(t, token)
	assert.Equal(t, "email atau password salah", err.Error())
	mockRepo.AssertExpectations(t)
}
