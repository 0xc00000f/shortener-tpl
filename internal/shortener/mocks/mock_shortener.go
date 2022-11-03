// Code generated by MockGen. DO NOT EDIT.
// Source: internal/shortener/shortener.go

// Package mock_shortener is a generated GoMock package.
package mock_shortener

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
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

// Get mocks base method.
func (m *MockShortener) Get(short string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", short)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockShortenerMockRecorder) Get(short interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockShortener)(nil).Get), short)
}

// GetAll mocks base method.
func (m *MockShortener) GetAll() (map[string]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAll")
	ret0, _ := ret[0].(map[string]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAll indicates an expected call of GetAll.
func (mr *MockShortenerMockRecorder) GetAll() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAll", reflect.TypeOf((*MockShortener)(nil).GetAll))
}

// Short mocks base method.
func (m *MockShortener) Short(long string) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Short", long)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Short indicates an expected call of Short.
func (mr *MockShortenerMockRecorder) Short(long interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Short", reflect.TypeOf((*MockShortener)(nil).Short), long)
}