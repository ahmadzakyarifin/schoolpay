package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewFinanceReportHandler(t *testing.T) {
	type args struct {
		s usecase.FinanceReportService
	}
	tests := []struct {
		name string
		args args
		want *FinanceReportHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFinanceReportHandler(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFinanceReportHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFinanceReportHandler_GetArrears(t *testing.T) {
	type fields struct {
		s usecase.FinanceReportService
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
			h := &FinanceReportHandler{
				s: tt.fields.s,
			}
			h.GetArrears(tt.args.c)
		})
	}
}

func TestFinanceReportHandler_ExportTrend(t *testing.T) {
	type fields struct {
		s usecase.FinanceReportService
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
			h := &FinanceReportHandler{
				s: tt.fields.s,
			}
			h.ExportTrend(tt.args.c)
		})
	}
}
