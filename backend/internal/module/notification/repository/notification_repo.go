package repository

import (
	"context"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/notification/domain"
	"github.com/uptrace/bun"
)

type NotificationRepo interface {
	Create(ctx context.Context, db bun.IDB, n *domain.Notification) error
	GetByUserID(ctx context.Context, userID uint) ([]domain.Notification, error)
	MarkAsRead(ctx context.Context, id uint) error
	UpdateStatusByWhatsappID(ctx context.Context, whatsappID string, status string) error
	GetStats(ctx context.Context) (map[string]int, error)
	GetEfficacyStats(ctx context.Context, channel string, start, end *time.Time) (map[string]int, error)
	GetDetailedLogs(ctx context.Context, page, limit int, status string, search string) ([]map[string]interface{}, int, error)
	FindByID(ctx context.Context, id uint) (*domain.Notification, error)
}

type notificationRepo struct {
	db *bun.DB
}

func NewNotificationRepo(db *bun.DB) NotificationRepo {
	return &notificationRepo{db: db}
}

func (r *notificationRepo) Create(ctx context.Context, db bun.IDB, n *domain.Notification) error {
	_, err := db.NewInsert().Model(n).Exec(ctx)
	return err
}

func (r *notificationRepo) FindByID(ctx context.Context, id uint) (*domain.Notification, error) {
	var n domain.Notification
	err := r.db.NewSelect().Model(&n).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &n, nil
}

func (r *notificationRepo) UpdateStatusByWhatsappID(ctx context.Context, whatsappID string, status string) error {
	_, err := r.db.NewUpdate().
		Model((*domain.Notification)(nil)).
		Set("delivery_status = ?", status).
		Where("whatsapp_id = ?", whatsappID).
		Exec(ctx)
	return err
}

func (r *notificationRepo) GetStats(ctx context.Context) (map[string]int, error) {
	var stats []struct {
		Status string `bun:"delivery_status"`
		Count  int    `bun:"count"`
	}

	err := r.db.NewSelect().
		Model((*domain.Notification)(nil)).
		ColumnExpr("delivery_status, COUNT(*) as count").
		Group("delivery_status").
		Scan(ctx, &stats)

	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, status := range []string{"PENDING", "SENT", "DELIVERED", "READ", "FAILED"} {
		result[status] = 0
	}
	for _, s := range stats {
		result[normalizeDeliveryStatus(s.Status)] += s.Count
	}
	return result, nil
}

func (r *notificationRepo) GetEfficacyStats(ctx context.Context, channel string, start, end *time.Time) (map[string]int, error) {
	var stats []struct {
		Status string `bun:"delivery_status"`
		Count  int    `bun:"count"`
	}

	q := r.db.NewSelect().
		Model((*domain.Notification)(nil)).
		ColumnExpr("delivery_status, COUNT(*) as count")

	if start != nil {
		q.Where("created_at >= ?", start)
	}
	if end != nil {
		q.Where("created_at <= ?", end)
	}

	if channel == "whatsapp" {
		q.Where("whatsapp_id IS NOT NULL AND whatsapp_id != ''")
	} else {
		q.Where("whatsapp_id IS NULL OR whatsapp_id = ''")
	}

	err := q.Group("delivery_status").Scan(ctx, &stats)
	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, status := range []string{"pending", "sent", "delivered", "read", "failed"} {
		result[status] = 0
	}
	for _, s := range stats {
		result[strings.ToLower(normalizeDeliveryStatus(s.Status))] += s.Count
	}
	return result, nil
}

func normalizeDeliveryStatus(status string) string {
	switch strings.ToUpper(strings.TrimSpace(status)) {
	case "SUCCESS":
		return "SENT"
	case "DELIVERED":
		return "DELIVERED"
	case "READ":
		return "READ"
	case "FAILED", "ERROR":
		return "FAILED"
	case "PENDING", "":
		return "PENDING"
	default:
		return "SENT"
	}
}

func (r *notificationRepo) GetByUserID(ctx context.Context, userID uint) ([]domain.Notification, error) {
	var ns []domain.Notification
	err := r.db.NewSelect().
		Model(&ns).
		Where("user_id = ?", userID).
		Order("created_at DESC").
		Scan(ctx)
	return ns, err
}

func (r *notificationRepo) MarkAsRead(ctx context.Context, id uint) error {
	_, err := r.db.NewUpdate().
		Model((*domain.Notification)(nil)).
		Set("is_read = TRUE").
		Where("id = ?", id).
		Exec(ctx)
	return err
}

func (r *notificationRepo) GetDetailedLogs(ctx context.Context, page, limit int, status string, search string) ([]map[string]interface{}, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	status = strings.TrimSpace(strings.ToLower(status))
	search = strings.TrimSpace(search)
	base := r.db.NewSelect().
		Model((*domain.Notification)(nil)).
		ColumnExpr("n.id, n.user_id, n.title, n.message, n.type, n.is_read, n.whatsapp_id, n.delivery_status, n.delivery_error, n.created_at, n.updated_at").
		ColumnExpr("COALESCE(u.name, '') as recipient_name").
		ColumnExpr("COALESCE(u.phone_number, '') as recipient_phone").
		ColumnExpr("COALESCE(u.email, '') as recipient_email").
		ColumnExpr("CASE WHEN n.whatsapp_id IS NULL OR n.whatsapp_id = '' THEN 'email' ELSE 'whatsapp' END as channel").
		Join("LEFT JOIN users u ON n.user_id = u.id")

	if status != "" {
		base.Where("LOWER(n.delivery_status) = ?", status)
	}
	if search != "" {
		like := "%" + search + "%"
		base.Where("u.name LIKE ? OR u.phone_number LIKE ? OR u.email LIKE ? OR n.title LIKE ? OR n.message LIKE ?", like, like, like, like, like)
	}

	total, err := base.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	var results []map[string]interface{}
	err = base.Order("n.created_at DESC").
		Limit(limit).
		Offset((page-1)*limit).
		Scan(ctx, &results)

	return results, total, err
}
