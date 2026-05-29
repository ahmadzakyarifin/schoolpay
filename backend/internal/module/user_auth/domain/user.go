package domain

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           uint              `bun:"id,pk,autoincrement" json:"id"`
	Name         string            `bun:"name" json:"name" binding:"required,min=2"`
	ImagePath    *string           `bun:"image_path" json:"image_path"`
	Email        string            `bun:"email,unique" json:"email" binding:"required,email"`
	PhoneNumber  string            `bun:"phone_number,unique" json:"phone_number" binding:"required,min=9,max=15"`
	PasswordHash *string           `bun:"password_hash" json:"-"`
	Role         string            `bun:"role" json:"role" binding:"required,oneof=admin parent"`
	NIK          *string           `bun:"nik,unique" json:"nik" binding:"-"`
	BirthDate    *utils.CustomDate `bun:"birth_date" json:"birth_date" binding:"-"`
	Address      *string           `bun:"address" json:"address" binding:"-"`
	Education    *string           `bun:"education" json:"education" binding:"-"`
	Occupation   *string           `bun:"occupation" json:"occupation" binding:"-"`
	Income       *string           `bun:"income" json:"income" binding:"-"`
	IsActive     bool              `bun:"is_active,default:true" json:"is_active"`
	Relation     string            `bun:",scanonly" json:"relation,omitempty"`
	StudentCount int               `bun:",scanonly" json:"student_count"`
	StudentNames *string           `bun:",scanonly" json:"student_names"`
	StudentIDs   []uint            `bun:"-" json:"student_ids,omitempty"`
	CreatedAt    time.Time         `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt    time.Time         `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt    *time.Time        `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at,omitempty"`
}

