package repository

import (
	"context"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

type FinanceReportRepo interface {
	GetPaymentTrendByPeriod(ctx context.Context, start, end *time.Time, interval string, academicYear int, classID, majorID uint) ([]map[string]interface{}, error)
	GetArrearsPaged(ctx context.Context, page, limit int, academicYear int, classID, majorID, billTypeID uint, search string, start, end *time.Time) ([]domain.ArrearRecord, int, error)
}

type financeReportRepo struct {
	db *bun.DB
}

func NewFinanceReportRepo(db *bun.DB) FinanceReportRepo {
	return &financeReportRepo{db: db}
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
	if err == nil {
		convertBytesToStrings(results)
	}
	return results, err
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

