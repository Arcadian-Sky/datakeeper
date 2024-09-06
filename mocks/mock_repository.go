// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/server/repository/repository.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	os "os"
	reflect "reflect"

	model "github.com/Arcadian-Sky/datakkeeper/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockFileRepository is a mock of FileRepository interface.
type MockFileRepository struct {
	ctrl     *gomock.Controller
	recorder *MockFileRepositoryMockRecorder
}

// MockFileRepositoryMockRecorder is the mock recorder for MockFileRepository.
type MockFileRepositoryMockRecorder struct {
	mock *MockFileRepository
}

// NewMockFileRepository creates a new mock instance.
func NewMockFileRepository(ctrl *gomock.Controller) *MockFileRepository {
	mock := &MockFileRepository{ctrl: ctrl}
	mock.recorder = &MockFileRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockFileRepository) EXPECT() *MockFileRepositoryMockRecorder {
	return m.recorder
}

// CreateContainer mocks base method.
func (m *MockFileRepository) CreateContainer(ctx context.Context, user *model.User) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateContainer", ctx, user)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateContainer indicates an expected call of CreateContainer.
func (mr *MockFileRepositoryMockRecorder) CreateContainer(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateContainer", reflect.TypeOf((*MockFileRepository)(nil).CreateContainer), ctx, user)
}

// DeleteFile mocks base method.
func (m *MockFileRepository) DeleteFile(ctx context.Context, fileID string, user *model.User) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFile", ctx, fileID, user)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockFileRepositoryMockRecorder) DeleteFile(ctx, fileID, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockFileRepository)(nil).DeleteFile), ctx, fileID, user)
}

// GetFile mocks base method.
func (m *MockFileRepository) GetFile(ctx context.Context, fileID string, user *model.User) (*os.File, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFile", ctx, fileID, user)
	ret0, _ := ret[0].(*os.File)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFile indicates an expected call of GetFile.
func (mr *MockFileRepositoryMockRecorder) GetFile(ctx, fileID, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFile", reflect.TypeOf((*MockFileRepository)(nil).GetFile), ctx, fileID, user)
}

// GetFileList mocks base method.
func (m *MockFileRepository) GetFileList(ctx context.Context, user *model.User) ([]model.FileItem, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileList", ctx, user)
	ret0, _ := ret[0].([]model.FileItem)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileList indicates an expected call of GetFileList.
func (mr *MockFileRepositoryMockRecorder) GetFileList(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileList", reflect.TypeOf((*MockFileRepository)(nil).GetFileList), ctx, user)
}

// UploadFile mocks base method.
func (m *MockFileRepository) UploadFile(ctx context.Context, user *model.User, objectName string, file *os.File) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", ctx, user, objectName, file)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockFileRepositoryMockRecorder) UploadFile(ctx, user, objectName, file interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockFileRepository)(nil).UploadFile), ctx, user, objectName, file)
}