package usecase

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/domain"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	userauthdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
	"github.com/xuri/excelize/v2"
)

type StudentService interface {
	GetPaginated(ctx context.Context, page, limit int, search, filter, status string, entryYear int, classID, majorID uint, sort string) ([]domain.Student, int, error)
	GetAcademicFilters(ctx context.Context) (map[string]interface{}, error)
	GetParents(ctx context.Context, studentID uint) ([]userauthdomain.User, error)
	GetStudentsByParentID(ctx context.Context, parentID uint) ([]domain.Student, error)
	Create(ctx context.Context, s *domain.Student) error
	Update(ctx context.Context, s *domain.Student) error
	Delete(ctx context.Context, id uint) error
	GetByID(ctx context.Context, id uint) (*domain.Student, error)
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	ToggleStatus(ctx context.Context, id uint) error
	ExportExcel(ctx context.Context, search, filter, status string, entryYear int, classID, majorID uint) ([]byte, error)
	BulkGraduate(ctx context.Context, classID uint, studentIDs []uint) error
	BulkPromote(ctx context.Context, sourceClassID, targetClassID uint, studentIDs []uint) error
	GetClassHistory(ctx context.Context, studentID uint) ([]domain.ClassHistory, error)
	BulkDelete(ctx context.Context, ids []uint) error
	Restore(ctx context.Context, id uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	CheckUnique(ctx context.Context, field string, value string, excludeID uint) (bool, error)
}

type studentService struct {
	db        *bun.DB
	repo      repository.StudentRepo
	userRepo  userauthrepo.UserRepo
	authRepo  userauthrepo.AuthRepo
	msg       utils.Messenger
	notiRepo  notificationrepo.NotificationRepo
	ayRepo    repository.AcademicYearRepo
	majorRepo repository.MajorRepo
	classRepo repository.ClassRepo
	cfg       *config.Config
	billRepo  financerepo.StudentBillRepo
	audit     auditusecase.AuditLogService
	jobChan   chan studentNotifyJob
}

type studentNotifyJob struct {
	user        *userauthdomain.User
	studentName string
	rawPassword string
}

func studentAuditValues(student *domain.Student) map[string]interface{} {
	if student == nil {
		return nil
	}

	birthDate := ""
	if !student.BirthDate.Time().IsZero() {
		birthDate = student.BirthDate.Time().Format("2006-01-02")
	}

	return map[string]interface{}{
		"name":         student.Name,
		"nik":          student.NIK,
		"nis":          student.NIS,
		"nisn":         student.NISN,
		"parent_id":    student.ParentID,
		"gender":       student.Gender,
		"birth_place":  student.BirthPlace,
		"birth_date":   birthDate,
		"email":        student.Email,
		"phone_number": student.PhoneNumber,
		"class_id":     student.ClassID,
		"major_id":     student.MajorID,
		"entry_year":   student.EntryYear,
		"address":      student.Address,
		"province":     student.Province,
		"city":         student.City,
		"district":     student.District,
		"village":      student.Village,
		"rt":           student.RT,
		"rw":           student.RW,
		"status":       student.Status,
	}
}

func NewStudentService(db *bun.DB, repo repository.StudentRepo, userRepo userauthrepo.UserRepo, authRepo userauthrepo.AuthRepo, msg utils.Messenger, noti notificationrepo.NotificationRepo, ay repository.AcademicYearRepo, jur repository.MajorRepo, cls repository.ClassRepo, billRepo financerepo.StudentBillRepo, cfg *config.Config, audit auditusecase.AuditLogService) StudentService {
	s := &studentService{
		db:        db,
		repo:      repo,
		userRepo:  userRepo,
		authRepo:  authRepo,
		msg:       msg,
		notiRepo:  noti,
		ayRepo:    ay,
		majorRepo: jur,
		classRepo: cls,
		billRepo:  billRepo,
		cfg:       cfg,
		audit:     audit,
		jobChan:   make(chan studentNotifyJob, 100),
	}
	for i := 0; i < 5; i++ {
		go s.notificationWorker()
	}
	return s
}

func (s *studentService) GetPaginated(ctx context.Context, page, limit int, search, filter, status string, entryYear int, classID, majorID uint, sort string) ([]domain.Student, int, error) {
	return s.repo.FindAllPaginated(ctx, page, limit, search, filter, status, entryYear, classID, majorID, sort)
}

func (s *studentService) GetAcademicFilters(ctx context.Context) (map[string]interface{}, error) {
	ayList, _, _ := s.ayRepo.FindAll(ctx, 1, 1000, "", "active", "")
	majorsFull, _, _ := s.majorRepo.FindAll(ctx, 1, 1000, "", "active", "")

	// Get all classes and their year associations
	var classes []domain.Class
	_ = s.db.NewSelect().Model(&classes).Where("deleted_at IS NULL").Scan(ctx)

	// Populate years for each major
	for i := range majorsFull {
		var yearIDs []uint
		_ = s.db.NewSelect().Table("academic_year_majors").Column("academic_year_id").Where("major_id = ?", majorsFull[i].ID).Scan(ctx, &yearIDs)
		majorsFull[i].YearIDs = yearIDs
	}

	// Populate years for each class
	for i := range classes {
		var yearIDs []uint
		_ = s.db.NewSelect().Table("academic_year_classes").Column("academic_year_id").Where("class_id = ?", classes[i].ID).Scan(ctx, &yearIDs)

		// Also include legacy academic_year_id if exists
		if classes[i].AcademicYearID != nil && *classes[i].AcademicYearID > 0 {
			found := false
			for _, id := range yearIDs {
				if id == *classes[i].AcademicYearID {
					found = true
					break
				}
			}
			if !found {
				yearIDs = append(yearIDs, *classes[i].AcademicYearID)
			}
		}
		classes[i].AcademicYearIDs = yearIDs
	}

	return map[string]interface{}{
		"years":   ayList,
		"majors":  majorsFull,
		"classes": classes,
	}, nil
}

func (s *studentService) GetParents(ctx context.Context, studentID uint) ([]userauthdomain.User, error) {
	return s.repo.GetParents(ctx, studentID)
}

func (s *studentService) GetStudentsByParentID(ctx context.Context, parentID uint) ([]domain.Student, error) {
	return s.repo.GetStudentsByParentID(ctx, parentID)
}

func (s *studentService) Create(ctx context.Context, student *domain.Student) error {
	student.Name = strings.Title(strings.ToLower(strings.TrimSpace(student.Name)))
	student.Email = strings.ToLower(strings.TrimSpace(student.Email))
	student.NIK = strings.TrimSpace(student.NIK)
	student.BirthPlace = strings.TrimSpace(student.BirthPlace)
	student.Religion = strings.TrimSpace(student.Religion)
	student.Address = strings.TrimSpace(student.Address)
	student.RT = strings.TrimSpace(student.RT)
	student.RW = strings.TrimSpace(student.RW)
	student.Description = strings.TrimSpace(student.Description)

	student.PhoneNumber = utils.NormalizePhoneNumber(student.PhoneNumber)

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		bizErrors := make(map[string][]string)

		// 2-Layer Validation: Ensure Class matches Major
		if student.ClassID > 0 && student.MajorID > 0 {
			var class domain.Class
			if err := tx.NewSelect().Model(&class).Where("id = ?", student.ClassID).Scan(ctx); err == nil {
				if class.MajorID != nil && *class.MajorID != student.MajorID {
					bizErrors["major_id"] = append(bizErrors["major_id"], "Jurusan tidak sesuai dengan kelas yang dipilih")
				}
			}
		}

		// 3. Data Integrity & Business Logic
		if student.BirthDate.Time().IsZero() {
			bizErrors["birth_date"] = append(bizErrors["birth_date"], "Tanggal lahir tidak valid atau kosong")
		}

		// Duplicate checks (Batch Check) - Run all, don't return early
		if student.NIK != "" {
			if existing, _ := s.repo.FindByNIK(ctx, student.NIK); existing != nil && (student.ID == 0 || existing.ID != student.ID) {
				bizErrors["nik"] = append(bizErrors["nik"], fmt.Sprintf("NIK '%s' sudah terdaftar sebagai Siswa", student.NIK))
			} else if existing, _ := s.userRepo.FindByNIK(ctx, student.NIK); existing != nil {
				bizErrors["nik"] = append(bizErrors["nik"], fmt.Sprintf("NIK '%s' sudah terdaftar sebagai Pengguna/Wali Murid", student.NIK))
			}
		}

		if student.NISN != "" {
			if existing, _ := s.repo.FindByNISN(ctx, student.NISN); existing != nil && (student.ID == 0 || existing.ID != student.ID) {
				bizErrors["nisn"] = append(bizErrors["nisn"], fmt.Sprintf("NISN '%s' sudah terdaftar", student.NISN))
			}
		}

		if student.NIS != "" {
			if existing, _ := s.repo.FindByNIS(ctx, student.NIS); existing != nil && (student.ID == 0 || existing.ID != student.ID) {
				bizErrors["nis"] = append(bizErrors["nis"], fmt.Sprintf("NIS '%s' sudah terdaftar", student.NIS))
			}
		}

		if student.Email != "" {
			if existing, _ := s.repo.FindByEmail(ctx, student.Email); existing != nil && (student.ID == 0 || existing.ID != student.ID) {
				bizErrors["email"] = append(bizErrors["email"], fmt.Sprintf("Email '%s' sudah terdaftar sebagai Siswa", student.Email))
			} else if existing, _ := s.userRepo.FindByEmail(ctx, student.Email); existing != nil {
				bizErrors["email"] = append(bizErrors["email"], fmt.Sprintf("Email '%s' sudah terdaftar sebagai Pengguna/Wali Murid", student.Email))
			}
		}

		if student.PhoneNumber != "" {
			if existing, _ := s.repo.FindByPhone(ctx, student.PhoneNumber); existing != nil && (student.ID == 0 || existing.ID != student.ID) {
				bizErrors["phone_number"] = append(bizErrors["phone_number"], fmt.Sprintf("Nomor WhatsApp '%s' sudah terdaftar sebagai Siswa", student.PhoneNumber))
			} else if existing, _ := s.userRepo.FindByPhone(ctx, student.PhoneNumber); existing != nil {
				bizErrors["phone_number"] = append(bizErrors["phone_number"], fmt.Sprintf("Nomor WhatsApp '%s' sudah terdaftar sebagai Pengguna/Wali Murid", student.PhoneNumber))
			}
		}

		if len(bizErrors) > 0 {
			return utils.NewBusinessMultiError(bizErrors)
		}

		if err := s.repo.Create(ctx, tx, student); err != nil {
			return err
		}
		if student.ClassID > 0 {
			if err := s.repo.AddClassHistory(ctx, tx, student.ID, student.ClassID); err != nil {
				return err
			}
		}

		if s.audit != nil {
			userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, tx, userID, userName, role, "CREATE", "students", student.ID, nil, studentAuditValues(student), ipAddress, userAgent)
		}
		return nil
	})
}

func (s *studentService) Update(ctx context.Context, student *domain.Student) error {
	// Normalization & Trimming
	student.Name = strings.Title(strings.ToLower(strings.TrimSpace(student.Name)))
	student.Email = strings.ToLower(strings.TrimSpace(student.Email))
	student.NIK = strings.TrimSpace(student.NIK)
	student.BirthPlace = strings.TrimSpace(student.BirthPlace)
	student.Religion = strings.TrimSpace(student.Religion)
	student.Address = strings.TrimSpace(student.Address)
	student.RT = strings.TrimSpace(student.RT)
	student.RW = strings.TrimSpace(student.RW)
	student.Description = strings.TrimSpace(student.Description)

	student.PhoneNumber = utils.NormalizePhoneNumber(student.PhoneNumber)

	bizErrors := make(map[string][]string)

	// 2-Layer Validation: Ensure Class matches Major
	if student.ClassID > 0 && student.MajorID > 0 {
		var class domain.Class
		if err := s.db.NewSelect().Model(&class).Where("id = ?", student.ClassID).Scan(ctx); err == nil {
			if class.MajorID != nil && *class.MajorID != student.MajorID {
				bizErrors["major_id"] = append(bizErrors["major_id"], "Jurusan tidak sesuai dengan kelas yang dipilih")
			}
		}
	}

	// 3. Data Integrity & Business Logic
	if student.BirthDate.Time().IsZero() {
		bizErrors["birth_date"] = append(bizErrors["birth_date"], "Tanggal lahir tidak valid atau kosong")
	}

	// Duplicate checks for Update (exclude current ID)
	if student.NIK != "" {
		if existing, _ := s.repo.FindByNIK(ctx, student.NIK); existing != nil && existing.ID != student.ID {
			bizErrors["nik"] = append(bizErrors["nik"], fmt.Sprintf("NIK '%s' sudah digunakan oleh Siswa lain", student.NIK))
		} else if existing, _ := s.userRepo.FindByNIK(ctx, student.NIK); existing != nil {
			bizErrors["nik"] = append(bizErrors["nik"], fmt.Sprintf("NIK '%s' sudah digunakan oleh Pengguna/Wali Murid lain", student.NIK))
		}
	}

	if student.NISN != "" {
		if existing, _ := s.repo.FindByNISN(ctx, student.NISN); existing != nil && existing.ID != student.ID {
			bizErrors["nisn"] = append(bizErrors["nisn"], fmt.Sprintf("NISN '%s' sudah digunakan oleh Siswa lain", student.NISN))
		}
	}

	if student.NIS != "" {
		if existing, _ := s.repo.FindByNIS(ctx, student.NIS); existing != nil && existing.ID != student.ID {
			bizErrors["nis"] = append(bizErrors["nis"], fmt.Sprintf("NIS '%s' sudah digunakan oleh Siswa lain", student.NIS))
		}
	}

	if student.Email != "" {
		if existing, _ := s.repo.FindByEmail(ctx, student.Email); existing != nil && existing.ID != student.ID {
			bizErrors["email"] = append(bizErrors["email"], fmt.Sprintf("Email '%s' sudah digunakan oleh Siswa lain", student.Email))
		} else if existing, _ := s.userRepo.FindByEmail(ctx, student.Email); existing != nil {
			bizErrors["email"] = append(bizErrors["email"], fmt.Sprintf("Email '%s' sudah digunakan oleh Pengguna/Wali Murid lain", student.Email))
		}
	}

	if student.PhoneNumber != "" {
		if existing, _ := s.repo.FindByPhone(ctx, student.PhoneNumber); existing != nil && existing.ID != student.ID {
			bizErrors["phone_number"] = append(bizErrors["phone_number"], fmt.Sprintf("Nomor WhatsApp '%s' sudah digunakan oleh Siswa lain", student.PhoneNumber))
		} else if existing, _ := s.userRepo.FindByPhone(ctx, student.PhoneNumber); existing != nil {
			bizErrors["phone_number"] = append(bizErrors["phone_number"], fmt.Sprintf("Nomor WhatsApp '%s' sudah digunakan oleh Pengguna/Wali Murid lain", student.PhoneNumber))
		}
	}

	if len(bizErrors) > 0 {
		return utils.NewBusinessMultiError(bizErrors)
	}

	existing, _ := s.repo.FindByID(ctx, student.ID)
	if existing == nil {
		return fmt.Errorf("siswa tidak ditemukan")
	}

	// Persist image_path if not provided (not updated)
	if student.ImagePath == nil || *student.ImagePath == "" {
		student.ImagePath = existing.ImagePath
	}

	if err := s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		if err := s.repo.Update(ctx, tx, student); err != nil {
			return err
		}

		// Edit data siswa hanya memperbaiki posisi akademik aktif.
		// Riwayat baru hanya dibuat oleh proses kenaikan kelas.
		if existing.ClassID != student.ClassID {
			if err := s.repo.UpdateActiveHistory(ctx, tx, student.ID, student.ClassID); err != nil {
				return err
			}
		}

		if s.audit != nil {
			userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, tx, userID, userName, role, "UPDATE", "students", student.ID, studentAuditValues(existing), studentAuditValues(student), ipAddress, userAgent)
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (s *studentService) Delete(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	if existing == nil {
		return fmt.Errorf("siswa tidak ditemukan")
	}

	if err := s.repo.Delete(ctx, id); err != nil {
		return err
	}

	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"status": existing.Status}
		newVals := map[string]interface{}{"status": "deleted"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "DELETE", "students", id, oldVals, newVals, ipAddress, userAgent)
	}
	return nil
}

func (s *studentService) GetByID(ctx context.Context, id uint) (*domain.Student, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *studentService) ToggleStatus(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	if existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"status": existing.Status}
		newStatus := "inactive"
		if existing.Status != "active" {
			newStatus = "active"
		}
		newVals := map[string]interface{}{"status": newStatus}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "UPDATE", "students", id, oldVals, newVals, ipAddress, userAgent)
	}
	return s.repo.ToggleStatus(ctx, id)
}

func (s *studentService) ExportExcel(ctx context.Context, search, filter, status string, entryYear int, classID, majorID uint) ([]byte, error) {
	list, _, _ := s.repo.FindAllPaginated(ctx, 1, 1000000, search, filter, status, entryYear, classID, majorID, "")
	f := excelize.NewFile()
	sheet := "Data Siswa"
	f.SetSheetName("Sheet1", sheet)

	// Headers aligned with import template — full field set for re-import compatibility
	headers := []string{
		"NIS", "NISN", "NAMA LENGKAP", "NIK", "L/P", "TEMPAT LAHIR", "TANGGAL LAHIR",
		"AGAMA", "EMAIL", "NO. WHATSAPP", "PROVINSI", "KOTA", "KECAMATAN", "KELURAHAN",
		"ALAMAT", "RT", "RW", "KELAS", "JURUSAN", "ANGKATAN", "NAMA WALI", "STATUS",
	}
	colWidths := []float64{18, 18, 28, 20, 8, 20, 20, 14, 28, 20, 20, 20, 18, 18, 30, 6, 6, 15, 18, 10, 28, 12}

	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Horizontal: "center"},
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
	})

	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
		colName, _ := excelize.ColumnNumberToName(i + 1)
		if i < len(colWidths) {
			f.SetColWidth(sheet, colName, colName, colWidths[i])
		} else {
			f.SetColWidth(sheet, colName, colName, 18)
		}
	}

	bodyStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	for i, stu := range list {
		row := i + 2
		birthDate := utils.FormatTimeID(stu.BirthDate.Time())
		phoneNumber := stu.PhoneNumber
		if len(phoneNumber) > 0 && phoneNumber[0] != '+' {
			phoneNumber = "+" + phoneNumber
		}
		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), stu.NIS)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), stu.NISN)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), stu.Name)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), stu.NIK)
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), stu.Gender)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), stu.BirthPlace)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), birthDate)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), stu.Religion)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), stu.Email)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), phoneNumber)
		f.SetCellValue(sheet, fmt.Sprintf("K%d", row), stu.Province)
		f.SetCellValue(sheet, fmt.Sprintf("L%d", row), stu.City)
		f.SetCellValue(sheet, fmt.Sprintf("M%d", row), stu.District)
		f.SetCellValue(sheet, fmt.Sprintf("N%d", row), stu.Village)
		f.SetCellValue(sheet, fmt.Sprintf("O%d", row), stu.Address)
		f.SetCellValue(sheet, fmt.Sprintf("P%d", row), stu.RT)
		f.SetCellValue(sheet, fmt.Sprintf("Q%d", row), stu.RW)
		f.SetCellValue(sheet, fmt.Sprintf("R%d", row), utils.StringValue(stu.ClassName))
		f.SetCellValue(sheet, fmt.Sprintf("S%d", row), utils.StringValue(stu.MajorName))
		f.SetCellValue(sheet, fmt.Sprintf("T%d", row), stu.EntryYear)
		f.SetCellValue(sheet, fmt.Sprintf("U%d", row), stu.ParentName)
		f.SetCellValue(sheet, fmt.Sprintf("V%d", row), strings.ToUpper(stu.Status))
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("V%d", row), bodyStyle)
	}

	buf, _ := f.WriteToBuffer()
	return buf.Bytes(), nil
}

func (s *studentService) BulkGraduate(ctx context.Context, classID uint, studentIDs []uint) error {
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		var graduateIDs []uint
		var failedGraduation []string
		var students []domain.Student

		q := tx.NewSelect().Model(&students).Where("status = 'active'")
		if len(studentIDs) > 0 {
			q.Where("id IN (?)", bun.In(studentIDs))
		}
		if classID > 0 {
			q.Where("class_id = ?", classID)
		}

		err := q.Scan(ctx)
		if err != nil {
			return err
		}
		if len(students) == 0 {
			return fmt.Errorf("tidak ada siswa aktif yang cocok untuk diluluskan")
		}

		for _, stu := range students {
			unpaid, _ := s.billRepo.FindUnpaidBillsByStudent(ctx, stu.ID)
			if len(unpaid) > 0 {
				failedGraduation = append(failedGraduation, stu.Name)
				continue
			}
			graduateIDs = append(graduateIDs, stu.ID)
		}

		if len(graduateIDs) == 0 && len(failedGraduation) > 0 {
			return fmt.Errorf("gagal meluluskan: semua siswa yang dipilih masih memiliki tunggakan aktif: %s", strings.Join(failedGraduation, ", "))
		}

		if len(graduateIDs) > 0 {
			_, err := tx.NewUpdate().Model((*domain.Student)(nil)).Set("status = 'graduated'").Where("id IN (?)", bun.In(graduateIDs)).Exec(ctx)
			if err != nil {
				return err
			}
			_, _ = tx.NewUpdate().Table("student_classes").Set("is_active = 0").Where("student_id IN (?)", bun.In(graduateIDs)).Exec(ctx)

			if s.audit != nil {
				userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
				for _, id := range graduateIDs {
					oldVals := map[string]interface{}{"status": "active"}
					newVals := map[string]interface{}{"status": "graduated"}
					_ = s.audit.Log(ctx, tx, userID, userName, role, "UPDATE", "students", id, oldVals, newVals, ipAddress, userAgent)
				}
			}
		}

		if len(failedGraduation) > 0 {
			return fmt.Errorf("sebagian berhasil (%d siswa), namun %d siswa berikut gagal diluluskan karena masih memiliki tunggakan aktif: %s", len(graduateIDs), len(failedGraduation), strings.Join(failedGraduation, ", "))
		}
		return nil
	})
}

func (s *studentService) BulkPromote(ctx context.Context, sourceClassID, targetClassID uint, studentIDs []uint) error {
	if sourceClassID == targetClassID {
		return fmt.Errorf("kelas asal dan tujuan tidak boleh sama")
	}
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		target, _ := s.classRepo.FindByID(ctx, targetClassID)
		if target == nil {
			return fmt.Errorf("kelas tujuan tidak ditemukan")
		}
		if !target.IsActive {
			return fmt.Errorf("kelas tujuan tidak aktif")
		}

		var promoteIDs []uint
		var failedPromotion []string
		var students []domain.Student

		q := tx.NewSelect().Model(&students).Where("status = 'active'")
		if len(studentIDs) > 0 {
			q.Where("id IN (?)", bun.In(studentIDs))
		}
		if sourceClassID > 0 {
			q.Where("class_id = ?", sourceClassID)
		}

		err := q.Scan(ctx)
		if err != nil {
			return err
		}
		if len(students) == 0 {
			return fmt.Errorf("tidak ada siswa aktif yang cocok untuk dipindahkan")
		}

		for _, stu := range students {
			unpaid, _ := s.billRepo.FindUnpaidBillsByStudent(ctx, stu.ID)
			if len(unpaid) > 0 {
				failedPromotion = append(failedPromotion, stu.Name)
				continue
			}
			promoteIDs = append(promoteIDs, stu.ID)
		}

		if len(promoteIDs) == 0 && len(failedPromotion) > 0 {
			return fmt.Errorf("gagal dipindahkan: semua siswa yang dipilih masih memiliki tunggakan aktif: %s", strings.Join(failedPromotion, ", "))
		}

		if len(promoteIDs) > 0 {
			_, _ = tx.NewUpdate().Table("student_classes").Set("is_active = 0").Where("is_active = 1 AND student_id IN (?)", bun.In(promoteIDs)).Exec(ctx)
			_, err = tx.NewUpdate().Model((*domain.Student)(nil)).Set("class_id = ?", targetClassID).Where("id IN (?)", bun.In(promoteIDs)).Exec(ctx)
			if err != nil {
				return err
			}

			for _, id := range promoteIDs {
				if err := s.repo.AddClassHistory(ctx, tx, id, targetClassID); err != nil {
					return err
				}
			}

			if s.audit != nil {
				userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
				for _, id := range promoteIDs {
					oldVals := map[string]interface{}{"class_id": sourceClassID}
					newVals := map[string]interface{}{"class_id": targetClassID}
					_ = s.audit.Log(ctx, tx, userID, userName, role, "UPDATE", "students", id, oldVals, newVals, ipAddress, userAgent)
				}
			}
		}

		if len(failedPromotion) > 0 {
			return fmt.Errorf("sebagian berhasil (%d siswa), namun %d siswa berikut gagal dipindahkan karena masih memiliki tunggakan aktif: %s", len(promoteIDs), len(failedPromotion), strings.Join(failedPromotion, ", "))
		}

		return nil
	})
}

func (s *studentService) GetClassHistory(ctx context.Context, studentID uint) ([]domain.ClassHistory, error) {
	return s.repo.GetClassHistory(ctx, studentID)
}

func (s *studentService) BulkDelete(ctx context.Context, ids []uint) error {
	if len(ids) == 0 {
		return nil
	}
	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		res, err := tx.NewUpdate().Model((*domain.Student)(nil)).Set("status = 'inactive'").Set("deleted_at = ?", time.Now()).Where("id IN (?) AND deleted_at IS NULL", bun.In(ids)).Exec(ctx)
		if err != nil {
			return err
		}
		rows, _ := res.RowsAffected()
		if rows == 0 {
			return fmt.Errorf("tidak ada data siswa yang perlu dihapus")
		}
		_, _ = tx.NewUpdate().Table("student_classes").Set("is_active = 0").Where("student_id IN (?) AND is_active = 1", bun.In(ids)).Exec(ctx)

		if s.audit != nil {
			userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			for _, id := range ids {
				oldVals := map[string]interface{}{"status": "active"}
				newVals := map[string]interface{}{"status": "deleted"}
				_ = s.audit.Log(ctx, tx, userID, userName, role, "DELETE", "students", id, oldVals, newVals, ipAddress, userAgent)
			}
		}

		return nil
	})
}

func (s *studentService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	student, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	bills, err := s.billRepo.FindUnpaidBillsByStudent(ctx, id)
	if err != nil {
		return nil, err
	}

	var messages []string
	if len(bills) > 0 {
		messages = append(messages, fmt.Sprintf("memiliki %d tagihan yang belum lunas", len(bills)))
	}
	if student.ParentName != "" {
		messages = append(messages, fmt.Sprintf("terhubung dengan orang tua (%s)", student.ParentName))
	}

	return map[string]interface{}{
		"has_dependencies": len(messages) > 0,
		"message":          strings.Join(messages, " dan "),
		"counts": map[string]int{
			"unpaid_bills": len(bills),
		},
	}, nil
}

func (s *studentService) Restore(ctx context.Context, id uint) error {
	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"status": "deleted"}
		newVals := map[string]interface{}{"status": "active"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "RESTORE", "students", id, oldVals, newVals, ipAddress, userAgent)
	}
	return s.repo.Restore(ctx, id)
}

func (s *studentService) BulkRestore(ctx context.Context, ids []uint) error {
	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		for _, id := range ids {
			oldVals := map[string]interface{}{"status": "deleted"}
			newVals := map[string]interface{}{"status": "active"}
			_ = s.audit.Log(ctx, s.db, userID, userName, role, "RESTORE", "students", id, oldVals, newVals, ipAddress, userAgent)
		}
	}
	return s.repo.BulkRestore(ctx, ids)
}

func (s *studentService) notificationWorker() {
	ctx := context.Background()
	for job := range s.jobChan {
		if !job.user.IsActive {
			token := utils.GenerateUUID()
			expiry := time.Now().Add(24 * 7 * time.Hour)
			_ = s.authRepo.SaveAuthToken(ctx, job.user.ID, token, "activation", expiry)
			link := fmt.Sprintf("%s/activate?token=%s", strings.TrimSuffix(s.cfg.FrontendURL, "/"), token)
			passInfo := ""
			if job.rawPassword != "" {
				passInfo = fmt.Sprintf("\n• Password sementara: *%s*", job.rawPassword)
			}
			msg := strings.Join([]string{
				"📢 *AKTIVASI AKUN SCHOOLPAY*",
				"",
				fmt.Sprintf("Halo *%s*,", job.user.Name),
				"",
				"Akun Anda telah didaftarkan dengan rincian:",
				"",
				fmt.Sprintf("• Siswa: *%s*", job.studentName),
				fmt.Sprintf("• Link aktivasi: %s%s", link, passInfo),
				"",
				"Tautan berlaku selama 7 hari.",
				"Silakan hubungi Admin Sekolah jika mengalami kendala.",
				"",
				"Terima kasih.",
			}, "\n")
			status := "sent"
			var deliveryErr *string
			var whatsappID *string
			if waID, err := s.msg.SendWhatsApp(job.user.PhoneNumber, msg); err != nil {
				status = "failed"
				deliveryErr = utils.StringPtr(err.Error())
			} else if waID != "" {
				whatsappID = utils.StringPtr(waID)
			}
			if whatsappID == nil {
				whatsappID = utils.StringPtr(fmt.Sprintf("local-wa-%d-%d", job.user.ID, time.Now().UnixNano()))
			}

			_ = s.notiRepo.Create(ctx, s.db, &notificationdomain.Notification{
				UserID:         job.user.ID,
				Title:          "Aktivasi Akun SchoolPay",
				Message:        msg,
				Type:           "auth",
				Channel:        "whatsapp",
				WhatsappID:     whatsappID,
				DeliveryStatus: status,
				DeliveryError:  deliveryErr,
			})
		}
	}
}

func (s *studentService) CheckUnique(ctx context.Context, field string, value string, excludeID uint) (bool, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return true, nil
	}

	switch field {
	case "nik":
		if existing, _ := s.repo.FindByNIK(ctx, value); existing != nil && existing.ID != excludeID {
			return false, nil
		}
		if existing, _ := s.userRepo.FindByNIK(ctx, value); existing != nil {
			return false, nil
		}
	case "nisn":
		if existing, _ := s.repo.FindByNISN(ctx, value); existing != nil && existing.ID != excludeID {
			return false, nil
		}
	case "nis":
		if existing, _ := s.repo.FindByNIS(ctx, value); existing != nil && existing.ID != excludeID {
			return false, nil
		}
	case "email":
		value = strings.ToLower(value)
		if existing, _ := s.repo.FindByEmail(ctx, value); existing != nil && existing.ID != excludeID {
			return false, nil
		}
		if existing, _ := s.userRepo.FindByEmail(ctx, value); existing != nil {
			return false, nil
		}
	case "phone_number":
		normalized := utils.NormalizePhoneNumber(value)
		if existing, _ := s.repo.FindByPhone(ctx, normalized); existing != nil && existing.ID != excludeID {
			return false, nil
		}
		if existing, _ := s.userRepo.FindByPhone(ctx, normalized); existing != nil {
			return false, nil
		}
	}

	return true, nil
}
