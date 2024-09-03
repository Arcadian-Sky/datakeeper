package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		grpc.WithStreamInterceptor(getStreamClientInterceptor(mstorage)),
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

// Регистрация нового пользователя.
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

// Аутентификация пользователя.
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

// Получение списка файлов
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
	for _, item := range res.Fileitem {
		data = append(data, model.FileItem{
			Hash: item.Key,
			Name: item.Name,
		})
	}

	return data, nil
}

// Удаление файла
func (gc *GRPCClient) DeleteFile(fileName string) error {
	if gc.Data == nil {
		return fmt.Errorf("GRPC client is not initialized")
	}

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pbsrv.DeleteFileRequest{
		Filename: fileName,
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
func (gc *GRPCClient) UploadFile(filePath string) error {
	fileName := filepath.Base(filePath)
	gc.log.Info("File name: ", fileName)
	// Открываем файл для чтения
	file, err := os.Open(filePath)
	if err != nil {
		gc.log.Trace("could not open file: ", err)
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	stream, err := gc.Data.UploadFile(context.Background())
	if err != nil {
		gc.log.Trace("error creating stream: ", err)
		return fmt.Errorf("error creating stream: %v", err)
	}

	// Read the file and send chunks
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			gc.log.Trace("Failed to read file: ", err)
		}

		if err := stream.Send(&pbsrv.FileChunk{
			Data:     buffer[:n],
			Filename: fileName,
		}); err != nil {
			gc.log.Trace("Failed to send chunk: ", err)
		}
	}
	// Close the stream and get the response
	status, err := stream.CloseAndRecv()
	if err != nil {
		gc.log.Trace("Failed to receive response: ", err)
	}

	gc.log.Info("Upload status:", status.Success, ", message: ", status.Message)

	return nil
}

func (gc *GRPCClient) GetFile(fileName string) error {

	stream, err := gc.Data.GetFile(context.Background(), &pbsrv.GetFileRequest{Name: fileName})
	if err != nil {
		gc.log.Fatalf("Ошибка при вызове GetFile: %v", err)
	}

	filePath := filepath.Join(gc.Storage.PfilesDir, fileName)
	// Создаём файл для записи полученных данных
	file, err := os.Create(filePath)
	if err != nil {
		gc.log.Trace("Не удалось создать файл: ", err)
		return err
	}
	defer file.Close()

	// Читаем поток данных и записываем в файл
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			gc.log.Trace("Ошибка при получении данных: ", err)
			return err
		}
		_, err = file.Write(chunk.Data)
		if err != nil {
			gc.log.Trace("Не удалось записать в файл: ", err)
			return err
		}
	}
	gc.log.Info("Файл успешно получен и сохранён:", filePath)
	return nil
}

func (gc *GRPCClient) SaveLoginPass(domain, login, pass string) error {
	if gc.Data == nil {
		return fmt.Errorf("GRPC client is not initialized")
	}
	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := gc.Data.SaveData(ctx, &pbsrv.SaveDataRequest{
		Data: &pbsrv.Data{
			Type:     pbsrv.DataType_DATA_TYPE_TYPE_LOGIN_PASSWORD,
			Password: pass,
			Login:    login,
			Title:    domain,
		},
	})

	if err != nil {
		gc.log.Debug("Error during save file : ", err)
		return err
	}
	gc.log.Trace(res)

	return nil
}

func (gc *GRPCClient) SaveCard(title, card string) error {

	if gc.Data == nil {
		return fmt.Errorf("GRPC client is not initialized")
	}

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := gc.Data.SaveData(ctx, &pbsrv.SaveDataRequest{
		Data: &pbsrv.Data{
			Type:  pbsrv.DataType_DATA_TYPE_TYPE_CREDIT_CARD,
			Card:  card,
			Title: title,
		},
	})

	if err != nil {
		gc.log.Debug("Error during save file : ", err)
		return err
	}
	gc.log.Trace(res)

	return nil
}

func (gc *GRPCClient) GetDataList() ([]model.Data, error) {
	var data []model.Data
	if gc.Data == nil {
		return data, fmt.Errorf("GRPC client is not initialized")
	}

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pbsrv.ListDataRequest{}
	// Отправляем запрос на сервер
	res, err := gc.Data.GetDataList(ctx, req)
	if err != nil {
		gc.log.Debug("Error during get list of files : ", err)
		return data, err
	}

	gc.log.Trace(res)
	for _, item := range res.Data {
		data = append(data, model.Data{
			ID:       item.Id,
			Title:    item.Title,
			Type:     item.Type.String(),
			Login:    item.Login,
			Card:     item.Card,
			Password: item.Password,
		})
	}

	return data, nil
}

func (gc *GRPCClient) Delete(id int64) error {
	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pbsrv.DeleteDataRequest{
		Dataid: id,
	}

	// Отправляем запрос на сервер
	res, err := gc.Data.DeleteData(ctx, req)
	if err != nil {
		gc.log.Debug("Error during get list of files : ", err)
		return err
	}

	gc.log.Trace(res)

	return nil
}
