package theme

import (
	"invitation-api/internal/database"
	themeDomain "invitation-api/internal/domain/theme"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for theme data operations
type Repository interface {
	Create(theme *themeDomain.Theme) error
	GetByID(id uuid.UUID) (*themeDomain.Theme, error)
	GetBySlug(slug string) (*themeDomain.Theme, error)
	GetBySlugUnscoped(slug string) (*themeDomain.Theme, error)
	Update(theme *themeDomain.Theme) error
	Delete(id uuid.UUID) error
	List(limit, offset int, activeOnly bool) ([]*themeDomain.Theme, error)
	ListByCategory(category string) ([]*themeDomain.Theme, error)
	Count(activeOnly bool) (int64, error)
}

// repository implements Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new theme repository
func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

// Create creates a new theme
func (r *repository) Create(theme *themeDomain.Theme) error {
	return r.db.Create(theme).Error
}

// GetByID retrieves a theme by ID
func (r *repository) GetByID(id uuid.UUID) (*themeDomain.Theme, error) {
	var theme themeDomain.Theme
	err := r.db.First(&theme, id).Error
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

// GetBySlug retrieves a theme by slug
func (r *repository) GetBySlug(slug string) (*themeDomain.Theme, error) {
	var theme themeDomain.Theme
	err := r.db.Where("slug = ?", slug).First(&theme).Error
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

// GetBySlugUnscoped retrieves a theme by slug including soft-deleted ones
func (r *repository) GetBySlugUnscoped(slug string) (*themeDomain.Theme, error) {
	var theme themeDomain.Theme
	err := r.db.Unscoped().Where("slug = ?", slug).First(&theme).Error
	if err != nil {
		return nil, err
	}
	return &theme, nil
}

// Update updates a theme
func (r *repository) Update(theme *themeDomain.Theme) error {
	return r.db.Save(theme).Error
}

// Delete soft deletes a theme
func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&themeDomain.Theme{}, id).Error
}

// List retrieves themes with pagination
func (r *repository) List(limit, offset int, activeOnly bool) ([]*themeDomain.Theme, error) {
	var themes []*themeDomain.Theme
	query := r.db.Model(&themeDomain.Theme{})

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Limit(limit).Offset(offset).Order("name").Find(&themes).Error
	return themes, err
}

// ListByCategory retrieves themes by category
func (r *repository) ListByCategory(category string) ([]*themeDomain.Theme, error) {
	var themes []*themeDomain.Theme
	err := r.db.Where("category = ? AND is_active = ?", category, true).Order("name").Find(&themes).Error
	return themes, err
}

// Count returns total number of themes
func (r *repository) Count(activeOnly bool) (int64, error) {
	var count int64
	query := r.db.Model(&themeDomain.Theme{})

	if activeOnly {
		query = query.Where("is_active = ?", true)
	}

	err := query.Count(&count).Error
	return count, err
}
