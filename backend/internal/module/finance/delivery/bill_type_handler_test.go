package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewBillTypeHandler(t *testing.T) {
	type args struct {
		s usecase.BillTypeService
	}
	tests := []struct {
		name string
		args args
		want *BillTypeHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBillTypeHandler(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBillTypeHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBillTypeHandler_Create(t *testing.T) {
	type fields struct {
		s usecase.BillTypeService
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
			h := &BillTypeHandler{
				s: tt.fields.s,
			}
			h.Create(tt.args.c)
		})
	}
}

func TestBillTypeHandler_GetAll(t *testing.T) {
	type fields struct {
		s usecase.BillTypeService
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
			h := &BillTypeHandler{
				s: tt.fields.s,
			}
			h.GetAll(tt.args.c)
		})
	}
}

func TestBillTypeHandler_Update(t *testing.T) {
	type fields struct {
		s usecase.BillTypeService
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
			h := &BillTypeHandler{
				s: tt.fields.s,
			}
			h.Update(tt.args.c)
		})
	}
}

func TestBillTypeHandler_Delete(t *testing.T) {
	type fields struct {
		s usecase.BillTypeService
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
			h := &BillTypeHandler{
				s: tt.fields.s,
			}
			h.Delete(tt.args.c)
		})
	}
}

func TestBillTypeHandler_Restore(t *testing.T) {
	type fields struct {
		s usecase.BillTypeService
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
			h := &BillTypeHandler{
				s: tt.fields.s,
			}
			h.Restore(tt.args.c)
		})
	}
}

func TestBillTypeHandler_ToggleStatus(t *testing.T) {
	type fields struct {
		s usecase.BillTypeService
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
			h := &BillTypeHandler{
				s: tt.fields.s,
			}
			h.ToggleStatus(tt.args.c)
		})
	}
}

func TestBillTypeHandler_BulkDelete(t *testing.T) {
	type fields struct {
		s usecase.BillTypeService
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
			h := &BillTypeHandler{
				s: tt.fields.s,
			}
			h.BulkDelete(tt.args.c)
		})
	}
}

func TestBillTypeHandler_BulkRestore(t *testing.T) {
	type fields struct {
		s usecase.BillTypeService
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
			h := &BillTypeHandler{
				s: tt.fields.s,
			}
			h.BulkRestore(tt.args.c)
		})
	}
}

func TestBillTypeHandler_GetDependencyInfo(t *testing.T) {
	type fields struct {
		s usecase.BillTypeService
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
			h := &BillTypeHandler{
				s: tt.fields.s,
			}
			h.GetDependencyInfo(tt.args.c)
		})
	}
}
