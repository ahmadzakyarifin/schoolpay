package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewAcademicYearHandler(t *testing.T) {
	type args struct {
		s usecase.AcademicYearService
	}
	tests := []struct {
		name string
		args args
		want *AcademicYearHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAcademicYearHandler(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAcademicYearHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAcademicYearHandler_Create(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.Create(tt.args.c)
		})
	}
}

func TestAcademicYearHandler_GetAll(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.GetAll(tt.args.c)
		})
	}
}

func TestAcademicYearHandler_Update(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.Update(tt.args.c)
		})
	}
}

func TestAcademicYearHandler_Delete(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.Delete(tt.args.c)
		})
	}
}

func TestAcademicYearHandler_Restore(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.Restore(tt.args.c)
		})
	}
}

func TestAcademicYearHandler_BulkDelete(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.BulkDelete(tt.args.c)
		})
	}
}

func TestAcademicYearHandler_BulkRestore(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.BulkRestore(tt.args.c)
		})
	}
}

func TestAcademicYearHandler_GetActive(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.GetActive(tt.args.c)
		})
	}
}

func TestAcademicYearHandler_AssignMajors(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.AssignMajors(tt.args.c)
		})
	}
}

func TestAcademicYearHandler_GetMajors(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.GetMajors(tt.args.c)
		})
	}
}

func TestAcademicYearHandler_AssignClasses(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.AssignClasses(tt.args.c)
		})
	}
}

func TestAcademicYearHandler_GetClasses(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.GetClasses(tt.args.c)
		})
	}
}

func TestAcademicYearHandler_GetDependencyInfo(t *testing.T) {
	type fields struct {
		s usecase.AcademicYearService
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
			h := &AcademicYearHandler{
				s: tt.fields.s,
			}
			h.GetDependencyInfo(tt.args.c)
		})
	}
}
