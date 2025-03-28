// Code generated by MockGen. DO NOT EDIT.
// Source: dashboard.go
//
// Generated by this command:
//
//	mockgen -source=dashboard.go -destination=mocks/dashboard.go -package=malak_mocks
//

// Package malak_mocks is a generated GoMock package.
package malak_mocks

import (
	context "context"
	reflect "reflect"

	malak "github.com/ayinke-llc/malak"
	uuid "github.com/google/uuid"
	gomock "go.uber.org/mock/gomock"
)

// MockDashboardRepository is a mock of DashboardRepository interface.
type MockDashboardRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDashboardRepositoryMockRecorder
	isgomock struct{}
}

// MockDashboardRepositoryMockRecorder is the mock recorder for MockDashboardRepository.
type MockDashboardRepositoryMockRecorder struct {
	mock *MockDashboardRepository
}

// NewMockDashboardRepository creates a new mock instance.
func NewMockDashboardRepository(ctrl *gomock.Controller) *MockDashboardRepository {
	mock := &MockDashboardRepository{ctrl: ctrl}
	mock.recorder = &MockDashboardRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDashboardRepository) EXPECT() *MockDashboardRepositoryMockRecorder {
	return m.recorder
}

// AddChart mocks base method.
func (m *MockDashboardRepository) AddChart(arg0 context.Context, arg1 *malak.DashboardChart) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddChart", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddChart indicates an expected call of AddChart.
func (mr *MockDashboardRepositoryMockRecorder) AddChart(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddChart", reflect.TypeOf((*MockDashboardRepository)(nil).AddChart), arg0, arg1)
}

// Create mocks base method.
func (m *MockDashboardRepository) Create(arg0 context.Context, arg1 *malak.Dashboard) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockDashboardRepositoryMockRecorder) Create(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockDashboardRepository)(nil).Create), arg0, arg1)
}

// Get mocks base method.
func (m *MockDashboardRepository) Get(arg0 context.Context, arg1 malak.FetchDashboardOption) (malak.Dashboard, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(malak.Dashboard)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockDashboardRepositoryMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockDashboardRepository)(nil).Get), arg0, arg1)
}

// GetCharts mocks base method.
func (m *MockDashboardRepository) GetCharts(arg0 context.Context, arg1 malak.FetchDashboardChartsOption) ([]malak.DashboardChart, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCharts", arg0, arg1)
	ret0, _ := ret[0].([]malak.DashboardChart)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCharts indicates an expected call of GetCharts.
func (mr *MockDashboardRepositoryMockRecorder) GetCharts(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCharts", reflect.TypeOf((*MockDashboardRepository)(nil).GetCharts), arg0, arg1)
}

// GetDashboardPositions mocks base method.
func (m *MockDashboardRepository) GetDashboardPositions(arg0 context.Context, arg1 uuid.UUID) ([]malak.DashboardChartPosition, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDashboardPositions", arg0, arg1)
	ret0, _ := ret[0].([]malak.DashboardChartPosition)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDashboardPositions indicates an expected call of GetDashboardPositions.
func (mr *MockDashboardRepositoryMockRecorder) GetDashboardPositions(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDashboardPositions", reflect.TypeOf((*MockDashboardRepository)(nil).GetDashboardPositions), arg0, arg1)
}

// List mocks base method.
func (m *MockDashboardRepository) List(arg0 context.Context, arg1 malak.ListDashboardOptions) ([]malak.Dashboard, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]malak.Dashboard)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockDashboardRepositoryMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockDashboardRepository)(nil).List), arg0, arg1)
}

// RemoveChart mocks base method.
func (m *MockDashboardRepository) RemoveChart(arg0 context.Context, arg1, arg2 uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RemoveChart", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// RemoveChart indicates an expected call of RemoveChart.
func (mr *MockDashboardRepositoryMockRecorder) RemoveChart(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RemoveChart", reflect.TypeOf((*MockDashboardRepository)(nil).RemoveChart), arg0, arg1, arg2)
}

// UpdateDashboardPositions mocks base method.
func (m *MockDashboardRepository) UpdateDashboardPositions(arg0 context.Context, arg1 uuid.UUID, arg2 []malak.DashboardChartPosition) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateDashboardPositions", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateDashboardPositions indicates an expected call of UpdateDashboardPositions.
func (mr *MockDashboardRepositoryMockRecorder) UpdateDashboardPositions(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateDashboardPositions", reflect.TypeOf((*MockDashboardRepository)(nil).UpdateDashboardPositions), arg0, arg1, arg2)
}
