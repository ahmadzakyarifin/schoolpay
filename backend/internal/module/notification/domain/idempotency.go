package domain

import (
	"time"

	"github.com/uptrace/bun"
)

const (
	IdempotencyStatusProcessing = "PROCESSING"
	IdempotencyStatusCompleted  = "COMPLETED"
)

type IdempotencyKey struct {
	bun.BaseModel `bun:"table:idempotency_keys,alias:ik"`

	Key             string    `bun:"key,pk" json:"key"`
	Status          string    `bun:"status,notnull,default:PROCESSING" json:"status"`
	RequestHash     string    `bun:"request_hash,notnull" json:"request_hash"`
	ResponsePayload string    `bun:"response_payload,nullzero" json:"response_payload"`
	CreatedAt       time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt       time.Time `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
}
