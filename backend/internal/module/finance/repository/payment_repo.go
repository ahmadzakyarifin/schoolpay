package repository

import (
	"context"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

type PaymentRepo interface {
	Create(ctx context.Context, db bun.IDB, p *domain.Payment) error
	CreateDetail(ctx context.Context, db bun.IDB, pd *domain.PaymentDetail) error
	FindByRef(ctx context.Context, ref string) (*domain.Payment, error)
	FindByID(ctx context.Context, id uint) (*domain.Payment, error)
	FindDetailsByPaymentID(ctx context.Context, paymentID uint) ([]domain.PaymentDetail, error)
	FindByStudent(ctx context.Context, studentID uint) ([]domain.Payment, error)
	DeleteDetailsByPaymentID(ctx context.Context, db bun.IDB, paymentID uint) error
}

type paymentRepo struct {
	db *bun.DB
}

func NewPaymentRepo(db *bun.DB) PaymentRepo {
	return &paymentRepo{db: db}
}

func (r *paymentRepo) Create(ctx context.Context, db bun.IDB, p *domain.Payment) error {
	_, err := db.NewInsert().Model(p).Exec(ctx)
	return err
}

func (r *paymentRepo) CreateDetail(ctx context.Context, db bun.IDB, pd *domain.PaymentDetail) error {
	_, err := db.NewInsert().Model(pd).Exec(ctx)
	return err
}

func (r *paymentRepo) FindByRef(ctx context.Context, ref string) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.NewSelect().Model(&p).Where("transaction_ref = ?", ref).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepo) FindByID(ctx context.Context, id uint) (*domain.Payment, error) {
	var p domain.Payment
	err := r.db.NewSelect().Model(&p).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *paymentRepo) FindDetailsByPaymentID(ctx context.Context, paymentID uint) ([]domain.PaymentDetail, error) {
	var list []domain.PaymentDetail
	err := r.db.NewSelect().
		Model(&list).
		ColumnExpr("pd.*").
		ColumnExpr("COALESCE(NULLIF(sb.name, ''), bt.name) as bill_type_name, COALESCE(sb.period, '') as period").
		Join("JOIN student_bills sb ON pd.student_bill_id = sb.id").
		Join("JOIN bill_types bt ON sb.bill_type_id = bt.id").
		Where("pd.payment_id = ?", paymentID).
		Scan(ctx)
	return list, err
}

func (r *paymentRepo) FindByStudent(ctx context.Context, studentID uint) ([]domain.Payment, error) {
	var list []domain.Payment
	err := r.db.NewSelect().
		Model(&list).
		Where("student_id = ? AND status = 'success'", studentID).
		Order("paid_at DESC").
		Scan(ctx)
	if err != nil {
		return list, err
	}
	for i := range list {
		details, detailErr := r.FindDetailsByPaymentID(ctx, list[i].ID)
		if detailErr != nil {
			return list, detailErr
		}
		list[i].Details = details
	}
	return list, err
}

func (r *paymentRepo) DeleteDetailsByPaymentID(ctx context.Context, db bun.IDB, paymentID uint) error {
	_, err := db.NewDelete().Model((*domain.PaymentDetail)(nil)).Where("payment_id = ?", paymentID).Exec(ctx)
	return err
}
