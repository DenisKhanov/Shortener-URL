// Code generated by MockGen. DO NOT EDIT.
// Source: storage.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// GetID mocks base method.
func (m *MockRepository) GetID() int {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetID")
	ret0, _ := ret[0].(int)
	return ret0
}

// GetID indicates an expected call of GetID.
func (mr *MockRepositoryMockRecorder) GetID() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetID", reflect.TypeOf((*MockRepository)(nil).GetID))
}

// GetOriginalURL mocks base method.
func (m *MockRepository) GetOriginalURL(shortURL string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalURL", shortURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetOriginalURL indicates an expected call of GetOriginalURL.
func (mr *MockRepositoryMockRecorder) GetOriginalURL(shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalURL", reflect.TypeOf((*MockRepository)(nil).GetOriginalURL), shortURL)
}

// GetShortURL mocks base method.
func (m *MockRepository) GetShortURL(originalURL string) (string, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortURL", originalURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetShortURL indicates an expected call of GetShortURL.
func (mr *MockRepositoryMockRecorder) GetShortURL(originalURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortURL", reflect.TypeOf((*MockRepository)(nil).GetShortURL), originalURL)
}

// StoreURL mocks base method.
func (m *MockRepository) StoreURL(originalURL, shortURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreURL", originalURL, shortURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreURL indicates an expected call of StoreURL.
func (mr *MockRepositoryMockRecorder) StoreURL(originalURL, shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreURL", reflect.TypeOf((*MockRepository)(nil).StoreURL), originalURL, shortURL)
}