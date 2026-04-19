package user

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents the core user entity for authentication
type User struct {
	ID        uuid.UUID      `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	Username  string         `gorm:"column:username;uniqueIndex" json:"username"`
	Password  string         `gorm:"column:password" json:"-"`
	Email     string         `gorm:"column:email;uniqueIndex" json:"email"`
	Role      string         `gorm:"column:role" json:"role"`
	IsActive  bool           `gorm:"column:is_active;default:true" json:"is_active"`
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	// Relationships
	Profile *UserProfile `gorm:"foreignKey:UserID" json:"profile,omitempty"`
}

// UserRole represents the role of a user in the system
type UserRole string

const (
	RoleAdmin UserRole = "admin"
	RoleUser  UserRole = "user"
)

// NewUser creates a new user instance with validation
func NewUser(username, password, email string, role UserRole) (*User, error) {
	user := &User{
		ID:        uuid.New(),
		Username:  username,
		Password:  password,
		Email:     email,
		Role:      string(role),
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := user.Validate(); err != nil {
		return nil, err
	}

	return user, nil
}

// Validate performs business rule validation
func (u *User) Validate() error {
	if strings.TrimSpace(u.Username) == "" {
		return errors.New("username is required")
	}

	if strings.TrimSpace(u.Password) == "" {
		return errors.New("password is required")
	}

	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if u.Email != "" && !isValidEmail(u.Email) {
		return errors.New("invalid email format")
	}

	if !isValidRole(u.Role) {
		return errors.New("invalid user role")
	}

	return nil
}

// IsAdmin checks if the user has admin role
func (u *User) IsAdmin() bool {
	return u.Role == string(RoleAdmin)
}

// IsDeleted checks if the user is soft deleted
func (u *User) IsDeleted() bool {
	return !u.DeletedAt.Time.IsZero()
}

// IsActiveAndNotDeleted checks if the user is both active and not deleted
func (u *User) IsActiveAndNotDeleted() bool {
	return u.IsActive && !u.IsDeleted()
}

// Deactivate deactivates the user account
func (u *User) Deactivate() {
	u.IsActive = false
	u.UpdatedAt = time.Now()
}

// Activate activates the user account
func (u *User) Activate() {
	u.IsActive = true
	u.UpdatedAt = time.Now()
}

// SoftDelete marks the user as deleted (soft delete)
func (u *User) SoftDelete() {
	u.DeletedAt = gorm.DeletedAt{Time: time.Now(), Valid: true}
	u.UpdatedAt = time.Now()
}

// Restore restores a soft-deleted user
func (u *User) Restore() {
	u.DeletedAt = gorm.DeletedAt{Time: time.Time{}, Valid: false}
	u.UpdatedAt = time.Now()
}

// BeforeCreate is a GORM hook that runs before creating a user
func (u *User) BeforeCreate(tx *gorm.DB) error {
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a user
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	u.UpdatedAt = time.Now()
	return nil
}

// Helper functions
func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isValidRole(role string) bool {
	return strings.TrimSpace(role) != ""
}
