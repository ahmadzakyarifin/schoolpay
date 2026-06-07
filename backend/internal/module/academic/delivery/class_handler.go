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

type ClassHandler struct {
	s usecase.ClassService
}

func NewClassHandler(s usecase.ClassService) *ClassHandler {
	return &ClassHandler{s: s}
}

func (h *ClassHandler) Create(c *gin.Context) {
	var j domain.Class
	if err := c.ShouldBindJSON(&j); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}
	if err := h.s.Create(c.Request.Context(), &j); err != nil {
		errMsg := err.Error()
		if strings.HasPrefix(strings.ToLower(errMsg), "gagal:") {
			field := "general"
			if strings.Contains(strings.ToLower(errMsg), "nama") {
				field = "name"
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
	helper.SuccessResponse(c, http.StatusCreated, "kelas berhasil dibuat", j)
}

func (h *ClassHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit > 1000 {
		limit = 1000
	}
	search := c.Query("search")
	status := c.Query("status")
	majorID := c.Query("major_id")
	ayID := c.Query("ay_id")
	sort := c.Query("sort")

	list, total, err := h.s.GetAll(c.Request.Context(), page, limit, search, status, majorID, ayID, sort)
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil mengambil data kelas", gin.H{
		"data":  list,
		"total": total,
	})
}

func (h *ClassHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var j domain.Class
	if err := c.ShouldBindJSON(&j); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}
	j.ID = uint(id)
	if err := h.s.Update(c.Request.Context(), &j); err != nil {
		errMsg := err.Error()
		if strings.HasPrefix(strings.ToLower(errMsg), "gagal:") {
			field := "general"
			if strings.Contains(strings.ToLower(errMsg), "nama") {
				field = "name"
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
	helper.SuccessResponse(c, http.StatusOK, "kelas berhasil diperbarui", j)
}

func (h *ClassHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "kelas berhasil dihapus", nil)
}

func (h *ClassHandler) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "kelas berhasil dipulihkan", nil)
}

func (h *ClassHandler) ToggleStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.ToggleStatus(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "status kelas berhasil diubah", nil)
}

func (h *ClassHandler) BulkDelete(c *gin.Context) {
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

func (h *ClassHandler) BulkRestore(c *gin.Context) {
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

func (h *ClassHandler) SuggestNextName(c *gin.Context) {
	name := c.Query("name")
	ayID, _ := strconv.Atoi(c.Query("ay_id"))
	majorID, _ := strconv.Atoi(c.Query("major_id"))
	excludeID, _ := strconv.Atoi(c.Query("exclude_id"))
	nextName, err := h.s.SuggestNextName(c.Request.Context(), name, uint(ayID), uint(majorID), uint(excludeID))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil mendapatkan saran nama", nextName)
}

func (h *ClassHandler) GetDependencyInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil", info)
}

func (h *ClassHandler) CheckUnique(c *gin.Context) {
	name := c.Query("name")
	majorID, _ := strconv.ParseUint(c.Query("major_id"), 10, 32)
	ayID, _ := strconv.ParseUint(c.Query("academic_year_id"), 10, 32)
	excludeID, _ := strconv.ParseUint(c.Query("exclude_id"), 10, 32)

	if name == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "name harus diisi")
		return
	}

	exists, err := h.s.CheckUnique(c.Request.Context(), name, uint(majorID), uint(ayID), uint(excludeID))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "berhasil", gin.H{"is_unique": !exists})
}
