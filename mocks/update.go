// Code generated by MockGen. DO NOT EDIT.
// Source: update.go
//
// Generated by this command:
//
//	mockgen -source=update.go -destination=mocks/update.go -package=malak_mocks
//

// Package malak_mocks is a generated GoMock package.
package malak_mocks

import (
	context "context"
	reflect "reflect"

	malak "github.com/ayinke-llc/malak"
	gomock "go.uber.org/mock/gomock"
)

// MockUpdateRepository is a mock of UpdateRepository interface.
type MockUpdateRepository struct {
	ctrl     *gomock.Controller
	recorder *MockUpdateRepositoryMockRecorder
}

// MockUpdateRepositoryMockRecorder is the mock recorder for MockUpdateRepository.
type MockUpdateRepositoryMockRecorder struct {
	mock *MockUpdateRepository
}

// NewMockUpdateRepository creates a new mock instance.
func NewMockUpdateRepository(ctrl *gomock.Controller) *MockUpdateRepository {
	mock := &MockUpdateRepository{ctrl: ctrl}
	mock.recorder = &MockUpdateRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUpdateRepository) EXPECT() *MockUpdateRepositoryMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUpdateRepository) Create(arg0 context.Context, arg1 *malak.Update) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockUpdateRepositoryMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUpdateRepository)(nil).Create), arg0, arg1)
}

// List mocks base method.
func (m *MockUpdateRepository) List(arg0 context.Context, arg1 malak.ListUpdateOptions) ([]malak.Update, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]malak.Update)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockUpdateRepositoryMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockUpdateRepository)(nil).List), arg0, arg1)
}