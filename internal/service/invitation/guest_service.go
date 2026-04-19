package invitation

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"
	"time"

	invitationDomain "invitation-api/internal/domain/invitation"

	"github.com/google/uuid"
	"github.com/xuri/excelize/v2"
)

// AddGuest adds a new guest
func (s *service) AddGuest(invitationID uuid.UUID, name, guestType string) (*invitationDomain.Guest, error) {
	guest := &invitationDomain.Guest{
		ID:           uuid.New(),
		InvitationID: invitationID,
		Name:         name,
		Type:         guestType,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.invitationRepo.CreateGuest(guest); err != nil {
		return nil, err
	}
	return guest, nil
}

// GetGuests retrieves guests
func (s *service) GetGuests(invitationID uuid.UUID) ([]*invitationDomain.Guest, error) {
	return s.invitationRepo.GetGuestsByInvitationID(invitationID)
}

// DeleteGuest deletes a guest
func (s *service) DeleteGuest(id uuid.UUID) error {
	return s.invitationRepo.DeleteGuest(id)
}

// ImportGuests imports guests from Excel
func (s *service) ImportGuests(invitationID uuid.UUID, file io.Reader) (int, error) {
	f, err := excelize.OpenReader(file)
	if err != nil {
		return 0, fmt.Errorf("failed to open excel file: %w", err)
	}
	defer f.Close()

	rows, err := f.GetRows("Sheet1")
	if err != nil {
		// Try getting the first sheet name if "Sheet1" fails
		sheets := f.GetSheetList()
		if len(sheets) == 0 {
			return 0, errors.New("no sheets found")
		}
		rows, err = f.GetRows(sheets[0])
		if err != nil {
			return 0, err
		}
	}

	count := 0
	for i, row := range rows {
		if i == 0 {
			continue // Skip header
		}
		if len(row) < 1 || row[0] == "" {
			continue
		}

		name := row[0]
		guestType := "Regular"
		if len(row) > 1 && row[1] != "" {
			guestType = row[1]
		}

		// Optional greeting
		var greeting *string
		if len(row) > 2 && row[2] != "" {
			val := row[2]
			greeting = &val
		}

		guest := &invitationDomain.Guest{
			ID:           uuid.New(),
			InvitationID: invitationID,
			Name:         name,
			Type:         guestType,
			Greeting:     greeting,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		if err := s.invitationRepo.CreateGuest(guest); err == nil {
			count++
		}
	}

	return count, nil
}

// GetGuestTemplate returns Excel template
func (s *service) GetGuestTemplate() (*bytes.Buffer, string) {
	f := excelize.NewFile()
	defer f.Close()

	// Headers
	f.SetCellValue("Sheet1", "A1", "Nama Tamu")
	f.SetCellValue("Sheet1", "B1", "Kategori (Regular/VIP)")
	f.SetCellValue("Sheet1", "C1", "Sapaan (Opsional)")

	// Example
	f.SetCellValue("Sheet1", "A2", "Budi Santoso")
	f.SetCellValue("Sheet1", "B2", "Regular")
	f.SetCellValue("Sheet1", "C2", "Di Tempat")

	// Adjust column width
	f.SetColWidth("Sheet1", "A", "A", 30)
	f.SetColWidth("Sheet1", "B", "B", 20)
	f.SetColWidth("Sheet1", "C", "C", 20)

	buf, _ := f.WriteToBuffer()
	return buf, "template_tamu.xlsx"
}

// ExportGuests exports guest list with URLs
func (s *service) ExportGuests(invitationID uuid.UUID) (*bytes.Buffer, string, error) {
	guests, err := s.invitationRepo.GetGuestsByInvitationID(invitationID)
	if err != nil {
		return nil, "", err
	}

	invitation, err := s.invitationRepo.GetByID(invitationID)
	if err != nil {
		return nil, "", err
	}

	// Get Subdomain from Tenant
	tenant, err := s.tenantRepo.GetByID(invitation.TenantID)
	if err != nil {
		return nil, "", fmt.Errorf("failed to get tenant: %w", err)
	}

	f := excelize.NewFile()
	defer f.Close()

	f.SetCellValue("Sheet1", "A1", "Nama Tamu")
	f.SetCellValue("Sheet1", "B1", "Kategori")
	f.SetCellValue("Sheet1", "C1", "Link Undangan")

	// Adjust column width
	f.SetColWidth("Sheet1", "A", "A", 30)
	f.SetColWidth("Sheet1", "B", "B", 20)
	f.SetColWidth("Sheet1", "C", "C", 60)

	// Base URL construction
	// Access localhost or production domain logic could be improved
	// Assuming default app structure here
	baseURL := "http://localhost:3000/" + tenant.Subdomain

	for i, guest := range guests {
		row := i + 2
		f.SetCellValue("Sheet1", fmt.Sprintf("A%d", row), guest.Name)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%d", row), guest.Type)

		encodedName := url.QueryEscape(guest.Name)
		link := fmt.Sprintf("%s?to=%s", baseURL, encodedName)

		f.SetCellValue("Sheet1", fmt.Sprintf("C%d", row), link)
	}

	buf, _ := f.WriteToBuffer()
	return buf, "daftar_tamu_" + tenant.Subdomain + ".xlsx", nil
}
