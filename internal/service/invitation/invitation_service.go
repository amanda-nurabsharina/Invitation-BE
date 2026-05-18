package invitation

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"time"

	invitationDomain "invitation-api/internal/domain/invitation"
	invitationRepo "invitation-api/internal/repository/invitation"
	tenantRepo "invitation-api/internal/repository/tenant"
	"invitation-api/pkg/config"

	"github.com/google/uuid"
)

// Service defines the interface for invitation business operations
type Service interface {
	// Invitation CRUD
	Create(tenantID uuid.UUID) (*invitationDomain.Invitation, error)
	GetByID(id uuid.UUID) (*invitationDomain.Invitation, error)
	GetByTenantID(tenantID uuid.UUID) (*invitationDomain.Invitation, error)
	Update(id uuid.UUID, data *UpdateInvitationRequest) (*invitationDomain.Invitation, error)
	Delete(id uuid.UUID) error
	Publish(id uuid.UUID) error
	Unpublish(id uuid.UUID) error

	// Gallery
	AddGalleryItem(invitationID uuid.UUID, imageURL string, caption *string) (*invitationDomain.Gallery, error)
	GetGallery(invitationID uuid.UUID) ([]*invitationDomain.Gallery, error)
	DeleteGalleryItem(id uuid.UUID) error
	ReorderGallery(items []GalleryOrderItem) error

	// RSVP
	GetRSVPs(invitationID uuid.UUID) ([]*invitationDomain.RSVP, error)
	GetRSVPStats(invitationID uuid.UUID) (*invitationRepo.RSVPStats, error)
	ExportRSVPs(invitationID uuid.UUID) ([]byte, error)

	// Gift Accounts
	AddGiftAccount(invitationID uuid.UUID, bankName, accountNumber, accountName string, qrImage *string) (*invitationDomain.GiftAccount, error)
	GetGiftAccounts(invitationID uuid.UUID) ([]*invitationDomain.GiftAccount, error)
	DeleteGiftAccount(id uuid.UUID) error

	// Guest Messages
	GetGuestMessages(invitationID uuid.UUID, approvedOnly bool) ([]*invitationDomain.GuestMessage, error)
	ApproveMessage(id uuid.UUID) error
	DeleteMessage(id uuid.UUID) error

	// Guests
	AddGuest(invitationID uuid.UUID, name, guestType string) (*invitationDomain.Guest, error)
	GetGuests(invitationID uuid.UUID) ([]*invitationDomain.Guest, error)
	DeleteGuest(id uuid.UUID) error
	ImportGuests(invitationID uuid.UUID, file io.Reader) (int, error)
	ExportGuests(invitationID uuid.UUID) (*bytes.Buffer, string, error)
	GetGuestTemplate() (*bytes.Buffer, string)
}

// UpdateInvitationRequest contains updateable invitation fields
type UpdateInvitationRequest struct {
	// Theme
	ThemeID *uuid.UUID `json:"theme_id"`

	// Groom Info
	GroomName     *string `json:"groom_name"`
	GroomNickname *string `json:"groom_nickname"`
	GroomFather   *string `json:"groom_father"`
	GroomMother   *string `json:"groom_mother"`
	GroomPhoto    *string `json:"groom_photo"`

	// Bride Info
	BrideName     *string `json:"bride_name"`
	BrideNickname *string `json:"bride_nickname"`
	BrideFather   *string `json:"bride_father"`
	BrideMother   *string `json:"bride_mother"`
	BridePhoto    *string `json:"bride_photo"`

	// Couple
	CouplePhoto *string `json:"couple_photo"`
	LoveStory   *string `json:"love_story"`

	// Akad Event
	AkadDate     *time.Time `json:"akad_date"`
	AkadTime     *string    `json:"akad_time"`
	AkadLocation *string    `json:"akad_location"`
	AkadAddress  *string    `json:"akad_address"`
	AkadMapsURL  *string    `json:"akad_maps_url"`

	// Reception Event
	ReceptionDate     *time.Time `json:"reception_date"`
	ReceptionTime     *string    `json:"reception_time"`
	ReceptionLocation *string    `json:"reception_location"`
	ReceptionAddress  *string    `json:"reception_address"`
	ReceptionMapsURL  *string    `json:"reception_maps_url"`

	// Settings
	CoverImage       *string `json:"cover_image"`
	BackgroundImage  *string `json:"background_image"`
	MusicURL         *string `json:"music_url"`
	CountdownEnabled *bool   `json:"countdown_enabled"`
	RSVPEnabled      *bool   `json:"rsvp_enabled"`
	GiftEnabled      *bool   `json:"gift_enabled"`
	GuestBookEnabled *bool   `json:"guest_book_enabled"`
	GalleryEnabled   *bool   `json:"gallery_enabled"`

	// SEO
	MetaTitle       *string `json:"meta_title"`
	MetaDescription *string `json:"meta_description"`
	OGImage         *string `json:"og_image"`

	// Status
	IsPublished *bool      `json:"is_published"`
	ExpiresAt   *time.Time `json:"expires_at"`

	// Custom Content
	// Custom Content
	CustomCSS   *string                `json:"custom_css"`
	CustomHTML  *string                `json:"custom_html"`
	ContentData map[string]interface{} `json:"content_data"`
	Barcode     *string                `json:"barcode"`
	Video       *string                `json:"video"`
}

// GalleryOrderItem for reordering gallery
type GalleryOrderItem struct {
	ID    uuid.UUID `json:"id"`
	Order int       `json:"order"`
}

// service implements Service interface
type service struct {
	invitationRepo invitationRepo.Repository
	tenantRepo     tenantRepo.Repository
	cfg            *config.Config
}

// NewService creates a new invitation service
func NewService(invitationRepo invitationRepo.Repository, tenantRepo tenantRepo.Repository, cfg *config.Config) Service {
	return &service{
		invitationRepo: invitationRepo,
		tenantRepo:     tenantRepo,
		cfg:            cfg,
	}
}

// Create creates a new invitation for a tenant
func (s *service) Create(tenantID uuid.UUID) (*invitationDomain.Invitation, error) {
	// Check if invitation already exists for tenant
	existing, _ := s.invitationRepo.GetByTenantID(tenantID)
	if existing != nil {
		return nil, errors.New("invitation already exists for this tenant")
	}

	invitation := invitationDomain.NewInvitation(tenantID)
	if err := s.invitationRepo.Create(invitation); err != nil {
		return nil, fmt.Errorf("failed to create invitation: %w", err)
	}

	return invitation, nil
}

// GetByID retrieves an invitation by ID
func (s *service) GetByID(id uuid.UUID) (*invitationDomain.Invitation, error) {
	return s.invitationRepo.GetByID(id)
}

// GetByTenantID retrieves an invitation by tenant ID
func (s *service) GetByTenantID(tenantID uuid.UUID) (*invitationDomain.Invitation, error) {
	return s.invitationRepo.GetByTenantID(tenantID)
}

// Update updates an invitation
func (s *service) Update(id uuid.UUID, data *UpdateInvitationRequest) (*invitationDomain.Invitation, error) {
	invitation, err := s.invitationRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if data.ThemeID != nil {
		invitation.ThemeID = *data.ThemeID
	}
	if data.GroomName != nil {
		invitation.GroomName = *data.GroomName
	}
	if data.GroomNickname != nil {
		invitation.GroomNickname = *data.GroomNickname
	}
	if data.GroomFather != nil {
		invitation.GroomFather = *data.GroomFather
	}
	if data.GroomMother != nil {
		invitation.GroomMother = *data.GroomMother
	}
	if data.GroomPhoto != nil {
		invitation.GroomPhoto = data.GroomPhoto
	}
	if data.BrideName != nil {
		invitation.BrideName = *data.BrideName
	}
	if data.BrideNickname != nil {
		invitation.BrideNickname = *data.BrideNickname
	}
	if data.BrideFather != nil {
		invitation.BrideFather = *data.BrideFather
	}
	if data.BrideMother != nil {
		invitation.BrideMother = *data.BrideMother
	}
	if data.BridePhoto != nil {
		invitation.BridePhoto = data.BridePhoto
	}
	if data.CouplePhoto != nil {
		invitation.CouplePhoto = data.CouplePhoto
	}
	if data.LoveStory != nil {
		invitation.LoveStory = data.LoveStory
	}
	if data.AkadDate != nil {
		invitation.AkadDate = data.AkadDate
	}
	if data.AkadTime != nil {
		invitation.AkadTime = data.AkadTime
	}
	if data.AkadLocation != nil {
		invitation.AkadLocation = data.AkadLocation
	}
	if data.AkadAddress != nil {
		invitation.AkadAddress = data.AkadAddress
	}
	if data.AkadMapsURL != nil {
		invitation.AkadMapsURL = data.AkadMapsURL
	}
	if data.ReceptionDate != nil {
		invitation.ReceptionDate = data.ReceptionDate
	}
	if data.ReceptionTime != nil {
		invitation.ReceptionTime = data.ReceptionTime
	}
	if data.ReceptionLocation != nil {
		invitation.ReceptionLocation = data.ReceptionLocation
	}
	if data.ReceptionAddress != nil {
		invitation.ReceptionAddress = data.ReceptionAddress
	}
	if data.ReceptionMapsURL != nil {
		invitation.ReceptionMapsURL = data.ReceptionMapsURL
	}
	if data.CoverImage != nil {
		invitation.CoverImage = data.CoverImage
	}
	if data.BackgroundImage != nil {
		invitation.BackgroundImage = data.BackgroundImage
	}
	if data.MusicURL != nil {
		invitation.MusicURL = data.MusicURL
	}
	if data.CountdownEnabled != nil {
		invitation.CountdownEnabled = *data.CountdownEnabled
	}
	if data.RSVPEnabled != nil {
		invitation.RSVPEnabled = *data.RSVPEnabled
	}
	if data.GiftEnabled != nil {
		invitation.GiftEnabled = *data.GiftEnabled
	}
	if data.GuestBookEnabled != nil {
		invitation.GuestBookEnabled = *data.GuestBookEnabled
	}
	if data.GalleryEnabled != nil {
		invitation.GalleryEnabled = *data.GalleryEnabled
	}
	if data.MetaTitle != nil {
		invitation.MetaTitle = data.MetaTitle
	}
	if data.MetaDescription != nil {
		invitation.MetaDescription = data.MetaDescription
	}
	if data.OGImage != nil {
		invitation.OGImage = data.OGImage
	}
	if data.ExpiresAt != nil {
		invitation.ExpiresAt = data.ExpiresAt
	}
	if data.CustomCSS != nil {
		invitation.CustomCSS = data.CustomCSS
	}
	if data.CustomHTML != nil {
		invitation.CustomHTML = data.CustomHTML
	}
	if data.Barcode != nil {
		invitation.Barcode = *data.Barcode
	}
	if data.Video != nil {
		invitation.Video = *data.Video
	}
	if data.ContentData != nil {
		// Merge or Replace? For simplicity, we just set the new map (assuming frontend sends full map)
		// Or creating a new map if nil
		if invitation.ContentData == nil {
			invitation.ContentData = make(map[string]interface{})
		}
		for k, v := range data.ContentData {
			invitation.ContentData[k] = v
		}
	}

	if err := s.invitationRepo.Update(invitation); err != nil {
		return nil, fmt.Errorf("failed to update invitation: %w", err)
	}

	return invitation, nil
}

// Delete deletes an invitation
func (s *service) Delete(id uuid.UUID) error {
	return s.invitationRepo.Delete(id)
}

// Publish publishes an invitation
func (s *service) Publish(id uuid.UUID) error {
	invitation, err := s.invitationRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Validate invitation has minimum required data
	if invitation.GroomName == "" || invitation.BrideName == "" {
		return errors.New("invitation must have groom and bride names before publishing")
	}

	invitation.Publish()
	return s.invitationRepo.Update(invitation)
}

// Unpublish unpublishes an invitation
func (s *service) Unpublish(id uuid.UUID) error {
	invitation, err := s.invitationRepo.GetByID(id)
	if err != nil {
		return err
	}

	invitation.Unpublish()
	return s.invitationRepo.Update(invitation)
}

// AddGalleryItem adds an image to gallery
func (s *service) AddGalleryItem(invitationID uuid.UUID, imageURL string, caption *string) (*invitationDomain.Gallery, error) {
	gallery := invitationDomain.NewGallery(invitationID, imageURL)
	if caption != nil {
		gallery.Caption = caption
	}

	if err := s.invitationRepo.CreateGalleryItem(gallery); err != nil {
		return nil, fmt.Errorf("failed to add gallery item: %w", err)
	}

	return gallery, nil
}

// GetGallery retrieves gallery for an invitation
func (s *service) GetGallery(invitationID uuid.UUID) ([]*invitationDomain.Gallery, error) {
	return s.invitationRepo.GetGalleryByInvitationID(invitationID)
}

// DeleteGalleryItem deletes a gallery item
func (s *service) DeleteGalleryItem(id uuid.UUID) error {
	return s.invitationRepo.DeleteGalleryItem(id)
}

// ReorderGallery reorders gallery items
func (s *service) ReorderGallery(items []GalleryOrderItem) error {
	for _, item := range items {
		if err := s.invitationRepo.UpdateGalleryOrder(item.ID, item.Order); err != nil {
			return err
		}
	}
	return nil
}

// GetRSVPs retrieves RSVPs for an invitation
func (s *service) GetRSVPs(invitationID uuid.UUID) ([]*invitationDomain.RSVP, error) {
	return s.invitationRepo.GetRSVPsByInvitationID(invitationID)
}

// GetRSVPStats returns RSVP statistics
func (s *service) GetRSVPStats(invitationID uuid.UUID) (*invitationRepo.RSVPStats, error) {
	return s.invitationRepo.GetRSVPStats(invitationID)
}

// ExportRSVPs exports RSVPs to CSV
func (s *service) ExportRSVPs(invitationID uuid.UUID) ([]byte, error) {
	rsvps, err := s.invitationRepo.GetRSVPsByInvitationID(invitationID)
	if err != nil {
		return nil, err
	}

	// Build CSV
	csv := "Nama,Kehadiran,Jumlah Tamu,Pesan,Tanggal\n"
	for _, rsvp := range rsvps {
		message := ""
		if rsvp.Message != nil {
			message = *rsvp.Message
		}
		csv += fmt.Sprintf("%s,%s,%d,%s,%s\n",
			rsvp.GuestName,
			rsvp.Attendance,
			rsvp.GuestCount,
			message,
			rsvp.CreatedAt.Format("2006-01-02 15:04:05"),
		)
	}

	return []byte(csv), nil
}

// AddGiftAccount adds a gift account
func (s *service) AddGiftAccount(invitationID uuid.UUID, bankName, accountNumber, accountName string, qrImage *string) (*invitationDomain.GiftAccount, error) {
	account := invitationDomain.NewGiftAccount(invitationID, bankName, accountNumber, accountName)
	if qrImage != nil {
		account.QRImage = qrImage
	}

	if err := s.invitationRepo.CreateGiftAccount(account); err != nil {
		return nil, fmt.Errorf("failed to add gift account: %w", err)
	}

	return account, nil
}

// GetGiftAccounts retrieves gift accounts
func (s *service) GetGiftAccounts(invitationID uuid.UUID) ([]*invitationDomain.GiftAccount, error) {
	return s.invitationRepo.GetGiftAccountsByInvitationID(invitationID)
}

// DeleteGiftAccount deletes a gift account
func (s *service) DeleteGiftAccount(id uuid.UUID) error {
	return s.invitationRepo.DeleteGiftAccount(id)
}

// GetGuestMessages retrieves guest messages
func (s *service) GetGuestMessages(invitationID uuid.UUID, approvedOnly bool) ([]*invitationDomain.GuestMessage, error) {
	return s.invitationRepo.GetGuestMessagesByInvitationID(invitationID, approvedOnly)
}

// ApproveMessage approves a guest message
func (s *service) ApproveMessage(id uuid.UUID) error {
	return s.invitationRepo.ApproveGuestMessage(id)
}

// DeleteMessage deletes a guest message
func (s *service) DeleteMessage(id uuid.UUID) error {
	return s.invitationRepo.DeleteGuestMessage(id)
}
