package user

import (
	"time"

	"github.com/google/uuid"
)

// UserProfile represents the user's profile information
type UserProfile struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	UserID    uuid.UUID `gorm:"type:uuid;column:user_id;not null;uniqueIndex" json:"user_id"`
	FullName  string    `gorm:"column:full_name" json:"full_name"`
	Phone     string    `gorm:"column:phone" json:"phone"`
	Address   string    `gorm:"column:address" json:"address"`
	Avatar    string    `gorm:"column:avatar" json:"avatar"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
}

// NewUserProfile creates a new user profile instance
func NewUserProfile(userID uuid.UUID, fullName, phone, address string) *UserProfile {
	return &UserProfile{
		ID:        uuid.New(),
		UserID:    userID,
		FullName:  fullName,
		Phone:     phone,
		Address:   address,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// TableName returns the table name for UserProfile
func (UserProfile) TableName() string {
	return "user_profile"
}
