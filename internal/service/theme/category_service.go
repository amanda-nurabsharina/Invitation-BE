package theme

import (
	"errors"

	"invitation-api/internal/domain/theme"
	repo "invitation-api/internal/repository/theme"

	"github.com/google/uuid"
)

type CategoryService interface {
	Create(req *CreateCategoryRequest) (*theme.ThemeCategory, error)
	Update(id uuid.UUID, req *UpdateCategoryRequest) (*theme.ThemeCategory, error)
	Delete(id uuid.UUID) error
	GetByID(id uuid.UUID) (*theme.ThemeCategory, error)
	GetAll() ([]theme.ThemeCategory, error)
}

type categoryService struct {
	repo repo.CategoryRepository
}

func NewCategoryService(repo repo.CategoryRepository) CategoryService {
	return &categoryService{repo: repo}
}

type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	DefaultHTML string `json:"default_html"`
	DefaultCSS  string `json:"default_css"`
	Description string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	DefaultHTML string `json:"default_html"`
	DefaultCSS  string `json:"default_css"`
	Description string `json:"description"`
}

func (s *categoryService) Create(req *CreateCategoryRequest) (*theme.ThemeCategory, error) {
	// Check if slug exists
	existing, _ := s.repo.GetBySlug(req.Slug)
	if existing != nil {
		return nil, errors.New("slug already exists")
	}

	category := theme.NewThemeCategory(req.Name, req.Slug, req.DefaultHTML, req.DefaultCSS)
	category.Description = req.Description

	if err := s.repo.Create(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) Update(id uuid.UUID, req *UpdateCategoryRequest) (*theme.ThemeCategory, error) {
	category, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if category == nil {
		return nil, errors.New("category not found")
	}

	category.Name = req.Name
	category.DefaultHTML = req.DefaultHTML
	category.DefaultCSS = req.DefaultCSS
	category.Description = req.Description

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}
	return category, nil
}

func (s *categoryService) Delete(id uuid.UUID) error {
	return s.repo.Delete(id)
}

func (s *categoryService) GetByID(id uuid.UUID) (*theme.ThemeCategory, error) {
	return s.repo.GetByID(id)
}

func (s *categoryService) GetAll() ([]theme.ThemeCategory, error) {
	return s.repo.GetAll()
}
