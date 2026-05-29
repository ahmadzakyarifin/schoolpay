package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/audit/domain"
	"github.com/uptrace/bun"
)

func TestNewAuditLogRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want AuditLogRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuditLogRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAuditLogRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_auditLogRepo_Log(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
		al  *domain.AuditLog
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
			r := &auditLogRepo{
				db: tt.fields.db,
			}
			if err := r.Log(tt.args.ctx, tt.args.db, tt.args.al); (err != nil) != tt.wantErr {
				t.Errorf("auditLogRepo.Log() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_auditLogRepo_FindAll(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx    context.Context
		filter map[string]interface{}
		page   int
		limit  int
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
			r := &auditLogRepo{
				db: tt.fields.db,
			}
			got, got1, err := r.FindAll(tt.args.ctx, tt.args.filter, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Fatalf("auditLogRepo.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("auditLogRepo.FindAll() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("auditLogRepo.FindAll() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_auditLogRepo_FindByEntity(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx        context.Context
		entityType string
		entityID   uint
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
			r := &auditLogRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByEntity(tt.args.ctx, tt.args.entityType, tt.args.entityID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("auditLogRepo.FindByEntity() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("auditLogRepo.FindByEntity() = %v, want %v", got, tt.want)
			}
		})
	}
}
