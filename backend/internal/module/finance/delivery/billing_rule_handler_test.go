package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewBillingRuleHandler(t *testing.T) {
	type args struct {
		s usecase.BillingRuleService
	}
	tests := []struct {
		name string
		args args
		want *BillingRuleHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBillingRuleHandler(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBillingRuleHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBillingRuleHandler_Create(t *testing.T) {
	type fields struct {
		s usecase.BillingRuleService
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
			h := &BillingRuleHandler{
				s: tt.fields.s,
			}
			h.Create(tt.args.c)
		})
	}
}

func TestBillingRuleHandler_GetAll(t *testing.T) {
	type fields struct {
		s usecase.BillingRuleService
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
			h := &BillingRuleHandler{
				s: tt.fields.s,
			}
			h.GetAll(tt.args.c)
		})
	}
}

func TestBillingRuleHandler_Update(t *testing.T) {
	type fields struct {
		s usecase.BillingRuleService
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
			h := &BillingRuleHandler{
				s: tt.fields.s,
			}
			h.Update(tt.args.c)
		})
	}
}

func TestBillingRuleHandler_Delete(t *testing.T) {
	type fields struct {
		s usecase.BillingRuleService
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
			h := &BillingRuleHandler{
				s: tt.fields.s,
			}
			h.Delete(tt.args.c)
		})
	}
}

func TestBillingRuleHandler_Restore(t *testing.T) {
	type fields struct {
		s usecase.BillingRuleService
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
			h := &BillingRuleHandler{
				s: tt.fields.s,
			}
			h.Restore(tt.args.c)
		})
	}
}

func TestBillingRuleHandler_ToggleStatus(t *testing.T) {
	type fields struct {
		s usecase.BillingRuleService
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
			h := &BillingRuleHandler{
				s: tt.fields.s,
			}
			h.ToggleStatus(tt.args.c)
		})
	}
}

func TestBillingRuleHandler_BulkDelete(t *testing.T) {
	type fields struct {
		s usecase.BillingRuleService
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
			h := &BillingRuleHandler{
				s: tt.fields.s,
			}
			h.BulkDelete(tt.args.c)
		})
	}
}

func TestBillingRuleHandler_BulkRestore(t *testing.T) {
	type fields struct {
		s usecase.BillingRuleService
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
			h := &BillingRuleHandler{
				s: tt.fields.s,
			}
			h.BulkRestore(tt.args.c)
		})
	}
}

func TestBillingRuleHandler_GetDependencyInfo(t *testing.T) {
	type fields struct {
		s usecase.BillingRuleService
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
			h := &BillingRuleHandler{
				s: tt.fields.s,
			}
			h.GetDependencyInfo(tt.args.c)
		})
	}
}
