// Code generated by protoc-gen-go-grpc-mock. DO NOT EDIT.
// source: proto/api/user/v1/user.proto

package user

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
)

// MockUserServiceClient is a mock of UserServiceClient interface.
type MockUserServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceClientMockRecorder
}

// MockUserServiceClientMockRecorder is the mock recorder for MockUserServiceClient.
type MockUserServiceClientMockRecorder struct {
	mock *MockUserServiceClient
}

// NewMockUserServiceClient creates a new mock instance.
func NewMockUserServiceClient(ctrl *gomock.Controller) *MockUserServiceClient {
	mock := &MockUserServiceClient{ctrl: ctrl}
	mock.recorder = &MockUserServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserServiceClient) EXPECT() *MockUserServiceClientMockRecorder {
	return m.recorder
}

// Authenticate mocks base method.
func (m *MockUserServiceClient) Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Authenticate", varargs...)
	ret0, _ := ret[0].(*AuthenticateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockUserServiceClientMockRecorder) Authenticate(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockUserServiceClient)(nil).Authenticate), varargs...)
}

// Register mocks base method.
func (m *MockUserServiceClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Register", varargs...)
	ret0, _ := ret[0].(*RegisterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockUserServiceClientMockRecorder) Register(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserServiceClient)(nil).Register), varargs...)
}

// MockUserServiceServer is a mock of UserServiceServer interface.
type MockUserServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUserServiceServerMockRecorder
}

// MockUserServiceServerMockRecorder is the mock recorder for MockUserServiceServer.
type MockUserServiceServerMockRecorder struct {
	mock *MockUserServiceServer
}

// NewMockUserServiceServer creates a new mock instance.
func NewMockUserServiceServer(ctrl *gomock.Controller) *MockUserServiceServer {
	mock := &MockUserServiceServer{ctrl: ctrl}
	mock.recorder = &MockUserServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserServiceServer) EXPECT() *MockUserServiceServerMockRecorder {
	return m.recorder
}

// Authenticate mocks base method.
func (m *MockUserServiceServer) Authenticate(ctx context.Context, in *AuthenticateRequest) (*AuthenticateResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Authenticate", ctx, in)
	ret0, _ := ret[0].(*AuthenticateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Authenticate indicates an expected call of Authenticate.
func (mr *MockUserServiceServerMockRecorder) Authenticate(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Authenticate", reflect.TypeOf((*MockUserServiceServer)(nil).Authenticate), ctx, in)
}

// Register mocks base method.
func (m *MockUserServiceServer) Register(ctx context.Context, in *RegisterRequest) (*RegisterResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Register", ctx, in)
	ret0, _ := ret[0].(*RegisterResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Register indicates an expected call of Register.
func (mr *MockUserServiceServerMockRecorder) Register(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Register", reflect.TypeOf((*MockUserServiceServer)(nil).Register), ctx, in)
}
