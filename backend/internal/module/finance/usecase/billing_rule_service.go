package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"

	academicdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/uptrace/bun"
)

type BillingRuleService interface {
	Create(ctx context.Context, br *domain.BillingRule) error
	GetAll(ctx context.Context) ([]domain.BillingRule, error)
	GetAllPaged(ctx context.Context, page, limit int, search, status, generateStatus, sort string) ([]domain.BillingRule, int, error)
	Update(ctx context.Context, br *domain.BillingRule) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetByID(ctx context.Context, id uint) (*domain.BillingRule, error)
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, billTypeID uint, targetType string, targetID uint, classID *uint, excludeID uint) (bool, error)
}

type billingRuleService struct {
	db       *bun.DB
	repo     repository.BillingRuleRepo
	notifSvc notificationusecase.FinanceNotificationService
	audit    auditusecase.AuditLogService
}

func NewBillingRuleService(db *bun.DB, repo repository.BillingRuleRepo, notifSvc notificationusecase.FinanceNotificationService, audit auditusecase.AuditLogService) BillingRuleService {
	return &billingRuleService{db: db, repo: repo, notifSvc: notifSvc, audit: audit}
}

func (s *billingRuleService) Create(ctx context.Context, br *domain.BillingRule) error {
	if br.Amount <= 0 {
		return errors.New("gagal: Nominal tagihan wajib diisi dengan angka lebih dari 0")
	}

	if br.StartDate == nil || br.EndDate == nil {
		return errors.New("gagal: Tanggal Masa Berlaku (Mulai & Sampai) wajib diisi dengan format yang valid")
	}

	if br.StartDate.After(*br.EndDate) {
		return errors.New("gagal: Tanggal selesai tidak boleh lebih awal dari tanggal mulai")
	}

	if br.DueDay == 0 {
		br.DueDay = 10
	}
	if br.DueDay < 1 || br.DueDay > 31 {
		return errors.New("gagal: Tanggal jatuh tempo bulanan harus berada di antara tanggal 1 sampai 31")
	}

	if br.AllowInstallment && (br.MaxInstallment == nil || *br.MaxInstallment <= 0) {
		return errors.New("gagal: Maksimal cicilan wajib diisi dengan angka lebih dari 0 jika opsi cicilan diaktifkan")
	}

	// Validasi BillType aktif
	var billType domain.BillType
	if err := s.db.NewSelect().Model(&billType).Where("id = ? AND deleted_at IS NULL", br.BillTypeID).Scan(ctx); err != nil {
		return errors.New("gagal: Jenis tagihan tidak ditemukan")
	} else if !billType.IsActive {
		return errors.New("gagal: Jenis tagihan yang dipilih sudah tidak aktif")
	}

	// Validasi Target aktif
	switch br.TargetType {
	case "academic_year":
		var ay academicdomain.AcademicYear
		if err := s.db.NewSelect().Model(&ay).Where("id = ? AND deleted_at IS NULL", br.TargetID).Scan(ctx); err != nil {
			return errors.New("gagal: Angkatan tidak ditemukan")
		} else if !ay.IsActive {
			return errors.New("gagal: Angkatan yang dipilih sudah tidak aktif")
		}
	case "major":
		var major academicdomain.Major
		if err := s.db.NewSelect().Model(&major).Where("id = ? AND deleted_at IS NULL", br.TargetID).Scan(ctx); err != nil {
			return errors.New("gagal: Jurusan tidak ditemukan")
		} else if !major.IsActive {
			return errors.New("gagal: Jurusan yang dipilih sudah tidak aktif")
		}
	case "class":
		var class academicdomain.Class
		if err := s.db.NewSelect().Model(&class).Where("id = ? AND deleted_at IS NULL", br.TargetID).Scan(ctx); err != nil {
			return errors.New("gagal: Kelas tidak ditemukan")
		} else if !class.IsActive {
			return errors.New("gagal: Kelas yang dipilih sudah tidak aktif")
		}
	}

	// Validasi Class aktif (jika class_id diisi sebagai filter tambahan pada target major)
	if br.ClassID != nil && *br.ClassID > 0 {
		var class academicdomain.Class
		if err := s.db.NewSelect().Model(&class).Where("id = ? AND deleted_at IS NULL", *br.ClassID).Scan(ctx); err != nil {
			return errors.New("gagal: Kelas tidak ditemukan")
		} else if !class.IsActive {
			return errors.New("gagal: Kelas yang dipilih sudah tidak aktif")
		}
		if br.TargetType == "major" && class.MajorID != nil && *class.MajorID != br.TargetID {
			return errors.New("gagal: Kelas spesifik tidak bernaung di bawah jurusan yang dipilih")
		}
	}

	exists, err := s.repo.Exists(ctx, br.BillTypeID, br.TargetType, br.TargetID, br.ClassID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("gagal: Aturan tagihan untuk target ini sudah terdaftar")
	}

	err = s.repo.Create(ctx, br)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		newVals := map[string]interface{}{"bill_type_id": br.BillTypeID, "amount": br.Amount, "target_type": br.TargetType, "target_id": br.TargetID, "class_id": br.ClassID, "period_type": br.PeriodType, "allow_installment": br.AllowInstallment, "due_day": br.DueDay, "start_date": br.StartDate, "end_date": br.EndDate}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "CREATE", "billing_rules", br.ID, nil, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *billingRuleService) GetAll(ctx context.Context) ([]domain.BillingRule, error) {
	return s.repo.FindAll(ctx)
}

func (s *billingRuleService) GetAllPaged(ctx context.Context, page, limit int, search, status, generateStatus, sort string) ([]domain.BillingRule, int, error) {
	return s.repo.FindAllPaged(ctx, page, limit, search, status, generateStatus, sort)
}

func (s *billingRuleService) GetByID(ctx context.Context, id uint) (*domain.BillingRule, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *billingRuleService) Update(ctx context.Context, br *domain.BillingRule) error {
	if br.Amount <= 0 {
		return errors.New("gagal: Nominal tagihan wajib diisi dengan angka lebih dari 0")
	}

	if br.StartDate == nil || br.EndDate == nil {
		return errors.New("gagal: Tanggal Masa Berlaku (Mulai & Sampai) wajib diisi dengan format yang valid")
	}

	if br.StartDate.After(*br.EndDate) {
		return errors.New("gagal: Tanggal selesai tidak boleh lebih awal dari tanggal mulai")
	}

	if br.DueDay == 0 {
		br.DueDay = 10
	}
	if br.DueDay < 1 || br.DueDay > 31 {
		return errors.New("gagal: Tanggal jatuh tempo bulanan harus berada di antara tanggal 1 sampai 31")
	}

	if br.AllowInstallment && (br.MaxInstallment == nil || *br.MaxInstallment <= 0) {
		return errors.New("gagal: Maksimal cicilan wajib diisi dengan angka lebih dari 0 jika opsi cicilan diaktifkan")
	}

	// Validasi BillType aktif
	var billType domain.BillType
	if err := s.db.NewSelect().Model(&billType).Where("id = ? AND deleted_at IS NULL", br.BillTypeID).Scan(ctx); err != nil {
		return errors.New("gagal: Jenis tagihan tidak ditemukan")
	} else if !billType.IsActive {
		return errors.New("gagal: Jenis tagihan yang dipilih sudah tidak aktif")
	}

	// Validasi Target aktif
	switch br.TargetType {
	case "academic_year":
		var ay academicdomain.AcademicYear
		if err := s.db.NewSelect().Model(&ay).Where("id = ? AND deleted_at IS NULL", br.TargetID).Scan(ctx); err != nil {
			return errors.New("gagal: Angkatan tidak ditemukan")
		} else if !ay.IsActive {
			return errors.New("gagal: Angkatan yang dipilih sudah tidak aktif")
		}
	case "major":
		var major academicdomain.Major
		if err := s.db.NewSelect().Model(&major).Where("id = ? AND deleted_at IS NULL", br.TargetID).Scan(ctx); err != nil {
			return errors.New("gagal: Jurusan tidak ditemukan")
		} else if !major.IsActive {
			return errors.New("gagal: Jurusan yang dipilih sudah tidak aktif")
		}
	case "class":
		var class academicdomain.Class
		if err := s.db.NewSelect().Model(&class).Where("id = ? AND deleted_at IS NULL", br.TargetID).Scan(ctx); err != nil {
			return errors.New("gagal: Kelas tidak ditemukan")
		} else if !class.IsActive {
			return errors.New("gagal: Kelas yang dipilih sudah tidak aktif")
		}
	}

	// Validasi Class aktif
	if br.ClassID != nil && *br.ClassID > 0 {
		var class academicdomain.Class
		if err := s.db.NewSelect().Model(&class).Where("id = ? AND deleted_at IS NULL", *br.ClassID).Scan(ctx); err != nil {
			return errors.New("gagal: Kelas tidak ditemukan")
		} else if !class.IsActive {
			return errors.New("gagal: Kelas yang dipilih sudah tidak aktif")
		}
		if br.TargetType == "major" && class.MajorID != nil && *class.MajorID != br.TargetID {
			return errors.New("gagal: Kelas spesifik tidak bernaung di bawah jurusan yang dipilih")
		}
	}

	exists, err := s.repo.ExistsExcludeID(ctx, br.BillTypeID, br.TargetType, br.TargetID, br.ClassID, br.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("gagal: Aturan tagihan serupa sudah ada untuk target yang sama")
	}

	existing, _ := s.repo.FindByID(ctx, br.ID)
	err = s.repo.Update(ctx, br)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"bill_type_id": existing.BillTypeID, "amount": existing.Amount, "target_type": existing.TargetType, "target_id": existing.TargetID, "class_id": existing.ClassID, "period_type": existing.PeriodType, "allow_installment": existing.AllowInstallment, "due_day": existing.DueDay, "start_date": existing.StartDate, "end_date": existing.EndDate}
		newVals := map[string]interface{}{"bill_type_id": br.BillTypeID, "amount": br.Amount, "target_type": br.TargetType, "target_id": br.TargetID, "class_id": br.ClassID, "period_type": br.PeriodType, "allow_installment": br.AllowInstallment, "due_day": br.DueDay, "start_date": br.StartDate, "end_date": br.EndDate}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "UPDATE", "billing_rules", br.ID, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *billingRuleService) Delete(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return fmt.Errorf("gagal: Aturan tagihan tidak ditemukan")
	}

	billsCount, err := s.repo.CountStudentBills(ctx, id)
	if err != nil {
		return err
	}
	if billsCount > 0 {
		return fmt.Errorf("gagal: Aturan tagihan tidak dapat dihapus karena sudah menerbitkan %d tagihan siswa yang belum ditarik. Nonaktifkan aturan untuk menghentikan penerbitan berikutnya, atau gunakan fitur Tarik Tagihan jika tagihan harus dibatalkan", billsCount)
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"bill_type_id": existing.BillTypeID, "status": "active"}
		newVals := map[string]interface{}{"status": "deleted"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "DELETE", "billing_rules", id, oldVals, newVals, ipAddress, userAgent)
	}

	return nil
}

func (s *billingRuleService) Restore(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	err := s.repo.Restore(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"bill_type_id": existing.BillTypeID, "status": "deleted"}
		newVals := map[string]interface{}{"status": "active"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "RESTORE", "billing_rules", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *billingRuleService) ToggleStatus(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	err := s.repo.ToggleStatus(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"is_active": existing.IsActive}
		newVals := map[string]interface{}{"is_active": !existing.IsActive}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "TOGGLE_STATUS", "billing_rules", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *billingRuleService) BulkDelete(ctx context.Context, ids []uint) error {
	for _, id := range ids {
		if err := s.Delete(ctx, id); err != nil {
			return err
		}
	}
	return nil
}

func (s *billingRuleService) BulkRestore(ctx context.Context, ids []uint) error {
	err := s.repo.BulkRestore(ctx, ids)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		for _, id := range ids {
			oldVals := map[string]interface{}{"status": "deleted"}
			newVals := map[string]interface{}{"status": "active"}
			_ = s.audit.Log(ctx, s.db, userID, userName, role, "BULK_RESTORE", "billing_rules", id, oldVals, newVals, ipAddress, userAgent)
		}
	}
	return err
}

func (s *billingRuleService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	billsCount, err := s.repo.CountStudentBills(ctx, id)
	if err != nil {
		return nil, err
	}

	var messages []string
	if billsCount > 0 {
		messages = append(messages, fmt.Sprintf("%d tagihan siswa yang belum ditarik", billsCount))
	}

	return map[string]interface{}{
		"has_dependencies": len(messages) > 0,
		"message":          strings.Join(messages, ", "),
		"counts": map[string]int{
			"student_bills": billsCount,
		},
	}, nil
}

func (s *billingRuleService) CheckUnique(ctx context.Context, billTypeID uint, targetType string, targetID uint, classID *uint, excludeID uint) (bool, error) {
	if excludeID == 0 {
		return s.repo.Exists(ctx, billTypeID, targetType, targetID, classID)
	}
	return s.repo.ExistsExcludeID(ctx, billTypeID, targetType, targetID, classID, excludeID)
}
