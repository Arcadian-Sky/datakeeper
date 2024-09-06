// Code generated by protoc-gen-go-grpc-mock. DO NOT EDIT.
// source: proto/api/service/v1/service.proto

package data

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
	metadata "google.golang.org/grpc/metadata"
)

// MockDataKeeperService_UploadFileClient is a mock of DataKeeperService_UploadFileClient interface.
type MockDataKeeperService_UploadFileClient struct {
	ctrl     *gomock.Controller
	recorder *MockDataKeeperService_UploadFileClientMockRecorder
}

// MockDataKeeperService_UploadFileClientMockRecorder is the mock recorder for MockDataKeeperService_UploadFileClient.
type MockDataKeeperService_UploadFileClientMockRecorder struct {
	mock *MockDataKeeperService_UploadFileClient
}

// NewMockDataKeeperService_UploadFileClient creates a new mock instance.
func NewMockDataKeeperService_UploadFileClient(ctrl *gomock.Controller) *MockDataKeeperService_UploadFileClient {
	mock := &MockDataKeeperService_UploadFileClient{ctrl: ctrl}
	mock.recorder = &MockDataKeeperService_UploadFileClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataKeeperService_UploadFileClient) EXPECT() *MockDataKeeperService_UploadFileClientMockRecorder {
	return m.recorder
}

// CloseAndRecv mocks base method.
func (m *MockDataKeeperService_UploadFileClient) CloseAndRecv() (*UploadStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseAndRecv")
	ret0, _ := ret[0].(*UploadStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CloseAndRecv indicates an expected call of CloseAndRecv.
func (mr *MockDataKeeperService_UploadFileClientMockRecorder) CloseAndRecv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseAndRecv", reflect.TypeOf((*MockDataKeeperService_UploadFileClient)(nil).CloseAndRecv))
}

// CloseSend mocks base method.
func (m *MockDataKeeperService_UploadFileClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockDataKeeperService_UploadFileClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockDataKeeperService_UploadFileClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockDataKeeperService_UploadFileClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDataKeeperService_UploadFileClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDataKeeperService_UploadFileClient)(nil).Context))
}

// Header mocks base method.
func (m *MockDataKeeperService_UploadFileClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockDataKeeperService_UploadFileClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockDataKeeperService_UploadFileClient)(nil).Header))
}

// RecvMsg mocks base method.
func (m *MockDataKeeperService_UploadFileClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockDataKeeperService_UploadFileClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockDataKeeperService_UploadFileClient)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockDataKeeperService_UploadFileClient) Send(arg0 *FileChunk) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockDataKeeperService_UploadFileClientMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockDataKeeperService_UploadFileClient)(nil).Send), arg0)
}

// SendMsg mocks base method.
func (m *MockDataKeeperService_UploadFileClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockDataKeeperService_UploadFileClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDataKeeperService_UploadFileClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method.
func (m *MockDataKeeperService_UploadFileClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockDataKeeperService_UploadFileClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockDataKeeperService_UploadFileClient)(nil).Trailer))
}

// MockDataKeeperService_UploadFileServer is a mock of DataKeeperService_UploadFileServer interface.
type MockDataKeeperService_UploadFileServer struct {
	ctrl     *gomock.Controller
	recorder *MockDataKeeperService_UploadFileServerMockRecorder
}

// MockDataKeeperService_UploadFileServerMockRecorder is the mock recorder for MockDataKeeperService_UploadFileServer.
type MockDataKeeperService_UploadFileServerMockRecorder struct {
	mock *MockDataKeeperService_UploadFileServer
}

// NewMockDataKeeperService_UploadFileServer creates a new mock instance.
func NewMockDataKeeperService_UploadFileServer(ctrl *gomock.Controller) *MockDataKeeperService_UploadFileServer {
	mock := &MockDataKeeperService_UploadFileServer{ctrl: ctrl}
	mock.recorder = &MockDataKeeperService_UploadFileServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataKeeperService_UploadFileServer) EXPECT() *MockDataKeeperService_UploadFileServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockDataKeeperService_UploadFileServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDataKeeperService_UploadFileServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDataKeeperService_UploadFileServer)(nil).Context))
}

// Recv mocks base method.
func (m *MockDataKeeperService_UploadFileServer) Recv() (*UploadStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*UploadStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockDataKeeperService_UploadFileServerMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockDataKeeperService_UploadFileServer)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockDataKeeperService_UploadFileServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockDataKeeperService_UploadFileServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockDataKeeperService_UploadFileServer)(nil).RecvMsg), arg0)
}

// SendAndClose mocks base method.
func (m *MockDataKeeperService_UploadFileServer) SendAndClose(arg0 *FileChunk) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendAndClose", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendAndClose indicates an expected call of SendAndClose.
func (mr *MockDataKeeperService_UploadFileServerMockRecorder) SendAndClose(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendAndClose", reflect.TypeOf((*MockDataKeeperService_UploadFileServer)(nil).SendAndClose), arg0)
}

// SendHeader mocks base method.
func (m *MockDataKeeperService_UploadFileServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockDataKeeperService_UploadFileServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockDataKeeperService_UploadFileServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m *MockDataKeeperService_UploadFileServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockDataKeeperService_UploadFileServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDataKeeperService_UploadFileServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method.
func (m *MockDataKeeperService_UploadFileServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockDataKeeperService_UploadFileServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockDataKeeperService_UploadFileServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockDataKeeperService_UploadFileServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockDataKeeperService_UploadFileServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockDataKeeperService_UploadFileServer)(nil).SetTrailer), arg0)
}

// MockDataKeeperService_GetFileClient is a mock of DataKeeperService_GetFileClient interface.
type MockDataKeeperService_GetFileClient struct {
	ctrl     *gomock.Controller
	recorder *MockDataKeeperService_GetFileClientMockRecorder
}

// MockDataKeeperService_GetFileClientMockRecorder is the mock recorder for MockDataKeeperService_GetFileClient.
type MockDataKeeperService_GetFileClientMockRecorder struct {
	mock *MockDataKeeperService_GetFileClient
}

// NewMockDataKeeperService_GetFileClient creates a new mock instance.
func NewMockDataKeeperService_GetFileClient(ctrl *gomock.Controller) *MockDataKeeperService_GetFileClient {
	mock := &MockDataKeeperService_GetFileClient{ctrl: ctrl}
	mock.recorder = &MockDataKeeperService_GetFileClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataKeeperService_GetFileClient) EXPECT() *MockDataKeeperService_GetFileClientMockRecorder {
	return m.recorder
}

// CloseSend mocks base method.
func (m *MockDataKeeperService_GetFileClient) CloseSend() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CloseSend")
	ret0, _ := ret[0].(error)
	return ret0
}

// CloseSend indicates an expected call of CloseSend.
func (mr *MockDataKeeperService_GetFileClientMockRecorder) CloseSend() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CloseSend", reflect.TypeOf((*MockDataKeeperService_GetFileClient)(nil).CloseSend))
}

// Context mocks base method.
func (m *MockDataKeeperService_GetFileClient) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDataKeeperService_GetFileClientMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDataKeeperService_GetFileClient)(nil).Context))
}

// Header mocks base method.
func (m *MockDataKeeperService_GetFileClient) Header() (metadata.MD, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Header")
	ret0, _ := ret[0].(metadata.MD)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Header indicates an expected call of Header.
func (mr *MockDataKeeperService_GetFileClientMockRecorder) Header() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Header", reflect.TypeOf((*MockDataKeeperService_GetFileClient)(nil).Header))
}

// Recv mocks base method.
func (m *MockDataKeeperService_GetFileClient) Recv() (*FileChunk, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Recv")
	ret0, _ := ret[0].(*FileChunk)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Recv indicates an expected call of Recv.
func (mr *MockDataKeeperService_GetFileClientMockRecorder) Recv() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Recv", reflect.TypeOf((*MockDataKeeperService_GetFileClient)(nil).Recv))
}

// RecvMsg mocks base method.
func (m *MockDataKeeperService_GetFileClient) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockDataKeeperService_GetFileClientMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockDataKeeperService_GetFileClient)(nil).RecvMsg), arg0)
}

// SendMsg mocks base method.
func (m *MockDataKeeperService_GetFileClient) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockDataKeeperService_GetFileClientMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDataKeeperService_GetFileClient)(nil).SendMsg), arg0)
}

// Trailer mocks base method.
func (m *MockDataKeeperService_GetFileClient) Trailer() metadata.MD {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Trailer")
	ret0, _ := ret[0].(metadata.MD)
	return ret0
}

// Trailer indicates an expected call of Trailer.
func (mr *MockDataKeeperService_GetFileClientMockRecorder) Trailer() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Trailer", reflect.TypeOf((*MockDataKeeperService_GetFileClient)(nil).Trailer))
}

// MockDataKeeperService_GetFileServer is a mock of DataKeeperService_GetFileServer interface.
type MockDataKeeperService_GetFileServer struct {
	ctrl     *gomock.Controller
	recorder *MockDataKeeperService_GetFileServerMockRecorder
}

// MockDataKeeperService_GetFileServerMockRecorder is the mock recorder for MockDataKeeperService_GetFileServer.
type MockDataKeeperService_GetFileServerMockRecorder struct {
	mock *MockDataKeeperService_GetFileServer
}

// NewMockDataKeeperService_GetFileServer creates a new mock instance.
func NewMockDataKeeperService_GetFileServer(ctrl *gomock.Controller) *MockDataKeeperService_GetFileServer {
	mock := &MockDataKeeperService_GetFileServer{ctrl: ctrl}
	mock.recorder = &MockDataKeeperService_GetFileServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataKeeperService_GetFileServer) EXPECT() *MockDataKeeperService_GetFileServerMockRecorder {
	return m.recorder
}

// Context mocks base method.
func (m *MockDataKeeperService_GetFileServer) Context() context.Context {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Context")
	ret0, _ := ret[0].(context.Context)
	return ret0
}

// Context indicates an expected call of Context.
func (mr *MockDataKeeperService_GetFileServerMockRecorder) Context() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Context", reflect.TypeOf((*MockDataKeeperService_GetFileServer)(nil).Context))
}

// RecvMsg mocks base method.
func (m *MockDataKeeperService_GetFileServer) RecvMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RecvMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// RecvMsg indicates an expected call of RecvMsg.
func (mr *MockDataKeeperService_GetFileServerMockRecorder) RecvMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RecvMsg", reflect.TypeOf((*MockDataKeeperService_GetFileServer)(nil).RecvMsg), arg0)
}

// Send mocks base method.
func (m *MockDataKeeperService_GetFileServer) Send(arg0 *FileChunk) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockDataKeeperService_GetFileServerMockRecorder) Send(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockDataKeeperService_GetFileServer)(nil).Send), arg0)
}

// SendHeader mocks base method.
func (m *MockDataKeeperService_GetFileServer) SendHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendHeader indicates an expected call of SendHeader.
func (mr *MockDataKeeperService_GetFileServerMockRecorder) SendHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendHeader", reflect.TypeOf((*MockDataKeeperService_GetFileServer)(nil).SendHeader), arg0)
}

// SendMsg mocks base method.
func (m *MockDataKeeperService_GetFileServer) SendMsg(arg0 interface{}) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendMsg", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendMsg indicates an expected call of SendMsg.
func (mr *MockDataKeeperService_GetFileServerMockRecorder) SendMsg(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendMsg", reflect.TypeOf((*MockDataKeeperService_GetFileServer)(nil).SendMsg), arg0)
}

// SetHeader mocks base method.
func (m *MockDataKeeperService_GetFileServer) SetHeader(arg0 metadata.MD) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetHeader", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetHeader indicates an expected call of SetHeader.
func (mr *MockDataKeeperService_GetFileServerMockRecorder) SetHeader(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetHeader", reflect.TypeOf((*MockDataKeeperService_GetFileServer)(nil).SetHeader), arg0)
}

// SetTrailer mocks base method.
func (m *MockDataKeeperService_GetFileServer) SetTrailer(arg0 metadata.MD) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetTrailer", arg0)
}

// SetTrailer indicates an expected call of SetTrailer.
func (mr *MockDataKeeperService_GetFileServerMockRecorder) SetTrailer(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetTrailer", reflect.TypeOf((*MockDataKeeperService_GetFileServer)(nil).SetTrailer), arg0)
}

// MockDataKeeperServiceClient is a mock of DataKeeperServiceClient interface.
type MockDataKeeperServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockDataKeeperServiceClientMockRecorder
}

// MockDataKeeperServiceClientMockRecorder is the mock recorder for MockDataKeeperServiceClient.
type MockDataKeeperServiceClientMockRecorder struct {
	mock *MockDataKeeperServiceClient
}

// NewMockDataKeeperServiceClient creates a new mock instance.
func NewMockDataKeeperServiceClient(ctrl *gomock.Controller) *MockDataKeeperServiceClient {
	mock := &MockDataKeeperServiceClient{ctrl: ctrl}
	mock.recorder = &MockDataKeeperServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataKeeperServiceClient) EXPECT() *MockDataKeeperServiceClientMockRecorder {
	return m.recorder
}

// DeleteData mocks base method.
func (m *MockDataKeeperServiceClient) DeleteData(ctx context.Context, in *DeleteDataRequest, opts ...grpc.CallOption) (*UploadStatus, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteData", varargs...)
	ret0, _ := ret[0].(*UploadStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteData indicates an expected call of DeleteData.
func (mr *MockDataKeeperServiceClientMockRecorder) DeleteData(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteData", reflect.TypeOf((*MockDataKeeperServiceClient)(nil).DeleteData), varargs...)
}

// DeleteFile mocks base method.
func (m *MockDataKeeperServiceClient) DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*UploadStatus, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "DeleteFile", varargs...)
	ret0, _ := ret[0].(*UploadStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockDataKeeperServiceClientMockRecorder) DeleteFile(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockDataKeeperServiceClient)(nil).DeleteFile), varargs...)
}

// GetDataList mocks base method.
func (m *MockDataKeeperServiceClient) GetDataList(ctx context.Context, in *ListDataRequest, opts ...grpc.CallOption) (*ListDataResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetDataList", varargs...)
	ret0, _ := ret[0].(*ListDataResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDataList indicates an expected call of GetDataList.
func (mr *MockDataKeeperServiceClientMockRecorder) GetDataList(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDataList", reflect.TypeOf((*MockDataKeeperServiceClient)(nil).GetDataList), varargs...)
}

// GetFile mocks base method.
func (m *MockDataKeeperServiceClient) GetFile(ctx context.Context, in *GetFileRequest, opts ...grpc.CallOption) (DataKeeperService_GetFileClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFile", varargs...)
	ret0, _ := ret[0].(DataKeeperService_GetFileClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFile indicates an expected call of GetFile.
func (mr *MockDataKeeperServiceClientMockRecorder) GetFile(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFile", reflect.TypeOf((*MockDataKeeperServiceClient)(nil).GetFile), varargs...)
}

// GetFileList mocks base method.
func (m *MockDataKeeperServiceClient) GetFileList(ctx context.Context, in *ListFileRequest, opts ...grpc.CallOption) (*ListFileResponse, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetFileList", varargs...)
	ret0, _ := ret[0].(*ListFileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileList indicates an expected call of GetFileList.
func (mr *MockDataKeeperServiceClientMockRecorder) GetFileList(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileList", reflect.TypeOf((*MockDataKeeperServiceClient)(nil).GetFileList), varargs...)
}

// SaveData mocks base method.
func (m *MockDataKeeperServiceClient) SaveData(ctx context.Context, in *SaveDataRequest, opts ...grpc.CallOption) (*UploadStatus, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SaveData", varargs...)
	ret0, _ := ret[0].(*UploadStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveData indicates an expected call of SaveData.
func (mr *MockDataKeeperServiceClientMockRecorder) SaveData(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveData", reflect.TypeOf((*MockDataKeeperServiceClient)(nil).SaveData), varargs...)
}

// UploadFile mocks base method.
func (m *MockDataKeeperServiceClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (DataKeeperService_UploadFileClient, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "UploadFile", varargs...)
	ret0, _ := ret[0].(DataKeeperService_UploadFileClient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockDataKeeperServiceClientMockRecorder) UploadFile(ctx interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockDataKeeperServiceClient)(nil).UploadFile), varargs...)
}

// MockDataKeeperServiceServer is a mock of DataKeeperServiceServer interface.
type MockDataKeeperServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockDataKeeperServiceServerMockRecorder
}

// MockDataKeeperServiceServerMockRecorder is the mock recorder for MockDataKeeperServiceServer.
type MockDataKeeperServiceServerMockRecorder struct {
	mock *MockDataKeeperServiceServer
}

// NewMockDataKeeperServiceServer creates a new mock instance.
func NewMockDataKeeperServiceServer(ctrl *gomock.Controller) *MockDataKeeperServiceServer {
	mock := &MockDataKeeperServiceServer{ctrl: ctrl}
	mock.recorder = &MockDataKeeperServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDataKeeperServiceServer) EXPECT() *MockDataKeeperServiceServerMockRecorder {
	return m.recorder
}

// DeleteData mocks base method.
func (m *MockDataKeeperServiceServer) DeleteData(ctx context.Context, in *DeleteDataRequest) (*UploadStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteData", ctx, in)
	ret0, _ := ret[0].(*UploadStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteData indicates an expected call of DeleteData.
func (mr *MockDataKeeperServiceServerMockRecorder) DeleteData(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteData", reflect.TypeOf((*MockDataKeeperServiceServer)(nil).DeleteData), ctx, in)
}

// DeleteFile mocks base method.
func (m *MockDataKeeperServiceServer) DeleteFile(ctx context.Context, in *DeleteFileRequest) (*UploadStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFile", ctx, in)
	ret0, _ := ret[0].(*UploadStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFile indicates an expected call of DeleteFile.
func (mr *MockDataKeeperServiceServerMockRecorder) DeleteFile(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFile", reflect.TypeOf((*MockDataKeeperServiceServer)(nil).DeleteFile), ctx, in)
}

// GetDataList mocks base method.
func (m *MockDataKeeperServiceServer) GetDataList(ctx context.Context, in *ListDataRequest) (*ListDataResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetDataList", ctx, in)
	ret0, _ := ret[0].(*ListDataResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetDataList indicates an expected call of GetDataList.
func (mr *MockDataKeeperServiceServerMockRecorder) GetDataList(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetDataList", reflect.TypeOf((*MockDataKeeperServiceServer)(nil).GetDataList), ctx, in)
}

// GetFile mocks base method.
func (m *MockDataKeeperServiceServer) GetFile(blob *GetFileRequest, server DataKeeperService_GetFileServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFile", blob, server)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetFile indicates an expected call of GetFile.
func (mr *MockDataKeeperServiceServerMockRecorder) GetFile(blob, server interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFile", reflect.TypeOf((*MockDataKeeperServiceServer)(nil).GetFile), blob, server)
}

// GetFileList mocks base method.
func (m *MockDataKeeperServiceServer) GetFileList(ctx context.Context, in *ListFileRequest) (*ListFileResponse, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFileList", ctx, in)
	ret0, _ := ret[0].(*ListFileResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFileList indicates an expected call of GetFileList.
func (mr *MockDataKeeperServiceServerMockRecorder) GetFileList(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFileList", reflect.TypeOf((*MockDataKeeperServiceServer)(nil).GetFileList), ctx, in)
}

// SaveData mocks base method.
func (m *MockDataKeeperServiceServer) SaveData(ctx context.Context, in *SaveDataRequest) (*UploadStatus, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveData", ctx, in)
	ret0, _ := ret[0].(*UploadStatus)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveData indicates an expected call of SaveData.
func (mr *MockDataKeeperServiceServerMockRecorder) SaveData(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveData", reflect.TypeOf((*MockDataKeeperServiceServer)(nil).SaveData), ctx, in)
}

// UploadFile mocks base method.
func (m *MockDataKeeperServiceServer) UploadFile(server DataKeeperService_UploadFileServer) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadFile", server)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadFile indicates an expected call of UploadFile.
func (mr *MockDataKeeperServiceServerMockRecorder) UploadFile(server interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadFile", reflect.TypeOf((*MockDataKeeperServiceServer)(nil).UploadFile), server)
}
