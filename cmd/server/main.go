package main

import (
	"log"

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

	// Public routes
	app.Get("/", handlers.HealthCheck)
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

	// Transcription routes (protected)
	transcriptions := api.Group("/transcriptions")
	transcriptions.Use(middleware.AuthMiddleware())
	transcriptions.Get("/", handlers.GetTranscriptions)
	transcriptions.Post("/upload", handlers.UploadFile)
	transcriptions.Post("/:id/process", handlers.ProcessTranscription)
	transcriptions.Get("/:id", handlers.GetTranscription)
	transcriptions.Put("/:id", handlers.UpdateTranscription)
	transcriptions.Delete("/:id", handlers.DeleteTranscription)
	transcriptions.Get("/:id/download", handlers.DownloadTranscription)

	// Start server
	port := config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
