package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
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
	RedisAddr         string
}

func LoadConfig() *Config {
	// Hanya muat .env jika belum ada environment variable (berarti bukan di Docker)
	if os.Getenv("DB_HOST") == "" {
		err := godotenv.Load()
		if err != nil {
			log.Println("warning : file .env tidak ditemukan")
		}
	}

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
		RedisAddr:         os.Getenv("REDIS_ADDR"),
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" || cfg.JWTSecret == "" {
		log.Fatal("CRITICAL: DB_HOST, DB_USER, DB_NAME, and JWT_SECRET must be set in environment variables")
	}
	if cfg.RedisAddr == "" {
		if cfg.AppEnv == "production" {
			log.Fatal("CRITICAL: REDIS_ADDR must be set in environment variables")
		}
		cfg.RedisAddr = "127.0.0.1:6379"
	}
	if cfg.AppEnv == "production" && cfg.WAHAWebhookSecret == "" {
		log.Fatal("CRITICAL: WAHA_WEBHOOK_SECRET must be set in production")
	}

	return cfg
}
