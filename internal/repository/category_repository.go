package repository

import (
	"todo-api/internal/domain"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll(userID uint) ([]domain.Category, error)
	FindByID(id, userID uint) (*domain.Category, error)
	Create(category *domain.Category) error
	Update(category *domain.Category) error
	Delete(id, userID uint) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db: db}
}

func (r *categoryRepository) FindAll(userID uint) ([]domain.Category, error) {
	var categories []domain.Category
	result := r.db.Where("user_id = ?", userID).Find(&categories)
	return categories, result.Error
}

func (r *categoryRepository) FindByID(id, userID uint) (*domain.Category, error) {
	var category domain.Category
	result := r.db.Where("id = ? AND user_id = ?", id, userID).First(&category)
	return &category, result.Error
}

func (r *categoryRepository) Create(category *domain.Category) error {
	result := r.db.Create(category)
	return result.Error
}

func (r *categoryRepository) Update(category *domain.Category) error {
	result := r.db.Save(category)
	return result.Error
}

func (r *categoryRepository) Delete(id, userID uint) error {
	result := r.db.Where("id = ? AND user_id = ?", id, userID).Delete(&domain.Category{})
	return result.Error
}
