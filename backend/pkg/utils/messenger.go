package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/smtp"
	"strings"
	"time"
)

type Messenger interface {
	SendWhatsApp(phone, message string) (string, error)
	SendEmail(to, subject, body string) error
}

type messenger struct {
	wahaURL  string
	wahaKey  string
	smtpHost string
	smtpPort string
	smtpUser string
	smtpPass string
}

func NewMessenger(wahaURL, wahaKey, smtpHost, smtpPort, smtpUser, smtpPass string) Messenger {
	return &messenger{
		wahaURL:  strings.TrimSuffix(wahaURL, "/"),
		wahaKey:  wahaKey,
		smtpHost: smtpHost,
		smtpPort: smtpPort,
		smtpUser: smtpUser,
		smtpPass: smtpPass,
	}
}

func (m *messenger) SendWhatsApp(phone, message string) (string, error) {
	if m.wahaURL == "" {
		return "", nil
	}

	url := fmt.Sprintf("%s/api/sendText", m.wahaURL)
	phone = strings.TrimSpace(phone)
	phone = strings.TrimLeft(phone, "+")

	if strings.HasPrefix(phone, "0") {
		phone = "62" + phone[1:]
	}
	if !strings.Contains(phone, "@") {
		phone = phone + "@c.us"
	}

	payload := map[string]interface{}{
		"chatId":  phone,
		"text":    message,
		"session": "default",
	}

	jsonPayload, _ := json.Marshal(payload)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	if m.wahaKey != "" {
		req.Header.Set("X-Api-Key", m.wahaKey)
	}

	client := &http.Client{Timeout: 20 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("gagal koneksi ke WAHA: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var errResp map[string]interface{}
		json.NewDecoder(resp.Body).Decode(&errResp)
		return "", fmt.Errorf("waha error %d: %v", resp.StatusCode, errResp)
	}

	var result struct {
		ID string `json:"id"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	return result.ID, nil
}

func (m *messenger) SendEmail(to, subject, body string) error {
	if m.smtpHost == "" {
		return nil
	}

	auth := smtp.PlainAuth("", m.smtpUser, m.smtpPass, m.smtpHost)
	header := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n", m.smtpUser, to, subject)
	message := []byte(header + body)
	addr := fmt.Sprintf("%s:%s", m.smtpHost, m.smtpPort)

	return smtp.SendMail(addr, auth, m.smtpUser, []string{to}, message)
}
