package delivery

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
)

type MajorHandler struct {
	s usecase.MajorService
}

func NewMajorHandler(s usecase.MajorService) *MajorHandler {
	return &MajorHandler{s: s}
}

func (h *MajorHandler) Create(c *gin.Context) {
	var j domain.Major
	if err := c.ShouldBindJSON(&j); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}
	if err := h.s.Create(c.Request.Context(), &j); err != nil {
		errMsg := err.Error()
		// Cek apakah ini error bisnis (duplikasi)
		if strings.HasPrefix(strings.ToLower(errMsg), "gagal:") {
			field := "general"
			if strings.Contains(strings.ToLower(errMsg), "kode") {
				field = "code"
			}
			if strings.Contains(strings.ToLower(errMsg), "nama") {
				field = "name"
			}

			// Ambil pesan setelah "gagal: "
			cleanMsg := errMsg
			if len(errMsg) > 7 {
				cleanMsg = strings.TrimSpace(errMsg[6:])
			}

			utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", map[string][]string{field: {cleanMsg}})
			return
		}
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "jurusan berhasil dibuat", j)
}

func (h *MajorHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 1000 {
		limit = 1000
	}
	search := c.Query("search")
	status := c.Query("status")
	sort := c.Query("sort")

	list, total, err := h.s.GetAll(c.Request.Context(), page, limit, search, status, sort)
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil mengambil data jurusan", gin.H{
		"data":  list,
		"total": total,
	})
}

func (h *MajorHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var j domain.Major
	if err := c.ShouldBindJSON(&j); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}
	j.ID = uint(id)
	if err := h.s.Update(c.Request.Context(), &j); err != nil {
		errMsg := err.Error()
		// Cek apakah ini error bisnis (duplikasi)
		if strings.HasPrefix(strings.ToLower(errMsg), "gagal:") {
			field := "general"
			if strings.Contains(strings.ToLower(errMsg), "kode") {
				field = "code"
			}
			if strings.Contains(strings.ToLower(errMsg), "nama") {
				field = "name"
			}

			// Ambil pesan setelah "gagal: "
			cleanMsg := errMsg
			if len(errMsg) > 7 {
				cleanMsg = strings.TrimSpace(errMsg[6:])
			}

			utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", map[string][]string{field: {cleanMsg}})
			return
		}
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "jurusan berhasil diperbarui", j)
}

func (h *MajorHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "jurusan berhasil dihapus", nil)
}

func (h *MajorHandler) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "jurusan berhasil dipulihkan", nil)
}

func (h *MajorHandler) ToggleStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.ToggleStatus(c.Request.Context(), uint(id)); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "status jurusan berhasil diubah", nil)
}

func (h *MajorHandler) BulkDelete(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	if err := h.s.BulkDelete(c.Request.Context(), req.IDs); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "data berhasil dihapus", nil)
}

func (h *MajorHandler) BulkRestore(c *gin.Context) {
	var req struct {
		IDs []uint `json:"ids" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	if err := h.s.BulkRestore(c.Request.Context(), req.IDs); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "data berhasil dipulihkan", nil)
}

func (h *MajorHandler) GetDependencyInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", info)
}

func (h *MajorHandler) CheckUnique(c *gin.Context) {
	field := c.Query("field")
	value := c.Query("value")
	excludeID, _ := strconv.ParseUint(c.Query("exclude_id"), 10, 32)

	if field == "" || value == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "field dan value harus diisi")
		return
	}

	exists, err := h.s.CheckUnique(c.Request.Context(), field, value, uint(excludeID))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "berhasil", gin.H{"is_unique": !exists})
}
