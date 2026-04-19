package tenant

import (
	"errors"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Tenant represents a client/customer in the multi-tenant system
type Tenant struct {
	ID           uuid.UUID       `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	UserID       uuid.UUID       `gorm:"type:uuid;column:user_id" json:"user_id"` // Owner
	Name         string          `gorm:"column:name;not null" json:"name"`
	Subdomain    string          `gorm:"column:subdomain;uniqueIndex;not null" json:"subdomain"`
	CustomDomain *string         `gorm:"column:custom_domain" json:"custom_domain,omitempty"`
	IsActive     bool            `gorm:"column:is_active;default:true" json:"is_active"`
	Plan         string          `gorm:"column:plan;default:basic" json:"plan"` // basic, premium, enterprise
	ExpiresAt    *time.Time      `gorm:"column:expires_at" json:"expires_at,omitempty"`
	Settings     *TenantSettings `gorm:"serializer:json;column:settings" json:"settings,omitempty"`
	CreatedAt    time.Time       `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time       `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt  `gorm:"column:deleted_at;index" json:"deleted_at"`

	// Relationships
	// Invitation *Invitation `gorm:"foreignKey:TenantID" json:"invitation,omitempty"`
}

// TenantSettings holds tenant-specific settings
type TenantSettings struct {
	PrimaryColor   string `json:"primary_color,omitempty"`
	SecondaryColor string `json:"secondary_color,omitempty"`
	Logo           string `json:"logo,omitempty"`
	Favicon        string `json:"favicon,omitempty"`
}

// Plan constants
const (
	PlanBasic      = "basic"
	PlanPremium    = "premium"
	PlanEnterprise = "enterprise"
)

// NewTenant creates a new tenant instance
func NewTenant(userID uuid.UUID, name, subdomain string) (*Tenant, error) {
	tenant := &Tenant{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      name,
		Subdomain: strings.ToLower(subdomain),
		IsActive:  true,
		Plan:      PlanBasic,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	if err := tenant.Validate(); err != nil {
		return nil, err
	}

	return tenant, nil
}

// Validate performs business rule validation
func (t *Tenant) Validate() error {
	if strings.TrimSpace(t.Name) == "" {
		return errors.New("tenant name is required")
	}

	if strings.TrimSpace(t.Subdomain) == "" {
		return errors.New("subdomain is required")
	}

	// Subdomain must be lowercase alphanumeric with optional hyphens
	subdomainRegex := regexp.MustCompile(`^[a-z0-9][a-z0-9-]*[a-z0-9]$|^[a-z0-9]$`)
	if !subdomainRegex.MatchString(t.Subdomain) {
		return errors.New("subdomain must be lowercase alphanumeric with optional hyphens")
	}

	// Reserved subdomains
	reserved := []string{"www", "api", "admin", "cms", "app", "mail", "smtp", "ftp", "cdn", "static", "assets"}
	for _, r := range reserved {
		if t.Subdomain == r {
			return errors.New("this subdomain is reserved")
		}
	}

	if len(t.Subdomain) < 3 {
		return errors.New("subdomain must be at least 3 characters")
	}

	if len(t.Subdomain) > 50 {
		return errors.New("subdomain must be less than 50 characters")
	}

	return nil
}

// IsExpired checks if the tenant subscription has expired
func (t *Tenant) IsExpired() bool {
	if t.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*t.ExpiresAt)
}

// IsAccessible checks if the tenant is active and not expired
func (t *Tenant) IsAccessible() bool {
	return t.IsActive && !t.IsExpired()
}

// TableName returns the table name for Tenant
func (Tenant) TableName() string {
	return "tenant"
}

// BeforeCreate is a GORM hook that runs before creating a tenant
func (t *Tenant) BeforeCreate(tx *gorm.DB) error {
	t.CreatedAt = time.Now()
	t.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate is a GORM hook that runs before updating a tenant
func (t *Tenant) BeforeUpdate(tx *gorm.DB) error {
	t.UpdatedAt = time.Now()
	return nil
}
