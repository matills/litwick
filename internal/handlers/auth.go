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

	return c.JSON(fiber.Map{
		"user": user,
	})
}

// HealthCheck returns the health status of the API
func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "ok",
		"message": "Litwick API is running",
	})
}
