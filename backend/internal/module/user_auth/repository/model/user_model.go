package model

import (
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/uptrace/bun"
)

type UserModel struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	ID           uint              `bun:"id,pk,autoincrement"`
	Name         string            `bun:"name"`
	ImagePath    *string           `bun:"image_path"`
	Email        string            `bun:"email,unique"`
	PhoneNumber  string            `bun:"phone_number,unique"`
	PasswordHash string            `bun:"password_hash,nullzero"`
	Role         string            `bun:"role"`
	NIK          *string           `bun:"nik,unique"`
	BirthDate    *time.Time        `bun:"birth_date"`
	Address      *string           `bun:"address"`
	Education    *string           `bun:"education"`
	Occupation   *string           `bun:"occupation"`
	Income       *string           `bun:"income"`
	IsActive     bool              `bun:"is_active,default:true"`
	Relation     string            `bun:",scanonly"`
	StudentCount int               `bun:",scanonly"`
	StudentNames *string           `bun:",scanonly"`
	CreatedAt    time.Time         `bun:"created_at,nullzero,notnull,default:current_timestamp"`
	UpdatedAt    time.Time         `bun:"updated_at,nullzero,notnull,default:current_timestamp"`
	DeletedAt    *time.Time        `bun:"deleted_at,soft_delete,nullzero"`
}

func (m *UserModel) ToDomain() *domain.User {
	if m == nil {
		return nil
	}
	return &domain.User{
		ID:           m.ID,
		Name:         m.Name,
		ImagePath:    m.ImagePath,
		Email:        m.Email,
		PhoneNumber:  m.PhoneNumber,
		PasswordHash: m.PasswordHash,
		Role:         m.Role,
		NIK:          m.NIK,
		BirthDate:    m.BirthDate,
		Address:      m.Address,
		Education:    m.Education,
		Occupation:   m.Occupation,
		Income:       m.Income,
		IsActive:     m.IsActive,
		Relation:     m.Relation,
		StudentCount: m.StudentCount,
		StudentNames: m.StudentNames,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
		DeletedAt:    m.DeletedAt,
	}
}

func FromDomain(d *domain.User) *UserModel {
	if d == nil {
		return nil
	}
	return &UserModel{
		ID:           d.ID,
		Name:         d.Name,
		ImagePath:    d.ImagePath,
		Email:        d.Email,
		PhoneNumber:  d.PhoneNumber,
		PasswordHash: d.PasswordHash,
		Role:         d.Role,
		NIK:          d.NIK,
		BirthDate:    d.BirthDate,
		Address:      d.Address,
		Education:    d.Education,
		Occupation:   d.Occupation,
		Income:       d.Income,
		IsActive:     d.IsActive,
		Relation:     d.Relation,
		StudentCount: d.StudentCount,
		StudentNames: d.StudentNames,
		CreatedAt:    d.CreatedAt,
		UpdatedAt:    d.UpdatedAt,
		DeletedAt:    d.DeletedAt,
	}
}
