package user

import (
	"fmt"
	"invitation-api/internal/database"
	userDomain "invitation-api/internal/domain/user"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for user data operations
type Repository interface {
	// User
	Create(user *userDomain.User) error
	GetByID(id uuid.UUID) (*userDomain.User, error)
	GetByEmail(email string) (*userDomain.User, error)
	GetByUsername(username string) (*userDomain.User, error)
	Update(user *userDomain.User) error
	SoftDelete(id uuid.UUID) error
	List(limit, offset int) ([]*userDomain.User, error)
	Count() (int64, error)
	CreateWithProfile(user *userDomain.User) error

	// Role Management
	CreateRole(role *userDomain.Role) error
	GetRoleByID(id uuid.UUID) (*userDomain.Role, error)
	GetRoleByName(name string) (*userDomain.Role, error)
	ListRoles() ([]*userDomain.Role, error)
	UpdateRole(role *userDomain.Role) error
	DeleteRole(id uuid.UUID) error

	// Menu Management
	CreateMenu(menu *userDomain.Menu) error
	GetMenuByID(id uuid.UUID) (*userDomain.Menu, error)
	ListMenus() ([]*userDomain.Menu, error)
	UpdateMenu(menu *userDomain.Menu) error
	DeleteMenu(id uuid.UUID) error

	// Role-Menu Assignment
	AssignMenuToRole(roleID, menuID uuid.UUID) error
	RemoveMenuFromRole(roleID, menuID uuid.UUID) error
	GetMenusByRole(roleID uuid.UUID) ([]*userDomain.Menu, error)
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new user repository instance
func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

// Create creates a new user in the database
func (r *repository) Create(user *userDomain.User) error {
	return r.db.Create(user).Error
}

// CreateWithProfile creates a new user with profile in the database
func (r *repository) CreateWithProfile(user *userDomain.User) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		var existingUser userDomain.User
		err := tx.Where("email = ? OR username = ?", user.Email, user.Username).First(&existingUser).Error
		if err == nil {
			return fmt.Errorf("user with email %s or username %s already exists", user.Email, user.Username)
		} else if err != gorm.ErrRecordNotFound {
			return err
		}

		// Store profile reference and remove from user to prevent auto-create
		profile := user.Profile
		user.Profile = nil

		// Create user without profile association
		if err := tx.Omit("Profile").Create(user).Error; err != nil {
			return err
		}

		// Create profile separately if exists
		if profile != nil {
			profile.UserID = user.ID
			if err := tx.Create(profile).Error; err != nil {
				return err
			}
			user.Profile = profile
		}

		return nil
	})
}

// GetByID retrieves a user by ID
func (r *repository) GetByID(id uuid.UUID) (*userDomain.User, error) {
	var user userDomain.User
	err := r.db.Preload("Profile").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *repository) GetByEmail(email string) (*userDomain.User, error) {
	var user userDomain.User
	err := r.db.Preload("Profile").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *repository) GetByUsername(username string) (*userDomain.User, error) {
	var user userDomain.User
	err := r.db.Preload("Profile").Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates an existing user
func (r *repository) Update(user *userDomain.User) error {
	return r.db.Save(user).Error
}

// SoftDelete soft deletes a user
func (r *repository) SoftDelete(id uuid.UUID) error {
	return r.db.Delete(&userDomain.User{}, id).Error
}

// List retrieves a list of users with pagination
func (r *repository) List(limit, offset int) ([]*userDomain.User, error) {
	var users []*userDomain.User
	err := r.db.Preload("Profile").Limit(limit).Offset(offset).Find(&users).Error
	return users, err
}

// Count returns the total count of users
func (r *repository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&userDomain.User{}).Count(&count).Error
	return count, err
}
