package seed_initial_user

import (
	"fmt"

	"invitation-api/internal/domain/user"
	userRepo "invitation-api/internal/repository/user"
	userService "invitation-api/internal/service/user"

	"github.com/google/uuid"
)

// defaultMenus defines the master list of CMS menus to seed
var defaultMenus = []struct {
	Title string
	Path  string
	Icon  string
	Order int
	Group string // logical grouping label (not stored, just for categorisation)
}{
	// Main
	{Title: "Dashboard", Path: "/dashboard", Icon: "dashboard", Order: 1, Group: "Main"},
	{Title: "Invitations", Path: "/invitations", Icon: "mail", Order: 2, Group: "Main"},
	{Title: "Themes", Path: "/themes", Icon: "palette", Order: 3, Group: "Main"},
	{Title: "Categories", Path: "/categories", Icon: "grid", Order: 4, Group: "Main"},
	{Title: "RSVP Responses", Path: "/responses", Icon: "clipboard", Order: 5, Group: "Main"},
	{Title: "Guest Messages", Path: "/messages", Icon: "chat", Order: 6, Group: "Main"},

	// Administration
	{Title: "Users", Path: "/users", Icon: "users", Order: 10, Group: "Administration"},
	{Title: "Roles", Path: "/roles", Icon: "key", Order: 11, Group: "Administration"},

	// Settings
	{Title: "Settings", Path: "/settings", Icon: "settings", Order: 20, Group: "Settings"},
}

// CreateInitialUser creates the initial admin user if it doesn't exist
func CreateInitialUser() {
	userRepository := userRepo.NewRepository()
	userSvc := userService.NewService(userRepository)

	// Check if admin user already exists by email
	existingUser, err := userSvc.GetUserByEmail("admin@invitation.com")
	if err == nil && existingUser != nil {
		// Ensure password and role are reset for dev
		existingUser.Password = "admin1234"
		existingUser.Role = "Super Admin"
		if updateErr := userRepository.Update(existingUser); updateErr != nil {
			fmt.Printf("⚠️  Failed to reset admin user: %v\n", updateErr)
		} else {
			fmt.Println("✅ Initial admin user updated (password reset & role set to Super Admin).")
		}
		return
	}

	// Also check by username
	existingUser, err = userSvc.GetUserByUsername("admin")
	if err == nil && existingUser != nil {
		existingUser.Password = "admin1234"
		existingUser.Role = "Super Admin"
		if updateErr := userRepository.Update(existingUser); updateErr != nil {
			fmt.Printf("⚠️  Failed to reset admin user: %v\n", updateErr)
		} else {
			fmt.Println("✅ Initial admin user updated (password reset & role set to Super Admin).")
		}
		return
	}

	// Create initial admin user
	res, err := userSvc.RegisterUser(
		"admin",
		"admin@invitation.com",
		"admin1234",
		"System Admin",
		user.UserRole("Super Admin"),
	)
	if err != nil {
		fmt.Printf("Warning: Failed to create initial admin user: %v\n", err)
		return
	}

	fmt.Printf("✅ Successfully created initial admin user: %v\n", res.Email)
}

// SeedDefaultMenus seeds all default menus if they don't exist yet
func SeedDefaultMenus() {
	userRepository := userRepo.NewRepository()

	// Seed missing menus
	for _, m := range defaultMenus {
		existing, _ := userRepository.ListMenus()
		found := false
		for _, em := range existing {
			if em.Path == m.Path {
				found = true
				break
			}
			for _, child := range em.Children {
				if child.Path == m.Path {
					found = true
					break
				}
			}
		}

		if found {
			continue
		}

		menu := user.NewMenu(m.Title, m.Path, m.Icon, m.Order, nil)
		if err := userRepository.CreateMenu(menu); err != nil {
			fmt.Printf("⚠️  Failed to seed menu '%s': %v\n", m.Title, err)
		} else {
			fmt.Printf("✅ Seeded menu: %s (%s)\n", m.Title, m.Path)
		}
	}

	// Cleanup: remove menus that are no longer in defaultMenus
	allMenus, _ := userRepository.ListMenus()
	for _, em := range allMenus {
		found := false
		for _, m := range defaultMenus {
			if em.Path == m.Path {
				found = true
				break
			}
		}
		if !found {
			if err := userRepository.DeleteMenu(em.ID); err != nil {
				fmt.Printf("⚠️  Failed to remove old menu '%s': %v\n", em.Title, err)
			} else {
				fmt.Printf("🗑️  Removed menu no longer in code: %s (%s)\n", em.Title, em.Path)
			}
		}
	}

	fmt.Println("✅ Default menus seeding completed")
}

// SeedSuperAdminRole creates a "Super Admin" role and assigns ALL menus to it
func SeedSuperAdminRole() {
	userRepository := userRepo.NewRepository()

	// 1. Check if "Super Admin" role already exists
	existingRole, err := userRepository.GetRoleByName("Super Admin")
	if err == nil && existingRole != nil {
		fmt.Println("✅ Super Admin role already exists. Ensuring all menus are assigned...")
		assignAllMenusToRole(userRepository, existingRole.ID)
		return
	}

	// 2. Create "Super Admin" role
	role := user.NewRole("Super Admin", "Full access to all features and menus")
	if err := userRepository.CreateRole(role); err != nil {
		fmt.Printf("⚠️  Failed to create Super Admin role: %v\n", err)
		return
	}

	fmt.Printf("✅ Created Super Admin role: %s\n", role.ID)

	// 3. Assign all menus to the role
	assignAllMenusToRole(userRepository, role.ID)
}

// assignAllMenusToRole assigns every menu in the DB to the given role
func assignAllMenusToRole(repo userRepo.Repository, roleID uuid.UUID) {
	// Get all menus (root level with children)
	menus, err := repo.ListMenus()
	if err != nil {
		fmt.Printf("⚠️  Failed to list menus for role assignment: %v\n", err)
		return
	}

	// Get currently assigned menus
	assignedMenus, _ := repo.GetMenusByRole(roleID)
	assignedMap := make(map[uuid.UUID]bool)
	for _, m := range assignedMenus {
		assignedMap[m.ID] = true
	}

	for _, menu := range menus {
		// Assign parent
		if !assignedMap[menu.ID] {
			if err := repo.AssignMenuToRole(roleID, menu.ID); err != nil {
				fmt.Printf("⚠️  Failed to assign menu '%s' to Super Admin: %v\n", menu.Title, err)
			}
		}
		// Assign children
		for _, child := range menu.Children {
			if !assignedMap[child.ID] {
				if err := repo.AssignMenuToRole(roleID, child.ID); err != nil {
					fmt.Printf("⚠️  Failed to assign menu '%s' to Super Admin: %v\n", child.Title, err)
				}
			}
		}
	}
	fmt.Println("✅ All menus assigned to Super Admin role")
}
