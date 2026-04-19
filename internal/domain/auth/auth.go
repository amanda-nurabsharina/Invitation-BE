package auth

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Token represents an authentication token
type Token struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uuid.UUID `json:"user_id" gorm:"not null"`
	Token     string    `json:"token" gorm:"uniqueIndex;not null"`
	Type      TokenType `json:"type" gorm:"not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	IsRevoked bool      `json:"is_revoked" gorm:"default:false"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TokenType represents the type of token
type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password" validate:"required,min=8"`
}

// RegisterRequest represents a registration request
type RegisterRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Username  string `json:"username" validate:"required"`
	Password  string `json:"password" validate:"required,min=8"`
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name"`
	Role      string `json:"role" validate:"required,oneof=admin user"`
}

// TokenResponse represents a token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int64  `json:"expires_in"`
}

// RefreshTokenRequest represents a refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// NewToken creates a new token instance
func NewToken(userID uuid.UUID, token string, tokenType TokenType, expiresAt time.Time) *Token {
	return &Token{
		UserID:    userID,
		Token:     token,
		Type:      tokenType,
		ExpiresAt: expiresAt,
		IsRevoked: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// Validate performs validation on the token
func (t *Token) Validate() error {
	if t.UserID == uuid.Nil {
		return errors.New("user ID is required")
	}

	if strings.TrimSpace(t.Token) == "" {
		return errors.New("token is required")
	}

	if !isValidTokenType(t.Type) {
		return errors.New("invalid token type")
	}

	if t.ExpiresAt.Before(time.Now()) {
		return errors.New("token has already expired")
	}

	return nil
}

// IsExpired checks if the token is expired
func (t *Token) IsExpired() bool {
	return time.Now().After(t.ExpiresAt)
}

// IsValid checks if the token is valid (not expired and not revoked)
func (t *Token) IsValid() bool {
	return !t.IsExpired() && !t.IsRevoked
}

// Revoke marks the token as revoked
func (t *Token) Revoke() {
	t.IsRevoked = true
	t.UpdatedAt = time.Now()
}

// ValidateLoginRequest validates login request
func (req *LoginRequest) Validate() error {
	// Either username or email is required
	if strings.TrimSpace(req.Username) == "" && strings.TrimSpace(req.Email) == "" {
		return errors.New("username or email is required")
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password is required")
	}

	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	return nil
}

// GetIdentifier returns username or email for login
func (req *LoginRequest) GetIdentifier() string {
	if strings.TrimSpace(req.Username) != "" {
		return req.Username
	}
	return req.Email
}

// ValidateRegisterRequest validates registration request
func (req *RegisterRequest) Validate() error {
	if strings.TrimSpace(req.Email) == "" {
		return errors.New("email is required")
	}

	if !isValidEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if strings.TrimSpace(req.Username) == "" {
		return errors.New("username is required")
	}

	if strings.TrimSpace(req.Password) == "" {
		return errors.New("password is required")
	}

	if len(req.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}

	if strings.TrimSpace(req.FirstName) == "" {
		return errors.New("first name is required")
	}

	if !isValidRole(req.Role) {
		return errors.New("invalid role")
	}

	return nil
}

// Helper functions
func isValidTokenType(tokenType TokenType) bool {
	switch tokenType {
	case TokenTypeAccess, TokenTypeRefresh:
		return true
	default:
		return false
	}
}

func isValidEmail(email string) bool {
	return strings.Contains(email, "@") && strings.Contains(email, ".")
}

func isValidRole(role string) bool {
	return strings.TrimSpace(role) != ""
}
