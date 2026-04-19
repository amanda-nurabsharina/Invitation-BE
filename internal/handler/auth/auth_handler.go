package auth

import (
	"invitation-api/internal/domain/auth"
	authRepo "invitation-api/internal/repository/auth"
	userRepo "invitation-api/internal/repository/user"
	authService "invitation-api/internal/service/auth"
	"invitation-api/pkg/config"
	"invitation-api/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Handler handles authentication HTTP requests
type Handler struct {
	authService authService.Service
}

// NewHandler creates a new auth handler instance
func NewHandler(cfg *config.Config) *Handler {
	// Initialize repositories
	authRepository := authRepo.NewRepository()
	userRepository := userRepo.NewRepository()

	// Initialize service
	authSvc := authService.NewService(authRepository, userRepository, cfg)

	return &Handler{
		authService: authSvc,
	}
}

// Register handles user registration
func (h *Handler) Register(c *fiber.Ctx) error {
	var req auth.RegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Register user
	tokenResponse, err := h.authService.Register(&req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "User registered successfully",
		"data":    tokenResponse,
	})
}

// Login handles user login
func (h *Handler) Login(c *fiber.Ctx) error {
	var req auth.LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Login user
	tokenResponse, err := h.authService.Login(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Login successful",
		"data":    tokenResponse,
	})
}

// RefreshToken handles token refresh
func (h *Handler) RefreshToken(c *fiber.Ctx) error {
	var req auth.RefreshTokenRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Refresh token
	tokenResponse, err := h.authService.RefreshToken(&req)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Token refreshed successfully",
		"data":    tokenResponse,
	})
}

// Logout handles user logout
func (h *Handler) Logout(c *fiber.Ctx) error {
	// Extract token from Authorization header
	authHeader := c.Get("Authorization")
	token, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Authorization header is required",
		})
	}

	// Logout user
	if err := h.authService.Logout(token); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Logout successful",
	})
}

// GetProfile handles getting user profile (protected route example)
func (h *Handler) GetProfile(c *fiber.Ctx) error {
	// Extract token from Authorization header
	authHeader := c.Get("Authorization")
	token, err := utils.ExtractTokenFromHeader(authHeader)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Authorization header is required",
		})
	}

	// Validate token
	claims, err := h.authService.ValidateToken(token)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid or expired token",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"error":   false,
		"message": "Profile retrieved successfully",
		"data": fiber.Map{
			"user_id":   claims.UserID,
			"email":     claims.Email,
			"role":      claims.Role,
			"full_name": claims.FullName,
		},
	})
}
