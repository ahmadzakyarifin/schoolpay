package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/dto"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/hibiken/asynq"
)

const TaskAuthEmail = "email:auth"

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest, jwtSecret string) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, req dto.RefreshTokenRequest, jwtSecret string) (*dto.LoginResponse, error)
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest, cfg *config.Config) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	r           repository.AuthRepo
	messenger   utils.Messenger
	audit       auditusecase.AuditLogService
	asynqClient *asynq.Client
}

type AuthEmailJob struct {
	Email   string `json:"email"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewAuthService(repo repository.AuthRepo, msg utils.Messenger, audit auditusecase.AuditLogService, asynqClient *asynq.Client) AuthService {
	s := &authService{
		r:           repo,
		messenger:   msg,
		audit:       audit,
		asynqClient: asynqClient,
	}
	return s
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest, jwtSecret string) (*dto.LoginResponse, error) {
	user, err := s.r.FindByEmail(ctx, req.Email)
	if err != nil {
		if s.audit != nil {
			_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, s.r.GetDB(), 0, "Anonymous", "anonymous", "LOGIN_FAILED", "auth", 0, nil, map[string]interface{}{"email": req.Email, "reason": "user_not_found_or_error"}, ipAddress, userAgent)
		}
		return nil, errors.New("email atau password salah")
	}

	// Cek status aktif dulu agar user tahu akunnya dicekal sebelum cek password
	if !user.IsActive {
		if s.audit != nil {
			_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, s.r.GetDB(), user.ID, user.Name, user.Role, "LOGIN_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email, "reason": "user_inactive"}, ipAddress, userAgent)
		}
		return nil, errors.New("akun Anda belum aktif atau telah dinonaktifkan. Silakan hubungi Admin.")
	}

	passHash := ""
	if user.PasswordHash != nil {
		passHash = *user.PasswordHash
	}
	if !utils.CheckPassword(req.Password, passHash) {
		if s.audit != nil {
			_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, s.r.GetDB(), user.ID, user.Name, user.Role, "LOGIN_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email, "reason": "incorrect_password"}, ipAddress, userAgent)
		}
		return nil, errors.New("email atau password salah")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role, jwtSecret)
	if err != nil {
		return nil, errors.New("gagal membuat access token")
	}
	refreshToken, expiry, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role, jwtSecret)
	if err != nil {
		return nil, errors.New("gagal membuat refresh token")
	}

	if err := s.r.SaveRefreshToken(ctx, user.ID, refreshToken, expiry); err != nil {
		return nil, errors.New("gagal menyimpan session")
	}

	if s.audit != nil {
		_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.r.GetDB(), user.ID, user.Name, user.Role, "LOGIN", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, ipAddress, userAgent)
	}

	return &dto.LoginResponse{
		AccessToken:        accessToken,
		RefreshToken:       refreshToken,
		RefreshTokenExpiry: expiry,
		User: dto.UserInfo{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenRequest, jwtSecret string) (*dto.LoginResponse, error) {
	claims, err := utils.ValidateToken(req.RefreshToken, jwtSecret)
	if err != nil {
		return nil, errors.New("token tidak valid")
	}

	storedUserID, err := s.r.FindRefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.New("session tidak ditemukan")
	}
	if storedUserID != claims.UserID {
		return nil, errors.New("session tidak valid")
	}

	user, err := s.r.FindUserByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user tidak ditemukan")
	}

	if !user.IsActive {
		return nil, errors.New("akun Anda telah dinonaktifkan. Silakan hubungi Admin.")
	}

	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role, jwtSecret)
	if err != nil {
		return nil, errors.New("gagal membuat access token baru")
	}
	newRefreshToken, expiry, err := utils.GenerateRefreshToken(user.ID, user.Email, user.Role, jwtSecret)
	if err != nil {
		return nil, errors.New("gagal membuat refresh token baru")
	}

	_ = s.r.DeleteRefreshToken(ctx, req.RefreshToken)
	if err := s.r.SaveRefreshToken(ctx, user.ID, newRefreshToken, expiry); err != nil {
		return nil, errors.New("gagal menyimpan session")
	}

	if s.audit != nil {
		_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.r.GetDB(), user.ID, user.Name, user.Role, "REFRESH_TOKEN", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, ipAddress, userAgent)
	}

	return &dto.LoginResponse{
		AccessToken:        newAccessToken,
		RefreshToken:       newRefreshToken,
		RefreshTokenExpiry: expiry,
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest, cfg *config.Config) error {
	user, err := s.r.FindByEmail(ctx, req.Email)
	if err != nil {
		if s.audit != nil {
			_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, s.r.GetDB(), 0, "Anonymous", "anonymous", "FORGOT_PASSWORD_FAILED", "auth", 0, nil, map[string]interface{}{"email": req.Email, "reason": "email_not_found_or_error"}, ipAddress, userAgent)
		}
		return nil
	}

	const resetPasswordTokenTTL = 15 * time.Minute

	token := utils.GenerateUUID()
	expiresAt := time.Now().Add(resetPasswordTokenTTL)

	if err := s.r.SaveAuthToken(ctx, user.ID, token, "reset_password", expiresAt); err != nil {
		return errors.New("gagal memproses permintaan")
	}

	link := fmt.Sprintf("%s/reset-password?token=%s", strings.TrimSuffix(cfg.FrontendURL, "/"), token)
	body := fmt.Sprintf(`
	<div style="font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 20px; text-align: center; background-color: #f8fafc; border-radius: 16px;">
		<div style="background-color: white; padding: 40px; border-radius: 12px; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);">
			<h2 style="color: #1e293b; margin-top: 0; font-size: 24px;">Reset Password SchoolPay</h2>
			<p style="color: #64748b; font-size: 16px; line-height: 1.6; margin-bottom: 30px;">
				Halo <strong>%s</strong>,<br><br>
				Kami menerima permintaan untuk menyetel ulang password akun SchoolPay Anda. Silakan klik tombol di bawah ini untuk membuat password baru:
			</p>
			<a href="%s" style="display: inline-block; background-color: #4f46e5; color: #ffffff; padding: 16px 32px; text-decoration: none; border-radius: 12px; font-weight: bold; font-size: 16px; box-shadow: 0 4px 6px -1px rgba(79, 70, 229, 0.3);">
				Setel Ulang Password
			</a>
			<p style="color: #94a3b8; font-size: 13px; margin-top: 40px; margin-bottom: 0;">
				Jika tombol di atas tidak berfungsi, Anda juga dapat menyalin dan menempelkan tautan berikut ke browser Anda:<br>
				<a href="%s" style="color: #4f46e5; text-decoration: underline; word-break: break-all;">%s</a>
			</p>
		</div>
		<p style="color: #94a3b8; font-size: 12px; margin-top: 20px;">
			Jika Anda tidak merasa meminta reset password, abaikan saja email ini.<br>
			Tautan ini hanya berlaku selama 15 menit demi keamanan akun Anda.
		</p>
	</div>
	`, user.Name, link, link, link)

	if s.audit != nil {
		_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.r.GetDB(), user.ID, user.Name, user.Role, "FORGOT_PASSWORD", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, ipAddress, userAgent)
	}

	job := AuthEmailJob{
		Email:   user.Email,
		Name:    user.Name,
		Subject: "Reset Password - SchoolPay",
		Body:    body,
	}
	payload, err := json.Marshal(job)
	if err != nil {
		return errors.New("gagal memproses permintaan email")
	}
	if _, err := s.asynqClient.Enqueue(asynq.NewTask(TaskAuthEmail, payload, asynq.MaxRetry(3))); err != nil {
		return errors.New("gagal menjadwalkan pengiriman email")
	}

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	userID, err := s.r.FindAuthToken(ctx, req.Token, "reset_password")
	if err != nil {
		if s.audit != nil {
			_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, s.r.GetDB(), 0, "Anonymous", "anonymous", "RESET_PASSWORD_FAILED", "auth", 0, nil, map[string]interface{}{"token": req.Token, "reason": "invalid_token_or_expired"}, ipAddress, userAgent)
		}
		return errors.New("token tidak valid")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		if s.audit != nil {
			user, _ := s.r.FindUserByID(ctx, userID)
			if user != nil {
				_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
				_ = s.audit.Log(ctx, s.r.GetDB(), user.ID, user.Name, user.Role, "RESET_PASSWORD_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email, "reason": "password_hash_failed"}, ipAddress, userAgent)
			}
		}
		return errors.New("gagal memproses password baru")
	}
	if err := s.r.UpdatePassword(ctx, userID, hashed); err != nil {
		if s.audit != nil {
			user, _ := s.r.FindUserByID(ctx, userID)
			if user != nil {
				_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
				_ = s.audit.Log(ctx, s.r.GetDB(), user.ID, user.Name, user.Role, "RESET_PASSWORD_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email, "reason": "database_update_failed"}, ipAddress, userAgent)
			}
		}
		return errors.New("gagal memperbarui password")
	}

	user, _ := s.r.FindUserByID(ctx, userID)
	if s.audit != nil && user != nil {
		_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
		_ = s.audit.Log(ctx, s.r.GetDB(), user.ID, user.Name, user.Role, "RESET_PASSWORD", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, ipAddress, userAgent)
	}

	_ = s.r.MarkTokenAsUsed(ctx, req.Token, "reset_password")
	_ = s.r.DeleteAllUserRefreshTokens(ctx, userID)

	return nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword string) error {
	user, err := s.r.FindUserByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user tidak ditemukan")
	}

	passHash := ""
	if user.PasswordHash != nil {
		passHash = *user.PasswordHash
	}
	if !utils.CheckPassword(currentPassword, passHash) {
		if s.audit != nil {
			_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, s.r.GetDB(), user.ID, user.Name, user.Role, "CHANGE_PASSWORD_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email, "reason": "incorrect_current_password"}, ipAddress, userAgent)
		}
		return errors.New("password saat ini tidak cocok")
	}

	hashed, err := utils.HashPassword(newPassword)
	if err != nil {
		if s.audit != nil {
			_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, s.r.GetDB(), user.ID, user.Name, user.Role, "CHANGE_PASSWORD_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email, "reason": "password_hash_failed"}, ipAddress, userAgent)
		}
		return errors.New("gagal memproses password baru")
	}
	err = s.r.UpdatePassword(ctx, userID, hashed)
	if err == nil {
		_ = s.r.DeleteAllUserRefreshTokens(ctx, userID)
		if s.audit != nil {
			_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, s.r.GetDB(), user.ID, user.Name, user.Role, "CHANGE_PASSWORD", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, ipAddress, userAgent)
		}
	} else {
		if s.audit != nil {
			_, _, _, ipAddress, userAgent := utils.GetAuditMeta(ctx)
			_ = s.audit.Log(ctx, s.r.GetDB(), user.ID, user.Name, user.Role, "CHANGE_PASSWORD_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email, "reason": "database_update_failed"}, ipAddress, userAgent)
		}
	}
	return err
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.r.DeleteRefreshToken(ctx, refreshToken)
}
