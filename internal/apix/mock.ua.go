// Code generated by MockGen. DO NOT EDIT.
// Source: ua.go

// Package apix is a generated GoMock package.
package apix

import (
	http "net/http"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockUserAgent is a mock of UserAgent interface.
type MockUserAgent struct {
	ctrl     *gomock.Controller
	recorder *MockUserAgentMockRecorder
}

// MockUserAgentMockRecorder is the mock recorder for MockUserAgent.
type MockUserAgentMockRecorder struct {
	mock *MockUserAgent
}

// NewMockUserAgent creates a new mock instance.
func NewMockUserAgent(ctrl *gomock.Controller) *MockUserAgent {
	mock := &MockUserAgent{ctrl: ctrl}
	mock.recorder = &MockUserAgentMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserAgent) EXPECT() *MockUserAgentMockRecorder {
	return m.recorder
}

// Do mocks base method.
func (m *MockUserAgent) Do(ctx Ctx, req *http.Request, skip int) (*http.Response, []byte, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", ctx, req, skip)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].([]byte)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// Do indicates an expected call of Do.
func (mr *MockUserAgentMockRecorder) Do(ctx, req, skip interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockUserAgent)(nil).Do), ctx, req, skip)
}

// Log mocks base method.
func (m *MockUserAgent) Log(ctx Ctx, resp *http.Response, body []byte) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Log", ctx, resp, body)
}

// Log indicates an expected call of Log.
func (mr *MockUserAgentMockRecorder) Log(ctx, resp, body interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Log", reflect.TypeOf((*MockUserAgent)(nil).Log), ctx, resp, body)
}
