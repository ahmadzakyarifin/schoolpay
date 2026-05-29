package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewUserHandler(t *testing.T) {
	type args struct {
		s   usecase.UserService
		cfg *config.Config
	}
	tests := []struct {
		name string
		args args
		want *UserHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserHandler(tt.args.s, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserHandler_GetAll(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.GetAll(tt.args.c)
		})
	}
}

func TestUserHandler_GetByID(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.GetByID(tt.args.c)
		})
	}
}

func TestUserHandler_Create(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Create(tt.args.c)
		})
	}
}

func TestUserHandler_Update(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Update(tt.args.c)
		})
	}
}

func TestUserHandler_Delete(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Delete(tt.args.c)
		})
	}
}

func TestUserHandler_ToggleStatus(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.ToggleStatus(tt.args.c)
		})
	}
}

func TestUserHandler_Activate(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Activate(tt.args.c)
		})
	}
}

func TestUserHandler_ResendNotification(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.ResendNotification(tt.args.c)
		})
	}
}

func TestUserHandler_BulkResendNotification(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.BulkResendNotification(tt.args.c)
		})
	}
}

func TestUserHandler_BulkDelete(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.BulkDelete(tt.args.c)
		})
	}
}

func TestUserHandler_Restore(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Restore(tt.args.c)
		})
	}
}

func TestUserHandler_BulkRestore(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.BulkRestore(tt.args.c)
		})
	}
}

func TestUserHandler_GetNotifications(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.GetNotifications(tt.args.c)
		})
	}
}

func TestUserHandler_Export(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Export(tt.args.c)
		})
	}
}

func TestUserHandler_GetDependencyInfo(t *testing.T) {
	type fields struct {
		s   usecase.UserService
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
			h := &UserHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.GetDependencyInfo(tt.args.c)
		})
	}
}

