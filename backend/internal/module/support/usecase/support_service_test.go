package usecase

import (
	"context"
	"reflect"
	"testing"

	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/repository"
	userdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	userrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/uptrace/bun"
)

func TestNewSupportService(t *testing.T) {
	type args struct {
		db    *bun.DB
		repo  repository.SupportRepo
		users userrepo.UserRepo
		wa    notificationusecase.WhatsAppService
		audit auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want SupportService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSupportService(tt.args.db, tt.args.repo, tt.args.users, tt.args.wa, tt.args.audit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSupportService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_normalizePhone(t *testing.T) {
	type args struct {
		phone string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := normalizePhone(tt.args.phone); got != tt.want {
				t.Errorf("normalizePhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_supportService_RecordIncoming(t *testing.T) {
	type fields struct {
		db    *bun.DB
		repo  repository.SupportRepo
		users userrepo.UserRepo
		wa    notificationusecase.WhatsAppService
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx     context.Context
		phone   string
		message string
		parent  *userdomain.User
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Conversation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &supportService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				users: tt.fields.users,
				wa:    tt.fields.wa,
				audit: tt.fields.audit,
			}
			got, err := s.RecordIncoming(tt.args.ctx, tt.args.phone, tt.args.message, tt.args.parent)
			if (err != nil) != tt.wantErr {
				t.Fatalf("supportService.RecordIncoming() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("supportService.RecordIncoming() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_supportService_List(t *testing.T) {
	type fields struct {
		db    *bun.DB
		repo  repository.SupportRepo
		users userrepo.UserRepo
		wa    notificationusecase.WhatsAppService
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx    context.Context
		status string
		page   int
		limit  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Conversation
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &supportService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				users: tt.fields.users,
				wa:    tt.fields.wa,
				audit: tt.fields.audit,
			}
			got, got1, err := s.List(tt.args.ctx, tt.args.status, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Fatalf("supportService.List() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("supportService.List() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("supportService.List() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_supportService_Messages(t *testing.T) {
	type fields struct {
		db    *bun.DB
		repo  repository.SupportRepo
		users userrepo.UserRepo
		wa    notificationusecase.WhatsAppService
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx            context.Context
		conversationID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Message
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &supportService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				users: tt.fields.users,
				wa:    tt.fields.wa,
				audit: tt.fields.audit,
			}
			got, err := s.Messages(tt.args.ctx, tt.args.conversationID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("supportService.Messages() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("supportService.Messages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_supportService_Reply(t *testing.T) {
	type fields struct {
		db    *bun.DB
		repo  repository.SupportRepo
		users userrepo.UserRepo
		wa    notificationusecase.WhatsAppService
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx            context.Context
		conversationID uint
		adminID        uint
		message        string
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
			s := &supportService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				users: tt.fields.users,
				wa:    tt.fields.wa,
				audit: tt.fields.audit,
			}
			if err := s.Reply(tt.args.ctx, tt.args.conversationID, tt.args.adminID, tt.args.message); (err != nil) != tt.wantErr {
				t.Errorf("supportService.Reply() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_supportService_Assign(t *testing.T) {
	type fields struct {
		db    *bun.DB
		repo  repository.SupportRepo
		users userrepo.UserRepo
		wa    notificationusecase.WhatsAppService
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx            context.Context
		conversationID uint
		adminID        uint
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
			s := &supportService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				users: tt.fields.users,
				wa:    tt.fields.wa,
				audit: tt.fields.audit,
			}
			if err := s.Assign(tt.args.ctx, tt.args.conversationID, tt.args.adminID); (err != nil) != tt.wantErr {
				t.Errorf("supportService.Assign() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_supportService_Close(t *testing.T) {
	type fields struct {
		db    *bun.DB
		repo  repository.SupportRepo
		users userrepo.UserRepo
		wa    notificationusecase.WhatsAppService
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx            context.Context
		conversationID uint
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
			s := &supportService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				users: tt.fields.users,
				wa:    tt.fields.wa,
				audit: tt.fields.audit,
			}
			if err := s.Close(tt.args.ctx, tt.args.conversationID); (err != nil) != tt.wantErr {
				t.Errorf("supportService.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
