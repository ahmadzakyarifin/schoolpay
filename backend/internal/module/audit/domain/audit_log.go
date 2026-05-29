package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type AuditLog struct {
	bun.BaseModel `bun:"table:audit_logs,alias:al"`

	ID         uint64                 `bun:"id,pk,autoincrement" json:"id"`
	UserID     uint                   `bun:"user_id" json:"user_id"`
	UserName   string                 `bun:"user_name" json:"user_name"`
	Role       string                 `bun:"role" json:"role"`
	Action     string                 `bun:"action" json:"action"`           // CREATE, UPDATE, DELETE, RESTORE
	EntityType string                 `bun:"entity_type" json:"entity_type"` // students, users, classes, majors, etc.
	EntityID   uint                   `bun:"entity_id" json:"entity_id"`
	OldValues  map[string]interface{} `bun:"old_values,type:json" json:"old_values,omitempty"`
	NewValues  map[string]interface{} `bun:"new_values,type:json" json:"new_values,omitempty"`
	IPAddress  string                 `bun:"ip_address" json:"ip_address"`
	UserAgent  string                 `bun:"user_agent" json:"user_agent"`
	CreatedAt  time.Time              `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}
