package main

import (
	"log"

	"github.com/matills/litwick/internal/config"
	"github.com/matills/litwick/internal/database"
)

func main() {
	config.Load()
	log.Println("Configuration loaded")

	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	log.Println("Database connected")

	log.Println("Fixing payments table - dropping old column mercado_pago_id...")

	query := "ALTER TABLE payments DROP COLUMN IF EXISTS mercado_pago_id"

	if err := database.DB.Exec(query).Error; err != nil {
		log.Printf("Warning: %v", err)
	} else {
		log.Printf("Success: Dropped mercado_pago_id column")
	}

	log.Println("Migration completed!")
}
