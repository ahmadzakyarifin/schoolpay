package delivery

import (
	"crypto/subtle"
	"encoding/json"
	"net/http"

	"github.com/ahmadzakyarifin/schoolpay/config"
	financeusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	webhookusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/webhook/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/gin-gonic/gin"
)

type WebhookHandler struct {
	s   webhookusecase.WebhookService
	pay financeusecase.PaymentService
	pg  financeusecase.PaymentGatewayService
	cfg *config.Config
}

func NewWebhookHandler(s webhookusecase.WebhookService, pay financeusecase.PaymentService, pg financeusecase.PaymentGatewayService, cfg *config.Config) *WebhookHandler {
	return &WebhookHandler{s: s, pay: pay, pg: pg, cfg: cfg}
}

func (h *WebhookHandler) HandlePayment(c *gin.Context) {
	var notification map[string]interface{}
	if err := c.ShouldBindJSON(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}

	orderID, _ := notification["order_id"].(string)
	statusCode, _ := notification["status_code"].(string)
	grossAmount, _ := notification["gross_amount"].(string)
	signatureKey, _ := notification["signature_key"].(string)

	if !h.pg.VerifySignature(orderID, statusCode, grossAmount, signatureKey) {
		c.JSON(http.StatusForbidden, gin.H{"error": "invalid signature"})
		return
	}

	if err := h.pay.HandleWebhook(c.Request.Context(), notification); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "ok"})
}

func (h *WebhookHandler) HandleWAHA(c *gin.Context) {
	if h.cfg != nil && h.cfg.WAHAWebhookSecret != "" {
		secret := c.GetHeader("X-SchoolPay-Webhook-Secret")
		if subtle.ConstantTimeCompare([]byte(secret), []byte(h.cfg.WAHAWebhookSecret)) != 1 {
			helper.ErrorResponse(c, http.StatusUnauthorized, "invalid webhook secret")
			return
		}
	}

	var payload json.RawMessage
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "invalid payload")
		return
	}

	if err := h.s.HandleWAHAWebhook(c.Request.Context(), payload); err != nil {
		helper.SuccessResponse(c, http.StatusOK, "logged with error", nil)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "webhook received", nil)
}
