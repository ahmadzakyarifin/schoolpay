package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
)

type MajorService interface {
	Create(ctx context.Context, j *domain.Major) error
	GetAll(ctx context.Context, page, limit int, search, status, sort string) ([]domain.Major, int, error)
	Update(ctx context.Context, j *domain.Major) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetByID(ctx context.Context, id uint) (*domain.Major, error)
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, field string, value string, excludeID uint) (bool, error)
}

type majorService struct {
	db    bun.IDB
	repo  repository.MajorRepo
	audit auditusecase.AuditLogService
}

func NewMajorService(db bun.IDB, repo repository.MajorRepo, audit auditusecase.AuditLogService) MajorService {
	return &majorService{db: db, repo: repo, audit: audit}
}

func majorDependencyMessage(classCount, studentCount, ayCount int) string {
	var messages []string
	if classCount > 0 {
		messages = append(messages, fmt.Sprintf("%d kelas aktif", classCount))
	}
	if studentCount > 0 {
		messages = append(messages, fmt.Sprintf("%d siswa aktif", studentCount))
	}
	if ayCount > 0 {
		messages = append(messages, fmt.Sprintf("%d angkatan", ayCount))
	}
	return strings.Join(messages, ", ")
}

func (s *majorService) dependencyCounts(ctx context.Context, id uint) (int, int, int, error) {
	classCount, err := s.repo.CountClasses(ctx, id)
	if err != nil {
		return 0, 0, 0, err
	}
	studentCount, err := s.repo.CountStudents(ctx, id)
	if err != nil {
		return 0, 0, 0, err
	}
	ayCount, err := s.repo.CountAcademicYears(ctx, id)
	if err != nil {
		return 0, 0, 0, err
	}
	return classCount, studentCount, ayCount, nil
}

func normalizeMajorCode(code *string) *string {
	if code == nil {
		return nil
	}
	trimmed := strings.TrimSpace(*code)
	if trimmed == "" {
		return nil
	}
	return &trimmed
}

func (s *majorService) Create(ctx context.Context, j *domain.Major) error {
	j.Name = strings.TrimSpace(j.Name)
	j.Code = normalizeMajorCode(j.Code)
	if j.Code != nil {
		exists, err := s.repo.ExistsByCode(ctx, *j.Code, 0)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("gagal: Kode jurusan '%s' sudah digunakan, termasuk pada data aktif atau riwayat hapus", *j.Code)
		}
	}

	exists, err := s.repo.Exists(ctx, j.Name, 0)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("gagal: Jurusan dengan nama '%s' sudah terdaftar", j.Name)
	}

	err = s.repo.Create(ctx, j)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		newVals := map[string]interface{}{
			"name": j.Name, "code": j.Code, "is_active": j.IsActive,
		}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "CREATE", "majors", j.ID, nil, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *majorService) GetAll(ctx context.Context, page, limit int, search, status, sort string) ([]domain.Major, int, error) {
	return s.repo.FindAll(ctx, page, limit, search, status, sort)
}

func (s *majorService) GetByID(ctx context.Context, id uint) (*domain.Major, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *majorService) Update(ctx context.Context, j *domain.Major) error {
	j.Name = strings.TrimSpace(j.Name)
	j.Code = normalizeMajorCode(j.Code)

	existing, err := s.repo.FindByID(ctx, j.ID)
	if err != nil || existing == nil {
		return fmt.Errorf("jurusan tidak ditemukan")
	}

	if j.Code != nil {
		exists, err := s.repo.ExistsByCode(ctx, *j.Code, j.ID)
		if err != nil {
			return err
		}
		if exists {
			return fmt.Errorf("gagal: Kode jurusan '%s' sudah digunakan, termasuk pada data aktif atau riwayat hapus", *j.Code)
		}
	}

	exists, err := s.repo.Exists(ctx, j.Name, j.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("gagal: Nama jurusan '%s' sudah digunakan", j.Name)
	}

	nameChanged := !strings.EqualFold(strings.TrimSpace(existing.Name), j.Name)
	if nameChanged {
		classCount, studentCount, ayCount, err := s.dependencyCounts(ctx, j.ID)
		if err != nil {
			return err
		}
		if classCount+studentCount+ayCount > 0 {
			return fmt.Errorf("gagal: Nama jurusan tidak dapat diubah karena masih terhubung dengan %s. Nonaktifkan jurusan jika hanya ingin menyembunyikan dari form", majorDependencyMessage(classCount, studentCount, ayCount))
		}
	}

	err = s.repo.Update(ctx, j)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"name": existing.Name, "code": existing.Code, "is_active": existing.IsActive}
		newVals := map[string]interface{}{"name": j.Name, "code": j.Code, "is_active": j.IsActive}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "UPDATE", "majors", j.ID, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *majorService) Delete(ctx context.Context, id uint) error {
	classCount, studentCount, ayCount, err := s.dependencyCounts(ctx, id)
	if err != nil {
		return err
	}
	if classCount+studentCount+ayCount > 0 {
		return fmt.Errorf("gagal: Jurusan tidak dapat dihapus karena masih terhubung dengan %s", majorDependencyMessage(classCount, studentCount, ayCount))
	}

	existing, _ := s.repo.FindByID(ctx, id)
	err = s.repo.Delete(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"is_active": existing.IsActive, "status": "active"}
		newVals := map[string]interface{}{"is_active": false, "status": "deleted"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "DELETE", "majors", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *majorService) Restore(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	err := s.repo.Restore(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"is_active": existing.IsActive, "status": "deleted"}
		newVals := map[string]interface{}{"is_active": true, "status": "active"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "RESTORE", "majors", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *majorService) ToggleStatus(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	err := s.repo.ToggleStatus(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"is_active": existing.IsActive}
		newVals := map[string]interface{}{"is_active": !existing.IsActive}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "TOGGLE_STATUS", "majors", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *majorService) BulkDelete(ctx context.Context, ids []uint) error {
	for _, id := range ids {
		classCount, studentCount, ayCount, err := s.dependencyCounts(ctx, id)
		if err != nil {
			return err
		}
		if classCount+studentCount+ayCount > 0 {
			return fmt.Errorf("gagal: Jurusan tidak dapat dihapus karena masih terhubung dengan %s", majorDependencyMessage(classCount, studentCount, ayCount))
		}
	}

	err := s.repo.BulkDelete(ctx, ids)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		for _, id := range ids {
			oldVals := map[string]interface{}{"status": "active"}
			newVals := map[string]interface{}{"status": "deleted"}
			_ = s.audit.Log(ctx, s.db, userID, userName, role, "BULK_DELETE", "majors", id, oldVals, newVals, ipAddress, userAgent)
		}
	}
	return err
}

func (s *majorService) BulkRestore(ctx context.Context, ids []uint) error {
	err := s.repo.BulkRestore(ctx, ids)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		for _, id := range ids {
			oldVals := map[string]interface{}{"status": "deleted"}
			newVals := map[string]interface{}{"status": "active"}
			_ = s.audit.Log(ctx, s.db, userID, userName, role, "BULK_RESTORE", "majors", id, oldVals, newVals, ipAddress, userAgent)
		}
	}
	return err
}
func (s *majorService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	classCount, studentCount, ayCount, err := s.dependencyCounts(ctx, id)
	if err != nil {
		return nil, err
	}

	message := majorDependencyMessage(classCount, studentCount, ayCount)

	return map[string]interface{}{
		"has_dependencies": message != "",
		"message":          message,
		"counts": map[string]int{
			"classes":        classCount,
			"students":       studentCount,
			"academic_years": ayCount,
		},
	}, nil
}

func (s *majorService) CheckUnique(ctx context.Context, field string, value string, excludeID uint) (bool, error) {
	switch field {
	case "code":
		return s.repo.ExistsByCode(ctx, value, excludeID)
	case "name":
		return s.repo.Exists(ctx, value, excludeID)
	default:
		return false, fmt.Errorf("field %s is not supported for uniqueness check", field)
	}
}
