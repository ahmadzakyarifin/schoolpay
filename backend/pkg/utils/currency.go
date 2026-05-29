package utils

import (
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func FormatCurrency(amount float64) string {
	p := message.NewPrinter(language.Indonesian)
	return p.Sprintf("Rp %.0f", amount)
}
