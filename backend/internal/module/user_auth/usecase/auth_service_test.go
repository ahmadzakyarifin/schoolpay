package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/dto"
	"github.com/ahmadzakyarifin/schoolpay/internal/mocks"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/usecase"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_Login_Success(t *testing.T) {
	repoMock := mocks.NewAuthRepo(t)
	auditMock := mocks.NewAuditLogService(t)

	svc := usecase.NewAuthService(repoMock, nil, auditMock)

	ctx := context.Background()
	password := "password123"
	hashedPassword, _ := utils.HashPassword(password)

	req := dto.LoginRequest{
		Email:    "test@schoolpay.id",
		Password: password,
	}

	user := &domain.User{
		ID:           1,
		Name:         "Test User",
		Email:        req.Email,
		Role:         "admin",
		IsActive:     true,
		PasswordHash: hashedPassword,
	}

	repoMock.On("FindByEmail", ctx, req.Email).Return(user, nil)
	repoMock.On("SaveRefreshToken", ctx, user.ID, mock.AnythingOfType("string"), mock.AnythingOfType("time.Time")).Return(nil)
	repoMock.On("GetDB").Return(nil)

	// Mock audit log
	auditMock.On("LogMeta", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)

	resp, err := svc.Login(ctx, req, dto.AuditMeta{}, "secret_key")

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, user.ID, resp.User.ID)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}

func TestAuthService_Login_InvalidEmail(t *testing.T) {
	repoMock := mocks.NewAuthRepo(t)
	svc := usecase.NewAuthService(repoMock, nil, nil)

	ctx := context.Background()
	req := dto.LoginRequest{
		Email:    "wrong@schoolpay.id",
		Password: "password123",
	}

	repoMock.On("FindByEmail", ctx, req.Email).Return(nil, errors.New("not found"))

	resp, err := svc.Login(ctx, req, dto.AuditMeta{}, "secret_key")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Equal(t, "email atau password salah", err.Error())
}

func TestAuthService_Login_InactiveUser(t *testing.T) {
	repoMock := mocks.NewAuthRepo(t)
	svc := usecase.NewAuthService(repoMock, nil, nil)

	ctx := context.Background()
	req := dto.LoginRequest{
		Email:    "inactive@schoolpay.id",
		Password: "password123",
	}

	user := &domain.User{
		ID:       2,
		Email:    req.Email,
		IsActive: false,
	}

	repoMock.On("FindByEmail", ctx, req.Email).Return(user, nil)

	resp, err := svc.Login(ctx, req, dto.AuditMeta{}, "secret_key")

	assert.Error(t, err)
	assert.Nil(t, resp)
	assert.Contains(t, err.Error(), "akun Anda belum aktif")
}
