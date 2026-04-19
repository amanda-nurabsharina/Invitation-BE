package htmx

import (
	"html/template"
	"log"
	"path/filepath"
	"strconv"

	invitationDomain "invitation-api/internal/domain/invitation"
	invitationRepo "invitation-api/internal/repository/invitation"
	tenantRepo "invitation-api/internal/repository/tenant"
	themeRepo "invitation-api/internal/repository/theme"
	"invitation-api/pkg/config"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// Handler handles htmx requests for client sites
type Handler struct {
	tenantRepo     tenantRepo.Repository
	invitationRepo invitationRepo.Repository
	themeRepo      themeRepo.Repository
	templates      *template.Template
	config         *config.Config
}

// NewHandler creates a new htmx handler
func NewHandler(cfg *config.Config) *Handler {
	// Load all templates
	tmpl, err := template.ParseGlob("internal/templates/**/*.html")
	if err != nil {
		log.Printf("Warning: Failed to load templates: %v", err)
		tmpl = template.New("")
	}

	return &Handler{
		tenantRepo:     tenantRepo.NewRepository(),
		invitationRepo: invitationRepo.NewRepository(),
		themeRepo:      themeRepo.NewRepository(),
		templates:      tmpl,
		config:         cfg,
	}
}

// RenderInvitation renders the full invitation page
func (h *Handler) RenderInvitation(c *fiber.Ctx) error {
	subdomain := c.Params("subdomain")
	if subdomain == "" {
		return c.Status(fiber.StatusNotFound).SendString("Invitation not found")
	}

	// Get tenant by subdomain
	tenant, err := h.tenantRepo.GetBySubdomain(subdomain)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Invitation not found")
	}

	// Check if tenant is accessible
	if !tenant.IsAccessible() {
		return c.Status(fiber.StatusForbidden).SendString("This invitation is not available")
	}

	// Get invitation for this tenant
	invitation, err := h.invitationRepo.GetByTenantID(tenant.ID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).SendString("Invitation not found")
	}

	// Check if published
	if !invitation.IsPublished {
		return c.Status(fiber.StatusForbidden).SendString("This invitation is not published yet")
	}

	// Check if expired
	if invitation.IsExpired() {
		return c.Status(fiber.StatusForbidden).SendString("This invitation has expired")
	}

	// Prepare template data
	data := h.prepareInvitationData(invitation, subdomain)
	data["ExpiresAt"] = invitation.ExpiresAt

	// Get guest name from query param
	guestName := c.Query("to")
	if guestName != "" {
		data["GuestName"] = guestName
	}

	// Determine theme template based on invitation's theme
	var tmpl *template.Template
	var errTmpl error

	if invitation.ThemeID != uuid.Nil {
		theme, err := h.themeRepo.GetByID(invitation.ThemeID)
		if err == nil && theme != nil {
			// Add theme CSS path to data
			data["ThemeCSS"] = "/assets/css/themes/" + theme.Slug + ".css"

			// Inject Theme Inline CSS (Global for theme)
			if theme.CustomCSS != nil && *theme.CustomCSS != "" {
				data["ThemeInlineCSS"] = template.CSS(*theme.CustomCSS)
			}

			// Priority 0: Check if Invitation has Custom HTML (Override Theme)
			if invitation.CustomHTML != nil && *invitation.CustomHTML != "" {
				tmpl, errTmpl = template.New("custom_" + invitation.ID.String()).Parse(*invitation.CustomHTML)
			} else {
				// Priority 1: Check if theme has Custom HTML in Database
				if theme.CustomHTML != nil && *theme.CustomHTML != "" {
					tmpl, errTmpl = template.New(theme.Slug).Parse(*theme.CustomHTML)
				} else {
					// Priority 2: Load from File System using TemplatePath
					templateFolder := theme.TemplatePath
					if templateFolder == "" {
						templateFolder = theme.Slug // Fallback to slug if TemplatePath is empty
					}
					themePath := filepath.Join("internal/templates/themes", templateFolder, "index.html")
					tmpl, errTmpl = template.ParseFiles(themePath)
				}
			}
		}
	}

	// Fallback to elegant theme if template is not loaded or error occurred
	if tmpl == nil || errTmpl != nil {
		if errTmpl != nil {
			log.Printf("Error loading theme template: %v", errTmpl)
		}

		if data["ThemeCSS"] == nil {
			data["ThemeCSS"] = "/assets/css/themes/elegant.css"
		}

		themePath := filepath.Join("internal/templates/themes/elegant/index.html")
		tmpl, err = template.ParseFiles(themePath)
		if err != nil {
			log.Printf("CRITICAL: Error parsing fallback template: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Error loading invitation system")
		}
	}

	// Render template
	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c.Response().BodyWriter(), data)
}

// RenderGallery renders gallery HTML fragment
func (h *Handler) RenderGallery(c *fiber.Ctx) error {
	subdomain := c.Params("subdomain")

	tenant, err := h.tenantRepo.GetBySubdomain(subdomain)
	if err != nil {
		return c.SendString("<p>Gallery not found</p>")
	}

	invitation, err := h.invitationRepo.GetByTenantID(tenant.ID)
	if err != nil {
		return c.SendString("<p>Gallery not found</p>")
	}

	gallery, err := h.invitationRepo.GetGalleryByInvitationID(invitation.ID)
	if err != nil {
		gallery = []*invitationDomain.Gallery{}
	}

	data := map[string]interface{}{
		"Gallery": gallery,
	}

	return h.renderComponent(c, "gallery.html", data)
}

// RenderMessages renders guest messages HTML fragment
func (h *Handler) RenderMessages(c *fiber.Ctx) error {
	subdomain := c.Params("subdomain")
	page := c.QueryInt("page", 1)
	limit := 10

	tenant, err := h.tenantRepo.GetBySubdomain(subdomain)
	if err != nil {
		return c.SendString("<p>Messages not found</p>")
	}

	invitation, err := h.invitationRepo.GetByTenantID(tenant.ID)
	if err != nil {
		return c.SendString("<p>Messages not found</p>")
	}

	// Get only approved messages
	messages, err := h.invitationRepo.GetGuestMessagesByInvitationID(invitation.ID, true)
	if err != nil {
		messages = []*invitationDomain.GuestMessage{}
	}

	// Paginate
	start := (page - 1) * limit
	end := start + limit
	if start > len(messages) {
		messages = []*invitationDomain.GuestMessage{}
	} else if end > len(messages) {
		messages = messages[start:]
	} else {
		messages = messages[start:end]
	}

	hasMore := end < len(messages)

	data := map[string]interface{}{
		"Messages":  messages,
		"Subdomain": subdomain,
		"HasMore":   hasMore,
		"NextPage":  page + 1,
	}

	return h.renderComponent(c, "messages.html", data)
}

// RenderGift renders gift accounts HTML fragment
func (h *Handler) RenderGift(c *fiber.Ctx) error {
	subdomain := c.Params("subdomain")

	tenant, err := h.tenantRepo.GetBySubdomain(subdomain)
	if err != nil {
		return c.SendString("<p>Gift info not found</p>")
	}

	invitation, err := h.invitationRepo.GetByTenantID(tenant.ID)
	if err != nil {
		return c.SendString("<p>Gift info not found</p>")
	}

	giftAccounts, err := h.invitationRepo.GetGiftAccountsByInvitationID(invitation.ID)
	if err != nil {
		giftAccounts = []*invitationDomain.GiftAccount{}
	}

	data := map[string]interface{}{
		"GiftAccounts": giftAccounts,
	}

	return h.renderComponent(c, "gift.html", data)
}

// SubmitRSVP handles RSVP form submission
func (h *Handler) SubmitRSVP(c *fiber.Ctx) error {
	subdomain := c.Params("subdomain")

	tenant, err := h.tenantRepo.GetBySubdomain(subdomain)
	if err != nil {
		return c.SendString("<div class='error'>Invitation not found</div>")
	}

	invitation, err := h.invitationRepo.GetByTenantID(tenant.ID)
	if err != nil {
		return c.SendString("<div class='error'>Invitation not found</div>")
	}

	// Parse form data
	guestName := c.FormValue("guest_name")
	attendance := c.FormValue("attendance")
	guestCountStr := c.FormValue("guest_count")
	message := c.FormValue("message")

	if guestName == "" || attendance == "" {
		return c.SendString("<div class='error' style='color: red; padding: 15px; background: #fef2f2; border-radius: 8px;'>Nama dan konfirmasi kehadiran wajib diisi</div>")
	}

	guestCount := 1
	if guestCountStr != "" {
		// Parse guest count
		if parsed, err := strconv.Atoi(guestCountStr); err == nil {
			guestCount = parsed
		}
	}

	// Create RSVP
	rsvp := invitationDomain.NewRSVP(invitation.ID, guestName, attendance, guestCount)
	if message != "" {
		rsvp.Message = &message

		// Also create a guest message
		guestMessage := invitationDomain.NewGuestMessage(invitation.ID, guestName, message)
		h.invitationRepo.CreateGuestMessage(guestMessage)
	}

	if err := h.invitationRepo.CreateRSVP(rsvp); err != nil {
		return c.SendString("<div class='error' style='color: red; padding: 15px; background: #fef2f2; border-radius: 8px;'>Gagal menyimpan konfirmasi. Silakan coba lagi.</div>")
	}

	return h.renderComponent(c, "rsvp_success.html", nil)
}

// Helper to render component templates
func (h *Handler) renderComponent(c *fiber.Ctx, name string, data interface{}) error {
	tmplPath := filepath.Join("internal/templates/components", name)
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		log.Printf("Error parsing component template %s: %v", name, err)
		return c.SendString("<p>Error loading content</p>")
	}

	c.Set("Content-Type", "text/html")
	return tmpl.Execute(c.Response().BodyWriter(), data)
}

// prepareInvitationData prepares data for invitation template
func (h *Handler) prepareInvitationData(inv *invitationDomain.Invitation, subdomain string) map[string]interface{} {
	data := map[string]interface{}{
		"Subdomain":        subdomain,
		"GroomName":        inv.GroomName,
		"GroomNickname":    inv.GroomNickname,
		"GroomFather":      inv.GroomFather,
		"GroomMother":      inv.GroomMother,
		"BrideName":        inv.BrideName,
		"BrideNickname":    inv.BrideNickname,
		"BrideFather":      inv.BrideFather,
		"BrideMother":      inv.BrideMother,
		"CountdownEnabled": inv.CountdownEnabled,
		"RSVPEnabled":      inv.RSVPEnabled,
		"GiftEnabled":      inv.GiftEnabled,
		"GuestBookEnabled": inv.GuestBookEnabled,
		"GalleryEnabled":   inv.GalleryEnabled,
	}

	// Optional fields
	if inv.GroomPhoto != nil {
		data["GroomPhoto"] = *inv.GroomPhoto
	}
	if inv.BridePhoto != nil {
		data["BridePhoto"] = *inv.BridePhoto
	}
	if inv.CouplePhoto != nil {
		data["CouplePhoto"] = *inv.CouplePhoto
	}
	if inv.MusicURL != nil {
		data["MusicURL"] = *inv.MusicURL
	}
	if inv.MetaTitle != nil {
		data["MetaTitle"] = *inv.MetaTitle
	} else {
		data["MetaTitle"] = inv.GroomNickname + " & " + inv.BrideNickname
	}
	if inv.MetaDescription != nil {
		data["MetaDescription"] = *inv.MetaDescription
	} else {
		data["MetaDescription"] = "Undangan pernikahan " + inv.GroomName + " & " + inv.BrideName
	}
	if inv.OGImage != nil {
		data["OGImage"] = *inv.OGImage
	}
	if inv.CustomCSS != nil {
		data["CustomCSS"] = template.CSS(*inv.CustomCSS)
	}
	if inv.CoverImage != nil {
		data["CoverImage"] = *inv.CoverImage
	}
	if inv.BackgroundImage != nil {
		data["BackgroundImage"] = *inv.BackgroundImage
	}
	if inv.LoveStory != nil {
		data["LoveStory"] = *inv.LoveStory
	}
	if inv.Barcode != "" {
		data["Barcode"] = inv.Barcode
	}
	if inv.Video != "" {
		data["Video"] = inv.Video
	}

	// Event dates
	if inv.ContentData != nil {
		data["ContentData"] = inv.ContentData
	}

	if inv.AkadDate != nil {
		data["AkadDate"] = inv.AkadDate
		data["AkadDateFormatted"] = inv.AkadDate.Format("Monday, 02 January 2006")
		data["WeddingDate"] = inv.AkadDate.Format("02.01.2006")
		data["WeddingDateISO"] = inv.AkadDate.Format("2006-01-02T15:04:05")
	}
	if inv.AkadTime != nil {
		data["AkadTime"] = *inv.AkadTime
	}
	if inv.AkadLocation != nil {
		data["AkadLocation"] = *inv.AkadLocation
	}
	if inv.AkadAddress != nil {
		data["AkadAddress"] = *inv.AkadAddress
	}
	if inv.AkadMapsURL != nil {
		data["AkadMapsURL"] = *inv.AkadMapsURL
	}

	if inv.ReceptionDate != nil {
		data["ReceptionDate"] = inv.ReceptionDate
		data["ReceptionDateFormatted"] = inv.ReceptionDate.Format("Monday, 02 January 2006")
	}
	if inv.ReceptionTime != nil {
		data["ReceptionTime"] = *inv.ReceptionTime
	}
	if inv.ReceptionLocation != nil {
		data["ReceptionLocation"] = *inv.ReceptionLocation
	}
	if inv.ReceptionAddress != nil {
		data["ReceptionAddress"] = *inv.ReceptionAddress
	}
	if inv.ReceptionMapsURL != nil {
		data["ReceptionMapsURL"] = *inv.ReceptionMapsURL
	}

	return data
}
