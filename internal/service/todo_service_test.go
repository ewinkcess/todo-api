package service

import (
	"errors"
	"testing"

	"todo-api/internal/domain"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockTodoRepository struct {
	mock.Mock
}

func (m *MockTodoRepository) FindAll(userID uint, query domain.PaginationQuery) ([]domain.Todo, int64, error) {
	args := m.Called(userID, query)
	return args.Get(0).([]domain.Todo), args.Get(1).(int64), args.Error(2)
}

func (m *MockTodoRepository) FindByID(id, userID uint) (*domain.Todo, error) {
	args := m.Called(id, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Todo), args.Error(1)
}

func (m *MockTodoRepository) Create(todo *domain.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Update(todo *domain.Todo) error {
	args := m.Called(todo)
	return args.Error(0)
}

func (m *MockTodoRepository) Delete(id, userID uint) error {
	args := m.Called(id, userID)
	return args.Error(0)
}

// ============================================
// Test GetAllTodos
// ============================================

func TestGetAllTodos_Success(t *testing.T) {
	mockRepo := new(MockTodoRepository)

	expectedTodos := []domain.Todo{
		{ID: 1, UserID: 1, Title: "Belajar Golang", Completed: false},
		{ID: 2, UserID: 1, Title: "Belajar Testing", Completed: true},
	}

	query := domain.PaginationQuery{Page: 1, Limit: 10}
	mockRepo.On("FindAll", uint(1), query).Return(expectedTodos, int64(2), nil)

	todoService := NewTodoService(mockRepo)

	result, err := todoService.GetAllTodos(uint(1), query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(2), result.Meta.TotalData)
	assert.Equal(t, 1, result.Meta.TotalPages)
	assert.Equal(t, expectedTodos, result.Data)
	mockRepo.AssertExpectations(t)
}

func TestGetAllTodos_WithSearch(t *testing.T) {
	mockRepo := new(MockTodoRepository)

	expectedTodos := []domain.Todo{
		{ID: 1, UserID: 1, Title: "Belajar Golang", Completed: false},
	}

	query := domain.PaginationQuery{Page: 1, Limit: 10, Search: "golang"}
	mockRepo.On("FindAll", uint(1), query).Return(expectedTodos, int64(1), nil)

	todoService := NewTodoService(mockRepo)

	result, err := todoService.GetAllTodos(uint(1), query)

	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, int64(1), result.Meta.TotalData)
	mockRepo.AssertExpectations(t)
}

// ============================================
// Test CreateTodo
// ============================================

func TestCreateTodo_Success(t *testing.T) {
	mockRepo := new(MockTodoRepository)

	mockRepo.On("Create", mock.AnythingOfType("*domain.Todo")).Return(nil)

	todoService := NewTodoService(mockRepo)

	todo, err := todoService.CreateTodo(uint(1), "Belajar Golang", "Deskripsi todo", nil)

	assert.NoError(t, err)
	assert.Equal(t, "Belajar Golang", todo.Title)
	assert.Equal(t, "Deskripsi todo", todo.Description)
	assert.Equal(t, uint(1), todo.UserID)
	assert.Nil(t, todo.CategoryID)
	assert.False(t, todo.Completed)
	mockRepo.AssertExpectations(t)
}

func TestCreateTodo_WithCategory(t *testing.T) {
	mockRepo := new(MockTodoRepository)

	mockRepo.On("Create", mock.AnythingOfType("*domain.Todo")).Return(nil)

	todoService := NewTodoService(mockRepo)

	categoryID := uint(1)
	todo, err := todoService.CreateTodo(uint(1), "Belajar Golang", "Deskripsi todo", &categoryID)

	assert.NoError(t, err)
	assert.Equal(t, "Belajar Golang", todo.Title)
	assert.Equal(t, &categoryID, todo.CategoryID)
	mockRepo.AssertExpectations(t)
}

func TestCreateTodo_EmptyTitle(t *testing.T) {
	mockRepo := new(MockTodoRepository)
	todoService := NewTodoService(mockRepo)

	todo, err := todoService.CreateTodo(uint(1), "", "Deskripsi todo", nil)

	assert.Error(t, err)
	assert.Nil(t, todo)
	assert.Equal(t, "judul todo tidak boleh kosong", err.Error())
}

// ============================================
// Test UpdateTodo
// ============================================

func TestUpdateTodo_Success(t *testing.T) {
	mockRepo := new(MockTodoRepository)

	existingTodo := &domain.Todo{
		ID:        1,
		UserID:    1,
		Title:     "Todo Lama",
		Completed: false,
	}

	mockRepo.On("FindByID", uint(1), uint(1)).Return(existingTodo, nil)
	mockRepo.On("Update", mock.AnythingOfType("*domain.Todo")).Return(nil)

	todoService := NewTodoService(mockRepo)

	todo, err := todoService.UpdateTodo(uint(1), uint(1), "Todo Baru", "Deskripsi baru", true)

	assert.NoError(t, err)
	assert.Equal(t, "Todo Baru", todo.Title)
	assert.Equal(t, true, todo.Completed)
	mockRepo.AssertExpectations(t)
}

func TestUpdateTodo_NotFound(t *testing.T) {
	mockRepo := new(MockTodoRepository)

	mockRepo.On("FindByID", uint(999), uint(1)).Return(nil, errors.New("record not found"))

	todoService := NewTodoService(mockRepo)

	todo, err := todoService.UpdateTodo(uint(999), uint(1), "Todo Baru", "Deskripsi baru", true)

	assert.Error(t, err)
	assert.Nil(t, todo)
	assert.Equal(t, "todo tidak ditemukan", err.Error())
	mockRepo.AssertExpectations(t)
}

// ============================================
// Test DeleteTodo
// ============================================

func TestDeleteTodo_Success(t *testing.T) {
	mockRepo := new(MockTodoRepository)

	existingTodo := &domain.Todo{
		ID:     1,
		UserID: 1,
		Title:  "Todo yang akan dihapus",
	}

	mockRepo.On("FindByID", uint(1), uint(1)).Return(existingTodo, nil)
	mockRepo.On("Delete", uint(1), uint(1)).Return(nil)

	todoService := NewTodoService(mockRepo)

	err := todoService.DeleteTodo(uint(1), uint(1))

	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestDeleteTodo_NotFound(t *testing.T) {
	mockRepo := new(MockTodoRepository)

	mockRepo.On("FindByID", uint(999), uint(1)).Return(nil, errors.New("record not found"))

	todoService := NewTodoService(mockRepo)

	err := todoService.DeleteTodo(uint(999), uint(1))

	assert.Error(t, err)
	assert.Equal(t, "todo tidak ditemukan", err.Error())
	mockRepo.AssertExpectations(t)
}
