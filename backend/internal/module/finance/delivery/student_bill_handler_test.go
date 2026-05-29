package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewStudentBillHandler(t *testing.T) {
	type args struct {
		s      usecase.StudentBillService
		paySvc usecase.PaymentService
	}
	tests := []struct {
		name string
		args args
		want *StudentBillHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStudentBillHandler(tt.args.s, tt.args.paySvc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStudentBillHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStudentBillHandler_GetAll(t *testing.T) {
	type fields struct {
		s      usecase.StudentBillService
		paySvc usecase.PaymentService
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
			h := &StudentBillHandler{
				s:      tt.fields.s,
				paySvc: tt.fields.paySvc,
			}
			h.GetAll(tt.args.c)
		})
	}
}

func TestStudentBillHandler_GetMyBills(t *testing.T) {
	type fields struct {
		s      usecase.StudentBillService
		paySvc usecase.PaymentService
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
			h := &StudentBillHandler{
				s:      tt.fields.s,
				paySvc: tt.fields.paySvc,
			}
			h.GetMyBills(tt.args.c)
		})
	}
}

func TestStudentBillHandler_Create(t *testing.T) {
	type fields struct {
		s      usecase.StudentBillService
		paySvc usecase.PaymentService
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
			h := &StudentBillHandler{
				s:      tt.fields.s,
				paySvc: tt.fields.paySvc,
			}
			h.Create(tt.args.c)
		})
	}
}

func TestStudentBillHandler_Update(t *testing.T) {
	type fields struct {
		s      usecase.StudentBillService
		paySvc usecase.PaymentService
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
			h := &StudentBillHandler{
				s:      tt.fields.s,
				paySvc: tt.fields.paySvc,
			}
			h.Update(tt.args.c)
		})
	}
}

func TestStudentBillHandler_Delete(t *testing.T) {
	type fields struct {
		s      usecase.StudentBillService
		paySvc usecase.PaymentService
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
			h := &StudentBillHandler{
				s:      tt.fields.s,
				paySvc: tt.fields.paySvc,
			}
			h.Delete(tt.args.c)
		})
	}
}

func TestStudentBillHandler_Generate(t *testing.T) {
	type fields struct {
		s      usecase.StudentBillService
		paySvc usecase.PaymentService
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
			h := &StudentBillHandler{
				s:      tt.fields.s,
				paySvc: tt.fields.paySvc,
			}
			h.Generate(tt.args.c)
		})
	}
}

func TestStudentBillHandler_BulkGenerate(t *testing.T) {
	type fields struct {
		s      usecase.StudentBillService
		paySvc usecase.PaymentService
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
			h := &StudentBillHandler{
				s:      tt.fields.s,
				paySvc: tt.fields.paySvc,
			}
			h.BulkGenerate(tt.args.c)
		})
	}
}

func TestStudentBillHandler_BulkCancel(t *testing.T) {
	type fields struct {
		s      usecase.StudentBillService
		paySvc usecase.PaymentService
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
			h := &StudentBillHandler{
				s:      tt.fields.s,
				paySvc: tt.fields.paySvc,
			}
			h.BulkCancel(tt.args.c)
		})
	}
}

func TestStudentBillHandler_Remind(t *testing.T) {
	type fields struct {
		s      usecase.StudentBillService
		paySvc usecase.PaymentService
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
			h := &StudentBillHandler{
				s:      tt.fields.s,
				paySvc: tt.fields.paySvc,
			}
			h.Remind(tt.args.c)
		})
	}
}

func TestStudentBillHandler_MarkAsPaidManual(t *testing.T) {
	type fields struct {
		s      usecase.StudentBillService
		paySvc usecase.PaymentService
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
			h := &StudentBillHandler{
				s:      tt.fields.s,
				paySvc: tt.fields.paySvc,
			}
			h.MarkAsPaidManual(tt.args.c)
		})
	}
}
