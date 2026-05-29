package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/uptrace/bun"
)

func TestNewMajorRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want MajorRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewMajorRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewMajorRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_majorRepo_Create(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &majorRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.j); (err != nil) != tt.wantErr {
				t.Errorf("majorRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorRepo_FindAll(t *testing.T) {
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
		want    []domain.Major
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &majorRepo{
				db: tt.fields.db,
			}
			got, got1, err := r.FindAll(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.status, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("majorRepo.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("majorRepo.FindAll() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("majorRepo.FindAll() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_majorRepo_FindByID(t *testing.T) {
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
		want    *domain.Major
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &majorRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("majorRepo.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("majorRepo.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_majorRepo_FindByName(t *testing.T) {
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
		want    *domain.Major
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &majorRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByName(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("majorRepo.FindByName() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("majorRepo.FindByName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_majorRepo_Update(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &majorRepo{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.ctx, tt.args.j); (err != nil) != tt.wantErr {
				t.Errorf("majorRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorRepo_Delete(t *testing.T) {
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
			r := &majorRepo{
				db: tt.fields.db,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("majorRepo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorRepo_Restore(t *testing.T) {
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
			r := &majorRepo{
				db: tt.fields.db,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("majorRepo.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorRepo_Exists(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx       context.Context
		name      string
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
			r := &majorRepo{
				db: tt.fields.db,
			}
			got, err := r.Exists(tt.args.ctx, tt.args.name, tt.args.excludeID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("majorRepo.Exists() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("majorRepo.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_majorRepo_ToggleStatus(t *testing.T) {
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
			r := &majorRepo{
				db: tt.fields.db,
			}
			if err := r.ToggleStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("majorRepo.ToggleStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorRepo_BulkDelete(t *testing.T) {
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
			r := &majorRepo{
				db: tt.fields.db,
			}
			if err := r.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("majorRepo.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorRepo_BulkRestore(t *testing.T) {
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
			r := &majorRepo{
				db: tt.fields.db,
			}
			if err := r.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("majorRepo.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_majorRepo_CountClasses(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx     context.Context
		majorID uint
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
			r := &majorRepo{
				db: tt.fields.db,
			}
			got, err := r.CountClasses(tt.args.ctx, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("majorRepo.CountClasses() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("majorRepo.CountClasses() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_majorRepo_CountStudents(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx     context.Context
		majorID uint
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
			r := &majorRepo{
				db: tt.fields.db,
			}
			got, err := r.CountStudents(tt.args.ctx, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("majorRepo.CountStudents() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("majorRepo.CountStudents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_majorRepo_CountAcademicYears(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx     context.Context
		majorID uint
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
			r := &majorRepo{
				db: tt.fields.db,
			}
			got, err := r.CountAcademicYears(tt.args.ctx, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("majorRepo.CountAcademicYears() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("majorRepo.CountAcademicYears() = %v, want %v", got, tt.want)
			}
		})
	}
}
