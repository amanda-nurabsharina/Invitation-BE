package invitation

import (
	"time"

	"github.com/google/uuid"
)

// Gallery represents an image in the invitation gallery
type Gallery struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	InvitationID uuid.UUID `gorm:"type:uuid;column:invitation_id;not null;index" json:"invitation_id"`
	ImageURL     string    `gorm:"column:image_url;not null" json:"image_url"`
	ThumbnailURL *string   `gorm:"column:thumbnail_url" json:"thumbnail_url,omitempty"`
	Caption      *string   `gorm:"column:caption" json:"caption,omitempty"`
	SortOrder    int       `gorm:"column:sort_order;default:0" json:"sort_order"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

// NewGallery creates a new gallery item
func NewGallery(invitationID uuid.UUID, imageURL string) *Gallery {
	return &Gallery{
		ID:           uuid.New(),
		InvitationID: invitationID,
		ImageURL:     imageURL,
		SortOrder:    0,
		CreatedAt:    time.Now(),
	}
}

// TableName returns the table name
func (Gallery) TableName() string {
	return "gallery"
}

// RSVP represents a guest's RSVP response
type RSVP struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	InvitationID uuid.UUID `gorm:"type:uuid;column:invitation_id;not null;index" json:"invitation_id"`
	GuestName    string    `gorm:"column:guest_name;not null" json:"guest_name"`
	Attendance   string    `gorm:"column:attendance;not null" json:"attendance"` // hadir, tidak_hadir, ragu
	GuestCount   int       `gorm:"column:guest_count;default:1" json:"guest_count"`
	Message      *string   `gorm:"column:message" json:"message,omitempty"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

// Attendance constants
const (
	AttendanceHadir      = "hadir"
	AttendanceTidakHadir = "tidak_hadir"
	AttendanceRagu       = "ragu"
)

// NewRSVP creates a new RSVP response
func NewRSVP(invitationID uuid.UUID, guestName, attendance string, guestCount int) *RSVP {
	return &RSVP{
		ID:           uuid.New(),
		InvitationID: invitationID,
		GuestName:    guestName,
		Attendance:   attendance,
		GuestCount:   guestCount,
		CreatedAt:    time.Now(),
	}
}

// TableName returns the table name
func (RSVP) TableName() string {
	return "rsvp"
}

// GiftAccount represents a bank account for digital gifts
type GiftAccount struct {
	ID            uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	InvitationID  uuid.UUID `gorm:"type:uuid;column:invitation_id;not null;index" json:"invitation_id"`
	BankName      string    `gorm:"column:bank_name;not null" json:"bank_name"`
	AccountNumber string    `gorm:"column:account_number;not null" json:"account_number"`
	AccountName   string    `gorm:"column:account_name;not null" json:"account_name"`
	QRImage       *string   `gorm:"column:qr_image" json:"qr_image,omitempty"`
	CreatedAt     time.Time `gorm:"column:created_at" json:"created_at"`
}

// NewGiftAccount creates a new gift account
func NewGiftAccount(invitationID uuid.UUID, bankName, accountNumber, accountName string) *GiftAccount {
	return &GiftAccount{
		ID:            uuid.New(),
		InvitationID:  invitationID,
		BankName:      bankName,
		AccountNumber: accountNumber,
		AccountName:   accountName,
		CreatedAt:     time.Now(),
	}
}

// TableName returns the table name
func (GiftAccount) TableName() string {
	return "gift_account"
}

// GuestMessage represents a message from a guest
type GuestMessage struct {
	ID           uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	InvitationID uuid.UUID `gorm:"type:uuid;column:invitation_id;not null;index" json:"invitation_id"`
	GuestName    string    `gorm:"column:guest_name;not null" json:"guest_name"`
	Message      string    `gorm:"column:message;not null" json:"message"`
	IsApproved   bool      `gorm:"column:is_approved;default:false" json:"is_approved"`
	CreatedAt    time.Time `gorm:"column:created_at" json:"created_at"`
}

// NewGuestMessage creates a new guest message
func NewGuestMessage(invitationID uuid.UUID, guestName, message string) *GuestMessage {
	return &GuestMessage{
		ID:           uuid.New(),
		InvitationID: invitationID,
		GuestName:    guestName,
		Message:      message,
		IsApproved:   false,
		CreatedAt:    time.Now(),
	}
}

// Approve approves the message
func (m *GuestMessage) Approve() {
	m.IsApproved = true
}

// TableName returns the table name
func (GuestMessage) TableName() string {
	return "guest_message"
}
