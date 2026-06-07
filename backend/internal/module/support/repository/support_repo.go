package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/domain"
	"github.com/uptrace/bun"
)

type SupportRepo interface {
	FindOpenByPhone(ctx context.Context, phone string) (*domain.Conversation, error)
	FindOpenByParentID(ctx context.Context, parentID uint) (*domain.Conversation, error)
	CreateConversation(ctx context.Context, db bun.IDB, c *domain.Conversation) error
	FindAll(ctx context.Context, status string, page, limit int) ([]domain.Conversation, int, error)
	FindByID(ctx context.Context, id uint) (*domain.Conversation, error)
	Assign(ctx context.Context, db bun.IDB, conversationID, adminID uint) error
	Close(ctx context.Context, db bun.IDB, conversationID uint) error
	UpdateStatus(ctx context.Context, db bun.IDB, conversationID uint, status string) error
}

type supportRepo struct{ db *bun.DB }

func NewSupportRepo(db *bun.DB) SupportRepo { return &supportRepo{db: db} }

func (r *supportRepo) FindOpenByPhone(ctx context.Context, phone string) (*domain.Conversation, error) {
	var c domain.Conversation
	err := r.db.NewSelect().Model(&c).Where("phone_number = ? AND status IN ('open', 'pending')", phone).Order("updated_at DESC").Limit(1).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *supportRepo) FindOpenByParentID(ctx context.Context, parentID uint) (*domain.Conversation, error) {
	var c domain.Conversation
	err := r.db.NewSelect().
		Model(&c).
		Where("parent_id = ? AND status IN ('open', 'pending')", parentID).
		Order("updated_at DESC").
		Limit(1).
		Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &c, nil
}

func (r *supportRepo) CreateConversation(ctx context.Context, db bun.IDB, c *domain.Conversation) error {
	_, err := db.NewInsert().Model(c).Exec(ctx)
	return err
}

func (r *supportRepo) FindAll(ctx context.Context, status string, page, limit int) ([]domain.Conversation, int, error) {
	var list []domain.Conversation
	q := r.db.NewSelect().Model(&list)
	if status != "" {
		q.Where("status = ?", status)
	}
	total, err := q.OrderExpr("created_at DESC").Limit(limit).Offset((page - 1) * limit).ScanAndCount(ctx)
	return list, total, err
}

func (r *supportRepo) FindByID(ctx context.Context, id uint) (*domain.Conversation, error) {
	var c domain.Conversation
	err := r.db.NewSelect().Model(&c).Where("id = ?", id).Scan(ctx)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *supportRepo) Assign(ctx context.Context, db bun.IDB, conversationID, adminID uint) error {
	result, err := db.NewUpdate().
		Model((*domain.Conversation)(nil)).
		Set("status = 'pending'").
		Set("updated_at = CURRENT_TIMESTAMP").
		Where("id = ?", conversationID).
		Exec(ctx)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("percakapan CS tidak ditemukan")
	}
	return nil
}

func (r *supportRepo) Close(ctx context.Context, db bun.IDB, conversationID uint) error {
	return r.UpdateStatus(ctx, db, conversationID, "closed")
}

func (r *supportRepo) UpdateStatus(ctx context.Context, db bun.IDB, conversationID uint, status string) error {
	result, err := db.NewUpdate().
		Model((*domain.Conversation)(nil)).
		Set("status = ?", status).
		Set("updated_at = CURRENT_TIMESTAMP").
		Where("id = ?", conversationID).
		Exec(ctx)
	if err != nil {
		return err
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return fmt.Errorf("percakapan CS tidak ditemukan")
	}
	return nil
}
