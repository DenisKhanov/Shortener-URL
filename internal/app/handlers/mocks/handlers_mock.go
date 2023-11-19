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

// GetBatchJSONShortURL mocks base method.
func (m *MockService) GetBatchJSONShortURL(ctx context.Context, batchURLRequests []models.URLRequest) ([]models.URLResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBatchJSONShortURL", ctx, batchURLRequests)
	ret0, _ := ret[0].([]models.URLResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetBatchJSONShortURL indicates an expected call of GetBatchJSONShortURL.
func (mr *MockServiceMockRecorder) GetBatchJSONShortURL(ctx, batchURLRequests interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBatchJSONShortURL", reflect.TypeOf((*MockService)(nil).GetBatchJSONShortURL), ctx, batchURLRequests)
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
