package usecase

import (
	"reflect"
	"testing"

	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
)

func TestNewFinanceNotificationService(t *testing.T) {
	type args struct {
		db       *bun.DB
		stuRepo  academicrepo.StudentRepo
		userRepo userauthrepo.UserRepo
		notiRepo notificationrepo.NotificationRepo
		msg      utils.Messenger
		audit    auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want FinanceNotificationService
	}{
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFinanceNotificationService(tt.args.db, tt.args.stuRepo, tt.args.userRepo, tt.args.notiRepo, tt.args.msg, tt.args.audit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFinanceNotificationService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_financeNotificationService_Notify(t *testing.T) {
	type fields struct {
		db       *bun.DB
		stuRepo  academicrepo.StudentRepo
		userRepo userauthrepo.UserRepo
		notiRepo notificationrepo.NotificationRepo
		msg      utils.Messenger
		audit    auditusecase.AuditLogService
	}
	type args struct {
		job FinanceNotifyJob
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &financeNotificationService{
				db:       tt.fields.db,
				stuRepo:  tt.fields.stuRepo,
				userRepo: tt.fields.userRepo,
				notiRepo: tt.fields.notiRepo,
				msg:      tt.fields.msg,
				audit:    tt.fields.audit,
			}
			s.Notify(tt.args.job)
		})
	}
}
