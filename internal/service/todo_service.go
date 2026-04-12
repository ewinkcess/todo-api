package service

import (
	"errors"
	"todo-api/internal/domain"
	"todo-api/internal/repository"
)

type TodoService interface {
	GetAllTodos(userID uint, query domain.PaginationQuery) (*domain.PaginationResponse, error)
	GetTodoByID(id, userID uint) (*domain.Todo, error)
	CreateTodo(userID uint, title, description string, categoryID *uint) (*domain.Todo, error)
	UpdateTodo(id, userID uint, title, description string, comleted bool) (*domain.Todo, error)
	DeleteTodo(id, userID uint) error
}

type todoService struct {
	repo repository.TodoRepository
}

func NewTodoService(repo repository.TodoRepository) TodoService {
	return &todoService{repo: repo}
}

func (s *todoService) GetAllTodos(userID uint, query domain.PaginationQuery) (*domain.PaginationResponse, error) {
	query.SetDefault()

	todos, total, err := s.repo.FindAll(userID, query)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / query.Limit
	if int(total)%query.Limit != 0 {
		totalPages++
	}

	response := &domain.PaginationResponse{
		Meta: domain.PaginationMeta{
			Page:       query.Page,
			Limit:      query.Limit,
			TotalData:  total,
			TotalPages: totalPages,
		},
		Data: todos,
	}
	return response, nil
}

func (s *todoService) GetTodoByID(id, userID uint) (*domain.Todo, error) {
	return s.repo.FindByID(id, userID)
}

func (s *todoService) CreateTodo(userID uint, title, description string, categoryID *uint) (*domain.Todo, error) {
	if title == "" {
		return nil, errors.New("judul todo tidak boleh kosong")
	}
	todo := &domain.Todo{
		UserID:      userID,
		CategoryID:  categoryID,
		Title:       title,
		Description: description,
		Completed:   false,
	}

	err := s.repo.Create(todo)
	if err != nil {
		return nil, err
	}
	return todo, nil
}

func (s *todoService) UpdateTodo(id, userID uint, title, description string, completed bool) (*domain.Todo, error) {
	todo, err := s.repo.FindByID(id, userID)
	if err != nil {
		return nil, errors.New("todo tidak ditemukan")
	}
	if title == "" {
		return nil, errors.New("judul todo tidak boleh kosong")
	}
	todo.Title = title
	todo.Description = description
	todo.Completed = completed
	err = s.repo.Update(todo)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

func (s *todoService) DeleteTodo(id, userID uint) error {
	_, err := s.repo.FindByID(id, userID)
	if err != nil {
		return errors.New("todo tidak ditemukan")
	}

	return s.repo.Delete(id, userID)
}
