package client

import (
	"context"

	pbsrv "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/service/v1"
	pb "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/user/v1"
	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GRPCClientInterface interface {
	Register(login, password string) error
	Authenticate(login, password string) error

	GetDataList() ([]model.Data, error)
	SaveLoginPass(domain, login, pass string) error
	SaveCard(title, card string) error
	Delete(id int64) error

	GetFileList() ([]model.FileItem, error)
	DeleteFile(fileName string) error
	UploadFile(filePath string) error
	GetFile(fileName string) error
}

type GRPCClient struct {
	log     *logrus.Logger
	User    pb.UserServiceClient
	Data    pbsrv.DataKeeperServiceClient
	Storage *MemStorage
}

// NewGclient initializes new Gclient
func NewGclient(clientConfig settings.ClientConfig, mstorage *MemStorage, lg *logrus.Logger) (GRPCClientInterface, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	var err error

	conn, err = grpc.NewClient(
		clientConfig.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(getUnaryClientInterceptor(mstorage)),
		grpc.WithStreamInterceptor(getStreamClientInterceptor(mstorage)),
	)
	if err != nil {
		lg.Debug("failed to connect to server: ", err)
	}

	return &GRPCClient{
		Storage: mstorage,
		log:     lg,
		User:    pb.NewUserServiceClient(conn),
		Data:    pbsrv.NewDataKeeperServiceClient(conn),
	}, conn
}

var (
	// we don't need to check the token for these methods.
	SkipCheckMethods = map[string]struct{}{
		"/proto.api.user.v1.UserService/Register":     {},
		"/proto.api.user.v1.UserService/Authenticate": {},
	}
)

func getUnaryClientInterceptor(mstorage *MemStorage) grpc.UnaryClientInterceptor {
	return func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker,
		opts ...grpc.CallOption,
	) error {
		_, ok := SkipCheckMethods[method]
		if !ok {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+mstorage.Token)
		}

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

func getStreamClientInterceptor(mstorage *MemStorage) grpc.StreamClientInterceptor {
	return func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		streamer grpc.Streamer,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		_, ok := SkipCheckMethods[method]
		if !ok {
			ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+mstorage.Token)
		}

		// Call the streamer to create the ClientStream
		clientStream, err := streamer(ctx, desc, cc, method, opts...)
		if err != nil {
			return nil, err
		}

		// Return the created ClientStream
		return clientStream, nil
	}
}
