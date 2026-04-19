package theme

import (
	"errors"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"invitation-api/internal/database"
	"invitation-api/internal/domain/theme"
)

type CategoryRepository interface {
	Create(category *theme.ThemeCategory) error
	Update(category *theme.ThemeCategory) error
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*theme.ThemeCategory, error)
	GetBySlug(slug string) (*theme.ThemeCategory, error)
	GetAll() ([]theme.ThemeCategory, error)
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
	return &categoryRepository{db: database.GetDB()}
}

func (r *categoryRepository) Create(category *theme.ThemeCategory) error {
	return r.db.Create(category).Error
}

func (r *categoryRepository) Update(category *theme.ThemeCategory) error {
	return r.db.Save(category).Error
}

func (r *categoryRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&theme.ThemeCategory{}, id).Error
}

func (r *categoryRepository) GetByID(id uuid.UUID) (*theme.ThemeCategory, error) {
	var category theme.ThemeCategory
	err := r.db.First(&category, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // Or return custom error
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetBySlug(slug string) (*theme.ThemeCategory, error) {
	var category theme.ThemeCategory
	err := r.db.First(&category, "slug = ?", slug).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &category, nil
}

func (r *categoryRepository) GetAll() ([]theme.ThemeCategory, error) {
	var categories []theme.ThemeCategory
	err := r.db.Find(&categories).Error
	return categories, err
}
