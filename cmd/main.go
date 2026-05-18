package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"invitation-api/internal/database"
	apiHandler "invitation-api/internal/handler/api"
	authHandler "invitation-api/internal/handler/auth"
	htmxHandler "invitation-api/internal/handler/htmx"
	"invitation-api/internal/middleware"
	seedInitialUserSvc "invitation-api/internal/service/seed_initial_user"
	seedThemesSvc "invitation-api/internal/service/seed_themes"
	"invitation-api/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Println("Configuration loaded:", cfg.App.Environment)
	log.Printf("STARTUP CONFIG: DB Name = %s", cfg.Database.DBName)

	// Connect to database
	if err := database.Connect(cfg); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Run migrations
	if err := database.AutoMigrate(); err != nil {
		log.Fatalf("Failed to run database migrations: %v", err)
	}

	// Seed initial data in development mode
	if cfg.IsDevelopment() {
		seedInitialUserSvc.CreateInitialUser()
		seedInitialUserSvc.SeedDefaultMenus()
		seedInitialUserSvc.SeedSuperAdminRole()
		seedThemesSvc.CreateInitialThemes()
		seedThemesSvc.CreateInitialCategories()
	}

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
		ErrorHandler: customErrorHandler,
		BodyLimit:    20 * 1024 * 1024, // 20MB for file uploads
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS, PATCH",
	}))

	// Serve static files
	log.Println("🗂️  Serving static files from /uploads")
	app.Static("/uploads", "./uploads")
	app.Static("/assets", "./public/assets")

	log.Println("\n\n\n=======================================================")
	log.Println(" INVITATION API - SERVER STARTING ")
	log.Println("=======================================================")

	// Setup routes
	setupRoutes(app, cfg)

	// Start server
	go func() {
		addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
		log.Printf("Server starting on %s", addr)
		if err := app.Listen(addr); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := app.ShutdownWithContext(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

// setupRoutes configures all API routes
func setupRoutes(app *fiber.App, cfg *config.Config) {
	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Invitation API is running",
			"version": "1.0.0",
		})
	})

	// Initialize handlers
	authH := authHandler.NewHandler(cfg)
	tenantH := apiHandler.NewTenantHandler(cfg)
	invitationH := apiHandler.NewInvitationHandler(cfg)
	themeH := apiHandler.NewThemeHandler()
	themeCategoryH := apiHandler.NewCategoryHandler()
	uploadH := apiHandler.NewUploadHandler("./uploads", cfg.App.BaseURL)
	htmxH := htmxHandler.NewHandler(cfg)

	// ==========================================
	// API v1 routes (JSON for CMS)
	// ==========================================
	api := app.Group("/api/v1")

	// Public auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authH.Register)
	auth.Post("/login", authH.Login)
	auth.Post("/refresh-token", authH.RefreshToken)

	// Public theme routes
	api.Get("/themes", themeH.List)
	api.Get("/themes/:id", themeH.GetByID)
	api.Get("/categories", themeCategoryH.GetAll)                      // Public read for dropdown
	api.Get("/downloads/guest-template", invitationH.GetGuestTemplate) // Public template download

	// Check subdomain availability (public)
	api.Get("/check-subdomain", tenantH.CheckSubdomain)

	// Protected routes (requires authentication)
	protected := api.Group("/", middleware.AuthMiddleware(cfg))
	protected.Post("/auth/logout", authH.Logout)
	protected.Get("/auth/profile", authH.GetProfile)

	// Tenant routes (for logged in users)
	protected.Get("/my-tenants", tenantH.GetMyTenants)
	protected.Post("/tenants", tenantH.Create)
	protected.Get("/tenants/:id", tenantH.GetByID)
	protected.Put("/tenants/:id", tenantH.Update)
	protected.Delete("/tenants/:id", tenantH.Delete)

	// Invitation routes
	protected.Post("/invitations", invitationH.Create)
	protected.Get("/invitations/:id", invitationH.GetByID)
	protected.Get("/tenants/:tenant_id/invitation", invitationH.GetByTenantID)
	protected.Put("/invitations/:id", invitationH.Update)
	protected.Delete("/invitations/:id", invitationH.Delete)
	protected.Post("/invitations/:id/publish", invitationH.Publish)
	protected.Post("/invitations/:id/unpublish", invitationH.Unpublish)
	protected.Get("/invitations/:id/dashboard", invitationH.GetDashboard)

	// Gallery routes
	protected.Get("/invitations/:id/gallery", invitationH.GetGallery)
	protected.Post("/invitations/:id/gallery", invitationH.AddGalleryItem)
	protected.Delete("/gallery/:gallery_id", invitationH.DeleteGalleryItem)
	protected.Put("/gallery/reorder", invitationH.ReorderGallery)

	// RSVP routes
	protected.Get("/invitations/:id/rsvp", invitationH.GetRSVPs)
	protected.Get("/invitations/:id/rsvp/export", invitationH.ExportRSVPs)

	// Guest List routes
	protected.Get("/invitations/:id/guests", invitationH.GetGuests)
	protected.Post("/invitations/:id/guests", invitationH.AddGuest)
	protected.Delete("/invitations/:id/guests/:guest_id", invitationH.DeleteGuest)
	protected.Post("/invitations/:id/guests/import", invitationH.ImportGuests)
	protected.Get("/invitations/:id/guests/export", invitationH.ExportGuests)

	// Gift account routes
	protected.Get("/invitations/:id/gift-accounts", invitationH.GetGiftAccounts)
	protected.Post("/invitations/:id/gift-accounts", invitationH.AddGiftAccount)
	protected.Delete("/gift-accounts/:gift_id", invitationH.DeleteGiftAccount)

	// Guest message routes
	protected.Get("/invitations/:id/messages", invitationH.GetGuestMessages)
	protected.Put("/messages/:message_id/approve", invitationH.ApproveMessage)
	protected.Delete("/messages/:message_id", invitationH.DeleteMessage)

	// Upload routes
	protected.Post("/upload/image", uploadH.UploadImage)
	protected.Post("/upload/images", uploadH.UploadMultipleImages)
	protected.Post("/upload/music", uploadH.UploadMusic)
	protected.Delete("/upload", uploadH.DeleteFile)

	// Admin only routes
	admin := protected.Group("/admin", middleware.AdminOnlyMiddleware())
	admin.Get("/tenants", tenantH.List)
	admin.Put("/tenants/:id/activate", tenantH.Activate)
	admin.Put("/tenants/:id/deactivate", tenantH.Deactivate)

	// Theme Management
	admin.Post("/themes", themeH.Create)
	admin.Put("/themes/:id", themeH.Update)
	admin.Delete("/themes/:id", themeH.Delete)
	admin.Put("/themes/:id/activate", themeH.Activate)
	admin.Put("/themes/:id/deactivate", themeH.Deactivate)

	// Master Data: Categories
	admin.Post("/categories", themeCategoryH.Create)
	admin.Put("/categories/:id", themeCategoryH.Update)
	admin.Delete("/categories/:id", themeCategoryH.Delete)
	admin.Get("/categories/:id", themeCategoryH.GetByID) // Admin detail view

	// User Management (Admin)
	userH := apiHandler.NewUserHandler()

	// Roles
	admin.Post("/roles", userH.CreateRole)
	admin.Get("/roles", userH.ListRoles)
	admin.Delete("/roles/:id", userH.DeleteRole)
	admin.Put("/roles/:id", userH.UpdateRole)

	// Menus
	admin.Post("/menus", userH.CreateMenu)
	admin.Get("/menus", userH.ListMenus)
	admin.Delete("/menus/:id", userH.DeleteMenu)
	admin.Put("/menus/:id", userH.UpdateMenu)

	// Role-Menu Assignment
	admin.Post("/roles/:role_id/menus", userH.AssignMenuToRole)
	admin.Delete("/roles/:role_id/menus/:menu_id", userH.RemoveMenuFromRole)
	admin.Get("/roles/:role_id/menus", userH.GetRoleMenus)

	// Users
	admin.Post("/users", userH.CreateUser)
	admin.Get("/users", userH.ListUsers)
	admin.Get("/users/:id", userH.GetUser)
	admin.Put("/users/:id", userH.UpdateUser)
	admin.Delete("/users/:id", userH.DeleteUser)

	// Current User Permissions (Protected, not just admin)
	protected.Get("/my-menus", userH.GetUserMenus)

	// ==========================================
	// Client site routes (htmx - HTML fragments)
	// ==========================================
	// Main invitation page
	app.Get("/:subdomain", htmxH.RenderInvitation)

	// htmx fragment endpoints
	app.Get("/:subdomain/gallery", htmxH.RenderGallery)
	app.Get("/:subdomain/messages", htmxH.RenderMessages)
	app.Get("/:subdomain/gift", htmxH.RenderGift)
	app.Post("/:subdomain/rsvp", htmxH.SubmitRSVP)
}

// customErrorHandler handles errors globally
func customErrorHandler(c *fiber.Ctx, err error) error {
	// Default 500 statuscode
	code := fiber.StatusInternalServerError

	// Retrieve the custom statuscode if it's an *fiber.Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Send custom error page
	return c.Status(code).JSON(fiber.Map{
		"error":   true,
		"message": err.Error(),
	})
}
