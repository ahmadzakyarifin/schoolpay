package repository

import (
	"context"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

type StudentBillRepo interface {
	ExistsByPeriod(ctx context.Context, studentID, billTypeID uint, period string) (bool, error)
	ExistsByPeriodExcludeID(ctx context.Context, studentID, billTypeID uint, period string, excludeID uint) (bool, error)
	Create(ctx context.Context, db bun.IDB, sb *domain.StudentBill) error
	FindByStudent(ctx context.Context, studentID uint) ([]domain.StudentBill, error)
	FindByParent(ctx context.Context, parentID uint) ([]domain.StudentBill, error)
	FindAll(ctx context.Context, search, sort string) ([]domain.StudentBill, error)
	FindByID(ctx context.Context, id uint) (*domain.StudentBill, error)
	UpdateStatus(ctx context.Context, db bun.IDB, id uint, status string, totalPaid float64) error
	Update(ctx context.Context, db bun.IDB, sb *domain.StudentBill) error
	FindUnpaidBillsByStudent(ctx context.Context, studentID uint) ([]domain.StudentBill, error)
	Delete(ctx context.Context, id uint) error
	FindForReminder(ctx context.Context, dueInDays int, overdueOnly bool) ([]domain.StudentBill, error)
}

type studentBillRepo struct {
	db *bun.DB
}

func NewStudentBillRepo(db *bun.DB) StudentBillRepo {
	return &studentBillRepo{db: db}
}

func (r *studentBillRepo) ExistsByPeriod(ctx context.Context, studentID, billTypeID uint, period string) (bool, error) {
	return r.db.NewSelect().Model((*domain.StudentBill)(nil)).
		Where("student_id = ? AND bill_type_id = ? AND COALESCE(period, '') = ? AND status != 'voided'", studentID, billTypeID, period).
		Exists(ctx)
}

func (r *studentBillRepo) ExistsByPeriodExcludeID(ctx context.Context, studentID, billTypeID uint, period string, excludeID uint) (bool, error) {
	return r.db.NewSelect().Model((*domain.StudentBill)(nil)).
		Where("student_id = ? AND bill_type_id = ? AND COALESCE(period, '') = ? AND status != 'voided' AND id != ?", studentID, billTypeID, period, excludeID).
		Exists(ctx)
}

func (r *studentBillRepo) Create(ctx context.Context, db bun.IDB, sb *domain.StudentBill) error {
	_, err := db.NewInsert().Model(sb).Exec(ctx)
	return err
}

func (r *studentBillRepo) FindByStudent(ctx context.Context, studentID uint) ([]domain.StudentBill, error) {
	var list []domain.StudentBill
	err := r.db.NewSelect().
		Model(&list).
		ColumnExpr("sb.*").
		ColumnExpr("s.deposit_balance as deposit_balance").
		ColumnExpr("bt.name as bill_type_name, br.allow_installment, br.max_installment, br.start_date as rule_start_date, br.end_date as rule_end_date").
		Join("JOIN students s ON sb.student_id = s.id").
		Join("JOIN bill_types bt ON sb.bill_type_id = bt.id").
		Join("LEFT JOIN billing_rules br ON sb.billing_rule_id = br.id").
		Where("sb.student_id = ?", studentID).
		Where("sb.status != ?", "voided").
		Order("sb.due_date ASC").
		Scan(ctx)
	return list, err
}

func (r *studentBillRepo) FindUnpaidBillsByStudent(ctx context.Context, studentID uint) ([]domain.StudentBill, error) {
	var list []domain.StudentBill
	err := r.db.NewSelect().
		Model(&list).
		ColumnExpr("sb.*").
		ColumnExpr("bt.name as bill_type_name, br.allow_installment, br.max_installment").
		Join("JOIN bill_types bt ON sb.bill_type_id = bt.id").
		Join("LEFT JOIN billing_rules br ON sb.billing_rule_id = br.id").
		Where("sb.student_id = ?", studentID).
		Where("sb.status IN ('unpaid', 'partial', 'overdue')").
		Order("sb.due_date ASC").
		Scan(ctx)
	return list, err
}

func (r *studentBillRepo) FindByParent(ctx context.Context, parentID uint) ([]domain.StudentBill, error) {
	var list []domain.StudentBill
	err := r.db.NewSelect().
		Model(&list).
		ColumnExpr("sb.*").
		ColumnExpr("s.name as student_name, s.deposit_balance as deposit_balance, bt.name as bill_type_name, br.allow_installment, br.max_installment").
		Join("JOIN students s ON sb.student_id = s.id").
		Join("JOIN bill_types bt ON sb.bill_type_id = bt.id").
		Join("LEFT JOIN billing_rules br ON sb.billing_rule_id = br.id").
		Join("JOIN parent_students ps ON s.id = ps.student_id").
		Where("ps.parent_id = ?", parentID).
		Where("sb.status != ?", "voided").
		Order("sb.due_date ASC").
		Scan(ctx)
	return list, err
}

func (r *studentBillRepo) FindAll(ctx context.Context, search, sort string) ([]domain.StudentBill, error) {
	var list []domain.StudentBill
	q := r.db.NewSelect().
		Model(&list).
		ColumnExpr("sb.*").
		ColumnExpr("s.name as student_name, s.deposit_balance as deposit_balance, bt.name as bill_type_name, br.allow_installment, br.max_installment").
		Join("JOIN students s ON sb.student_id = s.id").
		Join("JOIN bill_types bt ON sb.bill_type_id = bt.id").
		Join("LEFT JOIN billing_rules br ON sb.billing_rule_id = br.id").
		Where("sb.status != ?", "voided")

	if search != "" {
		s := "%" + search + "%"
		q.Where("s.name LIKE ? OR s.nis LIKE ? OR s.nisn LIKE ?", s, s, s)
	}

	switch sort {
	case "created_asc":
		q.Order("sb.created_at ASC")
	case "name_asc":
		q.Order("s.name ASC")
	case "name_desc":
		q.Order("s.name DESC")
	default:
		q.Order("sb.created_at DESC")
	}

	err := q.Scan(ctx)
	return list, err
}

func (r *studentBillRepo) FindByID(ctx context.Context, id uint) (*domain.StudentBill, error) {
	var sb domain.StudentBill
	err := r.db.NewSelect().Model(&sb).
		ColumnExpr("sb.*").
		ColumnExpr("bt.name as bill_type_name, br.allow_installment, br.max_installment").
		Join("JOIN bill_types bt ON sb.bill_type_id = bt.id").
		Join("LEFT JOIN billing_rules br ON sb.billing_rule_id = br.id").
		Where("sb.id = ?", id).Scan(ctx)
	return &sb, err
}

func (r *studentBillRepo) UpdateStatus(ctx context.Context, db bun.IDB, id uint, status string, totalPaid float64) error {
	_, err := db.NewUpdate().
		Model((*domain.StudentBill)(nil)).
		Set("status = ?", status).
		Set("total_paid = ?", totalPaid).
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *studentBillRepo) Update(ctx context.Context, db bun.IDB, sb *domain.StudentBill) error {
	_, err := db.NewUpdate().Model(sb).WherePK().Exec(ctx)
	return err
}

func (r *studentBillRepo) Delete(ctx context.Context, id uint) error {
	_, err := r.db.NewDelete().Model((*domain.StudentBill)(nil)).Where("id = ?", id).Exec(ctx)
	return err
}

func (r *studentBillRepo) FindForReminder(ctx context.Context, dueInDays int, overdueOnly bool) ([]domain.StudentBill, error) {
	var list []domain.StudentBill
	q := r.db.NewSelect().
		Model(&list).
		ColumnExpr("sb.*").
		ColumnExpr("s.name as student_name, bt.name as bill_type_name").
		Join("JOIN students s ON sb.student_id = s.id").
		Join("JOIN bill_types bt ON sb.bill_type_id = bt.id").
		Where("sb.status IN ('unpaid', 'partial', 'overdue')")

	if overdueOnly {
		q.Where("sb.due_date < NOW()").
			Where("(sb.last_notified_at IS NULL OR sb.last_notified_at < DATE_SUB(NOW(), INTERVAL 3 DAY))")
	} else {
		q.Where("DATE(sb.due_date) = DATE(DATE_ADD(NOW(), INTERVAL ? DAY))", dueInDays).
			Where("(sb.last_notified_at IS NULL OR DATE(sb.last_notified_at) != DATE(NOW()))")
	}

	err := q.Scan(ctx)
	return list, err
}
