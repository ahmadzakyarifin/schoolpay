package mocks

import (
	"context"
	"time"

	mock "github.com/stretchr/testify/mock"
)

type DashboardService struct {
	mock.Mock
}

func (_m *DashboardService) GetStats(ctx context.Context, period string, pStart, pEnd, pPrevStart, pPrevEnd *time.Time, filterEntryYear int, classID, majorID uint, page, limit int, billTypeID uint) (map[string]interface{}, error) {
	ret := _m.Called(ctx, period, pStart, pEnd, pPrevStart, pPrevEnd, filterEntryYear, classID, majorID, page, limit, billTypeID)

	var r0 map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, *time.Time, *time.Time, *time.Time, *time.Time, int, uint, uint, int, int, uint) (map[string]interface{}, error)); ok {
		return rf(ctx, period, pStart, pEnd, pPrevStart, pPrevEnd, filterEntryYear, classID, majorID, page, limit, billTypeID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, *time.Time, *time.Time, *time.Time, *time.Time, int, uint, uint, int, int, uint) map[string]interface{}); ok {
		r0 = rf(ctx, period, pStart, pEnd, pPrevStart, pPrevEnd, filterEntryYear, classID, majorID, page, limit, billTypeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, *time.Time, *time.Time, *time.Time, *time.Time, int, uint, uint, int, int, uint) error); ok {
		r1 = rf(ctx, period, pStart, pEnd, pPrevStart, pPrevEnd, filterEntryYear, classID, majorID, page, limit, billTypeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *DashboardService) GetCommunicationDetails(ctx context.Context, status, channel string, start, end *time.Time, filterEntryYear int, classID, majorID uint) ([]map[string]interface{}, error) {
	ret := _m.Called(ctx, status, channel, start, end, filterEntryYear, classID, majorID)

	var r0 []map[string]interface{}
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *time.Time, *time.Time, int, uint, uint) ([]map[string]interface{}, error)); ok {
		return rf(ctx, status, channel, start, end, filterEntryYear, classID, majorID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, *time.Time, *time.Time, int, uint, uint) []map[string]interface{}); ok {
		r0 = rf(ctx, status, channel, start, end, filterEntryYear, classID, majorID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]map[string]interface{})
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, *time.Time, *time.Time, int, uint, uint) error); ok {
		r1 = rf(ctx, status, channel, start, end, filterEntryYear, classID, majorID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *DashboardService) ExportGlobalReport(ctx context.Context, format, tab, search string, pStart, pEnd *time.Time, filterEntryYear int, classID, majorID, billTypeID uint) ([]byte, error) {
	ret := _m.Called(ctx, format, tab, search, pStart, pEnd, filterEntryYear, classID, majorID, billTypeID)

	var r0 []byte
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, *time.Time, *time.Time, int, uint, uint, uint) ([]byte, error)); ok {
		return rf(ctx, format, tab, search, pStart, pEnd, filterEntryYear, classID, majorID, billTypeID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, string, string, string, *time.Time, *time.Time, int, uint, uint, uint) []byte); ok {
		r0 = rf(ctx, format, tab, search, pStart, pEnd, filterEntryYear, classID, majorID, billTypeID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]byte)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, string, string, string, *time.Time, *time.Time, int, uint, uint, uint) error); ok {
		r1 = rf(ctx, format, tab, search, pStart, pEnd, filterEntryYear, classID, majorID, billTypeID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
