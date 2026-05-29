package utils

import (
	"database/sql/driver"
	"fmt"
	"strings"
	"time"
)

type CustomDate time.Time

func (d *CustomDate) UnmarshalParam(s string) error {
	if s == "" || s == "null" {
		return nil
	}

	formats := []string{
		"02/01/2006",
		"02-01-2006",
		"2006-01-02",
		time.RFC3339,
	}

	for _, f := range formats {
		t, err := time.Parse(f, s)
		if err == nil {
			// Check if the parsed date actually matches the input string (real calendar check)
			// time.Parse is already quite good, but this extra check ensures things like Feb 30 are rejected
			if t.Format(f) == s || f == time.RFC3339 || (f == "2006-01-02" && len(s) > 10) {
				*d = CustomDate(t)
				return nil
			}
		}
	}

	return fmt.Errorf("tanggal lahir tidak valid dalam kalender nyata")
}

func (d *CustomDate) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	return d.UnmarshalParam(s)
}

func (d CustomDate) MarshalJSON() ([]byte, error) {
	t := time.Time(d)
	if t.IsZero() {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf("\"%s\"", t.Format("2006-01-02"))), nil
}

func (d *CustomDate) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("invalid type for CustomDate: %T", value)
	}
	*d = CustomDate(t)
	return nil
}

func (d CustomDate) Value() (driver.Value, error) {
	t := time.Time(d)
	if t.IsZero() {
		return nil, nil
	}
	return t, nil
}

func (d CustomDate) Time() time.Time {
	return time.Time(d)
}

func FormatCustomDateID(d *CustomDate) string {
	if d == nil {
		return "-"
	}
	return FormatTimeID(time.Time(*d))
}
