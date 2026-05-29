package delivery

import (
	"reflect"
	"testing"

	financeusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	webhookusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/webhook/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewWebhookHandler(t *testing.T) {
	type args struct {
		s   webhookusecase.WebhookService
		pay financeusecase.PaymentService
		pg  financeusecase.PaymentGatewayService
	}
	tests := []struct {
		name string
		args args
		want *WebhookHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWebhookHandler(tt.args.s, tt.args.pay, tt.args.pg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebhookHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWebhookHandler_HandlePayment(t *testing.T) {
	type fields struct {
		s   webhookusecase.WebhookService
		pay financeusecase.PaymentService
		pg  financeusecase.PaymentGatewayService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &WebhookHandler{
				s:   tt.fields.s,
				pay: tt.fields.pay,
				pg:  tt.fields.pg,
			}
			h.HandlePayment(tt.args.c)
		})
	}
}

func TestWebhookHandler_HandleWAHA(t *testing.T) {
	type fields struct {
		s   webhookusecase.WebhookService
		pay financeusecase.PaymentService
		pg  financeusecase.PaymentGatewayService
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &WebhookHandler{
				s:   tt.fields.s,
				pay: tt.fields.pay,
				pg:  tt.fields.pg,
			}
			h.HandleWAHA(tt.args.c)
		})
	}
}
