// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/pedro00627/urblog/application (interfaces: LoadUsers)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	domain "github.com/pedro00627/urblog/domain"
)

// MockLoadUsers is a mock of LoadUsers interface.
type MockLoadUsers struct {
	ctrl     *gomock.Controller
	recorder *MockLoadUsersMockRecorder
}

// MockLoadUsersMockRecorder is the mock recorder for MockLoadUsers.
type MockLoadUsersMockRecorder struct {
	mock *MockLoadUsers
}

// NewMockLoadUsers creates a new mock instance.
func NewMockLoadUsers(ctrl *gomock.Controller) *MockLoadUsers {
	mock := &MockLoadUsers{ctrl: ctrl}
	mock.recorder = &MockLoadUsersMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockLoadUsers) EXPECT() *MockLoadUsersMockRecorder {
	return m.recorder
}

// Execute mocks base method.
func (m *MockLoadUsers) Execute(arg0 string) ([]domain.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Execute", arg0)
	ret0, _ := ret[0].([]domain.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Execute indicates an expected call of Execute.
func (mr *MockLoadUsersMockRecorder) Execute(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Execute", reflect.TypeOf((*MockLoadUsers)(nil).Execute), arg0)
}
