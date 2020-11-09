// Code generated by MockGen. DO NOT EDIT.
// Source: app.go

// Package app is a generated GoMock package.
package app

import (
	gomock "github.com/golang/mock/gomock"
	dom "github.com/powerman/go-monolith-example/internal/dom"
	sensitive "github.com/powerman/sensitive"
	reflect "reflect"
)

// MockAppl is a mock of Appl interface
type MockAppl struct {
	ctrl     *gomock.Controller
	recorder *MockApplMockRecorder
}

// MockApplMockRecorder is the mock recorder for MockAppl
type MockApplMockRecorder struct {
	mock *MockAppl
}

// NewMockAppl creates a new mock instance
func NewMockAppl(ctrl *gomock.Controller) *MockAppl {
	mock := &MockAppl{ctrl: ctrl}
	mock.recorder = &MockApplMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAppl) EXPECT() *MockApplMockRecorder {
	return m.recorder
}

// Register mocks base method
func (m *MockAppl) Register(arg0 Ctx, userID string, password sensitive.String, arg3 *User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", arg0, userID, password, arg3)
	ret0, _ := ret[0].(error)
	return ret0
}

// Register indicates an expected call of Register
func (mr *MockApplMockRecorder) Register(arg0, userID, password, arg3 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockAppl)(nil).Register), arg0, userID, password, arg3)
}

// LoginByUserID mocks base method
func (m *MockAppl) LoginByUserID(arg0 Ctx, userID string, password sensitive.String) (AccessToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginByUserID", arg0, userID, password)
	ret0, _ := ret[0].(AccessToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginByUserID indicates an expected call of LoginByUserID
func (mr *MockApplMockRecorder) LoginByUserID(arg0, userID, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginByUserID", reflect.TypeOf((*MockAppl)(nil).LoginByUserID), arg0, userID, password)
}

// LoginByEmail mocks base method
func (m *MockAppl) LoginByEmail(arg0 Ctx, email string, password sensitive.String) (AccessToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LoginByEmail", arg0, email, password)
	ret0, _ := ret[0].(AccessToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// LoginByEmail indicates an expected call of LoginByEmail
func (mr *MockApplMockRecorder) LoginByEmail(arg0, email, password interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LoginByEmail", reflect.TypeOf((*MockAppl)(nil).LoginByEmail), arg0, email, password)
}

// Authenticate mocks base method
func (m *MockAppl) Authenticate(arg0 Ctx, arg1 AccessToken) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", arg0, arg1)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authenticate indicates an expected call of Authenticate
func (mr *MockApplMockRecorder) Authenticate(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockAppl)(nil).Authenticate), arg0, arg1)
}

// Logout mocks base method
func (m *MockAppl) Logout(arg0 Ctx, arg1 AccessToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Logout", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Logout indicates an expected call of Logout
func (mr *MockApplMockRecorder) Logout(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Logout", reflect.TypeOf((*MockAppl)(nil).Logout), arg0, arg1)
}

// LogoutUser mocks base method
func (m *MockAppl) LogoutUser(arg0 Ctx, arg1 dom.UserName) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "LogoutUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// LogoutUser indicates an expected call of LogoutUser
func (mr *MockApplMockRecorder) LogoutUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "LogoutUser", reflect.TypeOf((*MockAppl)(nil).LogoutUser), arg0, arg1)
}

// MockRepo is a mock of Repo interface
type MockRepo struct {
	ctrl     *gomock.Controller
	recorder *MockRepoMockRecorder
}

// MockRepoMockRecorder is the mock recorder for MockRepo
type MockRepoMockRecorder struct {
	mock *MockRepo
}

// NewMockRepo creates a new mock instance
func NewMockRepo(ctrl *gomock.Controller) *MockRepo {
	mock := &MockRepo{ctrl: ctrl}
	mock.recorder = &MockRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepo) EXPECT() *MockRepoMockRecorder {
	return m.recorder
}

// AddUser mocks base method
func (m *MockRepo) AddUser(arg0 Ctx, arg1 User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddUser indicates an expected call of AddUser
func (mr *MockRepoMockRecorder) AddUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockRepo)(nil).AddUser), arg0, arg1)
}

// GetUser mocks base method
func (m *MockRepo) GetUser(arg0 Ctx, arg1 dom.UserName) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", arg0, arg1)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser
func (mr *MockRepoMockRecorder) GetUser(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockRepo)(nil).GetUser), arg0, arg1)
}

// GetUserByEmail mocks base method
func (m *MockRepo) GetUserByEmail(arg0 Ctx, arg1 string) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByEmail", arg0, arg1)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByEmail indicates an expected call of GetUserByEmail
func (mr *MockRepoMockRecorder) GetUserByEmail(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByEmail", reflect.TypeOf((*MockRepo)(nil).GetUserByEmail), arg0, arg1)
}

// GetUserByAccessToken mocks base method
func (m *MockRepo) GetUserByAccessToken(arg0 Ctx, arg1 AccessToken) (*User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByAccessToken", arg0, arg1)
	ret0, _ := ret[0].(*User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByAccessToken indicates an expected call of GetUserByAccessToken
func (mr *MockRepoMockRecorder) GetUserByAccessToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByAccessToken", reflect.TypeOf((*MockRepo)(nil).GetUserByAccessToken), arg0, arg1)
}

// AddAccessToken mocks base method
func (m *MockRepo) AddAccessToken(arg0 Ctx, arg1 dom.UserName) (AccessToken, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddAccessToken", arg0, arg1)
	ret0, _ := ret[0].(AccessToken)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddAccessToken indicates an expected call of AddAccessToken
func (mr *MockRepoMockRecorder) AddAccessToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddAccessToken", reflect.TypeOf((*MockRepo)(nil).AddAccessToken), arg0, arg1)
}

// DelAccessToken mocks base method
func (m *MockRepo) DelAccessToken(arg0 Ctx, arg1 AccessToken) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelAccessToken", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelAccessToken indicates an expected call of DelAccessToken
func (mr *MockRepoMockRecorder) DelAccessToken(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelAccessToken", reflect.TypeOf((*MockRepo)(nil).DelAccessToken), arg0, arg1)
}

// DelAccessTokens mocks base method
func (m *MockRepo) DelAccessTokens(arg0 Ctx, arg1 dom.UserName) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DelAccessTokens", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DelAccessTokens indicates an expected call of DelAccessTokens
func (mr *MockRepoMockRecorder) DelAccessTokens(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DelAccessTokens", reflect.TypeOf((*MockRepo)(nil).DelAccessTokens), arg0, arg1)
}
