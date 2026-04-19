package invitation

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Guest represents an invited guest
type Guest struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	InvitationID uuid.UUID `gorm:"type:uuid;column:invitation_id;not null;index" json:"invitation_id"`
	Name         string    `gorm:"column:name;not null" json:"name"`
	Type         string    `gorm:"column:type;default:regular" json:"type"` // e.g., VIP, Regular, Family

	// Optional specific greeting
	Greeting *string `gorm:"column:greeting" json:"greeting,omitempty"`

	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

// TableName returns the table name
func (Guest) TableName() string {
	return "invitation_guests"
}
