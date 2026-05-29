package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/redis/go-redis/v9"
	"github.com/uptrace/bun"
)

func TestNewAuthRepo(t *testing.T) {
	type args struct {
		db  *bun.DB
		rdb *redis.Client
	}
	tests := []struct {
		name string
		args args
		want AuthRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthRepo(tt.args.db, tt.args.rdb); got != tt.want {
				t.Errorf("NewAuthRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authRepo_GetDB(t *testing.T) {
	type fields struct {
		db  *bun.DB
		rdb *redis.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   *bun.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db:  tt.fields.db,
				rdb: tt.fields.rdb,
			}
			if got := r.GetDB(); got != tt.want {
				t.Errorf("authRepo.GetDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authRepo_FindByEmail(t *testing.T) {
	type fields struct {
		db  *bun.DB
		rdb *redis.Client
	}
	type args struct {
		ctx   context.Context
		email string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db:  tt.fields.db,
				rdb: tt.fields.rdb,
			}
			got, err := r.FindByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Fatalf("authRepo.FindByEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authRepo.FindByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authRepo_FindUserByID(t *testing.T) {
	type fields struct {
		db  *bun.DB
		rdb *redis.Client
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db:  tt.fields.db,
				rdb: tt.fields.rdb,
			}
			got, err := r.FindUserByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("authRepo.FindUserByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authRepo.FindUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authRepo_SaveRefreshToken(t *testing.T) {
	type fields struct {
		db  *bun.DB
		rdb *redis.Client
	}
	type args struct {
		ctx       context.Context
		userID    uint
		token     string
		expiresAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db:  tt.fields.db,
				rdb: tt.fields.rdb,
			}
			if err := r.SaveRefreshToken(tt.args.ctx, tt.args.userID, tt.args.token, tt.args.expiresAt); (err != nil) != tt.wantErr {
				t.Errorf("authRepo.SaveRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authRepo_FindRefreshToken(t *testing.T) {
	type fields struct {
		db  *bun.DB
		rdb *redis.Client
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db:  tt.fields.db,
				rdb: tt.fields.rdb,
			}
			got, err := r.FindRefreshToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Fatalf("authRepo.FindRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("authRepo.FindRefreshToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authRepo_DeleteRefreshToken(t *testing.T) {
	type fields struct {
		db  *bun.DB
		rdb *redis.Client
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db:  tt.fields.db,
				rdb: tt.fields.rdb,
			}
			if err := r.DeleteRefreshToken(tt.args.ctx, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("authRepo.DeleteRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authRepo_UpdatePassword(t *testing.T) {
	type fields struct {
		db  *bun.DB
		rdb *redis.Client
	}
	type args struct {
		ctx            context.Context
		userID         uint
		hashedPassword string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db:  tt.fields.db,
				rdb: tt.fields.rdb,
			}
			if err := r.UpdatePassword(tt.args.ctx, tt.args.userID, tt.args.hashedPassword); (err != nil) != tt.wantErr {
				t.Errorf("authRepo.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authRepo_SaveAuthToken(t *testing.T) {
	type fields struct {
		db  *bun.DB
		rdb *redis.Client
	}
	type args struct {
		ctx       context.Context
		userID    uint
		token     string
		tokenType string
		expiresAt time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db:  tt.fields.db,
				rdb: tt.fields.rdb,
			}
			if err := r.SaveAuthToken(tt.args.ctx, tt.args.userID, tt.args.token, tt.args.tokenType, tt.args.expiresAt); (err != nil) != tt.wantErr {
				t.Errorf("authRepo.SaveAuthToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authRepo_FindAuthToken(t *testing.T) {
	type fields struct {
		db  *bun.DB
		rdb *redis.Client
	}
	type args struct {
		ctx       context.Context
		token     string
		tokenType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    uint
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db:  tt.fields.db,
				rdb: tt.fields.rdb,
			}
			got, err := r.FindAuthToken(tt.args.ctx, tt.args.token, tt.args.tokenType)
			if (err != nil) != tt.wantErr {
				t.Fatalf("authRepo.FindAuthToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("authRepo.FindAuthToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authRepo_MarkTokenAsUsed(t *testing.T) {
	type fields struct {
		db  *bun.DB
		rdb *redis.Client
	}
	type args struct {
		ctx       context.Context
		token     string
		tokenType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db:  tt.fields.db,
				rdb: tt.fields.rdb,
			}
			if err := r.MarkTokenAsUsed(tt.args.ctx, tt.args.token, tt.args.tokenType); (err != nil) != tt.wantErr {
				t.Errorf("authRepo.MarkTokenAsUsed() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authRepo_DeleteAuthToken(t *testing.T) {
	type fields struct {
		db  *bun.DB
		rdb *redis.Client
	}
	type args struct {
		ctx       context.Context
		token     string
		tokenType string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db:  tt.fields.db,
				rdb: tt.fields.rdb,
			}
			if err := r.DeleteAuthToken(tt.args.ctx, tt.args.token, tt.args.tokenType); (err != nil) != tt.wantErr {
				t.Errorf("authRepo.DeleteAuthToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
