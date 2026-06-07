package usecase

import (
	"context"
	"fmt"
	"strings"

	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/uptrace/bun"
)

type BillTypeService interface {
	Create(ctx context.Context, b *domain.BillType) error
	GetAll(ctx context.Context) ([]domain.BillType, error)
	GetAllPaged(ctx context.Context, page, limit int, search, filterType, status, sort string) ([]domain.BillType, int, error)
	Update(ctx context.Context, b *domain.BillType) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetByID(ctx context.Context, id uint) (*domain.BillType, error)
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, name string, excludeID uint) (bool, error)
}

type billTypeService struct {
	db    bun.IDB
	repo  repository.BillTypeRepo
	audit auditusecase.AuditLogService
}

func NewBillTypeService(db bun.IDB, repo repository.BillTypeRepo, audit auditusecase.AuditLogService) BillTypeService {
	return &billTypeService{db: db, repo: repo, audit: audit}
}

func (s *billTypeService) Create(ctx context.Context, b *domain.BillType) error {
	b.Name = strings.TrimSpace(b.Name)
	if b.Name == "" {
		return fmt.Errorf("gagal: Nama jenis tagihan wajib diisi")
	}
	if b.DefaultAmount <= 0 {
		return fmt.Errorf("gagal: Nominal dasar wajib diisi dengan angka lebih dari 0")
	}

	exists, err := s.repo.Exists(ctx, b.Name)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("gagal: Jenis tagihan '%s' sudah ada dalam sistem", b.Name)
	}

	err = s.repo.Create(ctx, b)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		newVals := map[string]interface{}{"name": b.Name, "description": b.Description, "type": b.Type, "default_amount": b.DefaultAmount}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "CREATE", "bill_types", b.ID, nil, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *billTypeService) GetAll(ctx context.Context) ([]domain.BillType, error) {
	return s.repo.FindAll(ctx)
}

func (s *billTypeService) GetAllPaged(ctx context.Context, page, limit int, search, filterType, status, sort string) ([]domain.BillType, int, error) {
	return s.repo.FindAllPaged(ctx, page, limit, search, filterType, status, sort)
}

func (s *billTypeService) GetByID(ctx context.Context, id uint) (*domain.BillType, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *billTypeService) Update(ctx context.Context, b *domain.BillType) error {
	b.Name = strings.TrimSpace(b.Name)
	if b.Name == "" {
		return fmt.Errorf("gagal: Nama jenis tagihan wajib diisi")
	}
	if b.DefaultAmount <= 0 {
		return fmt.Errorf("gagal: Nominal dasar wajib diisi dengan angka lebih dari 0")
	}

	exists, err := s.repo.ExistsExcludeID(ctx, b.Name, b.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("gagal: Nama jenis tagihan '%s' sudah digunakan oleh tipe lain", b.Name)
	}

	existing, _ := s.repo.FindByID(ctx, b.ID)
	err = s.repo.Update(ctx, b)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"name": existing.Name, "description": existing.Description, "type": existing.Type, "default_amount": existing.DefaultAmount}
		newVals := map[string]interface{}{"name": b.Name, "description": b.Description, "type": b.Type, "default_amount": b.DefaultAmount}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "UPDATE", "bill_types", b.ID, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *billTypeService) Delete(ctx context.Context, id uint) error {
	rulesCount, err := s.repo.CountBillingRules(ctx, id)
	if err != nil {
		return err
	}
	if rulesCount > 0 {
		return fmt.Errorf("gagal: Jenis tagihan tidak dapat dihapus karena masih digunakan oleh %d aturan tagihan", rulesCount)
	}

	billsCount, err := s.repo.CountStudentBills(ctx, id)
	if err != nil {
		return err
	}
	if billsCount > 0 {
		return fmt.Errorf("gagal: Jenis tagihan tidak dapat dihapus karena masih memiliki %d tagihan siswa yang belum ditarik", billsCount)
	}

	existing, _ := s.repo.FindByID(ctx, id)
	err = s.repo.Delete(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"name": existing.Name, "status": "active"}
		newVals := map[string]interface{}{"status": "deleted"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "DELETE", "bill_types", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *billTypeService) Restore(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	err := s.repo.Restore(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"name": existing.Name, "status": "deleted"}
		newVals := map[string]interface{}{"status": "active"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "RESTORE", "bill_types", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *billTypeService) ToggleStatus(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	err := s.repo.ToggleStatus(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"is_active": existing.IsActive}
		newVals := map[string]interface{}{"is_active": !existing.IsActive}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "TOGGLE_STATUS", "bill_types", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *billTypeService) BulkDelete(ctx context.Context, ids []uint) error {
	for _, id := range ids {
		if err := s.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (s *billTypeService) BulkRestore(ctx context.Context, ids []uint) error {
	err := s.repo.BulkRestore(ctx, ids)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		for _, id := range ids {
			oldVals := map[string]interface{}{"status": "deleted"}
			newVals := map[string]interface{}{"status": "active"}
			_ = s.audit.Log(ctx, s.db, userID, userName, role, "BULK_RESTORE", "bill_types", id, oldVals, newVals, ipAddress, userAgent)
		}
	}
	return err
}

func (s *billTypeService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	rulesCount, err := s.repo.CountBillingRules(ctx, id)
	if err != nil {
		return nil, err
	}

	billsCount, err := s.repo.CountStudentBills(ctx, id)
	if err != nil {
		return nil, err
	}

	var messages []string
	if rulesCount > 0 {
		messages = append(messages, fmt.Sprintf("%d aturan tagihan", rulesCount))
	}
	if billsCount > 0 {
		messages = append(messages, fmt.Sprintf("%d tagihan siswa yang belum ditarik", billsCount))
	}

	return map[string]interface{}{
		"has_dependencies": len(messages) > 0,
		"message":          strings.Join(messages, ", "),
		"counts": map[string]int{
			"billing_rules": rulesCount,
			"student_bills": billsCount,
		},
	}, nil
}

func (s *billTypeService) CheckUnique(ctx context.Context, name string, excludeID uint) (bool, error) {
	if excludeID == 0 {
		return s.repo.Exists(ctx, name)
	}
	return s.repo.ExistsExcludeID(ctx, name, excludeID)
}
