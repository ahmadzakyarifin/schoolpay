package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type BackgroundJob struct {
	bun.BaseModel `bun:"table:background_jobs,alias:bj"`

	ID           uint      `bun:"id,pk,autoincrement" json:"id"`
	TaskName     string    `bun:"task_name" json:"task_name"`
	Payload      string    `bun:"payload" json:"payload"`
	Status       string    `bun:"status,default:pending" json:"status"` // pending, processing, completed, failed
	Attempts     int       `bun:"attempts,default:0" json:"attempts"`
	MaxAttempts  int       `bun:"max_attempts,default:3" json:"max_attempts"`
	ErrorMessage *string   `bun:"error_message" json:"error_message,omitempty"`
	CreatedAt    time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	ScheduledAt  time.Time `bun:"scheduled_at,nullzero,notnull,default:current_timestamp" json:"scheduled_at"`
}
