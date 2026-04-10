package repository

import (
	"todo-api/internal/domain"

	"gorm.io/gorm"
)

type TodoRepository interface {
	FindAll(userID uint) ([]domain.Todo, error)
	FindByID(id, userID uint) (*domain.Todo, error)
	Create(todo *domain.Todo) error
	Update(todo *domain.Todo) error
	Delete(id, userID uint) error
}

type todoRepository struct {
	db *gorm.DB
}

func NewTodoRepository(db *gorm.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) FindAll(userID uint) ([]domain.Todo, error) {
	var todos []domain.Todo
	result := r.db.Where("user_id = ?", userID).Find(&todos)
	return todos, result.Error
}

func (r *todoRepository) FindByID(id, userID uint) (*domain.Todo, error) {
	var todo domain.Todo
	result := r.db.Where("id = ? AND user_id = ?", id, userID).First(&todo)
	return &todo, result.Error
}

func (r *todoRepository) Create(todo *domain.Todo) error {
	result := r.db.Create(todo)
	return result.Error
}

func (r *todoRepository) Update(todo *domain.Todo) error {
	result := r.db.Save(todo)
	return result.Error
}

func (r *todoRepository) Delete(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id", id, userID).Delete(&domain.Todo{})
	return result.Error
}
