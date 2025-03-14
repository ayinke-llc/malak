package chart

import (
	"context"
	"io"
	"reflect"

	"github.com/adelowo/gulter"
	"go.uber.org/mock/gomock"
)

// MockStorage is a mock implementation of gulter.Storage
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageRecorder
}

type MockStorageRecorder struct {
	mock *MockStorage
}

func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageRecorder{mock}
	return mock
}

func (m *MockStorage) EXPECT() *MockStorageRecorder {
	return m.recorder
}

func (m *MockStorage) Upload(ctx context.Context, r io.Reader, opts *gulter.UploadFileOptions) (*gulter.UploadedFileMetadata, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Upload", ctx, r, opts)
	ret0, _ := ret[0].(*gulter.UploadedFileMetadata)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockStorageRecorder) Upload(ctx, r, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Upload", reflect.TypeOf((*MockStorage)(nil).Upload), ctx, r, opts)
}

func (m *MockStorage) Path(ctx context.Context, opts gulter.PathOptions) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Path", ctx, opts)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

func (mr *MockStorageRecorder) Path(ctx, opts interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Path", reflect.TypeOf((*MockStorage)(nil).Path), ctx, opts)
}

func (m *MockStorage) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

func (mr *MockStorageRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStorage)(nil).Close))
}
