package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"net/http"
	"net/smtp"
	"regexp"
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

var plainEmailBoldRE = regexp.MustCompile(`\*([^*\n]+)\*`)

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

	if !looksLikeHTML(body) {
		body = PlainTextEmailHTML(subject, body)
	}

	auth := smtp.PlainAuth("", m.smtpUser, m.smtpPass, m.smtpHost)
	header := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=\"utf-8\"\r\n\r\n", m.smtpUser, to, subject)
	message := []byte(header + body)
	addr := fmt.Sprintf("%s:%s", m.smtpHost, m.smtpPort)

	return smtp.SendMail(addr, auth, m.smtpUser, []string{to}, message)
}

func looksLikeHTML(body string) bool {
	trimmed := strings.TrimSpace(strings.ToLower(body))
	return strings.Contains(trimmed, "<html") ||
		strings.Contains(trimmed, "<body") ||
		strings.Contains(trimmed, "<div") ||
		strings.Contains(trimmed, "<table") ||
		strings.Contains(trimmed, "<p") ||
		strings.Contains(trimmed, "<br")
}

func PlainTextEmailHTML(subject, body string) string {
	escapedSubject := html.EscapeString(subject)
	bodyHTML := plainEmailBoldRE.ReplaceAllString(html.EscapeString(body), "<strong>$1</strong>")
	return fmt.Sprintf(`<!doctype html>
<html>
<body style="margin:0;background:#f8fafc;font-family:Arial,'Segoe UI',sans-serif;color:#0f172a;">
  <div style="max-width:640px;margin:0 auto;padding:32px 16px;">
    <div style="background:#ffffff;border:1px solid #e2e8f0;border-radius:14px;overflow:hidden;">
      <div style="padding:22px 26px;background:#4f46e5;color:#ffffff;">
        <div style="font-size:11px;font-weight:700;letter-spacing:.12em;text-transform:uppercase;opacity:.85;">SchoolPay</div>
        <h1 style="font-size:20px;line-height:1.35;margin:8px 0 0;">%s</h1>
      </div>
      <div style="padding:26px;">
        <div style="white-space:pre-line;font-size:14px;line-height:1.7;color:#334155;">%s</div>
      </div>
    </div>
    <p style="font-size:12px;line-height:1.6;color:#94a3b8;text-align:center;margin:18px 0 0;">
      Pesan otomatis dari SchoolPay. Mohon tidak membalas email ini.
    </p>
  </div>
</body>
</html>`, escapedSubject, bodyHTML)
}
