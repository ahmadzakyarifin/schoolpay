package delivery

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/websocket"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/gin-gonic/gin"
)

type SupportHandler struct {
	s   usecase.SupportService
	hub *websocket.Hub
}

func NewSupportHandler(s usecase.SupportService, hubs ...*websocket.Hub) *SupportHandler {
	var hub *websocket.Hub
	if len(hubs) > 0 {
		hub = hubs[0]
	}
	return &SupportHandler{s: s, hub: hub}
}

func (h *SupportHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	list, total, err := h.s.List(c.Request.Context(), status, page, limit)
	if err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	helper.SuccessResponse(c, http.StatusOK, "berhasil", gin.H{"data": list, "total": total})
}

func (h *SupportHandler) Assign(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		helper.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}
	adminID, _ := c.Get("user_id")
	if err := h.s.Assign(c.Request.Context(), uint(id), adminID.(uint)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	if h.hub != nil {
		h.hub.BroadcastToRoles("SUPPORT_CHAT_UPDATED", map[string]interface{}{"conversation_id": uint(id)}, "admin")
	}
	helper.SuccessResponse(c, http.StatusOK, "percakapan berhasil diambil", nil)
}

func (h *SupportHandler) Close(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		helper.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}
	if err := h.s.Close(c.Request.Context(), uint(id)); err != nil {
		helper.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	if h.hub != nil {
		h.hub.BroadcastToRoles("SUPPORT_CHAT_UPDATED", map[string]interface{}{"conversation_id": uint(id)}, "admin")
	}
	helper.SuccessResponse(c, http.StatusOK, "percakapan ditutup", nil)
}

func (h *SupportHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		helper.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	var payload struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		helper.ErrorResponse(c, http.StatusBadRequest, "status wajib diisi")
		return
	}

	if err := h.s.UpdateStatus(c.Request.Context(), uint(id), payload.Status); err != nil {
		helper.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	if h.hub != nil {
		h.hub.BroadcastToRoles("SUPPORT_CHAT_UPDATED", map[string]interface{}{"conversation_id": uint(id), "status": payload.Status}, "admin")
	}
	helper.SuccessResponse(c, http.StatusOK, "status percakapan diperbarui", nil)
}
