package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/uptrace/bun"
)

func TestNewUserRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want UserRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_Create(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
		u   *domain.User
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
			r := &userRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.db, tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("userRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepo_FindAll(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx  context.Context
		role string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepo{
				db: tt.fields.db,
			}
			got, err := r.FindAll(tt.args.ctx, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userRepo.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_FindByID(t *testing.T) {
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
		want    *domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userRepo.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_FindByEmail(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userRepo.FindByEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.FindByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_FindByPhone(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByPhone(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userRepo.FindByPhone() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.FindByPhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_FindByNIK(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		nik string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByNIK(tt.args.ctx, tt.args.nik)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userRepo.FindByNIK() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.FindByNIK() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_Update(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
		u   *domain.User
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
			r := &userRepo{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.ctx, tt.args.db, tt.args.u); (err != nil) != tt.wantErr {
				t.Errorf("userRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepo_Delete(t *testing.T) {
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
			r := &userRepo{
				db: tt.fields.db,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("userRepo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepo_BulkDelete(t *testing.T) {
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
			r := &userRepo{
				db: tt.fields.db,
			}
			if err := r.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("userRepo.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepo_Restore(t *testing.T) {
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
			r := &userRepo{
				db: tt.fields.db,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("userRepo.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepo_BulkRestore(t *testing.T) {
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
			r := &userRepo{
				db: tt.fields.db,
			}
			if err := r.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("userRepo.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepo_UpdatePassword(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx  context.Context
		id   uint
		hash string
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
			r := &userRepo{
				db: tt.fields.db,
			}
			if err := r.UpdatePassword(tt.args.ctx, tt.args.id, tt.args.hash); (err != nil) != tt.wantErr {
				t.Errorf("userRepo.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepo_Activate(t *testing.T) {
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
			r := &userRepo{
				db: tt.fields.db,
			}
			if err := r.Activate(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("userRepo.Activate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepo_ToggleStatus(t *testing.T) {
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
			r := &userRepo{
				db: tt.fields.db,
			}
			if err := r.ToggleStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("userRepo.ToggleStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userRepo_FindPaginated(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx    context.Context
		page   int
		limit  int
		search string
		role   string
		filter string
		status string
		sort   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.User
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepo{
				db: tt.fields.db,
			}
			got, got1, err := r.FindPaginated(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.role, tt.args.filter, tt.args.status, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userRepo.FindPaginated() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.FindPaginated() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("userRepo.FindPaginated() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_userRepo_CountActive(t *testing.T) {
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
		want    int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepo{
				db: tt.fields.db,
			}
			got, err := r.CountActive(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userRepo.CountActive() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("userRepo.CountActive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_CountActiveByPeriod(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx   context.Context
		start *time.Time
		end   *time.Time
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
			r := &userRepo{
				db: tt.fields.db,
			}
			got, err := r.CountActiveByPeriod(tt.args.ctx, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userRepo.CountActiveByPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("userRepo.CountActiveByPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userRepo_FindParentsByStudentID(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx       context.Context
		studentID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &userRepo{
				db: tt.fields.db,
			}
			got, err := r.FindParentsByStudentID(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userRepo.FindParentsByStudentID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userRepo.FindParentsByStudentID() = %v, want %v", got, tt.want)
			}
		})
	}
}
