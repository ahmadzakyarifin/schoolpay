package usecase

import (
	"crypto/sha512"
	"encoding/hex"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/config"
)

type PaymentGatewayService interface {
	VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool
}

type paymentGatewayService struct {
	cfg *config.Config
}

func NewPaymentGatewayService(cfg *config.Config) PaymentGatewayService {
	return &paymentGatewayService{cfg: cfg}
}

func (s *paymentGatewayService) VerifySignature(orderID, statusCode, grossAmount, signatureKey string) bool {
	payload := orderID + statusCode + grossAmount + s.cfg.MidtransServerKey
	hash := sha512.Sum512([]byte(payload))
	expected := hex.EncodeToString(hash[:])

	return strings.EqualFold(expected, signatureKey)
}
