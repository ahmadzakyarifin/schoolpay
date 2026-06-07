package delivery

import (
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewDashboardHandler(t *testing.T) {
	type args struct {
		dashSvc usecase.DashboardService
	}
	tests := []struct {
		name string
		args args
		want *DashboardHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDashboardHandler(tt.args.dashSvc)
			if (got == nil && tt.want != nil) || (got != nil && tt.want == nil) {
				t.Errorf("NewDashboardHandler() = %v, want %v", got, tt.want)
				return
			}
			if got != nil && tt.want != nil {
				if got.dashSvc != tt.want.dashSvc {
					t.Errorf("NewDashboardHandler() fields mismatch")
				}
			}
		})
	}
}

func TestDashboardHandler_GetStats(t *testing.T) {
	type fields struct {
		dashSvc usecase.DashboardService
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
			h := &DashboardHandler{
				dashSvc: tt.fields.dashSvc,
			}
			h.GetStats(tt.args.c)
		})
	}
}

func TestDashboardHandler_GetCommunicationDetails(t *testing.T) {
	type fields struct {
		dashSvc usecase.DashboardService
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
			h := &DashboardHandler{
				dashSvc: tt.fields.dashSvc,
			}
			h.GetCommunicationDetails(tt.args.c)
		})
	}
}

func TestDashboardHandler_ExportGlobalReport(t *testing.T) {
	type fields struct {
		dashSvc usecase.DashboardService
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
			h := &DashboardHandler{
				dashSvc: tt.fields.dashSvc,
			}
			h.ExportGlobalReport(tt.args.c)
		})
	}
}
