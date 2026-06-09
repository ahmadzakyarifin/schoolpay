package middleware

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	turnstileSiteVerifyURL = "https://challenges.cloudflare.com/turnstile/v0/siteverify"
	captchaVerifiedKey     = "captcha_verified"
)

type captchaPayload struct {
	CaptchaToken   string `json:"captcha_token"`
	TurnstileToken string `json:"turnstile_token"`
}

type turnstileResponse struct {
	Success    bool     `json:"success"`
	ErrorCodes []string `json:"error-codes"`
	Hostname   string   `json:"hostname"`
}

type turnstileVerifier struct {
	secretKey string
	client    *http.Client
}

var captchaChallengeEnabled bool

func CaptchaMiddleware(cfg *config.Config) gin.HandlerFunc {
	captchaChallengeEnabled = cfg != nil && cfg.CaptchaEnabled

	if cfg == nil || !cfg.CaptchaEnabled {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	verifier := &turnstileVerifier{
		secretKey: cfg.TurnstileSecretKey,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	return func(c *gin.Context) {
		if strings.TrimSpace(verifier.secretKey) == "" {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"status":  false,
				"message": "Konfigurasi CAPTCHA belum lengkap.",
				"data":    nil,
			})
			return
		}

		var payload captchaPayload
		_ = c.ShouldBindBodyWithJSON(&payload)
		restoreBody(c)

		token := captchaTokenFromPayload(payload)
		if token == "" {
			c.Next()
			return
		}

		if err := verifier.verify(c.Request.Context(), token, c.ClientIP()); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status":  false,
				"message": "Verifikasi CAPTCHA gagal. Silakan coba lagi.",
				"data":    nil,
			})
			return
		}

		c.Set(captchaVerifiedKey, true)
		c.Next()
	}
}

func captchaTokenFromPayload(payload captchaPayload) string {
	if token := strings.TrimSpace(payload.TurnstileToken); token != "" {
		return token
	}
	return strings.TrimSpace(payload.CaptchaToken)
}

func (v *turnstileVerifier) verify(ctx context.Context, token, remoteIP string) error {
	form := url.Values{}
	form.Set("secret", v.secretKey)
	form.Set("response", token)
	form.Set("idempotency_key", uuid.NewString())
	if strings.TrimSpace(remoteIP) != "" {
		form.Set("remoteip", remoteIP)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, turnstileSiteVerifyURL, strings.NewReader(form.Encode()))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := v.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("turnstile verification failed")
	}

	var parsed turnstileResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return err
	}
	if !parsed.Success {
		return errors.New("turnstile token rejected")
	}

	return nil
}
