package user

import (
	userDomain "invitation-api/internal/domain/user"

	"github.com/google/uuid"
)

// Role Management Methods

// CreateRole creates a new role
func (r *repository) CreateRole(role *userDomain.Role) error {
	return r.db.Create(role).Error
}

// GetRoleByID retrieves a role by ID
func (r *repository) GetRoleByID(id uuid.UUID) (*userDomain.Role, error) {
	var role userDomain.Role
	err := r.db.Preload("Menus").First(&role, id).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// GetRoleByName retrieves a role by name
func (r *repository) GetRoleByName(name string) (*userDomain.Role, error) {
	var role userDomain.Role
	err := r.db.Preload("Menus").Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// ListRoles retrieves all roles
func (r *repository) ListRoles() ([]*userDomain.Role, error) {
	var roles []*userDomain.Role
	err := r.db.Preload("Menus").Find(&roles).Error
	return roles, err
}

// UpdateRole updates a role
func (r *repository) UpdateRole(role *userDomain.Role) error {
	return r.db.Save(role).Error
}

// DeleteRole deletes a role
func (r *repository) DeleteRole(id uuid.UUID) error {
	return r.db.Delete(&userDomain.Role{}, id).Error
}

// Menu Management Methods

// CreateMenu creates a new menu
func (r *repository) CreateMenu(menu *userDomain.Menu) error {
	return r.db.Create(menu).Error
}

// GetMenuByID retrieves a menu by ID
func (r *repository) GetMenuByID(id uuid.UUID) (*userDomain.Menu, error) {
	var menu userDomain.Menu
	// Preload children for hierarchical view
	err := r.db.Preload("Children").First(&menu, id).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

// ListMenus retrieves all menus
func (r *repository) ListMenus() ([]*userDomain.Menu, error) {
	var menus []*userDomain.Menu
	// Get only root manus (parent_id users list) and preload children
	err := r.db.Where("parent_id IS NULL").Preload("Children").Order("\"order\" asc").Find(&menus).Error
	return menus, err
}

// UpdateMenu updates a menu
func (r *repository) UpdateMenu(menu *userDomain.Menu) error {
	return r.db.Save(menu).Error
}

// DeleteMenu deletes a menu
func (r *repository) DeleteMenu(id uuid.UUID) error {
	return r.db.Delete(&userDomain.Menu{}, id).Error
}

// Role-Menu Assignment Methods

// AssignMenuToRole assigns a menu to a role
func (r *repository) AssignMenuToRole(roleID, menuID uuid.UUID) error {
	// Use GORM association mode
	role := &userDomain.Role{ID: roleID}
	menu := &userDomain.Menu{ID: menuID}
	return r.db.Model(role).Association("Menus").Append(menu)
}

// RemoveMenuFromRole removes a menu from a role
func (r *repository) RemoveMenuFromRole(roleID, menuID uuid.UUID) error {
	role := &userDomain.Role{ID: roleID}
	menu := &userDomain.Menu{ID: menuID}
	return r.db.Model(role).Association("Menus").Delete(menu)
}

// GetMenusByRole retrieves menus assigned to a role
func (r *repository) GetMenusByRole(roleID uuid.UUID) ([]*userDomain.Menu, error) {
	var role userDomain.Role
	err := r.db.Preload("Menus", "parent_id IS NULL").Preload("Menus.Children").First(&role, roleID).Error
	if err != nil {
		return nil, err
	}

	// Filter menus to ensure hierarchy is respected if needed,
	// but Preload handles basic structure.
	// Note: GORM M2M preload with nested strict filtering can be complex.
	// For now return flattened or partially hierarchical list attached to role.
	return role.Menus, nil
}
