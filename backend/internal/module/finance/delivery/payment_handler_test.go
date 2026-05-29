package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewPaymentHandler(t *testing.T) {
	type args struct {
		s usecase.PaymentService
	}
	tests := []struct {
		name string
		args args
		want *PaymentHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPaymentHandler(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaymentHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaymentHandler_Process(t *testing.T) {
	type fields struct {
		s usecase.PaymentService
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
			h := &PaymentHandler{
				s: tt.fields.s,
			}
			h.Process(tt.args.c)
		})
	}
}

func TestPaymentHandler_CreateIntent(t *testing.T) {
	type fields struct {
		s usecase.PaymentService
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
			h := &PaymentHandler{
				s: tt.fields.s,
			}
			h.CreateIntent(tt.args.c)
		})
	}
}

func TestPaymentHandler_GetReceipt(t *testing.T) {
	type fields struct {
		s usecase.PaymentService
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
			h := &PaymentHandler{
				s: tt.fields.s,
			}
			h.GetReceipt(tt.args.c)
		})
	}
}

func TestPaymentHandler_GetHistory(t *testing.T) {
	type fields struct {
		s usecase.PaymentService
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
			h := &PaymentHandler{
				s: tt.fields.s,
			}
			h.GetHistory(tt.args.c)
		})
	}
}

func TestPaymentHandler_HandleWebhook(t *testing.T) {
	type fields struct {
		s usecase.PaymentService
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
			h := &PaymentHandler{
				s: tt.fields.s,
			}
			h.HandleWebhook(tt.args.c)
		})
	}
}
