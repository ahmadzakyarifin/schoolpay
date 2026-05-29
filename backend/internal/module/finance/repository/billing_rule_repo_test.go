package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

func TestNewBillingRuleRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want BillingRuleRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewBillingRuleRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewBillingRuleRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billingRuleRepo_Create(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		br  *domain.BillingRule
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
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.br); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleRepo_FindAll(t *testing.T) {
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
		want    []domain.BillingRule
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			got, err := r.FindAll(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billingRuleRepo.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billingRuleRepo.FindAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billingRuleRepo_FindAllPaged(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx            context.Context
		page           int
		limit          int
		search         string
		status         string
		generateStatus string
		sort           string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.BillingRule
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			got, got1, err := r.FindAllPaged(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.status, tt.args.generateStatus, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billingRuleRepo.FindAllPaged() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billingRuleRepo.FindAllPaged() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("billingRuleRepo.FindAllPaged() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_billingRuleRepo_FindByID(t *testing.T) {
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
		want    *domain.BillingRule
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billingRuleRepo.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("billingRuleRepo.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billingRuleRepo_Update(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		br  *domain.BillingRule
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
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			if err := r.Update(tt.args.ctx, tt.args.br); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleRepo.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleRepo_Delete(t *testing.T) {
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
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			if err := r.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleRepo.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleRepo_Restore(t *testing.T) {
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
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			if err := r.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleRepo.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleRepo_ToggleStatus(t *testing.T) {
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
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			if err := r.ToggleStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleRepo.ToggleStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleRepo_BulkDelete(t *testing.T) {
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
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			if err := r.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleRepo.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleRepo_BulkRestore(t *testing.T) {
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
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			if err := r.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("billingRuleRepo.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_billingRuleRepo_Exists(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx        context.Context
		btID       uint
		targetType string
		targetID   uint
		classID    *uint
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
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			got, err := r.Exists(tt.args.ctx, tt.args.btID, tt.args.targetType, tt.args.targetID, tt.args.classID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billingRuleRepo.Exists() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("billingRuleRepo.Exists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billingRuleRepo_ExistsExcludeID(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx        context.Context
		btID       uint
		targetType string
		targetID   uint
		classID    *uint
		excludeID  uint
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
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			got, err := r.ExistsExcludeID(tt.args.ctx, tt.args.btID, tt.args.targetType, tt.args.targetID, tt.args.classID, tt.args.excludeID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billingRuleRepo.ExistsExcludeID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("billingRuleRepo.ExistsExcludeID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_billingRuleRepo_CountStudentBills(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx    context.Context
		ruleID uint
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
			r := &billingRuleRepo{
				db: tt.fields.db,
			}
			got, err := r.CountStudentBills(tt.args.ctx, tt.args.ruleID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("billingRuleRepo.CountStudentBills() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("billingRuleRepo.CountStudentBills() = %v, want %v", got, tt.want)
			}
		})
	}
}
