package utils

import (
	"strings"

	"github.com/nyaruka/phonenumbers"
)

// NormalizePhoneNumber standarizes phone numbers to international format (starting with country code, e.g., 62)
func NormalizePhoneNumber(phone string) string {
	p := strings.TrimSpace(phone)
	p = strings.TrimPrefix(p, "+")

	p = strings.TrimPrefix(p, "0")

	if strings.HasPrefix(p, "620") {
		p = "62" + strings.TrimPrefix(p, "620")
	}

	// Default to 62 (Indonesia) if no common country code prefix is found
	if !strings.HasPrefix(p, "62") && !strings.HasPrefix(p, "60") && !strings.HasPrefix(p, "65") && !strings.HasPrefix(p, "1") {
		if p != "" {
			p = "62" + p
		}
	}

	return p
}

// ValidatePhoneNumber memeriksa validitas nomor telepon menggunakan standard libphonenumber Google
func ValidatePhoneNumber(phone string) bool {
	num, err := phonenumbers.Parse("+"+phone, "")
	if err != nil {
		return false
	}
	return phonenumbers.IsValidNumber(num)
}
