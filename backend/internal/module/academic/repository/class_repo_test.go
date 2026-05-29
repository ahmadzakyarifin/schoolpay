package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/uptrace/bun"
)

func TestNewClassRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want ClassRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClassRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClassRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_classRepo_Create(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &classRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("classRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classRepo_FindAll(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &classRepo{
				db: tt.fields.db,
			}
			got, got1, err := r.FindAll(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.status, tt.args.majorID, tt.args.ayID, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("classRepo.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("classRepo.FindAll() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("classRepo.FindAll() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_classRepo_FindByID(t *testing.T) {
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
		want    *domain.Class
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &classRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("classRepo.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("classRepo.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_classRepo_FindByName(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx  context.Context
		name string
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
			r := &classRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByName(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("classRepo.FindByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("classRepo.FindByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_classRepo_Update(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &classRepo{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.ctx, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("classRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classRepo_Delete(t *testing.T) {
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
			r := &classRepo{
				db: tt.fields.db,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("classRepo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classRepo_Restore(t *testing.T) {
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
			r := &classRepo{
				db: tt.fields.db,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("classRepo.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classRepo_ToggleStatus(t *testing.T) {
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
			r := &classRepo{
				db: tt.fields.db,
			}
			if err := r.ToggleStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("classRepo.ToggleStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classRepo_BulkDelete(t *testing.T) {
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
			r := &classRepo{
				db: tt.fields.db,
			}
			if err := r.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("classRepo.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classRepo_BulkRestore(t *testing.T) {
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
			r := &classRepo{
				db: tt.fields.db,
			}
			if err := r.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("classRepo.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_classRepo_Exists(t *testing.T) {
	type fields struct {
		db *bun.DB
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
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &classRepo{
				db: tt.fields.db,
			}
			got, err := r.Exists(tt.args.ctx, tt.args.name, tt.args.ayID, tt.args.majorID, tt.args.excludeID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("classRepo.Exists() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("classRepo.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_classRepo_CountStudents(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx     context.Context
		classID uint
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
			r := &classRepo{
				db: tt.fields.db,
			}
			got, err := r.CountStudents(tt.args.ctx, tt.args.classID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("classRepo.CountStudents() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("classRepo.CountStudents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_classRepo_CountAcademicYears(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx     context.Context
		classID uint
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
			r := &classRepo{
				db: tt.fields.db,
			}
			got, err := r.CountAcademicYears(tt.args.ctx, tt.args.classID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("classRepo.CountAcademicYears() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("classRepo.CountAcademicYears() = %v, want %v", got, tt.want)
			}
		})
	}
}
