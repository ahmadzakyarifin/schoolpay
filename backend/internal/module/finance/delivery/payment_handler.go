package delivery

import (
	"net/http"
	"strconv"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/finance/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
)

type PaymentHandler struct {
	s usecase.PaymentService
}

func NewPaymentHandler(s usecase.PaymentService) *PaymentHandler {
	return &PaymentHandler{s: s}
}

func (h *PaymentHandler) Process(c *gin.Context) {
	var req struct {
		StudentID       uint       `json:"student_id" binding:"required"`
		Amount          float64    `json:"amount" binding:"required"`
		DepositApplied  float64    `json:"deposit_applied"`
		Channel         string     `json:"channel" binding:"required"`
		Method          string     `json:"method" binding:"required"`
		Reference       string     `json:"reference"`
		BillIDs         []uint     `json:"bill_ids"`
		IsBypassRule    bool       `json:"is_bypass_rule"`
		BypassReason    string     `json:"bypass_reason"`
		Note            string     `json:"note"`
		ProofAttachment string     `json:"proof_attachment"`
		PaidAt          *time.Time `json:"paid_at"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}
	if req.Amount <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "nominal pembayaran harus lebih dari 0")
		return
	}
	if req.IsBypassRule && req.BypassReason == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "alasan bypass wajib diisi")
		return
	}

	p := &domain.Payment{
		Amount:         req.Amount,
		DepositApplied: req.DepositApplied,
		Channel:        req.Channel,
		Method:         req.Method,
		TransactionRef: req.Reference,
		IsBypassRule:   req.IsBypassRule,
		PaidAt:         req.PaidAt,
	}
	if req.BypassReason != "" {
		p.BypassReason = &req.BypassReason
	}
	if req.Note != "" {
		p.Note = &req.Note
	}
	if req.ProofAttachment != "" {
		p.ProofAttachment = &req.ProofAttachment
	}

	if err := h.s.Process(c.Request.Context(), req.StudentID, p, req.BillIDs); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "pembayaran berhasil diproses", p)
}

func (h *PaymentHandler) CreateIntent(c *gin.Context) {
	var req struct {
		StudentID    uint    `json:"student_id" binding:"required"`
		Amount         float64 `json:"amount" binding:"required"`
		DepositApplied float64 `json:"deposit_applied"`
		BillIDs        []uint  `json:"bill_ids" binding:"required"`
		IsBypassRule bool    `json:"is_bypass_rule"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}
	if req.Amount <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "nominal pembayaran harus lebih dari 0")
		return
	}

	p, err := h.s.CreateIntent(c.Request.Context(), req.StudentID, req.Amount, req.DepositApplied, req.BillIDs, req.IsBypassRule)
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "intent created", p)
}

func (h *PaymentHandler) GetReceipt(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}
	receipt, err := h.s.GetReceipt(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusNotFound, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", receipt)
}

func (h *PaymentHandler) GetHistory(c *gin.Context) {
	studentIDStr := c.Query("student_id")
	if studentIDStr == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "student_id diperlukan")
		return
	}

	id, _ := strconv.Atoi(studentIDStr)
	list, err := h.s.GetHistory(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", list)
}

func (h *PaymentHandler) HandleWebhook(c *gin.Context) {
	var payload map[string]interface{}
	if err := c.ShouldBindJSON(&payload); err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "invalid payload")
		return
	}

	if err := h.s.HandleWebhook(c.Request.Context(), payload); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	c.Status(http.StatusOK)
}
