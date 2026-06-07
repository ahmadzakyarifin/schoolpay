package usecase

import (
	"context"

	"github.com/ahmadzakyarifin/schoolpay/internal/dto"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/audit/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/audit/repository"
	userdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/uptrace/bun"
)

type AuditLogService interface {
	Log(ctx context.Context, db bun.IDB, userID uint, userName, role, action, entityType string, entityID uint, oldValues, newValues map[string]interface{}, ipAddress, userAgent string) error
	LogMeta(ctx context.Context, db bun.IDB, audit dto.AuditMeta, action, entityType string, entityID uint, oldValues, newValues map[string]interface{}, description string) error
	GetLogs(ctx context.Context, currentUser *userdomain.User, filter map[string]interface{}, page, limit int) ([]domain.AuditLog, int, error)
	GetEntityLogs(ctx context.Context, currentUser *userdomain.User, entityType string, entityID uint) ([]domain.AuditLog, error)
}

type auditLogService struct {
	repo repository.AuditLogRepo
}

func NewAuditLogService(repo repository.AuditLogRepo) AuditLogService {
	return &auditLogService{repo: repo}
}

func (s *auditLogService) Log(ctx context.Context, db bun.IDB, userID uint, userName, role, action, entityType string, entityID uint, oldValues, newValues map[string]interface{}, ipAddress, userAgent string) error {
	var uidPtr *uint
	if userID != 0 {
		uidPtr = &userID
	}
	al := &domain.AuditLog{
		UserID:     uidPtr,
		UserName:   userName,
		Role:       role,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		OldValues:  oldValues,
		NewValues:  newValues,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
	}
	return s.repo.Log(ctx, db, al)
}

func (s *auditLogService) LogMeta(ctx context.Context, db bun.IDB, audit dto.AuditMeta, action, entityType string, entityID uint, oldValues, newValues map[string]interface{}, description string) error {
	var userName string
	if v := ctx.Value("user_name"); v != nil {
		userName, _ = v.(string)
	}
	if userName == "" {
		userName = "System/Admin"
	}
	var role string
	if v := ctx.Value("role"); v != nil {
		role, _ = v.(string)
	}
	if role == "" {
		role = "admin"
	}

	al := &domain.AuditLog{
		UserID:      audit.UserID,
		UserName:    userName,
		Role:        role,
		Action:      action,
		EntityType:  entityType,
		EntityID:    entityID,
		OldValues:   oldValues,
		NewValues:   newValues,
		IPAddress:   audit.IP,
		UserAgent:   audit.UserAgent,
		Method:      audit.Method,
		Path:        audit.Path,
		DeviceID:    audit.DeviceID,
		AppPlatform: audit.AppPlatform,
		AppVersion:  audit.AppVersion,
		Description: description,
	}
	return s.repo.Log(ctx, db, al)
}

func (s *auditLogService) GetLogs(ctx context.Context, currentUser *userdomain.User, filter map[string]interface{}, page, limit int) ([]domain.AuditLog, int, error) {
	if filter == nil {
		filter = make(map[string]interface{})
	}

	// RBAC Grouping Strategy:
	// Admin Utama / Kepala Sekolah bisa melihat semua role.
	// Role spesifik (Bendahara, TU) hanya bisa melihat log dari role mereka sendiri.
	if currentUser != nil {
		if currentUser.Role != "admin" && currentUser.Role != "admin_utama" && currentUser.Role != "kepala_sekolah" {
			filter["role"] = currentUser.Role
		}
	}

	return s.repo.FindAll(ctx, filter, page, limit)
}

func (s *auditLogService) GetEntityLogs(ctx context.Context, currentUser *userdomain.User, entityType string, entityID uint) ([]domain.AuditLog, error) {
	return s.repo.FindByEntity(ctx, entityType, entityID)
}
