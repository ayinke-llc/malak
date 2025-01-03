// Code generated by MockGen. DO NOT EDIT.
// Source: plan.go
//
// Generated by this command:
//
//	mockgen -source=plan.go -destination=mocks/plan.go -package=malak_mocks
//

// Package malak_mocks is a generated GoMock package.
package malak_mocks

import (
	context "context"
	reflect "reflect"

	malak "github.com/ayinke-llc/malak"
	gomock "go.uber.org/mock/gomock"
)

// MockPlanRepository is a mock of PlanRepository interface.
type MockPlanRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPlanRepositoryMockRecorder
	isgomock struct{}
}

// MockPlanRepositoryMockRecorder is the mock recorder for MockPlanRepository.
type MockPlanRepositoryMockRecorder struct {
	mock *MockPlanRepository
}

// NewMockPlanRepository creates a new mock instance.
func NewMockPlanRepository(ctrl *gomock.Controller) *MockPlanRepository {
	mock := &MockPlanRepository{ctrl: ctrl}
	mock.recorder = &MockPlanRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPlanRepository) EXPECT() *MockPlanRepositoryMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockPlanRepository) Get(arg0 context.Context, arg1 *malak.FetchPlanOptions) (*malak.Plan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*malak.Plan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockPlanRepositoryMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockPlanRepository)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockPlanRepository) List(arg0 context.Context) ([]*malak.Plan, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0)
	ret0, _ := ret[0].([]*malak.Plan)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List.
func (mr *MockPlanRepositoryMockRecorder) List(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockPlanRepository)(nil).List), arg0)
}
