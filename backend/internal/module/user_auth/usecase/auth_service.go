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
	notificationdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/google/uuid"
	"github.com/uptrace/bun"
)

const TaskAuthEmail = "email:auth"

type AuthService interface {
	Login(ctx context.Context, req dto.LoginRequest, audit dto.AuditMeta, jwtSecret string) (*dto.LoginResponse, error)
	RefreshToken(ctx context.Context, refreshToken string, audit dto.AuditMeta, jwtSecret string) (*dto.LoginResponse, error)
	ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest, audit dto.AuditMeta, cfg *config.Config) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest, audit dto.AuditMeta) error
	ChangePassword(ctx context.Context, userID uint, req dto.ChangePasswordRequest, audit dto.AuditMeta) error
	Logout(ctx context.Context, refreshToken string) error
}

type authService struct {
	r         repository.AuthRepo
	messenger utils.Messenger
	audit     auditusecase.AuditLogService
}

type AuthEmailJob struct {
	UserID  uint   `json:"user_id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func NewAuthService(repo repository.AuthRepo, msg utils.Messenger, audit auditusecase.AuditLogService) AuthService {
	s := &authService{
		r:         repo,
		messenger: msg,
		audit:     audit,
	}
	return s
}

func (s *authService) Login(ctx context.Context, req dto.LoginRequest, audit dto.AuditMeta, jwtSecret string) (*dto.LoginResponse, error) {
	user, err := s.r.FindByEmail(ctx, req.Email)
	if err != nil {
		if s.audit != nil {
			_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "LOGIN_FAILED", "auth", 0, nil, map[string]interface{}{"email": req.Email}, "user_not_found_or_error")
		}
		return nil, errors.New("email atau password salah")
	}

	// Set UserID ke audit meta setelah user ditemukan
	audit.UserID = &user.ID

	// Cek status aktif dulu agar user tahu akunnya dicekal sebelum cek password
	if !user.IsActive {
		if s.audit != nil {
			_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "LOGIN_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, "user_inactive")
		}
		return nil, errors.New("akun Anda belum aktif atau telah dinonaktifkan. Silakan hubungi Admin.")
	}

	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		if s.audit != nil {
			_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "LOGIN_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, "incorrect_password")
		}
		return nil, errors.New("email atau password salah")
	}

	accessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role, jwtSecret)
	if err != nil {
		return nil, errors.New("gagal membuat access token")
	}
	refreshToken, expiry, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, errors.New("gagal membuat refresh token")
	}

	if err := s.r.SaveRefreshToken(ctx, user.ID, refreshToken, expiry); err != nil {
		return nil, errors.New("gagal menyimpan session")
	}

	if s.audit != nil {
		_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "LOGIN", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, "user berhasil login")
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

func (s *authService) RefreshToken(ctx context.Context, refreshToken string, audit dto.AuditMeta, jwtSecret string) (*dto.LoginResponse, error) {
	user, err := s.r.FindUserByRefreshToken(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("session tidak valid atau sudah kadaluarsa")
	}

	if !user.IsActive {
		return nil, errors.New("akun Anda telah dinonaktifkan. Silakan hubungi Admin.")
	}

	newAccessToken, err := utils.GenerateAccessToken(user.ID, user.Email, user.Role, jwtSecret)
	if err != nil {
		return nil, errors.New("gagal membuat access token baru")
	}
	newRefreshToken, expiry, err := utils.GenerateRefreshToken()
	if err != nil {
		return nil, errors.New("gagal membuat refresh token baru")
	}

	err = s.r.GetDB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		repoTx := s.r.WithTx(tx)
		if err := repoTx.RotateRefreshToken(ctx, refreshToken); err != nil {
			return err
		}
		if err := repoTx.SaveRefreshToken(ctx, user.ID, newRefreshToken, expiry); err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return nil, errors.New("gagal menyimpan session")
	}

	// Set UserID ke audit meta
	audit.UserID = &user.ID

	if s.audit != nil {
		_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "REFRESH_TOKEN", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, "refresh token berhasil")
	}

	return &dto.LoginResponse{
		AccessToken:        newAccessToken,
		RefreshToken:       newRefreshToken,
		RefreshTokenExpiry: expiry,
	}, nil
}

func (s *authService) ForgotPassword(ctx context.Context, req dto.ForgotPasswordRequest, audit dto.AuditMeta, cfg *config.Config) error {
	email := strings.TrimSpace(req.Email)
	user, err := s.r.FindByEmail(ctx, email)
	if err != nil {
		if s.audit != nil {
			_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "FORGOT_PASSWORD_FAILED", "auth", 0, nil, map[string]interface{}{"email": req.Email}, "email_not_found_or_error")
		}
		return nil
	}

	if !user.IsActive {
		if s.audit != nil {
			_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "FORGOT_PASSWORD_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, "user_inactive")
		}

		body := fmt.Sprintf(`
		<div style="font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; max-width: 600px; margin: 0 auto; padding: 40px 20px; text-align: center; background-color: #f8fafc; border-radius: 16px;">
			<div style="background-color: white; padding: 40px; border-radius: 12px; box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.05);">
				<h2 style="color: #e11d48; margin-top: 0; font-size: 24px;">Akun SchoolPay Dinonaktifkan</h2>
				<p style="color: #64748b; font-size: 16px; line-height: 1.6; margin-bottom: 30px; text-align: left;">
					Halo <strong>%s</strong>,<br><br>
					Kami mendeteksi adanya permintaan untuk menyetel ulang password akun SchoolPay Anda.<br><br>
					Namun, saat ini status akun Anda sedang <strong>tidak aktif / dinonaktifkan</strong>. Oleh karena itu, permintaan reset password tidak dapat diproses.<br><br>
					Jika Anda tidak melakukan permintaan ini, Anda dapat mengabaikan email ini. Jika Anda ingin mengaktifkan kembali akun Anda, silakan hubungi Administrator sekolah atau tim bantuan SchoolPay.
				</p>
			</div>
		</div>
		`, user.Name)

		job := AuthEmailJob{
			UserID:  user.ID,
			Email:   user.Email,
			Name:    user.Name,
			Subject: "Permintaan Reset Password - Akun Dinonaktifkan",
			Body:    body,
		}

		payload, err := json.Marshal(job)
		if err == nil {
			bj := &notificationdomain.BackgroundJob{
				TaskName: TaskAuthEmail,
				Payload:  string(payload),
				Status:   "pending",
			}
			if _, err := s.r.GetDB().NewInsert().Model(bj).Exec(ctx); err != nil {
				fmt.Printf("Error inserting background job: %v\n", err)
			}
		}

		return nil
	}

	audit.UserID = &user.ID

	token := uuid.New().String()
	// sisa waktu link hangus
	expiresAt := time.Now().Add(15 * time.Minute)

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
		_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "FORGOT_PASSWORD", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, "forgot password link generated")
	}

	job := AuthEmailJob{
		UserID:  user.ID,
		Email:   user.Email,
		Name:    user.Name,
		Subject: "Reset Password - SchoolPay",
		Body:    body,
	}
	payload, err := json.Marshal(job)
	if err != nil {
		return errors.New("gagal memproses permintaan email")
	}

	bj := &notificationdomain.BackgroundJob{
		TaskName: TaskAuthEmail,
		Payload:  string(payload),
		Status:   "pending",
	}
	if _, err := s.r.GetDB().NewInsert().Model(bj).Exec(ctx); err != nil {
		return errors.New("gagal menjadwalkan pengiriman email")
	}

	return nil
}

func (s *authService) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest, audit dto.AuditMeta) error {
	userID, err := s.r.FindAuthToken(ctx, req.Token, "reset_password")
	if err != nil {
		if s.audit != nil {
			_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "RESET_PASSWORD_FAILED", "auth", 0, nil, map[string]interface{}{"token": req.Token}, "invalid_token_or_expired")
		}
		return errors.New("token tidak valid")
	}

	// Set UserID ke audit meta
	audit.UserID = &userID

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		if s.audit != nil {
			_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "RESET_PASSWORD_FAILED", "auth", userID, nil, nil, "password_hash_failed")
		}
		return errors.New("gagal memproses password baru")
	}
	
	err = s.r.GetDB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		repoTx := s.r.WithTx(tx)
		if err := repoTx.UpdatePassword(ctx, userID, hashed); err != nil {
			return err
		}
		_ = repoTx.DeleteAuthToken(ctx, req.Token, "reset_password")
		_ = repoTx.DeleteAllUserRefreshTokens(ctx, userID)
		return nil
	})

	if err != nil {
		if s.audit != nil {
			_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "RESET_PASSWORD_FAILED", "auth", userID, nil, nil, "database_update_failed")
		}
		return errors.New("gagal memperbarui password")
	}

	user, _ := s.r.FindUserByID(ctx, userID)
	if s.audit != nil && user != nil {
		_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "RESET_PASSWORD", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, "password reset successfully")
	}

	return nil
}

func (s *authService) ChangePassword(ctx context.Context, userID uint, req dto.ChangePasswordRequest, audit dto.AuditMeta) error {
	user, err := s.r.FindUserByID(ctx, userID)
	if err != nil || user == nil {
		return errors.New("user tidak ditemukan")
	}

	// Set UserID ke audit meta
	audit.UserID = &userID

	if !utils.CheckPassword(req.CurrentPassword, user.PasswordHash) {
		if s.audit != nil {
			_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "CHANGE_PASSWORD_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, "incorrect_current_password")
		}
		return errors.New("password saat ini tidak cocok")
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		if s.audit != nil {
			_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "CHANGE_PASSWORD_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, "password_hash_failed")
		}
		return errors.New("gagal memproses password baru")
	}
	
	err = s.r.GetDB().RunInTx(ctx, nil, func(ctx context.Context, tx bun.Tx) error {
		repoTx := s.r.WithTx(tx)
		if err := repoTx.UpdatePassword(ctx, userID, hashed); err != nil {
			return err
		}
		return repoTx.DeleteAllUserRefreshTokens(ctx, userID)
	})
	
	if err == nil {
		if s.audit != nil {
			_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "CHANGE_PASSWORD", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, "password changed successfully")
		}
	} else {
		if s.audit != nil {
			_ = s.audit.LogMeta(ctx, s.r.GetDB(), audit, "CHANGE_PASSWORD_FAILED", "auth", user.ID, nil, map[string]interface{}{"email": user.Email}, "database_update_failed")
		}
	}
	return err
}

func (s *authService) Logout(ctx context.Context, refreshToken string) error {
	return s.r.DeleteRefreshToken(ctx, refreshToken)
}
