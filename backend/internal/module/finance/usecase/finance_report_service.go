package usecase

import (
	"context"
	"fmt"
	"time"

	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
	"github.com/xuri/excelize/v2"
)

type FinanceReportService interface {
	GetArrears(ctx context.Context, page, limit int, academicYear int, classID, majorID, billTypeID uint, search string, start, end *time.Time) ([]domain.ArrearRecord, int, error)
	ExportTrendExcel(ctx context.Context, start, end *time.Time, interval string, academicYear int, classID, majorID uint) ([]byte, error)
}

type financeReportService struct {
	db    bun.IDB
	repo  repository.FinanceReportRepo
	audit auditusecase.AuditLogService
}

func NewFinanceReportService(db bun.IDB, repo repository.FinanceReportRepo, audit auditusecase.AuditLogService) FinanceReportService {
	return &financeReportService{db: db, repo: repo, audit: audit}
}

func (s *financeReportService) GetArrears(ctx context.Context, page, limit int, academicYear int, classID, majorID, billTypeID uint, search string, start, end *time.Time) ([]domain.ArrearRecord, int, error) {
	return s.repo.GetArrearsPaged(ctx, page, limit, academicYear, classID, majorID, billTypeID, search, start, end)
}

func (s *financeReportService) ExportTrendExcel(ctx context.Context, start, end *time.Time, interval string, academicYear int, classID, majorID uint) ([]byte, error) {
	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		newVals := map[string]interface{}{"interval": interval, "academic_year": academicYear, "class_id": classID, "major_id": majorID}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "EXPORT_TREND_EXCEL", "finance_reports", 0, nil, newVals, ipAddress, userAgent)
	}

	trend, err := s.repo.GetPaymentTrendByPeriod(ctx, start, end, interval, academicYear, classID, majorID)
	if err != nil {
		return nil, err
	}

	f := excelize.NewFile()
	sheet := "Sheet1"
	f.SetSheetName("Sheet1", sheet)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})

	f.SetCellValue(sheet, "A1", "No")
	f.SetCellValue(sheet, "B1", "Periode")
	f.SetCellValue(sheet, "C1", "Total Pendapatan")
	f.SetCellStyle(sheet, "A1", "C1", headerStyle)

	for i, t := range trend {
		row := i + 2
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), t["date"])
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), t["total"])
	}

	buf, _ := f.WriteToBuffer()
	return buf.Bytes(), nil
}
