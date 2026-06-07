package helper

import (
	"context"

	"github.com/ahmadzakyarifin/schoolpay/internal/dto"
	"github.com/gin-gonic/gin"
)

func BuildAuditMeta(c *gin.Context) dto.AuditMeta {
	var userID *uint

	if value, exists := c.Get("user_id"); exists {
		if id, ok := value.(uint); ok {
			userID = &id
		}
	}
	

	return dto.AuditMeta{
		UserID:      userID,
		IP:          c.ClientIP(),
		UserAgent:   c.Request.UserAgent(),
		Method:      c.Request.Method,
		Path:        c.FullPath(),
		DeviceID:    c.GetHeader("X-Device-ID"),
		AppPlatform: c.GetHeader("X-App-Platform"),
		AppVersion:  c.GetHeader("X-App-Version"),
	}
}

func GetAuditMeta(ctx context.Context) (uint, string, string, string, string) {
	var userID uint
	if v := ctx.Value("user_id"); v != nil {
		userID, _ = v.(uint)
	}
	var userName string
	if v := ctx.Value("user_name"); v != nil {
		userName, _ = v.(string)
	}
	var role string
	if v := ctx.Value("role"); v != nil {
		role, _ = v.(string)
	}
	var ipAddress string
	if v := ctx.Value("ip_address"); v != nil {
		ipAddress, _ = v.(string)
	}
	var userAgent string
	if v := ctx.Value("user_agent"); v != nil {
		userAgent, _ = v.(string)
	}
	if userName == "" {
		userName = "System/Admin"
	}
	if role == "" {
		role = "admin"
	}
	return userID, userName, role, ipAddress, userAgent
}
