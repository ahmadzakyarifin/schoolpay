package dto

import (
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
)

type UserResponse struct {
	ID           uint       `json:"id"`
	Name         string     `json:"name"`
	ImagePath    *string    `json:"image_path"`
	Email        string     `json:"email"`
	PhoneNumber  string     `json:"phone_number"`
	Role         string     `json:"role"`
	NIK          *string    `json:"nik"`
	BirthDate    *string    `json:"birth_date"`
	Address      *string    `json:"address"`
	Education    *string    `json:"education"`
	Occupation   *string    `json:"occupation"`
	Income       *string    `json:"income"`
	IsActive     bool       `json:"is_active"`
	HasPassword  bool       `json:"has_password"`
	Relation     string     `json:"relation,omitempty"`
	StudentCount int        `json:"student_count"`
	StudentNames *string    `json:"student_names"`
	StudentIDs   []uint     `json:"student_ids,omitempty"`
	CreatedAt    string     `json:"created_at"`
	UpdatedAt    string     `json:"updated_at"`
}

func ToUserResponse(user domain.User) UserResponse {
	var birthDate *string
	if user.BirthDate != nil && !user.BirthDate.IsZero() {
		s := user.BirthDate.Format("2006-01-02")
		birthDate = &s
	}
	return UserResponse{
		ID:           user.ID,
		Name:         user.Name,
		ImagePath:    user.ImagePath,
		Email:        user.Email,
		PhoneNumber:  user.PhoneNumber,
		Role:         user.Role,
		NIK:          user.NIK,
		BirthDate:    birthDate,
		Address:      user.Address,
		Education:    user.Education,
		Occupation:   user.Occupation,
		Income:       user.Income,
		IsActive:     user.IsActive,
		HasPassword:  user.PasswordHash != "",
		Relation:     user.Relation,
		StudentCount: user.StudentCount,
		StudentNames: user.StudentNames,
		StudentIDs:   user.StudentIDs,
		CreatedAt:    user.CreatedAt.Format("02/01/2006"),
		UpdatedAt:    user.UpdatedAt.Format("02/01/2006"),
	}
}

func ToUserResponseList(users []domain.User) []UserResponse {
	res := make([]UserResponse, len(users))
	for i, u := range users {
		res[i] = ToUserResponse(u)
	}
	return res
}
