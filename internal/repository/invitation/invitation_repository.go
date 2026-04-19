package invitation

import (
	"invitation-api/internal/database"
	invitationDomain "invitation-api/internal/domain/invitation"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Repository defines the interface for invitation data operations
type Repository interface {
	// Invitation
	Create(invitation *invitationDomain.Invitation) error
	GetByID(id uuid.UUID) (*invitationDomain.Invitation, error)
	GetByTenantID(tenantID uuid.UUID) (*invitationDomain.Invitation, error)
	Update(invitation *invitationDomain.Invitation) error
	Delete(id uuid.UUID) error
	List(limit, offset int) ([]*invitationDomain.Invitation, error)

	// Gallery
	CreateGalleryItem(item *invitationDomain.Gallery) error
	GetGalleryByInvitationID(invitationID uuid.UUID) ([]*invitationDomain.Gallery, error)
	DeleteGalleryItem(id uuid.UUID) error
	UpdateGalleryOrder(id uuid.UUID, order int) error

	// RSVP
	CreateRSVP(rsvp *invitationDomain.RSVP) error
	GetRSVPsByInvitationID(invitationID uuid.UUID) ([]*invitationDomain.RSVP, error)
	GetRSVPStats(invitationID uuid.UUID) (*RSVPStats, error)

	// Gift Accounts
	CreateGiftAccount(account *invitationDomain.GiftAccount) error
	GetGiftAccountsByInvitationID(invitationID uuid.UUID) ([]*invitationDomain.GiftAccount, error)
	DeleteGiftAccount(id uuid.UUID) error

	// Guest Messages
	CreateGuestMessage(message *invitationDomain.GuestMessage) error
	GetGuestMessagesByInvitationID(invitationID uuid.UUID, approvedOnly bool) ([]*invitationDomain.GuestMessage, error)
	ApproveGuestMessage(id uuid.UUID) error
	DeleteGuestMessage(id uuid.UUID) error

	// Guests
	CreateGuest(guest *invitationDomain.Guest) error
	GetGuestsByInvitationID(invitationID uuid.UUID) ([]*invitationDomain.Guest, error)
	DeleteGuest(id uuid.UUID) error
}

// RSVPStats holds RSVP statistics
type RSVPStats struct {
	Total      int64 `json:"total"`
	Hadir      int64 `json:"hadir"`
	TidakHadir int64 `json:"tidak_hadir"`
	Ragu       int64 `json:"ragu"`
	GuestCount int64 `json:"guest_count"`
}

// repository implements Repository interface
type repository struct {
	db *gorm.DB
}

// NewRepository creates a new invitation repository
func NewRepository() Repository {
	return &repository{
		db: database.GetDB(),
	}
}

// Create creates a new invitation
func (r *repository) Create(invitation *invitationDomain.Invitation) error {
	return r.db.Create(invitation).Error
}

// GetByID retrieves an invitation by ID
func (r *repository) GetByID(id uuid.UUID) (*invitationDomain.Invitation, error) {
	var invitation invitationDomain.Invitation
	err := r.db.Preload("Gallery").Preload("GiftAccounts").First(&invitation, id).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// GetByTenantID retrieves invitation by tenant ID
func (r *repository) GetByTenantID(tenantID uuid.UUID) (*invitationDomain.Invitation, error) {
	var invitation invitationDomain.Invitation
	err := r.db.Preload("Gallery").Preload("GiftAccounts").
		Where("tenant_id = ?", tenantID).First(&invitation).Error
	if err != nil {
		return nil, err
	}
	return &invitation, nil
}

// Update updates an invitation
func (r *repository) Update(invitation *invitationDomain.Invitation) error {
	return r.db.Save(invitation).Error
}

// Delete soft deletes an invitation
func (r *repository) Delete(id uuid.UUID) error {
	return r.db.Delete(&invitationDomain.Invitation{}, id).Error
}

// List retrieves invitations with pagination
func (r *repository) List(limit, offset int) ([]*invitationDomain.Invitation, error) {
	var invitations []*invitationDomain.Invitation
	err := r.db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&invitations).Error
	return invitations, err
}

// CreateGalleryItem creates a new gallery item
func (r *repository) CreateGalleryItem(item *invitationDomain.Gallery) error {
	return r.db.Create(item).Error
}

// GetGalleryByInvitationID retrieves gallery for an invitation
func (r *repository) GetGalleryByInvitationID(invitationID uuid.UUID) ([]*invitationDomain.Gallery, error) {
	var gallery []*invitationDomain.Gallery
	err := r.db.Where("invitation_id = ?", invitationID).Order("sort_order").Find(&gallery).Error
	return gallery, err
}

// DeleteGalleryItem deletes a gallery item
func (r *repository) DeleteGalleryItem(id uuid.UUID) error {
	return r.db.Delete(&invitationDomain.Gallery{}, id).Error
}

// UpdateGalleryOrder updates gallery item order
func (r *repository) UpdateGalleryOrder(id uuid.UUID, order int) error {
	return r.db.Model(&invitationDomain.Gallery{}).Where("id = ?", id).Update("sort_order", order).Error
}

// CreateRSVP creates a new RSVP response
func (r *repository) CreateRSVP(rsvp *invitationDomain.RSVP) error {
	return r.db.Create(rsvp).Error
}

// GetRSVPsByInvitationID retrieves RSVPs for an invitation
func (r *repository) GetRSVPsByInvitationID(invitationID uuid.UUID) ([]*invitationDomain.RSVP, error) {
	var rsvps []*invitationDomain.RSVP
	err := r.db.Where("invitation_id = ?", invitationID).Order("created_at DESC").Find(&rsvps).Error
	return rsvps, err
}

// GetRSVPStats gets RSVP statistics
func (r *repository) GetRSVPStats(invitationID uuid.UUID) (*RSVPStats, error) {
	stats := &RSVPStats{}

	// Total responses
	r.db.Model(&invitationDomain.RSVP{}).Where("invitation_id = ?", invitationID).Count(&stats.Total)

	// By attendance
	r.db.Model(&invitationDomain.RSVP{}).Where("invitation_id = ? AND attendance = ?", invitationID, "hadir").Count(&stats.Hadir)
	r.db.Model(&invitationDomain.RSVP{}).Where("invitation_id = ? AND attendance = ?", invitationID, "tidak_hadir").Count(&stats.TidakHadir)
	r.db.Model(&invitationDomain.RSVP{}).Where("invitation_id = ? AND attendance = ?", invitationID, "ragu").Count(&stats.Ragu)

	// Total guest count (only from those attending)
	var guestCount int64
	r.db.Model(&invitationDomain.RSVP{}).Where("invitation_id = ? AND attendance = ?", invitationID, "hadir").
		Select("COALESCE(SUM(guest_count), 0)").Scan(&guestCount)
	stats.GuestCount = guestCount

	return stats, nil
}

// CreateGiftAccount creates a new gift account
func (r *repository) CreateGiftAccount(account *invitationDomain.GiftAccount) error {
	return r.db.Create(account).Error
}

// GetGiftAccountsByInvitationID retrieves gift accounts for an invitation
func (r *repository) GetGiftAccountsByInvitationID(invitationID uuid.UUID) ([]*invitationDomain.GiftAccount, error) {
	var accounts []*invitationDomain.GiftAccount
	err := r.db.Where("invitation_id = ?", invitationID).Find(&accounts).Error
	return accounts, err
}

// DeleteGiftAccount deletes a gift account
func (r *repository) DeleteGiftAccount(id uuid.UUID) error {
	return r.db.Delete(&invitationDomain.GiftAccount{}, id).Error
}

// CreateGuestMessage creates a new guest message
func (r *repository) CreateGuestMessage(message *invitationDomain.GuestMessage) error {
	return r.db.Create(message).Error
}

// GetGuestMessagesByInvitationID retrieves messages for an invitation
func (r *repository) GetGuestMessagesByInvitationID(invitationID uuid.UUID, approvedOnly bool) ([]*invitationDomain.GuestMessage, error) {
	var messages []*invitationDomain.GuestMessage
	query := r.db.Where("invitation_id = ?", invitationID)
	if approvedOnly {
		query = query.Where("is_approved = ?", true)
	}
	err := query.Order("created_at DESC").Find(&messages).Error
	return messages, err
}

// ApproveGuestMessage approves a message
func (r *repository) ApproveGuestMessage(id uuid.UUID) error {
	return r.db.Model(&invitationDomain.GuestMessage{}).Where("id = ?", id).Update("is_approved", true).Error
}

// DeleteGuestMessage deletes a message
func (r *repository) DeleteGuestMessage(id uuid.UUID) error {
	return r.db.Delete(&invitationDomain.GuestMessage{}, id).Error
}

// CreateGuest creates a new guest
func (r *repository) CreateGuest(guest *invitationDomain.Guest) error {
	return r.db.Create(guest).Error
}

// GetGuestsByInvitationID retrieves guests for an invitation
func (r *repository) GetGuestsByInvitationID(invitationID uuid.UUID) ([]*invitationDomain.Guest, error) {
	var guests []*invitationDomain.Guest
	err := r.db.Where("invitation_id = ?", invitationID).Order("created_at DESC").Find(&guests).Error
	return guests, err
}

// DeleteGuest deletes a guest
func (r *repository) DeleteGuest(id uuid.UUID) error {
	return r.db.Delete(&invitationDomain.Guest{}, id).Error
}
