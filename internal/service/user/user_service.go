package user

import (
	"errors"
	userDomain "invitation-api/internal/domain/user"
	userRepo "invitation-api/internal/repository/user"
	"time"

	"github.com/google/uuid"
)

// Service defines the interface for user service
type Service interface {
	RegisterUser(username, email, password, fullName string, role userDomain.UserRole) (*userDomain.User, error)
	GetUserByID(id uuid.UUID) (*userDomain.User, error)
	GetUserByEmail(email string) (*userDomain.User, error)
	GetUserByUsername(username string) (*userDomain.User, error)
	ValidateCredentials(identifier, password string) (*userDomain.User, error)
	ListUsers(limit, offset int) ([]*userDomain.User, error)
	UpdateUser(id uuid.UUID, username, email, role string) (*userDomain.User, error)
	SoftDeleteUser(id uuid.UUID) error

	// Role Management
	CreateRole(name, description string) (*userDomain.Role, error)
	ListRoles() ([]*userDomain.Role, error)
	GetRoleByID(id uuid.UUID) (*userDomain.Role, error)
	UpdateRole(id uuid.UUID, name, description string) (*userDomain.Role, error)
	DeleteRole(id uuid.UUID) error

	// Menu Management
	CreateMenu(title, path, icon string, order int, parentID *uuid.UUID) (*userDomain.Menu, error)
	ListMenus() ([]*userDomain.Menu, error)
	UpdateMenu(id uuid.UUID, title, path, icon string, order int) (*userDomain.Menu, error)
	DeleteMenu(id uuid.UUID) error

	// Role-Menu assignment
	AssignMenuToRole(roleID, menuID uuid.UUID) error
	RemoveMenuFromRole(roleID, menuID uuid.UUID) error
	GetMenusByRole(roleID uuid.UUID) ([]*userDomain.Menu, error)
}

// service implements Service interface
type service struct {
	userRepo userRepo.Repository
}

// NewService creates a new user service
func NewService(userRepo userRepo.Repository) Service {
	return &service{
		userRepo: userRepo,
	}
}

// RegisterUser creates a new user with profile
func (s *service) RegisterUser(username, email, password, fullName string, role userDomain.UserRole) (*userDomain.User, error) {
	// Validate role - validation logic here or rely on enum

	user, err := userDomain.NewUser(username, password, email, role)
	if err != nil {
		return nil, err
	}

	// Create profile
	profile := &userDomain.UserProfile{
		ID:        uuid.New(),
		UserID:    user.ID,
		FullName:  fullName,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	user.Profile = profile

	// Use CreateWithProfile from repo to ensure transaction
	if err := s.userRepo.CreateWithProfile(user); err != nil {
		return nil, err
	}

	// Re-fetch to ensure all fields populated if needed, or return constructed
	return user, nil
}

func (s *service) ValidateCredentials(identifier, password string) (*userDomain.User, error) {
	// Identifier can be email or username
	var user *userDomain.User
	var err error

	// Simple heuristic: if contains '@', assume email
	// Better: try both or specific check
	if s.isEmail(identifier) {
		user, err = s.userRepo.GetByEmail(identifier)
	} else {
		user, err = s.userRepo.GetByUsername(identifier)
	}

	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	// In real app use bcrypt.CompareHashAndPassword
	if user.Password != password {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	return user, nil
}

func (s *service) isEmail(input string) bool {
	// Simple check, can use regex
	for _, c := range input {
		if c == '@' {
			return true
		}
	}
	return false
}

func (s *service) GetUserByID(id uuid.UUID) (*userDomain.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *service) GetUserByEmail(email string) (*userDomain.User, error) {
	return s.userRepo.GetByEmail(email)
}

func (s *service) GetUserByUsername(username string) (*userDomain.User, error) {
	return s.userRepo.GetByUsername(username)
}

func (s *service) ListUsers(limit, offset int) ([]*userDomain.User, error) {
	return s.userRepo.List(limit, offset)
}

func (s *service) UpdateUser(id uuid.UUID, username, email, role string) (*userDomain.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.Username = username
	user.Email = email
	if role != "" {
		user.Role = role
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) SoftDeleteUser(id uuid.UUID) error {
	return s.userRepo.SoftDelete(id)
}

// --- Role Management ---

func (s *service) CreateRole(name, description string) (*userDomain.Role, error) {
	role := userDomain.NewRole(name, description)
	if err := s.userRepo.CreateRole(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *service) ListRoles() ([]*userDomain.Role, error) {
	return s.userRepo.ListRoles()
}

func (s *service) GetRoleByID(id uuid.UUID) (*userDomain.Role, error) {
	return s.userRepo.GetRoleByID(id)
}

func (s *service) UpdateRole(id uuid.UUID, name, description string) (*userDomain.Role, error) {
	role, err := s.userRepo.GetRoleByID(id)
	if err != nil {
		return nil, err
	}
	role.Name = name
	role.Description = description
	role.UpdatedAt = time.Now()

	if err := s.userRepo.UpdateRole(role); err != nil {
		return nil, err
	}
	return role, nil
}

func (s *service) DeleteRole(id uuid.UUID) error {
	return s.userRepo.DeleteRole(id)
}

// --- Menu Management ---

func (s *service) CreateMenu(title, path, icon string, order int, parentID *uuid.UUID) (*userDomain.Menu, error) {
	menu := userDomain.NewMenu(title, path, icon, order, parentID)
	if err := s.userRepo.CreateMenu(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *service) ListMenus() ([]*userDomain.Menu, error) {
	return s.userRepo.ListMenus()
}

func (s *service) UpdateMenu(id uuid.UUID, title, path, icon string, order int) (*userDomain.Menu, error) {
	menu, err := s.userRepo.GetMenuByID(id)
	if err != nil {
		return nil, err
	}
	menu.Title = title
	menu.Path = path
	menu.Icon = icon
	menu.Order = order
	menu.UpdatedAt = time.Now()

	if err := s.userRepo.UpdateMenu(menu); err != nil {
		return nil, err
	}
	return menu, nil
}

func (s *service) DeleteMenu(id uuid.UUID) error {
	return s.userRepo.DeleteMenu(id)
}

func (s *service) AssignMenuToRole(roleID, menuID uuid.UUID) error {
	return s.userRepo.AssignMenuToRole(roleID, menuID)
}

func (s *service) RemoveMenuFromRole(roleID, menuID uuid.UUID) error {
	return s.userRepo.RemoveMenuFromRole(roleID, menuID)
}

func (s *service) GetMenusByRole(roleID uuid.UUID) ([]*userDomain.Menu, error) {
	return s.userRepo.GetMenusByRole(roleID)
}
