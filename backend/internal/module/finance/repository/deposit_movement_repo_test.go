package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

func TestNewDepositMovementRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want DepositMovementRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewDepositMovementRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDepositMovementRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_depositMovementRepo_Create(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
		dm  *domain.DepositMovement
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
			r := &depositMovementRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.db, tt.args.dm); (err != nil) != tt.wantErr {
				t.Errorf("depositMovementRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_depositMovementRepo_FindByStudent(t *testing.T) {
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
		want    []domain.DepositMovement
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &depositMovementRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByStudent(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("depositMovementRepo.FindByStudent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("depositMovementRepo.FindByStudent() = %v, want %v", got, tt.want)
			}
		})
	}
}
