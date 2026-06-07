package usecase

import (
	"bytes"
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/jung-kurt/gofpdf"
	"github.com/xuri/excelize/v2"
)

type DashboardService interface {
	GetStats(ctx context.Context, period string, pStart, pEnd, pPrevStart, pPrevEnd *time.Time, filterEntryYear int, classID, majorID uint, page, limit int, billTypeID uint) (map[string]interface{}, error)
	GetCommunicationDetails(ctx context.Context, status, channel string, start, end *time.Time, filterEntryYear int, classID, majorID uint) ([]map[string]interface{}, error)
	ExportGlobalReport(ctx context.Context, format, tab, search string, pStart, pEnd *time.Time, filterEntryYear int, classID, majorID, billTypeID uint) ([]byte, error)
}

type dashboardService struct {
	dashRepo financerepo.DashboardRepo
	stuRepo  academicrepo.StudentRepo
	ayRepo   academicrepo.AcademicYearRepo
	audit    auditusecase.AuditLogService
}

func NewDashboardService(
	dashRepo financerepo.DashboardRepo,
	stuRepo academicrepo.StudentRepo,
	ayRepo academicrepo.AcademicYearRepo,
	audit auditusecase.AuditLogService,
) DashboardService {
	return &dashboardService{
		dashRepo: dashRepo,
		stuRepo:  stuRepo,
		ayRepo:   ayRepo,
		audit:    audit,
	}
}

func (s *dashboardService) GetStats(ctx context.Context, period string, pStart, pEnd, pPrevStart, pPrevEnd *time.Time, filterEntryYear int, classID, majorID uint, page, limit int, billTypeID uint) (map[string]interface{}, error) {
	// 1. Users Stats (0 if academic filters are active)
	var totalUsers, prevUsers, parentCount, adminCount int
	if filterEntryYear > 0 || classID > 0 || majorID > 0 {
		totalUsers = 0
		prevUsers = 0
		parentCount = 0
		adminCount = 0
	} else {
		_, _ = s.dashRepo.CountNewUsersByPeriod(ctx, pStart, pEnd)
		totalUsers, _ = s.dashRepo.CountTotalUsersUpTo(ctx, pEnd)
		prevUsers, _ = s.dashRepo.CountTotalUsersUpTo(ctx, pPrevEnd)
		parentCount, _ = s.dashRepo.CountUsersByRoleUpTo(ctx, pEnd, "parent")
		adminCount, _ = s.dashRepo.CountUsersByRoleUpTo(ctx, pEnd, "admin")
	}

	// 2. Students Stats (Accumulative running totals)
	totalStudents, _ := s.dashRepo.CountStudentsByStatusUpTo(ctx, pEnd, "active", filterEntryYear, classID, majorID)
	prevStudents, _ := s.dashRepo.CountStudentsByStatusUpTo(ctx, pPrevEnd, "active", filterEntryYear, classID, majorID)
	totalAllStudents, _ := s.dashRepo.CountStudentsByStatusUpTo(ctx, pEnd, "all", filterEntryYear, classID, majorID)
	activeStudentsCount, _ := s.dashRepo.CountStudentsByStatusUpTo(ctx, pEnd, "active", filterEntryYear, classID, majorID)
	inactiveStudentsCount, _ := s.dashRepo.CountStudentsByStatusUpTo(ctx, pEnd, "inactive", filterEntryYear, classID, majorID)

	currBills, _ := s.dashRepo.GetTotalBillsByPeriod(ctx, pStart, pEnd, filterEntryYear, classID, majorID)
	currPayments, _ := s.dashRepo.GetTotalPaymentsByPeriod(ctx, pStart, pEnd, filterEntryYear, classID, majorID)

	prevBills, _ := s.dashRepo.GetTotalBillsByPeriod(ctx, pPrevStart, pPrevEnd, filterEntryYear, classID, majorID)
	prevPayments, _ := s.dashRepo.GetTotalPaymentsByPeriod(ctx, pPrevStart, pPrevEnd, filterEntryYear, classID, majorID)

	// Gender stats should be accumulative (nil start date, pEnd end date)
	genderStats, _ := s.dashRepo.GetGenderStats(ctx, nil, pEnd, filterEntryYear, classID, majorID)
	majorStats, _ := s.stuRepo.CountByMajor(ctx, filterEntryYear, classID, majorID)
	classStats, _ := s.stuRepo.CountByClass(ctx, filterEntryYear, classID, majorID)

	// Dynamic trend interval based on period
	trendInterval := "daily"
	if period == "yearly" {
		trendInterval = "monthly"
	}
	paymentTrend, _ := s.dashRepo.GetPaymentTrendByPeriod(ctx, pStart, pEnd, trendInterval, filterEntryYear, classID, majorID)
	summary, _ := s.dashRepo.GetDashboardSummary(ctx, pStart, pEnd, filterEntryYear, classID, majorID)
	paymentMethods, _ := s.dashRepo.GetPaymentMethodsCount(ctx, pStart, pEnd, filterEntryYear, classID, majorID)

	// Efficacy stats are filtered by the selected period
	waEfficacy, _ := s.dashRepo.GetEfficacyStats(ctx, "whatsapp", pStart, pEnd, filterEntryYear, classID, majorID)
	emailEfficacy, _ := s.dashRepo.GetEfficacyStats(ctx, "email", pStart, pEnd, filterEntryYear, classID, majorID)

	// Critical bills are accumulative (pass nil start/end dates)
	overdueBills, _ := s.dashRepo.GetCriticalBills(ctx, "overdue", nil, nil, filterEntryYear, classID, majorID, 0)
	dueSoonBills, _ := s.dashRepo.GetCriticalBills(ctx, "due_soon", nil, nil, filterEntryYear, classID, majorID, 0)

	recentPayments, totalPaymentsCount, _ := s.dashRepo.GetRecentPayments(ctx, pStart, pEnd, filterEntryYear, classID, majorID, billTypeID, "", page, limit)
	recentNotifications, _ := s.dashRepo.GetRecentNotifications(ctx, pStart, pEnd, filterEntryYear, classID, majorID)

	activeAY, _, _ := s.ayRepo.FindAll(ctx, 1, 100, "", "", "")
	var availableYears []int
	for _, ay := range activeAY {
		availableYears = append(availableYears, ay.Year)
	}

	calculateGrowth := func(current, previous float64) float64 {
		if previous == 0 {
			if current > 0 {
				return 100
			}
			return 0
		}
		return ((current - previous) / previous) * 100
	}

	return map[string]interface{}{
		"recent_payments":      recentPayments,
		"recent_notifications": recentNotifications,
		"total_payments_count": totalPaymentsCount,
		"stats": map[string]interface{}{
			"users": map[string]interface{}{
				"total":           totalUsers,
				"new_this_period": totalUsers, // Frontend uses new_this_period as the main total card value
				"growth":          calculateGrowth(float64(totalUsers), float64(prevUsers)),
				"parent_count":    parentCount,
				"admin_count":     adminCount,
			},
			"students": map[string]interface{}{
				"total_all":      totalAllStudents,
				"total":          totalStudents,
				"growth":         calculateGrowth(float64(totalStudents), float64(prevStudents)),
				"active_count":   activeStudentsCount,
				"inactive_count": inactiveStudentsCount,
			},
			"bills": map[string]interface{}{
				"total":  currBills,
				"growth": calculateGrowth(currBills, prevBills),
			},
			"payments": map[string]interface{}{
				"total":  currPayments,
				"growth": calculateGrowth(currPayments, prevPayments),
			},
			"paid_amount":   summary["total_paid_amount"],
			"unpaid_amount": summary["total_unpaid_amount"],
			"paid_count":    summary["paid_count"],
			"unpaid_count":  summary["unpaid_count"],
		},
		"demographics": map[string]interface{}{
			"gender":          genderStats,
			"major":           majorStats,
			"class":           classStats,
			"payment_methods": paymentMethods,
			"whatsapp":        waEfficacy,
			"email":           emailEfficacy,
		},
		"payment_trend": paymentTrend,
		"entry_years":   availableYears,
		"critical_bills": map[string]interface{}{
			"overdue":  overdueBills,
			"due_soon": dueSoonBills,
		},
	}, nil
}

func (s *dashboardService) GetCommunicationDetails(ctx context.Context, status, channel string, start, end *time.Time, filterEntryYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	return s.dashRepo.GetCommunicationDetails(ctx, status, channel, start, end, filterEntryYear, classID, majorID)
}

func (s *dashboardService) ExportGlobalReport(ctx context.Context, format, tab, search string, pStart, pEnd *time.Time, filterEntryYear int, classID, majorID, billTypeID uint) ([]byte, error) {
	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		newVals := map[string]interface{}{
			"tab": tab, "search": search, "format": format, "entry_year": filterEntryYear,
			"class_id": classID, "major_id": majorID, "bill_type_id": billTypeID,
			"start_date": pStart, "end_date": pEnd,
		}
		_ = s.audit.Log(ctx, nil, userID, userName, role, "EXPORT_GLOBAL_FINANCE_REPORT", "finance_reports", 0, nil, newVals, ipAddress, userAgent)
	}

	if format == "pdf" {
		pdf := gofpdf.New("P", "mm", "A4", "")
		pdf.AddPage()
		pdf.SetFont("Arial", "B", 16)
		if tab == "arrears" {
			pdf.CellFormat(190, 10, "LAPORAN DATA TUNGGAKAN SISWA", "0", 1, "C", false, 0, "")
			pdf.Ln(5)
			pdf.SetFont("Arial", "B", 10)
			headers := []string{"No", "Nama Siswa", "Kelas", "Tagihan", "Sisa Tagihan", "Jatuh Tempo", "Status"}
			widths := []float64{8, 38, 23, 34, 28, 30, 29}
			for i, h := range headers {
				pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
			}
			pdf.Ln(-1)
			pdf.SetFont("Arial", "", 9)
			list, _, _ := s.dashRepo.GetArrearsPaged(ctx, 1, 1000000, filterEntryYear, classID, majorID, billTypeID, search, pStart, pEnd)
			for i, item := range list {
				pdf.CellFormat(widths[0], 7, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
				pdf.CellFormat(widths[1], 7, item.StudentName, "1", 0, "L", false, 0, "")
				pdf.CellFormat(widths[2], 7, item.ClassName, "1", 0, "L", false, 0, "")
				pdf.CellFormat(widths[3], 7, item.BillName, "1", 0, "L", false, 0, "")
				pdf.CellFormat(widths[4], 7, fmt.Sprintf("Rp %.0f", item.Amount-item.TotalPaid), "1", 0, "R", false, 0, "")
				pdf.CellFormat(widths[5], 7, item.DueDate.Format("02/01/2006"), "1", 0, "C", false, 0, "")
				pdf.CellFormat(widths[6], 7, dueStatusText(item.DueDate), "1", 0, "C", false, 0, "")
				pdf.Ln(-1)
			}
		} else {
			pdf.CellFormat(190, 10, "LAPORAN TRANSAKSI PEMBAYARAN", "0", 1, "C", false, 0, "")
			pdf.Ln(5)
			pdf.SetFont("Arial", "B", 10)
			headers := []string{"No", "Nama Siswa", "Jenis", "Total", "Saldo", "Tunai/Gateway", "Metode", "Tanggal"}
			widths := []float64{8, 32, 33, 23, 22, 28, 22, 24}
			for i, h := range headers {
				pdf.CellFormat(widths[i], 8, h, "1", 0, "C", false, 0, "")
			}
			pdf.Ln(-1)
			pdf.SetFont("Arial", "", 9)
			list, _, _ := s.dashRepo.GetRecentPayments(ctx, pStart, pEnd, filterEntryYear, classID, majorID, billTypeID, search, 1, 1000000)
			for i, item := range list {
				amt := mapNumber(item, "amount")
				depositApplied := mapNumber(item, "deposit_applied")
				cashOrGateway := mapNumber(item, "cash_or_gateway_amount")
				billTypes := mapString(item, "bill_type_names")
				if billTypes == "" {
					billTypes = "Deposit/Kustom"
				}
				billTypes = strings.ReplaceAll(billTypes, "||", ", ")
				method := ""
				if v, ok := item["method"].(string); ok {
					method = strings.ToUpper(v)
				}
				studentName := ""
				if v, ok := item["student_name"].(string); ok {
					studentName = v
				}
				paidAtStr := "-"
				if v, ok := item["created_at"].(time.Time); ok {
					paidAtStr = v.Format("02/01/2006 15:04:05")
				} else if v, ok := item["created_at"].(string); ok {
					paidAtStr = v
				}

				pdf.CellFormat(widths[0], 7, fmt.Sprintf("%d", i+1), "1", 0, "C", false, 0, "")
				pdf.CellFormat(widths[1], 7, studentName, "1", 0, "L", false, 0, "")
				pdf.CellFormat(widths[2], 7, billTypes, "1", 0, "L", false, 0, "")
				pdf.CellFormat(widths[3], 7, fmt.Sprintf("Rp %.0f", amt), "1", 0, "R", false, 0, "")
				pdf.CellFormat(widths[4], 7, fmt.Sprintf("Rp %.0f", depositApplied), "1", 0, "R", false, 0, "")
				pdf.CellFormat(widths[5], 7, fmt.Sprintf("Rp %.0f", cashOrGateway), "1", 0, "R", false, 0, "")
				pdf.CellFormat(widths[6], 7, method, "1", 0, "C", false, 0, "")
				pdf.CellFormat(widths[7], 7, paidAtStr, "1", 0, "C", false, 0, "")
				pdf.Ln(-1)
			}
		}

		var buf bytes.Buffer
		err := pdf.Output(&buf)
		if err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}

	f := excelize.NewFile()
	sheet := "Laporan"
	f.SetSheetName("Sheet1", sheet)

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	bodyStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	if tab == "arrears" {
		list, _, _ := s.dashRepo.GetArrearsPaged(ctx, 1, 1000000, filterEntryYear, classID, majorID, billTypeID, search, pStart, pEnd)
		headers := []string{"NO", "NAMA SISWA", "KELAS", "TAGIHAN", "PERIODE", "SISA TAGIHAN", "JATUH TEMPO", "STATUS JATUH TEMPO"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, h)
			f.SetCellStyle(sheet, cell, cell, headerStyle)
			f.SetColWidth(sheet, cell[:1], cell[:1], 22)
		}
		for i, item := range list {
			row := i + 2
			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), item.StudentName)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), item.ClassName)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), item.BillName)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), item.Period)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), item.Amount-item.TotalPaid)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), item.DueDate.Format("02/01/2006"))
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), dueStatusText(item.DueDate))
			f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("H%d", row), bodyStyle)
		}
	} else {
		list, _, _ := s.dashRepo.GetRecentPayments(ctx, pStart, pEnd, filterEntryYear, classID, majorID, billTypeID, search, 1, 1000000)
		headers := []string{"NO", "NAMA SISWA", "JENIS TAGIHAN", "TOTAL ALOKASI", "SALDO DIPAKAI", "TUNAI/GATEWAY", "METODE", "TANGGAL PEMBAYARAN"}
		for i, h := range headers {
			cell, _ := excelize.CoordinatesToCellName(i+1, 1)
			f.SetCellValue(sheet, cell, h)
			f.SetCellStyle(sheet, cell, cell, headerStyle)
			f.SetColWidth(sheet, cell[:1], cell[:1], 25)
		}
		for i, item := range list {
			row := i + 2
			amt := mapNumber(item, "amount")
			depositApplied := mapNumber(item, "deposit_applied")
			cashOrGateway := mapNumber(item, "cash_or_gateway_amount")
			billTypes := mapString(item, "bill_type_names")
			if billTypes == "" {
				billTypes = "Deposit/Kustom"
			}
			billTypes = strings.ReplaceAll(billTypes, "||", ", ")
			method := ""
			if v, ok := item["method"].(string); ok {
				method = strings.ToUpper(v)
			}
			studentName := ""
			if v, ok := item["student_name"].(string); ok {
				studentName = v
			}
			paidAtStr := "-"
			if v, ok := item["created_at"].(time.Time); ok {
				paidAtStr = v.Format("02/01/2006 15:04:05")
			} else if v, ok := item["created_at"].(string); ok {
				paidAtStr = v
			}

			f.SetCellValue(sheet, fmt.Sprintf("A%d", row), i+1)
			f.SetCellValue(sheet, fmt.Sprintf("B%d", row), studentName)
			f.SetCellValue(sheet, fmt.Sprintf("C%d", row), billTypes)
			f.SetCellValue(sheet, fmt.Sprintf("D%d", row), amt)
			f.SetCellValue(sheet, fmt.Sprintf("E%d", row), depositApplied)
			f.SetCellValue(sheet, fmt.Sprintf("F%d", row), cashOrGateway)
			f.SetCellValue(sheet, fmt.Sprintf("G%d", row), method)
			f.SetCellValue(sheet, fmt.Sprintf("H%d", row), paidAtStr)
			f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("H%d", row), bodyStyle)
		}
	}

	buf, err := f.WriteToBuffer()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func dueStatusText(dueDate time.Time) string {
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	due := time.Date(dueDate.Year(), dueDate.Month(), dueDate.Day(), 0, 0, 0, 0, time.Local)
	days := int(due.Sub(today).Hours() / 24)
	if days == 0 {
		return "Hari Ini"
	}
	if days > 0 {
		return fmt.Sprintf("H-%d", days)
	}
	return fmt.Sprintf("Telat %d Hari", -days)
}

func mapNumber(item map[string]interface{}, key string) float64 {
	if item[key] == nil {
		return 0
	}
	switch v := item[key].(type) {
	case []byte:
		n, _ := strconv.ParseFloat(string(v), 64)
		return n
	case string:
		n, _ := strconv.ParseFloat(v, 64)
		return n
	case float64:
		return v
	case float32:
		return float64(v)
	case int64:
		return float64(v)
	case int32:
		return float64(v)
	case int:
		return float64(v)
	default:
		return 0
	}
}

func mapString(item map[string]interface{}, key string) string {
	if item[key] == nil {
		return ""
	}
	switch v := item[key].(type) {
	case []byte:
		return string(v)
	case string:
		return v
	case fmt.Stringer:
		return v.String()
	default:
		return fmt.Sprint(v)
	}
}
