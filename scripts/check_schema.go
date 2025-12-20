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

	type ColumnInfo struct {
		ColumnName string
		DataType   string
		IsNullable string
	}

	var columns []ColumnInfo
	query := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_name = 'payments'
		ORDER BY ordinal_position
	`

	if err := database.DB.Raw(query).Scan(&columns).Error; err != nil {
		log.Fatal("Failed to query schema:", err)
	}

	log.Println("\nPayments table schema:")
	for _, col := range columns {
		log.Printf("  %s (%s) - Nullable: %s", col.ColumnName, col.DataType, col.IsNullable)
	}
}
