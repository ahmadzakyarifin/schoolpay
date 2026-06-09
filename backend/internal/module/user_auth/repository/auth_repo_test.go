package repository

import (
	"context"
	"testing"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/uptrace/bun"
)

func TestNewAuthRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want AuthRepo
	}{
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAuthRepo(tt.args.db); got != tt.want {
				t.Errorf("NewAuthRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authRepo_GetDB(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	tests := []struct {
		name   string
		fields fields
		want   *bun.DB
	}{
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db: tt.fields.db,
			}
			if got := r.GetDB(); got != tt.want {
				t.Errorf("authRepo.GetDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authRepo_FindByEmail(t *testing.T) {
	type fields struct {
		db *bun.DB
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
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByEmail(tt.args.ctx, tt.args.email)
			if (err != nil) != tt.wantErr {
				t.Fatalf("authRepo.FindByEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("authRepo.FindByEmail() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authRepo_FindUserByID(t *testing.T) {
	type fields struct {
		db *bun.DB
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
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db: tt.fields.db,
			}
			got, err := r.FindUserByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("authRepo.FindUserByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("authRepo.FindUserByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authRepo_SaveRefreshToken(t *testing.T) {
	type fields struct {
		db *bun.DB
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
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db: tt.fields.db,
			}
			if err := r.SaveRefreshToken(tt.args.ctx, tt.args.userID, tt.args.token, tt.args.expiresAt); (err != nil) != tt.wantErr {
				t.Errorf("authRepo.SaveRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authRepo_FindUserByRefreshToken(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx   context.Context
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.User
		wantErr bool
	}{
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db: tt.fields.db,
			}
			got, _, err := r.FindUserByRefreshToken(tt.args.ctx, tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Fatalf("authRepo.FindUserByRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if got != tt.want {
				t.Errorf("authRepo.FindUserByRefreshToken() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_authRepo_DeleteRefreshToken(t *testing.T) {
	type fields struct {
		db *bun.DB
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
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db: tt.fields.db,
			}
			if err := r.DeleteRefreshToken(tt.args.ctx, tt.args.token); (err != nil) != tt.wantErr {
				t.Errorf("authRepo.DeleteRefreshToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authRepo_UpdatePassword(t *testing.T) {
	type fields struct {
		db *bun.DB
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
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db: tt.fields.db,
			}
			if err := r.UpdatePassword(tt.args.ctx, tt.args.userID, tt.args.hashedPassword); (err != nil) != tt.wantErr {
				t.Errorf("authRepo.UpdatePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authRepo_SaveAuthToken(t *testing.T) {
	type fields struct {
		db *bun.DB
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
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db: tt.fields.db,
			}
			if err := r.SaveAuthToken(tt.args.ctx, tt.args.userID, tt.args.token, tt.args.tokenType, tt.args.expiresAt); (err != nil) != tt.wantErr {
				t.Errorf("authRepo.SaveAuthToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_authRepo_FindAuthToken(t *testing.T) {
	type fields struct {
		db *bun.DB
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
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db: tt.fields.db,
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

func Test_authRepo_DeleteAuthToken(t *testing.T) {
	type fields struct {
		db *bun.DB
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
		// Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &authRepo{
				db: tt.fields.db,
			}
			if err := r.DeleteAuthToken(tt.args.ctx, tt.args.token, tt.args.tokenType); (err != nil) != tt.wantErr {
				t.Errorf("authRepo.DeleteAuthToken() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
