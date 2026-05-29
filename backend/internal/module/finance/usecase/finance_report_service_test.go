package usecase

import (
	"context"
	"reflect"
	"testing"
	"time"

	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	"github.com/uptrace/bun"
)

func TestNewFinanceReportService(t *testing.T) {
	type args struct {
		db    bun.IDB
		repo  repository.FinanceReportRepo
		audit auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want FinanceReportService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewFinanceReportService(tt.args.db, tt.args.repo, tt.args.audit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFinanceReportService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_financeReportService_GetArrears(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.FinanceReportRepo
		audit auditusecase.AuditLogService
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
			s := &financeReportService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, got1, err := s.GetArrears(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.academicYear, tt.args.classID, tt.args.majorID, tt.args.billTypeID, tt.args.search, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Fatalf("financeReportService.GetArrears() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("financeReportService.GetArrears() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("financeReportService.GetArrears() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_financeReportService_ExportTrendExcel(t *testing.T) {
	type fields struct {
		db    bun.IDB
		repo  repository.FinanceReportRepo
		audit auditusecase.AuditLogService
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
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &financeReportService{
				db:    tt.fields.db,
				repo:  tt.fields.repo,
				audit: tt.fields.audit,
			}
			got, err := s.ExportTrendExcel(tt.args.ctx, tt.args.start, tt.args.end, tt.args.interval, tt.args.academicYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("financeReportService.ExportTrendExcel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("financeReportService.ExportTrendExcel() = %v, want %v", got, tt.want)
			}
		})
	}
}
