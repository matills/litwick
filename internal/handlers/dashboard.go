package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matills/litwick/internal/database"
	"github.com/matills/litwick/internal/middleware"
	"github.com/matills/litwick/internal/models"
)

// GetDashboard returns user's transcriptions and stats
func GetDashboard(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get all transcriptions for user
	var transcriptions []models.Transcription
	if err := database.DB.Where("user_id = ?", user.ID).
		Order("created_at DESC").
		Find(&transcriptions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch transcriptions",
		})
	}

	// Calculate stats
	stats := calculateStats(transcriptions, *user)

	return c.JSON(fiber.Map{
		"user":           user,
		"transcriptions": transcriptions,
		"stats":          stats,
	})
}

// GetTranscriptions returns paginated list of transcriptions
func GetTranscriptions(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "unauthorized",
		})
	}

	// Get pagination params
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 10)
	offset := (page - 1) * limit

	// Get transcriptions with pagination
	var transcriptions []models.Transcription
	var total int64

	database.DB.Model(&models.Transcription{}).Where("user_id = ?", user.ID).Count(&total)

	if err := database.DB.Where("user_id = ?", user.ID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&transcriptions).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to fetch transcriptions",
		})
	}

	return c.JSON(fiber.Map{
		"transcriptions": transcriptions,
		"pagination": fiber.Map{
			"page":       page,
			"limit":      limit,
			"total":      total,
			"totalPages": (total + int64(limit) - 1) / int64(limit),
		},
	})
}

type DashboardStats struct {
	TotalTranscriptions int `json:"total_transcriptions"`
	CompletedCount      int `json:"completed_count"`
	ProcessingCount     int `json:"processing_count"`
	FailedCount         int `json:"failed_count"`
	TotalMinutesUsed    int `json:"total_minutes_used"`
	CreditsRemaining    int `json:"credits_remaining"`
}

func calculateStats(transcriptions []models.Transcription, user models.User) DashboardStats {
	stats := DashboardStats{
		TotalTranscriptions: len(transcriptions),
		CreditsRemaining:    user.CreditsRemaining,
	}

	for _, t := range transcriptions {
		switch t.Status {
		case models.StatusCompleted:
			stats.CompletedCount++
			stats.TotalMinutesUsed += t.CreditsUsed
		case models.StatusProcessing, models.StatusPending:
			stats.ProcessingCount++
		case models.StatusFailed:
			stats.FailedCount++
		}
	}

	return stats
}
