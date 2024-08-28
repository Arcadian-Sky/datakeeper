package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	pbsrv "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/service/v1"
	pb "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/user/v1"
	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type GRPCClient struct {
	log     *logrus.Logger
	User    pb.UserServiceClient
	Data    pbsrv.DataKeeperServiceClient
	Storage *MemStorage
}

// NewGclient initializes new Gclient
func NewGclient(clientConfig settings.ClientConfig, mstorage *MemStorage, lg *logrus.Logger) (GRPCClient, *grpc.ClientConn) {
	var conn *grpc.ClientConn
	var err error

	conn, err = grpc.NewClient(
		clientConfig.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(getUnaryClientInterceptor(mstorage)),
	)
	if err != nil {
		lg.Debug("failed to connect to server: ", err)
	}

	return GRPCClient{
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

// // Регистрация нового пользователя.
// Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
// // Аутентификация пользователя.
// Authenticate(ctx context.Context, in *AuthenticateRequest, opts ...grpc.CallOption) (*AuthenticateResponse, error)
func (gc *GRPCClient) Register(login, password string) error {

	if gc.User == nil {
		return fmt.Errorf("GRPC client is not initialized")
	}
	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Формируем запрос на регистрацию
	req := &pb.RegisterRequest{
		Login:    login,
		Password: password,
	}

	// Отправляем запрос на сервер
	res, err := gc.User.Register(ctx, req)
	if err != nil {
		gc.log.Debug("Error during registration:", err)
		return err
	}
	gc.Storage.SetToken(res.AuthToken)

	// Обрабатываем ответ сервера
	if res.Success {
		gc.log.Info("Registration successful, auth token:", res.AuthToken)
	} else {
		gc.log.Info("Registration failed:", res.Message)
	}

	return nil
}

func (gc *GRPCClient) Authenticate(login, password string) error {
	if gc.User == nil {
		return fmt.Errorf("GRPC client is not initialized")
	}

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Формируем запрос на регистрацию
	req := &pb.AuthenticateRequest{
		Login:    login,
		Password: password,
	}

	// Отправляем запрос на сервер
	res, err := gc.User.Authenticate(ctx, req)
	if err != nil {
		gc.log.Debug("Error during authrntificate: ", err)
		return err
	}
	gc.Storage.SetToken(res.AuthToken)

	// Обрабатываем ответ сервера
	if res.Success {
		gc.log.Info("Authrntificate successful, auth login: ", gc.Storage.Login)
		gc.log.Info("Authrntificate successful, auth token: ", res.AuthToken)
	} else {
		gc.log.Info("Authrntificate failed: ", res.Message)
		return errors.New("Authrntificate failed: " + res.Message)
	}

	return nil
}

func (gc *GRPCClient) GetFileList() ([]model.FileItem, error) {
	var data []model.FileItem
	if gc.Data == nil {
		return data, fmt.Errorf("GRPC client is not initialized")
	}

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pbsrv.ListFileRequest{}
	// Отправляем запрос на сервер
	res, err := gc.Data.GetFileList(ctx, req)
	if err != nil {
		gc.log.Debug("Error during get list of files : ", err)
		return data, err
	}
	gc.log.Trace(res)
	for _, item := range res.FileItem {
		data = append(data, model.FileItem{
			Hash: item.Key,
			Name: item.Name,
		})
	}

	return data, nil
}

func (gc *GRPCClient) DeleteFile(fileName string) error {
	if gc.Data == nil {
		return fmt.Errorf("GRPC client is not initialized")
	}

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pbsrv.DeleteFileRequest{
		FileName: fileName,
	}
	// Отправляем запрос на сервер
	res, err := gc.Data.DeleteFile(ctx, req)
	if err != nil {
		gc.log.Debug("Error during delete file : ", err)
		return err
	}
	gc.log.Trace(res)

	return nil
}

// Отправка файлов на сервер
func (gc *GRPCClient) UploadFile(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[pbsrv.FileChunk, pbsrv.UploadStatus], error) {
	return nil, nil
}
