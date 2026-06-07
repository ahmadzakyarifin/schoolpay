package mocks

import (
	"context"
	"time"

	domain "github.com/ahmadzakyarifin/schoolpay/internal/module/finance/domain"
	mock "github.com/stretchr/testify/mock"
)

type DashboardRepo struct {
	mock.Mock
}

func (_m *DashboardRepo) CountNewUsersByPeriod(ctx context.Context, start, end *time.Time) (int, error) {
	ret := _m.Called(ctx, start, end)
	return ret.Get(0).(int), ret.Error(1)
}

func (_m *DashboardRepo) CountTotalUsersUpTo(ctx context.Context, end *time.Time) (int, error) {
	ret := _m.Called(ctx, end)
	return ret.Get(0).(int), ret.Error(1)
}

func (_m *DashboardRepo) CountUsersByRoleUpTo(ctx context.Context, end *time.Time, role string) (int, error) {
	ret := _m.Called(ctx, end, role)
	return ret.Get(0).(int), ret.Error(1)
}

func (_m *DashboardRepo) CountStudentsByStatusUpTo(ctx context.Context, end *time.Time, status string, academicYear int, classID, majorID uint) (int, error) {
	ret := _m.Called(ctx, end, status, academicYear, classID, majorID)
	return ret.Get(0).(int), ret.Error(1)
}

func (_m *DashboardRepo) GetGenderStats(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]int, error) {
	ret := _m.Called(ctx, start, end, academicYear, classID, majorID)
	return ret.Get(0).(map[string]int), ret.Error(1)
}

func (_m *DashboardRepo) GetEfficacyStats(ctx context.Context, channel string, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]int, error) {
	ret := _m.Called(ctx, channel, start, end, academicYear, classID, majorID)
	return ret.Get(0).(map[string]int), ret.Error(1)
}

func (_m *DashboardRepo) GetRecentNotifications(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	ret := _m.Called(ctx, start, end, academicYear, classID, majorID)
	return ret.Get(0).([]map[string]interface{}), ret.Error(1)
}

func (_m *DashboardRepo) GetCommunicationDetails(ctx context.Context, status, channel string, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	ret := _m.Called(ctx, status, channel, start, end, academicYear, classID, majorID)
	return ret.Get(0).([]map[string]interface{}), ret.Error(1)
}

func (_m *DashboardRepo) GetCriticalBills(ctx context.Context, status string, start, end *time.Time, academicYear int, classID, majorID, billTypeID uint) ([]map[string]interface{}, error) {
	ret := _m.Called(ctx, status, start, end, academicYear, classID, majorID, billTypeID)
	return ret.Get(0).([]map[string]interface{}), ret.Error(1)
}

func (_m *DashboardRepo) GetTotalBillsByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (float64, error) {
	ret := _m.Called(ctx, start, end, academicYear, classID, majorID)
	return ret.Get(0).(float64), ret.Error(1)
}

func (_m *DashboardRepo) GetTotalPaymentsByPeriod(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (float64, error) {
	ret := _m.Called(ctx, start, end, academicYear, classID, majorID)
	return ret.Get(0).(float64), ret.Error(1)
}

func (_m *DashboardRepo) GetPaymentTrendByPeriod(ctx context.Context, start, end *time.Time, interval string, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	ret := _m.Called(ctx, start, end, interval, academicYear, classID, majorID)
	return ret.Get(0).([]map[string]interface{}), ret.Error(1)
}

func (_m *DashboardRepo) GetDashboardSummary(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) (map[string]interface{}, error) {
	ret := _m.Called(ctx, start, end, academicYear, classID, majorID)
	return ret.Get(0).(map[string]interface{}), ret.Error(1)
}

func (_m *DashboardRepo) GetRecentPayments(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID, billTypeID uint, search string, page, limit int) ([]map[string]interface{}, int, error) {
	ret := _m.Called(ctx, start, end, academicYear, classID, majorID, billTypeID, search, page, limit)
	return ret.Get(0).([]map[string]interface{}), ret.Get(1).(int), ret.Error(2)
}

func (_m *DashboardRepo) GetPaymentMethodsCount(ctx context.Context, start, end *time.Time, academicYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	ret := _m.Called(ctx, start, end, academicYear, classID, majorID)
	return ret.Get(0).([]map[string]interface{}), ret.Error(1)
}

func (_m *DashboardRepo) GetArrearsPaged(ctx context.Context, page, limit int, academicYear int, classID, majorID, billTypeID uint, search string, start, end *time.Time) ([]domain.ArrearRecord, int, error) {
	ret := _m.Called(ctx, page, limit, academicYear, classID, majorID, billTypeID, search, start, end)
	return ret.Get(0).([]domain.ArrearRecord), ret.Get(1).(int), ret.Error(2)
}
