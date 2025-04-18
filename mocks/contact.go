// Code generated by MockGen. DO NOT EDIT.
// Source: contact.go
//
// Generated by this command:
//
//	mockgen -source=contact.go -destination=mocks/contact.go -package=malak_mocks
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

// MockContactRepository is a mock of ContactRepository interface.
type MockContactRepository struct {
	ctrl     *gomock.Controller
	recorder *MockContactRepositoryMockRecorder
	isgomock struct{}
}

// MockContactRepositoryMockRecorder is the mock recorder for MockContactRepository.
type MockContactRepositoryMockRecorder struct {
	mock *MockContactRepository
}

// NewMockContactRepository creates a new mock instance.
func NewMockContactRepository(ctrl *gomock.Controller) *MockContactRepository {
	mock := &MockContactRepository{ctrl: ctrl}
	mock.recorder = &MockContactRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockContactRepository) EXPECT() *MockContactRepositoryMockRecorder {
	return m.recorder
}

// All mocks base method.
func (m *MockContactRepository) All(arg0 context.Context, arg1 uuid.UUID) ([]malak.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "All", arg0, arg1)
	ret0, _ := ret[0].([]malak.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// All indicates an expected call of All.
func (mr *MockContactRepositoryMockRecorder) All(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "All", reflect.TypeOf((*MockContactRepository)(nil).All), arg0, arg1)
}

// Create mocks base method.
func (m *MockContactRepository) Create(arg0 context.Context, arg1 ...*malak.Contact) error {
	m.ctrl.T.Helper()
	varargs := []any{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Create", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create.
func (mr *MockContactRepositoryMockRecorder) Create(arg0 any, arg1 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockContactRepository)(nil).Create), varargs...)
}

// Delete mocks base method.
func (m *MockContactRepository) Delete(arg0 context.Context, arg1 *malak.Contact) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockContactRepositoryMockRecorder) Delete(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockContactRepository)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockContactRepository) Get(arg0 context.Context, arg1 malak.FetchContactOptions) (*malak.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*malak.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockContactRepositoryMockRecorder) Get(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockContactRepository)(nil).Get), arg0, arg1)
}

// List mocks base method.
func (m *MockContactRepository) List(arg0 context.Context, arg1 malak.ListContactOptions) ([]malak.Contact, int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "List", arg0, arg1)
	ret0, _ := ret[0].([]malak.Contact)
	ret1, _ := ret[1].(int64)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// List indicates an expected call of List.
func (mr *MockContactRepositoryMockRecorder) List(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockContactRepository)(nil).List), arg0, arg1)
}

// Overview mocks base method.
func (m *MockContactRepository) Overview(arg0 context.Context, arg1 uuid.UUID) (*malak.ContactOverview, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Overview", arg0, arg1)
	ret0, _ := ret[0].(*malak.ContactOverview)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Overview indicates an expected call of Overview.
func (mr *MockContactRepositoryMockRecorder) Overview(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Overview", reflect.TypeOf((*MockContactRepository)(nil).Overview), arg0, arg1)
}

// Search mocks base method.
func (m *MockContactRepository) Search(arg0 context.Context, arg1 malak.SearchContactOptions) ([]malak.Contact, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Search", arg0, arg1)
	ret0, _ := ret[0].([]malak.Contact)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Search indicates an expected call of Search.
func (mr *MockContactRepositoryMockRecorder) Search(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Search", reflect.TypeOf((*MockContactRepository)(nil).Search), arg0, arg1)
}

// Update mocks base method.
func (m *MockContactRepository) Update(arg0 context.Context, arg1 *malak.Contact) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Update indicates an expected call of Update.
func (mr *MockContactRepositoryMockRecorder) Update(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockContactRepository)(nil).Update), arg0, arg1)
}
