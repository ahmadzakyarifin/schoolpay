package delivery

import (
	"reflect"
	"testing"

	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func TestNewWhatsAppHandler(t *testing.T) {
	type args struct {
		s    usecase.WhatsAppService
		noti notificationrepo.NotificationRepo
		msg  utils.Messenger
		db   *bun.DB
	}
	tests := []struct {
		name string
		args args
		want *WhatsAppHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWhatsAppHandler(tt.args.s, tt.args.noti, tt.args.msg, tt.args.db, nil); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWhatsAppHandler() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWhatsAppHandler_GetStatus(t *testing.T) {
	type fields struct {
		s    usecase.WhatsAppService
		noti notificationrepo.NotificationRepo
		msg  utils.Messenger
		db   *bun.DB
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
			h := &WhatsAppHandler{
				s:    tt.fields.s,
				noti: tt.fields.noti,
				msg:  tt.fields.msg,
				db:   tt.fields.db,
			}
			h.GetStatus(tt.args.c)
		})
	}
}

func TestWhatsAppHandler_GetQR(t *testing.T) {
	type fields struct {
		s    usecase.WhatsAppService
		noti notificationrepo.NotificationRepo
		msg  utils.Messenger
		db   *bun.DB
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
			h := &WhatsAppHandler{
				s:    tt.fields.s,
				noti: tt.fields.noti,
				msg:  tt.fields.msg,
				db:   tt.fields.db,
			}
			h.GetQR(tt.args.c)
		})
	}
}

func TestWhatsAppHandler_GetStats(t *testing.T) {
	type fields struct {
		s    usecase.WhatsAppService
		noti notificationrepo.NotificationRepo
		msg  utils.Messenger
		db   *bun.DB
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
			h := &WhatsAppHandler{
				s:    tt.fields.s,
				noti: tt.fields.noti,
				msg:  tt.fields.msg,
				db:   tt.fields.db,
			}
			h.GetStats(tt.args.c)
		})
	}
}

func TestWhatsAppHandler_GetLogs(t *testing.T) {
	type fields struct {
		s    usecase.WhatsAppService
		noti notificationrepo.NotificationRepo
		msg  utils.Messenger
		db   *bun.DB
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
			h := &WhatsAppHandler{
				s:    tt.fields.s,
				noti: tt.fields.noti,
				msg:  tt.fields.msg,
				db:   tt.fields.db,
			}
			h.GetLogs(tt.args.c)
		})
	}
}

func TestWhatsAppHandler_GetChatHistory(t *testing.T) {
	type fields struct {
		s    usecase.WhatsAppService
		noti notificationrepo.NotificationRepo
		msg  utils.Messenger
		db   *bun.DB
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
			h := &WhatsAppHandler{
				s:    tt.fields.s,
				noti: tt.fields.noti,
				msg:  tt.fields.msg,
				db:   tt.fields.db,
			}
			h.GetChatHistory(tt.args.c)
		})
	}
}

func TestWhatsAppHandler_SendChatMessage(t *testing.T) {
	type fields struct {
		s    usecase.WhatsAppService
		noti notificationrepo.NotificationRepo
		msg  utils.Messenger
		db   *bun.DB
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
			h := &WhatsAppHandler{
				s:    tt.fields.s,
				noti: tt.fields.noti,
				msg:  tt.fields.msg,
				db:   tt.fields.db,
			}
			h.SendChatMessage(tt.args.c)
		})
	}
}

func TestWhatsAppHandler_ResendSpecificNotification(t *testing.T) {
	type fields struct {
		s    usecase.WhatsAppService
		noti notificationrepo.NotificationRepo
		msg  utils.Messenger
		db   *bun.DB
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
			h := &WhatsAppHandler{
				s:    tt.fields.s,
				noti: tt.fields.noti,
				msg:  tt.fields.msg,
				db:   tt.fields.db,
			}
			h.ResendSpecificNotification(tt.args.c)
		})
	}
}
