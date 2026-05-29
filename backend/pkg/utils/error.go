package utils

import "strings"

type BusinessError struct {
	Field   string
	Message string
}

func (e *BusinessError) Error() string {
	return e.Message
}

func NewBusinessError(field, message string) error {
	return &BusinessError{
		Field:   field,
		Message: message,
	}
}

type BusinessMultiError struct {
	Errors map[string][]string
}

func (e *BusinessMultiError) Error() string {
	if len(e.Errors) == 0 {
		return "validasi bisnis gagal"
	}
	var msgs []string
	for _, fieldErrs := range e.Errors {
		msgs = append(msgs, fieldErrs...)
	}
	return strings.Join(msgs, "; ")
}

func NewBusinessMultiError(errs map[string][]string) error {
	return &BusinessMultiError{Errors: errs}
}
