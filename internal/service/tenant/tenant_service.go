package tenant

import (
	"errors"
	"fmt"

	tenantDomain "invitation-api/internal/domain/tenant"
	tenantRepo "invitation-api/internal/repository/tenant"
	"invitation-api/pkg/config"

	"github.com/google/uuid"
)

// Service defines the interface for tenant business operations
type Service interface {
	Create(userID uuid.UUID, name, subdomain string) (*tenantDomain.Tenant, error)
	GetByID(id uuid.UUID) (*tenantDomain.Tenant, error)
	GetBySubdomain(subdomain string) (*tenantDomain.Tenant, error)
	GetByUserID(userID uuid.UUID) ([]*tenantDomain.Tenant, error)
	Update(id uuid.UUID, name string, customDomain *string) (*tenantDomain.Tenant, error)
	UpdatePlan(id uuid.UUID, plan string) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*tenantDomain.Tenant, int64, error)
	CheckSubdomainAvailability(subdomain string) (bool, error)
	Activate(id uuid.UUID) error
	Deactivate(id uuid.UUID) error
}

// service implements Service interface
type service struct {
	tenantRepo tenantRepo.Repository
	config     *config.Config
}

// NewService creates a new tenant service
func NewService(tenantRepo tenantRepo.Repository, cfg *config.Config) Service {
	return &service{
		tenantRepo: tenantRepo,
		config:     cfg,
	}
}

// Create creates a new tenant
func (s *service) Create(userID uuid.UUID, name, subdomain string) (*tenantDomain.Tenant, error) {
	// Check if subdomain is available
	exists, err := s.tenantRepo.SubdomainExists(subdomain)
	if err != nil {
		return nil, fmt.Errorf("failed to check subdomain: %w", err)
	}
	if exists {
		return nil, errors.New("subdomain is already taken")
	}

	// Create tenant
	tenant, err := tenantDomain.NewTenant(userID, name, subdomain)
	if err != nil {
		return nil, err
	}

	if err := s.tenantRepo.Create(tenant); err != nil {
		return nil, fmt.Errorf("failed to create tenant: %w", err)
	}

	return tenant, nil
}

// GetByID retrieves a tenant by ID
func (s *service) GetByID(id uuid.UUID) (*tenantDomain.Tenant, error) {
	return s.tenantRepo.GetByID(id)
}

// GetBySubdomain retrieves a tenant by subdomain
func (s *service) GetBySubdomain(subdomain string) (*tenantDomain.Tenant, error) {
	return s.tenantRepo.GetBySubdomain(subdomain)
}

// GetByUserID retrieves all tenants for a user
func (s *service) GetByUserID(userID uuid.UUID) ([]*tenantDomain.Tenant, error) {
	return s.tenantRepo.GetByUserID(userID)
}

// Update updates a tenant
func (s *service) Update(id uuid.UUID, name string, customDomain *string) (*tenantDomain.Tenant, error) {
	tenant, err := s.tenantRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		tenant.Name = name
	}
	if customDomain != nil {
		tenant.CustomDomain = customDomain
	}

	if err := s.tenantRepo.Update(tenant); err != nil {
		return nil, fmt.Errorf("failed to update tenant: %w", err)
	}

	return tenant, nil
}

// UpdatePlan updates tenant's plan
func (s *service) UpdatePlan(id uuid.UUID, plan string) error {
	tenant, err := s.tenantRepo.GetByID(id)
	if err != nil {
		return err
	}

	tenant.Plan = plan
	return s.tenantRepo.Update(tenant)
}

// Delete deletes a tenant
func (s *service) Delete(id uuid.UUID) error {
	return s.tenantRepo.Delete(id)
}

// List retrieves tenants with pagination
func (s *service) List(limit, offset int) ([]*tenantDomain.Tenant, int64, error) {
	tenants, err := s.tenantRepo.List(limit, offset)
	if err != nil {
		return nil, 0, err
	}

	count, err := s.tenantRepo.Count()
	if err != nil {
		return nil, 0, err
	}

	return tenants, count, nil
}

// CheckSubdomainAvailability checks if a subdomain is available
func (s *service) CheckSubdomainAvailability(subdomain string) (bool, error) {
	exists, err := s.tenantRepo.SubdomainExists(subdomain)
	if err != nil {
		return false, err
	}
	return !exists, nil
}

// Activate activates a tenant
func (s *service) Activate(id uuid.UUID) error {
	tenant, err := s.tenantRepo.GetByID(id)
	if err != nil {
		return err
	}

	tenant.IsActive = true
	return s.tenantRepo.Update(tenant)
}

// Deactivate deactivates a tenant
func (s *service) Deactivate(id uuid.UUID) error {
	tenant, err := s.tenantRepo.GetByID(id)
	if err != nil {
		return err
	}

	tenant.IsActive = false
	return s.tenantRepo.Update(tenant)
}
