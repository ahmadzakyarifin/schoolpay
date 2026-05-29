package usecase

import (
	"context"
	"reflect"
	"testing"

	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	"github.com/uptrace/bun"
)

func TestNewBillTypeService(t *testing.T) {
	type args struct {
		db    bun.IDB
		repo  repository.BillTypeRepo
		audit auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want BillTypeService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBillTypeService(tt.args.db, tt.args.repo, tt.args.audit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBillTypeService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billTypeService_Create(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.BillTypeRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		b   *domain.BillType
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
			s := &billTypeService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.Create(tt.args.ctx, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("billTypeService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeService_GetAll(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.BillTypeRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.BillType
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &billTypeService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, err := s.GetAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billTypeService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billTypeService.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billTypeService_GetAllPaged(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.BillTypeRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx        context.Context
		page       int
		limit      int
		search     string
		filterType string
		status     string
		sort       string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.BillType
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &billTypeService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, got1, err := s.GetAllPaged(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.filterType, tt.args.status, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billTypeService.GetAllPaged() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billTypeService.GetAllPaged() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("billTypeService.GetAllPaged() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_billTypeService_GetByID(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.BillTypeRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.BillType
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &billTypeService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, err := s.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billTypeService.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billTypeService.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billTypeService_Update(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.BillTypeRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		b   *domain.BillType
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
			s := &billTypeService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.Update(tt.args.ctx, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("billTypeService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeService_Delete(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.BillTypeRepo
		audit auditusecase.AuditLogService
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
			s := &billTypeService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("billTypeService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeService_Restore(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.BillTypeRepo
		audit auditusecase.AuditLogService
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
			s := &billTypeService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("billTypeService.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeService_ToggleStatus(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.BillTypeRepo
		audit auditusecase.AuditLogService
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
			s := &billTypeService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.ToggleStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("billTypeService.ToggleStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeService_BulkDelete(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.BillTypeRepo
		audit auditusecase.AuditLogService
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
			s := &billTypeService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("billTypeService.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeService_BulkRestore(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.BillTypeRepo
		audit auditusecase.AuditLogService
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
			s := &billTypeService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("billTypeService.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeService_GetDependencyInfo(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.BillTypeRepo
		audit auditusecase.AuditLogService
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
			s := &billTypeService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, err := s.GetDependencyInfo(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billTypeService.GetDependencyInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billTypeService.GetDependencyInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
