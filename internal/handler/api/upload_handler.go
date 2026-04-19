package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

// UploadHandler handles file upload requests
type UploadHandler struct {
	uploadDir string
	baseURL   string
}

// NewUploadHandler creates a new upload handler
func NewUploadHandler(uploadDir, baseURL string) *UploadHandler {
	// Create upload directories if they don't exist
	dirs := []string{
		filepath.Join(uploadDir, "photos"),
		filepath.Join(uploadDir, "gallery"),
		filepath.Join(uploadDir, "music"),
		filepath.Join(uploadDir, "qr"),
	}
	for _, dir := range dirs {
		os.MkdirAll(dir, 0755)
	}

	// Trim trailing slash from baseURL
	baseURL = strings.TrimRight(baseURL, "/")

	return &UploadHandler{
		uploadDir: uploadDir,
		baseURL:   baseURL,
	}
}

// AllowedImageTypes contains allowed image MIME types
var AllowedImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
}

// AllowedAudioTypes contains allowed audio MIME types
var AllowedAudioTypes = map[string]bool{
	"audio/mpeg": true,
	"audio/mp3":  true,
	"audio/wav":  true,
	"audio/ogg":  true,
}

// UploadImage handles image upload
func (h *UploadHandler) UploadImage(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "No file uploaded",
		})
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if !AllowedImageTypes[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid file type. Allowed: JPEG, PNG, GIF, WebP",
		})
	}

	// Validate file size (max 5MB)
	if file.Size > 5*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "File too large. Maximum size is 5MB",
		})
	}

	// Get upload type from query
	uploadType := c.Query("type", "photos")
	if uploadType != "photos" && uploadType != "gallery" && uploadType != "qr" {
		uploadType = "photos"
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	savePath := filepath.Join(h.uploadDir, uploadType, filename)

	// Save file
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to save file",
		})
	}

	// Generate URL
	fileURL := fmt.Sprintf("%s/uploads/%s/%s", h.baseURL, uploadType, filename)

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "File uploaded successfully",
		"data": fiber.Map{
			"url":      fileURL,
			"filename": filename,
			"size":     file.Size,
			"type":     contentType,
		},
	})
}

// UploadMusic handles music file upload
func (h *UploadHandler) UploadMusic(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "No file uploaded",
		})
	}

	// Validate file type
	contentType := file.Header.Get("Content-Type")
	if !AllowedAudioTypes[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid file type. Allowed: MP3, WAV, OGG",
		})
	}

	// Validate file size (max 10MB)
	if file.Size > 10*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "File too large. Maximum size is 10MB",
		})
	}

	// Generate unique filename
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
	savePath := filepath.Join(h.uploadDir, "music", filename)

	// Save file
	if err := c.SaveFile(file, savePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to save file",
		})
	}

	// Generate URL
	fileURL := fmt.Sprintf("%s/uploads/music/%s", h.baseURL, filename)

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "Music uploaded successfully",
		"data": fiber.Map{
			"url":      fileURL,
			"filename": filename,
			"size":     file.Size,
			"type":     contentType,
		},
	})
}

// UploadMultipleImages handles multiple image upload
func (h *UploadHandler) UploadMultipleImages(c *fiber.Ctx) error {
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "No files uploaded",
		})
	}

	files := form.File["files"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "No files uploaded",
		})
	}

	if len(files) > 20 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Maximum 20 files allowed per upload",
		})
	}

	uploadType := c.Query("type", "gallery")
	var uploadedFiles []fiber.Map
	var errors []string

	for _, file := range files {
		// Validate file type
		contentType := file.Header.Get("Content-Type")
		if !AllowedImageTypes[contentType] {
			errors = append(errors, fmt.Sprintf("%s: Invalid file type", file.Filename))
			continue
		}

		// Validate file size (max 5MB)
		if file.Size > 5*1024*1024 {
			errors = append(errors, fmt.Sprintf("%s: File too large", file.Filename))
			continue
		}

		// Generate unique filename
		ext := filepath.Ext(file.Filename)
		filename := fmt.Sprintf("%s_%d%s", uuid.New().String(), time.Now().Unix(), ext)
		savePath := filepath.Join(h.uploadDir, uploadType, filename)

		// Save file
		if err := c.SaveFile(file, savePath); err != nil {
			errors = append(errors, fmt.Sprintf("%s: Failed to save", file.Filename))
			continue
		}

		// Generate URL
		fileURL := fmt.Sprintf("%s/uploads/%s/%s", h.baseURL, uploadType, filename)

		uploadedFiles = append(uploadedFiles, fiber.Map{
			"url":           fileURL,
			"filename":      filename,
			"original_name": file.Filename,
			"size":          file.Size,
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": fmt.Sprintf("Uploaded %d of %d files", len(uploadedFiles), len(files)),
		"data":    uploadedFiles,
		"errors":  errors,
	})
}

// DeleteFile deletes an uploaded file
func (h *UploadHandler) DeleteFile(c *fiber.Ctx) error {
	fileURL := c.Query("url")
	if fileURL == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "File URL is required",
		})
	}

	// Extract path from URL
	// Expected format: /uploads/{type}/{filename}
	parts := strings.Split(fileURL, "/uploads/")
	if len(parts) != 2 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid file URL",
		})
	}

	filePath := filepath.Join(h.uploadDir, parts[1])

	// Security check - make sure path is within upload directory
	absUploadDir, _ := filepath.Abs(h.uploadDir)
	absFilePath, _ := filepath.Abs(filePath)
	if !strings.HasPrefix(absFilePath, absUploadDir) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   true,
			"message": "Invalid file path",
		})
	}

	// Delete file
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error":   true,
				"message": "File not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   true,
			"message": "Failed to delete file",
		})
	}

	return c.JSON(fiber.Map{
		"error":   false,
		"message": "File deleted successfully",
	})
}
