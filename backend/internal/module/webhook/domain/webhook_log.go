package domain

import (
	"encoding/json"
	"time"

	"github.com/uptrace/bun"
)

type WebhookLog struct {
	bun.BaseModel `bun:"table:webhook_logs,alias:wl"`

	ID        uint            `bun:"id,pk,autoincrement" json:"id"`
	Provider  string          `bun:"provider" json:"provider"`
	EventID   string          `bun:"event_id,unique" json:"event_id"`
	Payload   json.RawMessage `bun:"payload" json:"payload"`
	Status    string          `bun:"status" json:"status"`
	CreatedAt time.Time       `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time       `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time      `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at,omitempty"`
}
