package delivery

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/usecase"
	userdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/websocket"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
)

type SupportHandler struct {
	s   usecase.SupportService
	hub *websocket.Hub
}

type supportConversationLookup interface {
	Conversation(ctx context.Context, conversationID uint) (*domain.Conversation, error)
}

func NewSupportHandler(s usecase.SupportService, hubs ...*websocket.Hub) *SupportHandler {
	var hub *websocket.Hub
	if len(hubs) > 0 {
		hub = hubs[0]
	}
	return &SupportHandler{s: s, hub: hub}
}

func (h *SupportHandler) broadcast(conversationID uint, parentID *uint) {
	if h.hub == nil {
		return
	}
	data := map[string]interface{}{"conversation_id": conversationID}
	if parentID != nil {
		data["parent_id"] = *parentID
		h.hub.BroadcastToUser("SUPPORT_CHAT_UPDATED", data, *parentID)
	}
	h.hub.BroadcastToRoles("SUPPORT_CHAT_UPDATED", data, "admin")
}

func (h *SupportHandler) List(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	status := c.Query("status")
	list, total, err := h.s.List(c.Request.Context(), status, page, limit)
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", gin.H{"data": list, "total": total})
}

func (h *SupportHandler) Messages(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}
	list, err := h.s.Messages(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", list)
}

func (h *SupportHandler) Reply(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}
	var req struct {
		Message string `json:"message" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}
	adminID, _ := c.Get("user_id")
	if err := h.s.Reply(c.Request.Context(), uint(id), adminID.(uint), req.Message); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	var parentID *uint
	if lookup, ok := h.s.(supportConversationLookup); ok {
		if conv, err := lookup.Conversation(c.Request.Context(), uint(id)); err == nil && conv != nil {
			parentID = conv.ParentID
		}
	}
	h.broadcast(uint(id), parentID)
	utils.SuccessResponse(c, http.StatusOK, "pesan berhasil dikirim", nil)
}

func (h *SupportHandler) Assign(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}
	adminID, _ := c.Get("user_id")
	if err := h.s.Assign(c.Request.Context(), uint(id), adminID.(uint)); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	h.broadcast(uint(id), nil)
	utils.SuccessResponse(c, http.StatusOK, "percakapan berhasil diambil", nil)
}

func (h *SupportHandler) Close(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil || id <= 0 {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}
	if err := h.s.Close(c.Request.Context(), uint(id)); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	h.broadcast(uint(id), nil)
	utils.SuccessResponse(c, http.StatusOK, "percakapan ditutup", nil)
}

func (h *SupportHandler) ParentConversation(c *gin.Context) {
	userID, _ := c.Get("user_id")
	conv, err := h.s.ParentConversation(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", conv)
}

func (h *SupportHandler) ParentMessages(c *gin.Context) {
	userID, _ := c.Get("user_id")
	messages, err := h.s.ParentMessages(c.Request.Context(), userID.(uint))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	utils.SuccessResponse(c, http.StatusOK, "berhasil", messages)
}

func (h *SupportHandler) ParentSendMessage(c *gin.Context) {
	var req struct {
		Topic   string `json:"topic"`
		Message string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorValidationResponse(c, http.StatusBadRequest, "validasi gagal", utils.GetValidationErrors(err))
		return
	}
	req.Topic = strings.TrimSpace(req.Topic)
	req.Message = strings.TrimSpace(req.Message)

	userRaw, exists := c.Get("user")
	if !exists {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user tidak ditemukan")
		return
	}
	parent, ok := userRaw.(*userdomain.User)
	if !ok || parent == nil {
		utils.ErrorResponse(c, http.StatusUnauthorized, "user tidak valid")
		return
	}

	conv, err := h.s.ParentSendMessage(c.Request.Context(), parent, req.Topic, req.Message)
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}
	h.broadcast(conv.ID, conv.ParentID)
	utils.SuccessResponse(c, http.StatusCreated, "pesan berhasil dikirim ke CS Admin", conv)
}
