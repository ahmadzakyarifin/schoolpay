package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Notification struct {
	bun.BaseModel `bun:"table:notifications,alias:n"`

	ID             uint       `bun:"id,pk,autoincrement" json:"id"`
	UserID         uint       `bun:"user_id" json:"user_id"`
	Title          string     `bun:"title" json:"title"`
	Message        string     `bun:"message" json:"message"`
	Type           string     `bun:"type" json:"type"`
	Channel        string     `bun:"channel,default:system" json:"channel"`
	IsRead         bool       `bun:"is_read,default:false" json:"is_read"`
	WhatsappID     *string    `bun:"whatsapp_id" json:"whatsapp_id"`
	DeliveryStatus string     `bun:"delivery_status,default:pending" json:"delivery_status"`
	DeliveryError  *string    `bun:"delivery_error" json:"delivery_error,omitempty"`
	CreatedAt      time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt      *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at"`
}
