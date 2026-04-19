package middleware

import (
	authRepo "invitation-api/internal/repository/auth"
	userRepository "invitation-api/internal/repository/user"
	authService "invitation-api/internal/service/auth"
	"invitation-api/pkg/config"
	"invitation-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// AuthMiddleware creates authentication middleware
func AuthMiddleware(cfg *config.Config) fiber.Handler {
	// Initialize repositories and services
	authRepository := authRepo.NewRepository()
	userRepo := userRepository.NewRepository()
	authSvc := authService.NewService(authRepository, userRepo, cfg)

	return func(c *fiber.Ctx) error {
		// Extract token from Authorization header or Query param
		authHeader := c.Get("Authorization")
		token, err := utils.ExtractTokenFromHeader(authHeader)
		if err != nil {
			// Fallback: Check for "token" query parameter
			token = c.Query("token")
			if token == "" {
				// DEBUG: Print path where auth failed
				// fmt.Printf("Auth Failed for Path: %s, Headers: %v\n", c.Path(), c.GetReqHeaders())
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error":   true,
					"message": "Authorization header is required",
				})
			}
		}

		// Validate token
		claims, err := authSvc.ValidateToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid or expired token",
			})
		}

		// Store user information in context
		c.Locals("user_id", claims.UserID.String())
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)
		c.Locals("user_full_name", claims.FullName)

		return c.Next()
	}
}

// RoleMiddleware creates role-based access control middleware
func RoleMiddleware(allowedRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("user_role").(string)
		if !ok {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error":   true,
				"message": "Insufficient permissions",
			})
		}

		// Check if user role is in allowed roles
		for _, role := range allowedRoles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error":   true,
			"message": "Insufficient permissions",
		})
	}
}

// AdminOnlyMiddleware creates middleware that only allows admin access
func AdminOnlyMiddleware() fiber.Handler {
	return RoleMiddleware("admin", "Super Admin")
}

// GetUserIDFromContext extracts the user ID from the fiber context
func GetUserIDFromContext(c *fiber.Ctx) (uuid.UUID, error) {
	userIDInterface := c.Locals("user_id")
	if userIDInterface == nil {
		return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "User authentication required")
	}

	// Try to get as UUID first
	if userID, ok := userIDInterface.(uuid.UUID); ok {
		return userID, nil
	}

	// Try to get as string and parse
	if userIDStr, ok := userIDInterface.(string); ok {
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid user authentication")
		}
		return userID, nil
	}

	return uuid.Nil, fiber.NewError(fiber.StatusUnauthorized, "Invalid user authentication")
}
