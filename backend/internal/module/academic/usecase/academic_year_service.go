package usecase

import (
	"context"
	"fmt"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/uptrace/bun"
)

type AcademicYearService interface {
	Create(ctx context.Context, ay *domain.AcademicYear) error
	GetAll(ctx context.Context, page, limit int, search, status, sort string) ([]domain.AcademicYear, int, error)
	Update(ctx context.Context, ay *domain.AcademicYear) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetActive(ctx context.Context) ([]domain.AcademicYear, error)
	AssignMajors(ctx context.Context, ayID uint, majorIDs []uint) error
	GetMajorsByYear(ctx context.Context, ayID uint) ([]domain.Major, error)
	AssignClasses(ctx context.Context, ayID uint, classIDs []uint) error
	GetClassesByYear(ctx context.Context, ayID uint) ([]domain.Class, error)
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, year int, excludeID uint) (bool, error)
}

type academicYearService struct {
	db    bun.IDB
	repo  repository.AcademicYearRepo
	audit auditusecase.AuditLogService
}

func NewAcademicYearService(db bun.IDB, repo repository.AcademicYearRepo, audit auditusecase.AuditLogService) AcademicYearService {
	return &academicYearService{db: db, repo: repo, audit: audit}
}

func academicYearDependencyMessage(majorCount, classCount, studentCount int) string {
	var messages []string
	if majorCount > 0 {
		messages = append(messages, fmt.Sprintf("%d jurusan", majorCount))
	}
	if classCount > 0 {
		messages = append(messages, fmt.Sprintf("%d kelas", classCount))
	}
	if studentCount > 0 {
		messages = append(messages, fmt.Sprintf("%d siswa aktif", studentCount))
	}
	return strings.Join(messages, ", ")
}

func (s *academicYearService) dependencyCounts(ctx context.Context, ay *domain.AcademicYear) (int, int, int, error) {
	if ay == nil {
		return 0, 0, 0, fmt.Errorf("angkatan tidak ditemukan")
	}
	majors, err := s.repo.GetMajorsByYear(ctx, ay.ID)
	if err != nil {
		return 0, 0, 0, err
	}
	classes, err := s.repo.GetClassesByYear(ctx, ay.ID)
	if err != nil {
		return 0, 0, 0, err
	}
	studentCount, err := s.repo.CountStudentsByEntryYear(ctx, ay.Year)
	if err != nil {
		return 0, 0, 0, err
	}
	return len(majors), len(classes), studentCount, nil
}

func (s *academicYearService) Create(ctx context.Context, ay *domain.AcademicYear) error {
	if len(ay.MajorIDs) == 0 {
		return fmt.Errorf("gagal: Pilih minimal satu jurusan")
	}
	if len(ay.ClassIDs) == 0 {
		return fmt.Errorf("gagal: Pilih minimal satu kelas")
	}

	exists, err := s.repo.Exists(ctx, ay.Year, 0)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("gagal: Angkatan tahun %d sudah terdaftar", ay.Year)
	}

	if err := s.repo.Create(ctx, ay); err != nil {
		return err
	}

	if err := s.repo.AssignMajors(ctx, ay.ID, ay.MajorIDs); err != nil {
		return err
	}
	if err := s.repo.AssignClasses(ctx, ay.ID, ay.ClassIDs); err != nil {
		return err
	}

	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		newVals := map[string]interface{}{
			"year": ay.Year, "major_ids": ay.MajorIDs, "class_ids": ay.ClassIDs, "is_active": ay.IsActive,
		}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "CREATE", "academic_years", ay.ID, nil, newVals, ipAddress, userAgent)
	}
	return nil
}

func (s *academicYearService) GetAll(ctx context.Context, page, limit int, search, status, sort string) ([]domain.AcademicYear, int, error) {
	return s.repo.FindAll(ctx, page, limit, search, status, sort)
}

func (s *academicYearService) Update(ctx context.Context, ay *domain.AcademicYear) error {
	if len(ay.MajorIDs) == 0 {
		return fmt.Errorf("gagal: Pilih minimal satu jurusan")
	}
	if len(ay.ClassIDs) == 0 {
		return fmt.Errorf("gagal: Pilih minimal satu kelas")
	}

	exists, err := s.repo.Exists(ctx, ay.Year, ay.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("gagal: Angkatan tahun %d sudah digunakan", ay.Year)
	}

	existing, _ := s.repo.FindByID(ctx, ay.ID)
	if err := s.repo.Update(ctx, ay); err != nil {
		return err
	}

	if err := s.repo.AssignMajors(ctx, ay.ID, ay.MajorIDs); err != nil {
		return err
	}
	if err := s.repo.AssignClasses(ctx, ay.ID, ay.ClassIDs); err != nil {
		return err
	}

	if existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"year": existing.Year, "is_active": existing.IsActive}
		newVals := map[string]interface{}{"year": ay.Year, "major_ids": ay.MajorIDs, "class_ids": ay.ClassIDs, "is_active": ay.IsActive}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "UPDATE", "academic_years", ay.ID, oldVals, newVals, ipAddress, userAgent)
	}
	return nil
}

func (s *academicYearService) Delete(ctx context.Context, id uint) error {
	ay, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	majorCount, classCount, studentCount, err := s.dependencyCounts(ctx, ay)
	if err != nil {
		return err
	}
	if majorCount+classCount+studentCount > 0 {
		return fmt.Errorf("gagal: Angkatan tahun %d tidak dapat dihapus karena masih terhubung dengan %s", ay.Year, academicYearDependencyMessage(majorCount, classCount, studentCount))
	}

	err = s.repo.Delete(ctx, id)
	if err == nil && ay != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"year": ay.Year, "is_active": ay.IsActive, "status": "active"}
		newVals := map[string]interface{}{"year": ay.Year, "is_active": false, "status": "deleted"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "DELETE", "academic_years", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *academicYearService) Restore(ctx context.Context, id uint) error {
	ay, _ := s.repo.FindByID(ctx, id)
	err := s.repo.Restore(ctx, id)
	if err == nil && ay != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"year": ay.Year, "is_active": ay.IsActive, "status": "deleted"}
		newVals := map[string]interface{}{"year": ay.Year, "is_active": true, "status": "active"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "RESTORE", "academic_years", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *academicYearService) BulkDelete(ctx context.Context, ids []uint) error {
	for _, id := range ids {
		ay, err := s.repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		majorCount, classCount, studentCount, err := s.dependencyCounts(ctx, ay)
		if err != nil {
			return err
		}
		if majorCount+classCount+studentCount > 0 {
			return fmt.Errorf("gagal: Angkatan tahun %d tidak dapat dihapus karena masih terhubung dengan %s", ay.Year, academicYearDependencyMessage(majorCount, classCount, studentCount))
		}
	}

	err := s.repo.BulkDelete(ctx, ids)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		for _, id := range ids {
			oldVals := map[string]interface{}{"status": "active"}
			newVals := map[string]interface{}{"status": "deleted"}
			_ = s.audit.Log(ctx, s.db, userID, userName, role, "BULK_DELETE", "academic_years", id, oldVals, newVals, ipAddress, userAgent)
		}
	}
	return err
}

func (s *academicYearService) BulkRestore(ctx context.Context, ids []uint) error {
	err := s.repo.BulkRestore(ctx, ids)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		for _, id := range ids {
			oldVals := map[string]interface{}{"status": "deleted"}
			newVals := map[string]interface{}{"status": "active"}
			_ = s.audit.Log(ctx, s.db, userID, userName, role, "BULK_RESTORE", "academic_years", id, oldVals, newVals, ipAddress, userAgent)
		}
	}
	return err
}

func (s *academicYearService) GetActive(ctx context.Context) ([]domain.AcademicYear, error) {
	return s.repo.FindActive(ctx)
}

func (s *academicYearService) AssignMajors(ctx context.Context, ayID uint, majorIDs []uint) error {
	return s.repo.AssignMajors(ctx, ayID, majorIDs)
}

func (s *academicYearService) GetMajorsByYear(ctx context.Context, ayID uint) ([]domain.Major, error) {
	return s.repo.GetMajorsByYear(ctx, ayID)
}

func (s *academicYearService) AssignClasses(ctx context.Context, ayID uint, classIDs []uint) error {
	return s.repo.AssignClasses(ctx, ayID, classIDs)
}

func (s *academicYearService) GetClassesByYear(ctx context.Context, ayID uint) ([]domain.Class, error) {
	return s.repo.GetClassesByYear(ctx, ayID)
}

func (s *academicYearService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	ay, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	majorCount, classCount, studentCount, err := s.dependencyCounts(ctx, ay)
	if err != nil {
		return nil, err
	}
	message := academicYearDependencyMessage(majorCount, classCount, studentCount)

	return map[string]interface{}{
		"has_dependencies": message != "",
		"message":          message,
		"counts": map[string]int{
			"majors":   majorCount,
			"classes":  classCount,
			"students": studentCount,
		},
	}, nil
}

func (s *academicYearService) CheckUnique(ctx context.Context, year int, excludeID uint) (bool, error) {
	return s.repo.Exists(ctx, year, excludeID)
}
