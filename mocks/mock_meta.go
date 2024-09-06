// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/server/repository/meta.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	model "github.com/Arcadian-Sky/datakkeeper/internal/model"
	gomock "github.com/golang/mock/gomock"
)

// MockDataRepository is a mock of DataRepository interface.
type MockDataRepository struct {
	ctrl     *gomock.Controller
	recorder *MockDataRepositoryMockRecorder
}

// MockDataRepositoryMockRecorder is the mock recorder for MockDataRepository.
type MockDataRepositoryMockRecorder struct {
	mock *MockDataRepository
}

// NewMockDataRepository creates a new mock instance.
func NewMockDataRepository(ctrl *gomock.Controller) *MockDataRepository {
	mock := &MockDataRepository{ctrl: ctrl}
	mock.recorder = &MockDataRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataRepository) EXPECT() *MockDataRepositoryMockRecorder {
	return m.recorder
}

// GetList mocks base method.
func (m *MockDataRepository) GetList(ctx context.Context, user *model.User) ([]model.Data, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetList", ctx, user)
	ret0, _ := ret[0].([]model.Data)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetList indicates an expected call of GetList.
func (mr *MockDataRepositoryMockRecorder) GetList(ctx, user interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetList", reflect.TypeOf((*MockDataRepository)(nil).GetList), ctx, user)
}

// Save mocks base method.
func (m *MockDataRepository) Save(ctx context.Context, data *model.Data) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", ctx, data)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Save indicates an expected call of Save.
func (mr *MockDataRepositoryMockRecorder) Save(ctx, data interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockDataRepository)(nil).Save), ctx, data)
}