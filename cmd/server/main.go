package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/matills/litwick/internal/config"
	"github.com/matills/litwick/internal/database"
	"github.com/matills/litwick/internal/handlers"
	"github.com/matills/litwick/internal/middleware"
)

func main() {
	// Load configuration
	config.Load()
	log.Println("Configuration loaded")

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected")

	// Run migrations
	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrations completed")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024, // 500MB for file uploads
	})

	// Middleware
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.FrontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	// Health check
	app.Get("/health", handlers.HealthCheck)

	// API routes
	api := app.Group("/api")

	// Auth routes (protected)
	auth := api.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	auth.Get("/me", handlers.GetMe)

	// Dashboard routes (protected)
	dashboard := api.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware())
	dashboard.Get("/", handlers.GetDashboard)

	// Upload route (protected)
	upload := api.Group("/upload")
	upload.Use(middleware.AuthMiddleware())
	upload.Post("/", handlers.UploadFile)

	// Transcription routes (protected)
	transcriptions := api.Group("/transcriptions")
	transcriptions.Use(middleware.AuthMiddleware())
	transcriptions.Get("/", handlers.GetTranscriptions)
	transcriptions.Post("/:id/process", handlers.ProcessTranscription)
	transcriptions.Get("/:id", handlers.GetTranscription)
	transcriptions.Put("/:id", handlers.UpdateTranscription)
	transcriptions.Delete("/:id", handlers.DeleteTranscription)
	transcriptions.Get("/:id/download", handlers.DownloadTranscription)

	// Serve static files from frontend/dist
	distPath := "./frontend/dist"

	// Check if dist directory exists
	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		log.Println("Warning: frontend/dist directory not found. Please build the frontend with 'cd frontend && npm run build'")
	} else {
		// Serve static files
		app.Static("/", distPath)

		// SPA fallback - serve index.html for all non-API routes
		app.Use(func(c *fiber.Ctx) error {
			// If the request is not for an API route and not a static file
			// serve the index.html (for client-side routing)
			return c.SendFile(distPath + "/index.html")
		})
	}

	// Start server
	port := config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
