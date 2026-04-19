package api

import (
	"strconv"

	themeRepo "invitation-api/internal/repository/theme"
	themeService "invitation-api/internal/service/theme"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// ThemeHandler handles theme API requests
type ThemeHandler struct {
	themeService themeService.Service
}

// NewThemeHandler creates a new theme handler
func NewThemeHandler() *ThemeHandler {
	themeRepository := themeRepo.NewRepository()
	themeSvc := themeService.NewService(themeRepository)

	return &ThemeHandler{
		themeService: themeSvc,
	}
}

// CreateThemeRequest represents create theme request
type CreateThemeRequest struct {
	Name         string `json:"name" validate:"required"`
	Slug         string `json:"slug" validate:"required"`
	TemplatePath string `json:"template_path"`
	CustomHTML   string `json:"custom_html"`
	CustomCSS    string `json:"custom_css"`
	Category     string `json:"category"`
	IsPremium    bool   `json:"is_premium"`
}

// Create creates a new theme (admin only)
func (h *ThemeHandler) Create(c *fiber.Ctx) error {
	var req CreateThemeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	category := req.Category
	if category == "" {
		category = "elegant"
	}

	templatePath := req.TemplatePath
	if templatePath == "" {
		templatePath = "db_custom"
	}

	theme, err := h.themeService.Create(req.Name, req.Slug, templatePath, category, req.IsPremium, req.CustomHTML, req.CustomCSS)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Theme created successfully",
		"data":    theme,
	})
}

// GetByID retrieves a theme by ID
func (h *ThemeHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid theme ID",
		})
	}

	theme, err := h.themeService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Theme not found",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  theme,
	})
}

// List retrieves all themes
func (h *ThemeHandler) List(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	activeOnly := c.Query("active_only", "true") == "true"
	category := c.Query("category", "")

	if category != "" {
		themeList, err := h.themeService.ListByCategory(category)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Failed to retrieve themes",
			})
		}
		return c.JSON(fiber.Map{
			"error": false,
			"data":  themeList,
		})
	}

	themeList, total, err := h.themeService.List(limit, offset, activeOnly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve themes",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  themeList,
		"meta": fiber.Map{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// Update updates a theme
func (h *ThemeHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid theme ID",
		})
	}

	var req themeService.UpdateThemeRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	theme, err := h.themeService.Update(id, &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Theme updated successfully",
		"data":    theme,
	})
}

// Delete deletes a theme
func (h *ThemeHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid theme ID",
		})
	}

	if err := h.themeService.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete theme",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Theme deleted successfully",
	})
}

// Activate activates a theme
func (h *ThemeHandler) Activate(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid theme ID",
		})
	}

	if err := h.themeService.Activate(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to activate theme",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Theme activated successfully",
	})
}

// Deactivate deactivates a theme
func (h *ThemeHandler) Deactivate(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid theme ID",
		})
	}

	if err := h.themeService.Deactivate(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to deactivate theme",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Theme deactivated successfully",
	})
}
