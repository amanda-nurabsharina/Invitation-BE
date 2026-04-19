package auth

import (
	"invitation-api/internal/database"
	"invitation-api/internal/domain/auth"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for auth data operations
type Repository interface {
	CreateToken(token *auth.Token) error
	GetTokenByValue(tokenValue string) (*auth.Token, error)
	GetTokensByUserID(userID uuid.UUID) ([]*auth.Token, error)
	GetValidTokensByUserID(userID uuid.UUID) ([]*auth.Token, error)
	RevokeToken(tokenID uint) error
	RevokeAllUserTokens(userID uuid.UUID) error
	DeleteExpiredTokens() error
	DeleteToken(tokenID uint) error
}

// repository implements the Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new auth repository instance
func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

// CreateToken creates a new token in the database
func (r *repository) CreateToken(token *auth.Token) error {
	return r.db.Create(token).Error
}

// GetTokenByValue retrieves a token by its value
func (r *repository) GetTokenByValue(tokenValue string) (*auth.Token, error) {
	var token auth.Token
	err := r.db.Where("token = ?", tokenValue).First(&token).Error
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// GetTokensByUserID retrieves all tokens for a user
func (r *repository) GetTokensByUserID(userID uuid.UUID) ([]*auth.Token, error) {
	var tokens []*auth.Token
	err := r.db.Where("user_id = ?", userID).Find(&tokens).Error
	return tokens, err
}

// GetValidTokensByUserID retrieves all valid (not expired and not revoked) tokens for a user
func (r *repository) GetValidTokensByUserID(userID uuid.UUID) ([]*auth.Token, error) {
	var tokens []*auth.Token
	err := r.db.Where("user_id = ? AND is_revoked = ? AND expires_at > ?",
		userID, false, time.Now()).Find(&tokens).Error
	return tokens, err
}

// RevokeToken marks a token as revoked
func (r *repository) RevokeToken(tokenID uint) error {
	return r.db.Model(&auth.Token{}).Where("id = ?", tokenID).Update("is_revoked", true).Error
}

// RevokeAllUserTokens revokes all tokens for a specific user
func (r *repository) RevokeAllUserTokens(userID uuid.UUID) error {
	return r.db.Model(&auth.Token{}).Where("user_id = ?", userID).Update("is_revoked", true).Error
}

// DeleteExpiredTokens removes all expired tokens from the database
func (r *repository) DeleteExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&auth.Token{}).Error
}

// DeleteToken permanently deletes a token
func (r *repository) DeleteToken(tokenID uint) error {
	return r.db.Delete(&auth.Token{}, tokenID).Error
}
