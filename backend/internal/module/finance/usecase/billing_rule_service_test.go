package usecase

import (
	"context"
	"reflect"
	"testing"

	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	"github.com/uptrace/bun"
)

func TestNewBillingRuleService(t *testing.T) {
	type args struct {
		db       *bun.DB
		repo     repository.BillingRuleRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want BillingRuleService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBillingRuleService(tt.args.db, tt.args.repo, tt.args.notifSvc, tt.args.audit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBillingRuleService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billingRuleService_Create(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.BillingRuleRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		br  *domain.BillingRule
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
			s := &billingRuleService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.Create(tt.args.ctx, tt.args.br); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleService_GetAll(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.BillingRuleRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.BillingRule
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &billingRuleService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			got, err := s.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billingRuleService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billingRuleService.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billingRuleService_GetAllPaged(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.BillingRuleRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx            context.Context
		page           int
		limit          int
		search         string
		status         string
		generateStatus string
		sort           string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.BillingRule
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &billingRuleService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			got, got1, err := s.GetAllPaged(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.status, tt.args.generateStatus, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billingRuleService.GetAllPaged() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billingRuleService.GetAllPaged() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("billingRuleService.GetAllPaged() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_billingRuleService_GetByID(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.BillingRuleRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.BillingRule
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &billingRuleService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			got, err := s.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billingRuleService.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billingRuleService.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billingRuleService_Update(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.BillingRuleRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		br  *domain.BillingRule
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
			s := &billingRuleService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.Update(tt.args.ctx, tt.args.br); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleService_Delete(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.BillingRuleRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		id  uint
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
			s := &billingRuleService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleService_Restore(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.BillingRuleRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		id  uint
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
			s := &billingRuleService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleService.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleService_ToggleStatus(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.BillingRuleRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		id  uint
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
			s := &billingRuleService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.ToggleStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleService.ToggleStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleService_BulkDelete(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.BillingRuleRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		ids []uint
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
			s := &billingRuleService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleService.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleService_BulkRestore(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.BillingRuleRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		ids []uint
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
			s := &billingRuleService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			if err := s.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleService.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleService_GetDependencyInfo(t *testing.T) {
	type fields struct {
		db       *bun.DB
		repo     repository.BillingRuleRepo
		notifSvc notificationusecase.FinanceNotificationService
		audit    auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &billingRuleService{
				db:       tt.fields.db,
				repo:     tt.fields.repo,
				notifSvc: tt.fields.notifSvc,
				audit:    tt.fields.audit,
			}
			got, err := s.GetDependencyInfo(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billingRuleService.GetDependencyInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billingRuleService.GetDependencyInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
