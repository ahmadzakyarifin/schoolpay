package repository

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/notification/domain"
	"github.com/uptrace/bun"
)

func TestNewNotificationRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want NotificationRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewNotificationRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewNotificationRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notificationRepo_Create(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
		n   *domain.Notification
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
			r := &notificationRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.db, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("notificationRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_notificationRepo_FindByID(t *testing.T) {
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
		want    *domain.Notification
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &notificationRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("notificationRepo.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("notificationRepo.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notificationRepo_UpdateStatusByWhatsappID(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx        context.Context
		whatsappID string
		status     string
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
			r := &notificationRepo{
				db: tt.fields.db,
			}
			if err := r.UpdateStatusByWhatsappID(tt.args.ctx, tt.args.whatsappID, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("notificationRepo.UpdateStatusByWhatsappID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_notificationRepo_GetStats(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &notificationRepo{
				db: tt.fields.db,
			}
			got, err := r.GetStats(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("notificationRepo.GetStats() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("notificationRepo.GetStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notificationRepo_GetEfficacyStats(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx     context.Context
		channel string
		start   *time.Time
		end     *time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &notificationRepo{
				db: tt.fields.db,
			}
			got, err := r.GetEfficacyStats(tt.args.ctx, tt.args.channel, tt.args.start, tt.args.end)
			if (err != nil) != tt.wantErr {
				t.Fatalf("notificationRepo.GetEfficacyStats() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("notificationRepo.GetEfficacyStats() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notificationRepo_GetByUserID(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx    context.Context
		userID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Notification
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &notificationRepo{
				db: tt.fields.db,
			}
			got, err := r.GetByUserID(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("notificationRepo.GetByUserID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("notificationRepo.GetByUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_notificationRepo_MarkAsRead(t *testing.T) {
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
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &notificationRepo{
				db: tt.fields.db,
			}
			if err := r.MarkAsRead(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("notificationRepo.MarkAsRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_notificationRepo_GetDetailedLogs(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx    context.Context
		page   int
		limit  int
		status string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []map[string]interface{}
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &notificationRepo{
				db: tt.fields.db,
			}
			got, got1, err := r.GetDetailedLogs(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.status, "", "")
			if (err != nil) != tt.wantErr {
				t.Fatalf("notificationRepo.GetDetailedLogs() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("notificationRepo.GetDetailedLogs() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("notificationRepo.GetDetailedLogs() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
