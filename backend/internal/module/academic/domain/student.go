package domain

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
)

type Student struct {
	bun.BaseModel `bun:"table:students,alias:s"`

	ID             uint             `bun:"id,pk,autoincrement" json:"id" form:"id"`
	NIK            string           `bun:"nik,unique" json:"nik" form:"nik" binding:"required,custom_nik"`
	NIS            string           `bun:"nis,unique,nullzero" json:"nis" form:"nis" binding:"omitempty,min=3,max=20"`
	NISN           string           `bun:"nisn,unique" json:"nisn" form:"nisn" binding:"required,min=10,max=10"`
	ParentID       uint             `bun:"parent_id" json:"parent_id" form:"parent_id" binding:"required"`
	ClassID        uint             `bun:"class_id" json:"class_id" form:"class_id" binding:"required"`
	MajorID        uint             `bun:"major_id" json:"major_id" form:"major_id" binding:"required"`
	Name           string           `bun:"name" json:"name" form:"name" binding:"required,min=2"`
	Gender         string           `bun:"gender" json:"gender" form:"gender" binding:"required"`
	BirthPlace     string           `bun:"birth_place" json:"birth_place" form:"birth_place" binding:"required"`
	BirthDate      utils.CustomDate `bun:"birth_date" json:"birth_date" form:"birth_date" binding:"required"`
	Religion       string           `bun:"religion" json:"religion" form:"religion" binding:"required"`
	Address        string           `bun:"address" json:"address" form:"address" binding:"omitempty"`
	RT             string           `bun:"rt" json:"rt" form:"rt" binding:"required,min=1,max=5"`
	RW             string           `bun:"rw" json:"rw" form:"rw" binding:"required,min=1,max=5"`
	Village        string           `bun:"village" json:"village" form:"village" binding:"required"`
	District       string           `bun:"district" json:"district" form:"district" binding:"required"`
	City           string           `bun:"city" json:"city" form:"city" binding:"required"`
	Province       string           `bun:"province" json:"province" form:"province" binding:"required"`
	PhoneNumber    string           `bun:"phone_number" json:"phone_number" form:"phone_number" binding:"required,min=9,max=15"`
	EntryYear      int              `bun:"entry_year" json:"entry_year" form:"entry_year" binding:"required"`
	Email          string           `bun:"email,unique" json:"email" form:"email" binding:"required,email"`
	ImagePath      *string          `bun:"image_path" json:"image_path" form:"-"`
	Status         string           `bun:"status,default:active" json:"status" form:"status" binding:"required"`
	Description    string           `bun:"description" json:"description" form:"description" binding:"omitempty"`
	DepositBalance float64          `bun:"deposit_balance,default:0" json:"deposit_balance" form:"-"`
	CreatedAt      time.Time        `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time        `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt      *time.Time       `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at"`

	// Join Fields (not in DB table)
	ParentName string  `bun:",scanonly" json:"parent_name,omitempty"`
	ClassName  *string `bun:",scanonly" json:"class_name"`
	MajorName  *string `bun:",scanonly" json:"major_name"`
}
