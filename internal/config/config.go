package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                   string
	Environment            string
	DatabaseURL            string
	SupabaseURL            string
	SupabaseAnonKey        string
	SupabaseServiceKey     string
	SupabaseJWTSecret      string
	AssemblyAIAPIKey       string
	StorageBucket          string
	StripeSecretKey        string
	StripeWebhookSecret    string
	MercadoPagoAccessToken string
	WebhookURL             string
	FrontendURL            string
}

var AppConfig *Config

func Load() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	AppConfig = &Config{
		Port:                   getEnv("PORT", "8080"),
		Environment:            getEnv("ENVIRONMENT", "development"),
		DatabaseURL:            getEnv("DATABASE_URL", ""),
		SupabaseURL:            getEnv("SUPABASE_URL", ""),
		SupabaseAnonKey:        getEnv("SUPABASE_ANON_KEY", ""),
		SupabaseServiceKey:     getEnv("SUPABASE_SERVICE_KEY", ""),
		SupabaseJWTSecret:      getEnv("SUPABASE_JWT_SECRET", ""),
		AssemblyAIAPIKey:       getEnv("ASSEMBLYAI_API_KEY", ""),
		StorageBucket:          getEnv("STORAGE_BUCKET", "litwick-uploads"),
		StripeSecretKey:        getEnv("STRIPE_SECRET_KEY", ""),
		StripeWebhookSecret:    getEnv("STRIPE_WEBHOOK_SECRET", ""),
		MercadoPagoAccessToken: getEnv("MERCADOPAGO_ACCESS_TOKEN", ""),
		WebhookURL:             getEnv("WEBHOOK_URL", "http://localhost:8080"),
		FrontendURL:            getEnv("FRONTEND_URL", "http://localhost:5173"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
