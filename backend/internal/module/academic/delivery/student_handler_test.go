package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewStudentHandler(t *testing.T) {
	type args struct {
		s   usecase.StudentService
		cfg *config.Config
	}
	tests := []struct {
		name string
		args args
		want *StudentHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStudentHandler(tt.args.s, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStudentHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStudentHandler_GetAll(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.GetAll(tt.args.c)
		})
	}
}

func TestStudentHandler_GetParents(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.GetParents(tt.args.c)
		})
	}
}

func TestStudentHandler_GetFilters(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.GetFilters(tt.args.c)
		})
	}
}

func TestStudentHandler_Create(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Create(tt.args.c)
		})
	}
}

func TestStudentHandler_Update(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Update(tt.args.c)
		})
	}
}

func TestStudentHandler_GetByID(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.GetByID(tt.args.c)
		})
	}
}

func TestStudentHandler_ToggleStatus(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.ToggleStatus(tt.args.c)
		})
	}
}

func TestStudentHandler_Delete(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Delete(tt.args.c)
		})
	}
}

func TestStudentHandler_GetClassHistory(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.GetClassHistory(tt.args.c)
		})
	}
}

func TestStudentHandler_GetMyStudents(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.GetMyStudents(tt.args.c)
		})
	}
}

func TestStudentHandler_GetByParentID(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.GetByParentID(tt.args.c)
		})
	}
}

func TestStudentHandler_Export(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Export(tt.args.c)
		})
	}
}

func TestStudentHandler_BulkGraduate(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.BulkGraduate(tt.args.c)
		})
	}
}

func TestStudentHandler_BulkDelete(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.BulkDelete(tt.args.c)
		})
	}
}

func TestStudentHandler_Restore(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Restore(tt.args.c)
		})
	}
}

func TestStudentHandler_BulkRestore(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.BulkRestore(tt.args.c)
		})
	}
}

func TestStudentHandler_BulkPromote(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.BulkPromote(tt.args.c)
		})
	}
}

func TestStudentHandler_ensureDir(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			if err := h.ensureDir(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("StudentHandler.ensureDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestStudentHandler_GetDependencyInfo(t *testing.T) {
	type fields struct {
		s   usecase.StudentService
		cfg *config.Config
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
			h := &StudentHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.GetDependencyInfo(tt.args.c)
		})
	}
}

