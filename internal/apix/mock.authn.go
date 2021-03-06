// Code generated by MockGen. DO NOT EDIT.
// Source: authn.go

// Package apix is a generated GoMock package.
package apix

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	dom "github.com/powerman/go-monolith-example/internal/dom"
)

// MockAuthn is a mock of Authn interface.
type MockAuthn struct {
	ctrl     *gomock.Controller
	recorder *MockAuthnMockRecorder
}

// MockAuthnMockRecorder is the mock recorder for MockAuthn.
type MockAuthnMockRecorder struct {
	mock *MockAuthn
}

// NewMockAuthn creates a new mock instance.
func NewMockAuthn(ctrl *gomock.Controller) *MockAuthn {
	mock := &MockAuthn{ctrl: ctrl}
	mock.recorder = &MockAuthnMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockAuthn) EXPECT() *MockAuthnMockRecorder {
	return m.recorder
}

// Authenticate mocks base method.
func (m *MockAuthn) Authenticate(arg0 Ctx, arg1 AccessToken) (dom.Auth, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", arg0, arg1)
	ret0, _ := ret[0].(dom.Auth)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockAuthnMockRecorder) Authenticate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockAuthn)(nil).Authenticate), arg0, arg1)
}
