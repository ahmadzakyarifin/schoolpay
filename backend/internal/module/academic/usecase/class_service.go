package usecase

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
)

type ClassService interface {
	Create(ctx context.Context, c *domain.Class) error
	GetAll(ctx context.Context, page, limit int, search, status, majorID, ayID, sort string) ([]domain.Class, int, error)
	Update(ctx context.Context, c *domain.Class) error
	Delete(ctx context.Context, id uint) error
	Restore(ctx context.Context, id uint) error
	ToggleStatus(ctx context.Context, id uint) error
	BulkDelete(ctx context.Context, ids []uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetByID(ctx context.Context, id uint) (*domain.Class, error)
	SuggestNextName(ctx context.Context, name string, ayID, majorID, excludeID uint) (string, error)
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, name string, majorID, ayID, excludeID uint) (bool, error)
}

type classService struct {
	db        bun.IDB
	repo      repository.ClassRepo
	majorRepo repository.MajorRepo
	audit     auditusecase.AuditLogService
}

func NewClassService(db bun.IDB, repo repository.ClassRepo, majorRepo repository.MajorRepo, audit auditusecase.AuditLogService) ClassService {
	return &classService{db: db, repo: repo, majorRepo: majorRepo, audit: audit}
}

func classDependencyMessage(studentCount, ayCount int) string {
	var messages []string
	if studentCount > 0 {
		messages = append(messages, fmt.Sprintf("%d siswa aktif", studentCount))
	}
	if ayCount > 0 {
		messages = append(messages, fmt.Sprintf("%d angkatan", ayCount))
	}
	return strings.Join(messages, ", ")
}

func (s *classService) dependencyCounts(ctx context.Context, id uint) (int, int, error) {
	studentCount, err := s.repo.CountStudents(ctx, id)
	if err != nil {
		return 0, 0, err
	}
	ayCount, err := s.repo.CountAcademicYears(ctx, id)
	if err != nil {
		return 0, 0, err
	}
	return studentCount, ayCount, nil
}

func (s *classService) Create(ctx context.Context, c *domain.Class) error {
	c.Name = strings.TrimSpace(c.Name)
	if c.MajorID == nil || *c.MajorID == 0 {
		return fmt.Errorf("gagal: Jurusan wajib dipilih")
	}

	majorObj, err := s.majorRepo.FindByID(ctx, *c.MajorID)
	if err != nil {
		return fmt.Errorf("gagal: Jurusan tidak ditemukan")
	}
	if !majorObj.IsActive {
		return fmt.Errorf("gagal: Jurusan yang dipilih sudah tidak aktif")
	}

	ayID := uint(0)
	if c.AcademicYearID != nil {
		ayID = *c.AcademicYearID
	}

	exists, err := s.repo.Exists(ctx, c.Name, ayID, *c.MajorID, 0)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("gagal: Kelas dengan nama '%s' sudah ada di angkatan dan jurusan yang sama", c.Name)
	}

	err = s.repo.Create(ctx, c)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		newVals := map[string]interface{}{
			"name": c.Name, "major_id": c.MajorID, "academic_year_id": c.AcademicYearID, "is_active": c.IsActive,
		}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "CREATE", "classes", c.ID, nil, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *classService) GetAll(ctx context.Context, page, limit int, search, status, majorID, ayID, sort string) ([]domain.Class, int, error) {
	return s.repo.FindAll(ctx, page, limit, search, status, majorID, ayID, sort)
}

func (s *classService) GetByID(ctx context.Context, id uint) (*domain.Class, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *classService) Update(ctx context.Context, c *domain.Class) error {
	c.Name = strings.TrimSpace(c.Name)
	if c.MajorID == nil || *c.MajorID == 0 {
		return fmt.Errorf("gagal: Jurusan wajib dipilih")
	}

	majorObj, err := s.majorRepo.FindByID(ctx, *c.MajorID)
	if err != nil {
		return fmt.Errorf("gagal: Jurusan tidak ditemukan")
	}
	if !majorObj.IsActive {
		return fmt.Errorf("gagal: Jurusan yang dipilih sudah tidak aktif")
	}

	ayID := uint(0)
	if c.AcademicYearID != nil {
		ayID = *c.AcademicYearID
	}

	exists, err := s.repo.Exists(ctx, c.Name, ayID, *c.MajorID, c.ID)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("gagal: Nama kelas '%s' sudah digunakan di angkatan dan jurusan yang sama", c.Name)
	}

	existing, _ := s.repo.FindByID(ctx, c.ID)
	err = s.repo.Update(ctx, c)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"name": existing.Name, "major_id": existing.MajorID, "academic_year_id": existing.AcademicYearID, "is_active": existing.IsActive}
		newVals := map[string]interface{}{"name": c.Name, "major_id": c.MajorID, "academic_year_id": c.AcademicYearID, "is_active": c.IsActive}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "UPDATE", "classes", c.ID, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *classService) Delete(ctx context.Context, id uint) error {
	studentCount, ayCount, err := s.dependencyCounts(ctx, id)
	if err != nil {
		return err
	}
	if studentCount+ayCount > 0 {
		return fmt.Errorf("gagal: Kelas tidak dapat dihapus karena masih terhubung dengan %s", classDependencyMessage(studentCount, ayCount))
	}

	existing, _ := s.repo.FindByID(ctx, id)
	err = s.repo.Delete(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"is_active": existing.IsActive, "status": "active"}
		newVals := map[string]interface{}{"is_active": false, "status": "deleted"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "DELETE", "classes", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *classService) Restore(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	err := s.repo.Restore(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"is_active": existing.IsActive, "status": "deleted"}
		newVals := map[string]interface{}{"is_active": true, "status": "active"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "RESTORE", "classes", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *classService) ToggleStatus(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	err := s.repo.ToggleStatus(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"is_active": existing.IsActive}
		newVals := map[string]interface{}{"is_active": !existing.IsActive}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "TOGGLE_STATUS", "classes", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *classService) BulkDelete(ctx context.Context, ids []uint) error {
	for _, id := range ids {
		studentCount, ayCount, err := s.dependencyCounts(ctx, id)
		if err != nil {
			return err
		}
		if studentCount+ayCount > 0 {
			return fmt.Errorf("gagal: Kelas tidak dapat dihapus karena masih terhubung dengan %s", classDependencyMessage(studentCount, ayCount))
		}
	}

	err := s.repo.BulkDelete(ctx, ids)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		for _, id := range ids {
			oldVals := map[string]interface{}{"status": "active"}
			newVals := map[string]interface{}{"status": "deleted"}
			_ = s.audit.Log(ctx, s.db, userID, userName, role, "BULK_DELETE", "classes", id, oldVals, newVals, ipAddress, userAgent)
		}
	}
	return err
}

func (s *classService) BulkRestore(ctx context.Context, ids []uint) error {
	err := s.repo.BulkRestore(ctx, ids)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		for _, id := range ids {
			oldVals := map[string]interface{}{"status": "deleted"}
			newVals := map[string]interface{}{"status": "active"}
			_ = s.audit.Log(ctx, s.db, userID, userName, role, "BULK_RESTORE", "classes", id, oldVals, newVals, ipAddress, userAgent)
		}
	}
	return err
}

func (s *classService) SuggestNextName(ctx context.Context, name string, ayID, majorID, excludeID uint) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", nil
	}

	re := regexp.MustCompile(`^(.*?)(\d+)$`)
	matches := re.FindStringSubmatch(name)

	var baseName string
	var currentNum int

	if len(matches) == 3 {
		baseName = matches[1]
		currentNum, _ = strconv.Atoi(matches[2])
	} else {
		baseName = name
		if !strings.HasSuffix(baseName, " ") {
			baseName += " "
		}
		currentNum = 0
	}

	for i := currentNum + 1; i <= currentNum+100; i++ {
		nextName := strings.TrimSpace(baseName + strconv.Itoa(i))
		// Check using the current class scope when provided.
		exists, err := s.repo.Exists(ctx, nextName, ayID, majorID, excludeID)
		if err != nil {
			return "", err
		}
		if !exists {
			return nextName, nil
		}
	}

	return "", nil
}
func (s *classService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	classObj, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	studentCount, err := s.repo.CountStudents(ctx, id)
	if err != nil {
		return nil, err
	}

	ayCount, err := s.repo.CountAcademicYears(ctx, id)
	if err != nil {
		return nil, err
	}

	message := classDependencyMessage(studentCount, ayCount)
	var messages []string
	if message != "" {
		messages = append(messages, message)
	}
	if classObj.MajorName != nil && *classObj.MajorName != "" {
		messages = append(messages, fmt.Sprintf("bagian dari jurusan %s", *classObj.MajorName))
	}

	return map[string]interface{}{
		"has_dependencies": len(messages) > 0,
		"message":          strings.Join(messages, ", "),
		"counts": map[string]int{
			"students":       studentCount,
			"academic_years": ayCount,
		},
	}, nil
}

func (s *classService) CheckUnique(ctx context.Context, name string, majorID, ayID, excludeID uint) (bool, error) {
	return s.repo.Exists(ctx, name, ayID, majorID, excludeID)
}
