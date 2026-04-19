package user

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Role represents a system role
type Role struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	Name        string         `gorm:"column:name;uniqueIndex;not null" json:"name"` // e.g., "admin", "user", "editor"
	Description string         `gorm:"column:description" json:"description"`
	Menus       []*Menu        `gorm:"many2many:role_menus;" json:"menus"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// Menu represents a navigation menu item in the CMS
type Menu struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	Title       string         `gorm:"column:title;not null" json:"title"`
	Path        string         `gorm:"column:path;not null" json:"path"` // Route path
	Icon        string         `gorm:"column:icon" json:"icon"`
	ParentID    *uuid.UUID     `gorm:"column:parent_id" json:"parent_id,omitempty"`
	Order       int            `gorm:"column:order;default:0" json:"order"`
	IsActive    bool           `gorm:"column:is_active;default:true" json:"is_active"`
	Roles       []*Role        `gorm:"many2many:role_menus;" json:"-"`
	Children    []*Menu        `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// NewRole creates a new role instance
func NewRole(name, description string) *Role {
	return &Role{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// NewMenu creates a new menu instance
func NewMenu(title, path, icon string, order int, parentID *uuid.UUID) *Menu {
	return &Menu{
		ID:        uuid.New(),
		Title:     title,
		Path:      path,
		Icon:      icon,
		Order:     order,
		ParentID:  parentID,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// TableName overrides the table name
func (Role) TableName() string {
	return "role"
}

// TableName overrides the table name
func (Menu) TableName() string {
	return "menu"
}
