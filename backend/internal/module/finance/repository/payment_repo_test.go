package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

func TestNewPaymentRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want PaymentRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPaymentRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPaymentRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paymentRepo_Create(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
		p   *domain.Payment
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
			r := &paymentRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.db, tt.args.p); (err != nil) != tt.wantErr {
				t.Errorf("paymentRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_paymentRepo_CreateDetail(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
		pd  *domain.PaymentDetail
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
			r := &paymentRepo{
				db: tt.fields.db,
			}
			if err := r.CreateDetail(tt.args.ctx, tt.args.db, tt.args.pd); (err != nil) != tt.wantErr {
				t.Errorf("paymentRepo.CreateDetail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_paymentRepo_FindByRef(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		ref string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Payment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &paymentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByRef(tt.args.ctx, tt.args.ref)
			if (err != nil) != tt.wantErr {
				t.Fatalf("paymentRepo.FindByRef() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paymentRepo.FindByRef() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paymentRepo_FindByID(t *testing.T) {
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
		want    *domain.Payment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &paymentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("paymentRepo.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paymentRepo.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paymentRepo_FindDetailsByPaymentID(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx       context.Context
		paymentID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.PaymentDetail
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &paymentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindDetailsByPaymentID(tt.args.ctx, tt.args.paymentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("paymentRepo.FindDetailsByPaymentID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paymentRepo.FindDetailsByPaymentID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paymentRepo_FindByStudent(t *testing.T) {
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
		want    []domain.Payment
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &paymentRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByStudent(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("paymentRepo.FindByStudent() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("paymentRepo.FindByStudent() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_paymentRepo_DeleteDetailsByPaymentID(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx       context.Context
		db        bun.IDB
		paymentID uint
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
			r := &paymentRepo{
				db: tt.fields.db,
			}
			if err := r.DeleteDetailsByPaymentID(tt.args.ctx, tt.args.db, tt.args.paymentID); (err != nil) != tt.wantErr {
				t.Errorf("paymentRepo.DeleteDetailsByPaymentID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
