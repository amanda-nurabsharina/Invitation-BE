package api

import (
	"strconv"

	"invitation-api/internal/domain/user"
	userRepo "invitation-api/internal/repository/user"
	userService "invitation-api/internal/service/user"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UserHandler handles user management API requests
type UserHandler struct {
	userService userService.Service
}

// NewUserHandler creates a new user handler
func NewUserHandler() *UserHandler {
	userRepository := userRepo.NewRepository()
	userSvc := userService.NewService(userRepository)

	return &UserHandler{
		userService: userSvc,
	}
}

// CreateUserRequest represents create user request
type CreateUserRequest struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
	FullName string `json:"full_name"`
	Role     string `json:"role" validate:"required"` // "admin" or "user"
}

// UpdateUserRequest represents update user request
type UpdateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

// CreateRoleRequest represents create role request
type CreateRoleRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

// CreateMenuRequest represents create menu request
type CreateMenuRequest struct {
	Title    string     `json:"title" validate:"required"`
	Path     string     `json:"path" validate:"required"`
	Icon     string     `json:"icon"`
	Order    int        `json:"order"`
	ParentID *uuid.UUID `json:"parent_id"`
}

// UpdateRoleRequest
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// UpdateMenuRequest
type UpdateMenuRequest struct {
	Title string `json:"title"`
	Path  string `json:"path"`
	Icon  string `json:"icon"`
	Order int    `json:"order"`
}

// AssignMenuRequest represents assigning menu to role
type AssignMenuRequest struct {
	MenuID uuid.UUID `json:"menu_id" validate:"required"`
}

// --- User Handlers ---

// CreateUser creates a new user (admin only)
func (h *UserHandler) CreateUser(c *fiber.Ctx) error {
	var req CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Accept any role name from the request
	userRole := user.UserRole(req.Role)

	fullName := req.FullName
	if fullName == "" {
		fullName = req.Username
	}

	newUser, err := h.userService.RegisterUser(req.Username, req.Email, req.Password, fullName, userRole)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "User created successfully",
		"data":    newUser,
	})
}

// ListUsers lists all users
func (h *UserHandler) ListUsers(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	users, err := h.userService.ListUsers(limit, offset)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to list users",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  users,
	})
}

// GetUser retrieves a user by ID
func (h *UserHandler) GetUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid ID"})
	}

	user, err := h.userService.GetUserByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": true, "message": "User not found"})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  user,
	})
}

// UpdateUser updates a user
func (h *UserHandler) UpdateUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid ID"})
	}

	var req UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid body"})
	}

	updatedUser, err := h.userService.UpdateUser(id, req.Username, req.Email, req.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "User updated",
		"data":    updatedUser,
	})
}

// DeleteUser soft deletes a user
func (h *UserHandler) DeleteUser(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid ID"})
	}

	if err := h.userService.SoftDeleteUser(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to delete user"})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "User deleted",
	})
}

// --- Role Handlers ---

// CreateRole creates a new role
func (h *UserHandler) CreateRole(c *fiber.Ctx) error {
	var req CreateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid body"})
	}

	role, err := h.userService.CreateRole(req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Role created",
		"data":    role,
	})
}

func (h *UserHandler) UpdateRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid ID"})
	}

	var req UpdateRoleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid body"})
	}

	role, err := h.userService.UpdateRole(id, req.Name, req.Description)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"error": false, "message": "Role updated", "data": role})
}

func (h *UserHandler) ListRoles(c *fiber.Ctx) error {
	roles, err := h.userService.ListRoles()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to list roles"})
	}
	return c.JSON(fiber.Map{"error": false, "data": roles})
}

func (h *UserHandler) DeleteRole(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid ID"})
	}
	if err := h.userService.DeleteRole(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to delete role"})
	}
	return c.JSON(fiber.Map{"error": false, "message": "Role deleted"})
}

// --- Menu Handlers ---

// CreateMenu creates a new menu
func (h *UserHandler) CreateMenu(c *fiber.Ctx) error {
	var req CreateMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid body"})
	}

	menu, err := h.userService.CreateMenu(req.Title, req.Path, req.Icon, req.Order, req.ParentID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Menu created",
		"data":    menu,
	})
}

func (h *UserHandler) UpdateMenu(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid ID"})
	}

	var req UpdateMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid body"})
	}

	menu, err := h.userService.UpdateMenu(id, req.Title, req.Path, req.Icon, req.Order)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"error": false, "message": "Menu updated", "data": menu})
}

func (h *UserHandler) ListMenus(c *fiber.Ctx) error {
	menus, err := h.userService.ListMenus()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to list menus"})
	}
	return c.JSON(fiber.Map{"error": false, "data": menus})
}

func (h *UserHandler) DeleteMenu(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid ID"})
	}
	if err := h.userService.DeleteMenu(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": "Failed to delete menu"})
	}
	return c.JSON(fiber.Map{"error": false, "message": "Menu deleted"})
}

// --- Assignment Handlers ---

// AssignMenuToRole assigns a menu to a role
func (h *UserHandler) AssignMenuToRole(c *fiber.Ctx) error {
	roleIdStr := c.Params("role_id")
	roleId, err := uuid.Parse(roleIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid Role ID"})
	}

	var req AssignMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid body"})
	}

	if err := h.userService.AssignMenuToRole(roleId, req.MenuID); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"error": false, "message": "Menu assigned to role"})
}

// RemoveMenuFromRole
func (h *UserHandler) RemoveMenuFromRole(c *fiber.Ctx) error {
	roleIdStr := c.Params("role_id")
	roleId, err := uuid.Parse(roleIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid Role ID"})
	}

	menuIdStr := c.Params("menu_id")
	menuId, err := uuid.Parse(menuIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid Menu ID"})
	}

	if err := h.userService.RemoveMenuFromRole(roleId, menuId); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"error": false, "message": "Menu removed from role"})
}

// GetRoleMenus
func (h *UserHandler) GetRoleMenus(c *fiber.Ctx) error {
	roleIdStr := c.Params("role_id")
	roleId, err := uuid.Parse(roleIdStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": true, "message": "Invalid Role ID"})
	}

	menus, err := h.userService.GetMenusByRole(roleId)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"error": false, "data": menus})
}

// GetUserMenus returns menus for the current authenticated user's role
func (h *UserHandler) GetUserMenus(c *fiber.Ctx) error {
	userRoleName := c.Locals("user_role").(string)

	// Map JWT role name to DB role name
	// Search for matching role in DB
	// The JWT role name could be "Super Admin", "admin", "user", etc.
	searchNames := []string{userRoleName}
	// Backward compat: old "admin" JWT tokens should also match "Super Admin"
	if userRoleName == "admin" {
		searchNames = append(searchNames, "Super Admin")
	}

	roles, _ := h.userService.ListRoles()
	var roleID uuid.UUID
	found := false
	for _, searchName := range searchNames {
		for _, r := range roles {
			if r.Name == searchName {
				roleID = r.ID
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		// Fallback: admin gets everything
		if userRoleName == "admin" {
			menus, _ := h.userService.ListMenus()
			return c.JSON(fiber.Map{"error": false, "data": menus})
		}
		return c.JSON(fiber.Map{"error": false, "data": []user.Menu{}})
	}

	menus, err := h.userService.GetMenusByRole(roleID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": true, "message": err.Error()})
	}

	return c.JSON(fiber.Map{"error": false, "data": menus})
}
