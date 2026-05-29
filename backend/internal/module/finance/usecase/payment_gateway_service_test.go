package usecase

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/config"
)

func TestNewPaymentGatewayService(t *testing.T) {
	type args struct {
		cfg *config.Config
	}
	tests := []struct {
		name string
		args args
		want PaymentGatewayService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPaymentGatewayService(tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaymentGatewayService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paymentGatewayService_VerifySignature(t *testing.T) {
	type fields struct {
		cfg *config.Config
	}
	type args struct {
		orderID      string
		statusCode   string
		grossAmount  string
		signatureKey string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &paymentGatewayService{
				cfg: tt.fields.cfg,
			}
			if got := s.VerifySignature(tt.args.orderID, tt.args.statusCode, tt.args.grossAmount, tt.args.signatureKey); got != tt.want {
				t.Errorf("paymentGatewayService.VerifySignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
