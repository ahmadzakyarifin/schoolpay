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

func TestNewAcademicYearService(t *testing.T) {
	type args struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
		audit auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want AcademicYearService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAcademicYearService(tt.args.db, tt.args.repo, tt.args.audit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAcademicYearService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_academicYearService_Create(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		ay  *domain.AcademicYear
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
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.Create(tt.args.ctx, tt.args.ay); (err != nil) != tt.wantErr {
				t.Errorf("academicYearService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearService_GetAll(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
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
		want    []domain.AcademicYear
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, got1, err := s.GetAll(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.status, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("academicYearService.GetAll() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("academicYearService.GetAll() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_academicYearService_Update(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
		ay  *domain.AcademicYear
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
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.Update(tt.args.ctx, tt.args.ay); (err != nil) != tt.wantErr {
				t.Errorf("academicYearService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearService_Delete(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
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
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("academicYearService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearService_Restore(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
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
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("academicYearService.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearService_BulkDelete(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
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
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("academicYearService.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearService_BulkRestore(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
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
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("academicYearService.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearService_GetActive(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.AcademicYear
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, err := s.GetActive(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearService.GetActive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("academicYearService.GetActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_academicYearService_AssignMajors(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx      context.Context
		ayID     uint
		majorIDs []uint
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
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.AssignMajors(tt.args.ctx, tt.args.ayID, tt.args.majorIDs); (err != nil) != tt.wantErr {
				t.Errorf("academicYearService.AssignMajors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearService_GetMajorsByYear(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx  context.Context
		ayID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Major
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, err := s.GetMajorsByYear(tt.args.ctx, tt.args.ayID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearService.GetMajorsByYear() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("academicYearService.GetMajorsByYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_academicYearService_AssignClasses(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx      context.Context
		ayID     uint
		classIDs []uint
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
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			if err := s.AssignClasses(tt.args.ctx, tt.args.ayID, tt.args.classIDs); (err != nil) != tt.wantErr {
				t.Errorf("academicYearService.AssignClasses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearService_GetClassesByYear(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
		audit auditusecase.AuditLogService
	}
	type args struct {
		ctx  context.Context
		ayID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Class
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, err := s.GetClassesByYear(tt.args.ctx, tt.args.ayID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearService.GetClassesByYear() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("academicYearService.GetClassesByYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_academicYearService_GetDependencyInfo(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.AcademicYearRepo
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
			s := &academicYearService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, err := s.GetDependencyInfo(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearService.GetDependencyInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("academicYearService.GetDependencyInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
