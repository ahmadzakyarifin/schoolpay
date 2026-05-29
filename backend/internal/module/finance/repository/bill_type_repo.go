package repository

import (
	"context"
	"fmt"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

type BillTypeRepo interface {
	Create(ctx context.Context, b *domain.BillType) error
	FindAll(ctx context.Context) ([]domain.BillType, error)
	FindAllPaged(ctx context.Context, page, limit int, search, filterType, status, sort string) ([]domain.BillType, int, error)
	Update(ctx context.Context, b *domain.BillType) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	Exists(ctx context.Context, name string) (bool, error)
	ExistsExcludeID(ctx context.Context, name string, id uint) (bool, error)
	FindByID(ctx context.Context, id uint) (*domain.BillType, error)
	CountBillingRules(ctx context.Context, billTypeID uint) (int, error)
	CountStudentBills(ctx context.Context, billTypeID uint) (int, error)
}

type billTypeRepo struct {
	db *bun.DB
}

func NewBillTypeRepo(db *bun.DB) BillTypeRepo {
	return &billTypeRepo{db: db}
}

func (r *billTypeRepo) Create(ctx context.Context, b *domain.BillType) error {
	_, err := r.db.NewInsert().Model(b).Exec(ctx)
	return err
}

func (r *billTypeRepo) FindAll(ctx context.Context) ([]domain.BillType, error) {
	var list []domain.BillType
	err := r.db.NewSelect().Model(&list).Order("name ASC").Scan(ctx)
	return list, err
}

func (r *billTypeRepo) FindAllPaged(ctx context.Context, page, limit int, search, filterType, status, sort string) ([]domain.BillType, int, error) {
	var list []domain.BillType
	q := r.db.NewSelect().Model(&list).
		ColumnExpr("bt.*").
		ColumnExpr("(SELECT COUNT(*) FROM billing_rules WHERE bill_type_id = bt.id AND deleted_at IS NULL) AS rule_count")

	if status == "trash" {
		q.WhereAllWithDeleted().Where("bt.deleted_at IS NOT NULL")
	} else {
		q.Where("bt.deleted_at IS NULL")
		switch status {
		case "active":
			q.Where("bt.is_active = ?", true)
		case "inactive":
			q.Where("bt.is_active = ?", false)
		}
	}

	if filterType != "" {
		q.Where("bt.type = ?", filterType)
	}

	if search != "" {
		s := "%" + search + "%"
		q.Where("bt.name LIKE ? OR bt.description LIKE ?", s, s)
	}

	switch sort {
	case "created_asc":
		q.Order("bt.created_at ASC")
	case "name_asc":
		q.Order("bt.name ASC")
	case "name_desc":
		q.Order("bt.name DESC")
	default:
		q.Order("bt.created_at DESC")
	}

	total, err := q.Limit(limit).
		Offset((page - 1) * limit).
		ScanAndCount(ctx)

	return list, total, err
}

func (r *billTypeRepo) Update(ctx context.Context, b *domain.BillType) error {
	_, err := r.db.NewUpdate().Model(b).WherePK().Exec(ctx)
	return err
}

func (r *billTypeRepo) Delete(ctx context.Context, id uint) error {
	res, err := r.db.NewDelete().Model((*domain.BillType)(nil)).Where("id = ? AND deleted_at IS NULL", id).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("jenis tagihan tidak ditemukan atau sudah terhapus")
	}
	return nil
}

func (r *billTypeRepo) Restore(ctx context.Context, id uint) error {
	res, err := r.db.NewUpdate().
		Model((*domain.BillType)(nil)).
		WhereAllWithDeleted().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Set("deleted_at = NULL").
		Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("jenis tagihan tidak ditemukan di riwayat penghapusan")
	}
	return nil
}

func (r *billTypeRepo) ToggleStatus(ctx context.Context, id uint) error {
	_, err := r.db.NewUpdate().
		Model((*domain.BillType)(nil)).
		Set("is_active = NOT is_active").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *billTypeRepo) BulkDelete(ctx context.Context, ids []uint) error {
	res, err := r.db.NewDelete().Model((*domain.BillType)(nil)).Where("id IN (?) AND deleted_at IS NULL", bun.In(ids)).Exec(ctx)
	if err != nil {
		return err
	}
	rows, _ := res.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("tidak ada data yang dapat dihapus")
	}
	return nil
}

func (r *billTypeRepo) BulkRestore(ctx context.Context, ids []uint) error {
	res, err := r.db.NewUpdate().
		Model((*domain.BillType)(nil)).
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

func (r *billTypeRepo) Exists(ctx context.Context, name string) (bool, error) {
	return r.db.NewSelect().Model((*domain.BillType)(nil)).
		WhereAllWithDeleted().
		Where("name = ?", name).
		Exists(ctx)
}

func (r *billTypeRepo) ExistsExcludeID(ctx context.Context, name string, id uint) (bool, error) {
	return r.db.NewSelect().Model((*domain.BillType)(nil)).
		WhereAllWithDeleted().
		Where("name = ? AND id != ?", name, id).
		Exists(ctx)
}

func (r *billTypeRepo) FindByID(ctx context.Context, id uint) (*domain.BillType, error) {
	b := new(domain.BillType)
	err := r.db.NewSelect().Model(b).Where("id = ?", id).Scan(ctx)
	return b, err
}

func (r *billTypeRepo) CountBillingRules(ctx context.Context, billTypeID uint) (int, error) {
	return r.db.NewSelect().Model((*domain.BillingRule)(nil)).Where("bill_type_id = ?", billTypeID).Count(ctx)
}

func (r *billTypeRepo) CountStudentBills(ctx context.Context, billTypeID uint) (int, error) {
	return r.db.NewSelect().Model((*domain.StudentBill)(nil)).Where("bill_type_id = ? AND status != ?", billTypeID, "voided").Count(ctx)
}
