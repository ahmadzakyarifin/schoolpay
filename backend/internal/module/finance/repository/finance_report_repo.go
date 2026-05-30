package repository

import (
	"context"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

type FinanceReportRepo interface {
	GetTotalBillsByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (float64, error)
	GetTotalPaymentsByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (float64, error)
	GetPaymentTrendByPeriod(ctx context.Context, start, end *time.Time, interval string, academicYear int, classID, majorID uint) ([]map[string]interface{}, error)
	GetDashboardSummary(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]interface{}, error)
	GetArrearsPaged(ctx context.Context, page, limit int, academicYear int, classID, majorID, billTypeID uint, search string, start, end *time.Time) ([]domain.ArrearRecord, int, error)
	GetCriticalBills(ctx context.Context, status string, limit int, academicYear int, classID, majorID, billTypeID uint) ([]domain.CriticalBillRecord, error)
	GetRecentPayments(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID, billTypeID uint, search string, page, limit int) ([]map[string]interface{}, int, error)
	GetPaymentMethodsCount(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error)
}

type financeReportRepo struct {
	db *bun.DB
}

func NewFinanceReportRepo(db *bun.DB) FinanceReportRepo {
	return &financeReportRepo{db: db}
}

func (r *financeReportRepo) GetTotalBillsByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (float64, error) {
	var total float64
	q := r.db.NewSelect().
		Model((*domain.StudentBill)(nil)).
		ColumnExpr("COALESCE(SUM(sb.amount - sb.total_paid), 0)").
		Join("JOIN students s ON sb.student_id = s.id").
		Where("sb.status != 'paid' AND sb.status != 'voided'")

	if start != nil {
		q.Where("sb.created_at >= ?", start)
	}
	if end != nil {
		q.Where("sb.created_at <= ?", end)
	}
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
		q.Where("s.major_id = ?", majorID)
	}

	err := q.Scan(ctx, &total)
	return total, err
}

func (r *financeReportRepo) GetTotalPaymentsByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (float64, error) {
	var total float64
	q := r.db.NewSelect().
		Model((*domain.Payment)(nil)).
		ColumnExpr("COALESCE(SUM(p.amount), 0)").
		Join("JOIN students s ON p.student_id = s.id").
		Where("p.status = 'success'")

	if start != nil {
		q.Where("p.paid_at >= ?", start)
	}
	if end != nil {
		q.Where("p.paid_at <= ?", end)
	}
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
		q.Where("s.major_id = ?", majorID)
	}

	err := q.Scan(ctx, &total)
	return total, err
}

func (r *financeReportRepo) GetPaymentTrendByPeriod(ctx context.Context, start, end *time.Time, interval string, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	var datePart string
	switch interval {
	case "daily":
		datePart = "DATE(p.paid_at)"
	case "monthly":
		datePart = "DATE_FORMAT(p.paid_at, '%Y-%m')"
	case "yearly":
		datePart = "YEAR(p.paid_at)"
	default:
		datePart = "DATE(p.paid_at)"
	}

	q := r.db.NewSelect().
		Model((*domain.Payment)(nil)).
		ColumnExpr(datePart + " as date, SUM(p.amount) as total").
		Join("JOIN students s ON p.student_id = s.id").
		Where("p.status = 'success'")

	if start != nil {
		q.Where("p.paid_at >= ?", start)
	}
	if end != nil {
		q.Where("p.paid_at <= ?", end)
	}
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
		q.Where("s.major_id = ?", majorID)
	}

	var results []map[string]interface{}
	err := q.GroupExpr(datePart).Order("date ASC").Scan(ctx, &results)
	return results, err
}

func (r *financeReportRepo) GetDashboardSummary(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	qUnpaid := r.db.NewSelect().
		Model((*domain.StudentBill)(nil)).
		ColumnExpr("COALESCE(SUM(sb.amount - sb.total_paid), 0) as unpaid_amount").
		ColumnExpr("COUNT(sb.id) as unpaid_count").
		Join("JOIN students s ON sb.student_id = s.id").
		Where("sb.status != 'paid' AND sb.status != 'voided'")

	if end != nil {
		qUnpaid.Where("sb.created_at <= ?", end)
	}
	if academicYear > 0 {
		if academicYear < 1000 {
			qUnpaid.Where("s.entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
		} else {
			qUnpaid.Where("s.entry_year = ?", academicYear)
		}
	}
	if classID > 0 {
		qUnpaid.Where("s.class_id = ?", classID)
	}
	if majorID > 0 {
		qUnpaid.Where("s.major_id = ?", majorID)
	}

	var summaryUnpaid struct {
		UnpaidAmount float64 `bun:"unpaid_amount"`
		UnpaidCount  int     `bun:"unpaid_count"`
	}
	_ = qUnpaid.Scan(ctx, &summaryUnpaid)

	qPaid := r.db.NewSelect().
		Model((*domain.Payment)(nil)).
		ColumnExpr("COALESCE(SUM(p.amount), 0) as paid_amount").
		ColumnExpr("COUNT(p.id) as paid_count").
		Join("JOIN students s ON p.student_id = s.id").
		Where("p.status = 'success'")

	if start != nil {
		qPaid.Where("p.paid_at >= ?", start)
	}
	if end != nil {
		qPaid.Where("p.paid_at <= ?", end)
	}
	if academicYear > 0 {
		if academicYear < 1000 {
			qPaid.Where("s.entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
		} else {
			qPaid.Where("s.entry_year = ?", academicYear)
		}
	}
	if classID > 0 {
		qPaid.Where("s.class_id = ?", classID)
	}
	if majorID > 0 {
		qPaid.Where("s.major_id = ?", majorID)
	}

	var summaryPaid struct {
		PaidAmount float64 `bun:"paid_amount"`
		PaidCount  int     `bun:"paid_count"`
	}
	_ = qPaid.Scan(ctx, &summaryPaid)

	res["total_paid_amount"] = summaryPaid.PaidAmount
	res["total_unpaid_amount"] = summaryUnpaid.UnpaidAmount
	res["paid_count"] = summaryPaid.PaidCount
	res["unpaid_count"] = summaryUnpaid.UnpaidCount

	return res, nil
}

func (r *financeReportRepo) GetArrearsPaged(ctx context.Context, page, limit int, academicYear int, classID, majorID, billTypeID uint, search string, start, end *time.Time) ([]domain.ArrearRecord, int, error) {
	var list []domain.ArrearRecord
	q := r.db.NewSelect().
		Model((*domain.StudentBill)(nil)).
		ColumnExpr("sb.id as id, s.name as student_name, c.name as class_name, bt.name as bill_name, COALESCE(sb.period, '') as period, sb.amount as amount, sb.total_paid as total_paid, sb.status as status, sb.due_date as due_date").
		ColumnExpr("COALESCE(sb.period_start_date, STR_TO_DATE(CONCAT(sb.period, '-01'), '%Y-%m-%d'), br.start_date) as start_date").
		ColumnExpr("COALESCE(sb.period_end_date, LAST_DAY(STR_TO_DATE(CONCAT(sb.period, '-01'), '%Y-%m-%d')), sb.end_date, br.end_date) as end_date").
		Join("JOIN students s ON sb.student_id = s.id").
		Join("JOIN classes c ON s.class_id = c.id").
		Join("JOIN bill_types bt ON sb.bill_type_id = bt.id").
		Join("LEFT JOIN billing_rules br ON sb.billing_rule_id = br.id").
		Where("sb.status != 'paid' AND sb.status != 'voided'")

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
		q.Where("s.major_id = ?", majorID)
	}
	if billTypeID > 0 {
		q.Where("sb.bill_type_id = ?", billTypeID)
	}
	if search != "" {
		q.Where("s.name LIKE ? OR s.nis LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if start != nil {
		q.Where("sb.due_date >= ?", start)
	}
	if end != nil {
		q.Where("sb.due_date <= ?", end)
	}

	total, err := q.Order("sb.due_date ASC").
		Limit(limit).
		Offset((page-1)*limit).
		ScanAndCount(ctx, &list)

	return list, total, err
}

func (r *financeReportRepo) GetRecentPayments(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID, billTypeID uint, search string, page, limit int) ([]map[string]interface{}, int, error) {
	var list []map[string]interface{}
	q := r.db.NewSelect().
		Model((*domain.Payment)(nil)).
		ColumnExpr("p.id, p.transaction_ref, p.channel, p.created_by, p.amount, COALESCE(p.deposit_applied, 0) as deposit_applied, (p.amount - COALESCE(p.deposit_applied, 0)) as cash_or_gateway_amount, p.method, p.paid_at as created_at, s.name as student_name, COALESCE(c.name, '') as class_name").
		ColumnExpr("CAST(COALESCE((SELECT GROUP_CONCAT(COALESCE(NULLIF(sb.name, ''), IF(COALESCE(sb.period, '') != '', CONCAT(bt.name, ' ', sb.period), bt.name)) ORDER BY sb.due_date ASC SEPARATOR '||') FROM payment_details pd JOIN student_bills sb ON pd.student_bill_id = sb.id JOIN bill_types bt ON sb.bill_type_id = bt.id WHERE pd.payment_id = p.id), '') AS CHAR) as bill_type_names").
		ColumnExpr("CAST(COALESCE((SELECT GROUP_CONCAT(CONCAT(COALESCE(NULLIF(sb.name, ''), IF(COALESCE(sb.period, '') != '', CONCAT(bt.name, ' ', sb.period), bt.name)), '::', COALESCE(sb.period, ''), '::', FORMAT(pd.amount, 0)) ORDER BY sb.due_date ASC SEPARATOR '||') FROM payment_details pd JOIN student_bills sb ON pd.student_bill_id = sb.id JOIN bill_types bt ON sb.bill_type_id = bt.id WHERE pd.payment_id = p.id), '') AS CHAR) as bill_type_details").
		ColumnExpr("(SELECT COUNT(*) FROM payment_details pd WHERE pd.payment_id = p.id) as bill_item_count").
		Join("JOIN students s ON p.student_id = s.id").
		Join("LEFT JOIN classes c ON s.class_id = c.id").
		Where("p.status = 'success'")

	if billTypeID > 0 {
		q.Join("JOIN payment_details pd ON pd.payment_id = p.id").
			Join("JOIN student_bills sb ON pd.student_bill_id = sb.id").
			Where("sb.bill_type_id = ?", billTypeID)
	}

	if start != nil {
		q.Where("p.paid_at >= ?", start)
	}
	if end != nil {
		q.Where("p.paid_at <= ?", end)
	}
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
		q.Where("s.major_id = ?", majorID)
	}
	if search != "" {
		q.Where("s.name LIKE ? OR s.nis LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	total, err := q.Group("p.id", "p.transaction_ref", "p.channel", "p.created_by", "p.amount", "p.deposit_applied", "p.method", "p.paid_at", "s.name", "c.name").
		Order("p.paid_at DESC").
		Limit(limit).
		Offset((page-1)*limit).
		ScanAndCount(ctx, &list)

	return list, total, err
}

func (r *financeReportRepo) GetCriticalBills(ctx context.Context, status string, limit int, academicYear int, classID, majorID, billTypeID uint) ([]domain.CriticalBillRecord, error) {
	var results []domain.CriticalBillRecord
	q := r.db.NewSelect().
		Model((*domain.StudentBill)(nil)).
		ColumnExpr("sb.id, sb.student_id, sb.bill_type_id, sb.amount, sb.total_paid, sb.due_date").
		ColumnExpr("s.name as student_name").
		ColumnExpr("bt.name as bill_type_name").
		ColumnExpr("COALESCE(u.name, '') as parent_name").
		ColumnExpr("COALESCE(u.phone_number, '') as parent_phone").
		ColumnExpr("COALESCE(sb.period_start_date, STR_TO_DATE(CONCAT(sb.period, '-01'), '%Y-%m-%d'), br.start_date) as start_date").
		ColumnExpr("COALESCE(sb.period_end_date, LAST_DAY(STR_TO_DATE(CONCAT(sb.period, '-01'), '%Y-%m-%d')), sb.end_date, br.end_date) as end_date").
		Join("JOIN students s ON sb.student_id = s.id").
		Join("JOIN bill_types bt ON sb.bill_type_id = bt.id").
		Join("LEFT JOIN users u ON s.parent_id = u.id").
		Join("LEFT JOIN billing_rules br ON sb.billing_rule_id = br.id")

	switch status {
	case "overdue":
		q.Where("sb.status != 'paid' AND sb.status != 'voided' AND DATE(sb.due_date) < CURDATE()")
	case "due_soon":
		q.Where("sb.status != 'paid' AND sb.status != 'voided' AND DATE(sb.due_date) >= CURDATE()")
	}

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
		q.Where("s.major_id = ?", majorID)
	}
	if billTypeID > 0 {
		q.Where("sb.bill_type_id = ?", billTypeID)
	}

	err := q.Order("sb.due_date ASC").Limit(limit).Scan(ctx, &results)
	return results, err
}

func (r *financeReportRepo) GetPaymentMethodsCount(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	q := r.db.NewSelect().
		Model((*domain.Payment)(nil)).
		ColumnExpr("p.method as method, COUNT(*) as count").
		Join("JOIN students s ON p.student_id = s.id").
		Where("p.status = 'success'")

	if start != nil {
		q.Where("p.paid_at >= ?", start)
	}
	if end != nil {
		q.Where("p.paid_at <= ?", end)
	}
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
		q.Where("s.major_id = ?", majorID)
	}

	var results []map[string]interface{}
	err := q.Group("p.method").Scan(ctx, &results)
	return results, err
}
