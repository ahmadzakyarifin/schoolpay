package delivery

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/gin-gonic/gin"
)

type AcademicYearHandler struct {
	s usecase.AcademicYearService
}

func NewAcademicYearHandler(s usecase.AcademicYearService) *AcademicYearHandler {
	return &AcademicYearHandler{s: s}
}

func (h *AcademicYearHandler) Create(c *gin.Context) {
	var ay domain.AcademicYear
	if err := c.ShouldBindJSON(&ay); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}
	if err := h.s.Create(c.Request.Context(), &ay); err != nil {
		errMsg := err.Error()
		if strings.HasPrefix(strings.ToLower(errMsg), "gagal:") {
			field := "general"
			if strings.Contains(strings.ToLower(errMsg), "jurusan") {
				field = "major_ids"
			} else if strings.Contains(strings.ToLower(errMsg), "kelas") {
				field = "class_ids"
			} else if strings.Contains(strings.ToLower(errMsg), "angkatan") {
				field = "year"
			}
			cleanMsg := errMsg
			if len(errMsg) > 7 {
				cleanMsg = strings.TrimSpace(errMsg[6:])
			}
			helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", map[string][]string{field: {cleanMsg}})
			return
		}
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "angkatan berhasil dibuat", ay)
}

func (h *AcademicYearHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	status := c.Query("status")
	sort := c.Query("sort")

	list, total, err := h.s.GetAll(c.Request.Context(), page, limit, search, status, sort)
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil mengambil data angkatan", gin.H{
		"data":  list,
		"total": total,
	})
}

func (h *AcademicYearHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var ay domain.AcademicYear
	if err := c.ShouldBindJSON(&ay); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}
	ay.ID = uint(id)
	if err := h.s.Update(c.Request.Context(), &ay); err != nil {
		errMsg := err.Error()
		if strings.HasPrefix(strings.ToLower(errMsg), "gagal:") {
			field := "general"
			if strings.Contains(strings.ToLower(errMsg), "jurusan") {
				field = "major_ids"
			} else if strings.Contains(strings.ToLower(errMsg), "kelas") {
				field = "class_ids"
			} else if strings.Contains(strings.ToLower(errMsg), "angkatan") {
				field = "year"
			}
			cleanMsg := errMsg
			if len(errMsg) > 7 {
				cleanMsg = strings.TrimSpace(errMsg[6:])
			}
			helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", map[string][]string{field: {cleanMsg}})
			return
		}
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "angkatan berhasil diperbarui", ay)
}

func (h *AcademicYearHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "angkatan berhasil dihapus", nil)
}

func (h *AcademicYearHandler) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "angkatan berhasil dipulihkan", nil)
}

func (h *AcademicYearHandler) BulkDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}
	if err := h.s.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "data berhasil dihapus", nil)
}

func (h *AcademicYearHandler) BulkRestore(c *gin.Context) {
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
	helper.SuccessResponse(c, http.StatusOK, "data berhasil dipulihkan", nil)
}

func (h *AcademicYearHandler) GetActive(c *gin.Context) {
	list, err := h.s.GetActive(c.Request.Context())
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil mengambil data angkatan aktif", list)
}

func (h *AcademicYearHandler) AssignMajors(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		MajorIDs []uint `json:"major_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "data tidak valid")
		return
	}
	if err := h.s.AssignMajors(c.Request.Context(), uint(id), req.MajorIDs); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "jurusan berhasil dikaitkan ke angkatan", nil)
}

func (h *AcademicYearHandler) GetMajors(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	list, err := h.s.GetMajorsByYear(c.Request.Context(), uint(id))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil mengambil data jurusan untuk angkatan", list)
}

func (h *AcademicYearHandler) AssignClasses(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var req struct {
		ClassIDs []uint `json:"class_ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "data tidak valid")
		return
	}
	if err := h.s.AssignClasses(c.Request.Context(), uint(id), req.ClassIDs); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "kelas berhasil dikaitkan ke angkatan", nil)
}

func (h *AcademicYearHandler) GetClasses(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	list, err := h.s.GetClassesByYear(c.Request.Context(), uint(id))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil mengambil data kelas untuk angkatan", list)
}

func (h *AcademicYearHandler) GetDependencyInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil", info)
}

func (h *AcademicYearHandler) CheckUnique(c *gin.Context) {
	yearStr := c.Query("year")
	excludeID, _ := strconv.ParseUint(c.Query("exclude_id"), 10, 32)

	if yearStr == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "year harus diisi")
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "format year tidak valid")
		return
	}

	exists, err := h.s.CheckUnique(c.Request.Context(), year, uint(excludeID))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "berhasil", gin.H{"is_unique": !exists})
}
