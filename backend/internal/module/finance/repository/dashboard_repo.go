package repository

import (
	"context"
	"reflect"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

type DashboardRepo interface {
	CountNewUsersByPeriod(ctx context.Context, start, end *time.Time) (int, error)
	CountTotalUsersUpTo(ctx context.Context, end *time.Time) (int, error)
	CountUsersByRoleUpTo(ctx context.Context, end *time.Time, role string) (int, error)
	CountStudentsByStatusUpTo(ctx context.Context, end *time.Time, status string, academicYear int, classID, majorID uint) (int, error)
	GetGenderStats(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]int, error)
	GetEfficacyStats(ctx context.Context, channel string, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]int, error)
	GetRecentNotifications(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error)
	GetCommunicationDetails(ctx context.Context, status, channel string, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error)
	GetCriticalBills(ctx context.Context, status string, start, end *time.Time, academicYear int, classID, majorID, billTypeID uint) ([]map[string]interface{}, error)

	GetTotalBillsByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (float64, error)
	GetTotalPaymentsByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (float64, error)
	GetPaymentTrendByPeriod(ctx context.Context, start, end *time.Time, interval string, academicYear int, classID, majorID uint) ([]map[string]interface{}, error)
	GetDashboardSummary(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]interface{}, error)
	GetRecentPayments(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID, billTypeID uint, search string, page, limit int) ([]map[string]interface{}, int, error)
	GetPaymentMethodsCount(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error)
	GetArrearsPaged(ctx context.Context, page, limit int, academicYear int, classID, majorID, billTypeID uint, search string, start, end *time.Time) ([]domain.ArrearRecord, int, error)
}

type dashboardRepo struct {
	db *bun.DB
}

func NewDashboardRepo(db *bun.DB) DashboardRepo {
	return &dashboardRepo{db: db}
}

func convertBytesToStrings(slice []map[string]interface{}) {
	for _, item := range slice {
		for k, v := range item {
			if v == nil {
				continue
			}

			// first convert []byte to string
			if b, ok := v.([]byte); ok {
				item[k] = string(b)
				continue
			}

			// check for database/sql Null types using reflection
			val := reflect.ValueOf(v)
			if val.Kind() == reflect.Ptr {
				if val.IsNil() {
					item[k] = nil
					continue
				}
				val = val.Elem()
			}

			if val.Kind() == reflect.Struct {
				validField := val.FieldByName("Valid")
				if validField.IsValid() && validField.Kind() == reflect.Bool {
					if !validField.Bool() {
						if val.FieldByName("String").IsValid() {
							item[k] = ""
						} else {
							item[k] = nil
						}
					} else {
						// Extract the value field: check commonly used names like String, Int64, Int32, Float64, Bool, Time
						for _, fieldName := range []string{"String", "Int64", "Int32", "Float64", "Bool", "Time"} {
							f := val.FieldByName(fieldName)
							if f.IsValid() {
								item[k] = f.Interface()
								break
							}
						}
					}
				}
			}
		}
	}
}

func (r *dashboardRepo) CountNewUsersByPeriod(ctx context.Context, start, end *time.Time) (int, error) {
	q := r.db.NewSelect().TableExpr("users AS u").Where("u.is_active = 1 AND u.deleted_at IS NULL")
	if start != nil {
		q.Where("u.created_at >= ?", start)
	}
	if end != nil {
		q.Where("u.created_at <= ?", end)
	}
	return q.Count(ctx)
}

func (r *dashboardRepo) CountTotalUsersUpTo(ctx context.Context, end *time.Time) (int, error) {
	q := r.db.NewSelect().TableExpr("users AS u").Where("u.is_active = 1 AND u.deleted_at IS NULL")
	if end != nil {
		q.Where("u.created_at <= ?", end)
	}
	return q.Count(ctx)
}

func (r *dashboardRepo) CountUsersByRoleUpTo(ctx context.Context, end *time.Time, role string) (int, error) {
	q := r.db.NewSelect().TableExpr("users AS u").Where("u.is_active = 1 AND u.deleted_at IS NULL AND u.role = ?", role)
	if end != nil {
		q.Where("u.created_at <= ?", end)
	}
	return q.Count(ctx)
}

func (r *dashboardRepo) CountStudentsByStatusUpTo(ctx context.Context, end *time.Time, status string, academicYear int, classID, majorID uint) (int, error) {
	q := r.db.NewSelect().TableExpr("students AS s").Where("s.deleted_at IS NULL")
	switch status {
	case "active":
		q.Where("s.status = 'active'")
	case "inactive":
		q.Where("s.status != 'active'")
	}
	if end != nil {
		q.Where("s.created_at <= ?", end)
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
	return q.Count(ctx)
}

func (r *dashboardRepo) GetGenderStats(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]int, error) {
	var genderStatsMap []struct {
		Gender string `bun:"gender"`
		Count  int    `bun:"count"`
	}
	gq := r.db.NewSelect().
		TableExpr("students AS s").
		ColumnExpr("s.gender, COUNT(*) as count").
		Where("s.status = 'active' AND s.deleted_at IS NULL")

	if start != nil {
		gq.Where("s.created_at >= ?", start)
	}
	if end != nil {
		gq.Where("s.created_at <= ?", end)
	}
	if academicYear > 0 {
		if academicYear < 1000 {
			gq.Where("s.entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
		} else {
			gq.Where("s.entry_year = ?", academicYear)
		}
	}
	if classID > 0 {
		gq.Where("s.class_id = ?", classID)
	}
	if majorID > 0 {
		gq.Where("s.major_id = ?", majorID)
	}

	err := gq.Group("s.gender").Scan(ctx, &genderStatsMap)
	if err != nil {
		return nil, err
	}

	genderStats := map[string]int{"L": 0, "P": 0}
	for _, s := range genderStatsMap {
		genderStats[s.Gender] = s.Count
	}
	return genderStats, nil
}

func (r *dashboardRepo) GetEfficacyStats(ctx context.Context, channel string, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]int, error) {
	var stats []struct {
		Status string `bun:"delivery_status"`
		Count  int    `bun:"count"`
	}
	eq := r.db.NewSelect().
		TableExpr("notifications AS n").
		ColumnExpr("n.delivery_status, COUNT(DISTINCT n.id) as count").
		Join("LEFT JOIN users u ON n.user_id = u.id")

	if start != nil {
		eq.Where("n.created_at >= ?", start)
	}
	if end != nil {
		eq.Where("n.created_at <= ?", end)
	}
	if channel == "whatsapp" {
		eq.Where("LOWER(COALESCE(NULLIF(n.channel, ''), CASE WHEN n.whatsapp_id IS NOT NULL AND n.whatsapp_id != '' THEN 'whatsapp' ELSE 'email' END)) = 'whatsapp'")
	} else {
		eq.Where("LOWER(COALESCE(NULLIF(n.channel, ''), CASE WHEN n.whatsapp_id IS NOT NULL AND n.whatsapp_id != '' THEN 'whatsapp' ELSE 'email' END)) = 'email'")
	}
	if academicYear > 0 || classID > 0 || majorID > 0 {
		eq.Join("JOIN students s ON s.parent_id = u.id")
		if academicYear > 0 {
			if academicYear < 1000 {
				eq.Where("s.entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
			} else {
				eq.Where("s.entry_year = ?", academicYear)
			}
		}
		if classID > 0 {
			eq.Where("s.class_id = ?", classID)
		}
		if majorID > 0 {
			eq.Where("s.major_id = ?", majorID)
		}
	}
	err := eq.Group("n.delivery_status").Scan(ctx, &stats)
	if err != nil {
		return nil, err
	}

	normalizeStatus := func(status string) string {
		switch strings.ToUpper(strings.TrimSpace(status)) {
		case "SUCCESS":
			return "sent"
		case "DELIVERED":
			return "delivered"
		case "READ":
			return "read"
		case "FAILED", "ERROR":
			return "failed"
		case "PENDING", "":
			return "pending"
		default:
			return "sent"
		}
	}

	result := make(map[string]int)
	statuses := []string{"pending", "sent", "failed"}
	if channel == "whatsapp" {
		statuses = []string{"pending", "sent", "delivered", "read", "failed"}
	}
	for _, status := range statuses {
		result[status] = 0
	}
	for _, s := range stats {
		normalized := strings.ToLower(normalizeStatus(s.Status))
		if channel != "whatsapp" && (normalized == "delivered" || normalized == "read") {
			normalized = "sent"
		}
		if _, ok := result[normalized]; ok {
			result[normalized] += s.Count
		}
	}
	return result, nil
}

func (r *dashboardRepo) GetRecentNotifications(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	var recentNotifications []map[string]interface{}
	nq := r.db.NewSelect().
		TableExpr("notifications AS n").
		ColumnExpr("n.id, n.user_id, n.title, n.message, n.type, n.is_read, n.whatsapp_id, n.delivery_error, n.created_at, n.updated_at").
		ColumnExpr("CASE WHEN n.delivery_status IN ('SUCCESS') THEN 'sent' WHEN n.delivery_status IS NULL OR n.delivery_status = '' THEN 'pending' ELSE LOWER(n.delivery_status) END as delivery_status").
		ColumnExpr("COALESCE(u.name, '') as recipient_name").
		ColumnExpr("COALESCE(u.phone_number, '') as recipient_phone").
		ColumnExpr("COALESCE(u.email, '') as recipient_email").
		ColumnExpr("COALESCE(NULLIF(n.channel, ''), CASE WHEN n.whatsapp_id IS NULL OR n.whatsapp_id = '' THEN 'email' ELSE 'whatsapp' END) as channel").
		Join("LEFT JOIN users u ON n.user_id = u.id")

	if start != nil {
		nq.Where("n.created_at >= ?", start)
	}
	if end != nil {
		nq.Where("n.created_at <= ?", end)
	}
	if academicYear > 0 || classID > 0 || majorID > 0 {
		nq.Join("JOIN students s ON s.parent_id = u.id")
		if academicYear > 0 {
			if academicYear < 1000 {
				nq.Where("s.entry_year = (SELECT year FROM academic_years WHERE id = ?)", academicYear)
			} else {
				nq.Where("s.entry_year = ?", academicYear)
			}
		}
		if classID > 0 {
			nq.Where("s.class_id = ?", classID)
		}
		if majorID > 0 {
			nq.Where("s.major_id = ?", majorID)
		}
		nq.Group("n.id")
	}
	err := nq.Order("n.created_at DESC").Limit(5).Scan(ctx, &recentNotifications)
	if err == nil {
		convertBytesToStrings(recentNotifications)
	}
	return recentNotifications, err
}

func (r *dashboardRepo) GetCommunicationDetails(ctx context.Context, status, channel string, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	var list []map[string]interface{}
	q := r.db.NewSelect().
		TableExpr("notifications AS n").
		ColumnExpr("n.id, s.name as student_name, u.name as recipient_name, u.phone_number as recipient_phone, u.email as recipient_email, n.title, n.message, n.delivery_status, n.delivery_error, n.created_at, n.updated_at").
		ColumnExpr("COALESCE(NULLIF(n.channel, ''), CASE WHEN n.whatsapp_id IS NULL OR n.whatsapp_id = '' THEN 'email' ELSE 'whatsapp' END) as channel").
		Join("JOIN users u ON n.user_id = u.id").
		Join("LEFT JOIN students s ON s.parent_id = u.id")

	if status != "" && status != "all" {
		q.Where("LOWER(n.delivery_status) = ?", status)
	}

	if start != nil && end != nil {
		q.Where("n.created_at >= ? AND n.created_at <= ?", start, end)
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

	if channel == "whatsapp" {
		q.Where("LOWER(COALESCE(NULLIF(n.channel, ''), CASE WHEN n.whatsapp_id IS NOT NULL AND n.whatsapp_id != '' THEN 'whatsapp' ELSE 'email' END)) = 'whatsapp'")
	} else {
		q.Where("LOWER(COALESCE(NULLIF(n.channel, ''), CASE WHEN n.whatsapp_id IS NOT NULL AND n.whatsapp_id != '' THEN 'whatsapp' ELSE 'email' END)) = 'email'")
	}

	err := q.Group("n.id").Order("n.created_at DESC").Scan(ctx, &list)
	if err == nil {
		convertBytesToStrings(list)
	}
	return list, err
}

func (r *dashboardRepo) GetCriticalBills(ctx context.Context, status string, start, end *time.Time, academicYear int, classID, majorID, billTypeID uint) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	q := r.db.NewSelect().
		TableExpr("student_bills AS sb").
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

	if start != nil {
		q.Where("sb.due_date >= ?", start)
	}
	if end != nil {
		q.Where("sb.due_date <= ?", end)
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

	err := q.Order("sb.due_date ASC").Limit(5).Scan(ctx, &results)
	if err == nil {
		convertBytesToStrings(results)
	}
	return results, err
}

func (r *dashboardRepo) GetTotalBillsByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (float64, error) {
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

func (r *dashboardRepo) GetTotalPaymentsByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (float64, error) {
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

func (r *dashboardRepo) GetPaymentTrendByPeriod(ctx context.Context, start, end *time.Time, interval string, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
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

func (r *dashboardRepo) GetDashboardSummary(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]interface{}, error) {
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

func (r *dashboardRepo) GetRecentPayments(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID, billTypeID uint, search string, page, limit int) ([]map[string]interface{}, int, error) {
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

	if err == nil {
		convertBytesToStrings(list)
	}
	return list, total, err
}

func (r *dashboardRepo) GetPaymentMethodsCount(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
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
	if err == nil {
		convertBytesToStrings(results)
	}
	return results, err
}

func (r *dashboardRepo) GetArrearsPaged(ctx context.Context, page, limit int, academicYear int, classID, majorID, billTypeID uint, search string, start, end *time.Time) ([]domain.ArrearRecord, int, error) {
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
