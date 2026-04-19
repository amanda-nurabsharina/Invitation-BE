package api

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	invitationRepo "invitation-api/internal/repository/invitation"
	tenantRepo "invitation-api/internal/repository/tenant"
	invitationService "invitation-api/internal/service/invitation"
	tenantService "invitation-api/internal/service/tenant"
	"invitation-api/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// TenantHandler handles tenant API requests
type TenantHandler struct {
	tenantService     tenantService.Service
	invitationService invitationService.Service
}

// NewTenantHandler creates a new tenant handler
func NewTenantHandler(cfg *config.Config) *TenantHandler {
	tenantRepository := tenantRepo.NewRepository()
	tenantSvc := tenantService.NewService(tenantRepository, cfg)

	invitationRepository := invitationRepo.NewRepository()
	invitationSvc := invitationService.NewService(invitationRepository, tenantRepository)

	return &TenantHandler{
		tenantService:     tenantSvc,
		invitationService: invitationSvc,
	}
}

// CreateTenantRequest represents create tenant request
type CreateTenantRequest struct {
	Name      string `json:"name" validate:"required"`
	Subdomain string `json:"subdomain"` // optional, will be generated if empty
}

// toSlug converts a string to a slug
func toSlug(s string) string {
	s = strings.ToLower(s)
	reg := regexp.MustCompile(`[^a-z0-9]+`)
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	if len(s) > 30 {
		s = s[:30]
	}
	// Append random string to ensure uniqueness if too short or likely common
	return fmt.Sprintf("%s-%s", s, uuid.New().String()[:6])
}

// UpdateTenantRequest represents update tenant request
type UpdateTenantRequest struct {
	Name         string  `json:"name"`
	CustomDomain *string `json:"custom_domain"`
}

// Create creates a new tenant
func (h *TenantHandler) Create(c *fiber.Ctx) error {
	var req CreateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Auto-generate subdomain if empty
	if req.Subdomain == "" {
		req.Subdomain = toSlug(req.Name)
	}

	// Get user ID from context (set by auth middleware)
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "User not authenticated",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid user ID",
		})
	}

	tenant, err := h.tenantService.Create(userID, req.Name, req.Subdomain)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	// Create empty invitation
	h.invitationService.Create(tenant.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Tenant created successfully",
		"data":    tenant,
	})
}

// GetByID retrieves a tenant by ID
func (h *TenantHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid tenant ID",
		})
	}

	tenant, err := h.tenantService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Tenant not found",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  tenant,
	})
}

// GetMyTenants retrieves all tenants for the current user
func (h *TenantHandler) GetMyTenants(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "User not authenticated",
		})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid user ID",
		})
	}

	tenants, err := h.tenantService.GetByUserID(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve tenants",
		})
	}

	// Enrich with invitation status
	var result []map[string]interface{}
	for _, t := range tenants {
		status := "Draft"
		var expiresAt *time.Time
		isPublished := false

		inv, err := h.invitationService.GetByTenantID(t.ID)
		if err == nil && inv != nil {
			expiresAt = inv.ExpiresAt
			isPublished = inv.IsPublished

			if inv.IsExpired() {
				status = "Expired"
			} else if inv.IsPublished {
				status = "Published"
			}
		}

		// Convert to map to avoid creating new struct type if possible, or use standard struct
		tData := map[string]interface{}{
			"id":                t.ID,
			"name":              t.Name,
			"subdomain":         t.Subdomain,
			"is_active":         t.IsActive,
			"created_at":        t.CreatedAt,
			"invitation_status": status,
			"expires_at":        expiresAt,
			"is_published":      isPublished,
		}
		result = append(result, tData)
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  result,
	})
}

// List retrieves all tenants (admin only)
func (h *TenantHandler) List(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	tenants, total, err := h.tenantService.List(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve tenants",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  tenants,
		"meta": fiber.Map{
			"total":  total,
			"limit":  limit,
			"offset": offset,
		},
	})
}

// Update updates a tenant
func (h *TenantHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid tenant ID",
		})
	}

	var req UpdateTenantRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	tenant, err := h.tenantService.Update(id, req.Name, req.CustomDomain)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Tenant updated successfully",
		"data":    tenant,
	})
}

// Delete deletes a tenant
func (h *TenantHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid tenant ID",
		})
	}

	if err := h.tenantService.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete tenant",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Tenant deleted successfully",
	})
}

// CheckSubdomain checks if a subdomain is available
func (h *TenantHandler) CheckSubdomain(c *fiber.Ctx) error {
	subdomain := c.Query("subdomain")
	if subdomain == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Subdomain is required",
		})
	}

	available, err := h.tenantService.CheckSubdomainAvailability(subdomain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to check subdomain",
		})
	}

	return c.JSON(fiber.Map{
		"error":     false,
		"available": available,
	})
}

// Activate activates a tenant
func (h *TenantHandler) Activate(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid tenant ID",
		})
	}

	if err := h.tenantService.Activate(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to activate tenant",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Tenant activated successfully",
	})
}

// Deactivate deactivates a tenant
func (h *TenantHandler) Deactivate(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid tenant ID",
		})
	}

	if err := h.tenantService.Deactivate(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to deactivate tenant",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Tenant deactivated successfully",
	})
}
