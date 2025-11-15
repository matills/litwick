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

	// Get plan limits for the user
	limits := user.GetPlanLimits()

	// Validate file size based on plan
	// Rough estimate: 1 minute of audio ≈ 1MB for MP3, so we use max minutes * 2MB as safety
	maxSize := int64(limits.MaxMinutesPerFile * 2 * 1024 * 1024)
	if maxSize > 500*1024*1024 {
		maxSize = 500 * 1024 * 1024 // Hard limit at 500MB
	}

	if file.Size > maxSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":       fmt.Sprintf("file too large for your plan (max %d minutes per file)", limits.MaxMinutesPerFile),
			"max_size":    maxSize,
			"plan":        user.Plan,
			"max_minutes": limits.MaxMinutesPerFile,
		})
	}

	// Check if user has any credits remaining
	if user.CreditsRemaining <= 0 {
		return c.Status(fiber.StatusPaymentRequired).JSON(fiber.Map{
			"error":             "insufficient credits",
			"credits_remaining": user.CreditsRemaining,
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
		Status:   models.StatusProcessing,
		Language: language,
	}

	if err := database.DB.Create(&transcription).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create transcription record",
		})
	}

	// Start transcription process in background
	go processTranscriptionAsync(transcription.ID, user.ID)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":       "file uploaded successfully, transcription started",
		"transcription": transcription,
	})
}
