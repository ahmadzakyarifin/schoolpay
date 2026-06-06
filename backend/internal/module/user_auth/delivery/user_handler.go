package delivery

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/dto"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/gin-gonic/gin"
)

type userRequest struct {
	Name        string `json:"name" binding:"required,min=2"`
	Email       string `json:"email" binding:"required,email"`
	PhoneNumber string `json:"phone_number" binding:"required,min=9,max=15"`
	Role        string `json:"role" binding:"required,oneof=admin parent"`
	NIK         string `json:"nik" binding:"omitempty,numeric,len=16"`
	BirthDate   string `json:"birth_date" binding:"omitempty,custom_date"`
	Address     string `json:"address" binding:"omitempty"`
	Education   string `json:"education" binding:"omitempty"`
	Occupation  string `json:"occupation" binding:"omitempty"`
	Income      string `json:"income" binding:"omitempty"`
	IsActive    *bool  `json:"is_active" binding:"omitempty"`
}

func (req *userRequest) ToDomain() domain.User {
	var birthDate *time.Time
	if req.BirthDate != "" {
		if parsedDate, err := time.Parse("02/01/2006", req.BirthDate); err == nil {
			birthDate = &parsedDate
		}
	}

	var isActive bool = true
	if req.IsActive != nil {
		isActive = *req.IsActive
	}

	var nik *string
	if req.NIK != "" {
		nik = &req.NIK
	}

	var address *string
	if req.Address != "" {
		address = &req.Address
	}

	var education *string
	if req.Education != "" {
		education = &req.Education
	}

	var occupation *string
	if req.Occupation != "" {
		occupation = &req.Occupation
	}

	var income *string
	if req.Income != "" {
		income = &req.Income
	}

	return domain.User{
		Name:        req.Name,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Role:        req.Role,
		NIK:         nik,
		BirthDate:   birthDate,
		Address:     address,
		Education:   education,
		Occupation:  occupation,
		Income:      income,
		IsActive:    isActive,
	}
}

type UserHandler struct {
	s   usecase.UserService
	cfg *config.Config
}

func NewUserHandler(s usecase.UserService, cfg *config.Config) *UserHandler {
	return &UserHandler{s: s, cfg: cfg}
}

func (h *UserHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	filter := c.Query("filter")
	role := c.Query("role")
	status := c.Query("status")
	sort := c.Query("sort")

	users, total, err := h.s.GetPaginated(c.Request.Context(), page, limit, search, role, filter, status, sort)
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "berhasil mengambil data", gin.H{
		"users": dto.ToUserResponseList(users),
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *UserHandler) GetByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := h.s.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "user tidak ditemukan")
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil mengambil data user", dto.ToUserResponse(*user))
}

func (h *UserHandler) Create(c *gin.Context) {
	var req userRequest
	allErrors := make(map[string][]string)

	if err := c.ShouldBindJSON(&req); err != nil {
		allErrors = helper.GetValidationErrors(err)
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", allErrors)
		return
	}

	user := req.ToDomain()

	if err := h.s.Create(c.Request.Context(), &user, h.cfg); err != nil {
		if bmErr, ok := err.(*utils.BusinessMultiError); ok {
			for k, v := range bmErr.Errors {
				allErrors[k] = append(allErrors[k], v...)
			}
		} else if bErr, ok := err.(*utils.BusinessError); ok {
			allErrors[bErr.Field] = append(allErrors[bErr.Field], bErr.Message)
		} else {
			// If it's a real system error (not validation), return 500
			if len(allErrors) == 0 {
				helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
				return
			}
		}
	}

	//  Return All Errors if any
	if len(allErrors) > 0 {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", allErrors)
		return
	}

	helper.SuccessResponse(c, http.StatusCreated, "user berhasil dibuat, notifikasi dikirim", dto.ToUserResponse(user))
}

func (h *UserHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req userRequest
	allErrors := make(map[string][]string)

	if err := c.ShouldBindJSON(&req); err != nil {
		allErrors = helper.GetValidationErrors(err)
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", allErrors)
		return
	}
	user := req.ToDomain()
	user.ID = uint(id)

	if err := h.s.Update(c.Request.Context(), &user, h.cfg); err != nil {
		if bmErr, ok := err.(*utils.BusinessMultiError); ok {
			for k, v := range bmErr.Errors {
				allErrors[k] = append(allErrors[k], v...)
			}
		} else if bErr, ok := err.(*utils.BusinessError); ok {
			allErrors[bErr.Field] = append(allErrors[bErr.Field], bErr.Message)
		} else {
			if len(allErrors) == 0 {
				helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
				return
			}
		}
	}

	// Return All if any
	if len(allErrors) > 0 {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", allErrors)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "user berhasil diperbarui", dto.ToUserResponse(user))
}

func (h *UserHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 1. Proteksi Diri Sendiri
	currentUserID, _ := c.Get("user_id")
	if uint(id) == currentUserID.(uint) {
		helper.ErrorResponse(c, http.StatusBadRequest, "Anda tidak diperbolehkan menghapus akun Anda sendiri")
		return
	}

	// 2. Sinkronisasi FE: Role Admin TIDAK BOLEH dihapus
	user, err := h.s.GetByID(c.Request.Context(), uint(id))
	if err != nil {
		helper.ErrorResponse(c, http.StatusNotFound, "user tidak ditemukan")
		return
	}
	if user.Role == "admin" {
		helper.ErrorResponse(c, http.StatusBadRequest, "Akun dengan role Admin tidak dapat dihapus demi alasan keamanan data")
		return
	}

	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "user berhasil dihapus", nil)
}

func (h *UserHandler) ToggleStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	// 1. Cek Apakah User Eksis
	user, err := h.s.GetByID(c.Request.Context(), uint(id))
	if err != nil || user == nil {
		helper.ErrorResponse(c, http.StatusNotFound, "user tidak ditemukan")
		return
	}

	// 2. Proteksi Diri Sendiri
	currentUserID, _ := c.Get("user_id")
	if uint(id) == currentUserID.(uint) {
		helper.ErrorResponse(c, http.StatusBadRequest, "Anda tidak diperbolehkan menonaktifkan akun Anda sendiri")
		return
	}

	// 3. Proteksi Admin Terakhir (Safety Measure)
	if user.Role == "admin" && user.IsActive {
		// Jika ini admin dan mau di-nonaktifkan, cek jumlah admin aktif lainnya
		admins, _ := h.s.GetAll(c.Request.Context(), "admin")
		activeAdmins := 0
		for _, a := range admins {
			if a.IsActive {
				activeAdmins++
			}
		}
		if activeAdmins <= 1 {
			helper.ErrorResponse(c, http.StatusBadRequest, "Gagal: Ini adalah satu-satunya akun Admin yang aktif di sistem")
			return
		}
	}

	if err := h.s.ToggleStatus(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "status user berhasil diubah", nil)
}

func (h *UserHandler) Activate(c *gin.Context) {
	var req struct {
		Token    string `json:"token" binding:"required"`
		Password string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}

	ctx := c.Request.Context()
	user, err := h.s.ActivateAccount(ctx, req.Token, req.Password)
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	accessToken, _ := utils.GenerateAccessToken(user.ID, user.Email, user.Role, h.cfg.JWTSecret)
	refreshToken, expiry, _ := utils.GenerateRefreshToken()
	if err := h.s.SaveRefreshToken(ctx, user.ID, refreshToken, expiry); err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "gagal menyimpan session")
		return
	}

	// Set Refresh Token in HttpOnly Cookie
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("refresh_token", refreshToken, int(time.Until(expiry).Seconds()), "/", "", false, true)

	helper.SuccessResponse(c, http.StatusOK, "akun berhasil diaktifkan", gin.H{
		"access_token": accessToken,
		"user": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (h *UserHandler) ResendNotification(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var req struct {
		Channel string `json:"channel"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		// If JSON fails, fallback to query for backward compatibility
		req.Channel = c.Query("channel")
	}

	if err := h.s.ResendNotification(c.Request.Context(), uint(id), req.Channel, h.cfg); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "link aktivasi berhasil dikirim ulang", nil)
}

func (h *UserHandler) BulkResendNotification(c *gin.Context) {
	var req struct {
		IDs     []uint `json:"ids" binding:"required"`
		Channel string `json:"channel"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}

	result, err := h.s.BulkResendNotification(c.Request.Context(), req.IDs, req.Channel, h.cfg)
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, fmt.Sprintf("%d link aktivasi diproses", result.Sent), result)
}

func (h *UserHandler) BulkDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	currentUserID, _ := c.Get("user_id")
	for _, id := range req.IDs {
		// 1. Proteksi Diri Sendiri
		if id == currentUserID.(uint) {
			helper.ErrorResponse(c, http.StatusBadRequest, "Operasi dibatalkan: Terdapat akun Anda sendiri dalam daftar hapus")
			return
		}

		// 2. Proteksi Role Admin (Sinkronisasi FE)
		user, err := h.s.GetByID(c.Request.Context(), id)
		if err == nil && user.Role == "admin" {
			helper.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Operasi dibatalkan: Akun %s adalah Admin dan tidak boleh dihapus", user.Name))
			return
		}
	}

	if err := h.s.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, fmt.Sprintf("%d pengguna berhasil dihapus", len(req.IDs)), nil)
}

func (h *UserHandler) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "akun berhasil dipulihkan", nil)
}

func (h *UserHandler) BulkRestore(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	if err := h.s.BulkRestore(c.Request.Context(), req.IDs); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, fmt.Sprintf("%d pengguna berhasil dipulihkan", len(req.IDs)), nil)
}

func (h *UserHandler) GetNotifications(c *gin.Context) {
	userID, _ := c.Get("user_id")
	ns, err := h.s.GetNotifications(c.Request.Context(), userID.(uint))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil", ns)
}

func (h *UserHandler) Export(c *gin.Context) {
	search := c.Query("search")
	role := c.Query("role")
	filter := c.Query("filter")
	status := c.Query("status")

	data, err := h.s.ExportExcel(c.Request.Context(), search, role, filter, status)
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	c.Header("Content-Disposition", "attachment; filename=Daftar_Pengguna.xlsx")
	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", data)
}

func (h *UserHandler) GetDependencyInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil", info)
}

func (h *UserHandler) CheckUnique(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")
	excludeIDStr := c.DefaultQuery("exclude_id", "0")
	excludeID, _ := strconv.Atoi(excludeIDStr)

	isUnique, err := h.s.CheckUnique(c.Request.Context(), field, value, uint(excludeID))
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "berhasil", gin.H{
		"is_unique": isUnique,
	})
}
