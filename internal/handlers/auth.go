package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/matills/litwick/internal/middleware"
)

// GetMe returns the authenticated user's information
func GetMe(c *fiber.Ctx) error {
	user := middleware.GetUser(c)
	if user == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "user not found",
		})
	}

	// Get plan limits
	limits := user.GetPlanLimits()

	return c.JSON(fiber.Map{
		"user": user,
		"plan_limits": fiber.Map{
			"monthly_minutes":        limits.MonthlyMinutes,
			"max_minutes_per_file":   limits.MaxMinutesPerFile,
			"max_file_uploads_month": limits.MaxFileUploadsMonth,
		},
	})
}

// HealthCheck returns the health status of the API
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"message": "Litwick API is running",
	})
}
