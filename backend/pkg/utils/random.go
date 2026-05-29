package utils

import (
	"github.com/google/uuid"
	"math/rand"
	"time"
)

func GenerateUUID() string {
	return uuid.New().String()
}

func GenerateRandomString(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	b := make([]byte, n)
	for i := range b {
		b[i] = charset[random.Intn(len(charset))]
	}
	return string(b)
}
