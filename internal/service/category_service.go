package service

import (
	"errors"

	"todo-api/internal/domain"
	"todo-api/internal/repository"
)

type CategoryService interface {
	GetAllCategories(userID uint) ([]domain.Category, error)
	GetCategoryByID(id, userID uint) (*domain.Category, error)
	CreateCategory(userID uint, name, color string) (*domain.Category, error)
	UpdateCategory(id, userID uint, name, color string) (*domain.Category, error)
	DeleteCategory(id, userID uint) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

func (s *categoryService) GetAllCategories(userID uint) ([]domain.Category, error) {
	return s.repo.FindAll(userID)
}

func (s *categoryService) GetCategoryByID(id, userID uint) (*domain.Category, error) {
	return s.repo.FindByID(id, userID)
}

func (s *categoryService) CreateCategory(userID uint, name, color string) (*domain.Category, error) {
	if name == "" {
		return nil, errors.New("nama kategori tidak boleh kosong")
	}

	if color == "" {
		color = "#000000"
	}

	category := &domain.Category{
		UserID: userID,
		Name:   name,
		Color:  color,
	}

	err := s.repo.Create(category)
	if err != nil {
		return nil, errors.New("gagal membuat kategori")
	}

	return category, nil
}

func (s *categoryService) UpdateCategory(id, userID uint, name, color string) (*domain.Category, error) {
	category, err := s.repo.FindByID(id, userID)
	if err != nil {
		return nil, errors.New("kategori tidak ditemukan")
	}
	if name == "" {
		return nil, errors.New("nama kategori tidak boleh kosong")
	}
	category.Name = name
	category.Color = color
	err = s.repo.Update(category)
	if err != nil {
		return nil, errors.New("gagal mengupdate kategori")
	}
	return category, nil
}

func (s *categoryService) DeleteCategory(id, userID uint) error {
	_, err := s.repo.FindByID(id, userID)
	if err != nil {
		return errors.New("kategori tidak ditemukan")
	}
	return s.repo.Delete(id, userID)
}
