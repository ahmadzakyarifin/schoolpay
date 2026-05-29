package utils

import (
	"fmt"
	"time"
)

func FormatDateID(dateStr *string) string {
	if dateStr == nil || *dateStr == "" {
		return "-"
	}

	formats := []string{"2006-01-02", "2006-01-02T15:04:05Z", "02/01/2006", "02-01-2006"}
	var t time.Time
	var err error
	for _, f := range formats {
		t, err = time.Parse(f, *dateStr)
		if err == nil {
			break
		}
	}

	if err != nil {
		return *dateStr
	}

	return FormatTimeID(t)
}

func FormatTimeID(t time.Time) string {
	if t.IsZero() {
		return "-"
	}
	return fmt.Sprintf("%02d/%02d/%d", t.Day(), int(t.Month()), t.Year())
}

func IsValidDate(dateStr string) bool {
	if dateStr == "" {
		return true
	}

	formats := []string{"2006-01-02", "02/01/2006", "02-01-2006"}
	for _, f := range formats {
		t, err := time.Parse(f, dateStr)
		if err == nil {
			// Check if it's a real date (e.g. not Feb 30)
			if t.Format(f) == dateStr {
				return true
			}
		}
	}
	return false
}

// ParseDateToDB converts various date formats to YYYY-MM-DD for database storage.
func ParseDateToDB(dateStr string) string {
	if dateStr == "" {
		return ""
	}
	formats := []string{"2006-01-02", "02/01/2006", "02-01-2006", "2006-01-02T15:04:05Z"}
	for _, f := range formats {
		t, err := time.Parse(f, dateStr)
		if err == nil {
			return t.Format("2006-01-02")
		}
	}
	return dateStr
}
