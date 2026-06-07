package dto

type AuditMeta struct {
	UserID      *uint  `json:"user_id,omitempty"`
	IP          string `json:"ip"`
	UserAgent   string `json:"user_agent"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	DeviceID    string `json:"device_id"`
	AppPlatform string `json:"app_platform"`
	AppVersion  string `json:"app_version"`
}
