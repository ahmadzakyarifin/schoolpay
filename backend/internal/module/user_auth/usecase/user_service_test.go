package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/config"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	notificationdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/domain"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
)

func TestNewUserService(t *testing.T) {
	type args struct {
		db       *bun.DB
		repo     repository.UserRepo
		authRepo repository.AuthRepo
		msg      utils.Messenger
		noti     notificationrepo.NotificationRepo
		stuRepo  academicrepo.StudentRepo
		audit    auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want UserService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewUserService(tt.args.db, tt.args.repo, tt.args.authRepo, tt.args.msg, tt.args.noti, tt.args.stuRepo, tt.args.audit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewUserService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_GetAll(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	type args struct {
		ctx  context.Context
		role string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			got, err := s.GetAll(tt.args.ctx, tt.args.role)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userService.GetAll() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.GetAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_GetByID(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
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
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			got, err := s.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userService.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_GetPaginated(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	type args struct {
		ctx    context.Context
		page   int
		limit  int
		search string
		role   string
		filter string
		status string
		sort   string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.User
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			got, got1, err := s.GetPaginated(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.role, tt.args.filter, tt.args.status, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userService.GetPaginated() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.GetPaginated() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("userService.GetPaginated() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_userService_Create(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	type args struct {
		ctx  context.Context
		user *domain.User
		cfg  *config.Config
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
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			if err := s.Create(tt.args.ctx, tt.args.user, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("userService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_Update(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	type args struct {
		ctx  context.Context
		user *domain.User
		cfg  *config.Config
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
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			if err := s.Update(tt.args.ctx, tt.args.user, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("userService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_Delete(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
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
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("userService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_BulkDelete(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	type args struct {
		ctx context.Context
		ids []uint
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
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			if err := s.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("userService.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_Restore(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
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
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			if err := s.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("userService.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_BulkRestore(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	type args struct {
		ctx context.Context
		ids []uint
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
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			if err := s.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("userService.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_ActivateAccount(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	type args struct {
		ctx      context.Context
		token    string
		password string
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
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			got, err := s.ActivateAccount(tt.args.ctx, tt.args.token, tt.args.password)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userService.ActivateAccount() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.ActivateAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_ToggleStatus(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
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
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			if err := s.ToggleStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("userService.ToggleStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_ResendNotification(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	type args struct {
		ctx     context.Context
		id      uint
		channel string
		cfg     *config.Config
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
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			if err := s.ResendNotification(tt.args.ctx, tt.args.id, tt.args.channel, tt.args.cfg); (err != nil) != tt.wantErr {
				t.Errorf("userService.ResendNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_userService_BulkResendNotification(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	type args struct {
		ctx     context.Context
		ids     []uint
		channel string
		cfg     *config.Config
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *BulkResendResult
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			got, err := s.BulkResendNotification(tt.args.ctx, tt.args.ids, tt.args.channel, tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userService.BulkResendNotification() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.BulkResendNotification() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_ExportExcel(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	type args struct {
		ctx    context.Context
		search string
		role   string
		filter string
		status string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			got, err := s.ExportExcel(tt.args.ctx, tt.args.search, tt.args.role, tt.args.filter, tt.args.status)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userService.ExportExcel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.ExportExcel() = %v, want %v", got, tt.want)
			}
		})
	}
}


func Test_userService_GetNotifications(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	type args struct {
		ctx    context.Context
		userID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []notificationdomain.Notification
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			got, err := s.GetNotifications(tt.args.ctx, tt.args.userID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userService.GetNotifications() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.GetNotifications() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_userService_notificationWorker(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			s.notificationWorker()
		})
	}
}

func Test_userService_GetDependencyInfo(t *testing.T) {
	type fields struct {
		db          *bun.DB
		repo        repository.UserRepo
		authRepo    repository.AuthRepo
		messenger   utils.Messenger
		notiRepo    notificationrepo.NotificationRepo
		studentRepo academicrepo.StudentRepo
		audit       auditusecase.AuditLogService
		jobChan     chan userNotifyJob
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &userService{
				db:          tt.fields.db,
				repo:        tt.fields.repo,
				authRepo:    tt.fields.authRepo,
				messenger:   tt.fields.messenger,
				notiRepo:    tt.fields.notiRepo,
				studentRepo: tt.fields.studentRepo,
				audit:       tt.fields.audit,
				jobChan:     tt.fields.jobChan,
			}
			got, err := s.GetDependencyInfo(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("userService.GetDependencyInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("userService.GetDependencyInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}
