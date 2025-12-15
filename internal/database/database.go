package database

import (
	"fmt"
	"log"

	"github.com/matills/litwick/internal/config"
	"github.com/matills/litwick/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func Connect() error {
	var err error

	logLevel := logger.Info
	if config.AppConfig.Environment == "production" {
		logLevel = logger.Error
	}

	DB, err = gorm.Open(postgres.Open(config.AppConfig.DatabaseURL), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	log.Println("Database connected successfully")
	return nil
}

func Migrate() error {
	log.Println("Running database migrations...")

	err := DB.AutoMigrate(
		&models.User{},
		&models.Transcription{},
		&models.CreditTransaction{},
		&models.Payment{},
	)

	if err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}

	log.Println("Database migrations completed successfully")
	return nil
}

func Close() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
