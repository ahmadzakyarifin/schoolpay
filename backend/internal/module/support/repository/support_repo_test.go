package repository

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/internal/module/support/domain"
	"github.com/uptrace/bun"
)

func TestNewSupportRepo(t *testing.T) {
	type args struct {
		db *bun.DB
	}
	tests := []struct {
		name string
		args args
		want SupportRepo
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSupportRepo(tt.args.db); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewSupportRepo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_supportRepo_FindOpenByPhone(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx   context.Context
		phone string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Conversation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &supportRepo{
				db: tt.fields.db,
			}
			got, err := r.FindOpenByPhone(tt.args.ctx, tt.args.phone)
			if (err != nil) != tt.wantErr {
				t.Fatalf("supportRepo.FindOpenByPhone() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("supportRepo.FindOpenByPhone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_supportRepo_CreateConversation(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
		c   *domain.Conversation
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
			r := &supportRepo{
				db: tt.fields.db,
			}
			if err := r.CreateConversation(tt.args.ctx, tt.args.db, tt.args.c); (err != nil) != tt.wantErr {
				t.Errorf("supportRepo.CreateConversation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_supportRepo_UpdateConversationPreview(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx         context.Context
		db          bun.IDB
		id          uint
		lastMessage string
		unreadDelta int
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
			r := &supportRepo{
				db: tt.fields.db,
			}
			if err := r.UpdateConversationPreview(tt.args.ctx, tt.args.db, tt.args.id, tt.args.lastMessage, tt.args.unreadDelta); (err != nil) != tt.wantErr {
				t.Errorf("supportRepo.UpdateConversationPreview() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_supportRepo_CreateMessage(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx context.Context
		db  bun.IDB
		m   *domain.Message
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
			r := &supportRepo{
				db: tt.fields.db,
			}
			if err := r.CreateMessage(tt.args.ctx, tt.args.db, tt.args.m); (err != nil) != tt.wantErr {
				t.Errorf("supportRepo.CreateMessage() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_supportRepo_FindAll(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx    context.Context
		status string
		page   int
		limit  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Conversation
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &supportRepo{
				db: tt.fields.db,
			}
			got, got1, err := r.FindAll(tt.args.ctx, tt.args.status, tt.args.page, tt.args.limit)
			if (err != nil) != tt.wantErr {
				t.Fatalf("supportRepo.FindAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("supportRepo.FindAll() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("supportRepo.FindAll() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_supportRepo_FindMessages(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx            context.Context
		conversationID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Message
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &supportRepo{
				db: tt.fields.db,
			}
			got, err := r.FindMessages(tt.args.ctx, tt.args.conversationID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("supportRepo.FindMessages() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("supportRepo.FindMessages() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_supportRepo_FindByID(t *testing.T) {
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
		want    *domain.Conversation
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &supportRepo{
				db: tt.fields.db,
			}
			got, err := r.FindByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("supportRepo.FindByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("supportRepo.FindByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_supportRepo_Assign(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx            context.Context
		db             bun.IDB
		conversationID uint
		adminID        uint
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
			r := &supportRepo{
				db: tt.fields.db,
			}
			if err := r.Assign(tt.args.ctx, tt.args.db, tt.args.conversationID, tt.args.adminID); (err != nil) != tt.wantErr {
				t.Errorf("supportRepo.Assign() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_supportRepo_MarkRead(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx            context.Context
		db             bun.IDB
		conversationID uint
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
			r := &supportRepo{
				db: tt.fields.db,
			}
			if err := r.MarkRead(tt.args.ctx, tt.args.db, tt.args.conversationID); (err != nil) != tt.wantErr {
				t.Errorf("supportRepo.MarkRead() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_supportRepo_Close(t *testing.T) {
	type fields struct {
		db *bun.DB
	}
	type args struct {
		ctx            context.Context
		db             bun.IDB
		conversationID uint
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
			r := &supportRepo{
				db: tt.fields.db,
			}
			if err := r.Close(tt.args.ctx, tt.args.db, tt.args.conversationID); (err != nil) != tt.wantErr {
				t.Errorf("supportRepo.Close() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
