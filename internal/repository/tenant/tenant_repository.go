package tenant

import (
	"invitation-api/internal/database"
	tenantDomain "invitation-api/internal/domain/tenant"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for tenant data operations
type Repository interface {
	Create(tenant *tenantDomain.Tenant) error
	GetByID(id uuid.UUID) (*tenantDomain.Tenant, error)
	GetBySubdomain(subdomain string) (*tenantDomain.Tenant, error)
	GetByUserID(userID uuid.UUID) ([]*tenantDomain.Tenant, error)
	Update(tenant *tenantDomain.Tenant) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*tenantDomain.Tenant, error)
	Count() (int64, error)
	SubdomainExists(subdomain string) (bool, error)
}

// repository implements Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new tenant repository instance
func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

// Create creates a new tenant
func (r *repository) Create(tenant *tenantDomain.Tenant) error {
	return r.db.Create(tenant).Error
}

// GetByID retrieves a tenant by ID
func (r *repository) GetByID(id uuid.UUID) (*tenantDomain.Tenant, error) {
	var tenant tenantDomain.Tenant
	err := r.db.First(&tenant, id).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetBySubdomain retrieves a tenant by subdomain
func (r *repository) GetBySubdomain(subdomain string) (*tenantDomain.Tenant, error) {
	var tenant tenantDomain.Tenant
	err := r.db.Where("subdomain = ?", subdomain).First(&tenant).Error
	if err != nil {
		return nil, err
	}
	return &tenant, nil
}

// GetByUserID retrieves all tenants for a user
func (r *repository) GetByUserID(userID uuid.UUID) ([]*tenantDomain.Tenant, error) {
	var tenants []*tenantDomain.Tenant
	err := r.db.Where("user_id = ?", userID).Find(&tenants).Error
	return tenants, err
}

// Update updates a tenant
func (r *repository) Update(tenant *tenantDomain.Tenant) error {
	return r.db.Save(tenant).Error
}

// Delete soft deletes a tenant
func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&tenantDomain.Tenant{}, id).Error
}

// List retrieves tenants with pagination
func (r *repository) List(limit, offset int) ([]*tenantDomain.Tenant, error) {
	var tenants []*tenantDomain.Tenant
	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&tenants).Error
	return tenants, err
}

// Count returns total number of tenants
func (r *repository) Count() (int64, error) {
	var count int64
	err := r.db.Model(&tenantDomain.Tenant{}).Count(&count).Error
	return count, err
}

// SubdomainExists checks if subdomain is already taken
func (r *repository) SubdomainExists(subdomain string) (bool, error) {
	var count int64
	err := r.db.Model(&tenantDomain.Tenant{}).Where("subdomain = ?", subdomain).Count(&count).Error
	return count > 0, err
}
