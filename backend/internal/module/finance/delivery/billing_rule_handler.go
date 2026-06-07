package delivery

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/gin-gonic/gin"
)

type BillingRuleHandler struct {
	s usecase.BillingRuleService
}

func NewBillingRuleHandler(s usecase.BillingRuleService) *BillingRuleHandler {
	return &BillingRuleHandler{s: s}
}

func (h *BillingRuleHandler) Create(c *gin.Context) {
	var br domain.BillingRule
	if err := c.ShouldBindJSON(&br); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}
	if err := h.s.Create(c.Request.Context(), &br); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusCreated, "aturan tagihan berhasil dibuat", br)
}

func (h *BillingRuleHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	status := c.Query("status")
	generateStatus := c.Query("generate_status")
	sort := c.Query("sort")

	list, total, err := h.s.GetAllPaged(c.Request.Context(), page, limit, search, status, generateStatus, sort)
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil mengambil data aturan tagihan", gin.H{
		"data":  list,
		"total": total,
	})
}

func (h *BillingRuleHandler) Update(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var br domain.BillingRule
	if err := c.ShouldBindJSON(&br); err != nil {
		helper.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", helper.GetValidationErrors(err))
		return
	}
	br.ID = uint(id)

	if err := h.s.Update(c.Request.Context(), &br); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "aturan tagihan berhasil diperbarui", br)
}

func (h *BillingRuleHandler) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Delete(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "aturan tagihan berhasil dihapus", nil)
}

func (h *BillingRuleHandler) Restore(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.Restore(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "aturan tagihan berhasil dipulihkan", nil)
}

func (h *BillingRuleHandler) ToggleStatus(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	if err := h.s.ToggleStatus(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "status aturan tagihan berhasil diubah", nil)
}

func (h *BillingRuleHandler) BulkDelete(c *gin.Context) {
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

	helper.SuccessResponse(c, http.StatusOK, fmt.Sprintf("%d aturan tagihan berhasil dihapus", len(req.IDs)), nil)
}

func (h *BillingRuleHandler) BulkRestore(c *gin.Context) {
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

	helper.SuccessResponse(c, http.StatusOK, fmt.Sprintf("%d aturan tagihan berhasil dipulihkan", len(req.IDs)), nil)
}

func (h *BillingRuleHandler) GetDependencyInfo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	info, err := h.s.GetDependencyInfo(c.Request.Context(), uint(id))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil", info)
}

func (h *BillingRuleHandler) CheckUnique(c *gin.Context) {
	billTypeID, _ := strconv.ParseUint(c.Query("bill_type_id"), 10, 32)
	targetType := c.Query("target_type")
	targetID, _ := strconv.ParseUint(c.Query("target_id"), 10, 32)
	classIDStr := c.Query("class_id")
	excludeID, _ := strconv.ParseUint(c.Query("exclude_id"), 10, 32)

	var classID *uint
	if classIDStr != "" {
		cid, _ := strconv.ParseUint(classIDStr, 10, 32)
		cidUint := uint(cid)
		classID = &cidUint
	}

	exists, err := h.s.CheckUnique(c.Request.Context(), uint(billTypeID), targetType, uint(targetID), classID, uint(excludeID))
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "berhasil", gin.H{"is_unique": !exists})
}
