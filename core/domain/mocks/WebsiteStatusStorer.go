// Code generated by mockery v2.16.0. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/sainak/status-checker/core/domain"
	mock "github.com/stretchr/testify/mock"
)

// WebsiteStatusStorer is an autogenerated mock type for the WebsiteStatusStorer type
type WebsiteStatusStorer struct {
	mock.Mock
}

// DropWebsite provides a mock function with given fields: ctx, id
func (_m *WebsiteStatusStorer) DropWebsite(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertStatus provides a mock function with given fields: ctx, status
func (_m *WebsiteStatusStorer) InsertStatus(ctx context.Context, status *domain.Status) error {
	ret := _m.Called(ctx, status)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Status) error); ok {
		r0 = rf(ctx, status)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// InsertWebsite provides a mock function with given fields: ctx, website
func (_m *WebsiteStatusStorer) InsertWebsite(ctx context.Context, website *domain.Website) error {
	ret := _m.Called(ctx, website)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Website) error); ok {
		r0 = rf(ctx, website)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// QueryStatusesByWebsiteID provides a mock function with given fields: ctx, websiteID, cursor, num
func (_m *WebsiteStatusStorer) QueryStatusesByWebsiteID(ctx context.Context, websiteID int64, cursor string, num int64) ([]domain.Status, string, error) {
	ret := _m.Called(ctx, websiteID, cursor, num)

	var r0 []domain.Status
	if rf, ok := ret.Get(0).(func(context.Context, int64, string, int64) []domain.Status); ok {
		r0 = rf(ctx, websiteID, cursor, num)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Status)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, int64, string, int64) string); ok {
		r1 = rf(ctx, websiteID, cursor, num)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, int64, string, int64) error); ok {
		r2 = rf(ctx, websiteID, cursor, num)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// QueryWebsiteStatusByID provides a mock function with given fields: ctx, id
func (_m *WebsiteStatusStorer) QueryWebsiteStatusByID(ctx context.Context, id int64) (domain.WebsiteStatus, error) {
	ret := _m.Called(ctx, id)

	var r0 domain.WebsiteStatus
	if rf, ok := ret.Get(0).(func(context.Context, int64) domain.WebsiteStatus); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(domain.WebsiteStatus)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// QueryWebsites provides a mock function with given fields: ctx, cursor, num, filters
func (_m *WebsiteStatusStorer) QueryWebsites(ctx context.Context, cursor string, num int64, filters map[string]string) ([]domain.Website, string, error) {
	ret := _m.Called(ctx, cursor, num, filters)

	var r0 []domain.Website
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, map[string]string) []domain.Website); ok {
		r0 = rf(ctx, cursor, num, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.Website)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string, int64, map[string]string) string); ok {
		r1 = rf(ctx, cursor, num, filters)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, int64, map[string]string) error); ok {
		r2 = rf(ctx, cursor, num, filters)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// QueryWebsitesStatus provides a mock function with given fields: ctx, cursor, num, filters
func (_m *WebsiteStatusStorer) QueryWebsitesStatus(ctx context.Context, cursor string, num int64, filters map[string]string) ([]domain.WebsiteStatus, string, error) {
	ret := _m.Called(ctx, cursor, num, filters)

	var r0 []domain.WebsiteStatus
	if rf, ok := ret.Get(0).(func(context.Context, string, int64, map[string]string) []domain.WebsiteStatus); ok {
		r0 = rf(ctx, cursor, num, filters)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.WebsiteStatus)
		}
	}

	var r1 string
	if rf, ok := ret.Get(1).(func(context.Context, string, int64, map[string]string) string); ok {
		r1 = rf(ctx, cursor, num, filters)
	} else {
		r1 = ret.Get(1).(string)
	}

	var r2 error
	if rf, ok := ret.Get(2).(func(context.Context, string, int64, map[string]string) error); ok {
		r2 = rf(ctx, cursor, num, filters)
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

type mockConstructorTestingTNewWebsiteStatusStorer interface {
	mock.TestingT
	Cleanup(func())
}

// NewWebsiteStatusStorer creates a new instance of WebsiteStatusStorer. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewWebsiteStatusStorer(t mockConstructorTestingTNewWebsiteStatusStorer) *WebsiteStatusStorer {
	mock := &WebsiteStatusStorer{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}