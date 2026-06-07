package config

import (
	"log"
	"os"
)

type Config struct {
	DBUser            string
	DBPass            string
	DBHost            string
	DBPort            string
	DBName            string
	Port              string
	AppEnv            string
	JWTSecret         string
	SMTPHost          string
	SMTPPort          string
	SMTPEmail         string
	SMTPPass          string
	WAHAURL           string
	WAHAApiKey        string
	WAHAWebhookSecret string
	FrontendURL       string
	MidtransServerKey string
	MidtransClientKey string
	MidtransIsSandbox bool
}

func LoadConfig() *Config {


	cfg := &Config{
		DBUser:            os.Getenv("DB_USER"),
		DBPass:            os.Getenv("DB_PASS"),
		DBHost:            os.Getenv("DB_HOST"),
		DBPort:            os.Getenv("DB_PORT"),
		DBName:            os.Getenv("DB_NAME"),
		Port:              os.Getenv("PORT"),
		AppEnv:            os.Getenv("APP_ENV"),
		JWTSecret:         os.Getenv("JWT_SECRET"),
		SMTPHost:          os.Getenv("SMTP_HOST"),
		SMTPPort:          os.Getenv("SMTP_PORT"),
		SMTPEmail:         os.Getenv("SMTP_EMAIL"),
		SMTPPass:          os.Getenv("SMTP_PASS"),
		WAHAURL:           os.Getenv("WAHA_URL"),
		WAHAApiKey:        os.Getenv("WAHA_API_KEY"),
		WAHAWebhookSecret: os.Getenv("WAHA_WEBHOOK_SECRET"),
		FrontendURL:       os.Getenv("FRONTEND_URL"),
		MidtransServerKey: os.Getenv("MIDTRANS_SERVER_KEY"),
		MidtransClientKey: os.Getenv("MIDTRANS_CLIENT_KEY"),
		MidtransIsSandbox: os.Getenv("MIDTRANS_IS_SANDBOX") == "true",
	}

	// Validasi khusus untuk lingkungan production demi menjaga keamanan webhook
	if cfg.AppEnv == "production" && cfg.WAHAWebhookSecret == "" {
		log.Fatal("CRITICAL: WAHA_WEBHOOK_SECRET must be set in production")
	}

	return cfg
}
