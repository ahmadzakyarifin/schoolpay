package usecase

import (
	"context"
	"testing"
	"time"

	academicdomain "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/domain"
	academicrepo "github.com/ahmadzakyarifin/schoolpay/internal/module/academic/repository"
	domain "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	financerepo "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/repository"
)

// mockDashboardRepo is a manual mock for DashboardRepo
type mockDashboardRepo struct {
	financerepo.DashboardRepo
}

func (m *mockDashboardRepo) CountNewUsersByPeriod(ctx context.Context, start, end *time.Time) (int, error) {
	return 5, nil
}
func (m *mockDashboardRepo) CountTotalUsersUpTo(ctx context.Context, end *time.Time) (int, error) {
	return 100, nil
}
func (m *mockDashboardRepo) CountUsersByRoleUpTo(ctx context.Context, end *time.Time, role string) (int, error) {
	return 50, nil
}
func (m *mockDashboardRepo) CountStudentsByStatusUpTo(ctx context.Context, end *time.Time, status string, academicYear int, classID, majorID uint) (int, error) {
	return 200, nil
}
func (m *mockDashboardRepo) GetGenderStats(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]int, error) {
	return map[string]int{"L": 110, "P": 90}, nil
}
func (m *mockDashboardRepo) GetEfficacyStats(ctx context.Context, channel string, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]int, error) {
	return map[string]int{"sent": 10, "failed": 0}, nil
}
func (m *mockDashboardRepo) GetRecentNotifications(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"title": "Test notification"}}, nil
}
func (m *mockDashboardRepo) GetCommunicationDetails(ctx context.Context, status, channel string, start, end *time.Time, filterEntryYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"status": status}}, nil
}
func (m *mockDashboardRepo) GetCriticalBills(ctx context.Context, status string, start, end *time.Time, academicYear int, classID, majorID, billTypeID uint) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"id": uint(1)}}, nil
}
func (m *mockDashboardRepo) GetTotalBillsByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (float64, error) {
	return 15000000, nil
}
func (m *mockDashboardRepo) GetTotalPaymentsByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (float64, error) {
	return 10000000, nil
}
func (m *mockDashboardRepo) GetPaymentTrendByPeriod(ctx context.Context, start, end *time.Time, interval string, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"date": "2026-06-01", "total": float64(10000000)}}, nil
}
func (m *mockDashboardRepo) GetDashboardSummary(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]interface{}, error) {
	return map[string]interface{}{
		"total_paid_amount":   float64(10000000),
		"total_unpaid_amount": float64(5000000),
		"paid_count":          10,
		"unpaid_count":        5,
	}, nil
}
func (m *mockDashboardRepo) GetRecentPayments(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID, billTypeID uint, search string, page, limit int) ([]map[string]interface{}, int, error) {
	return []map[string]interface{}{{"amount": 1000000}}, 1, nil
}
func (m *mockDashboardRepo) GetPaymentMethodsCount(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	return []map[string]interface{}{{"method": "cash", "count": 10}}, nil
}
func (m *mockDashboardRepo) GetArrearsPaged(ctx context.Context, page, limit int, academicYear int, classID, majorID, billTypeID uint, search string, start, end *time.Time) ([]domain.ArrearRecord, int, error) {
	return []domain.ArrearRecord{{StudentName: "Student 1"}}, 1, nil
}

// mockStudentRepo stub
type mockStudentRepo struct {
	academicrepo.StudentRepo
}

func (m *mockStudentRepo) CountByMajor(ctx context.Context, academicYear int, classID uint, majorID uint) (map[string]int, error) {
	return map[string]int{"IPA": 100}, nil
}
func (m *mockStudentRepo) CountByClass(ctx context.Context, academicYear int, classID uint, majorID uint) (map[string]int, error) {
	return map[string]int{"X-1": 30}, nil
}

// mockAcademicYearRepo stub
type mockAcademicYearRepo struct {
	academicrepo.AcademicYearRepo
}

func (m *mockAcademicYearRepo) FindAll(ctx context.Context, page, limit int, search, order, sort string) ([]academicdomain.AcademicYear, int, error) {
	return []academicdomain.AcademicYear{{Year: 2026}}, 1, nil
}

func TestDashboardService_GetStats(t *testing.T) {
	dashRepo := &mockDashboardRepo{}
	stuRepo := &mockStudentRepo{}
	ayRepo := &mockAcademicYearRepo{}

	service := NewDashboardService(dashRepo, stuRepo, ayRepo, nil)
	ctx := context.Background()

	now := time.Now()
	prev := now.AddDate(0, -1, 0)
	stats, err := service.GetStats(ctx, "monthly", &now, &now, &prev, &prev, 0, 0, 0, 1, 10, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if stats["total_payments_count"] != 1 {
		t.Errorf("expected total payments count to be 1, got %v", stats["total_payments_count"])
	}

	statsMap, ok := stats["stats"].(map[string]interface{})
	if !ok {
		t.Fatalf("stats field not found or not map")
	}

	usersMap := statsMap["users"].(map[string]interface{})
	if usersMap["total"] != 100 {
		t.Errorf("expected total users to be 100, got %v", usersMap["total"])
	}
}

func TestDashboardService_GetCommunicationDetails(t *testing.T) {
	dashRepo := &mockDashboardRepo{}
	stuRepo := &mockStudentRepo{}
	ayRepo := &mockAcademicYearRepo{}

	service := NewDashboardService(dashRepo, stuRepo, ayRepo, nil)
	ctx := context.Background()

	details, err := service.GetCommunicationDetails(ctx, "sent", "whatsapp", nil, nil, 0, 0, 0)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(details) != 1 || details[0]["status"] != "sent" {
		t.Errorf("expected communication status 'sent', got %v", details)
	}
}
