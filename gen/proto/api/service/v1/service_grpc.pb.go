// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: proto/api/service/v1/service.proto

package data

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DataKeeperService_SaveData_FullMethodName    = "/proto.api.service.v1.DataKeeperService/SaveData"
	DataKeeperService_GetDataList_FullMethodName = "/proto.api.service.v1.DataKeeperService/GetDataList"
	DataKeeperService_DeleteData_FullMethodName  = "/proto.api.service.v1.DataKeeperService/DeleteData"
	DataKeeperService_GetFileList_FullMethodName = "/proto.api.service.v1.DataKeeperService/GetFileList"
	DataKeeperService_UploadFile_FullMethodName  = "/proto.api.service.v1.DataKeeperService/UploadFile"
	DataKeeperService_GetFile_FullMethodName     = "/proto.api.service.v1.DataKeeperService/GetFile"
	DataKeeperService_DeleteFile_FullMethodName  = "/proto.api.service.v1.DataKeeperService/DeleteFile"
)

// DataKeeperServiceClient is the client API for DataKeeperService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Определение gRPC-сервиса для управления данными
type DataKeeperServiceClient interface {
	// Хранение новых данных на сервере (кроме файлов)
	SaveData(ctx context.Context, in *SaveDataRequest, opts ...grpc.CallOption) (*UploadStatus, error)
	GetDataList(ctx context.Context, in *ListDataRequest, opts ...grpc.CallOption) (*ListDataResponse, error)
	DeleteData(ctx context.Context, in *DeleteDataRequest, opts ...grpc.CallOption) (*UploadStatus, error)
	// Отправка файлов на сервер
	GetFileList(ctx context.Context, in *ListFileRequest, opts ...grpc.CallOption) (*ListFileResponse, error)
	UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[FileChunk, UploadStatus], error)
	GetFile(ctx context.Context, in *GetFileRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[FileChunk], error)
	DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*UploadStatus, error)
}

type dataKeeperServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataKeeperServiceClient(cc grpc.ClientConnInterface) DataKeeperServiceClient {
	return &dataKeeperServiceClient{cc}
}

func (c *dataKeeperServiceClient) SaveData(ctx context.Context, in *SaveDataRequest, opts ...grpc.CallOption) (*UploadStatus, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadStatus)
	err := c.cc.Invoke(ctx, DataKeeperService_SaveData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataKeeperServiceClient) GetDataList(ctx context.Context, in *ListDataRequest, opts ...grpc.CallOption) (*ListDataResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListDataResponse)
	err := c.cc.Invoke(ctx, DataKeeperService_GetDataList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataKeeperServiceClient) DeleteData(ctx context.Context, in *DeleteDataRequest, opts ...grpc.CallOption) (*UploadStatus, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadStatus)
	err := c.cc.Invoke(ctx, DataKeeperService_DeleteData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataKeeperServiceClient) GetFileList(ctx context.Context, in *ListFileRequest, opts ...grpc.CallOption) (*ListFileResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListFileResponse)
	err := c.cc.Invoke(ctx, DataKeeperService_GetFileList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataKeeperServiceClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[FileChunk, UploadStatus], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DataKeeperService_ServiceDesc.Streams[0], DataKeeperService_UploadFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[FileChunk, UploadStatus]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DataKeeperService_UploadFileClient = grpc.ClientStreamingClient[FileChunk, UploadStatus]

func (c *dataKeeperServiceClient) GetFile(ctx context.Context, in *GetFileRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[FileChunk], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DataKeeperService_ServiceDesc.Streams[1], DataKeeperService_GetFile_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[GetFileRequest, FileChunk]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DataKeeperService_GetFileClient = grpc.ServerStreamingClient[FileChunk]

func (c *dataKeeperServiceClient) DeleteFile(ctx context.Context, in *DeleteFileRequest, opts ...grpc.CallOption) (*UploadStatus, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UploadStatus)
	err := c.cc.Invoke(ctx, DataKeeperService_DeleteFile_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataKeeperServiceServer is the server API for DataKeeperService service.
// All implementations should embed UnimplementedDataKeeperServiceServer
// for forward compatibility.
//
// Определение gRPC-сервиса для управления данными
type DataKeeperServiceServer interface {
	// Хранение новых данных на сервере (кроме файлов)
	SaveData(context.Context, *SaveDataRequest) (*UploadStatus, error)
	GetDataList(context.Context, *ListDataRequest) (*ListDataResponse, error)
	DeleteData(context.Context, *DeleteDataRequest) (*UploadStatus, error)
	// Отправка файлов на сервер
	GetFileList(context.Context, *ListFileRequest) (*ListFileResponse, error)
	UploadFile(grpc.ClientStreamingServer[FileChunk, UploadStatus]) error
	GetFile(*GetFileRequest, grpc.ServerStreamingServer[FileChunk]) error
	DeleteFile(context.Context, *DeleteFileRequest) (*UploadStatus, error)
}

// UnimplementedDataKeeperServiceServer should be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDataKeeperServiceServer struct{}

func (UnimplementedDataKeeperServiceServer) SaveData(context.Context, *SaveDataRequest) (*UploadStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SaveData not implemented")
}
func (UnimplementedDataKeeperServiceServer) GetDataList(context.Context, *ListDataRequest) (*ListDataResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDataList not implemented")
}
func (UnimplementedDataKeeperServiceServer) DeleteData(context.Context, *DeleteDataRequest) (*UploadStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteData not implemented")
}
func (UnimplementedDataKeeperServiceServer) GetFileList(context.Context, *ListFileRequest) (*ListFileResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetFileList not implemented")
}
func (UnimplementedDataKeeperServiceServer) UploadFile(grpc.ClientStreamingServer[FileChunk, UploadStatus]) error {
	return status.Errorf(codes.Unimplemented, "method UploadFile not implemented")
}
func (UnimplementedDataKeeperServiceServer) GetFile(*GetFileRequest, grpc.ServerStreamingServer[FileChunk]) error {
	return status.Errorf(codes.Unimplemented, "method GetFile not implemented")
}
func (UnimplementedDataKeeperServiceServer) DeleteFile(context.Context, *DeleteFileRequest) (*UploadStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteFile not implemented")
}
func (UnimplementedDataKeeperServiceServer) testEmbeddedByValue() {}

// UnsafeDataKeeperServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataKeeperServiceServer will
// result in compilation errors.
type UnsafeDataKeeperServiceServer interface {
	mustEmbedUnimplementedDataKeeperServiceServer()
}

func RegisterDataKeeperServiceServer(s grpc.ServiceRegistrar, srv DataKeeperServiceServer) {
	// If the following call pancis, it indicates UnimplementedDataKeeperServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DataKeeperService_ServiceDesc, srv)
}

func _DataKeeperService_SaveData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataKeeperServiceServer).SaveData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataKeeperService_SaveData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataKeeperServiceServer).SaveData(ctx, req.(*SaveDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataKeeperService_GetDataList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataKeeperServiceServer).GetDataList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataKeeperService_GetDataList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataKeeperServiceServer).GetDataList(ctx, req.(*ListDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataKeeperService_DeleteData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataKeeperServiceServer).DeleteData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataKeeperService_DeleteData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataKeeperServiceServer).DeleteData(ctx, req.(*DeleteDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataKeeperService_GetFileList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataKeeperServiceServer).GetFileList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataKeeperService_GetFileList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataKeeperServiceServer).GetFileList(ctx, req.(*ListFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataKeeperService_UploadFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DataKeeperServiceServer).UploadFile(&grpc.GenericServerStream[FileChunk, UploadStatus]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DataKeeperService_UploadFileServer = grpc.ClientStreamingServer[FileChunk, UploadStatus]

func _DataKeeperService_GetFile_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetFileRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(DataKeeperServiceServer).GetFile(m, &grpc.GenericServerStream[GetFileRequest, FileChunk]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DataKeeperService_GetFileServer = grpc.ServerStreamingServer[FileChunk]

func _DataKeeperService_DeleteFile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteFileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataKeeperServiceServer).DeleteFile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataKeeperService_DeleteFile_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataKeeperServiceServer).DeleteFile(ctx, req.(*DeleteFileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DataKeeperService_ServiceDesc is the grpc.ServiceDesc for DataKeeperService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataKeeperService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.api.service.v1.DataKeeperService",
	HandlerType: (*DataKeeperServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveData",
			Handler:    _DataKeeperService_SaveData_Handler,
		},
		{
			MethodName: "GetDataList",
			Handler:    _DataKeeperService_GetDataList_Handler,
		},
		{
			MethodName: "DeleteData",
			Handler:    _DataKeeperService_DeleteData_Handler,
		},
		{
			MethodName: "GetFileList",
			Handler:    _DataKeeperService_GetFileList_Handler,
		},
		{
			MethodName: "DeleteFile",
			Handler:    _DataKeeperService_DeleteFile_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "UploadFile",
			Handler:       _DataKeeperService_UploadFile_Handler,
			ClientStreams: true,
		},
		{
			StreamName:    "GetFile",
			Handler:       _DataKeeperService_GetFile_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "proto/api/service/v1/service.proto",
}
