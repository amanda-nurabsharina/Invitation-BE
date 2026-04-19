package theme

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Theme represents an invitation template theme
type Theme struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	Name         string         `gorm:"column:name;not null" json:"name"`
	Slug         string         `gorm:"column:slug;uniqueIndex;not null" json:"slug"`
	Description  *string        `gorm:"column:description" json:"description,omitempty"`
	Category     string         `gorm:"column:category;default:general" json:"category"` // elegant, rustic, modern, floral, minimalist
	PreviewImage *string        `gorm:"column:preview_image" json:"preview_image,omitempty"`
	TemplatePath string         `gorm:"column:template_path" json:"template_path"` // Made optional if CustomHTML is present
	CustomHTML   *string        `gorm:"column:custom_html;type:text" json:"custom_html,omitempty"`
	CustomCSS    *string        `gorm:"column:custom_css;type:text" json:"custom_css,omitempty"`
	IsPremium    bool           `gorm:"column:is_premium;default:false" json:"is_premium"`
	IsActive     bool           `gorm:"column:is_active;default:true" json:"is_active"`
	Settings     *ThemeSettings `gorm:"serializer:json;column:settings" json:"settings,omitempty"`
	Colors       *ThemeColors   `gorm:"serializer:json;column:colors" json:"colors,omitempty"`
	Fonts        *ThemeFonts    `gorm:"serializer:json;column:fonts" json:"fonts,omitempty"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// ThemeCategory represents a master template category
type ThemeCategory struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	Name        string         `gorm:"column:name;not null" json:"name"`
	Slug        string         `gorm:"column:slug;uniqueIndex;not null" json:"slug"`
	Description string         `gorm:"column:description" json:"description"`
	DefaultHTML string         `gorm:"column:default_html;type:text" json:"default_html"`
	DefaultCSS  string         `gorm:"column:default_css;type:text" json:"default_css"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func NewThemeCategory(name, slug, defaultHTML, defaultCSS string) *ThemeCategory {
	return &ThemeCategory{
		ID:          uuid.New(),
		Name:        name,
		Slug:        slug,
		DefaultHTML: defaultHTML,
		DefaultCSS:  defaultCSS,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// ThemeSettings holds theme-specific settings
type ThemeSettings struct {
	HasMusic          bool   `json:"has_music"`
	HasGallery        bool   `json:"has_gallery"`
	HasCountdown      bool   `json:"has_countdown"`
	HasSnowEffect     bool   `json:"has_snow_effect"`
	HasFloatingHearts bool   `json:"has_floating_hearts"`
	AnimationStyle    string `json:"animation_style,omitempty"` // fade, slide, zoom
	GalleryStyle      string `json:"gallery_style,omitempty"`   // grid, carousel, masonry
}

// ThemeColors holds theme color palette
type ThemeColors struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Accent     string `json:"accent"`
	Background string `json:"background"`
	Text       string `json:"text"`
	TextMuted  string `json:"text_muted"`
}

// ThemeFonts holds theme typography
type ThemeFonts struct {
	Heading    string `json:"heading"`
	Body       string `json:"body"`
	Accent     string `json:"accent"`
	HeadingURL string `json:"heading_url,omitempty"`
	BodyURL    string `json:"body_url,omitempty"`
	AccentURL  string `json:"accent_url,omitempty"`
}

// Category constants
const (
	CategoryElegant     = "elegant"
	CategoryRustic      = "rustic"
	CategoryModern      = "modern"
	CategoryFloral      = "floral"
	CategoryMinimalist  = "minimalist"
	CategoryTraditional = "traditional"
)

// NewTheme creates a new theme instance
func NewTheme(name, slug, templatePath string) (*Theme, error) {
	theme := &Theme{
		ID:           uuid.New(),
		Name:         name,
		Slug:         strings.ToLower(slug),
		TemplatePath: templatePath,
		Category:     CategoryElegant,
		IsPremium:    false,
		IsActive:     true,
		Settings: &ThemeSettings{
			HasMusic:     true,
			HasGallery:   true,
			HasCountdown: true,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := theme.Validate(); err != nil {
		return nil, err
	}

	return theme, nil
}

// Validate performs business rule validation
func (t *Theme) Validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return errors.New("theme name is required")
	}

	if strings.TrimSpace(t.Slug) == "" {
		return errors.New("theme slug is required")
	}

	// TemplatePath OR CustomHTML must be present
	hasTemplate := strings.TrimSpace(t.TemplatePath) != ""
	hasCustomHTML := t.CustomHTML != nil && strings.TrimSpace(*t.CustomHTML) != ""

	if !hasTemplate && !hasCustomHTML {
		return errors.New("either template path or custom html is required")
	}

	return nil
}

// TableName returns the table name
func (Theme) TableName() string {
	return "theme"
}

// BeforeCreate GORM hook
func (t *Theme) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate GORM hook
func (t *Theme) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}
