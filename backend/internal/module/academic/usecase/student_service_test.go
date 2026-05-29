package usecase

import (
	"context"
	"reflect"
	"testing"

	"github.com/ahmadzakyarifin/schoolpay/config"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	"github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	userauthdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/domain"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/ahmadzakyarifin/schoolpay/pkg/utils"
	"github.com/uptrace/bun"
)

func TestNewStudentService(t *testing.T) {
	type args struct {
		db       *bun.DB
		repo     repository.StudentRepo
		userRepo userauthrepo.UserRepo
		authRepo userauthrepo.AuthRepo
		msg      utils.Messenger
		noti     notificationrepo.NotificationRepo
		ay       repository.AcademicYearRepo
		jur      repository.MajorRepo
		cls      repository.ClassRepo
		billRepo financerepo.StudentBillRepo
		cfg      *config.Config
		audit    auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want StudentService
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewStudentService(tt.args.db, tt.args.repo, tt.args.userRepo, tt.args.authRepo, tt.args.msg, tt.args.noti, tt.args.ay, tt.args.jur, tt.args.cls, tt.args.billRepo, tt.args.cfg, tt.args.audit); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewStudentService() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentService_GetPaginated(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
	}
	type args struct {
		ctx       context.Context
		page      int
		limit     int
		search    string
		filter    string
		status    string
		entryYear int
		classID   uint
		majorID   uint
		sort      string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Student
		want1   int
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			got, got1, err := s.GetPaginated(tt.args.ctx, tt.args.page, tt.args.limit, tt.args.search, tt.args.filter, tt.args.status, tt.args.entryYear, tt.args.classID, tt.args.majorID, tt.args.sort)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentService.GetPaginated() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentService.GetPaginated() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("studentService.GetPaginated() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_studentService_GetAcademicFilters(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
	}
	type args struct {
		ctx context.Context
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
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			got, err := s.GetAcademicFilters(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentService.GetAcademicFilters() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentService.GetAcademicFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentService_GetParents(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
	}
	type args struct {
		ctx       context.Context
		studentID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []userauthdomain.User
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			got, err := s.GetParents(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentService.GetParents() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentService.GetParents() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentService_GetStudentsByParentID(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
	}
	type args struct {
		ctx      context.Context
		parentID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.Student
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			got, err := s.GetStudentsByParentID(tt.args.ctx, tt.args.parentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentService.GetStudentsByParentID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentService.GetStudentsByParentID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentService_Create(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
	}
	type args struct {
		ctx     context.Context
		student *domain.Student
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
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			if err := s.Create(tt.args.ctx, tt.args.student); (err != nil) != tt.wantErr {
				t.Errorf("studentService.Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentService_Update(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
	}
	type args struct {
		ctx     context.Context
		student *domain.Student
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
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			if err := s.Update(tt.args.ctx, tt.args.student); (err != nil) != tt.wantErr {
				t.Errorf("studentService.Update() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentService_Delete(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
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
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			if err := s.Delete(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("studentService.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentService_GetByID(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
	}
	type args struct {
		ctx context.Context
		id  uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *domain.Student
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			got, err := s.GetByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentService.GetByID() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentService.GetByID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentService_ToggleStatus(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
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
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			if err := s.ToggleStatus(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("studentService.ToggleStatus() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentService_ExportExcel(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
	}
	type args struct {
		ctx       context.Context
		search    string
		filter    string
		status    string
		entryYear int
		classID   uint
		majorID   uint
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
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			got, err := s.ExportExcel(tt.args.ctx, tt.args.search, tt.args.filter, tt.args.status, tt.args.entryYear, tt.args.classID, tt.args.majorID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentService.ExportExcel() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentService.ExportExcel() = %v, want %v", got, tt.want)
			}
		})
	}
}


func Test_studentService_BulkGraduate(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
	}
	type args struct {
		ctx        context.Context
		classID    uint
		studentIDs []uint
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
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			if err := s.BulkGraduate(tt.args.ctx, tt.args.classID, tt.args.studentIDs); (err != nil) != tt.wantErr {
				t.Errorf("studentService.BulkGraduate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentService_BulkPromote(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
	}
	type args struct {
		ctx           context.Context
		sourceClassID uint
		targetClassID uint
		studentIDs    []uint
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
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			if err := s.BulkPromote(tt.args.ctx, tt.args.sourceClassID, tt.args.targetClassID, tt.args.studentIDs); (err != nil) != tt.wantErr {
				t.Errorf("studentService.BulkPromote() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentService_GetClassHistory(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
	}
	type args struct {
		ctx       context.Context
		studentID uint
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []domain.ClassHistory
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			got, err := s.GetClassHistory(tt.args.ctx, tt.args.studentID)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentService.GetClassHistory() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentService.GetClassHistory() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentService_BulkDelete(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
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
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			if err := s.BulkDelete(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("studentService.BulkDelete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentService_GetDependencyInfo(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
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
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			got, err := s.GetDependencyInfo(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("studentService.GetDependencyInfo() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("studentService.GetDependencyInfo() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_studentService_Restore(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
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
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			if err := s.Restore(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("studentService.Restore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentService_BulkRestore(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
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
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			if err := s.BulkRestore(tt.args.ctx, tt.args.ids); (err != nil) != tt.wantErr {
				t.Errorf("studentService.BulkRestore() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_studentService_notificationWorker(t *testing.T) {
	type fields struct {
		db        *bun.DB
		repo      repository.StudentRepo
		userRepo  userauthrepo.UserRepo
		authRepo  userauthrepo.AuthRepo
		msg       utils.Messenger
		notiRepo  notificationrepo.NotificationRepo
		ayRepo    repository.AcademicYearRepo
		majorRepo repository.MajorRepo
		classRepo repository.ClassRepo
		cfg       *config.Config
		billRepo  financerepo.StudentBillRepo
		audit     auditusecase.AuditLogService
		jobChan   chan studentNotifyJob
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &studentService{
				db:        tt.fields.db,
				repo:      tt.fields.repo,
				userRepo:  tt.fields.userRepo,
				authRepo:  tt.fields.authRepo,
				msg:       tt.fields.msg,
				notiRepo:  tt.fields.notiRepo,
				ayRepo:    tt.fields.ayRepo,
				majorRepo: tt.fields.majorRepo,
				classRepo: tt.fields.classRepo,
				cfg:       tt.fields.cfg,
				billRepo:  tt.fields.billRepo,
				audit:     tt.fields.audit,
				jobChan:   tt.fields.jobChan,
			}
			s.notificationWorker()
		})
	}
}
