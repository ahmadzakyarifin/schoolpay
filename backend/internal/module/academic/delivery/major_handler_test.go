package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewMajorHandler(t *testing.T) {
	type args struct {
		s usecase.MajorService
	}
	tests := []struct {
		name string
		args args
		want *MajorHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMajorHandler(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMajorHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMajorHandler_Create(t *testing.T) {
	type fields struct {
		s usecase.MajorService
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
			h := &MajorHandler{
				s: tt.fields.s,
			}
			h.Create(tt.args.c)
		})
	}
}

func TestMajorHandler_GetAll(t *testing.T) {
	type fields struct {
		s usecase.MajorService
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
			h := &MajorHandler{
				s: tt.fields.s,
			}
			h.GetAll(tt.args.c)
		})
	}
}

func TestMajorHandler_Update(t *testing.T) {
	type fields struct {
		s usecase.MajorService
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
			h := &MajorHandler{
				s: tt.fields.s,
			}
			h.Update(tt.args.c)
		})
	}
}

func TestMajorHandler_Delete(t *testing.T) {
	type fields struct {
		s usecase.MajorService
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
			h := &MajorHandler{
				s: tt.fields.s,
			}
			h.Delete(tt.args.c)
		})
	}
}

func TestMajorHandler_Restore(t *testing.T) {
	type fields struct {
		s usecase.MajorService
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
			h := &MajorHandler{
				s: tt.fields.s,
			}
			h.Restore(tt.args.c)
		})
	}
}

func TestMajorHandler_ToggleStatus(t *testing.T) {
	type fields struct {
		s usecase.MajorService
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
			h := &MajorHandler{
				s: tt.fields.s,
			}
			h.ToggleStatus(tt.args.c)
		})
	}
}

func TestMajorHandler_BulkDelete(t *testing.T) {
	type fields struct {
		s usecase.MajorService
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
			h := &MajorHandler{
				s: tt.fields.s,
			}
			h.BulkDelete(tt.args.c)
		})
	}
}

func TestMajorHandler_BulkRestore(t *testing.T) {
	type fields struct {
		s usecase.MajorService
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
			h := &MajorHandler{
				s: tt.fields.s,
			}
			h.BulkRestore(tt.args.c)
		})
	}
}

func TestMajorHandler_GetDependencyInfo(t *testing.T) {
	type fields struct {
		s usecase.MajorService
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
			h := &MajorHandler{
				s: tt.fields.s,
			}
			h.GetDependencyInfo(tt.args.c)
		})
	}
}
