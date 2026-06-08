package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/usecase"
	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

func TestNewAuthHandler(t *testing.T) {
	type args struct {
		service     usecase.AuthService
		cfg         *config.Config
		redisClient *redis.Client
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
			if got := NewAuthHandler(tt.args.service, tt.args.cfg, tt.args.redisClient); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuthHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAuthHandler_Login(t *testing.T) {
	type fields struct {
		s           usecase.AuthService
		cfg         *config.Config
		redisClient *redis.Client
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
				s:           tt.fields.s,
				cfg:         tt.fields.cfg,
				redisClient: tt.fields.redisClient,
			}
			h.Login(tt.args.c)
		})
	}
}

func TestAuthHandler_RefreshToken(t *testing.T) {
	type fields struct {
		s           usecase.AuthService
		cfg         *config.Config
		redisClient *redis.Client
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
				s:           tt.fields.s,
				cfg:         tt.fields.cfg,
				redisClient: tt.fields.redisClient,
			}
			h.RefreshToken(tt.args.c)
		})
	}
}

func TestAuthHandler_ForgotPassword(t *testing.T) {
	type fields struct {
		s           usecase.AuthService
		cfg         *config.Config
		redisClient *redis.Client
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
				s:           tt.fields.s,
				cfg:         tt.fields.cfg,
				redisClient: tt.fields.redisClient,
			}
			h.ForgotPassword(tt.args.c)
		})
	}
}

func TestAuthHandler_ResetPassword(t *testing.T) {
	type fields struct {
		s           usecase.AuthService
		cfg         *config.Config
		redisClient *redis.Client
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
				s:           tt.fields.s,
				cfg:         tt.fields.cfg,
				redisClient: tt.fields.redisClient,
			}
			h.ResetPassword(tt.args.c)
		})
	}
}

func TestAuthHandler_Logout(t *testing.T) {
	type fields struct {
		s           usecase.AuthService
		cfg         *config.Config
		redisClient *redis.Client
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
				s:           tt.fields.s,
				cfg:         tt.fields.cfg,
				redisClient: tt.fields.redisClient,
			}
			h.Logout(tt.args.c)
		})
	}
}

func TestAuthHandler_ChangePassword(t *testing.T) {
	type fields struct {
		s           usecase.AuthService
		cfg         *config.Config
		redisClient *redis.Client
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
				s:           tt.fields.s,
				cfg:         tt.fields.cfg,
				redisClient: tt.fields.redisClient,
			}
			h.ChangePassword(tt.args.c)
		})
	}
}
