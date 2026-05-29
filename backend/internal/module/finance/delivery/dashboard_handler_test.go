package delivery

import (
	"testing"

	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	auditusecase "github.com/ahmadzakyarifin/schoolpay/internal/module/audit/usecase"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
	notificationrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/notification/repository"
	userauthrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/user_auth/repository"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
)

func TestNewDashboardHandler(t *testing.T) {
	type args struct {
		db        *bun.DB
		userRepo  userauthrepo.UserRepo
		stuRepo   academicrepo.StudentRepo
		finRepo   financerepo.FinanceReportRepo
		ayRepo    academicrepo.AcademicYearRepo
		notifRepo notificationrepo.NotificationRepo
		audit     auditusecase.AuditLogService
	}
	tests := []struct {
		name string
		args args
		want *DashboardHandler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDashboardHandler(tt.args.db, tt.args.userRepo, tt.args.stuRepo, tt.args.finRepo, tt.args.ayRepo, tt.args.notifRepo, tt.args.audit)
			if (got == nil && tt.want != nil) || (got != nil && tt.want == nil) {
				t.Errorf("NewDashboardHandler() = %v, want %v", got, tt.want)
				return
			}
			if got != nil && tt.want != nil {
				if got.db != tt.want.db || got.userRepo != tt.want.userRepo || got.stuRepo != tt.want.stuRepo || got.finRepo != tt.want.finRepo || got.ayRepo != tt.want.ayRepo || got.notifRepo != tt.want.notifRepo || got.audit != tt.want.audit {
					t.Errorf("NewDashboardHandler() fields mismatch")
				}
			}
		})
	}
}

func TestDashboardHandler_GetStats(t *testing.T) {
	type fields struct {
		db        *bun.DB
		userRepo  userauthrepo.UserRepo
		stuRepo   academicrepo.StudentRepo
		finRepo   financerepo.FinanceReportRepo
		ayRepo    academicrepo.AcademicYearRepo
		notifRepo notificationrepo.NotificationRepo
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DashboardHandler{
				db:        tt.fields.db,
				userRepo:  tt.fields.userRepo,
				stuRepo:   tt.fields.stuRepo,
				finRepo:   tt.fields.finRepo,
				ayRepo:    tt.fields.ayRepo,
				notifRepo: tt.fields.notifRepo,
			}
			h.GetStats(tt.args.c)
		})
	}
}

func TestDashboardHandler_GetCommunicationDetails(t *testing.T) {
	type fields struct {
		db        *bun.DB
		userRepo  userauthrepo.UserRepo
		stuRepo   academicrepo.StudentRepo
		finRepo   financerepo.FinanceReportRepo
		ayRepo    academicrepo.AcademicYearRepo
		notifRepo notificationrepo.NotificationRepo
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DashboardHandler{
				db:        tt.fields.db,
				userRepo:  tt.fields.userRepo,
				stuRepo:   tt.fields.stuRepo,
				finRepo:   tt.fields.finRepo,
				ayRepo:    tt.fields.ayRepo,
				notifRepo: tt.fields.notifRepo,
			}
			h.GetCommunicationDetails(tt.args.c)
		})
	}
}

func Test_calculateGrowth(t *testing.T) {
	type args struct {
		current  float64
		previous float64
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := calculateGrowth(tt.args.current, tt.args.previous); got != tt.want {
				t.Errorf("calculateGrowth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDashboardHandler_ExportGlobalReport(t *testing.T) {
	type fields struct {
		db        *bun.DB
		userRepo  userauthrepo.UserRepo
		stuRepo   academicrepo.StudentRepo
		finRepo   financerepo.FinanceReportRepo
		ayRepo    academicrepo.AcademicYearRepo
		notifRepo notificationrepo.NotificationRepo
	}
	type args struct {
		c *gin.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DashboardHandler{
				db:        tt.fields.db,
				userRepo:  tt.fields.userRepo,
				stuRepo:   tt.fields.stuRepo,
				finRepo:   tt.fields.finRepo,
				ayRepo:    tt.fields.ayRepo,
				notifRepo: tt.fields.notifRepo,
			}
			h.ExportGlobalReport(tt.args.c)
		})
	}
}
