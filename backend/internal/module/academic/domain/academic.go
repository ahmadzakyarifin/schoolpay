package domain

import (
	"time"

	"github.com/uptrace/bun"
)

type Major struct {
	bun.BaseModel `bun:"table:majors,alias:m"`

	ID        uint       `bun:"id,pk,autoincrement" json:"id"`
	Code      *string    `bun:"code" json:"code,omitempty"`
	Name      string     `bun:"name" json:"name" binding:"required"`
	IsActive  bool       `bun:"is_active,default:true" json:"is_active"`
	CreatedAt time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at,omitempty"`

	YearIDs           []uint   `bun:"-" json:"year_ids,omitempty"`
	AcademicYearCount int      `bun:"-" json:"academic_year_count"`
	ClassCount        int      `bun:"-" json:"class_count"`
	StudentCount      int      `bun:"-" json:"student_count"`
	ClassNames        []string `bun:"-" json:"class_names"`
}

type Class struct {
	bun.BaseModel `bun:"table:classes,alias:c"`

	ID             uint       `bun:"id,pk,autoincrement" json:"id"`
	Name           string     `bun:"name" json:"name" binding:"required"`
	MajorID        *uint      `bun:"major_id" json:"major_id" binding:"required"`
	AcademicYearID *uint      `bun:"academic_year_id" json:"academic_year_id"`
	Grade          int        `bun:"grade" json:"grade"`
	IsActive       bool       `bun:"is_active,default:true" json:"is_active"`
	CreatedAt      time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt      time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt      *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at,omitempty"`

	// Join
	MajorName        *string `bun:",scanonly" json:"major_name"`
	AcademicYearName *int    `bun:",scanonly" json:"academic_year_name"`
	AcademicYearIDs   []uint `bun:"-" json:"academic_year_ids,omitempty"`
	AcademicYearCount int    `bun:"-" json:"academic_year_count"`
	StudentCount      int    `bun:"-" json:"student_count"`
}

type AcademicYearMajor struct {
	bun.BaseModel `bun:"table:academic_year_majors,alias:aym"`

	AcademicYearID uint `bun:"academic_year_id,pk"`
	MajorID        uint `bun:"major_id,pk"`
}

type AcademicYearClass struct {
	bun.BaseModel `bun:"table:academic_year_classes,alias:ayc"`

	AcademicYearID uint `bun:"academic_year_id,pk"`
	ClassID        uint `bun:"class_id,pk"`
}

type StudentClass struct {
	bun.BaseModel `bun:"table:student_classes,alias:sc"`

	ID        uint      `bun:"id,pk,autoincrement" json:"id"`
	StudentID uint      `bun:"student_id" json:"student_id"`
	ClassID   uint      `bun:"class_id" json:"class_id"`
	IsActive  bool      `bun:"is_active,default:true" json:"is_active"`
	CreatedAt time.Time `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
}

type ClassHistory struct {
	bun.BaseModel `bun:"-"`

	ID           uint      `json:"id"`
	StudentID    uint      `json:"student_id"`
	ClassID      uint      `json:"class_id"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	ClassName    string    `json:"class_name"`
	Grade        string    `json:"grade"`
	AcademicYear string    `json:"academic_year"`
}

type AcademicYear struct {
	bun.BaseModel `bun:"table:academic_years,alias:ay"`

	ID        uint       `bun:"id,pk,autoincrement" json:"id"`
	Year      int        `bun:"year" json:"year" binding:"required"`
	IsActive  bool       `bun:"is_active,default:true" json:"is_active"`
	CreatedAt time.Time  `bun:"created_at,nullzero,notnull,default:current_timestamp" json:"created_at"`
	UpdatedAt time.Time  `bun:"updated_at,nullzero,notnull,default:current_timestamp" json:"updated_at"`
	DeletedAt *time.Time `bun:"deleted_at,soft_delete,nullzero" json:"deleted_at,omitempty"`

	// Summary Fields
	MajorCount   int      `bun:"-" json:"major_count"`
	ClassCount   int      `bun:"-" json:"class_count"`
	StudentCount int      `bun:"-" json:"student_count"`
	MajorNames   []string `bun:"-" json:"major_names"`
	MajorIDs   []uint   `bun:"-" json:"major_ids"`
	ClassIDs   []uint   `bun:"-" json:"class_ids"`
	ClassNames []string `bun:"-" json:"class_names"`
}
