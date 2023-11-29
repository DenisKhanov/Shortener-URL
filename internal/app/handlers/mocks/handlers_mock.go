// Code generated by MockGen. DO NOT EDIT.
// Source: handlers.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	models "github.com/DenisKhanov/shorterURL/internal/app/models"
	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// AsyncDeleteUserURLs mocks base method.
func (m *MockService) AsyncDeleteUserURLs(ctx context.Context, URLSToDel []string) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "AsyncDeleteUserURLs", ctx, URLSToDel)
}

// AsyncDeleteUserURLs indicates an expected call of AsyncDeleteUserURLs.
func (mr *MockServiceMockRecorder) AsyncDeleteUserURLs(ctx, URLSToDel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AsyncDeleteUserURLs", reflect.TypeOf((*MockService)(nil).AsyncDeleteUserURLs), ctx, URLSToDel)
}

// DelUserURLS mocks base method.
func (m *MockService) DelUserURLS(ctx context.Context, URLSToDel []string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelUserURLS", ctx, URLSToDel)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelUserURLS indicates an expected call of DelUserURLS.
func (mr *MockServiceMockRecorder) DelUserURLS(ctx, URLSToDel interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelUserURLS", reflect.TypeOf((*MockService)(nil).DelUserURLS), ctx, URLSToDel)
}

// GetBatchShortURL mocks base method.
func (m *MockService) GetBatchShortURL(ctx context.Context, batchURLRequests []models.URLRequest) ([]models.URLResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatchShortURL", ctx, batchURLRequests)
	ret0, _ := ret[0].([]models.URLResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBatchShortURL indicates an expected call of GetBatchShortURL.
func (mr *MockServiceMockRecorder) GetBatchShortURL(ctx, batchURLRequests interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatchShortURL", reflect.TypeOf((*MockService)(nil).GetBatchShortURL), ctx, batchURLRequests)
}

// GetOriginalURL mocks base method.
func (m *MockService) GetOriginalURL(ctx context.Context, shortURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalURL", ctx, shortURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginalURL indicates an expected call of GetOriginalURL.
func (mr *MockServiceMockRecorder) GetOriginalURL(ctx, shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalURL", reflect.TypeOf((*MockService)(nil).GetOriginalURL), ctx, shortURL)
}

// GetShortURL mocks base method.
func (m *MockService) GetShortURL(ctx context.Context, url string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortURL", ctx, url)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortURL indicates an expected call of GetShortURL.
func (mr *MockServiceMockRecorder) GetShortURL(ctx, url interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortURL", reflect.TypeOf((*MockService)(nil).GetShortURL), ctx, url)
}

// GetUserURLS mocks base method.
func (m *MockService) GetUserURLS(ctx context.Context) ([]models.URL, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserURLS", ctx)
	ret0, _ := ret[0].([]models.URL)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserURLS indicates an expected call of GetUserURLS.
func (mr *MockServiceMockRecorder) GetUserURLS(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserURLS", reflect.TypeOf((*MockService)(nil).GetUserURLS), ctx)
}
