package delivery

import (
	"net/http"
	"strconv"

	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/notification/usecase"
	userauthdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

type WhatsAppHandler struct {
	s    usecase.WhatsAppService
	noti  notificationrepo.NotificationRepo
	msg   utils.Messenger
	db    *bun.DB
	audit auditusecase.AuditLogService
}

func NewWhatsAppHandler(s usecase.WhatsAppService, noti notificationrepo.NotificationRepo, msg utils.Messenger, db *bun.DB, audit auditusecase.AuditLogService) *WhatsAppHandler {
	return &WhatsAppHandler{s: s, noti: noti, msg: msg, db: db, audit: audit}
}

func (h *WhatsAppHandler) GetStatus(c *gin.Context) {
	status, err := h.s.GetStatus()
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", gin.H{
		"status": status,
	})
}

func (h *WhatsAppHandler) GetQR(c *gin.Context) {
	qr, err := h.s.GetQR()
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	c.Data(http.StatusOK, "image/png", qr)
}


func (h *WhatsAppHandler) Logout(c *gin.Context) {
	if err := h.s.LogoutSession(); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	if h.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(c.Request.Context())
		_ = h.audit.Log(c.Request.Context(), h.db, userID, userName, role, "LOGOUT_WHATSAPP_SESSION", "whatsapp", 0, nil, map[string]interface{}{"status": "logged_out"}, ipAddress, userAgent)
	}
	utils.SuccessResponse(c, http.StatusOK, "WhatsApp berhasil logout", nil)
}

func (h *WhatsAppHandler) Restart(c *gin.Context) {
	if err := h.s.StopSession(); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	if err := h.s.StartSession(); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	go h.s.RegisterWebhook()

	if h.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(c.Request.Context())
		_ = h.audit.Log(c.Request.Context(), h.db, userID, userName, role, "RESTART_WHATSAPP_SESSION", "whatsapp", 0, nil, map[string]interface{}{"status": "restarting"}, ipAddress, userAgent)
	}

	utils.SuccessResponse(c, http.StatusOK, "WhatsApp session sedang direstart", gin.H{
		"status": "STARTING",
	})
}

func (h *WhatsAppHandler) GetStats(c *gin.Context) {
	stats, err := h.noti.GetStats(c.Request.Context())
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", stats)
}

func (h *WhatsAppHandler) GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	status := c.Query("status")
	search := c.Query("search")
	channel := c.Query("channel")

	list, total, err := h.noti.GetDetailedLogs(c.Request.Context(), page, limit, status, search, channel)
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", gin.H{
		"data":  list,
		"total": total,
	})
}

func (h *WhatsAppHandler) GetChatHistory(c *gin.Context) {
	phone := c.Param("phone")
	if phone == "" {
		utils.ErrorResponse(c, http.StatusBadRequest, "phone parameter is required")
		return
	}

	data, err := h.s.GetChatHistory(phone)
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "success", data)
}

func (h *WhatsAppHandler) SendChatMessage(c *gin.Context) {
	phone := c.Param("phone")
	var req struct {
		Message string `json:"message" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponseRaw(c, http.StatusBadRequest, err)
		return
	}

	if err := h.s.SendChatMessage(phone, req.Message); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}
	if h.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(c.Request.Context())
		_ = h.audit.Log(c.Request.Context(), h.db, userID, userName, role, "SEND_WHATSAPP_CHAT", "notifications", 0, nil, map[string]interface{}{"phone": phone}, ipAddress, userAgent)
	}

	utils.SuccessResponse(c, http.StatusOK, "message sent", nil)
}

func (h *WhatsAppHandler) ResendSpecificNotification(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		utils.ErrorResponse(c, http.StatusBadRequest, "ID tidak valid")
		return
	}

	noti, err := h.noti.FindByID(c.Request.Context(), uint(id))
	if err != nil {
		utils.ErrorResponseRaw(c, http.StatusNotFound, err)
		return
	}

	var user userauthdomain.User
	if err := h.db.NewSelect().Model(&user).Where("id = ?", noti.UserID).Scan(c.Request.Context()); err != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, err)
		return
	}

	var sendErr error
	var whatsappID string

	if noti.WhatsappID != nil && *noti.WhatsappID != "" && user.PhoneNumber != "" {
		wID, err := h.msg.SendWhatsApp(user.PhoneNumber, noti.Message)
		if err != nil {
			sendErr = err
		} else {
			whatsappID = wID
		}
	} else if user.Email != "" {
		err := h.msg.SendEmail(user.Email, noti.Title, noti.Message)
		if err != nil {
			sendErr = err
		}
	} else {
		if user.PhoneNumber != "" {
			wID, err := h.msg.SendWhatsApp(user.PhoneNumber, noti.Message)
			if err != nil {
				sendErr = err
			} else {
				whatsappID = wID
			}
		} else if user.Email != "" {
			err := h.msg.SendEmail(user.Email, noti.Title, noti.Message)
			if err != nil {
				sendErr = err
			}
		}
	}

	if sendErr != nil {
		noti.DeliveryStatus = "failed"
		noti.DeliveryError = utils.StringPtr(sendErr.Error())
	} else {
		noti.DeliveryStatus = "sent"
		noti.DeliveryError = nil
		if whatsappID != "" {
			noti.WhatsappID = &whatsappID
		}
	}

	_, _ = h.db.NewUpdate().Model(noti).Column("delivery_status", "delivery_error", "whatsapp_id", "updated_at").WherePK().Exec(c.Request.Context())
	if h.audit != nil {
		userID, userName, role, ipAddress, userAgent := utils.GetAuditMeta(c.Request.Context())
		_ = h.audit.Log(c.Request.Context(), h.db, userID, userName, role, "RESEND_NOTIFICATION", "notifications", noti.ID, nil, map[string]interface{}{"status": noti.DeliveryStatus, "user_id": noti.UserID, "error": func() string { if noti.DeliveryError != nil { return *noti.DeliveryError }; return "" }()}, ipAddress, userAgent)
	}

	if sendErr != nil {
		utils.ErrorResponseRaw(c, http.StatusInternalServerError, sendErr)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "Notifikasi berhasil dikirim ulang", noti)
}
