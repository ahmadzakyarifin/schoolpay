package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/uptrace/bun"
)

func TestNewAcademicYearRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want AcademicYearRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAcademicYearRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAcademicYearRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_academicYearRepo_Create(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.ay); (err != nil) != tt.wantErr {
				t.Errorf("academicYearRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearRepo_FindAll(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			got, got1, err := r.FindAll(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.status, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearRepo.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("academicYearRepo.FindAll() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("academicYearRepo.FindAll() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_academicYearRepo_FindByID(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.AcademicYear
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearRepo.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("academicYearRepo.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_academicYearRepo_FindByYear(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx  context.Context
		year int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.AcademicYear
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByYear(tt.args.ctx, tt.args.year)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearRepo.FindByYear() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("academicYearRepo.FindByYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_academicYearRepo_Update(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.ctx, tt.args.ay); (err != nil) != tt.wantErr {
				t.Errorf("academicYearRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearRepo_Delete(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("academicYearRepo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearRepo_Restore(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("academicYearRepo.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearRepo_BulkDelete(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			if err := r.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("academicYearRepo.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearRepo_BulkRestore(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			if err := r.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("academicYearRepo.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearRepo_Exists(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx       context.Context
		year      int
		excludeID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			got, err := r.Exists(tt.args.ctx, tt.args.year, tt.args.excludeID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearRepo.Exists() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("academicYearRepo.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_academicYearRepo_AssignMajors(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			if err := r.AssignMajors(tt.args.ctx, tt.args.ayID, tt.args.majorIDs); (err != nil) != tt.wantErr {
				t.Errorf("academicYearRepo.AssignMajors() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearRepo_GetMajorsByYear(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			got, err := r.GetMajorsByYear(tt.args.ctx, tt.args.ayID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearRepo.GetMajorsByYear() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("academicYearRepo.GetMajorsByYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_academicYearRepo_AssignClasses(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			if err := r.AssignClasses(tt.args.ctx, tt.args.ayID, tt.args.classIDs); (err != nil) != tt.wantErr {
				t.Errorf("academicYearRepo.AssignClasses() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_academicYearRepo_GetClassesByYear(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			got, err := r.GetClassesByYear(tt.args.ctx, tt.args.ayID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearRepo.GetClassesByYear() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("academicYearRepo.GetClassesByYear() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_academicYearRepo_FindActive(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			got, err := r.FindActive(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearRepo.FindActive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("academicYearRepo.FindActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_academicYearRepo_CountStudentsByEntryYear(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx  context.Context
		year int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &academicYearRepo{
				db: tt.fields.db,
			}
			got, err := r.CountStudentsByEntryYear(tt.args.ctx, tt.args.year)
			if (err != nil) != tt.wantErr {
				t.Fatalf("academicYearRepo.CountStudentsByEntryYear() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("academicYearRepo.CountStudentsByEntryYear() = %v, want %v", got, tt.want)
			}
		})
	}
}
