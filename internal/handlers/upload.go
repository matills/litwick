package handlers

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/matills/litwick/internal/database"
	"github.com/matills/litwick/internal/middleware"
	"github.com/matills/litwick/internal/models"
	"github.com/matills/litwick/internal/services"
)

var allowedExtensions = map[string]bool{
	".mp3":  true,
	".mp4":  true,
	".wav":  true,
	".m4a":  true,
	".flac": true,
	".ogg":  true,
	".webm": true,
	".avi":  true,
	".mov":  true,
	".mkv":  true,
}

// UploadFile handles file upload and creates a transcription job
func UploadFile(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "no file uploaded",
		})
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(file.Filename))
	if !allowedExtensions[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("unsupported file type: %s", ext),
		})
	}

	// Validate file size (max 500MB)
	maxSize := int64(500 * 1024 * 1024) // 500MB
	if file.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file too large (max 500MB)",
		})
	}

	// Open file
	fileReader, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to read file",
		})
	}
	defer fileReader.Close()

	// Upload to S3
	storageService, err := services.NewStorageService(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to initialize storage",
		})
	}

	fileURL, err := storageService.UploadFile(c.Context(), fileReader, file.Filename, file.Header.Get("Content-Type"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to upload file",
		})
	}

	// Get language from form (optional, default to Spanish)
	language := c.FormValue("language")
	if language == "" {
		language = "es"
	}

	// Create transcription record
	transcription := models.Transcription{
		ID:       uuid.New(),
		UserID:   user.ID,
		FileName: file.Filename,
		FileURL:  fileURL,
		FileSize: file.Size,
		Status:   models.StatusPending,
		Language: language,
	}

	if err := database.DB.Create(&transcription).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create transcription record",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":       "file uploaded successfully",
		"transcription": transcription,
	})
}
