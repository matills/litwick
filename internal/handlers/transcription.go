package handlers

import (
	"context"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/matills/litwick/internal/database"
	"github.com/matills/litwick/internal/middleware"
	"github.com/matills/litwick/internal/models"
	"github.com/matills/litwick/internal/services"
)

// ProcessTranscription starts the transcription process with AssemblyAI
func ProcessTranscription(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get transcription ID from URL
	transcriptionID := c.Params("id")
	tid, err := uuid.Parse(transcriptionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid transcription ID",
		})
	}

	// Get transcription from database
	var transcription models.Transcription
	if err := database.DB.Where("id = ? AND user_id = ?", tid, user.ID).First(&transcription).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "transcription not found",
		})
	}

	// Check if already processing or completed
	if transcription.Status != models.StatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("transcription already %s", transcription.Status),
		})
	}

	// Update status to processing
	transcription.Status = models.StatusProcessing
	database.DB.Save(&transcription)

	// Start transcription in background
	go processTranscriptionAsync(transcription.ID, user.ID)

	return c.JSON(fiber.Map{
		"message":       "transcription started",
		"transcription": transcription,
	})
}

// processTranscriptionAsync handles the async transcription process
func processTranscriptionAsync(transcriptionID, userID uuid.UUID) {
	ctx := context.Background()

	var transcription models.Transcription
	if err := database.DB.First(&transcription, transcriptionID).Error; err != nil {
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return
	}

	// Initialize AssemblyAI service
	aaiService := services.NewAssemblyAIService()

	// Create transcription
	result, err := aaiService.CreateTranscription(ctx, transcription.FileURL, transcription.Language)
	if err != nil {
		transcription.Status = models.StatusFailed
		transcription.ErrorMessage = err.Error()
		database.DB.Save(&transcription)
		return
	}

	transcription.AssemblyAIID = result.ID
	database.DB.Save(&transcription)

	// Wait for completion (max 30 minutes)
	result, err = aaiService.WaitForCompletion(ctx, result.ID, 30*time.Minute)
	if err != nil {
		transcription.Status = models.StatusFailed
		transcription.ErrorMessage = err.Error()
		database.DB.Save(&transcription)
		return
	}

	// Calculate duration in minutes and check credits
	durationMinutes := (result.Duration / 1000) / 60
	if durationMinutes == 0 {
		durationMinutes = 1 // Minimum 1 minute
	}

	if !user.HasCredits(durationMinutes) {
		transcription.Status = models.StatusFailed
		transcription.ErrorMessage = "insufficient credits"
		database.DB.Save(&transcription)
		return
	}

	// Get SRT content
	srtContent, err := aaiService.GetSRT(ctx, result.ID)
	if err != nil {
		// Continue even if SRT fails
		srtContent = ""
	}

	// Update transcription
	now := time.Now()
	transcription.Status = models.StatusCompleted
	transcription.TranscriptText = &result.Text
	transcription.SRTContent = &srtContent
	transcription.Duration = result.Duration / 1000 // Convert to seconds
	transcription.CreditsUsed = durationMinutes
	transcription.CompletedAt = &now
	database.DB.Save(&transcription)

	// Deduct credits from user
	user.DeductCredits(durationMinutes)
	database.DB.Save(&user)

	// Create credit transaction record
	transaction := models.CreditTransaction{
		UserID:          user.ID,
		TranscriptionID: &transcription.ID,
		Type:            models.TransactionDebit,
		Amount:          durationMinutes,
		BalanceBefore:   user.CreditsRemaining + durationMinutes,
		BalanceAfter:    user.CreditsRemaining,
		Description:     fmt.Sprintf("Transcription: %s", transcription.FileName),
	}
	database.DB.Create(&transaction)
}

// GetTranscription returns a single transcription
func GetTranscription(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	transcriptionID := c.Params("id")
	tid, err := uuid.Parse(transcriptionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid transcription ID",
		})
	}

	var transcription models.Transcription
	if err := database.DB.Where("id = ? AND user_id = ?", tid, user.ID).First(&transcription).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "transcription not found",
		})
	}

	return c.JSON(transcription)
}

// DownloadTranscription returns transcription in requested format
func DownloadTranscription(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	transcriptionID := c.Params("id")
	tid, err := uuid.Parse(transcriptionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid transcription ID",
		})
	}

	format := c.Query("format", "txt") // txt, srt

	var transcription models.Transcription
	if err := database.DB.Where("id = ? AND user_id = ?", tid, user.ID).First(&transcription).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "transcription not found",
		})
	}

	if transcription.Status != models.StatusCompleted {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "transcription not completed",
		})
	}

	var content string
	var contentType string
	var filename string

	switch format {
	case "srt":
		if transcription.SRTContent != nil {
			content = *transcription.SRTContent
		}
		contentType = "application/x-subrip"
		filename = fmt.Sprintf("%s.srt", transcription.FileName)
	case "txt":
		fallthrough
	default:
		if transcription.TranscriptText != nil {
			content = *transcription.TranscriptText
		}
		contentType = "text/plain"
		filename = fmt.Sprintf("%s.txt", transcription.FileName)
	}

	c.Set("Content-Type", contentType)
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return c.SendString(content)
}

// UpdateTranscription allows editing the transcript text
func UpdateTranscription(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	transcriptionID := c.Params("id")
	tid, err := uuid.Parse(transcriptionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid transcription ID",
		})
	}

	var transcription models.Transcription
	if err := database.DB.Where("id = ? AND user_id = ?", tid, user.ID).First(&transcription).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "transcription not found",
		})
	}

	// Parse request body
	type UpdateRequest struct {
		TranscriptText string `json:"transcript_text"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	// Update transcript
	transcription.TranscriptText = &req.TranscriptText
	if err := database.DB.Save(&transcription).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update transcription",
		})
	}

	return c.JSON(transcription)
}

// DeleteTranscription deletes a transcription
func DeleteTranscription(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	transcriptionID := c.Params("id")
	tid, err := uuid.Parse(transcriptionID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid transcription ID",
		})
	}

	var transcription models.Transcription
	if err := database.DB.Where("id = ? AND user_id = ?", tid, user.ID).First(&transcription).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "transcription not found",
		})
	}

	// Delete from database
	if err := database.DB.Delete(&transcription).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete transcription",
		})
	}

	// TODO: Delete file from S3 (optional, can be done async)

	return c.JSON(fiber.Map{
		"message": "transcription deleted successfully",
	})
}
