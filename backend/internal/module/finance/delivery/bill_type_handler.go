package delivery

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/gin-gonic/gin"
)

type BillTypeHandler struct {
	s usecase.BillTypeService
}

func NewBillTypeHandler(s usecase.BillTypeService) *BillTypeHandler {
	return &BillTypeHandler{s: s}
}

func (h *BillTypeHandler) Create(c *gin.Context) {
	var b domain.BillType
	if err := c.ShouldBindJSON(&b); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}
	if err := h.s.Create(c.Request.Context(), &b); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "jenis tagihan berhasil dibuat", b)
}

func (h *BillTypeHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	filterType := c.Query("type")
	status := c.Query("status")
	sort := c.Query("sort")

	list, total, err := h.s.GetAllPaged(c.Request.Context(), page, limit, search, filterType, status, sort)
	if err != nil {
		log.Printf("[ERROR] BillType.GetAll: %v", err)
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil mengambil data jenis tagihan", gin.H{
		"data":  list,
		"total": total,
	})
}

func (h *BillTypeHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var b domain.BillType
	if err := c.ShouldBindJSON(&b); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}
	b.ID = uint(id)

	if err := h.s.Update(c.Request.Context(), &b); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "jenis tagihan berhasil diperbarui", b)
}

func (h *BillTypeHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "jenis tagihan berhasil dihapus", nil)
}

func (h *BillTypeHandler) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "jenis tagihan berhasil dipulihkan", nil)
}

func (h *BillTypeHandler) ToggleStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.ToggleStatus(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "status jenis tagihan berhasil diubah", nil)
}

func (h *BillTypeHandler) BulkDelete(c *gin.Context) {
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

	helper.SuccessResponse(c, http.StatusOK, fmt.Sprintf("%d jenis tagihan berhasil dihapus", len(req.IDs)), nil)
}

func (h *BillTypeHandler) BulkRestore(c *gin.Context) {
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

	helper.SuccessResponse(c, http.StatusOK, fmt.Sprintf("%d jenis tagihan berhasil dipulihkan", len(req.IDs)), nil)
}

func (h *BillTypeHandler) GetDependencyInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil", info)
}

func (h *BillTypeHandler) CheckUnique(c *gin.Context) {
	name := c.Query("name")
	excludeID, _ := strconv.ParseUint(c.Query("exclude_id"), 10, 32)

	if name == "" {
		helper.ErrorResponse(c, http.StatusBadRequest, "name harus diisi")
		return
	}

	exists, err := h.s.CheckUnique(c.Request.Context(), name, uint(excludeID))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "berhasil", gin.H{"is_unique": !exists})
}
