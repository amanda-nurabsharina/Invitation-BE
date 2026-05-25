package theme

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	themeDomain "invitation-api/internal/domain/theme"
	themeRepo "invitation-api/internal/repository/theme"

	"github.com/google/uuid"
)

// Service defines the interface for theme business operations
type Service interface {
	Create(name, slug, templatePath, category string, isPremium bool, customHTML, customCSS string) (*themeDomain.Theme, error)
	GetByID(id uuid.UUID) (*themeDomain.Theme, error)
	GetBySlug(slug string) (*themeDomain.Theme, error)
	Update(id uuid.UUID, data *UpdateThemeRequest) (*themeDomain.Theme, error)
	Delete(id uuid.UUID) error
	List(limit, offset int, activeOnly bool) ([]*themeDomain.Theme, int64, error)
	ListByCategory(category string) ([]*themeDomain.Theme, error)
	Activate(id uuid.UUID) error
	Deactivate(id uuid.UUID) error
}

// UpdateThemeRequest contains updateable theme fields
type UpdateThemeRequest struct {
	Name         *string                    `json:"name"`
	Description  *string                    `json:"description"`
	Category     *string                    `json:"category"`
	PreviewImage *string                    `json:"preview_image"`
	IsPremium    *bool                      `json:"is_premium"`
	TemplatePath *string                    `json:"template_path"`
	CustomHTML   *string                    `json:"custom_html"`
	CustomCSS    *string                    `json:"custom_css"`
	Settings     *themeDomain.ThemeSettings `json:"settings"`
	Colors       *themeDomain.ThemeColors   `json:"colors"`
	Fonts        *themeDomain.ThemeFonts    `json:"fonts"`
}

// service implements Service interface
type service struct {
	themeRepo themeRepo.Repository
}

// NewService creates a new theme service
func NewService(themeRepo themeRepo.Repository) Service {
	return &service{
		themeRepo: themeRepo,
	}
}

// Create creates a new theme
func (s *service) Create(name, slug, templatePath, category string, isPremium bool, customHTML, customCSS string) (*themeDomain.Theme, error) {
	// Check if slug is unique
	existing, _ := s.themeRepo.GetBySlug(slug)
	if existing != nil {
		return nil, errors.New("theme with this slug already exists")
	}

	theme, err := themeDomain.NewTheme(name, slug, templatePath)
	if err != nil {
		return nil, err
	}

	theme.Category = category
	theme.IsPremium = isPremium

	if customHTML != "" {
		theme.CustomHTML = &customHTML
	}
	if customCSS != "" {
		theme.CustomCSS = &customCSS
	}

	if err := s.themeRepo.Create(theme); err != nil {
		return nil, fmt.Errorf("failed to create theme: %w", err)
	}

	return theme, nil
}

// GetByID retrieves a theme by ID
func (s *service) GetByID(id uuid.UUID) (*themeDomain.Theme, error) {
	theme, err := s.themeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// If CustomHTML is empty but TemplatePath is set, try to load from file system
	if (theme.CustomHTML == nil || *theme.CustomHTML == "") && theme.TemplatePath != "" {
		// Use internal/templates/themes path standard
		htmlContent, err := os.ReadFile(filepath.Join("internal", "templates", "themes", theme.TemplatePath, "index.html"))
		if err == nil {
			content := string(htmlContent)
			theme.CustomHTML = &content
		} else {
			// Fallback: If template path fails, use slug?
			// Or maybe user meant `themes/` folder relative to root?
			// Let's stick to internal/templates first as confirmed by directory listing.
		}

		// Also try CSS (if exists in theme folder)
		cssContent, err := os.ReadFile(filepath.Join("internal", "templates", "themes", theme.TemplatePath, "style.css"))
		if err == nil {
			content := string(cssContent)
			theme.CustomCSS = &content
		}
	}

	return theme, nil
}

// GetBySlug retrieves a theme by slug
func (s *service) GetBySlug(slug string) (*themeDomain.Theme, error) {
	return s.themeRepo.GetBySlug(slug)
}

// Update updates a theme
func (s *service) Update(id uuid.UUID, data *UpdateThemeRequest) (*themeDomain.Theme, error) {
	theme, err := s.themeRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if data.Name != nil {
		theme.Name = *data.Name
	}
	if data.Description != nil {
		theme.Description = data.Description
	}
	if data.Category != nil {
		theme.Category = *data.Category
	}
	if data.PreviewImage != nil {
		theme.PreviewImage = data.PreviewImage
	}
	if data.IsPremium != nil {
		theme.IsPremium = *data.IsPremium
	}
	if data.TemplatePath != nil {
		theme.TemplatePath = *data.TemplatePath
	}
	if data.CustomHTML != nil {
		theme.CustomHTML = data.CustomHTML
	}
	if data.CustomCSS != nil {
		theme.CustomCSS = data.CustomCSS
	}
	if data.Settings != nil {
		theme.Settings = data.Settings
	}
	if data.Colors != nil {
		theme.Colors = data.Colors
	}
	if data.Fonts != nil {
		theme.Fonts = data.Fonts
	}

	if err := s.themeRepo.Update(theme); err != nil {
		return nil, fmt.Errorf("failed to update theme: %w", err)
	}

	return theme, nil
}

// Delete deletes a theme
func (s *service) Delete(id uuid.UUID) error {
	return s.themeRepo.Delete(id)
}

// List retrieves themes with pagination
func (s *service) List(limit, offset int, activeOnly bool) ([]*themeDomain.Theme, int64, error) {
	themes, err := s.themeRepo.List(limit, offset, activeOnly)
	if err != nil {
		return nil, 0, err
	}

	for _, theme := range themes {
		if (theme.CustomHTML == nil || *theme.CustomHTML == "") && theme.TemplatePath != "" {
			htmlContent, err := os.ReadFile(filepath.Join("internal", "templates", "themes", theme.TemplatePath, "index.html"))
			if err == nil {
				content := string(htmlContent)
				theme.CustomHTML = &content
			}
			cssContent, err := os.ReadFile(filepath.Join("internal", "templates", "themes", theme.TemplatePath, "style.css"))
			if err == nil {
				content := string(cssContent)
				theme.CustomCSS = &content
			}
		}
	}

	count, err := s.themeRepo.Count(activeOnly)
	if err != nil {
		return nil, 0, err
	}

	return themes, count, nil
}

// ListByCategory retrieves themes by category
func (s *service) ListByCategory(category string) ([]*themeDomain.Theme, error) {
	themes, err := s.themeRepo.ListByCategory(category)
	if err != nil {
		return nil, err
	}

	for _, theme := range themes {
		if (theme.CustomHTML == nil || *theme.CustomHTML == "") && theme.TemplatePath != "" {
			htmlContent, err := os.ReadFile(filepath.Join("internal", "templates", "themes", theme.TemplatePath, "index.html"))
			if err == nil {
				content := string(htmlContent)
				theme.CustomHTML = &content
			}
			cssContent, err := os.ReadFile(filepath.Join("internal", "templates", "themes", theme.TemplatePath, "style.css"))
			if err == nil {
				content := string(cssContent)
				theme.CustomCSS = &content
			}
		}
	}

	return themes, nil
}

// Activate activates a theme
func (s *service) Activate(id uuid.UUID) error {
	theme, err := s.themeRepo.GetByID(id)
	if err != nil {
		return err
	}

	theme.IsActive = true
	return s.themeRepo.Update(theme)
}

// Deactivate deactivates a theme
func (s *service) Deactivate(id uuid.UUID) error {
	theme, err := s.themeRepo.GetByID(id)
	if err != nil {
		return err
	}

	theme.IsActive = false
	return s.themeRepo.Update(theme)
}
