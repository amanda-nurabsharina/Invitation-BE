package invitation

import (
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Invitation represents a wedding invitation
type Invitation struct {
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;column:id" json:"id"`
	TenantID uuid.UUID `gorm:"type:uuid;column:tenant_id;not null;index" json:"tenant_id"`
	ThemeID  uuid.UUID `gorm:"type:uuid;column:theme_id" json:"theme_id"`

	// Groom Information
	GroomName     string  `gorm:"column:groom_name" json:"groom_name"`
	GroomNickname string  `gorm:"column:groom_nickname" json:"groom_nickname"`
	GroomFather   string  `gorm:"column:groom_father" json:"groom_father"`
	GroomMother   string  `gorm:"column:groom_mother" json:"groom_mother"`
	GroomPhoto    *string `gorm:"column:groom_photo" json:"groom_photo,omitempty"`

	// Bride Information
	BrideName     string  `gorm:"column:bride_name" json:"bride_name"`
	BrideNickname string  `gorm:"column:bride_nickname" json:"bride_nickname"`
	BrideFather   string  `gorm:"column:bride_father" json:"bride_father"`
	BrideMother   string  `gorm:"column:bride_mother" json:"bride_mother"`
	BridePhoto    *string `gorm:"column:bride_photo" json:"bride_photo,omitempty"`

	// Couple
	CouplePhoto *string `gorm:"column:couple_photo" json:"couple_photo,omitempty"`
	LoveStory   *string `gorm:"column:love_story" json:"love_story,omitempty"`

	// Akad Event
	AkadDate     *time.Time `gorm:"column:akad_date" json:"akad_date,omitempty"`
	AkadTime     *string    `gorm:"column:akad_time" json:"akad_time,omitempty"`
	AkadLocation *string    `gorm:"column:akad_location" json:"akad_location,omitempty"`
	AkadAddress  *string    `gorm:"column:akad_address" json:"akad_address,omitempty"`
	AkadMapsURL  *string    `gorm:"column:akad_maps_url" json:"akad_maps_url,omitempty"`

	// Reception Event
	ReceptionDate     *time.Time `gorm:"column:reception_date" json:"reception_date,omitempty"`
	ReceptionTime     *string    `gorm:"column:reception_time" json:"reception_time,omitempty"`
	ReceptionLocation *string    `gorm:"column:reception_location" json:"reception_location,omitempty"`
	ReceptionAddress  *string    `gorm:"column:reception_address" json:"reception_address,omitempty"`
	ReceptionMapsURL  *string    `gorm:"column:reception_maps_url" json:"reception_maps_url,omitempty"`

	// Additional Events (stored as JSON)
	AdditionalEvents []AdditionalEvent `gorm:"serializer:json;column:additional_events" json:"additional_events,omitempty"`

	// Settings
	CoverImage       *string `gorm:"column:cover_image" json:"cover_image,omitempty"`
	BackgroundImage  *string `gorm:"column:background_image" json:"background_image,omitempty"`
	MusicURL         *string `gorm:"column:music_url" json:"music_url,omitempty"`
	CountdownEnabled bool    `gorm:"column:countdown_enabled;default:true" json:"countdown_enabled"`
	RSVPEnabled      bool    `gorm:"column:rsvp_enabled;default:true" json:"rsvp_enabled"`
	GiftEnabled      bool    `gorm:"column:gift_enabled;default:false" json:"gift_enabled"`
	GuestBookEnabled bool    `gorm:"column:guest_book_enabled;default:true" json:"guest_book_enabled"`
	GalleryEnabled   bool    `gorm:"column:gallery_enabled;default:true" json:"gallery_enabled"`

	// Media
	Barcode string `gorm:"column:barcode" json:"barcode"` // URL to barcode image
	Video   string `gorm:"column:video" json:"video"`     // URL to video file or embed

	// SEO & Styling
	MetaTitle       *string `gorm:"column:meta_title" json:"meta_title,omitempty"`
	MetaDescription *string `gorm:"column:meta_description" json:"meta_description,omitempty"`
	OGImage         *string `gorm:"column:og_image" json:"og_image,omitempty"`
	CustomCSS       *string `gorm:"column:custom_css;type:text" json:"custom_css,omitempty"`
	CustomHTML      *string `gorm:"column:custom_html;type:text" json:"custom_html,omitempty"`

	// Dynamic Content (Images, Texts mapped from Template)
	ContentData map[string]interface{} `gorm:"serializer:json;column:content_data" json:"content_data,omitempty"`

	// Status
	IsPublished bool       `gorm:"column:is_published;default:false" json:"is_published"`
	PublishedAt *time.Time `gorm:"column:published_at" json:"published_at,omitempty"`
	ExpiresAt   *time.Time `gorm:"column:expires_at" json:"expires_at,omitempty"`

	// Timestamps
	CreatedAt time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`

	// Relationships
	Guests       []Guest        `gorm:"foreignKey:InvitationID" json:"guests,omitempty"`
	Gallery      []Gallery      `gorm:"foreignKey:InvitationID" json:"gallery,omitempty"`
	RSVPs        []RSVP         `gorm:"foreignKey:InvitationID" json:"rsvps,omitempty"`
	GiftAccounts []GiftAccount  `gorm:"foreignKey:InvitationID" json:"gift_accounts,omitempty"`
	Messages     []GuestMessage `gorm:"foreignKey:InvitationID" json:"messages,omitempty"`
}

// AdditionalEvent for custom events
type AdditionalEvent struct {
	Name     string `json:"name"`
	Date     string `json:"date"`
	Time     string `json:"time"`
	Location string `json:"location"`
	Address  string `json:"address"`
	MapsURL  string `json:"maps_url,omitempty"`
}

// NewInvitation creates a new invitation instance
func NewInvitation(tenantID uuid.UUID) *Invitation {
	return &Invitation{
		ID:               uuid.New(),
		TenantID:         tenantID,
		CountdownEnabled: true,
		RSVPEnabled:      true,
		GiftEnabled:      false,
		GuestBookEnabled: true,
		GalleryEnabled:   true,
		IsPublished:      false,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}
}

// Validate performs business rule validation
func (i *Invitation) Validate() error {
	if i.TenantID == uuid.Nil {
		return errors.New("tenant ID is required")
	}

	if strings.TrimSpace(i.GroomName) == "" && strings.TrimSpace(i.BrideName) == "" {
		return errors.New("at least one name (groom or bride) is required")
	}

	return nil
}

// Publish publishes the invitation
func (i *Invitation) Publish() {
	now := time.Now()
	i.IsPublished = true
	i.PublishedAt = &now
	i.UpdatedAt = now
}

// Unpublish unpublishes the invitation
func (i *Invitation) Unpublish() {
	i.IsPublished = false
	i.PublishedAt = nil
	i.UpdatedAt = time.Now()
}

// IsExpired checks if the invitation has expired
func (i *Invitation) IsExpired() bool {
	if i.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*i.ExpiresAt)
}

// TableName returns the table name
func (Invitation) TableName() string {
	return "invitation"
}

// BeforeCreate GORM hook
func (i *Invitation) BeforeCreate(tx *gorm.DB) error {
	i.CreatedAt = time.Now()
	i.UpdatedAt = time.Now()
	return nil
}

// BeforeUpdate GORM hook
func (i *Invitation) BeforeUpdate(tx *gorm.DB) error {
	i.UpdatedAt = time.Now()
	return nil
}
