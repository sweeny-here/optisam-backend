// Code generated by MockGen. DO NOT EDIT.
// Source: optisam-backend/auth-service/pkg/repository/v1 (interfaces: Repository)

package mock

import (
	context "context"
	v1 "optisam-backend/auth-service/pkg/repository/v1"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// CheckPassword mocks base method
func (m *MockRepository) CheckPassword(arg0 context.Context, arg1, arg2 string) (bool, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CheckPassword", arg0, arg1, arg2)
	ret0, _ := ret[0].(bool)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CheckPassword indicates an expected call of CheckPassword
func (mr *MockRepositoryMockRecorder) CheckPassword(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CheckPassword", reflect.TypeOf((*MockRepository)(nil).CheckPassword), arg0, arg1, arg2)
}

// IncreaseFailedLoginCount mocks base method
func (m *MockRepository) IncreaseFailedLoginCount(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IncreaseFailedLoginCount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// IncreaseFailedLoginCount indicates an expected call of IncreaseFailedLoginCount
func (mr *MockRepositoryMockRecorder) IncreaseFailedLoginCount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IncreaseFailedLoginCount", reflect.TypeOf((*MockRepository)(nil).IncreaseFailedLoginCount), arg0, arg1)
}

// ResetLoginCount mocks base method
func (m *MockRepository) ResetLoginCount(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ResetLoginCount", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ResetLoginCount indicates an expected call of ResetLoginCount
func (mr *MockRepositoryMockRecorder) ResetLoginCount(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ResetLoginCount", reflect.TypeOf((*MockRepository)(nil).ResetLoginCount), arg0, arg1)
}

// UserInfo mocks base method
func (m *MockRepository) UserInfo(arg0 context.Context, arg1 string) (*v1.UserInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserInfo", arg0, arg1)
	ret0, _ := ret[0].(*v1.UserInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserInfo indicates an expected call of UserInfo
func (mr *MockRepositoryMockRecorder) UserInfo(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserInfo", reflect.TypeOf((*MockRepository)(nil).UserInfo), arg0, arg1)
}

// UserOwnedGroupsDirect mocks base method
func (m *MockRepository) UserOwnedGroupsDirect(arg0 context.Context, arg1 string) ([]*v1.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UserOwnedGroupsDirect", arg0, arg1)
	ret0, _ := ret[0].([]*v1.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UserOwnedGroupsDirect indicates an expected call of UserOwnedGroupsDirect
func (mr *MockRepositoryMockRecorder) UserOwnedGroupsDirect(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UserOwnedGroupsDirect", reflect.TypeOf((*MockRepository)(nil).UserOwnedGroupsDirect), arg0, arg1)
}
