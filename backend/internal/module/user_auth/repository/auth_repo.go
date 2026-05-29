package repository

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

type AuthRepo interface {
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindUserByID(ctx context.Context, id uint) (*domain.User, error)
	UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error

	// Refresh Tokens
	SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error
	FindRefreshToken(ctx context.Context, token string) (uint, error)
	DeleteRefreshToken(ctx context.Context, token string) error
	DeleteAllUserRefreshTokens(ctx context.Context, userID uint) error

	// Reset Password
	SaveAuthToken(ctx context.Context, userID uint, token string, tokenType string, expiresAt time.Time) error
	FindAuthToken(ctx context.Context, token string, tokenType string) (uint, error)
	MarkTokenAsUsed(ctx context.Context, token string, tokenType string) error
	DeleteAuthToken(ctx context.Context, token string, tokenType string) error
	GetDB() *bun.DB
}

type authRepo struct {
	db  *bun.DB
	rdb *redis.Client
}

func NewAuthRepo(db *bun.DB, rdb *redis.Client) AuthRepo {
	return &authRepo{db: db, rdb: rdb}
}

func (r *authRepo) GetDB() *bun.DB {
	return r.db
}

func (r *authRepo) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User
	err := r.db.NewSelect().Model(&user).Where("email = ?", email).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepo) FindUserByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := r.db.NewSelect().Model(&user).Where("id = ?", id).Scan(ctx)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepo) SaveRefreshToken(ctx context.Context, userID uint, token string, expiresAt time.Time) error {
	duration := time.Until(expiresAt)
	if duration <= 0 {
		return fmt.Errorf("token already expired")
	}

	pipe := r.rdb.Pipeline()
	pipe.Set(ctx, "refresh_token:"+token, userID, duration)
	userSetKey := fmt.Sprintf("user_refresh_tokens:%d", userID)
	pipe.SAdd(ctx, userSetKey, token)
	pipe.Expire(ctx, userSetKey, duration)
	_, err := pipe.Exec(ctx)
	return err
}

func (r *authRepo) FindRefreshToken(ctx context.Context, token string) (uint, error) {
	userIDStr, err := r.rdb.Get(ctx, "refresh_token:"+token).Result()
	if err != nil {
		return 0, err // redis.Nil means not found
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, err
	}
	return uint(userID), nil
}

func (r *authRepo) DeleteRefreshToken(ctx context.Context, token string) error {
	userIDStr, err := r.rdb.Get(ctx, "refresh_token:"+token).Result()
	if err == nil {
		r.rdb.SRem(ctx, fmt.Sprintf("user_refresh_tokens:%s", userIDStr), token)
	}
	return r.rdb.Del(ctx, "refresh_token:"+token).Err()
}

func (r *authRepo) DeleteAllUserRefreshTokens(ctx context.Context, userID uint) error {
	userSetKey := fmt.Sprintf("user_refresh_tokens:%d", userID)
	tokens, err := r.rdb.SMembers(ctx, userSetKey).Result()
	if err != nil {
		return err
	}

	if len(tokens) > 0 {
		var keys []string
		for _, t := range tokens {
			keys = append(keys, "refresh_token:"+t)
		}
		r.rdb.Del(ctx, keys...)
	}
	return r.rdb.Del(ctx, userSetKey).Err()
}

func (r *authRepo) UpdatePassword(ctx context.Context, userID uint, hashedPassword string) error {
	_, err := r.db.NewUpdate().
		Model((*domain.User)(nil)).
		Set("password_hash = ?", hashedPassword).
		Where("id = ?", userID).
		Exec(ctx)
	return err
}

func (r *authRepo) SaveAuthToken(ctx context.Context, userID uint, token string, tokenType string, expiresAt time.Time) error {
	duration := time.Until(expiresAt)
	if duration <= 0 {
		return fmt.Errorf("token already expired")
	}
	// Store token in Redis: key="auth_token:<type>:<token>", value="<userID>"
	key := fmt.Sprintf("auth_token:%s:%s", tokenType, token)
	return r.rdb.Set(ctx, key, userID, duration).Err()
}

func (r *authRepo) FindAuthToken(ctx context.Context, token string, tokenType string) (uint, error) {
	key := fmt.Sprintf("auth_token:%s:%s", tokenType, token)
	userIDStr, err := r.rdb.Get(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		return 0, err
	}
	return uint(userID), nil
}

func (r *authRepo) MarkTokenAsUsed(ctx context.Context, token string, tokenType string) error {
	r.rdb.Del(ctx, fmt.Sprintf("auth_token:%s:%s", tokenType, token))
	return nil
}

func (r *authRepo) DeleteAuthToken(ctx context.Context, token string, tokenType string) error {
	return r.rdb.Del(ctx, fmt.Sprintf("auth_token:%s:%s", tokenType, token)).Err()
}
