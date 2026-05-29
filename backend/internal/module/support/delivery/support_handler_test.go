package delivery

import (
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/usecase"
	"github.com/gin-gonic/gin"
)

func TestNewSupportHandler(t *testing.T) {
	type args struct {
		s usecase.SupportService
	}
	tests := []struct {
		name string
		args args
		want *SupportHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSupportHandler(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSupportHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSupportHandler_List(t *testing.T) {
	type fields struct {
		s usecase.SupportService
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
			h := &SupportHandler{
				s: tt.fields.s,
			}
			h.List(tt.args.c)
		})
	}
}

func TestSupportHandler_Messages(t *testing.T) {
	type fields struct {
		s usecase.SupportService
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
			h := &SupportHandler{
				s: tt.fields.s,
			}
			h.Messages(tt.args.c)
		})
	}
}

func TestSupportHandler_Reply(t *testing.T) {
	type fields struct {
		s usecase.SupportService
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
			h := &SupportHandler{
				s: tt.fields.s,
			}
			h.Reply(tt.args.c)
		})
	}
}

func TestSupportHandler_Assign(t *testing.T) {
	type fields struct {
		s usecase.SupportService
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
			h := &SupportHandler{
				s: tt.fields.s,
			}
			h.Assign(tt.args.c)
		})
	}
}

func TestSupportHandler_Close(t *testing.T) {
	type fields struct {
		s usecase.SupportService
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
			h := &SupportHandler{
				s: tt.fields.s,
			}
			h.Close(tt.args.c)
		})
	}
}
