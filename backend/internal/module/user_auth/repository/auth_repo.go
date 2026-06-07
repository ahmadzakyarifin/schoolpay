package repository

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository/model"
	"github.com/uptrace/bun"
)

type AuthRepo interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindUserByID(ctx context.Context, id uint) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error

	// Refresh Tokens
	SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error
	FindUserByRefreshToken(ctx context.Context, token string) (*domain.User, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteAllUserRefreshTokens(ctx context.Context, userID uint) error

	// Reset / Activation Tokens
	SaveAuthToken(ctx context.Context, userID uint, token string, tokenType string, expiresAt time.Time) error
	FindAuthToken(ctx context.Context, token string, tokenType string) (uint, error)
	DeleteAuthToken(ctx context.Context, token string, tokenType string) error
	GetDB() bun.IDB
	WithTx(tx bun.Tx) AuthRepo
	RotateRefreshToken(ctx context.Context, token string) error
}

type authRepo struct {
	db bun.IDB
}

type RefreshTokenModel struct {
	bun.BaseModel `bun:"table:refresh_tokens"`

	ID        uint      `bun:"id,pk,autoincrement"`
	UserID    uint      `bun:"user_id"`
	Token     string    `bun:"token"`
	ExpiresAt time.Time `bun:"expires_at"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,default:current_timestamp"`
}

type AuthTokenModel struct {
	bun.BaseModel `bun:"table:auth_tokens"`

	ID        uint      `bun:"id,pk,autoincrement"`
	UserID    uint      `bun:"user_id"`
	Token     string    `bun:"token"`
	Type      string    `bun:"type"`
	ExpiresAt time.Time `bun:"expires_at"`
	CreatedAt time.Time `bun:"created_at,default:current_timestamp"`
	UpdatedAt time.Time `bun:"updated_at,default:current_timestamp"`
}

func hashToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

func NewAuthRepo(db bun.IDB) AuthRepo {
	return &authRepo{db: db}
}

func (r *authRepo) WithTx(tx bun.Tx) AuthRepo {
	return &authRepo{db: tx}
}

func (r *authRepo) GetDB() bun.IDB {
	return r.db
}

func (r *authRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var u model.UserModel
	err := r.db.NewSelect().Model(&u).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}

func (r *authRepo) FindUserByID(ctx context.Context, id uint) (*domain.User, error) {
	var u model.UserModel
	err := r.db.NewSelect().Model(&u).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}

func (r *authRepo) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	hashedToken := hashToken(token)

	_, _ = r.db.NewDelete().Model((*RefreshTokenModel)(nil)).
		Where("user_id = ? AND expires_at < ?", userID, time.Now()).
		Exec(ctx)

	rt := &RefreshTokenModel{
		UserID:    userID,
		Token:     hashedToken,
		ExpiresAt: expiresAt,
	}

	_, err := r.db.NewInsert().Model(rt).Exec(ctx)
	return err
}

func (r *authRepo) FindUserByRefreshToken(ctx context.Context, token string) (*domain.User, error) {
	hashedToken := hashToken(token)
	var u model.UserModel
	err := r.db.NewSelect().Model(&u).
		Join("JOIN refresh_tokens AS rt ON rt.user_id = user_model.id").
		Where("rt.token = ? AND rt.expires_at > ?", hashedToken, time.Now()).
		Scan(ctx)
	if err != nil {
		return nil, err
	}
	return u.ToDomain(), nil
}

func (r *authRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	hashedToken := hashToken(token)
	_, err := r.db.NewDelete().Model((*RefreshTokenModel)(nil)).
		Where("token = ?", hashedToken).
		Exec(ctx)
	return err
}

func (r *authRepo) RotateRefreshToken(ctx context.Context, token string) error {
	hashedToken := hashToken(token)
	// Grace period: Token is still valid for 30 seconds
	gracePeriod := time.Now().Add(30 * time.Second)
	
	_, err := r.db.NewUpdate().Model((*RefreshTokenModel)(nil)).
		Set("expires_at = ?", gracePeriod).
		Where("token = ?", hashedToken).
		Exec(ctx)
	return err
}

func (r *authRepo) DeleteAllUserRefreshTokens(ctx context.Context, userID uint) error {
	_, err := r.db.NewDelete().Model((*RefreshTokenModel)(nil)).
		Where("user_id = ?", userID).
		Exec(ctx)
	return err
}

func (r *authRepo) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	_, err := r.db.NewUpdate().
		Model((*model.UserModel)(nil)).
		Set("password_hash = ?", hashedPassword).
		Where("id = ?", userID).
		Exec(ctx)
	return err
}

func (r *authRepo) SaveAuthToken(ctx context.Context, userID uint, token string, tokenType string, expiresAt time.Time) error {
	hashedToken := hashToken(token)
	at := &AuthTokenModel{
		UserID:    userID,
		Token:     hashedToken,
		Type:      tokenType,
		ExpiresAt: expiresAt,
	}
	_, err := r.db.NewInsert().Model(at).Exec(ctx)
	return err
}

func (r *authRepo) FindAuthToken(ctx context.Context, token string, tokenType string) (uint, error) {
	hashedToken := hashToken(token)
	var at AuthTokenModel
	err := r.db.NewSelect().Model(&at).
		Where("token = ? AND type = ? AND expires_at > ?", hashedToken, tokenType, time.Now()).
		Scan(ctx)
	if err != nil {
		return 0, err
	}
	return at.UserID, nil
}

func (r *authRepo) DeleteAuthToken(ctx context.Context, token string, tokenType string) error {
	hashedToken := hashToken(token)
	_, err := r.db.NewDelete().Model((*AuthTokenModel)(nil)).
		Where("token = ? AND type = ?", hashedToken, tokenType).
		Exec(ctx)
	return err
}
