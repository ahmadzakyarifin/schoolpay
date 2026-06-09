package dto

import "time"

type LoginRequest struct {
	Email          string `json:"email" binding:"required,email"`
	Password       string `json:"password" binding:"required,min=6"`
	CaptchaToken   string `json:"captcha_token"`
	TurnstileToken string `json:"turnstile_token"`
}

type LoginResponse struct {
	AccessToken        string    `json:"access_token"`
	RefreshToken       string    `json:"-"`
	RefreshTokenExpiry time.Time `json:"-"`
	User               UserInfo  `json:"user"`
}

type UserInfo struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}
