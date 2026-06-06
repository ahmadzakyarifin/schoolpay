package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/config"
	academicdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	notificationdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/domain"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/uptrace/bun"
	"github.com/xuri/excelize/v2"
)

type UserService interface {
	GetAll(ctx context.Context, role string) ([]domain.User, error)
	GetPaginated(ctx context.Context, page, limit int, search, role, filter, status, sort string) ([]domain.User, int, error)
	Create(ctx context.Context, user *domain.User, cfg *config.Config) error
	Update(ctx context.Context, user *domain.User, cfg *config.Config) error
	Delete(ctx context.Context, id uint) error
	ActivateAccount(ctx context.Context, token string, password string) (*domain.User, error)
	SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error
	ToggleStatus(ctx context.Context, id uint) error
	ResendNotification(ctx context.Context, id uint, channel string, cfg *config.Config) error
	BulkResendNotification(ctx context.Context, ids []uint, channel string, cfg *config.Config) (*BulkResendResult, error)
	ExportExcel(ctx context.Context, search, role, filter, status string) ([]byte, error)
	GetNotifications(ctx context.Context, userID uint) ([]notificationdomain.Notification, error)
	GetByID(ctx context.Context, id uint) (*domain.User, error)
	BulkDelete(ctx context.Context, ids []uint) error
	Restore(ctx context.Context, id uint) error
	BulkRestore(ctx context.Context, ids []uint) error
	GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error)
	CheckUnique(ctx context.Context, field string, value string, excludeID uint) (bool, error)
}

type userService struct {
	db          *bun.DB
	repo        repository.UserRepo
	authRepo    repository.AuthRepo
	messenger   utils.Messenger
	notiRepo    notificationrepo.NotificationRepo
	studentRepo academicrepo.StudentRepo
	audit       auditusecase.AuditLogService
}

func NewUserService(db *bun.DB, repo repository.UserRepo, authRepo repository.AuthRepo, msg utils.Messenger, noti notificationrepo.NotificationRepo, stuRepo academicrepo.StudentRepo, audit auditusecase.AuditLogService) UserService {
	return &userService{
		db:          db,
		repo:        repo,
		authRepo:    authRepo,
		messenger:   msg,
		notiRepo:    noti,
		studentRepo: stuRepo,
		audit:       audit,
	}
}

func (s *userService) GetAll(ctx context.Context, role string) ([]domain.User, error) {
	return s.repo.FindAll(ctx, role)
}

func (s *userService) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	return s.repo.FindByID(ctx, id)
}

func (s *userService) GetPaginated(ctx context.Context, page, limit int, search, role, filter, status, sort string) ([]domain.User, int, error) {
	return s.repo.FindPaginated(ctx, page, limit, search, role, filter, status, sort)
}

func userAuditValues(user *domain.User) map[string]interface{} {
	if user == nil {
		return nil
	}

	var birthDate interface{}
	if user.BirthDate != nil && !user.BirthDate.IsZero() {
		birthDate = user.BirthDate.Format("2006-01-02")
	}

	return map[string]interface{}{
		"name":         user.Name,
		"email":        user.Email,
		"phone_number": user.PhoneNumber,
		"role":         user.Role,
		"nik":          user.NIK,
		"birth_date":   birthDate,
		"address":      user.Address,
		"education":    user.Education,
		"occupation":   user.Occupation,
		"income":       user.Income,
		"is_active":    user.IsActive,
	}
}

func (s *userService) Create(ctx context.Context, user *domain.User, cfg *config.Config) error {
	// Normalization & Trimming
	user.Name = strings.Title(strings.ToLower(strings.TrimSpace(user.Name)))
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	user.PhoneNumber = utils.NormalizePhoneNumber(user.PhoneNumber)

	if user.NIK != nil {
		*user.NIK = strings.TrimSpace(*user.NIK)
	}
	if user.Address != nil {
		*user.Address = strings.TrimSpace(*user.Address)
	}
	if user.Education != nil {
		*user.Education = strings.TrimSpace(*user.Education)
	}
	if user.Occupation != nil {
		*user.Occupation = strings.TrimSpace(*user.Occupation)
	}
	if user.Income != nil {
		*user.Income = strings.TrimSpace(*user.Income)
	}

	return s.db.RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		bizErrors := make(map[string][]string)

		if existing, _ := s.repo.FindByEmail(ctx, user.Email); existing != nil {
			bizErrors["email"] = append(bizErrors["email"], fmt.Sprintf("Email '%s' sudah terdaftar sebagai Pengguna/Wali Murid", user.Email))
		} else if existing, _ := s.studentRepo.FindByEmail(ctx, user.Email); existing != nil {
			bizErrors["email"] = append(bizErrors["email"], fmt.Sprintf("Email '%s' sudah terdaftar sebagai Siswa", user.Email))
		}

		if existing, _ := s.repo.FindByPhone(ctx, user.PhoneNumber); existing != nil {
			bizErrors["phone_number"] = append(bizErrors["phone_number"], fmt.Sprintf("Nomor WhatsApp '%s' sudah terdaftar sebagai Pengguna/Wali Murid", user.PhoneNumber))
		} else if existing, _ := s.studentRepo.FindByPhone(ctx, user.PhoneNumber); existing != nil {
			bizErrors["phone_number"] = append(bizErrors["phone_number"], fmt.Sprintf("Nomor WhatsApp '%s' sudah terdaftar sebagai Siswa", user.PhoneNumber))
		}

		if user.NIK != nil && *user.NIK != "" {
			if existing, _ := s.repo.FindByNIK(ctx, *user.NIK); existing != nil {
				bizErrors["nik"] = append(bizErrors["nik"], fmt.Sprintf("NIK '%s' sudah terdaftar sebagai Pengguna/Wali Murid", *user.NIK))
			} else if existing, _ := s.studentRepo.FindByNIK(ctx, *user.NIK); existing != nil {
				bizErrors["nik"] = append(bizErrors["nik"], fmt.Sprintf("NIK '%s' sudah terdaftar sebagai Siswa", *user.NIK))
			}
		}

		if len(bizErrors) > 0 {
			return utils.NewBusinessMultiError(bizErrors)
		}

		if err := s.repo.Create(ctx, tx, user); err != nil {
			return err
		}

		if strings.EqualFold(user.Role, "parent") {
			for _, sid := range user.StudentIDs {
				_, _ = tx.NewUpdate().Model((*academicdomain.Student)(nil)).Set("parent_id = ?", user.ID).Where("id = ?", sid).Exec(ctx)
			}
		}

		if s.audit != nil {
			userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, s.db, userID, userName, role, "CREATE", "users", user.ID, nil, userAuditValues(user), ipAddress, userAgent)
		}

		if user.IsActive {
			payloadBytes, _ := json.Marshal(map[string]interface{}{"user_id": user.ID, "channel": ""})
			job := &notificationdomain.BackgroundJob{
				TaskName: "auth:user_activation",
				Payload:  string(payloadBytes),
				Status:   "pending",
			}
			_, _ = tx.NewInsert().Model(job).Exec(ctx)
		}
		return nil
	})
}

func (s *userService) Update(ctx context.Context, user *domain.User, cfg *config.Config) error {
	// Normalization & Trimming
	user.Name = strings.Title(strings.ToLower(strings.TrimSpace(user.Name)))
	user.Email = strings.ToLower(strings.TrimSpace(user.Email))

	user.PhoneNumber = utils.NormalizePhoneNumber(user.PhoneNumber)

	if user.NIK != nil {
		*user.NIK = strings.TrimSpace(*user.NIK)
	}
	if user.Address != nil {
		*user.Address = strings.TrimSpace(*user.Address)
	}
	if user.Education != nil {
		*user.Education = strings.TrimSpace(*user.Education)
	}
	if user.Occupation != nil {
		*user.Occupation = strings.TrimSpace(*user.Occupation)
	}
	if user.Income != nil {
		*user.Income = strings.TrimSpace(*user.Income)
	}

	// Business Logic Validation (Batch Check)
	bizErrors := make(map[string][]string)

	if existing, _ := s.repo.FindByEmail(ctx, user.Email); existing != nil && existing.ID != user.ID {
		bizErrors["email"] = append(bizErrors["email"], fmt.Sprintf("Email '%s' sudah digunakan oleh Pengguna lain", user.Email))
	} else if existing, _ := s.studentRepo.FindByEmail(ctx, user.Email); existing != nil {
		bizErrors["email"] = append(bizErrors["email"], fmt.Sprintf("Email '%s' sudah digunakan oleh Siswa lain", user.Email))
	}

	if existing, _ := s.repo.FindByPhone(ctx, user.PhoneNumber); existing != nil && existing.ID != user.ID {
		bizErrors["phone_number"] = append(bizErrors["phone_number"], fmt.Sprintf("Nomor WhatsApp '%s' sudah digunakan oleh Pengguna lain", user.PhoneNumber))
	} else if existing, _ := s.studentRepo.FindByPhone(ctx, user.PhoneNumber); existing != nil {
		bizErrors["phone_number"] = append(bizErrors["phone_number"], fmt.Sprintf("Nomor WhatsApp '%s' sudah digunakan oleh Siswa lain", user.PhoneNumber))
	}

	if user.NIK != nil && *user.NIK != "" {
		if existing, _ := s.repo.FindByNIK(ctx, *user.NIK); existing != nil && existing.ID != user.ID {
			bizErrors["nik"] = append(bizErrors["nik"], fmt.Sprintf("NIK '%s' sudah digunakan oleh Pengguna lain", *user.NIK))
		} else if existing, _ := s.studentRepo.FindByNIK(ctx, *user.NIK); existing != nil {
			bizErrors["nik"] = append(bizErrors["nik"], fmt.Sprintf("NIK '%s' sudah digunakan oleh Siswa lain", *user.NIK))
		}
	}

	if len(bizErrors) > 0 {
		return utils.NewBusinessMultiError(bizErrors)
	}

	existing, _ := s.repo.FindByID(ctx, user.ID)

	if err := s.repo.Update(ctx, s.db, user); err != nil {
		return err
	}

	if existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "UPDATE", "users", user.ID, userAuditValues(existing), userAuditValues(user), ipAddress, userAgent)
	}

	return nil
}

func (s *userService) ensureDeletableUser(user *domain.User) error {
	if user == nil {
		return errors.New("pengguna tidak ditemukan")
	}

	if strings.EqualFold(user.Role, "parent") && user.StudentCount > 0 {
		return fmt.Errorf("wali murid %s tidak dapat dihapus karena masih terhubung dengan %d siswa aktif. Pindahkan atau lepaskan relasi siswa terlebih dahulu", user.Name, user.StudentCount)
	}

	return nil
}

func (s *userService) Delete(ctx context.Context, id uint) error {
	existing, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if err := s.ensureDeletableUser(existing); err != nil {
		return err
	}

	err = s.repo.Delete(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"is_active": existing.IsActive, "status": "active"}
		newVals := map[string]interface{}{"is_active": false, "status": "deleted"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "DELETE", "users", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *userService) BulkDelete(ctx context.Context, ids []uint) error {
	for _, id := range ids {
		existing, err := s.repo.FindByID(ctx, id)
		if err != nil {
			return err
		}
		if err := s.ensureDeletableUser(existing); err != nil {
			return err
		}
	}

	err := s.repo.BulkDelete(ctx, ids)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		for _, id := range ids {
			oldVals := map[string]interface{}{"status": "active"}
			newVals := map[string]interface{}{"status": "deleted"}
			_ = s.audit.Log(ctx, s.db, userID, userName, role, "BULK_DELETE", "users", id, oldVals, newVals, ipAddress, userAgent)
		}
	}
	return err
}

func (s *userService) Restore(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	err := s.repo.Restore(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"is_active": existing.IsActive, "status": "deleted"}
		newVals := map[string]interface{}{"is_active": true, "status": "active"}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "RESTORE", "users", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *userService) BulkRestore(ctx context.Context, ids []uint) error {
	err := s.repo.BulkRestore(ctx, ids)
	if err == nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		for _, id := range ids {
			oldVals := map[string]interface{}{"status": "deleted"}
			newVals := map[string]interface{}{"status": "active"}
			_ = s.audit.Log(ctx, s.db, userID, userName, role, "BULK_RESTORE", "users", id, oldVals, newVals, ipAddress, userAgent)
		}
	}
	return err
}

func (s *userService) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	return s.authRepo.SaveRefreshToken(ctx, userID, token, expiresAt)
}

func (s *userService) ActivateAccount(ctx context.Context, token string, password string) (*domain.User, error) {
	userID, err := s.authRepo.FindAuthToken(ctx, token, "activation")
	if err != nil {
		return nil, errors.New("link aktivasi tidak valid atau sudah kedaluwarsa")
	}

	hashed, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("gagal memproses password")
	}
	_ = s.repo.UpdatePassword(ctx, userID, hashed)

	_ = s.authRepo.DeleteAuthToken(ctx, token, "activation")
	user, err := s.repo.FindByID(ctx, userID)
	if err == nil && user != nil && s.audit != nil {
		userIDMeta, userNameMeta, roleMeta, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.db, userIDMeta, userNameMeta, roleMeta, "ACTIVATE_ACCOUNT", "users", user.ID, nil, map[string]interface{}{"is_active": true}, ipAddress, userAgent)
	}
	return user, err
}

func (s *userService) ToggleStatus(ctx context.Context, id uint) error {
	existing, _ := s.repo.FindByID(ctx, id)
	err := s.repo.ToggleStatus(ctx, id)
	if err == nil && existing != nil && s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		oldVals := map[string]interface{}{"is_active": existing.IsActive}
		newVals := map[string]interface{}{"is_active": !existing.IsActive}
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "TOGGLE_STATUS", "users", id, oldVals, newVals, ipAddress, userAgent)
	}
	return err
}

func (s *userService) ResendNotification(ctx context.Context, id uint, channel string, cfg *config.Config) error {
	channel = normalizeActivationChannel(channel)
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// Validasi Bersama
	if !user.IsActive {
		return fmt.Errorf("Gagal: Akun %s sedang Non-Aktif", user.Name)
	}
	if user.Role == "parent" && user.StudentCount == 0 {
		return fmt.Errorf("Gagal: Wali Murid %s belum terhubung ke siswa manapun", user.Name)
	}
	if userHasPassword(user) {
		return fmt.Errorf("Gagal: Akun %s sudah aktif dan sudah memiliki password. Gunakan fitur reset password jika diperlukan", user.Name)
	}

	if s.audit != nil {
		userID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.db, userID, userName, role, "RESEND_ACTIVATION", "users", id, nil, map[string]interface{}{"channel": channel}, ipAddress, userAgent)
	}

	payloadBytes, _ := json.Marshal(map[string]interface{}{"user_id": user.ID, "channel": channel})
	job := &notificationdomain.BackgroundJob{
		TaskName: "auth:user_activation",
		Payload:  string(payloadBytes),
		Status:   "pending",
	}
	_, _ = s.db.NewInsert().Model(job).Exec(ctx)
	return nil
}

type BulkResendResult struct {
	Total  int      `json:"total"`
	Sent   int      `json:"sent"`
	Failed int      `json:"failed"`
	Errors []string `json:"errors"`
}

func (s *userService) BulkResendNotification(ctx context.Context, ids []uint, channel string, cfg *config.Config) (*BulkResendResult, error) {
	channel = normalizeActivationChannel(channel)
	result := &BulkResendResult{
		Total:  len(ids),
		Errors: []string{},
	}

	for _, id := range ids {
		user, err := s.repo.FindByID(ctx, id)
		if err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("ID %d: Data tidak ditemukan", id))
			continue
		}

		// Validasi Instan
		if !user.IsActive {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: Akun sedang Non-Aktif", user.Name))
			continue
		}
		if user.Role == "parent" && user.StudentCount == 0 {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: Wali Murid belum memiliki data anak", user.Name))
			continue
		}
		if userHasPassword(user) {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("%s: Akun sudah aktif dan sudah memiliki password", user.Name))
			continue
		}

		// Jika lolos, kirim ke job queue
		result.Sent++
		if s.audit != nil {
			userIDMeta, userNameMeta, roleMeta, ipAddress, userAgent := helper.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, s.db, userIDMeta, userNameMeta, roleMeta, "BULK_RESEND_ACTIVATION", "users", id, nil, map[string]interface{}{"channel": channel}, ipAddress, userAgent)
		}
		payloadBytes, _ := json.Marshal(map[string]interface{}{"user_id": user.ID, "channel": channel})
		job := &notificationdomain.BackgroundJob{
			TaskName: "auth:user_activation",
			Payload:  string(payloadBytes),
			Status:   "pending",
		}
		_, _ = s.db.NewInsert().Model(job).Exec(ctx)
	}

	return result, nil
}

func userHasPassword(user *domain.User) bool {
	return user != nil && strings.TrimSpace(user.PasswordHash) != ""
}

func normalizeActivationChannel(channel string) string {
	switch strings.ToLower(strings.TrimSpace(channel)) {
	case "email":
		return "email"
	case "whatsapp", "wa":
		return "whatsapp"
	default:
		return ""
	}
}

func (s *userService) ExportExcel(ctx context.Context, search, role, filter, status string) ([]byte, error) {
	if s.audit != nil {
		userIDMeta, userNameMeta, roleMeta, ipAddress, userAgent := helper.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.db, userIDMeta, userNameMeta, roleMeta, "EXPORT_EXCEL", "users", 0, nil, map[string]interface{}{"search": search, "role": role, "filter": filter, "status": status}, ipAddress, userAgent)
	}
	users, _, _ := s.repo.FindPaginated(ctx, 1, 1000000, search, role, filter, status, "")
	f := excelize.NewFile()
	sheet := "Daftar Pengguna"
	f.SetSheetName("Sheet1", sheet)

	// Premium Headers with Border — aligned with import template fields
	headers := []string{"NAMA LENGKAP", "EMAIL", "NOMOR WHATSAPP", "ROLE", "NIK", "TANGGAL LAHIR", "ALAMAT", "PENDIDIKAN", "PEKERJAAN", "PENGHASILAN", "JML ANAK", "STATUS"}
	colWidths := []float64{28, 28, 20, 12, 20, 20, 30, 22, 25, 16, 10, 12}
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
			f.SetColWidth(sheet, colName, colName, 20)
		}
	}

	// Body Style with Border
	bodyStyle, _ := f.NewStyle(&excelize.Style{
		Border: []excelize.Border{
			{Type: "left", Color: "000000", Style: 1},
			{Type: "top", Color: "000000", Style: 1},
			{Type: "bottom", Color: "000000", Style: 1},
			{Type: "right", Color: "000000", Style: 1},
		},
		Alignment: &excelize.Alignment{Vertical: "center"},
	})

	for i, u := range users {
		row := i + 2
		activeStatus := "Aktif"
		if !u.IsActive {
			activeStatus = "Non-Aktif"
		}

		nik := utils.StringValue(u.NIK)
		if nik == "" {
			nik = "-"
		}

		var birthDate string
		if u.BirthDate != nil && !u.BirthDate.IsZero() && u.Role != "admin" {
			birthDate = u.BirthDate.Format("02/01/2006")
		} else {
			birthDate = "-"
		}

		address := utils.StringValue(u.Address)
		if address == "" {
			address = "-"
		}

		education := utils.StringValue(u.Education)
		if education == "" {
			education = "-"
		}

		occupation := utils.StringValue(u.Occupation)
		if occupation == "" {
			occupation = "-"
		}

		income := utils.StringValue(u.Income)
		if income == "" {
			income = "-"
		}

		f.SetCellValue(sheet, fmt.Sprintf("A%d", row), u.Name)
		f.SetCellValue(sheet, fmt.Sprintf("B%d", row), u.Email)
		f.SetCellValue(sheet, fmt.Sprintf("C%d", row), u.PhoneNumber)
		f.SetCellValue(sheet, fmt.Sprintf("D%d", row), strings.ToUpper(u.Role))
		f.SetCellValue(sheet, fmt.Sprintf("E%d", row), nik)
		f.SetCellValue(sheet, fmt.Sprintf("F%d", row), birthDate)
		f.SetCellValue(sheet, fmt.Sprintf("G%d", row), address)
		f.SetCellValue(sheet, fmt.Sprintf("H%d", row), education)
		f.SetCellValue(sheet, fmt.Sprintf("I%d", row), occupation)
		f.SetCellValue(sheet, fmt.Sprintf("J%d", row), income)
		f.SetCellValue(sheet, fmt.Sprintf("K%d", row), u.StudentCount)
		f.SetCellValue(sheet, fmt.Sprintf("L%d", row), activeStatus)

		// Apply body style to the entire row (A-L)
		f.SetCellStyle(sheet, fmt.Sprintf("A%d", row), fmt.Sprintf("L%d", row), bodyStyle)
	}

	buf, _ := f.WriteToBuffer()
	return buf.Bytes(), nil
}

func (s *userService) GetNotifications(ctx context.Context, userID uint) ([]notificationdomain.Notification, error) {
	return s.notiRepo.GetByUserID(ctx, userID)
}

func FormatActivationStudentList(raw string) string {
	parts := strings.Split(raw, "||")
	lines := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		if split := strings.SplitN(part, "::", 2); len(split) == 2 {
			part = split[1]
		}
		lines = append(lines, fmt.Sprintf("• *%s*", part))
	}
	if len(lines) == 0 {
		return "• *Data siswa terhubung*"
	}
	return strings.Join(lines, "\n")
}

func (s *userService) GetDependencyInfo(ctx context.Context, id uint) (map[string]interface{}, error) {
	user, err := s.repo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	var messages []string
	if user.Role == "parent" && user.StudentCount > 0 {
		messages = append(messages, fmt.Sprintf("terhubung dengan %d siswa", user.StudentCount))
	}

	return map[string]interface{}{
		"has_dependencies": len(messages) > 0,
		"message":          strings.Join(messages, " dan "),
		"counts": map[string]int{
			"students": user.StudentCount,
		},
	}, nil
}

func (s *userService) CheckUnique(ctx context.Context, field string, value string, excludeID uint) (bool, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return true, nil
	}

	switch field {
	case "nik":
		if existing, _ := s.repo.FindByNIK(ctx, value); existing != nil && existing.ID != excludeID {
			return false, nil
		}
		if existing, _ := s.studentRepo.FindByNIK(ctx, value); existing != nil {
			return false, nil
		}
	case "email":
		value = strings.ToLower(value)
		if existing, _ := s.repo.FindByEmail(ctx, value); existing != nil && existing.ID != excludeID {
			return false, nil
		}
		if existing, _ := s.studentRepo.FindByEmail(ctx, value); existing != nil {
			return false, nil
		}
	case "phone_number":
		normalized := utils.NormalizePhoneNumber(value)
		if existing, _ := s.repo.FindByPhone(ctx, normalized); existing != nil && existing.ID != excludeID {
			return false, nil
		}
		if existing, _ := s.studentRepo.FindByPhone(ctx, normalized); existing != nil {
			return false, nil
		}
	}

	return true, nil
}
