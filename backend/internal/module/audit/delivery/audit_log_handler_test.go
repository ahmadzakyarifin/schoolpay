package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewAuditLogHandler(t *testing.T) {
	type args struct {
		s usecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want *AuditLogHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuditLogHandler(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuditLogHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuditLogHandler_GetLogs(t *testing.T) {
	type fields struct {
		service usecase.AuditLogService
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
			h := &AuditLogHandler{
				service: tt.fields.service,
			}
			h.GetLogs(tt.args.c)
		})
	}
}

func TestAuditLogHandler_GetEntityLogs(t *testing.T) {
	type fields struct {
		service usecase.AuditLogService
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
			h := &AuditLogHandler{
				service: tt.fields.service,
			}
			h.GetEntityLogs(tt.args.c)
		})
	}
}
