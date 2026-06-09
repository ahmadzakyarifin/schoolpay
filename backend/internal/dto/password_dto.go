package dto

type ForgotPasswordRequest struct {
	Email          string `json:"email" binding:"required,email"`
	CaptchaToken   string `json:"captcha_token"`
	TurnstileToken string `json:"turnstile_token"`
}

type ResetPasswordRequest struct {
	Token           string `json:"token" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
	CaptchaToken    string `json:"captcha_token"`
	TurnstileToken  string `json:"turnstile_token"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" binding:"required"`
	Password        string `json:"password" binding:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" binding:"required,eqfield=Password"`
}
