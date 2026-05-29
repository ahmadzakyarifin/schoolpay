package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/config"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/websocket"
	"github.com/uptrace/bun"
)

func TestNewPaymentService(t *testing.T) {
	type args struct {
		db       *bun.DB
		repo     repository.PaymentRepo
		sbRepo   repository.StudentBillRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		cfg      *config.Config
		hub      *websocket.Hub
		audit    auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want PaymentService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPaymentService(tt.args.db, tt.args.repo, tt.args.sbRepo, tt.args.stuRepo, tt.args.notifSvc, tt.args.cfg, tt.args.hub, tt.args.audit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaymentService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paymentService_CreateIntent(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.PaymentRepo
		sbRepo   repository.StudentBillRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		cfg      *config.Config
		hub      *websocket.Hub
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx          context.Context
		studentID    uint
		amount       float64
		depositApplied float64
		billIDs      []uint
		isBypassRule bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Payment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &paymentService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				sbRepo:   tt.fields.sbRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				cfg:      tt.fields.cfg,
				hub:      tt.fields.hub,
				audit:    tt.fields.audit,
			}
			got, err := s.CreateIntent(tt.args.ctx, tt.args.studentID, tt.args.amount, tt.args.depositApplied, tt.args.billIDs, tt.args.isBypassRule)
			if (err != nil) != tt.wantErr {
				t.Fatalf("paymentService.CreateIntent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paymentService.CreateIntent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paymentService_HandleWebhook(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.PaymentRepo
		sbRepo   repository.StudentBillRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		cfg      *config.Config
		hub      *websocket.Hub
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx     context.Context
		payload map[string]interface{}
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
			s := &paymentService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				sbRepo:   tt.fields.sbRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				cfg:      tt.fields.cfg,
				hub:      tt.fields.hub,
				audit:    tt.fields.audit,
			}
			if err := s.HandleWebhook(tt.args.ctx, tt.args.payload); (err != nil) != tt.wantErr {
				t.Errorf("paymentService.HandleWebhook() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_paymentService_GetReceipt(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.PaymentRepo
		sbRepo   repository.StudentBillRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		cfg      *config.Config
		hub      *websocket.Hub
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx       context.Context
		paymentID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Receipt
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &paymentService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				sbRepo:   tt.fields.sbRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				cfg:      tt.fields.cfg,
				hub:      tt.fields.hub,
				audit:    tt.fields.audit,
			}
			got, err := s.GetReceipt(tt.args.ctx, tt.args.paymentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("paymentService.GetReceipt() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paymentService.GetReceipt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paymentService_Process(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.PaymentRepo
		sbRepo   repository.StudentBillRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		cfg      *config.Config
		hub      *websocket.Hub
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx       context.Context
		studentID uint
		p         *domain.Payment
		billIDs   []uint
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
			s := &paymentService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				sbRepo:   tt.fields.sbRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				cfg:      tt.fields.cfg,
				hub:      tt.fields.hub,
				audit:    tt.fields.audit,
			}
			if err := s.Process(tt.args.ctx, tt.args.studentID, tt.args.p, tt.args.billIDs); (err != nil) != tt.wantErr {
				t.Errorf("paymentService.Process() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_paymentService_GetHistory(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.PaymentRepo
		sbRepo   repository.StudentBillRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		cfg      *config.Config
		hub      *websocket.Hub
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx       context.Context
		studentID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Payment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &paymentService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				sbRepo:   tt.fields.sbRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				cfg:      tt.fields.cfg,
				hub:      tt.fields.hub,
				audit:    tt.fields.audit,
			}
			got, err := s.GetHistory(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("paymentService.GetHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paymentService.GetHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paymentService_checkMidtransStatus(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.PaymentRepo
		sbRepo   repository.StudentBillRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		cfg      *config.Config
		hub      *websocket.Hub
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx     context.Context
		orderID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &paymentService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				sbRepo:   tt.fields.sbRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				cfg:      tt.fields.cfg,
				hub:      tt.fields.hub,
				audit:    tt.fields.audit,
			}
			got, got1, _, err := s.checkMidtransStatus(tt.args.ctx, tt.args.orderID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("paymentService.checkMidtransStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("paymentService.checkMidtransStatus() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("paymentService.checkMidtransStatus() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_paymentService_startReconciler(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.PaymentRepo
		sbRepo   repository.StudentBillRepo
		stuRepo  academicrepo.StudentRepo
		notifSvc notificationusecase.FinanceNotificationService
		cfg      *config.Config
		hub      *websocket.Hub
		audit    auditusecase.AuditLogService
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &paymentService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				sbRepo:   tt.fields.sbRepo,
				stuRepo:  tt.fields.stuRepo,
				notifSvc: tt.fields.notifSvc,
				cfg:      tt.fields.cfg,
				hub:      tt.fields.hub,
				audit:    tt.fields.audit,
			}
			s.startReconciler()
		})
	}
}
