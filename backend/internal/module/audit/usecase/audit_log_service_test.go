package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/audit/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/audit/repository"
	userdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/uptrace/bun"
)

func TestNewAuditLogService(t *testing.T) {
	type args struct {
		repo repository.AuditLogRepo
	}
	tests := []struct {
		name string
		args args
		want AuditLogService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuditLogService(tt.args.repo); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuditLogService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_auditLogService_Log(t *testing.T) {
	type fields struct {
		repo repository.AuditLogRepo
	}
	type args struct {
		ctx        context.Context
		db         bun.IDB
		userID     uint
		userName   string
		role       string
		action     string
		entityType string
		entityID   uint
		oldValues  map[string]interface{}
		newValues  map[string]interface{}
		ipAddress  string
		userAgent  string
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
			s := &auditLogService{
				repo: tt.fields.repo,
			}
			if err := s.Log(tt.args.ctx, tt.args.db, tt.args.userID, tt.args.userName, tt.args.role, tt.args.action, tt.args.entityType, tt.args.entityID, tt.args.oldValues, tt.args.newValues, tt.args.ipAddress, tt.args.userAgent); (err != nil) != tt.wantErr {
				t.Errorf("auditLogService.Log() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_auditLogService_GetLogs(t *testing.T) {
	type fields struct {
		repo repository.AuditLogRepo
	}
	type args struct {
		ctx         context.Context
		currentUser *userdomain.User
		filter      map[string]interface{}
		page        int
		limit       int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.AuditLog
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditLogService{
				repo: tt.fields.repo,
			}
			got, got1, err := s.GetLogs(tt.args.ctx, tt.args.currentUser, tt.args.filter, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Fatalf("auditLogService.GetLogs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("auditLogService.GetLogs() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("auditLogService.GetLogs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_auditLogService_GetEntityLogs(t *testing.T) {
	type fields struct {
		repo repository.AuditLogRepo
	}
	type args struct {
		ctx         context.Context
		currentUser *userdomain.User
		entityType  string
		entityID    uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.AuditLog
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &auditLogService{
				repo: tt.fields.repo,
			}
			got, err := s.GetEntityLogs(tt.args.ctx, tt.args.currentUser, tt.args.entityType, tt.args.entityID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("auditLogService.GetEntityLogs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("auditLogService.GetEntityLogs() = %v, want %v", got, tt.want)
			}
		})
	}
}
