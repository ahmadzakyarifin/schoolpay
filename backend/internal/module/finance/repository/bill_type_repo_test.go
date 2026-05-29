package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

func TestNewBillTypeRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want BillTypeRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBillTypeRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBillTypeRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billTypeRepo_Create(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("billTypeRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeRepo_FindAll(t *testing.T) {
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
		want    []domain.BillType
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			got, err := r.FindAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billTypeRepo.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billTypeRepo.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billTypeRepo_FindAllPaged(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			got, got1, err := r.FindAllPaged(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.filterType, tt.args.status, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billTypeRepo.FindAllPaged() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billTypeRepo.FindAllPaged() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("billTypeRepo.FindAllPaged() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_billTypeRepo_Update(t *testing.T) {
	type fields struct {
		db *bun.DB
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
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.ctx, tt.args.b); (err != nil) != tt.wantErr {
				t.Errorf("billTypeRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeRepo_Delete(t *testing.T) {
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
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("billTypeRepo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeRepo_Restore(t *testing.T) {
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
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("billTypeRepo.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeRepo_ToggleStatus(t *testing.T) {
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
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			if err := r.ToggleStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("billTypeRepo.ToggleStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeRepo_BulkDelete(t *testing.T) {
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
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			if err := r.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("billTypeRepo.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeRepo_BulkRestore(t *testing.T) {
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
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			if err := r.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("billTypeRepo.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billTypeRepo_Exists(t *testing.T) {
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
		want    bool
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			got, err := r.Exists(tt.args.ctx, tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billTypeRepo.Exists() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("billTypeRepo.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billTypeRepo_ExistsExcludeID(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx  context.Context
		name string
		id   uint
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
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			got, err := r.ExistsExcludeID(tt.args.ctx, tt.args.name, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billTypeRepo.ExistsExcludeID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("billTypeRepo.ExistsExcludeID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billTypeRepo_FindByID(t *testing.T) {
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
		want    *domain.BillType
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billTypeRepo.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billTypeRepo.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billTypeRepo_CountBillingRules(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx        context.Context
		billTypeID uint
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
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			got, err := r.CountBillingRules(tt.args.ctx, tt.args.billTypeID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billTypeRepo.CountBillingRules() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("billTypeRepo.CountBillingRules() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billTypeRepo_CountStudentBills(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx        context.Context
		billTypeID uint
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
			r := &billTypeRepo{
				db: tt.fields.db,
			}
			got, err := r.CountStudentBills(tt.args.ctx, tt.args.billTypeID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billTypeRepo.CountStudentBills() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("billTypeRepo.CountStudentBills() = %v, want %v", got, tt.want)
			}
		})
	}
}
