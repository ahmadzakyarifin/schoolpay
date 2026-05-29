package utils

import (
	"context"
)

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
