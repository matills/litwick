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

	log.Println("Fixing payments table constraints...")

	queries := []string{
		"ALTER TABLE payments ALTER COLUMN mercado_pago_payment_id DROP NOT NULL",
		"ALTER TABLE payments ALTER COLUMN preference_id DROP NOT NULL",
		"ALTER TABLE payments ALTER COLUMN payment_method DROP NOT NULL",
	}

	for _, query := range queries {
		if err := database.DB.Exec(query).Error; err != nil {
			log.Printf("Warning: %v", err)
		} else {
			log.Printf("Success: %s", query)
		}
	}

	log.Println("Migration completed!")
}
