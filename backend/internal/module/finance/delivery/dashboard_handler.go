package delivery

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	academicdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financedomain "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	userauthdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/jung-kurt/gofpdf"
	"github.com/uptrace/bun"
	"github.com/xuri/excelize/v2"
)

type DashboardHandler struct {
	db        *bun.DB
	userRepo  userauthrepo.UserRepo
	stuRepo   academicrepo.StudentRepo
	finRepo   financerepo.FinanceReportRepo
	ayRepo    academicrepo.AcademicYearRepo
	notifRepo notificationrepo.NotificationRepo
	audit     auditusecase.AuditLogService
}

func NewDashboardHandler(
	db *bun.DB,
	userRepo userauthrepo.UserRepo,
	stuRepo academicrepo.StudentRepo,
	finRepo financerepo.FinanceReportRepo,
	ayRepo academicrepo.AcademicYearRepo,
	notifRepo notificationrepo.NotificationRepo,
	audit auditusecase.AuditLogService,
) *DashboardHandler {
	return &DashboardHandler{db, userRepo, stuRepo, finRepo, ayRepo, notifRepo, audit}
}

func (h *DashboardHandler) countUsersForDashboard(ctx context.Context, start, end *time.Time, search string) (int, error) {
	q := h.db.NewSelect().Model((*userauthdomain.User)(nil))
	if start != nil {
		q.Where("u.created_at >= ?", start)
	}
	if end != nil {
		q.Where("u.created_at <= ?", end)
	}
	if search != "" {
		like := "%" + search + "%"
		q.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("u.name LIKE ?", like).
				WhereOr("u.email LIKE ?", like).
				WhereOr("u.phone_number LIKE ?", like).
				WhereOr("EXISTS (SELECT 1 FROM students st WHERE st.parent_id = u.id AND st.name LIKE ? AND st.deleted_at IS NULL)", like)
		})
	}
	return q.Count(ctx)
}

func (h *DashboardHandler) countStudentsForDashboard(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint, search string) (int, error) {
	q := h.db.NewSelect().
		Model((*academicdomain.Student)(nil)).
		Join("LEFT JOIN classes c ON s.class_id = c.id").
		Where("s.status = 'active'")
	if start != nil {
		q.Where("s.created_at >= ?", start)
	}
	if end != nil {
		q.Where("s.created_at <= ?", end)
	}
	applyDashboardStudentFilters(q, academicYear, classID, majorID, search)
	return q.Count(ctx)
}

func (h *DashboardHandler) sumBillsForDashboard(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint, search string) (float64, error) {
	var total float64
	q := h.db.NewSelect().
		Model((*financedomain.StudentBill)(nil)).
		ColumnExpr("COALESCE(SUM(sb.amount - sb.total_paid), 0)").
		Join("JOIN students s ON sb.student_id = s.id").
		Join("LEFT JOIN classes c ON s.class_id = c.id").
		Join("LEFT JOIN bill_types bt ON sb.bill_type_id = bt.id").
		Where("sb.status != 'paid' AND sb.status != 'voided'")
	if start != nil {
		q.Where("sb.created_at >= ?", start)
	}
	if end != nil {
		q.Where("sb.created_at <= ?", end)
	}
	applyDashboardStudentFilters(q, academicYear, classID, majorID, "")
	applyDashboardBillSearch(q, search)
	return total, q.Scan(ctx, &total)
}

func (h *DashboardHandler) sumPaymentsForDashboard(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint, search string) (float64, error) {
	var total float64
	q := h.db.NewSelect().
		Model((*financedomain.Payment)(nil)).
		ColumnExpr("COALESCE(SUM(p.amount), 0)").
		Join("JOIN students s ON p.student_id = s.id").
		Join("LEFT JOIN classes c ON s.class_id = c.id").
		Where("p.status = 'success'")
	if start != nil {
		q.Where("p.paid_at >= ?", start)
	}
	if end != nil {
		q.Where("p.paid_at <= ?", end)
	}
	applyDashboardStudentFilters(q, academicYear, classID, majorID, "")
	applyDashboardPaymentSearch(q, search)
	return total, q.Scan(ctx, &total)
}

func (h *DashboardHandler) dashboardSummaryForSearch(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint, search string) map[string]interface{} {
	res := map[string]interface{}{
		"total_paid_amount":   0.0,
		"total_unpaid_amount": 0.0,
		"paid_count":          0,
		"unpaid_count":        0,
	}

	var unpaid struct {
		Amount float64 `bun:"amount"`
		Count  int     `bun:"count"`
	}
	qUnpaid := h.db.NewSelect().
		Model((*financedomain.StudentBill)(nil)).
		ColumnExpr("COALESCE(SUM(sb.amount - sb.total_paid), 0) as amount").
		ColumnExpr("COUNT(sb.id) as count").
		Join("JOIN students s ON sb.student_id = s.id").
		Join("LEFT JOIN classes c ON s.class_id = c.id").
		Join("LEFT JOIN bill_types bt ON sb.bill_type_id = bt.id").
		Where("sb.status != 'paid' AND sb.status != 'voided'")
	if end != nil {
		qUnpaid.Where("sb.created_at <= ?", end)
	}
	applyDashboardStudentFilters(qUnpaid, academicYear, classID, majorID, "")
	applyDashboardBillSearch(qUnpaid, search)
	_ = qUnpaid.Scan(ctx, &unpaid)

	var paid struct {
		Amount float64 `bun:"amount"`
		Count  int     `bun:"count"`
	}
	qPaid := h.db.NewSelect().
		Model((*financedomain.Payment)(nil)).
		ColumnExpr("COALESCE(SUM(p.amount), 0) as amount").
		ColumnExpr("COUNT(p.id) as count").
		Join("JOIN students s ON p.student_id = s.id").
		Join("LEFT JOIN classes c ON s.class_id = c.id").
		Where("p.status = 'success'")
	if start != nil {
		qPaid.Where("p.paid_at >= ?", start)
	}
	if end != nil {
		qPaid.Where("p.paid_at <= ?", end)
	}
	applyDashboardStudentFilters(qPaid, academicYear, classID, majorID, "")
	applyDashboardPaymentSearch(qPaid, search)
	_ = qPaid.Scan(ctx, &paid)

	res["total_paid_amount"] = paid.Amount
	res["total_unpaid_amount"] = unpaid.Amount
	res["paid_count"] = paid.Count
	res["unpaid_count"] = unpaid.Count
	return res
}

func applyDashboardStudentFilters(q *bun.SelectQuery, academicYear int, classID, majorID uint, search string) {
	if academicYear > 0 {
		if academicYear < 1000 {
			q.Where("s.entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
		} else {
			q.Where("s.entry_year = ?", academicYear)
		}
	}
	if classID > 0 {
		q.Where("s.class_id = ?", classID)
	}
	if majorID > 0 {
		q.Where("s.major_id = ? OR c.major_id = ?", majorID, majorID)
	}
	if search != "" {
		like := "%" + search + "%"
		q.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("s.name LIKE ?", like).
				WhereOr("s.nis LIKE ?", like).
				WhereOr("s.nisn LIKE ?", like).
				WhereOr("s.email LIKE ?", like).
				WhereOr("s.phone_number LIKE ?", like)
		})
	}
}

func applyDashboardBillSearch(q *bun.SelectQuery, search string) {
	if search != "" {
		like := "%" + search + "%"
		q.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("s.name LIKE ?", like).
				WhereOr("s.nis LIKE ?", like).
				WhereOr("s.nisn LIKE ?", like).
				WhereOr("s.email LIKE ?", like).
				WhereOr("s.phone_number LIKE ?", like).
				WhereOr("bt.name LIKE ?", like).
				WhereOr("sb.name LIKE ?", like)
		})
	}
}

func applyDashboardPaymentSearch(q *bun.SelectQuery, search string) {
	if search != "" {
		like := "%" + search + "%"
		q.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("s.name LIKE ?", like).
				WhereOr("s.nis LIKE ?", like).
				WhereOr("s.nisn LIKE ?", like).
				WhereOr("s.email LIKE ?", like).
				WhereOr("s.phone_number LIKE ?", like).
				WhereOr("p.transaction_ref LIKE ?", like)
		})
	}
}

func (h *DashboardHandler) GetStats(c *gin.Context) {
	ctx := c.Request.Context()
	period := c.DefaultQuery("period", "all")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")

	now := time.Now()
	var start, end, prevStart, prevEnd time.Time

	// 1. Resolve Filter Params
	var filterEntryYear int
	if c.Query("academic_year_id") != "" {
		filterEntryYear, _ = strconv.Atoi(c.Query("academic_year_id"))
	} else if c.Query("entry_year") != "" {
		filterEntryYear, _ = strconv.Atoi(c.Query("entry_year"))
	}
	classID, _ := strconv.ParseUint(c.Query("class_id"), 10, 32)
	majorID, _ := strconv.ParseUint(c.Query("major_id"), 10, 32)

	// 2. Handle Periods
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")


	if period == "custom" && startDateStr != "" && endDateStr != "" {
		s, _ := time.Parse("2006-01-02", startDateStr)
		e, _ := time.Parse("2006-01-02", endDateStr)
		start = time.Date(s.Year(), s.Month(), s.Day(), 0, 0, 0, 0, time.Local)
		end = time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, time.Local)
		diff := end.Sub(start)
		prevEnd = start
		prevStart = start.Add(-diff)
	}

	refTime := now
	refDateStr := c.Query("ref_date")
	if refDateStr != "" {
		switch period {
		case "monthly":
			if t, err := time.Parse("2006-01", refDateStr); err == nil {
				refTime = t
			}
		case "daily":
			if t, err := time.Parse("2006-01-02", refDateStr); err == nil {
				refTime = t
			}
		case "yearly":
			if t, err := time.Parse("2006", refDateStr); err == nil {
				refTime = t
			}
		}
	}

	if period != "custom" && period != "all" {
		switch period {
		case "daily":
			start = time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 0, 0, 0, 0, refTime.Location())
			end = time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, start.Location())
			prevStart = start.Add(-24 * time.Hour)
			prevEnd = time.Date(prevStart.Year(), prevStart.Month(), prevStart.Day(), 23, 59, 59, 0, prevStart.Location())
		case "yearly":
			start = time.Date(refTime.Year(), 1, 1, 0, 0, 0, 0, refTime.Location())
			end = time.Date(start.Year(), 12, 31, 23, 59, 59, 0, start.Location())
			prevStart = start.AddDate(-1, 0, 0)
			prevEnd = time.Date(prevStart.Year(), 12, 31, 23, 59, 59, 0, prevStart.Location())
		default: // monthly
			start = time.Date(refTime.Year(), refTime.Month(), 1, 0, 0, 0, 0, refTime.Location())
			end = start.AddDate(0, 1, 0).Add(-time.Second)
			prevStart = start.AddDate(0, -1, 0)
			prevEnd = prevStart.AddDate(0, 1, 0).Add(-time.Second)
		}
	}

	var pStart, pEnd, pPrevStart, pPrevEnd *time.Time
	if period != "all" && !(period == "custom" && (startDateStr == "" || endDateStr == "")) {
		pStart = &start
		pEnd = &end
		pPrevStart = &prevStart
		pPrevEnd = &prevEnd
	}

	// 3. Fetch Stats
	currUsers, _ := h.userRepo.CountActiveByPeriod(ctx, pStart, pEnd)
	totalUsers, _ := h.userRepo.CountActiveByPeriod(ctx, nil, nil)
	currStudents, _ := h.stuRepo.CountActiveByPeriod(ctx, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID))
	totalStudents, _ := h.stuRepo.CountActiveByPeriod(ctx, nil, nil, filterEntryYear, uint(classID), uint(majorID))
	currBills, _ := h.finRepo.GetTotalBillsByPeriod(ctx, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID))
	currPayments, _ := h.finRepo.GetTotalPaymentsByPeriod(ctx, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID))

	prevUsers, _ := h.userRepo.CountActiveByPeriod(ctx, pPrevStart, pPrevEnd)
	prevStudents, _ := h.stuRepo.CountActiveByPeriod(ctx, pPrevStart, pPrevEnd, filterEntryYear, uint(classID), uint(majorID))
	prevBills, _ := h.finRepo.GetTotalBillsByPeriod(ctx, pPrevStart, pPrevEnd, filterEntryYear, uint(classID), uint(majorID))
	prevPayments, _ := h.finRepo.GetTotalPaymentsByPeriod(ctx, pPrevStart, pPrevEnd, filterEntryYear, uint(classID), uint(majorID))

	search = strings.TrimSpace(search)
	if search != "" {
		currUsers, _ = h.countUsersForDashboard(ctx, pStart, pEnd, search)
		totalUsers, _ = h.countUsersForDashboard(ctx, nil, nil, search)
		currStudents, _ = h.countStudentsForDashboard(ctx, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID), search)
		totalStudents, _ = h.countStudentsForDashboard(ctx, nil, nil, filterEntryYear, uint(classID), uint(majorID), search)
		currBills, _ = h.sumBillsForDashboard(ctx, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID), search)
		currPayments, _ = h.sumPaymentsForDashboard(ctx, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID), search)
		prevUsers, _ = h.countUsersForDashboard(ctx, pPrevStart, pPrevEnd, search)
		prevStudents, _ = h.countStudentsForDashboard(ctx, pPrevStart, pPrevEnd, filterEntryYear, uint(classID), uint(majorID), search)
		prevBills, _ = h.sumBillsForDashboard(ctx, pPrevStart, pPrevEnd, filterEntryYear, uint(classID), uint(majorID), search)
		prevPayments, _ = h.sumPaymentsForDashboard(ctx, pPrevStart, pPrevEnd, filterEntryYear, uint(classID), uint(majorID), search)
	}

	// Demographics
	genderStats, _ := h.stuRepo.CountByGender(ctx, filterEntryYear, uint(classID), uint(majorID))
	majorStats, _ := h.stuRepo.CountByMajor(ctx, filterEntryYear, uint(classID), uint(majorID))
	classStats, _ := h.stuRepo.CountByClass(ctx, filterEntryYear, uint(classID), uint(majorID))

	// Trend
	paymentTrend, _ := h.finRepo.GetPaymentTrendByPeriod(ctx, pStart, pEnd, "daily", filterEntryYear, uint(classID), uint(majorID))

	// Summary
	summary, _ := h.finRepo.GetDashboardSummary(ctx, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID))
	if search != "" {
		summary = h.dashboardSummaryForSearch(ctx, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID), search)
	}

	// Methods & Efficacy
	paymentMethods, _ := h.finRepo.GetPaymentMethodsCount(ctx, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID))
	waEfficacy, _ := h.notifRepo.GetEfficacyStats(ctx, "whatsapp", pStart, pEnd)
	emailEfficacy, _ := h.notifRepo.GetEfficacyStats(ctx, "email", pStart, pEnd)

	// Critical Bills
	overdueBills, _ := h.finRepo.GetCriticalBills(ctx, "overdue", 5, filterEntryYear, uint(classID), uint(majorID), 0)
	dueSoonBills, _ := h.finRepo.GetCriticalBills(ctx, "due_soon", 5, filterEntryYear, uint(classID), uint(majorID), 0)

	billTypeID, _ := strconv.ParseUint(c.Query("bill_type_id"), 10, 32)

	// Recent Payments for Reports View
	recentPayments, totalPaymentsCount, _ := h.finRepo.GetRecentPayments(ctx, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID), uint(billTypeID), search, page, limit)

	// Recent Notifications
	recentNotifications, _, _ := h.notifRepo.GetDetailedLogs(ctx, 1, 5, "", search, "")

	// Available academic years
	activeAY, _, _ := h.ayRepo.FindAll(ctx, 1, 100, "", "", "")
	var availableYears []int
	for _, ay := range activeAY {
		availableYears = append(availableYears, ay.Year)
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"recent_payments":      recentPayments,
			"recent_notifications": recentNotifications,
			"total_payments_count": totalPaymentsCount,
			"stats": gin.H{
				"users": gin.H{
					"total":           totalUsers,
					"new_this_period": currUsers,
					"growth":          calculateGrowth(float64(currUsers), float64(prevUsers)),
				},
				"students": gin.H{
					"total_all": totalStudents,
					"total":     currStudents,
					"growth":    calculateGrowth(float64(currStudents), float64(prevStudents)),
				},
				"bills": gin.H{
					"total":  currBills,
					"growth": calculateGrowth(currBills, prevBills),
				},
				"payments": gin.H{
					"total":  currPayments,
					"growth": calculateGrowth(currPayments, prevPayments),
				},
				"paid_amount":   summary["total_paid_amount"],
				"unpaid_amount": summary["total_unpaid_amount"],
				"paid_count":    summary["paid_count"],
				"unpaid_count":  summary["unpaid_count"],
			},
			"demographics": gin.H{
				"gender":          genderStats,
				"major":           majorStats,
				"class":           classStats,
				"payment_methods": paymentMethods,
				"whatsapp":        waEfficacy,
				"email":           emailEfficacy,
			},
			"payment_trend": paymentTrend,
			"entry_years":   availableYears,
			"critical_bills": gin.H{
				"overdue":  overdueBills,
				"due_soon": dueSoonBills,
			},
		},
	})
}

func (h *DashboardHandler) GetCommunicationDetails(c *gin.Context) {
	ctx := c.Request.Context()
	status := strings.TrimSpace(strings.ToLower(c.Query("status")))
	channel := c.DefaultQuery("channel", "whatsapp")
	period := c.DefaultQuery("period", "all")

	// 1. Resolve Filter Params
	var filterEntryYear int
	if c.Query("academic_year_id") != "" {
		filterEntryYear, _ = strconv.Atoi(c.Query("academic_year_id"))
	} else if c.Query("entry_year") != "" {
		filterEntryYear, _ = strconv.Atoi(c.Query("entry_year"))
	}
	classID, _ := strconv.ParseUint(c.Query("class_id"), 10, 32)
	majorID, _ := strconv.ParseUint(c.Query("major_id"), 10, 32)

	// 2. Handle Periods
	var start, end time.Time
	now := time.Now()
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if period == "custom" && startDateStr != "" && endDateStr != "" {
		s, _ := time.Parse("2006-01-02", startDateStr)
		e, _ := time.Parse("2006-01-02", endDateStr)
		start = time.Date(s.Year(), s.Month(), s.Day(), 0, 0, 0, 0, time.Local)
		end = time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, time.Local)
	}

	refTime := now
	refDateStr := c.Query("ref_date")
	if refDateStr != "" {
		switch period {
		case "monthly":
			if t, err := time.Parse("2006-01", refDateStr); err == nil {
				refTime = t
			}
		case "daily":
			if t, err := time.Parse("2006-01-02", refDateStr); err == nil {
				refTime = t
			}
		case "yearly":
			if t, err := time.Parse("2006", refDateStr); err == nil {
				refTime = t
			}
		}
	}

	if period != "custom" && period != "all" {
		switch period {
		case "daily":
			start = time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 0, 0, 0, 0, refTime.Location())
			end = time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, start.Location())
		case "yearly":
			start = time.Date(refTime.Year(), 1, 1, 0, 0, 0, 0, refTime.Location())
			end = time.Date(start.Year(), 12, 31, 23, 59, 59, 0, start.Location())
		default: // monthly
			start = time.Date(refTime.Year(), refTime.Month(), 1, 0, 0, 0, 0, refTime.Location())
			end = start.AddDate(0, 1, 0).Add(-time.Second)
		}
	}

	q := h.db.NewSelect().
		TableExpr("notifications AS n").
		ColumnExpr("n.id, s.name as student_name, u.name as recipient_name, u.phone_number as recipient_phone, u.email as recipient_email, n.title, n.message, n.delivery_status, n.delivery_error, n.created_at, n.updated_at").
		ColumnExpr("CASE WHEN n.whatsapp_id IS NULL OR n.whatsapp_id = '' THEN 'email' ELSE 'whatsapp' END as channel").
		Join("JOIN users u ON n.user_id = u.id").
		Join("LEFT JOIN students s ON s.parent_id = u.id")

	if status != "" && status != "all" {
		q.Where("LOWER(n.delivery_status) = ?", status)
	}

	if period != "all" && !(period == "custom" && (startDateStr == "" || endDateStr == "")) {
		q.Where("n.created_at >= ? AND n.created_at <= ?", start, end)
	}

	if filterEntryYear > 0 {
		if filterEntryYear < 1000 {
			q.Where("s.entry_year = (SELECT year FROM academic_years WHERE id = ?)", filterEntryYear)
		} else {
			q.Where("s.entry_year = ?", filterEntryYear)
		}
	}
	if classID > 0 {
		q.Where("s.class_id = ?", classID)
	}
	if majorID > 0 {
		q.Where("s.major_id = ?", majorID)
	}

	if channel == "whatsapp" {
		q.Where("n.whatsapp_id IS NOT NULL AND n.whatsapp_id != ''")
	} else {
		q.Where("n.whatsapp_id IS NULL OR n.whatsapp_id = ''")
	}

	type CommDetail struct {
		ID             uint      `json:"id" bun:"id"`
		StudentName    *string   `json:"student_name" bun:"student_name"`
		RecipientName  string    `json:"recipient_name" bun:"recipient_name"`
		RecipientPhone string    `json:"recipient_phone" bun:"recipient_phone"`
		RecipientEmail string    `json:"recipient_email" bun:"recipient_email"`
		Channel        string    `json:"channel" bun:"channel"`
		Title          string    `json:"title" bun:"title"`
		Message        string    `json:"message" bun:"message"`
		DeliveryStatus string    `json:"delivery_status" bun:"delivery_status"`
		DeliveryError  *string   `json:"delivery_error" bun:"delivery_error"`
		CreatedAt      time.Time `json:"created_at" bun:"created_at"`
		UpdatedAt      time.Time `json:"updated_at" bun:"updated_at"`
	}

	var list []CommDetail
	err := q.Group("n.id").Order("n.created_at DESC").Scan(ctx, &list)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
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

func calculateGrowth(current, previous float64) float64 {
	if previous == 0 {
		if current > 0 {
			return 100
		}
		return 0
	}
	return ((current - previous) / previous) * 100
}

func (h *DashboardHandler) ExportGlobalReport(c *gin.Context) {
	ctx := c.Request.Context()
	tab := c.Query("tab")
	search := c.Query("search")
	period := c.DefaultQuery("period", "all")

	now := time.Now()
	var start, end time.Time

	// 1. Resolve Filter Params
	var filterEntryYear int
	if c.Query("academic_year_id") != "" {
		filterEntryYear, _ = strconv.Atoi(c.Query("academic_year_id"))
	} else if c.Query("entry_year") != "" {
		filterEntryYear, _ = strconv.Atoi(c.Query("entry_year"))
	}
	classID, _ := strconv.ParseUint(c.Query("class_id"), 10, 32)
	majorID, _ := strconv.ParseUint(c.Query("major_id"), 10, 32)
	billTypeID, _ := strconv.ParseUint(c.Query("bill_type_id"), 10, 32)

	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if h.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		newVals := map[string]interface{}{
			"tab": tab, "search": search, "period": period, "entry_year": filterEntryYear,
			"class_id": classID, "major_id": majorID, "bill_type_id": billTypeID,
			"start_date": startDateStr, "end_date": endDateStr,
		}
		_ = h.audit.Log(ctx, h.db, userID, userName, role, "EXPORT_GLOBAL_FINANCE_REPORT", "finance_reports", 0, nil, newVals, ipAddress, userAgent)
	}

	if period == "custom" && startDateStr != "" && endDateStr != "" {
		s, _ := time.Parse("2006-01-02", startDateStr)
		e, _ := time.Parse("2006-01-02", endDateStr)
		start = time.Date(s.Year(), s.Month(), s.Day(), 0, 0, 0, 0, time.Local)
		end = time.Date(e.Year(), e.Month(), e.Day(), 23, 59, 59, 0, time.Local)
	}

	refTime := now
	refDateStr := c.Query("ref_date")
	if refDateStr != "" {
		switch period {
		case "monthly":
			if t, err := time.Parse("2006-01", refDateStr); err == nil {
				refTime = t
			}
		case "daily":
			if t, err := time.Parse("2006-01-02", refDateStr); err == nil {
				refTime = t
			}
		case "yearly":
			if t, err := time.Parse("2006", refDateStr); err == nil {
				refTime = t
			}
		}
	}

	if period != "custom" && period != "all" {
		switch period {
		case "daily":
			start = time.Date(refTime.Year(), refTime.Month(), refTime.Day(), 0, 0, 0, 0, refTime.Location())
			end = time.Date(start.Year(), start.Month(), start.Day(), 23, 59, 59, 0, start.Location())
		case "yearly":
			start = time.Date(refTime.Year(), 1, 1, 0, 0, 0, 0, refTime.Location())
			end = time.Date(start.Year(), 12, 31, 23, 59, 59, 0, start.Location())
		default: // monthly
			start = time.Date(refTime.Year(), refTime.Month(), 1, 0, 0, 0, 0, refTime.Location())
			end = start.AddDate(0, 1, 0).Add(-time.Second)
		}
	}

	var pStart, pEnd *time.Time
	if period != "all" && !(period == "custom" && (startDateStr == "" || endDateStr == "")) {
		pStart = &start
		pEnd = &end
	}

	format := c.Query("format")
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
			list, _, _ := h.finRepo.GetArrearsPaged(ctx, 1, 1000000, filterEntryYear, uint(classID), uint(majorID), uint(billTypeID), search, pStart, pEnd)
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
			list, _, _ := h.finRepo.GetRecentPayments(ctx, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID), uint(billTypeID), search, 1, 1000000)
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
		_ = pdf.Output(&buf)
		c.Header("Content-Disposition", "attachment; filename=Laporan_Keuangan.pdf")
		c.Header("Content-Type", "application/pdf")
		c.Data(http.StatusOK, "application/pdf", buf.Bytes())
		return
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
		// Export Arrears
		list, _, _ := h.finRepo.GetArrearsPaged(ctx, 1, 1000000, filterEntryYear, uint(classID), uint(majorID), uint(billTypeID), search, pStart, pEnd)
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
		// Export Payments
		list, _, _ := h.finRepo.GetRecentPayments(ctx, pStart, pEnd, filterEntryYear, uint(classID), uint(majorID), uint(billTypeID), search, 1, 1000000)
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

	buf, _ := f.WriteToBuffer()

	c.Header("Content-Disposition", "attachment; filename=Laporan_Keuangan.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", buf.Bytes())
}
