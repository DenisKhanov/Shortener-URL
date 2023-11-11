// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	models "github.com/DenisKhanov/shorterURL/internal/app/models"
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

// GetOriginalURLFromDB mocks base method.
func (m *MockRepository) GetOriginalURLFromDB(shortURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOriginalURLFromDB", shortURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOriginalURLFromDB indicates an expected call of GetOriginalURLFromDB.
func (mr *MockRepositoryMockRecorder) GetOriginalURLFromDB(shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOriginalURLFromDB", reflect.TypeOf((*MockRepository)(nil).GetOriginalURLFromDB), shortURL)
}

// GetShortURLFromDB mocks base method.
func (m *MockRepository) GetShortURLFromDB(originalURL string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortURLFromDB", originalURL)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortURLFromDB indicates an expected call of GetShortURLFromDB.
func (mr *MockRepositoryMockRecorder) GetShortURLFromDB(originalURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortURLFromDB", reflect.TypeOf((*MockRepository)(nil).GetShortURLFromDB), originalURL)
}

// StoreURLInDB mocks base method.
func (m *MockRepository) StoreURLInDB(originalURL, shortURL string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreURLInDB", originalURL, shortURL)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreURLInDB indicates an expected call of StoreURLInDB.
func (mr *MockRepositoryMockRecorder) StoreURLInDB(originalURL, shortURL interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreURLInDB", reflect.TypeOf((*MockRepository)(nil).StoreURLInDB), originalURL, shortURL)
}

// MockURLInMemoryRepository is a mock of URLInMemoryRepository interface.
type MockURLInMemoryRepository struct {
	ctrl     *gomock.Controller
	recorder *MockURLInMemoryRepositoryMockRecorder
}

// MockURLInMemoryRepositoryMockRecorder is the mock recorder for MockURLInMemoryRepository.
type MockURLInMemoryRepositoryMockRecorder struct {
	mock *MockURLInMemoryRepository
}

// NewMockURLInMemoryRepository creates a new mock instance.
func NewMockURLInMemoryRepository(ctrl *gomock.Controller) *MockURLInMemoryRepository {
	mock := &MockURLInMemoryRepository{ctrl: ctrl}
	mock.recorder = &MockURLInMemoryRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLInMemoryRepository) EXPECT() *MockURLInMemoryRepositoryMockRecorder {
	return m.recorder
}

// SaveBatchToFile mocks base method.
func (m *MockURLInMemoryRepository) SaveBatchToFile() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveBatchToFile")
	ret0, _ := ret[0].(error)
	return ret0
}

// SaveBatchToFile indicates an expected call of SaveBatchToFile.
func (mr *MockURLInMemoryRepositoryMockRecorder) SaveBatchToFile() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveBatchToFile", reflect.TypeOf((*MockURLInMemoryRepository)(nil).SaveBatchToFile))
}

// MockURLInDBRepository is a mock of URLInDBRepository interface.
type MockURLInDBRepository struct {
	ctrl     *gomock.Controller
	recorder *MockURLInDBRepositoryMockRecorder
}

// MockURLInDBRepositoryMockRecorder is the mock recorder for MockURLInDBRepository.
type MockURLInDBRepositoryMockRecorder struct {
	mock *MockURLInDBRepository
}

// NewMockURLInDBRepository creates a new mock instance.
func NewMockURLInDBRepository(ctrl *gomock.Controller) *MockURLInDBRepository {
	mock := &MockURLInDBRepository{ctrl: ctrl}
	mock.recorder = &MockURLInDBRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockURLInDBRepository) EXPECT() *MockURLInDBRepositoryMockRecorder {
	return m.recorder
}

// GetShortBatchURLFromDB mocks base method.
func (m *MockURLInDBRepository) GetShortBatchURLFromDB(batchURLRequests []models.URLRequest) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetShortBatchURLFromDB", batchURLRequests)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetShortBatchURLFromDB indicates an expected call of GetShortBatchURLFromDB.
func (mr *MockURLInDBRepositoryMockRecorder) GetShortBatchURLFromDB(batchURLRequests interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetShortBatchURLFromDB", reflect.TypeOf((*MockURLInDBRepository)(nil).GetShortBatchURLFromDB), batchURLRequests)
}

// StoreBatchURLInDB mocks base method.
func (m *MockURLInDBRepository) StoreBatchURLInDB(batchURLtoStores map[string]string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StoreBatchURLInDB", batchURLtoStores)
	ret0, _ := ret[0].(error)
	return ret0
}

// StoreBatchURLInDB indicates an expected call of StoreBatchURLInDB.
func (mr *MockURLInDBRepositoryMockRecorder) StoreBatchURLInDB(batchURLtoStores interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StoreBatchURLInDB", reflect.TypeOf((*MockURLInDBRepository)(nil).StoreBatchURLInDB), batchURLtoStores)
}

// MockEncoder is a mock of Encoder interface.
type MockEncoder struct {
	ctrl     *gomock.Controller
	recorder *MockEncoderMockRecorder
}

// MockEncoderMockRecorder is the mock recorder for MockEncoder.
type MockEncoderMockRecorder struct {
	mock *MockEncoder
}

// NewMockEncoder creates a new mock instance.
func NewMockEncoder(ctrl *gomock.Controller) *MockEncoder {
	mock := &MockEncoder{ctrl: ctrl}
	mock.recorder = &MockEncoderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEncoder) EXPECT() *MockEncoderMockRecorder {
	return m.recorder
}

// CryptoBase62Encode mocks base method.
func (m *MockEncoder) CryptoBase62Encode() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CryptoBase62Encode")
	ret0, _ := ret[0].(string)
	return ret0
}

// CryptoBase62Encode indicates an expected call of CryptoBase62Encode.
func (mr *MockEncoderMockRecorder) CryptoBase62Encode() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CryptoBase62Encode", reflect.TypeOf((*MockEncoder)(nil).CryptoBase62Encode))
}
