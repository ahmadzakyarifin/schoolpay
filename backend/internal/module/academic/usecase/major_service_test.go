package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/uptrace/bun"
)

func TestNewMajorService(t *testing.T) {
	type args struct {
		db    bun.IDB
		repo  repository.MajorRepo
		audit auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want MajorService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMajorService(tt.args.db, tt.args.repo, tt.args.audit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMajorService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_majorService_Create(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.MajorRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		j   *domain.Major
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
			s := &majorService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.Create(tt.args.ctx, tt.args.j); (err != nil) != tt.wantErr {
				t.Errorf("majorService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorService_GetAll(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.MajorRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx    context.Context
		page   int
		limit  int
		search string
		status string
		sort   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Major
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &majorService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, got1, err := s.GetAll(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.status, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("majorService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("majorService.GetAll() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("majorService.GetAll() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_majorService_GetByID(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.MajorRepo
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
		want    *domain.Major
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &majorService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, err := s.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("majorService.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("majorService.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_majorService_Update(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.MajorRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		j   *domain.Major
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
			s := &majorService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.Update(tt.args.ctx, tt.args.j); (err != nil) != tt.wantErr {
				t.Errorf("majorService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorService_Delete(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.MajorRepo
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
			s := &majorService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("majorService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorService_Restore(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.MajorRepo
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
			s := &majorService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("majorService.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorService_ToggleStatus(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.MajorRepo
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
			s := &majorService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.ToggleStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("majorService.ToggleStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorService_BulkDelete(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.MajorRepo
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
			s := &majorService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("majorService.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorService_BulkRestore(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.MajorRepo
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
			s := &majorService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("majorService.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorService_GetDependencyInfo(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.MajorRepo
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
			s := &majorService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, err := s.GetDependencyInfo(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("majorService.GetDependencyInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("majorService.GetDependencyInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
