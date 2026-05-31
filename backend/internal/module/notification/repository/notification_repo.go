package repository

import (
	"context"
	"fmt"
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
	GetDetailedLogs(ctx context.Context, page, limit int, status string, search string, channel string) ([]map[string]interface{}, int, error)
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
	status = normalizeDeliveryStatus(status)
	rank := deliveryStatusRank(status)

	_, err := r.db.NewUpdate().
		Model((*domain.Notification)(nil)).
		Set("delivery_status = ?", status).
		Where("whatsapp_id = ?", whatsappID).
		Where("(? = 'FAILED' AND UPPER(COALESCE(delivery_status, '')) NOT IN ('READ', 'DELIVERED')) OR (? != 'FAILED' AND CASE UPPER(COALESCE(delivery_status, '')) WHEN 'READ' THEN 4 WHEN 'DELIVERED' THEN 3 WHEN 'SENT' THEN 2 WHEN 'SUCCESS' THEN 2 WHEN 'PENDING' THEN 1 WHEN '' THEN 0 ELSE 0 END <= ?)", status, status, rank).
		Exec(ctx)
	return err
}

func (r *notificationRepo) GetStats(ctx context.Context) (map[string]int, error) {
	var stats []struct {
		Status  string `bun:"delivery_status"`
		Channel string `bun:"channel"`
		Count   int    `bun:"count"`
	}

	err := r.db.NewSelect().
		Model((*domain.Notification)(nil)).
		ColumnExpr("delivery_status").
		ColumnExpr("LOWER(COALESCE(NULLIF(channel, ''), CASE WHEN whatsapp_id IS NOT NULL AND whatsapp_id != '' THEN 'whatsapp' ELSE 'email' END)) as channel").
		ColumnExpr("COUNT(*) as count").
		Group("delivery_status", "channel").
		Scan(ctx, &stats)

	if err != nil {
		return nil, err
	}

	result := make(map[string]int)
	for _, status := range []string{"PENDING", "SENT", "DELIVERED", "READ", "FAILED"} {
		result[status] = 0
	}
	for _, s := range stats {
		normalized := normalizeDeliveryStatus(s.Status)
		if s.Channel != "whatsapp" && (normalized == "DELIVERED" || normalized == "READ") {
			normalized = "SENT"
		}
		result[normalized] += s.Count
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
		q.Where("LOWER(COALESCE(NULLIF(channel, ''), CASE WHEN whatsapp_id IS NOT NULL AND whatsapp_id != '' THEN 'whatsapp' ELSE 'email' END)) = 'whatsapp'")
	} else {
		q.Where("LOWER(COALESCE(NULLIF(channel, ''), CASE WHEN whatsapp_id IS NOT NULL AND whatsapp_id != '' THEN 'whatsapp' ELSE 'email' END)) = 'email'")
	}

	err := q.Group("delivery_status").Scan(ctx, &stats)
	if err != nil {
		return nil, err
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
		normalized := strings.ToLower(normalizeDeliveryStatus(s.Status))
		if channel != "whatsapp" && (normalized == "delivered" || normalized == "read") {
			normalized = "sent"
		}
		if _, ok := result[normalized]; ok {
			result[normalized] += s.Count
		}
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

func deliveryStatusRank(status string) int {
	switch normalizeDeliveryStatus(status) {
	case "READ":
		return 4
	case "DELIVERED":
		return 3
	case "SENT":
		return 2
	case "PENDING":
		return 1
	default:
		return 0
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

func (r *notificationRepo) GetDetailedLogs(ctx context.Context, page, limit int, status string, search string, channel string) ([]map[string]interface{}, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	status = strings.TrimSpace(strings.ToLower(status))
	search = strings.TrimSpace(search)
	channel = strings.TrimSpace(strings.ToLower(channel))
	statusExpr := "LOWER(CASE WHEN n.delivery_status IN ('SUCCESS') THEN 'sent' WHEN n.delivery_status IS NULL OR n.delivery_status = '' THEN 'pending' ELSE n.delivery_status END)"
	base := r.db.NewSelect().
		Model((*domain.Notification)(nil)).
		ColumnExpr("n.id, n.user_id, n.title, n.message, n.type, n.channel, n.is_read, n.whatsapp_id, n.delivery_error, n.created_at, n.updated_at").
		ColumnExpr("CASE WHEN n.delivery_status IN ('SUCCESS') THEN 'sent' WHEN n.delivery_status IS NULL OR n.delivery_status = '' THEN 'pending' ELSE LOWER(n.delivery_status) END as delivery_status").
		ColumnExpr("COALESCE(u.name, '') as recipient_name").
		ColumnExpr("COALESCE(u.phone_number, '') as recipient_phone").
		ColumnExpr("COALESCE(u.email, '') as recipient_email").
		ColumnExpr("COALESCE(NULLIF(n.channel, ''), CASE WHEN n.whatsapp_id IS NULL OR n.whatsapp_id = '' THEN 'email' ELSE 'whatsapp' END) as channel").
		Join("LEFT JOIN users u ON n.user_id = u.id")

	if status != "" {
		if channel == "email" && status == "sent" {
			base.Where(statusExpr+" IN (?)", bun.In([]string{"sent", "delivered", "read"}))
		} else {
			base.Where(statusExpr+" = ?", status)
		}
	}
	if channel == "whatsapp" {
		base.Where("LOWER(COALESCE(NULLIF(n.channel, ''), CASE WHEN n.whatsapp_id IS NOT NULL AND n.whatsapp_id != '' THEN 'whatsapp' ELSE 'email' END)) = 'whatsapp'")
	} else if channel == "email" {
		base.Where("LOWER(COALESCE(NULLIF(n.channel, ''), CASE WHEN n.whatsapp_id IS NOT NULL AND n.whatsapp_id != '' THEN 'whatsapp' ELSE 'email' END)) = 'email'")
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

	for _, row := range results {
		if strings.EqualFold(fmt.Sprint(row["channel"]), "email") {
			status := strings.ToLower(fmt.Sprint(row["delivery_status"]))
			if status == "delivered" || status == "read" {
				row["delivery_status"] = "sent"
			}
		}
	}

	return results, total, err
}
