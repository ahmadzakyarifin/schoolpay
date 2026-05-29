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

func TestNewClassService(t *testing.T) {
	type args struct {
		db        bun.IDB
		repo      repository.ClassRepo
		majorRepo repository.MajorRepo
		audit     auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want ClassService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClassService(tt.args.db, tt.args.repo, tt.args.majorRepo, tt.args.audit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClassService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_classService_Create(t *testing.T) {
	type fields struct {
		db        bun.IDB
		repo      repository.ClassRepo
		majorRepo repository.MajorRepo
		audit     auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		c   *domain.Class
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
			s := &classService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				majorRepo: tt.fields.majorRepo,
				audit:     tt.fields.audit,
			}
			if err := s.Create(tt.args.ctx, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("classService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classService_GetAll(t *testing.T) {
	type fields struct {
		db        bun.IDB
		repo      repository.ClassRepo
		majorRepo repository.MajorRepo
		audit     auditusecase.AuditLogService
	}
	type args struct {
		ctx     context.Context
		page    int
		limit   int
		search  string
		status  string
		majorID string
		ayID    string
		sort    string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Class
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &classService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				majorRepo: tt.fields.majorRepo,
				audit:     tt.fields.audit,
			}
			got, got1, err := s.GetAll(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.status, tt.args.majorID, tt.args.ayID, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("classService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("classService.GetAll() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("classService.GetAll() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_classService_GetByID(t *testing.T) {
	type fields struct {
		db        bun.IDB
		repo      repository.ClassRepo
		majorRepo repository.MajorRepo
		audit     auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Class
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &classService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				majorRepo: tt.fields.majorRepo,
				audit:     tt.fields.audit,
			}
			got, err := s.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("classService.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("classService.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_classService_Update(t *testing.T) {
	type fields struct {
		db        bun.IDB
		repo      repository.ClassRepo
		majorRepo repository.MajorRepo
		audit     auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		c   *domain.Class
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
			s := &classService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				majorRepo: tt.fields.majorRepo,
				audit:     tt.fields.audit,
			}
			if err := s.Update(tt.args.ctx, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("classService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classService_Delete(t *testing.T) {
	type fields struct {
		db        bun.IDB
		repo      repository.ClassRepo
		majorRepo repository.MajorRepo
		audit     auditusecase.AuditLogService
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
			s := &classService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				majorRepo: tt.fields.majorRepo,
				audit:     tt.fields.audit,
			}
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("classService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classService_Restore(t *testing.T) {
	type fields struct {
		db        bun.IDB
		repo      repository.ClassRepo
		majorRepo repository.MajorRepo
		audit     auditusecase.AuditLogService
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
			s := &classService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				majorRepo: tt.fields.majorRepo,
				audit:     tt.fields.audit,
			}
			if err := s.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("classService.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classService_ToggleStatus(t *testing.T) {
	type fields struct {
		db        bun.IDB
		repo      repository.ClassRepo
		majorRepo repository.MajorRepo
		audit     auditusecase.AuditLogService
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
			s := &classService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				majorRepo: tt.fields.majorRepo,
				audit:     tt.fields.audit,
			}
			if err := s.ToggleStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("classService.ToggleStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classService_BulkDelete(t *testing.T) {
	type fields struct {
		db        bun.IDB
		repo      repository.ClassRepo
		majorRepo repository.MajorRepo
		audit     auditusecase.AuditLogService
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
			s := &classService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				majorRepo: tt.fields.majorRepo,
				audit:     tt.fields.audit,
			}
			if err := s.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("classService.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classService_BulkRestore(t *testing.T) {
	type fields struct {
		db        bun.IDB
		repo      repository.ClassRepo
		majorRepo repository.MajorRepo
		audit     auditusecase.AuditLogService
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
			s := &classService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				majorRepo: tt.fields.majorRepo,
				audit:     tt.fields.audit,
			}
			if err := s.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("classService.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classService_SuggestNextName(t *testing.T) {
	type fields struct {
		db        bun.IDB
		repo      repository.ClassRepo
		majorRepo repository.MajorRepo
		audit     auditusecase.AuditLogService
	}
	type args struct {
		ctx       context.Context
		name      string
		ayID      uint
		majorID   uint
		excludeID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &classService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				majorRepo: tt.fields.majorRepo,
				audit:     tt.fields.audit,
			}
			got, err := s.SuggestNextName(tt.args.ctx, tt.args.name, tt.args.ayID, tt.args.majorID, tt.args.excludeID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("classService.SuggestNextName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("classService.SuggestNextName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_classService_GetDependencyInfo(t *testing.T) {
	type fields struct {
		db        bun.IDB
		repo      repository.ClassRepo
		majorRepo repository.MajorRepo
		audit     auditusecase.AuditLogService
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
			s := &classService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				majorRepo: tt.fields.majorRepo,
				audit:     tt.fields.audit,
			}
			got, err := s.GetDependencyInfo(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("classService.GetDependencyInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("classService.GetDependencyInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
