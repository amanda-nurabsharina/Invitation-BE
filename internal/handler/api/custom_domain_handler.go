package api

import (
	tenantRepo "invitation-api/internal/repository/tenant"
	customDomainService "invitation-api/internal/service/custom_domain"
	"invitation-api/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// CustomDomainHandler handles custom domain API requests
type CustomDomainHandler struct {
	customDomainService *customDomainService.Service
}

// NewCustomDomainHandler creates a new custom domain handler
func NewCustomDomainHandler(cfg *config.Config) *CustomDomainHandler {
	tenantRepository := tenantRepo.NewRepository()
	domainSvc := customDomainService.NewService(tenantRepository, cfg)

	return &CustomDomainHandler{
		customDomainService: domainSvc,
	}
}

// SetupDomainRequest represents custom domain setup request
type SetupDomainRequest struct {
	TenantID     string `json:"tenant_id" validate:"required"`
	CustomDomain string `json:"custom_domain" validate:"required"`
}

// Setup sets up a custom domain for a tenant
func (h *CustomDomainHandler) Setup(c *fiber.Ctx) error {
	var req SetupDomainRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	tenantID, err := uuid.Parse(req.TenantID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid tenant ID",
		})
	}

	if err := h.customDomainService.SetupCustomDomain(tenantID, req.CustomDomain); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Custom domain setup initiated. DNS verification may take up to 24 hours.",
		"data": fiber.Map{
			"custom_domain": req.CustomDomain,
			"instructions": []string{
				"Add a CNAME record pointing your domain to invitation.yourdomain.com",
				"Wait for DNS propagation (up to 24 hours)",
				"SSL certificate will be automatically issued once DNS is verified",
			},
		},
	})
}

// Remove removes a custom domain from a tenant
func (h *CustomDomainHandler) Remove(c *fiber.Ctx) error {
	tenantIDStr := c.Params("tenant_id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid tenant ID",
		})
	}

	if err := h.customDomainService.RemoveCustomDomain(tenantID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Custom domain removed successfully",
	})
}

// VerifyDNS verifies DNS configuration for a domain
func (h *CustomDomainHandler) VerifyDNS(c *fiber.Ctx) error {
	domain := c.Query("domain")
	if domain == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Domain is required",
		})
	}

	result, err := h.customDomainService.VerifyDomainDNS(domain)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to verify DNS",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  result,
	})
}
