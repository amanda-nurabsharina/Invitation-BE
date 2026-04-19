package seed_themes

import (
	"log"

	themeDomain "invitation-api/internal/domain/theme"
	themeRepo "invitation-api/internal/repository/theme"
)

// CreateInitialThemes seeds default themes
func CreateInitialThemes() {
	repo := themeRepo.NewRepository()

	themes := []struct {
		Name         string
		Slug         string
		TemplatePath string
		Category     string
		IsPremium    bool
		Colors       *themeDomain.ThemeColors
		Description  string
	}{
		{
			Name:         "Elegant Gold",
			Slug:         "elegant",
			TemplatePath: "elegant",
			Category:     "elegant",
			IsPremium:    false,
			Description:  "Tema elegan dengan nuansa emas klasik yang mewah",
			Colors: &themeDomain.ThemeColors{
				Primary:    "#9b7b5c",
				Secondary:  "#f5f0eb",
				Accent:     "#d4af37",
				Background: "#f5f0eb",
				Text:       "#3d3d3d",
				TextMuted:  "#6b6b6b",
			},
		},
		{
			Name:         "Modern V2 Split",
			Slug:         "modern_v2",
			TemplatePath: "modern_v2",
			Category:     "modern",
			IsPremium:    false,
			Description:  "Tema modern dengan layout split screen dan sidebar navigation",
			Colors: &themeDomain.ThemeColors{
				Primary:    "#2c3e50",
				Secondary:  "#ffffff",
				Accent:     "#34495e",
				Background: "#f8f9fa",
				Text:       "#2c3e50",
				TextMuted:  "#7f8c8d",
			},
		},
		{
			Name:         "Rustic Garden",
			Slug:         "rustic",
			TemplatePath: "rustic",
			Category:     "rustic",
			IsPremium:    false,
			Description:  "Tema rustic dengan sentuhan natural dan hangat",
			Colors: &themeDomain.ThemeColors{
				Primary:    "#8b7355",
				Secondary:  "#f5f1eb",
				Accent:     "#b8860b",
				Background: "#faf8f5",
				Text:       "#4a4a4a",
				TextMuted:  "#7a7a7a",
			},
		},
		{
			Name:         "Modern Minimalist",
			Slug:         "modern",
			TemplatePath: "modern",
			Category:     "modern",
			IsPremium:    true,
			Description:  "Tema modern dengan desain bersih dan minimalis",
			Colors: &themeDomain.ThemeColors{
				Primary:    "#2c3e50",
				Secondary:  "#ecf0f1",
				Accent:     "#e74c3c",
				Background: "#ffffff",
				Text:       "#2c3e50",
				TextMuted:  "#7f8c8d",
			},
		},
		{
			Name:         "Floral Romance",
			Slug:         "floral",
			TemplatePath: "floral",
			Category:     "floral",
			IsPremium:    true,
			Description:  "Tema romantis dengan ornamen bunga yang indah",
			Colors: &themeDomain.ThemeColors{
				Primary:    "#d4a5a5",
				Secondary:  "#fff5f5",
				Accent:     "#e8b4b4",
				Background: "#fffafa",
				Text:       "#5a4a4a",
				TextMuted:  "#8a7a7a",
			},
		},
		{
			Name:         "Javanese Traditional",
			Slug:         "javanese",
			TemplatePath: "javanese",
			Category:     "traditional",
			IsPremium:    true,
			Description:  "Tema tradisional Jawa dengan batik dan ornamen klasik",
			Colors: &themeDomain.ThemeColors{
				Primary:    "#8B4513",
				Secondary:  "#FFF8DC",
				Accent:     "#DAA520",
				Background: "#FFFAF0",
				Text:       "#4a3728",
				TextMuted:  "#8b7355",
			},
		},
		{
			Name:         "Luxury Gold",
			Slug:         "luxury_gold",
			TemplatePath: "luxury",
			Category:     "luxury",
			IsPremium:    true,
			Description:  "Tema mewah dengan sentuhan emas, animasi premium, slideshow, countdown, love story, dan fitur amplop digital",
			Colors: &themeDomain.ThemeColors{
				Primary:    "#C5A880",
				Secondary:  "#2A2A2A",
				Accent:     "#d2d0cc",
				Background: "#FAF8F5",
				Text:       "#2A2A2A",
				TextMuted:  "#6b6b6b",
			},
		},
	}

	for _, t := range themes {
		// Check if theme already exists
		existing, _ := repo.GetBySlug(t.Slug)
		if existing != nil {
			// Fix TemplatePath if different
			if existing.TemplatePath != t.TemplatePath {
				existing.TemplatePath = t.TemplatePath
				if err := repo.Update(existing); err != nil {
					log.Printf("Failed to update theme %s template path: %v", t.Slug, err)
				} else {
					log.Printf("✅ Theme %s template path updated to: %s", t.Slug, t.TemplatePath)
				}
			}
			log.Printf("Theme %s already exists, skipping", t.Slug)
			continue
		}

		theme, err := themeDomain.NewTheme(t.Name, t.Slug, t.TemplatePath)
		if err != nil {
			log.Printf("Failed to create theme %s: %v", t.Name, err)
			continue
		}

		theme.Category = t.Category
		theme.IsPremium = t.IsPremium
		theme.Colors = t.Colors
		desc := t.Description
		theme.Description = &desc
		theme.Settings = &themeDomain.ThemeSettings{
			HasMusic:       true,
			HasGallery:     true,
			HasCountdown:   true,
			AnimationStyle: "fade",
			GalleryStyle:   "grid",
		}

		if err := repo.Create(theme); err != nil {
			log.Printf("Failed to save theme %s: %v", t.Name, err)
			continue
		}

		log.Printf("✅ Theme %s created successfully", t.Name)
	}
}
