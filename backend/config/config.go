package config

import (
	"log"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	DBUser                string
	DBPass                string
	DBHost                string
	DBPort                string
	DBName                string
	Port                  string
	AppEnv                string
	JWTSecret             string
	SMTPHost              string
	SMTPPort              string
	SMTPEmail             string
	SMTPPass              string
	WAHAURL               string
	WAHAApiKey            string
	WAHAWebhookSecret     string
	FrontendURL           string
	MidtransServerKey     string
	MidtransClientKey     string
	MidtransIsSandbox     bool
	RedisHost             string
	RedisPort             string
	RedisPass             string
	CaptchaEnabled        bool
	TurnstileSecretKey    string
	RateLimitBlockDelayMs int
}

func LoadConfig() *Config {

	cfg := &Config{
		DBUser:                os.Getenv("DB_USER"),
		DBPass:                os.Getenv("DB_PASS"),
		DBHost:                os.Getenv("DB_HOST"),
		DBPort:                os.Getenv("DB_PORT"),
		DBName:                os.Getenv("DB_NAME"),
		Port:                  os.Getenv("PORT"),
		AppEnv:                os.Getenv("APP_ENV"),
		JWTSecret:             os.Getenv("JWT_SECRET"),
		SMTPHost:              os.Getenv("SMTP_HOST"),
		SMTPPort:              os.Getenv("SMTP_PORT"),
		SMTPEmail:             os.Getenv("SMTP_EMAIL"),
		SMTPPass:              os.Getenv("SMTP_PASS"),
		WAHAURL:               os.Getenv("WAHA_URL"),
		WAHAApiKey:            os.Getenv("WAHA_API_KEY"),
		WAHAWebhookSecret:     os.Getenv("WAHA_WEBHOOK_SECRET"),
		FrontendURL:           os.Getenv("FRONTEND_URL"),
		MidtransServerKey:     os.Getenv("MIDTRANS_SERVER_KEY"),
		MidtransClientKey:     os.Getenv("MIDTRANS_CLIENT_KEY"),
		MidtransIsSandbox:     os.Getenv("MIDTRANS_IS_SANDBOX") == "true",
		RedisHost:             os.Getenv("REDIS_HOST"),
		RedisPort:             os.Getenv("REDIS_PORT"),
		RedisPass:             os.Getenv("REDIS_PASS"),
		CaptchaEnabled:        envBool("CAPTCHA_ENABLED") || os.Getenv("TURNSTILE_SECRET_KEY") != "",
		TurnstileSecretKey:    os.Getenv("TURNSTILE_SECRET_KEY"),
		RateLimitBlockDelayMs: envInt("RATE_LIMIT_BLOCK_DELAY_MS", 0),
	}

	// Validasi khusus untuk lingkungan production demi menjaga keamanan webhook
	if cfg.AppEnv == "production" && cfg.WAHAWebhookSecret == "" {
		log.Fatal("CRITICAL: WAHA_WEBHOOK_SECRET must be set in production")
	}

	if cfg.AppEnv == "production" && cfg.CaptchaEnabled && cfg.TurnstileSecretKey == "" {
		log.Fatal("CRITICAL: TURNSTILE_SECRET_KEY must be set when CAPTCHA is enabled in production")
	}

	return cfg
}


func envBool(key string) bool {
	value := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
	return value == "true" || value == "1" || value == "yes" || value == "on"
}

func envInt(key string, fallback int) int {
	if parsed, err := strconv.Atoi(strings.TrimSpace(os.Getenv(key))); err == nil {
		return parsed
	}
	return fallback
}
