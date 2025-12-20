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

	log.Println("Cleaning up old columns from payments table...")

	oldColumns := []string{
		"subscription_id",
		"plan_type",
		"payment_type",
		"payment_method_id",
		"external_reference",
		"mercado_pago_details",
		"processed_at",
	}

	for _, col := range oldColumns {
		query := "ALTER TABLE payments DROP COLUMN IF EXISTS " + col
		if err := database.DB.Exec(query).Error; err != nil {
			log.Printf("Warning dropping %s: %v", col, err)
		} else {
			log.Printf("Success: Dropped %s column", col)
		}
	}

	log.Println("Cleanup completed!")
}
