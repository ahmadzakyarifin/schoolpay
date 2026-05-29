package repository

import (
	"context"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/domain"
	"github.com/uptrace/bun"
)

type SupportRepo interface {
	FindOpenByPhone(ctx context.Context, phone string) (*domain.Conversation, error)
	CreateConversation(ctx context.Context, db bun.IDB, c *domain.Conversation) error
	UpdateConversationPreview(ctx context.Context, db bun.IDB, id uint, lastMessage string, unreadDelta int) error
	CreateMessage(ctx context.Context, db bun.IDB, m *domain.Message) error
	FindAll(ctx context.Context, status string, page, limit int) ([]domain.Conversation, int, error)
	FindMessages(ctx context.Context, conversationID uint) ([]domain.Message, error)
	FindByID(ctx context.Context, id uint) (*domain.Conversation, error)
	Assign(ctx context.Context, db bun.IDB, conversationID, adminID uint) error
	MarkRead(ctx context.Context, db bun.IDB, conversationID uint) error
	Close(ctx context.Context, db bun.IDB, conversationID uint) error
}

type supportRepo struct{ db *bun.DB }

func NewSupportRepo(db *bun.DB) SupportRepo { return &supportRepo{db: db} }

func (r *supportRepo) FindOpenByPhone(ctx context.Context, phone string) (*domain.Conversation, error) {
	var c domain.Conversation
	err := r.db.NewSelect().Model(&c).Where("phone_number = ? AND status IN ('open', 'pending')", phone).Order("updated_at DESC").Limit(1).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *supportRepo) CreateConversation(ctx context.Context, db bun.IDB, c *domain.Conversation) error {
	_, err := db.NewInsert().Model(c).Exec(ctx)
	return err
}

func (r *supportRepo) UpdateConversationPreview(ctx context.Context, db bun.IDB, id uint, lastMessage string, unreadDelta int) error {
	now := time.Now()
	q := db.NewUpdate().Model((*domain.Conversation)(nil)).
		Set("last_message = ?", lastMessage).
		Set("last_message_at = ?", now).
		Set("updated_at = ?", now).
		Where("id = ?", id)
	if unreadDelta != 0 {
		q.Set("unread_count = unread_count + ?", unreadDelta)
	}
	_, err := q.Exec(ctx)
	return err
}

func (r *supportRepo) CreateMessage(ctx context.Context, db bun.IDB, m *domain.Message) error {
	_, err := db.NewInsert().Model(m).Exec(ctx)
	return err
}

func (r *supportRepo) FindAll(ctx context.Context, status string, page, limit int) ([]domain.Conversation, int, error) {
	var list []domain.Conversation
	q := r.db.NewSelect().Model(&list)
	if status != "" {
		q.Where("status = ?", status)
	}
	total, err := q.OrderExpr("COALESCE(last_message_at, created_at) DESC").Limit(limit).Offset((page - 1) * limit).ScanAndCount(ctx)
	return list, total, err
}

func (r *supportRepo) FindMessages(ctx context.Context, conversationID uint) ([]domain.Message, error) {
	var list []domain.Message
	err := r.db.NewSelect().Model(&list).Where("conversation_id = ?", conversationID).Order("created_at ASC").Scan(ctx)
	return list, err
}

func (r *supportRepo) FindByID(ctx context.Context, id uint) (*domain.Conversation, error) {
	var c domain.Conversation
	err := r.db.NewSelect().Model(&c).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *supportRepo) Assign(ctx context.Context, db bun.IDB, conversationID, adminID uint) error {
	_, err := db.NewUpdate().Model((*domain.Conversation)(nil)).Set("assigned_admin_id = ?", adminID).Set("status = 'pending'").Where("id = ?", conversationID).Exec(ctx)
	return err
}

func (r *supportRepo) MarkRead(ctx context.Context, db bun.IDB, conversationID uint) error {
	_, err := db.NewUpdate().Model((*domain.Conversation)(nil)).Set("unread_count = 0").Where("id = ?", conversationID).Exec(ctx)
	return err
}

func (r *supportRepo) Close(ctx context.Context, db bun.IDB, conversationID uint) error {
	now := time.Now()
	_, err := db.NewUpdate().Model((*domain.Conversation)(nil)).Set("status = 'closed'").Set("closed_at = ?", now).Where("id = ?", conversationID).Exec(ctx)
	return err
}
