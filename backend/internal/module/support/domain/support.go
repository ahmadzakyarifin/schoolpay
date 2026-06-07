package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Conversation struct {
	bun.BaseModel `bun:"table:support_conversations,alias:sc"`

	ID          uint      `bun:"id,pk,autoincrement" json:"id"`
	ParentID    *uint     `bun:"parent_id" json:"parent_id,omitempty"`
	PhoneNumber string    `bun:"phone_number" json:"phone_number"`
	ParentName  *string   `bun:"parent_name" json:"parent_name,omitempty"`
	Status      string    `bun:"status,default:open" json:"status"`
	CreatedAt   time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt   time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`

	// Custom transient field
	StudentNames   string `bun:"-" json:"student_names,omitempty"`
	WhatsAppWebURL string `bun:"-" json:"whatsapp_web_url,omitempty"`
}
