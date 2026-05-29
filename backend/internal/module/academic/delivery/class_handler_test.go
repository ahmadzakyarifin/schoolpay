package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewClassHandler(t *testing.T) {
	type args struct {
		s usecase.ClassService
	}
	tests := []struct {
		name string
		args args
		want *ClassHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClassHandler(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClassHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClassHandler_Create(t *testing.T) {
	type fields struct {
		s usecase.ClassService
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
			h := &ClassHandler{
				s: tt.fields.s,
			}
			h.Create(tt.args.c)
		})
	}
}

func TestClassHandler_GetAll(t *testing.T) {
	type fields struct {
		s usecase.ClassService
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
			h := &ClassHandler{
				s: tt.fields.s,
			}
			h.GetAll(tt.args.c)
		})
	}
}

func TestClassHandler_Update(t *testing.T) {
	type fields struct {
		s usecase.ClassService
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
			h := &ClassHandler{
				s: tt.fields.s,
			}
			h.Update(tt.args.c)
		})
	}
}

func TestClassHandler_Delete(t *testing.T) {
	type fields struct {
		s usecase.ClassService
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
			h := &ClassHandler{
				s: tt.fields.s,
			}
			h.Delete(tt.args.c)
		})
	}
}

func TestClassHandler_Restore(t *testing.T) {
	type fields struct {
		s usecase.ClassService
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
			h := &ClassHandler{
				s: tt.fields.s,
			}
			h.Restore(tt.args.c)
		})
	}
}

func TestClassHandler_ToggleStatus(t *testing.T) {
	type fields struct {
		s usecase.ClassService
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
			h := &ClassHandler{
				s: tt.fields.s,
			}
			h.ToggleStatus(tt.args.c)
		})
	}
}

func TestClassHandler_BulkDelete(t *testing.T) {
	type fields struct {
		s usecase.ClassService
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
			h := &ClassHandler{
				s: tt.fields.s,
			}
			h.BulkDelete(tt.args.c)
		})
	}
}

func TestClassHandler_BulkRestore(t *testing.T) {
	type fields struct {
		s usecase.ClassService
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
			h := &ClassHandler{
				s: tt.fields.s,
			}
			h.BulkRestore(tt.args.c)
		})
	}
}

func TestClassHandler_SuggestNextName(t *testing.T) {
	type fields struct {
		s usecase.ClassService
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
			h := &ClassHandler{
				s: tt.fields.s,
			}
			h.SuggestNextName(tt.args.c)
		})
	}
}

func TestClassHandler_GetDependencyInfo(t *testing.T) {
	type fields struct {
		s usecase.ClassService
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
			h := &ClassHandler{
				s: tt.fields.s,
			}
			h.GetDependencyInfo(tt.args.c)
		})
	}
}
