package repository

import (
	"context"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/webhook/domain"
	"github.com/uptrace/bun"
)

type WebhookRepo interface {
	Create(ctx context.Context, log *domain.WebhookLog) error
	UpdateStatus(ctx context.Context, eventID string, status string) error
}

type webhookRepo struct {
	db *bun.DB
}

func NewWebhookRepo(db *bun.DB) WebhookRepo {
	return &webhookRepo{db: db}
}

func (r *webhookRepo) Create(ctx context.Context, log *domain.WebhookLog) error {
	_, err := r.db.NewInsert().Model(log).Exec(ctx)
	return err
}

func (r *webhookRepo) UpdateStatus(ctx context.Context, eventID string, status string) error {
	_, err := r.db.NewUpdate().
		Model((*domain.WebhookLog)(nil)).
		Set("status = ?", status).
		Where("event_id = ?", eventID).
		Exec(ctx)
	return err
}
