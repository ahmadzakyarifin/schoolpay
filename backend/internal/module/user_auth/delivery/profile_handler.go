package delivery

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/dto"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type ProfileHandler struct {
	db          *bun.DB
	userRepo    repository.UserRepo
	studentRepo academicrepo.StudentRepo
	audit       auditusecase.AuditLogService
}

func NewProfileHandler(db *bun.DB, userRepo repository.UserRepo, studentRepo academicrepo.StudentRepo, audit auditusecase.AuditLogService) *ProfileHandler {
	return &ProfileHandler{
		db:          db,
		userRepo:    userRepo,
		studentRepo: studentRepo,
		audit:       audit,
	}
}

type updateOwnProfileRequest struct {
	Name        string `json:"name" binding:"required,min=2"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required,min=9,max=20"`
	NIK         string `json:"nik"`
	BirthDate   string `json:"birth_date" binding:"omitempty,custom_date"`
	Address     string `json:"address"`
	Education   string `json:"education"`
	Occupation  string `json:"occupation"`
	Income      string `json:"income"`
}

func (h *ProfileHandler) Me(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "session tidak valid")
		return
	}

	user, err := h.userRepo.FindByID(c.Request.Context(), userID.(uint))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusNotFound, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "berhasil", dto.ToUserResponse(*user))
}

func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	var req updateOwnProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}

	userID, ok := c.Get("user_id")
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "session tidak valid")
		return
	}

	existing, err := h.userRepo.FindByID(c.Request.Context(), userID.(uint))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusNotFound, err)
		return
	}

	updated := *existing
	updated.Name = strings.TrimSpace(req.Name)
	updated.Email = strings.ToLower(strings.TrimSpace(req.Email))
	updated.PhoneNumber = utils.NormalizePhoneNumber(req.PhoneNumber)
	updated.NIK = nilIfBlank(req.NIK)
	updated.Address = nilIfBlank(req.Address)
	updated.Education = nilIfBlank(req.Education)
	updated.Occupation = nilIfBlank(req.Occupation)
	updated.Income = nilIfBlank(req.Income)
	if req.BirthDate != "" {
		parsedDate, _ := time.Parse("02/01/2006", req.BirthDate)
		updated.BirthDate = &parsedDate
	}
	updated.UpdatedAt = time.Now()

	bizErrors := make(map[string][]string)
	if !utils.ValidatePhoneNumber(updated.PhoneNumber) {
		bizErrors["phone_number"] = append(bizErrors["phone_number"], "Nomor WhatsApp tidak valid.")
	}
	if updated.NIK != nil && !isDigits(*updated.NIK, 16) {
		bizErrors["nik"] = append(bizErrors["nik"], "NIK harus 16 digit angka.")
	}

	if other, _ := h.userRepo.FindByEmail(c.Request.Context(), updated.Email); other != nil && other.ID != updated.ID {
		bizErrors["email"] = append(bizErrors["email"], fmt.Sprintf("Email '%s' sudah digunakan pengguna lain.", updated.Email))
	} else if otherStudent, _ := h.studentRepo.FindByEmail(c.Request.Context(), updated.Email); otherStudent != nil {
		bizErrors["email"] = append(bizErrors["email"], fmt.Sprintf("Email '%s' sudah digunakan siswa.", updated.Email))
	}

	if other, _ := h.userRepo.FindByPhone(c.Request.Context(), updated.PhoneNumber); other != nil && other.ID != updated.ID {
		bizErrors["phone_number"] = append(bizErrors["phone_number"], fmt.Sprintf("Nomor WhatsApp '%s' sudah digunakan pengguna lain.", updated.PhoneNumber))
	} else if otherStudent, _ := h.studentRepo.FindByPhone(c.Request.Context(), updated.PhoneNumber); otherStudent != nil {
		bizErrors["phone_number"] = append(bizErrors["phone_number"], fmt.Sprintf("Nomor WhatsApp '%s' sudah digunakan siswa.", updated.PhoneNumber))
	}

	if updated.NIK != nil {
		if other, _ := h.userRepo.FindByNIK(c.Request.Context(), *updated.NIK); other != nil && other.ID != updated.ID {
			bizErrors["nik"] = append(bizErrors["nik"], fmt.Sprintf("NIK '%s' sudah digunakan pengguna lain.", *updated.NIK))
		} else if otherStudent, _ := h.studentRepo.FindByNIK(c.Request.Context(), *updated.NIK); otherStudent != nil {
			bizErrors["nik"] = append(bizErrors["nik"], fmt.Sprintf("NIK '%s' sudah digunakan siswa.", *updated.NIK))
		}
	}

	if len(bizErrors) > 0 {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi bisnis gagal", bizErrors)
		return
	}

	if err := h.userRepo.Update(c.Request.Context(), h.db, &updated); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	if h.audit != nil {
		auditUserID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(c.Request.Context())
		_ = h.audit.Log(c.Request.Context(), h.db, auditUserID, userName, role, "UPDATE_OWN_PROFILE", "users", updated.ID, profileAuditValues(existing), profileAuditValues(&updated), ipAddress, userAgent)
	}

	user, _ := h.userRepo.FindByID(c.Request.Context(), updated.ID)
	helper.SuccessResponse(c, http.StatusOK, "profil berhasil diperbarui", dto.ToUserResponse(*user))
}

func (h *ProfileHandler) UploadPhoto(c *gin.Context) {
	userID, ok := c.Get("user_id")
	if !ok {
		helper.ErrorResponse(c, http.StatusUnauthorized, "session tidak valid")
		return
	}

	existing, err := h.userRepo.FindByID(c.Request.Context(), userID.(uint))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusNotFound, err)
		return
	}

	file, err := c.FormFile("photo")
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "file foto tidak ditemukan")
		return
	}
	if file.Size > 2*1024*1024 {
		helper.ErrorResponse(c, http.StatusBadRequest, "ukuran foto maksimal 2MB")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".webp" {
		helper.ErrorResponse(c, http.StatusBadRequest, "format foto harus jpg, jpeg, png, atau webp")
		return
	}

	uploadDir := filepath.Join("public", "uploads", "users")
	if err := utils.EnsureDir(uploadDir); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "gagal menyiapkan folder upload")
		return
	}

	filename := fmt.Sprintf("user-%d-%d%s", existing.ID, time.Now().UnixNano(), ext)
	savePath := filepath.Join(uploadDir, filename)
	if err := c.SaveUploadedFile(file, savePath); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "gagal menyimpan foto")
		return
	}

	updated := *existing
	publicPath := "/uploads/users/" + filename
	updated.ImagePath = &publicPath
	updated.UpdatedAt = time.Now()
	if err := h.userRepo.Update(c.Request.Context(), h.db, &updated); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	if h.audit != nil {
		auditUserID, userName, role, ipAddress, userAgent := helper.GetAuditMeta(c.Request.Context())
		oldVals := map[string]interface{}{"image_path": existing.ImagePath}
		newVals := map[string]interface{}{"image_path": updated.ImagePath}
		_ = h.audit.Log(c.Request.Context(), h.db, auditUserID, userName, role, "UPDATE_OWN_PROFILE_PHOTO", "users", updated.ID, oldVals, newVals, ipAddress, userAgent)
	}

	user, _ := h.userRepo.FindByID(c.Request.Context(), updated.ID)
	helper.SuccessResponse(c, http.StatusOK, "foto profil berhasil diperbarui", dto.ToUserResponse(*user))
}

func nilIfBlank(value string) *string {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	return &value
}

func isDigits(value string, size int) bool {
	if len(value) != size {
		return false
	}
	for _, r := range value {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

func profileAuditValues(user *domain.User) map[string]interface{} {
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
		"nik":          user.NIK,
		"birth_date":   birthDate,
		"address":      user.Address,
		"education":    user.Education,
		"occupation":   user.Occupation,
		"income":       user.Income,
		"image_path":   user.ImagePath,
	}
}
