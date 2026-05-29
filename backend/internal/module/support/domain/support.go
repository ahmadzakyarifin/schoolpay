package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Conversation struct {
	bun.BaseModel `bun:"table:support_conversations,alias:sc"`

	ID              uint       `bun:"id,pk,autoincrement" json:"id"`
	ParentID        *uint      `bun:"parent_id" json:"parent_id,omitempty"`
	PhoneNumber     string     `bun:"phone_number" json:"phone_number"`
	ParentName      *string    `bun:"parent_name" json:"parent_name,omitempty"`
	Status          string     `bun:"status,default:open" json:"status"`
	AssignedAdminID *uint      `bun:"assigned_admin_id" json:"assigned_admin_id,omitempty"`
	Subject         *string    `bun:"subject" json:"subject,omitempty"`
	UnreadCount     int        `bun:"unread_count" json:"unread_count"`
	LastMessage     *string    `bun:"last_message" json:"last_message,omitempty"`
	LastMessageAt   *time.Time `bun:"last_message_at" json:"last_message_at,omitempty"`
	CreatedAt       time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	ClosedAt        *time.Time `bun:"closed_at" json:"closed_at,omitempty"`

	Messages []Message `bun:"rel:has-many,join:id=conversation_id" json:"messages,omitempty"`
}

type Message struct {
	bun.BaseModel `bun:"table:support_messages,alias:sm"`

	ID                uint      `bun:"id,pk,autoincrement" json:"id"`
	ConversationID    uint      `bun:"conversation_id" json:"conversation_id"`
	SenderType        string    `bun:"sender_type" json:"sender_type"`
	SenderID          *uint     `bun:"sender_id" json:"sender_id,omitempty"`
	Message           string    `bun:"message" json:"message"`
	WhatsappMessageID *string   `bun:"whatsapp_message_id" json:"whatsapp_message_id,omitempty"`
	DeliveryStatus    string    `bun:"delivery_status,default:received" json:"delivery_status"`
	CreatedAt         time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}
