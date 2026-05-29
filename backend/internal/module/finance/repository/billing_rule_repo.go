package repository

import (
	"context"
	"fmt"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

type BillingRuleRepo interface {
	Create(ctx context.Context, br *domain.BillingRule) error
	FindAll(ctx context.Context) ([]domain.BillingRule, error)
	FindAllPaged(ctx context.Context, page, limit int, search, status, generateStatus, sort string) ([]domain.BillingRule, int, error)
	FindByID(ctx context.Context, id uint) (*domain.BillingRule, error)
	Update(ctx context.Context, br *domain.BillingRule) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	Exists(ctx context.Context, btID uint, targetType string, targetID uint, classID *uint) (bool, error)
	ExistsExcludeID(ctx context.Context, btID uint, targetType string, targetID uint, classID *uint, excludeID uint) (bool, error)
	CountStudentBills(ctx context.Context, ruleID uint) (int, error)
}

type billingRuleRepo struct {
	db *bun.DB
}

func NewBillingRuleRepo(db *bun.DB) BillingRuleRepo {
	return &billingRuleRepo{db: db}
}

func (r *billingRuleRepo) Create(ctx context.Context, br *domain.BillingRule) error {
	_, err := r.db.NewInsert().Model(br).Exec(ctx)
	return err
}

func (r *billingRuleRepo) FindAll(ctx context.Context) ([]domain.BillingRule, error) {
	var list []domain.BillingRule
	err := r.db.NewSelect().
		Model(&list).
		ColumnExpr("br.*").
		ColumnExpr("bt.name as bill_type_name").
		ColumnExpr("COALESCE(CASE WHEN br.target_type = 'class' THEN (SELECT name FROM classes WHERE id = br.target_id) WHEN br.target_type = 'major' THEN (SELECT name FROM majors WHERE id = br.target_id) WHEN br.class_id IS NOT NULL THEN (SELECT name FROM classes WHERE id = br.class_id) ELSE NULL END, '') as class_name").
		ColumnExpr("(SELECT COUNT(*) FROM student_bills sb WHERE sb.billing_rule_id = br.id AND sb.deleted_at IS NULL AND sb.status != 'voided') as bill_count").
		Join("JOIN bill_types bt ON br.bill_type_id = bt.id").
		Where("br.deleted_at IS NULL").
		Order("br.created_at DESC").
		Scan(ctx)
	return list, err
}

func (r *billingRuleRepo) FindAllPaged(ctx context.Context, page, limit int, search, status, generateStatus, sort string) ([]domain.BillingRule, int, error) {
	var list []domain.BillingRule
	q := r.db.NewSelect().
		Model(&list).
		ColumnExpr("br.*").
		ColumnExpr("bt.name as bill_type_name").
		ColumnExpr("COALESCE(CASE WHEN br.target_type = 'class' THEN (SELECT name FROM classes WHERE id = br.target_id) WHEN br.target_type = 'major' THEN (SELECT name FROM majors WHERE id = br.target_id) WHEN br.class_id IS NOT NULL THEN (SELECT name FROM classes WHERE id = br.class_id) ELSE NULL END, '') as class_name").
		ColumnExpr("(SELECT COUNT(*) FROM student_bills sb WHERE sb.billing_rule_id = br.id AND sb.deleted_at IS NULL AND sb.status != 'voided') as bill_count").
		Join("LEFT JOIN bill_types bt ON br.bill_type_id = bt.id")

	if status == "trash" {
		q.WhereAllWithDeleted().Where("br.deleted_at IS NOT NULL")
	} else {
		q.Where("br.deleted_at IS NULL")
		switch status {
		case "active":
			q.Where("br.is_active = ?", true)
		case "inactive":
			q.Where("br.is_active = ?", false)
		}
	}

	switch generateStatus {
	case "generated":
		q.Where("(SELECT COUNT(*) FROM student_bills sb WHERE sb.billing_rule_id = br.id AND sb.deleted_at IS NULL AND sb.status != 'voided') > 0")
	case "not_generated":
		q.Where("(SELECT COUNT(*) FROM student_bills sb WHERE sb.billing_rule_id = br.id AND sb.deleted_at IS NULL AND sb.status != 'voided') = 0")
	}

	if search != "" {
		s := "%" + search + "%"
		q.Where("bt.name LIKE ? OR br.target_type LIKE ? OR COALESCE(CASE WHEN br.target_type = 'class' THEN (SELECT name FROM classes WHERE id = br.target_id) WHEN br.target_type = 'major' THEN (SELECT name FROM majors WHERE id = br.target_id) WHEN br.class_id IS NOT NULL THEN (SELECT name FROM classes WHERE id = br.class_id) ELSE NULL END, '') LIKE ?", s, s, s)
	}

	switch sort {
	case "created_asc":
		q.Order("br.created_at ASC")
	case "name_asc":
		q.Order("bt.name ASC")
	case "name_desc":
		q.Order("bt.name DESC")
	default:
		q.Order("br.created_at DESC")
	}

	total, err := q.Limit(limit).
		Offset((page - 1) * limit).
		ScanAndCount(ctx)

	return list, total, err
}

func (r *billingRuleRepo) FindByID(ctx context.Context, id uint) (*domain.BillingRule, error) {
	var br domain.BillingRule
	err := r.db.NewSelect().Model(&br).Where("id = ?", id).Scan(ctx)
	return &br, err
}

func (r *billingRuleRepo) Update(ctx context.Context, br *domain.BillingRule) error {
	_, err := r.db.NewUpdate().Model(br).WherePK().Exec(ctx)
	return err
}

func (r *billingRuleRepo) Delete(ctx context.Context, id uint) error {
	res, err := r.db.NewDelete().Model((*domain.BillingRule)(nil)).Where("id = ? AND deleted_at IS NULL", id).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("aturan tagihan tidak ditemukan atau sudah terhapus")
	}
	return nil
}

func (r *billingRuleRepo) Restore(ctx context.Context, id uint) error {
	res, err := r.db.NewUpdate().
		Model((*domain.BillingRule)(nil)).
		WhereAllWithDeleted().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Set("deleted_at = NULL").
		Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("aturan tagihan tidak ditemukan di riwayat penghapusan")
	}
	return nil
}

func (r *billingRuleRepo) ToggleStatus(ctx context.Context, id uint) error {
	_, err := r.db.NewUpdate().
		Model((*domain.BillingRule)(nil)).
		Set("is_active = NOT is_active").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *billingRuleRepo) BulkDelete(ctx context.Context, ids []uint) error {
	res, err := r.db.NewDelete().Model((*domain.BillingRule)(nil)).Where("id IN (?) AND deleted_at IS NULL", bun.In(ids)).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tidak ada data yang dapat dihapus")
	}
	return nil
}

func (r *billingRuleRepo) BulkRestore(ctx context.Context, ids []uint) error {
	res, err := r.db.NewUpdate().
		Model((*domain.BillingRule)(nil)).
		WhereAllWithDeleted().
		Set("deleted_at = NULL").
		Where("id IN (?) AND deleted_at IS NOT NULL", bun.In(ids)).
		Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tidak ada data yang perlu dipulihkan")
	}
	return nil
}

func (r *billingRuleRepo) Exists(ctx context.Context, btID uint, targetType string, targetID uint, classID *uint) (bool, error) {
	q := r.db.NewSelect().Model((*domain.BillingRule)(nil)).
		WhereAllWithDeleted().
		Where("bill_type_id = ? AND target_type = ? AND target_id = ?", btID, targetType, targetID)

	if classID != nil {
		q.Where("class_id = ?", *classID)
	} else {
		q.Where("class_id IS NULL")
	}

	return q.Exists(ctx)
}

func (r *billingRuleRepo) ExistsExcludeID(ctx context.Context, btID uint, targetType string, targetID uint, classID *uint, excludeID uint) (bool, error) {
	q := r.db.NewSelect().Model((*domain.BillingRule)(nil)).
		WhereAllWithDeleted().
		Where("bill_type_id = ? AND target_type = ? AND target_id = ? AND id != ?", btID, targetType, targetID, excludeID)

	if classID != nil {
		q.Where("class_id = ?", *classID)
	} else {
		q.Where("class_id IS NULL")
	}

	return q.Exists(ctx)
}

func (r *billingRuleRepo) CountStudentBills(ctx context.Context, ruleID uint) (int, error) {
	return r.db.NewSelect().Model((*domain.StudentBill)(nil)).Where("billing_rule_id = ? AND status != 'voided'", ruleID).Count(ctx)
}
