package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matills/litwick/internal/database"
	"github.com/matills/litwick/internal/middleware"
)

func GetMe(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

func UpdateSettings(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	type SettingsRequest struct {
		DefaultLanguage     *string `json:"default_language"`
		DefaultExportFormat *string `json:"default_export_format"`
		IncludeTimestamps   *bool   `json:"include_timestamps"`
		DetectSpeakers      *bool   `json:"detect_speakers"`
		EmailNotifications  *bool   `json:"email_notifications"`
		PromotionalEmails   *bool   `json:"promotional_emails"`
	}

	var req SettingsRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid request body",
		})
	}

	if req.DefaultLanguage != nil {
		user.DefaultLanguage = *req.DefaultLanguage
	}
	if req.DefaultExportFormat != nil {
		user.DefaultExportFormat = *req.DefaultExportFormat
	}
	if req.IncludeTimestamps != nil {
		user.IncludeTimestamps = *req.IncludeTimestamps
	}
	if req.DetectSpeakers != nil {
		user.DetectSpeakers = *req.DetectSpeakers
	}
	if req.EmailNotifications != nil {
		user.EmailNotifications = *req.EmailNotifications
	}
	if req.PromotionalEmails != nil {
		user.PromotionalEmails = *req.PromotionalEmails
	}

	if err := database.DB.Save(user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to update settings",
		})
	}

	return c.JSON(fiber.Map{
		"user": user,
	})
}

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"message": "Litwick API is running",
	})
}
