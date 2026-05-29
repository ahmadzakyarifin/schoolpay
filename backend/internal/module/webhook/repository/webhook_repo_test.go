package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/webhook/domain"
	"github.com/uptrace/bun"
)

func TestNewWebhookRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want WebhookRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewWebhookRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewWebhookRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_webhookRepo_Create(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		log *domain.WebhookLog
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
			r := &webhookRepo{
				db: tt.fields.db,
			}
			if err := r.Create(tt.args.ctx, tt.args.log); (err != nil) != tt.wantErr {
				t.Errorf("webhookRepo.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_webhookRepo_UpdateStatus(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx     context.Context
		eventID string
		status  string
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
			r := &webhookRepo{
				db: tt.fields.db,
			}
			if err := r.UpdateStatus(tt.args.ctx, tt.args.eventID, tt.args.status); (err != nil) != tt.wantErr {
				t.Errorf("webhookRepo.UpdateStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
