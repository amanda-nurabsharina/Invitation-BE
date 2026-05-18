package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	App      AppConfig
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Port         string
	Host         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

// AppConfig holds application configuration
type AppConfig struct {
	Environment string
	LogLevel    string
	BaseURL     string
	BaseDomain  string // Base domain for subdomain resolution (e.g., "wedding.com")
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// Load .env.local first if it exists for local development overrides
	_ = godotenv.Load(".env.local")

	// Load .env file next for fallback / production configuration
	if err := godotenv.Load(); err != nil {
		// It's okay if .env file doesn't exist in production
		fmt.Println("No .env file found, using environment variables")
	}

	config := &Config{}

	// Server configuration
	config.Server = ServerConfig{
		Port:         getEnv("SERVER_PORT", "8080"),
		Host:         getEnv("SERVER_HOST", "0.0.0.0"),
		ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
		WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
		IdleTimeout:  getDurationEnv("SERVER_IDLE_TIMEOUT", 120*time.Second),
	}

	// Database configuration
	config.Database = DatabaseConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "admin12345"),
		DBName:   getEnv("DB_NAME", "invitation_db"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// JWT configuration
	config.JWT = JWTConfig{
		SecretKey:       getEnv("JWT_SECRET_KEY", "your-secret-key-change-in-production"),
		AccessTokenTTL:  getDurationEnv("JWT_ACCESS_TOKEN_TTL", 60*time.Minute),
		RefreshTokenTTL: getDurationEnv("JWT_REFRESH_TOKEN_TTL", 7*24*time.Hour),
		Issuer:          getEnv("JWT_ISSUER", "invitation-api"),
	}

	// App configuration
	config.App = AppConfig{
		Environment: getEnv("APP_ENV", "development"),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		BaseURL:     getEnv("APP_BASE_URL", "http://localhost:8080"),
		BaseDomain:  getEnv("APP_BASE_DOMAIN", "localhost"),
	}

	return config, nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		c.Database.Host,
		c.Database.Port,
		c.Database.User,
		c.Database.Password,
		c.Database.DBName,
		c.Database.SSLMode,
	)
}

// IsDevelopment checks if the application is running in development mode
func (c *Config) IsDevelopment() bool {
	return c.App.Environment == "development"
}

// IsProduction checks if the application is running in production mode
func (c *Config) IsProduction() bool {
	return c.App.Environment == "production"
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}

func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
