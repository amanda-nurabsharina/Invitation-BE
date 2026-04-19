package middleware

import (
	"strings"

	tenantRepo "invitation-api/internal/repository/tenant"
	"invitation-api/pkg/config"

	"github.com/gofiber/fiber/v2"
)

// TenantMiddleware resolves tenant from subdomain
func TenantMiddleware(cfg *config.Config) fiber.Handler {
	tenantRepository := tenantRepo.NewRepository()

	return func(c *fiber.Ctx) error {
		// Get host from request
		host := c.Hostname()

		// Extract subdomain from host
		subdomain := extractSubdomain(host, cfg.App.BaseDomain)

		// Skip tenant resolution for main domain, api, admin, cms
		if subdomain == "" || subdomain == "www" || subdomain == "api" || subdomain == "admin" || subdomain == "cms" {
			return c.Next()
		}

		// Find tenant by subdomain
		tenant, err := tenantRepository.GetBySubdomain(subdomain)
		if err != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "Invitation not found",
			})
		}

		// Check if tenant is accessible
		if !tenant.IsAccessible() {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "This invitation is not accessible",
			})
		}

		// Store tenant info in context
		c.Locals("tenant_id", tenant.ID.String())
		c.Locals("tenant_subdomain", tenant.Subdomain)
		c.Locals("tenant", tenant)

		return c.Next()
	}
}

// TenantRequired middleware ensures tenant is present in context
func TenantRequired() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Locals("tenant_id")
		if tenantID == nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Tenant context required",
			})
		}
		return c.Next()
	}
}

// extractSubdomain extracts subdomain from full host
// e.g., "john-jane.wedding.com" -> "john-jane"
func extractSubdomain(host, baseDomain string) string {
	// Remove port if present
	if colonIdx := strings.Index(host, ":"); colonIdx != -1 {
		host = host[:colonIdx]
	}

	// If baseDomain is not configured, try to extract first part
	if baseDomain == "" {
		parts := strings.Split(host, ".")
		if len(parts) >= 3 {
			return parts[0]
		}
		return ""
	}

	// Remove base domain to get subdomain
	if strings.HasSuffix(host, "."+baseDomain) {
		subdomain := strings.TrimSuffix(host, "."+baseDomain)
		// Handle nested subdomains
		parts := strings.Split(subdomain, ".")
		return parts[len(parts)-1]
	}

	return ""
}

// GetTenantIDFromContext extracts tenant ID from context
func GetTenantIDFromContext(c *fiber.Ctx) string {
	tenantID, ok := c.Locals("tenant_id").(string)
	if !ok {
		return ""
	}
	return tenantID
}
