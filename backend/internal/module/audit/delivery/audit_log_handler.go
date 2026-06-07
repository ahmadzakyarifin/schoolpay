package delivery

import (
	"net/http"
	"strconv"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	userdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/helper"
	"github.com/gin-gonic/gin"
)

type AuditLogHandler struct {
	service usecase.AuditLogService
}

func NewAuditLogHandler(s usecase.AuditLogService) *AuditLogHandler {
	return &AuditLogHandler{service: s}
}

func (h *AuditLogHandler) GetLogs(c *gin.Context) {
	currentUser, _ := c.Get("user")
	userObj, _ := currentUser.(*userdomain.User)

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	filter := make(map[string]interface{})
	if action := c.Query("action"); action != "" {
		filter["action"] = action
	}
	if entityType := c.Query("entity_type"); entityType != "" {
		filter["entity_type"] = entityType
	}
	if role := c.Query("role"); role != "" {
		filter["role"] = role
	}
	if userName := c.Query("user_name"); userName != "" {
		filter["user_name"] = userName
	}
	if sort := c.Query("sort"); sort != "" {
		filter["sort"] = sort
	}

	logs, total, err := h.service.GetLogs(c.Request.Context(), userObj, filter, page, limit)
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil log audit")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Berhasil mengambil log audit", gin.H{
		"logs":  logs,
		"total": total,
		"page":  page,
		"limit": limit,
	})
}

func (h *AuditLogHandler) GetEntityLogs(c *gin.Context) {
	currentUser, _ := c.Get("user")
	userObj, _ := currentUser.(*userdomain.User)

	entityType := c.Param("entityType")
	entityIDStr := c.Param("entityID")
	entityID, _ := strconv.Atoi(entityIDStr)

	logs, err := h.service.GetEntityLogs(c.Request.Context(), userObj, entityType, uint(entityID))
	if err != nil {
		helper.ErrorResponse(c, http.StatusInternalServerError, "Gagal mengambil log spesifik entitas")
		return
	}

	helper.SuccessResponse(c, http.StatusOK, "Berhasil mengambil log spesifik entitas", logs)
}
