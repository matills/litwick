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

func ProcessTranscription(c *fiber.Ctx) error {
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

	if transcription.Status != models.StatusPending {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": fmt.Sprintf("transcription already %s", transcription.Status),
		})
	}

	transcription.Status = models.StatusProcessing
	database.DB.Save(&transcription)

	go processTranscriptionAsync(transcription.ID, user.ID)

	return c.JSON(fiber.Map{
		"message":       "transcription started",
		"transcription": transcription,
	})
}

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

	aaiService := services.NewAssemblyAIService()

	result, err := aaiService.CreateTranscription(ctx, transcription.FileURL, transcription.Language)
	if err != nil {
		transcription.Status = models.StatusFailed
		transcription.ErrorMessage = err.Error()
		database.DB.Save(&transcription)
		return
	}

	transcription.AssemblyAIID = result.ID
	database.DB.Save(&transcription)

	result, err = aaiService.WaitForCompletion(ctx, result.ID, 30*time.Minute)
	if err != nil {
		transcription.Status = models.StatusFailed
		transcription.ErrorMessage = err.Error()
		database.DB.Save(&transcription)
		return
	}

	durationMinutes := (result.Duration / 1000) / 60
	if durationMinutes == 0 {
		durationMinutes = 1
	}

	if !user.HasCredits(durationMinutes) {
		transcription.Status = models.StatusFailed
		transcription.ErrorMessage = "insufficient credits"
		database.DB.Save(&transcription)
		return
	}

	srtContent, err := aaiService.GetSRT(ctx, result.ID)
	if err != nil {
		srtContent = ""
	}

	vttContent, err := aaiService.GetVTT(ctx, result.ID)
	if err != nil {
		vttContent = ""
	}

	now := time.Now()
	transcription.Status = models.StatusCompleted
	transcription.TranscriptText = &result.Text
	transcription.SRTContent = &srtContent
	transcription.VTTContent = &vttContent
	transcription.Duration = result.Duration / 1000 // Convert to seconds
	transcription.CreditsUsed = durationMinutes
	transcription.CompletedAt = &now
	database.DB.Save(&transcription)

	user.DeductCredits(durationMinutes)
	database.DB.Save(&user)

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

	format := c.Query("format", "txt")

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
	case "vtt":
		if transcription.VTTContent != nil {
			content = *transcription.VTTContent
		}
		contentType = "text/vtt"
		filename = fmt.Sprintf("%s.vtt", transcription.FileName)
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

	type UpdateRequest struct {
		TranscriptText string `json:"transcript_text"`
	}

	var req UpdateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	transcription.TranscriptText = &req.TranscriptText
	if err := database.DB.Save(&transcription).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update transcription",
		})
	}

	return c.JSON(transcription)
}

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

	if err := database.DB.Delete(&transcription).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to delete transcription",
		})
	}

	go func() {
		ctx := context.Background()
		storageService, err := services.NewStorageService(ctx)
		if err != nil {
			fmt.Printf("Failed to initialize storage service: %v\n", err)
			return
		}

		filePath := storageService.ExtractFilePathFromURL(transcription.FileURL)

		if err := storageService.DeleteFile(ctx, filePath); err != nil {
			fmt.Printf("Failed to delete file from storage: %v\n", err)
		}
	}()

	return c.JSON(fiber.Map{
		"message": "transcription deleted successfully",
	})
}
