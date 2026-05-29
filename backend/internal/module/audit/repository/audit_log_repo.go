package repository

import (
	"context"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/audit/domain"
	"github.com/uptrace/bun"
)

type AuditLogRepo interface {
	Log(ctx context.Context, db bun.IDB, al *domain.AuditLog) error
	FindAll(ctx context.Context, filter map[string]interface{}, page, limit int) ([]domain.AuditLog, int, error)
	FindByEntity(ctx context.Context, entityType string, entityID uint) ([]domain.AuditLog, error)
}

type auditLogRepo struct {
	db *bun.DB
}

func NewAuditLogRepo(db *bun.DB) AuditLogRepo {
	return &auditLogRepo{db: db}
}

func (r *auditLogRepo) Log(ctx context.Context, db bun.IDB, al *domain.AuditLog) error {
	if db == nil {
		db = r.db
	}
	_, err := db.NewInsert().Model(al).Exec(ctx)
	return err
}

func (r *auditLogRepo) FindAll(ctx context.Context, filter map[string]interface{}, page, limit int) ([]domain.AuditLog, int, error) {
	var logs []domain.AuditLog

	q := r.db.NewSelect().Model(&logs)

	if action, ok := filter["action"].(string); ok && action != "" {
		q.Where("action = ?", action)
	}
	if entityType, ok := filter["entity_type"].(string); ok && entityType != "" {
		q.Where("entity_type = ?", entityType)
	}
	if role, ok := filter["role"].(string); ok && role != "" {
		q.Where("role = ?", role)
	}
	if userName, ok := filter["user_name"].(string); ok && userName != "" {
		search := "%" + userName + "%"
		q.WhereGroup(" AND ", func(sq *bun.SelectQuery) *bun.SelectQuery {
			return sq.Where("user_name LIKE ?", search).
				WhereOr("action LIKE ?", search).
				WhereOr("entity_type LIKE ?", search)
		})
	}

	total, err := q.Count(ctx)
	if err != nil {
		return nil, 0, err
	}

	sort, _ := filter["sort"].(string)
	switch sort {
	case "created_asc":
		q.Order("created_at ASC", "id ASC")
	case "action_asc":
		q.Order("action ASC", "created_at DESC")
	case "entity_asc":
		q.Order("entity_type ASC", "created_at DESC")
	case "user_asc":
		q.Order("user_name ASC", "created_at DESC")
	default:
		q.Order("created_at DESC", "id DESC")
	}

	err = q.
		Limit(limit).
		Offset((page - 1) * limit).
		Scan(ctx)

	return logs, total, err
}

func (r *auditLogRepo) FindByEntity(ctx context.Context, entityType string, entityID uint) ([]domain.AuditLog, error) {
	var logs []domain.AuditLog
	err := r.db.NewSelect().
		Model(&logs).
		Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Order("created_at DESC").
		Scan(ctx)
	return logs, err
}
