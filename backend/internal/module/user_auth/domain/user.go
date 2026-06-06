package domain

import (
	"time"
)

type User struct {
	ID           uint
	Name         string
	ImagePath    *string
	Email        string
	PhoneNumber  string
	PasswordHash string
	Role         string
	NIK          *string
	BirthDate    *time.Time
	Address      *string
	Education    *string
	Occupation   *string
	Income       *string
	IsActive     bool
	Relation     string
	StudentCount int
	StudentNames *string
	StudentIDs   []uint
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    *time.Time
}
