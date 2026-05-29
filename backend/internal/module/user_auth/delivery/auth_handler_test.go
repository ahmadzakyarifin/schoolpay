package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewAuthHandler(t *testing.T) {
	type args struct {
		service usecase.AuthService
		cfg     *config.Config
	}
	tests := []struct {
		name string
		args args
		want *AuthHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthHandler(tt.args.service, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	type fields struct {
		s   usecase.AuthService
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
			h := &AuthHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Login(tt.args.c)
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	type fields struct {
		s   usecase.AuthService
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
			h := &AuthHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.RefreshToken(tt.args.c)
		})
	}
}

func TestAuthHandler_ForgotPassword(t *testing.T) {
	type fields struct {
		s   usecase.AuthService
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
			h := &AuthHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.ForgotPassword(tt.args.c)
		})
	}
}

func TestAuthHandler_ResetPassword(t *testing.T) {
	type fields struct {
		s   usecase.AuthService
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
			h := &AuthHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.ResetPassword(tt.args.c)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	type fields struct {
		s   usecase.AuthService
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
			h := &AuthHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.Logout(tt.args.c)
		})
	}
}

func TestAuthHandler_ChangePassword(t *testing.T) {
	type fields struct {
		s   usecase.AuthService
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
			h := &AuthHandler{
				s:   tt.fields.s,
				cfg: tt.fields.cfg,
			}
			h.ChangePassword(tt.args.c)
		})
	}
}
