package usecase

import (
	"context"
	"encoding/json"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/config"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	supportusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/support/usecase"
	userauthdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	webhookrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/webhook/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/websocket"
)

func TestNewWebhookService(t *testing.T) {
	type args struct {
		repo     webhookrepo.WebhookRepo
		wa       notificationusecase.WhatsAppService
		notiRepo notificationrepo.NotificationRepo
		sbRepo   financerepo.StudentBillRepo
		payRepo  financerepo.PaymentRepo
		stuRepo  academicrepo.StudentRepo
		userRepo userauthrepo.UserRepo
		hub      *websocket.Hub
		support  supportusecase.SupportService
		cfg      *config.Config
	}
	tests := []struct {
		name string
		args args
		want WebhookService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWebhookService(tt.args.repo, tt.args.wa, tt.args.notiRepo, tt.args.sbRepo, tt.args.payRepo, tt.args.stuRepo, tt.args.userRepo, tt.args.hub, tt.args.support, tt.args.cfg); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebhookService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webhookService_HandleWAHAWebhook(t *testing.T) {
	type fields struct {
		repo     webhookrepo.WebhookRepo
		wa       notificationusecase.WhatsAppService
		notiRepo notificationrepo.NotificationRepo
		sbRepo   financerepo.StudentBillRepo
		payRepo  financerepo.PaymentRepo
		stuRepo  academicrepo.StudentRepo
		userRepo userauthrepo.UserRepo
		hub      *websocket.Hub
		support  supportusecase.SupportService
		cfg      *config.Config
	}
	type args struct {
		ctx     context.Context
		payload json.RawMessage
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
			s := &webhookService{
				repo:     tt.fields.repo,
				wa:       tt.fields.wa,
				notiRepo: tt.fields.notiRepo,
				sbRepo:   tt.fields.sbRepo,
				payRepo:  tt.fields.payRepo,
				stuRepo:  tt.fields.stuRepo,
				userRepo: tt.fields.userRepo,
				hub:      tt.fields.hub,
				support:  tt.fields.support,
				cfg:      tt.fields.cfg,
			}
			if err := s.HandleWAHAWebhook(tt.args.ctx, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("webhookService.HandleWAHAWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webhookService_handleIncomingMessage(t *testing.T) {
	type fields struct {
		repo     webhookrepo.WebhookRepo
		wa       notificationusecase.WhatsAppService
		notiRepo notificationrepo.NotificationRepo
		sbRepo   financerepo.StudentBillRepo
		payRepo  financerepo.PaymentRepo
		stuRepo  academicrepo.StudentRepo
		userRepo userauthrepo.UserRepo
		hub      *websocket.Hub
		support  supportusecase.SupportService
		cfg      *config.Config
	}
	type args struct {
		ctx     context.Context
		payload map[string]interface{}
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
			s := &webhookService{
				repo:     tt.fields.repo,
				wa:       tt.fields.wa,
				notiRepo: tt.fields.notiRepo,
				sbRepo:   tt.fields.sbRepo,
				payRepo:  tt.fields.payRepo,
				stuRepo:  tt.fields.stuRepo,
				userRepo: tt.fields.userRepo,
				hub:      tt.fields.hub,
				support:  tt.fields.support,
				cfg:      tt.fields.cfg,
			}
			s.handleIncomingMessage(tt.args.ctx, tt.args.payload)
		})
	}
}

func Test_webhookService_sendMenu(t *testing.T) {
	type fields struct {
		repo     webhookrepo.WebhookRepo
		wa       notificationusecase.WhatsAppService
		notiRepo notificationrepo.NotificationRepo
		sbRepo   financerepo.StudentBillRepo
		payRepo  financerepo.PaymentRepo
		stuRepo  academicrepo.StudentRepo
		userRepo userauthrepo.UserRepo
		hub      *websocket.Hub
		support  supportusecase.SupportService
		cfg      *config.Config
	}
	type args struct {
		to   string
		name string
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
			s := &webhookService{
				repo:     tt.fields.repo,
				wa:       tt.fields.wa,
				notiRepo: tt.fields.notiRepo,
				sbRepo:   tt.fields.sbRepo,
				payRepo:  tt.fields.payRepo,
				stuRepo:  tt.fields.stuRepo,
				userRepo: tt.fields.userRepo,
				hub:      tt.fields.hub,
				support:  tt.fields.support,
				cfg:      tt.fields.cfg,
			}
			s.sendMenu(tt.args.to, tt.args.name)
		})
	}
}

func Test_webhookService_handleCekTagihan(t *testing.T) {
	type fields struct {
		repo     webhookrepo.WebhookRepo
		wa       notificationusecase.WhatsAppService
		notiRepo notificationrepo.NotificationRepo
		sbRepo   financerepo.StudentBillRepo
		payRepo  financerepo.PaymentRepo
		stuRepo  academicrepo.StudentRepo
		userRepo userauthrepo.UserRepo
		hub      *websocket.Hub
		support  supportusecase.SupportService
		cfg      *config.Config
	}
	type args struct {
		ctx  context.Context
		to   string
		user *userauthdomain.User
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
			s := &webhookService{
				repo:     tt.fields.repo,
				wa:       tt.fields.wa,
				notiRepo: tt.fields.notiRepo,
				sbRepo:   tt.fields.sbRepo,
				payRepo:  tt.fields.payRepo,
				stuRepo:  tt.fields.stuRepo,
				userRepo: tt.fields.userRepo,
				hub:      tt.fields.hub,
				support:  tt.fields.support,
				cfg:      tt.fields.cfg,
			}
			s.handleCekTagihan(tt.args.ctx, tt.args.to, tt.args.user)
		})
	}
}

func Test_webhookService_handleCekTunggakan(t *testing.T) {
	type fields struct {
		repo     webhookrepo.WebhookRepo
		wa       notificationusecase.WhatsAppService
		notiRepo notificationrepo.NotificationRepo
		sbRepo   financerepo.StudentBillRepo
		payRepo  financerepo.PaymentRepo
		stuRepo  academicrepo.StudentRepo
		userRepo userauthrepo.UserRepo
		hub      *websocket.Hub
		support  supportusecase.SupportService
		cfg      *config.Config
	}
	type args struct {
		ctx  context.Context
		to   string
		user *userauthdomain.User
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
			s := &webhookService{
				repo:     tt.fields.repo,
				wa:       tt.fields.wa,
				notiRepo: tt.fields.notiRepo,
				sbRepo:   tt.fields.sbRepo,
				payRepo:  tt.fields.payRepo,
				stuRepo:  tt.fields.stuRepo,
				userRepo: tt.fields.userRepo,
				hub:      tt.fields.hub,
				support:  tt.fields.support,
				cfg:      tt.fields.cfg,
			}
			s.handleCekTunggakan(tt.args.ctx, tt.args.to, tt.args.user)
		})
	}
}

func Test_webhookService_handleCekPembayaran(t *testing.T) {
	type fields struct {
		repo     webhookrepo.WebhookRepo
		wa       notificationusecase.WhatsAppService
		notiRepo notificationrepo.NotificationRepo
		sbRepo   financerepo.StudentBillRepo
		payRepo  financerepo.PaymentRepo
		stuRepo  academicrepo.StudentRepo
		userRepo userauthrepo.UserRepo
		hub      *websocket.Hub
		support  supportusecase.SupportService
		cfg      *config.Config
	}
	type args struct {
		ctx  context.Context
		to   string
		user *userauthdomain.User
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
			s := &webhookService{
				repo:     tt.fields.repo,
				wa:       tt.fields.wa,
				notiRepo: tt.fields.notiRepo,
				sbRepo:   tt.fields.sbRepo,
				payRepo:  tt.fields.payRepo,
				stuRepo:  tt.fields.stuRepo,
				userRepo: tt.fields.userRepo,
				hub:      tt.fields.hub,
				support:  tt.fields.support,
				cfg:      tt.fields.cfg,
			}
			s.handleCekPembayaran(tt.args.ctx, tt.args.to, tt.args.user)
		})
	}
}

func Test_webhookService_sendInstruction(t *testing.T) {
	type fields struct {
		repo     webhookrepo.WebhookRepo
		wa       notificationusecase.WhatsAppService
		notiRepo notificationrepo.NotificationRepo
		sbRepo   financerepo.StudentBillRepo
		payRepo  financerepo.PaymentRepo
		stuRepo  academicrepo.StudentRepo
		userRepo userauthrepo.UserRepo
		hub      *websocket.Hub
		support  supportusecase.SupportService
		cfg      *config.Config
	}
	type args struct {
		to string
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
			s := &webhookService{
				repo:     tt.fields.repo,
				wa:       tt.fields.wa,
				notiRepo: tt.fields.notiRepo,
				sbRepo:   tt.fields.sbRepo,
				payRepo:  tt.fields.payRepo,
				stuRepo:  tt.fields.stuRepo,
				userRepo: tt.fields.userRepo,
				hub:      tt.fields.hub,
				support:  tt.fields.support,
				cfg:      tt.fields.cfg,
			}
			s.sendInstruction(tt.args.to)
		})
	}
}

