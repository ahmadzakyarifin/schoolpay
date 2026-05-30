package delivery

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
)

type StudentBillHandler struct {
	s      usecase.StudentBillService
	paySvc usecase.PaymentService
}

func NewStudentBillHandler(s usecase.StudentBillService, paySvc usecase.PaymentService) *StudentBillHandler {
	return &StudentBillHandler{s: s, paySvc: paySvc}
}

func (h *StudentBillHandler) GetAll(c *gin.Context) {
	ctx := c.Request.Context()
	studentIDStr := c.Query("student_id")
	if studentIDStr != "" {
		id, err := strconv.Atoi(studentIDStr)
		if err != nil || id <= 0 {
			utils.ErrorResponse(c, http.StatusBadRequest, "student_id tidak valid")
			return
		}
		list, err := h.s.GetByStudent(ctx, uint(id))
		if err != nil {
			utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
			return
		}
		utils.SuccessResponse(c, http.StatusOK, "berhasil", list)
		return
	}

	search := c.Query("search")
	sort := c.Query("sort")
	list, err := h.s.GetAll(ctx, search, sort)
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", list)
}

func (h *StudentBillHandler) GetMyBills(c *gin.Context) {
	userID, _ := c.Get("user_id")
	list, err := h.s.GetByParent(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", list)
}

func (h *StudentBillHandler) GetStudentSummaries(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	search := c.Query("search")
	sort := c.Query("sort")
	status := c.Query("status")

	list, total, err := h.s.GetStudentSummaries(c.Request.Context(), search, sort, status, page, limit)
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", gin.H{
		"data":  list,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *StudentBillHandler) Create(c *gin.Context) {
	var b domain.StudentBill
	if err := c.ShouldBindJSON(&b); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}
	if err := h.s.Create(c.Request.Context(), &b); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusCreated, "tagihan berhasil dibuat", b)
}

func (h *StudentBillHandler) Update(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}
	var b domain.StudentBill
	if err := c.ShouldBindJSON(&b); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}
	b.ID = uint(id)

	if err := h.s.Update(c.Request.Context(), &b); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "tagihan berhasil diperbarui", b)
}

func (h *StudentBillHandler) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}
	var req struct {
		Reason string `json:"reason"`
	}
	if c.Request.ContentLength > 0 {
		if err := c.ShouldBindJSON(&req); err != nil {
			utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
			return
		}
	}
	if strings.TrimSpace(req.Reason) == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "alasan reset/void tagihan wajib diisi untuk kebutuhan audit")
		return
	}
	if err := h.s.Delete(c.Request.Context(), uint(id), strings.TrimSpace(req.Reason)); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "tagihan berhasil dihapus", nil)
}

func (h *StudentBillHandler) Generate(c *gin.Context) {
	var req struct {
		RuleID           uint   `json:"rule_id" binding:"required"`
		CustomReason     string `json:"custom_reason"`
		CustomMessage    string `json:"custom_message"`
		SkipNotification bool   `json:"skip_notification"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}

	if err := h.s.GenerateFromRule(c.Request.Context(), req.RuleID, req.CustomReason, req.CustomMessage, req.SkipNotification); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "tagihan berhasil digenerate", nil)
}

func (h *StudentBillHandler) BulkGenerate(c *gin.Context) {
	var req struct {
		RuleIDs          []uint `json:"rule_ids" binding:"required"`
		CustomReason     string `json:"custom_reason"`
		CustomMessage    string `json:"custom_message"`
		SkipNotification bool   `json:"skip_notification"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}

	if err := h.s.BulkGenerateFromRules(c.Request.Context(), req.RuleIDs, req.CustomReason, req.CustomMessage, req.SkipNotification); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "tagihan berhasil digenerate", nil)
}

func (h *StudentBillHandler) BulkCancel(c *gin.Context) {
	var req struct {
		RuleIDs          []uint `json:"rule_ids" binding:"required"`
		CustomReason     string `json:"custom_reason"`
		CustomMessage    string `json:"custom_message"`
		SkipNotification bool   `json:"skip_notification"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}
	if strings.TrimSpace(req.CustomReason) == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "alasan penarikan tagihan wajib diisi untuk kebutuhan audit")
		return
	}

	if err := h.s.CancelGeneratedBills(c.Request.Context(), req.RuleIDs, req.CustomReason, req.CustomMessage, req.SkipNotification); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "tagihan berhasil ditarik", nil)
}

func (h *StudentBillHandler) Remind(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}
	if err := h.s.SendManualReminder(c.Request.Context(), uint(id)); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "pengingat berhasil dijadwalkan", nil)
}

func (h *StudentBillHandler) MarkAsPaidManual(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}
	ctx := c.Request.Context()

	var req struct {
		Amount          float64    `json:"amount"`
		PaymentMethod   string     `json:"payment_method"`
		PaidAt          *time.Time `json:"paid_at"`
		ReferenceNumber string     `json:"reference_number"`
		Note            string     `json:"note"`
		ProofAttachment string     `json:"proof_attachment"`
		IsBypassRule    bool       `json:"is_bypass_rule"`
		BypassReason    string     `json:"bypass_reason"`
		Reason          string     `json:"reason"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}

	bill, err := h.s.GetByID(ctx, uint(id))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusNotFound, err)
		return
	}

	if bill.Status == "voided" {
		utils.ErrorResponse(c, http.StatusBadRequest, "tagihan sudah dibatalkan (void)")
		return
	}

	remaining := bill.Amount - bill.TotalPaid
	if remaining <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "tagihan sudah lunas")
		return
	}

	payAmount := remaining
	if req.Amount > 0 {
		if req.Amount > remaining {
			utils.ErrorResponse(c, http.StatusBadRequest, "nominal pembayaran melebihi sisa tagihan")
			return
		}
		payAmount = req.Amount
	}

	if payAmount <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "nominal pembayaran harus lebih dari 0")
		return
	}
	if req.IsBypassRule && req.BypassReason == "" && req.Reason == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "alasan bypass wajib diisi")
		return
	}

	methodName := "Tunai"
	if req.PaymentMethod != "" {
		methodName = req.PaymentMethod
	}

	refNum := req.ReferenceNumber
	if refNum == "" {
		refNum = fmt.Sprintf("MANUAL-UPD-%d-%d", bill.ID, time.Now().Unix())
	}

	paidAtTime := time.Now()
	if req.PaidAt != nil {
		paidAtTime = *req.PaidAt
	}

	adminID, adminName, _, _, _ := utils.GetAuditMeta(ctx)
	createdBy := "SYSTEM"
	if adminName != "" {
		createdBy = fmt.Sprintf("ADMIN-%d (%s)", adminID, adminName)
	}

	p := &domain.Payment{
		StudentID:      bill.StudentID,
		Amount:         payAmount,
		Channel:        "cash",
		Method:         methodName,
		Status:         "success",
		TransactionRef: refNum,
		PaidAt:         &paidAtTime,
		IsBypassRule:   req.IsBypassRule,
		CreatedBy:      createdBy,
	}

	if req.BypassReason != "" {
		p.BypassReason = &req.BypassReason
	} else if req.Reason != "" {
		p.BypassReason = &req.Reason
	}
	if req.ProofAttachment != "" {
		p.ProofAttachment = &req.ProofAttachment
	}
	if req.Note != "" {
		p.Note = &req.Note
	}

	if err := h.paySvc.Process(ctx, bill.StudentID, p, []uint{bill.ID}); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "pembayaran manual berhasil dicatat", p)
}
