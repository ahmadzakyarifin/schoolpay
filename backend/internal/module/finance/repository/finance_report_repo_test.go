package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

func TestNewFinanceReportRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want FinanceReportRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFinanceReportRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFinanceReportRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_financeReportRepo_GetTotalBillsByPeriod(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		start        *time.Time
		end          *time.Time
		academicYear int
		classID      uint
		majorID      uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &financeReportRepo{
				db: tt.fields.db,
			}
			got, err := r.GetTotalBillsByPeriod(tt.args.ctx, tt.args.start, tt.args.end, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("financeReportRepo.GetTotalBillsByPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("financeReportRepo.GetTotalBillsByPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_financeReportRepo_GetTotalPaymentsByPeriod(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		start        *time.Time
		end          *time.Time
		academicYear int
		classID      uint
		majorID      uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &financeReportRepo{
				db: tt.fields.db,
			}
			got, err := r.GetTotalPaymentsByPeriod(tt.args.ctx, tt.args.start, tt.args.end, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("financeReportRepo.GetTotalPaymentsByPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("financeReportRepo.GetTotalPaymentsByPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_financeReportRepo_GetPaymentTrendByPeriod(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		start        *time.Time
		end          *time.Time
		interval     string
		academicYear int
		classID      uint
		majorID      uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &financeReportRepo{
				db: tt.fields.db,
			}
			got, err := r.GetPaymentTrendByPeriod(tt.args.ctx, tt.args.start, tt.args.end, tt.args.interval, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("financeReportRepo.GetPaymentTrendByPeriod() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("financeReportRepo.GetPaymentTrendByPeriod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_financeReportRepo_GetDashboardSummary(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		start        *time.Time
		end          *time.Time
		academicYear int
		classID      uint
		majorID      uint
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
			r := &financeReportRepo{
				db: tt.fields.db,
			}
			got, err := r.GetDashboardSummary(tt.args.ctx, tt.args.start, tt.args.end, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("financeReportRepo.GetDashboardSummary() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("financeReportRepo.GetDashboardSummary() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_financeReportRepo_GetArrearsPaged(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		page         int
		limit        int
		academicYear int
		classID      uint
		majorID      uint
		billTypeID   uint
		search       string
		start        *time.Time
		end          *time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.ArrearRecord
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &financeReportRepo{
				db: tt.fields.db,
			}
			got, got1, err := r.GetArrearsPaged(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.academicYear, tt.args.classID, tt.args.majorID, tt.args.billTypeID, tt.args.search, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Fatalf("financeReportRepo.GetArrearsPaged() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("financeReportRepo.GetArrearsPaged() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("financeReportRepo.GetArrearsPaged() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_financeReportRepo_GetRecentPayments(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		start        *time.Time
		end          *time.Time
		academicYear int
		classID      uint
		majorID      uint
		billTypeID   uint
		search       string
		page         int
		limit        int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &financeReportRepo{
				db: tt.fields.db,
			}
			got, got1, err := r.GetRecentPayments(tt.args.ctx, tt.args.start, tt.args.end, tt.args.academicYear, tt.args.classID, tt.args.majorID, tt.args.billTypeID, tt.args.search, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Fatalf("financeReportRepo.GetRecentPayments() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("financeReportRepo.GetRecentPayments() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("financeReportRepo.GetRecentPayments() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_financeReportRepo_GetCriticalBills(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		status       string
		limit        int
		academicYear int
		classID      uint
		majorID      uint
		billTypeID   uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.CriticalBillRecord
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &financeReportRepo{
				db: tt.fields.db,
			}
			got, err := r.GetCriticalBills(tt.args.ctx, tt.args.status, tt.args.limit, tt.args.academicYear, tt.args.classID, tt.args.majorID, tt.args.billTypeID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("financeReportRepo.GetCriticalBills() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("financeReportRepo.GetCriticalBills() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_financeReportRepo_GetPaymentMethodsCount(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx          context.Context
		start        *time.Time
		end          *time.Time
		academicYear int
		classID      uint
		majorID      uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &financeReportRepo{
				db: tt.fields.db,
			}
			got, err := r.GetPaymentMethodsCount(tt.args.ctx, tt.args.start, tt.args.end, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("financeReportRepo.GetPaymentMethodsCount() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("financeReportRepo.GetPaymentMethodsCount() = %v, want %v", got, tt.want)
			}
		})
	}
}
