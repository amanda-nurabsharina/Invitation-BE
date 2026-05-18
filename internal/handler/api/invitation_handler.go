package api

import (
	"fmt"
	"time"

	invitationRepo "invitation-api/internal/repository/invitation"
	tenantRepo "invitation-api/internal/repository/tenant"
	invitationService "invitation-api/internal/service/invitation"
	"invitation-api/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// InvitationHandler handles invitation API requests
type InvitationHandler struct {
	invitationService invitationService.Service
}

// NewInvitationHandler creates a new invitation handler
func NewInvitationHandler(cfg *config.Config) *InvitationHandler {
	invitationRepository := invitationRepo.NewRepository()
	tenantRepository := tenantRepo.NewRepository()
	invitationSvc := invitationService.NewService(invitationRepository, tenantRepository, cfg)

	return &InvitationHandler{
		invitationService: invitationSvc,
	}
}

// CreateInvitationRequest represents create invitation request
type CreateInvitationRequest struct {
	TenantID uuid.UUID `json:"tenant_id" validate:"required"`
}

// UpdateInvitationRequest represents update invitation request
type UpdateInvitationRequest struct {
	ThemeID           *string                `json:"theme_id"`
	GroomName         *string                `json:"groom_name"`
	GroomNickname     *string                `json:"groom_nickname"`
	GroomFather       *string                `json:"groom_father"`
	GroomMother       *string                `json:"groom_mother"`
	GroomPhoto        *string                `json:"groom_photo"`
	BrideName         *string                `json:"bride_name"`
	BrideNickname     *string                `json:"bride_nickname"`
	BrideFather       *string                `json:"bride_father"`
	BrideMother       *string                `json:"bride_mother"`
	BridePhoto        *string                `json:"bride_photo"`
	CouplePhoto       *string                `json:"couple_photo"`
	LoveStory         *string                `json:"love_story"`
	AkadDate          *string                `json:"akad_date"`
	AkadTime          *string                `json:"akad_time"`
	AkadLocation      *string                `json:"akad_location"`
	AkadAddress       *string                `json:"akad_address"`
	AkadMapsURL       *string                `json:"akad_maps_url"`
	ReceptionDate     *string                `json:"reception_date"`
	ReceptionTime     *string                `json:"reception_time"`
	ReceptionLocation *string                `json:"reception_location"`
	ReceptionAddress  *string                `json:"reception_address"`
	ReceptionMapsURL  *string                `json:"reception_maps_url"`
	MusicURL          *string                `json:"music_url"`
	CountdownEnabled  *bool                  `json:"countdown_enabled"`
	RSVPEnabled       *bool                  `json:"rsvp_enabled"`
	GiftEnabled       *bool                  `json:"gift_enabled"`
	GuestBookEnabled  *bool                  `json:"guest_book_enabled"`
	GalleryEnabled    *bool                  `json:"gallery_enabled"`
	MetaTitle         *string                `json:"meta_title"`
	MetaDescription   *string                `json:"meta_description"`
	OGImage           *string                `json:"og_image"`
	CustomCSS         *string                `json:"custom_css"`
	CustomHTML        *string                `json:"custom_html"`
	ContentData       map[string]interface{} `json:"content_data"`
	ExpiresAt         *string                `json:"expires_at"` // RFC3339 format
	Barcode           *string                `json:"barcode"`
	Video             *string                `json:"video"`
}

// Create creates a new invitation
func (h *InvitationHandler) Create(c *fiber.Ctx) error {
	var req CreateInvitationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	invitation, err := h.invitationService.Create(req.TenantID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Invitation created successfully",
		"data":    invitation,
	})
}

// GetByID retrieves an invitation by ID
func (h *InvitationHandler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	invitation, err := h.invitationService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Invitation not found",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  invitation,
	})
}

// GetByTenantID retrieves an invitation by tenant ID
func (h *InvitationHandler) GetByTenantID(c *fiber.Ctx) error {
	tenantIDStr := c.Params("tenant_id")
	tenantID, err := uuid.Parse(tenantIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid tenant ID",
		})
	}

	invitation, err := h.invitationService.GetByTenantID(tenantID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Invitation not found",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  invitation,
	})
}

// Update updates an invitation
func (h *InvitationHandler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	var req UpdateInvitationRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	// Convert to service request
	updateReq := &invitationService.UpdateInvitationRequest{}

	if req.ThemeID != nil {
		themeID, _ := uuid.Parse(*req.ThemeID)
		updateReq.ThemeID = &themeID
	}
	updateReq.GroomName = req.GroomName
	updateReq.GroomNickname = req.GroomNickname
	updateReq.GroomFather = req.GroomFather
	updateReq.GroomMother = req.GroomMother
	updateReq.GroomPhoto = req.GroomPhoto
	updateReq.BrideName = req.BrideName
	updateReq.BrideNickname = req.BrideNickname
	updateReq.BrideFather = req.BrideFather
	updateReq.BrideMother = req.BrideMother
	updateReq.BridePhoto = req.BridePhoto
	updateReq.CouplePhoto = req.CouplePhoto
	updateReq.LoveStory = req.LoveStory
	updateReq.AkadTime = req.AkadTime
	updateReq.AkadLocation = req.AkadLocation
	updateReq.AkadAddress = req.AkadAddress
	updateReq.AkadMapsURL = req.AkadMapsURL
	updateReq.ReceptionTime = req.ReceptionTime
	updateReq.ReceptionLocation = req.ReceptionLocation
	updateReq.ReceptionAddress = req.ReceptionAddress
	updateReq.ReceptionMapsURL = req.ReceptionMapsURL
	updateReq.MusicURL = req.MusicURL
	updateReq.CountdownEnabled = req.CountdownEnabled
	updateReq.RSVPEnabled = req.RSVPEnabled
	updateReq.GiftEnabled = req.GiftEnabled
	updateReq.GuestBookEnabled = req.GuestBookEnabled
	updateReq.GalleryEnabled = req.GalleryEnabled
	updateReq.MetaTitle = req.MetaTitle
	updateReq.MetaDescription = req.MetaDescription
	updateReq.OGImage = req.OGImage
	updateReq.CustomCSS = req.CustomCSS
	updateReq.CustomHTML = req.CustomHTML
	updateReq.ContentData = req.ContentData
	updateReq.Barcode = req.Barcode
	updateReq.Video = req.Video

	// Parse dates
	if req.AkadDate != nil {
		akadDate, err := time.Parse("2006-01-02", *req.AkadDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid Akad Date format: " + err.Error(),
			})
		}
		updateReq.AkadDate = &akadDate
	}
	if req.ReceptionDate != nil {
		receptionDate, err := time.Parse("2006-01-02", *req.ReceptionDate)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid Reception Date format: " + err.Error(),
			})
		}
		updateReq.ReceptionDate = &receptionDate
	}
	if req.ExpiresAt != nil {
		// Try parsing RFC3339 which covers ISO8601
		expiresAt, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			fmt.Println("Error parsing expires_at:", err) // Debug log
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   true,
				"message": "Invalid Expires At format: " + err.Error(),
			})
		}
		updateReq.ExpiresAt = &expiresAt
	}

	invitation, err := h.invitationService.Update(id, updateReq)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Invitation updated successfully",
		"data":    invitation,
	})
}

// Delete deletes an invitation
func (h *InvitationHandler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	if err := h.invitationService.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete invitation",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Invitation deleted successfully",
	})
}

// Publish publishes an invitation
func (h *InvitationHandler) Publish(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	if err := h.invitationService.Publish(id); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Invitation published successfully",
	})
}

// Unpublish unpublishes an invitation
func (h *InvitationHandler) Unpublish(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	if err := h.invitationService.Unpublish(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to unpublish invitation",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Invitation unpublished successfully",
	})
}

// GetRSVPs retrieves RSVPs for an invitation
func (h *InvitationHandler) GetRSVPs(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	rsvps, err := h.invitationService.GetRSVPs(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve RSVPs",
		})
	}

	stats, _ := h.invitationService.GetRSVPStats(id)

	return c.JSON(fiber.Map{
		"error": false,
		"data":  rsvps,
		"stats": stats,
	})
}

// ExportRSVPs exports RSVPs to CSV
func (h *InvitationHandler) ExportRSVPs(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	csv, err := h.invitationService.ExportRSVPs(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to export RSVPs",
		})
	}

	c.Set("Content-Type", "text/csv")
	c.Set("Content-Disposition", "attachment; filename=rsvp_export.csv")
	return c.Send(csv)
}

// GetGuestMessages retrieves guest messages
func (h *InvitationHandler) GetGuestMessages(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	approvedOnly := c.Query("approved_only", "false") == "true"
	messages, err := h.invitationService.GetGuestMessages(id, approvedOnly)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve messages",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  messages,
	})
}

// ApproveMessage approves a guest message
func (h *InvitationHandler) ApproveMessage(c *fiber.Ctx) error {
	idStr := c.Params("message_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid message ID",
		})
	}

	if err := h.invitationService.ApproveMessage(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to approve message",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Message approved successfully",
	})
}

// DeleteMessage deletes a guest message
func (h *InvitationHandler) DeleteMessage(c *fiber.Ctx) error {
	idStr := c.Params("message_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid message ID",
		})
	}

	if err := h.invitationService.DeleteMessage(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete message",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Message deleted successfully",
	})
}

// AddGalleryItemRequest represents add gallery item request
type AddGalleryItemRequest struct {
	ImageURL string  `json:"image_url" validate:"required"`
	Caption  *string `json:"caption"`
}

// AddGalleryItem adds an image to gallery
func (h *InvitationHandler) AddGalleryItem(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	var req AddGalleryItemRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	gallery, err := h.invitationService.AddGalleryItem(id, req.ImageURL, req.Caption)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Gallery item added successfully",
		"data":    gallery,
	})
}

// GetGallery retrieves gallery for an invitation
func (h *InvitationHandler) GetGallery(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	gallery, err := h.invitationService.GetGallery(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve gallery",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  gallery,
	})
}

// DeleteGalleryItem deletes a gallery item
func (h *InvitationHandler) DeleteGalleryItem(c *fiber.Ctx) error {
	idStr := c.Params("gallery_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid gallery ID",
		})
	}

	if err := h.invitationService.DeleteGalleryItem(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete gallery item",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Gallery item deleted successfully",
	})
}

// ReorderGalleryRequest represents gallery reorder request
type ReorderGalleryRequest struct {
	Items []invitationService.GalleryOrderItem `json:"items" validate:"required"`
}

// ReorderGallery reorders gallery items
func (h *InvitationHandler) ReorderGallery(c *fiber.Ctx) error {
	var req ReorderGalleryRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	if err := h.invitationService.ReorderGallery(req.Items); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to reorder gallery",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Gallery reordered successfully",
	})
}

// AddGiftAccountRequest represents add gift account request
type AddGiftAccountRequest struct {
	BankName      string  `json:"bank_name" validate:"required"`
	AccountNumber string  `json:"account_number" validate:"required"`
	AccountName   string  `json:"account_name" validate:"required"`
	QRImage       *string `json:"qr_image"`
}

// AddGiftAccount adds a gift account
func (h *InvitationHandler) AddGiftAccount(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	var req AddGiftAccountRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	account, err := h.invitationService.AddGiftAccount(id, req.BankName, req.AccountNumber, req.AccountName, req.QRImage)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Gift account added successfully",
		"data":    account,
	})
}

// GetGiftAccounts retrieves gift accounts
func (h *InvitationHandler) GetGiftAccounts(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	accounts, err := h.invitationService.GetGiftAccounts(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to retrieve gift accounts",
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  accounts,
	})
}

// DeleteGiftAccount deletes a gift account
func (h *InvitationHandler) DeleteGiftAccount(c *fiber.Ctx) error {
	idStr := c.Params("gift_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid gift account ID",
		})
	}

	if err := h.invitationService.DeleteGiftAccount(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete gift account",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Gift account deleted successfully",
	})
}

// GetDashboard returns dashboard statistics
func (h *InvitationHandler) GetDashboard(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	invitation, err := h.invitationService.GetByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error":   true,
			"message": "Invitation not found",
		})
	}

	rsvpStats, _ := h.invitationService.GetRSVPStats(id)
	gallery, _ := h.invitationService.GetGallery(id)
	messages, _ := h.invitationService.GetGuestMessages(id, false)

	return c.JSON(fiber.Map{
		"error": false,
		"data": fiber.Map{
			"invitation":    invitation,
			"rsvp_stats":    rsvpStats,
			"gallery_count": len(gallery),
			"message_count": len(messages),
		},
	})
}

// AddGuestRequest represents add guest request
type AddGuestRequest struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type"`
}

// AddGuest adds a new guest
func (h *InvitationHandler) AddGuest(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	var req AddGuestRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid request body",
		})
	}

	guest, err := h.invitationService.AddGuest(id, req.Name, req.Type)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"error":   false,
		"message": "Guest added successfully",
		"data":    guest,
	})
}

// GetGuests retrieves guests
func (h *InvitationHandler) GetGuests(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	guests, err := h.invitationService.GetGuests(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error": false,
		"data":  guests,
	})
}

// DeleteGuest deletes a guest
func (h *InvitationHandler) DeleteGuest(c *fiber.Ctx) error {
	idStr := c.Params("guest_id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid guest ID",
		})
	}

	if err := h.invitationService.DeleteGuest(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Guest deleted successfully",
	})
}

// ImportGuests imports guests from Excel
func (h *InvitationHandler) ImportGuests(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "File is required",
		})
	}

	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}
	defer src.Close()

	count, err := h.invitationService.ImportGuests(id, src)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": fmt.Sprintf("Imported %d guests", count),
	})
}

// ExportGuests exports guest list with URLs
func (h *InvitationHandler) ExportGuests(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid invitation ID",
		})
	}

	buf, filename, err := h.invitationService.ExportGuests(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": err.Error(),
		})
	}

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.SendStream(buf)
}

// GetGuestTemplate returns Excel template
func (h *InvitationHandler) GetGuestTemplate(c *fiber.Ctx) error {
	buf, filename := h.invitationService.GetGuestTemplate()

	c.Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	return c.SendStream(buf)
}
