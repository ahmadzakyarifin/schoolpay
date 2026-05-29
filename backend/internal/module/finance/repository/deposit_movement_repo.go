package repository

import (
	"context"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/uptrace/bun"
)

type DepositMovementRepo interface {
	Create(ctx context.Context, db bun.IDB, dm *domain.DepositMovement) error
	FindByStudent(ctx context.Context, studentID uint) ([]domain.DepositMovement, error)
}

type depositMovementRepo struct {
	db *bun.DB
}

func NewDepositMovementRepo(db *bun.DB) DepositMovementRepo {
	return &depositMovementRepo{db: db}
}

func (r *depositMovementRepo) Create(ctx context.Context, db bun.IDB, dm *domain.DepositMovement) error {
	_, err := db.NewInsert().Model(dm).Exec(ctx)
	return err
}

func (r *depositMovementRepo) FindByStudent(ctx context.Context, studentID uint) ([]domain.DepositMovement, error) {
	var list []domain.DepositMovement
	err := r.db.NewSelect().
		Model(&list).
		ColumnExpr("dm.*").
		ColumnExpr("s.name as student_name").
		Join("JOIN students s ON dm.student_id = s.id").
		Where("dm.student_id = ?", studentID).
		Order("dm.created_at DESC").
		Scan(ctx)
	return list, err
}
