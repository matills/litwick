package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/matills/litwick/internal/database"
	"github.com/matills/litwick/internal/models"
	"github.com/matills/litwick/internal/services"
)

// AuthMiddleware verifies Supabase JWT token and loads user
func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get Authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "missing authorization header",
			})
		}

		// Extract token
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid authorization format",
			})
		}

		// Verify token with Supabase
		supabaseUserID, email, err := services.VerifySupabaseToken(token)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		// Find or create user in our database
		var user models.User
		result := database.DB.Where("supabase_user_id = ?", supabaseUserID).First(&user)

		if result.Error != nil {
			// User doesn't exist, create new user
			user = models.User{
				SupabaseUserID:  supabaseUserID,
				Email:           email,
				CreditsRemaining: 300, // 5 hours free
				Plan:            "free",
			}
			if err := database.DB.Create(&user).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "failed to create user",
				})
			}
		}

		// Store user in context
		c.Locals("user", user)
		c.Locals("userID", user.ID.String())

		return c.Next()
	}
}

// GetUser retrieves the authenticated user from context
func GetUser(c *fiber.Ctx) *models.User {
	user, ok := c.Locals("user").(models.User)
	if !ok {
		return nil
	}
	return &user
}
