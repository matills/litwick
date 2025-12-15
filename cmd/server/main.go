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
	config.Load()
	log.Println("Configuration loaded")

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected")

	if err := database.Migrate(); err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
	log.Println("Database migrations completed")

	app := fiber.New(fiber.Config{
		BodyLimit: 500 * 1024 * 1024,
	})

	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins:     config.AppConfig.FrontendURL,
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, PUT, DELETE, OPTIONS",
		AllowCredentials: true,
	}))

	app.Get("/health", handlers.HealthCheck)

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Use(middleware.AuthMiddleware())
	auth.Get("/me", handlers.GetMe)
	auth.Put("/settings", handlers.UpdateSettings)

	dashboard := api.Group("/dashboard")
	dashboard.Use(middleware.AuthMiddleware())
	dashboard.Get("/", handlers.GetDashboard)

	upload := api.Group("/upload")
	upload.Use(middleware.AuthMiddleware())
	upload.Post("/", handlers.UploadFile)

	transcriptions := api.Group("/transcriptions")
	transcriptions.Use(middleware.AuthMiddleware())
	transcriptions.Get("/", handlers.GetTranscriptions)
	transcriptions.Post("/:id/process", handlers.ProcessTranscription)
	transcriptions.Get("/:id", handlers.GetTranscription)
	transcriptions.Put("/:id", handlers.UpdateTranscription)
	transcriptions.Delete("/:id", handlers.DeleteTranscription)
	transcriptions.Get("/:id/download", handlers.DownloadTranscription)

	payments := api.Group("/payments")
	payments.Get("/packages", handlers.GetCreditPackages)
	payments.Post("/webhook", handlers.WebhookMercadoPago)
	payments.Use(middleware.AuthMiddleware())
	payments.Post("/create", handlers.CreatePayment)
	payments.Get("/history", handlers.GetPaymentHistory)
	payments.Get("/success", handlers.ProcessPaymentSuccess)

	distPath := "./frontend/dist"

	if _, err := os.Stat(distPath); os.IsNotExist(err) {
		log.Println("Warning: frontend/dist directory not found. Please build the frontend with 'cd frontend && npm run build'")
	} else {
		app.Static("/", distPath)

		app.Use(func(c *fiber.Ctx) error {
			return c.SendFile(distPath + "/index.html")
		})
	}

	port := config.AppConfig.Port
	log.Printf("Server starting on port %s", port)
	if err := app.Listen(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
