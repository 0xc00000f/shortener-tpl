// Code generated by MockGen. DO NOT EDIT.
// Source: internal/shortener/shortener.go

// Package mock_shortener is a generated GoMock package.
package mock_shortener

import (
	context "context"
	reflect "reflect"

	models "github.com/0xc00000f/shortener-tpl/internal/models"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockShortener is a mock of Shortener interface.
type MockShortener struct {
	ctrl     *gomock.Controller
	recorder *MockShortenerMockRecorder
}

// MockShortenerMockRecorder is the mock recorder for MockShortener.
type MockShortenerMockRecorder struct {
	mock *MockShortener
}

// NewMockShortener creates a new mock instance.
func NewMockShortener(ctrl *gomock.Controller) *MockShortener {
	mock := &MockShortener{ctrl: ctrl}
	mock.recorder = &MockShortenerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockShortener) EXPECT() *MockShortenerMockRecorder {
	return m.recorder
}

// Delete mocks base method.
func (m *MockShortener) Delete(ctx context.Context, data []models.URL) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, data)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockShortenerMockRecorder) Delete(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockShortener)(nil).Delete), ctx, data)
}

// Get mocks base method.
func (m *MockShortener) Get(ctx context.Context, short string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, short)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockShortenerMockRecorder) Get(ctx, short interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockShortener)(nil).Get), ctx, short)
}

// GetAll mocks base method.
func (m *MockShortener) GetAll(ctx context.Context, userID uuid.UUID) (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll", ctx, userID)
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockShortenerMockRecorder) GetAll(ctx, userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockShortener)(nil).GetAll), ctx, userID)
}

// Short mocks base method.
func (m *MockShortener) Short(ctx context.Context, userID uuid.UUID, long string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Short", ctx, userID, long)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Short indicates an expected call of Short.
func (mr *MockShortenerMockRecorder) Short(ctx, userID, long interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Short", reflect.TypeOf((*MockShortener)(nil).Short), ctx, userID, long)
}
