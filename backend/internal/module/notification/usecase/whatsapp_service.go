package usecase

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type WhatsAppService interface {
	GetStatus() (string, error)
	GetQR() ([]byte, error)
	StartSession() error
	StopSession() error
	LogoutSession() error
	RegisterWebhook() error
	GetChatHistory(phone string) (interface{}, error)
	SendChatMessage(phone, message string) error
}

type whatsappService struct{}

func NewWhatsAppService() WhatsAppService {
	return &whatsappService{}
}

func (s *whatsappService) GetStatus() (string, error) {
	wahaURL := os.Getenv("WAHA_URL")
	apiKey := os.Getenv("WAHA_API_KEY")

	fmt.Printf("[WA-DEBUG] Checking status at %s/api/sessions/default\n", wahaURL)

	client := &http.Client{Timeout: 8 * time.Second}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/sessions/default", wahaURL), nil)
	if err != nil {
		return "OFFLINE", err
	}

	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey) // Try both headers
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("[WA-DEBUG] Connection error: %v\n", err)
		return "OFFLINE", err
	}
	defer resp.Body.Close()

	fmt.Printf("[WA-DEBUG] WAHA Response Status: %d\n", resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		return "OFFLINE", nil
	}

	var result struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "OFFLINE", err
	}

	// Normalize status
	status := result.Status
	if status == "" {
		status = "OFFLINE"
	}

	// Auto-start if stopped or failed
	if status == "STOPPED" || status == "FAILED" {
		fmt.Printf("[WA-DEBUG] Session is %s, auto-restarting...\n", status)
		if status == "FAILED" {
			s.StopSession()
		}
		if err := s.StartSession(); err == nil {
			go s.RegisterWebhook()
		}
		return "STARTING", nil
	}

	return status, nil
}

func (s *whatsappService) StartSession() error {
	wahaURL := os.Getenv("WAHA_URL")
	apiKey := os.Getenv("WAHA_API_KEY")

	client := &http.Client{Timeout: 10 * time.Second}
	// Try to start the default session
	payload := strings.NewReader(`{"name": "default", "config": {"proxy": null}}`)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sessions/start", wahaURL), payload)
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to start WAHA session, status code: %d", resp.StatusCode)
	}

	return nil
}

func (s *whatsappService) StopSession() error {
	wahaURL := os.Getenv("WAHA_URL")
	apiKey := os.Getenv("WAHA_API_KEY")

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sessions/stop", wahaURL), strings.NewReader(`{"name": "default"}`))
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}

func (s *whatsappService) LogoutSession() error {
	wahaURL := os.Getenv("WAHA_URL")
	apiKey := os.Getenv("WAHA_API_KEY")

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sessions/default/logout", wahaURL), nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode != http.StatusNotFound {
		return fmt.Errorf("failed to logout WAHA session, status code: %d", resp.StatusCode)
	}

	fallbackPayload := strings.NewReader(`{"name":"default"}`)
	fallbackReq, err := http.NewRequest("POST", fmt.Sprintf("%s/api/sessions/logout", wahaURL), fallbackPayload)
	if err != nil {
		return err
	}
	fallbackReq.Header.Set("X-Api-Key", apiKey)
	fallbackReq.Header.Set("Authorization", "Bearer "+apiKey)
	fallbackReq.Header.Set("Content-Type", "application/json")
	fallbackResp, err := client.Do(fallbackReq)
	if err != nil {
		return err
	}
	defer fallbackResp.Body.Close()
	if fallbackResp.StatusCode >= 400 {
		return fmt.Errorf("failed to logout WAHA session, status code: %d", fallbackResp.StatusCode)
	}
	return nil
}

func (s *whatsappService) GetQR() ([]byte, error) {
	wahaURL := os.Getenv("WAHA_URL")
	apiKey := os.Getenv("WAHA_API_KEY")

	// 1. Check status first
	status, _ := s.GetStatus()
	switch status {
	case "WORKING", "CONNECTED":
		return nil, fmt.Errorf("ALREADY_CONNECTED")
	case "OFFLINE", "STOPPED", "FAILED":
		fmt.Println("[WA-DEBUG] Session is stopped/offline, starting it first...")
		if err := s.StartSession(); err == nil {
			go s.RegisterWebhook()
		}
		time.Sleep(2 * time.Second)
	}

	client := &http.Client{}
	var qrBody []byte
	for i := 0; i < 20; i++ {
		if i == 5 {
			fmt.Println("[WA-DEBUG] QR still not ready, attempting session RESTART...")
			s.StopSession()
			time.Sleep(2 * time.Second)
			if err := s.StartSession(); err == nil {
				go s.RegisterWebhook()
			}
			time.Sleep(3 * time.Second)
		}

		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/default/auth/qr?format=image", wahaURL), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-Api-Key", apiKey)
		req.Header.Set("Authorization", "Bearer "+apiKey)

		resp, err := client.Do(req)
		if err == nil {
			if resp.StatusCode == http.StatusOK {
				qrBody, _ = io.ReadAll(resp.Body)
				resp.Body.Close()
				fmt.Println("[WA-DEBUG] QR Code successfully fetched!")
				break
			}
			resp.Body.Close()
		}
		fmt.Printf("[WA-DEBUG] QR not ready yet (attempt %d/20), waiting...\n", i+1)
		time.Sleep(1 * time.Second)
	}

	if qrBody == nil {
		return nil, fmt.Errorf("QR code not ready after multiple attempts")
	}

	return qrBody, nil
}

func (s *whatsappService) RegisterWebhook() error {
	wahaURL := os.Getenv("WAHA_URL")
	apiKey := os.Getenv("WAHA_API_KEY")
	webhookURL := os.Getenv("WHATSAPP_HOOK_URL")
	if webhookURL == "" {
		webhookURL = "http://schoolpay_be:8080/wa-webhook"
	}

	client := &http.Client{Timeout: 5 * time.Second}

	// Retry up to 5 times with delay
	for i := 0; i < 5; i++ {
		fmt.Printf("[WA-DEBUG] Registering Webhook to %s (attempt %d/5)...\n", webhookURL, i+1)

		payload := strings.NewReader(fmt.Sprintf(`{
			"url": "%s",
			"events": ["session.status", "message", "message.ack"],
			"hmac": null
		}`, webhookURL))

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/webhooks", wahaURL), payload)
		if err != nil {
			time.Sleep(5 * time.Second)
			continue
		}

		req.Header.Set("X-Api-Key", apiKey)
		req.Header.Set("Authorization", "Bearer "+apiKey)
		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("[WA-DEBUG] Registration failed: %v, retrying...\n", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusCreated {
			fmt.Println("[WA-DEBUG] Webhook Registered Successfully!")
			resp.Body.Close()
			return nil
		}

		if resp.StatusCode == http.StatusNotFound {
			fmt.Println("[WA-DEBUG] Webhook API not found. WAHA Core might not support it.")
			resp.Body.Close()
			return nil
		}

		resp.Body.Close()
		time.Sleep(5 * time.Second)
	}

	return fmt.Errorf("failed to register webhook after 5 attempts")
}

func (s *whatsappService) GetChatHistory(phone string) (interface{}, error) {
	wahaURL := os.Getenv("WAHA_URL")
	apiKey := os.Getenv("WAHA_API_KEY")

	phone = strings.TrimSpace(phone)
	phone = strings.TrimLeft(phone, "+")
	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}
	chatID := phone + "@c.us"

	url := fmt.Sprintf("%s/api/default/chats/%s/messages?limit=20", wahaURL, chatID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	return data, nil
}

func (s *whatsappService) SendChatMessage(phone, message string) error {
	wahaURL := os.Getenv("WAHA_URL")
	apiKey := os.Getenv("WAHA_API_KEY")
	if wahaURL == "" {
		return fmt.Errorf("WAHA_URL belum dikonfigurasi")
	}

	phone = strings.TrimSpace(phone)
	phone = strings.TrimLeft(phone, "+")
	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}
	if phone == "" || strings.TrimSpace(message) == "" {
		return fmt.Errorf("nomor dan pesan wajib diisi")
	}

	payloadBytes, err := json.Marshal(map[string]string{
		"chatId":  phone + "@c.us",
		"text":    message,
		"session": "default",
	})
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/sendText", wahaURL)
	req, err := http.NewRequest("POST", url, strings.NewReader(string(payloadBytes)))
	if err != nil {
		return err
	}

	req.Header.Set("X-Api-Key", apiKey)
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("failed to send message, status code: %d", resp.StatusCode)
	}

	return nil
}
