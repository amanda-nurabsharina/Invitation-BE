package database

import (
	"fmt"
	"log"
	"time"

	authDomain "invitation-api/internal/domain/auth"
	invitationDomain "invitation-api/internal/domain/invitation"
	tenantDomain "invitation-api/internal/domain/tenant"
	themeDomain "invitation-api/internal/domain/theme"
	userDomain "invitation-api/internal/domain/user"
	"invitation-api/pkg/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// DB holds the database connection
var DB *gorm.DB

// Connect establishes a connection to the database
func Connect(config *config.Config) error {
	var err error

	// Configure GORM logger
	gormLogger := logger.Default
	if config.IsDevelopment() {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	// Connect to database
	DB, err = gorm.Open(postgres.Open(config.GetDSN()), &gorm.Config{
		Logger: gormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get underlying sql.DB object
	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Configure connection pool
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connected successfully")
	return nil
}

// AutoMigrate runs database migrations
func AutoMigrate() error {
	if DB == nil {
		return fmt.Errorf("database connection not established")
	}

	// Run GORM's auto migrations
	err := DB.AutoMigrate(
		// User & Auth
		&userDomain.User{},
		&userDomain.UserProfile{},
		&authDomain.Token{},
		&userDomain.Role{},
		&userDomain.Menu{},

		// Multi-tenant
		&tenantDomain.Tenant{},

		// Theme
		&themeDomain.Theme{},
		&themeDomain.ThemeCategory{},

		// Invitation
		&invitationDomain.Invitation{},
		&invitationDomain.Gallery{},
		&invitationDomain.RSVP{},
		&invitationDomain.GiftAccount{},
		&invitationDomain.GuestMessage{},
		&invitationDomain.Guest{},
	)

	if err != nil {
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

// Close closes the database connection
func Close() error {
	if DB == nil {
		return nil
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	return sqlDB.Close()
}

// GetDB returns the database instance
func GetDB() *gorm.DB {
	return DB
}
