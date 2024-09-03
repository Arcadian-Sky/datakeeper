package router

import (
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	pbservice "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/service/v1"
	pbuser "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/user/v1"
	"github.com/sirupsen/logrus"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/Arcadian-Sky/datakkeeper/internal/server/repository"
	"github.com/Arcadian-Sky/datakkeeper/internal/server/router/interceptor"
	"github.com/Arcadian-Sky/datakkeeper/internal/server/router/jwtrule"
	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCServer struct {
	cfg         *settings.InitedFlags
	log         *logrus.Logger
	reposervice repository.FileRepository
	repouser    repository.UserRepository
	repodata    repository.DataRepository
	serv        *grpc.Server
	// tokenKey
	pbservice.UnimplementedDataKeeperServiceServer
	pbuser.UnimplementedUserServiceServer
}

// InitGRPCServer initializes a new gRPC server.
func InitGRPCServer(cf *settings.InitedFlags, lg *logrus.Logger, rs *repository.FileRepository, ru *repository.UserRepository, rd *repository.DataRepository) (*GRPCServer, error) {
	// creates a gRPC server
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.UnaryInterceptor(lg, cf.SecretKey)),
		grpc.StreamInterceptor(interceptor.StreamInterceptor(lg, cf.SecretKey)),
	)

	ob := &GRPCServer{
		cfg:         cf,
		log:         lg,
		repodata:    *rd,
		reposervice: *rs,
		repouser:    *ru,
		serv:        s,
	}
	// register the service
	pbservice.RegisterDataKeeperServiceServer(s, ob)
	pbuser.RegisterUserServiceServer(s, ob)

	return ob, nil
}

const (
	FormatErrUserRegisterUserFailed      = "failed to register user (Register): %s"
	FormatErrUserRegisterContainerFailed = "failed to register user (CreateContainer): %s"
	FormatErrUserRegisterTokenFailed     = "failed to register user (cant generate token): %s"
)

// Регистрация нового пользователя.
func (s *GRPCServer) Register(ctx context.Context, in *pbuser.RegisterRequest) (*pbuser.RegisterResponse, error) {
	if in.Login == `` {
		return nil, status.Errorf(codes.InvalidArgument, "user is not set")
	}
	r := ""
	// encPass := service.EncryptPass(in.Password)
	user := model.User{Login: in.Login, Password: in.Password}

	id, err := s.repouser.Register(ctx, &user)
	if err != nil {
		return nil, status.Errorf(getCode(err), fmt.Sprintf(FormatErrUserRegisterUserFailed, err.Error()))
	}
	str := fmt.Sprintf("user %s (userid: %d) was created\n", user.Login, id)
	r += str
	s.log.Info(s)

	user.ID = id
	_, err = s.reposervice.CreateContainer(ctx, &user)
	if err != nil {
		return nil, status.Errorf(getCode(err), fmt.Sprintf(FormatErrUserRegisterContainerFailed, err.Error()))
	}
	str = fmt.Sprintf("bucket container %s (userid: %d) was created\n", user.Bucket, id)
	r += str
	s.log.Info(str)

	// generate JWT
	userJWT, err := jwtrule.Generate(user.ID, s.cfg.SecretKey)
	if err != nil {
		return nil, status.Errorf(codes.Internal, fmt.Sprintf(FormatErrUserRegisterTokenFailed, err.Error()))
	}

	bSuccess := false
	if user.ID > 0 && user.Bucket != "" {
		bSuccess = true
	}

	return &pbuser.RegisterResponse{Success: bSuccess, Message: r, AuthToken: userJWT.Token}, nil
}

// Аутентификация пользователя.
func (s *GRPCServer) Authenticate(ctx context.Context, in *pbuser.AuthenticateRequest) (*pbuser.AuthenticateResponse, error) {
	if in.Login == `` || in.Password == `` {
		return nil, status.Error(codes.InvalidArgument, "invalid argument")
	}
	user := &model.User{
		Login:    in.Login,
		Password: in.Password,
	}
	mess := ""
	user, err := s.repouser.Auth(ctx, user)
	if err != nil {
		s.log.Println(err)
		return nil, status.Error(codes.Internal, "failed to auth user")
	}
	mess += fmt.Sprintf("authorized as userID: %v ", user.ID)

	// generate JWT
	userJWT, err := jwtrule.Generate(user.ID, s.cfg.SecretKey)
	if err != nil {
		e := fmt.Sprintf("cant generate token: %s", err.Error())
		s.log.Info(e)
		return nil, status.Errorf(codes.Internal, "cant generate token: "+err.Error())
	}

	mess += "token generated"

	bSuccess := false
	if user.ID > 0 && userJWT.Token != "" {
		bSuccess = true
	}

	s.log.Trace(mess)

	return &pbuser.AuthenticateResponse{
		Success:   bSuccess,
		AuthToken: userJWT.Token,
		Message:   mess,
	}, nil
}

// Start launch the server.
func (g *GRPCServer) Start() error {
	g.log.Info("g.cfg.Endpoint: ", g.cfg.Endpoint, "")
	// determines the server port
	listen, err := net.Listen("tcp", g.cfg.Endpoint)
	if err != nil {
		return err
	}
	g.log.Info("Сервер gRPC начал работу")
	// listen for gRPC requests
	return g.serv.Serve(listen)
}

// ShutDown graceful stops the server.
func (g *GRPCServer) ShutDown() error {
	g.serv.GracefulStop()
	return nil
}

func (s *GRPCServer) UploadFile(stream pbservice.DataKeeperService_UploadFileServer) error {
	// Obtain context from the stream
	ctx := stream.Context()
	uID := jwtrule.GetUserIDFromCTX(ctx)
	s.log.Trace("uID: ", uID)

	user := &model.User{
		ID: uID,
	}

	// Create a temporary file to store uploaded data
	tmpFile, err := os.CreateTemp("", "uploaded-*.tmp")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	defer os.Remove(tmpFile.Name())

	objectName := ""
	// Read chunks from the stream and write to the temp file
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to receive chunk: %w", err)
		}

		if _, err := tmpFile.Write(chunk.Data); err != nil {
			return fmt.Errorf("failed to write chunk to temp file: %w", err)
		}

		objectName = chunk.Filename
	}

	// Close the temp file
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Upload the file to MinIO
	file, err := os.Open(tmpFile.Name())
	if err != nil {
		return fmt.Errorf("failed to open temp file: %w", err)
	}
	defer file.Close()

	err = s.reposervice.UploadFile(ctx, user, objectName, file)
	if err != nil {
		return fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	//Update User
	user.LastUpdate = time.Now()
	_, err = s.repouser.SetLastUpdate(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to SetLastUpdate: %v", err)
	}
	// Send response to client
	return stream.SendAndClose(&pbservice.UploadStatus{
		Success: true,
		Message: "File uploaded successfully",
	})
}

func (s *GRPCServer) GetFileList(ctx context.Context, in *pbservice.ListFileRequest) (*pbservice.ListFileResponse, error) {
	uID := jwtrule.GetUserIDFromCTX(ctx)
	s.log.Trace("uID: ", uID)

	user := &model.User{
		ID: uID,
	}

	var data []model.FileItem
	data, err := s.reposervice.GetFileList(ctx, user)
	if err != nil {
		s.log.Println(err)
		return nil, status.Error(codes.Internal, "failed to get user files")
	}
	var resp []*pbservice.FileItem
	for _, it := range data {
		resp = append(resp, &pbservice.FileItem{
			Key:  it.Hash,
			Name: it.Name,
		})
	}

	return &pbservice.ListFileResponse{Fileitem: resp}, nil
}

func (s *GRPCServer) SaveData(ctx context.Context, in *pbservice.SaveDataRequest) (*pbservice.UploadStatus, error) {

	uID := jwtrule.GetUserIDFromCTX(ctx)
	user := &model.User{
		ID: uID,
	}
	s.log.Trace("uID: ", uID)

	data := model.Data{
		Type:     getType(in.Data.Type),
		UserID:   uID,
		Title:    in.Data.Title,
		Card:     in.Data.Card,
		Login:    in.Data.Login,
		Password: in.Data.Password,
	}

	fmt.Printf("data: %v\n", data)

	_, err := s.repodata.Save(ctx, &data)
	if err != nil {
		s.log.Println(err)
		return nil, status.Error(codes.Internal, "failed to get user files")
	}

	//Update User
	user.LastUpdate = time.Now()
	s.log.Println(user)

	_, err = s.repouser.SetLastUpdate(ctx, user)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to SetLastUpdate: "+err.Error())
	}

	return &pbservice.UploadStatus{Success: true, Message: "empty"}, nil
}

func (s *GRPCServer) GetDataList(ctx context.Context, in *pbservice.ListDataRequest) (*pbservice.ListDataResponse, error) {
	uID := jwtrule.GetUserIDFromCTX(ctx)
	s.log.Trace("uID: ", uID)
	user := &model.User{
		ID: uID,
	}
	data, err := s.repodata.GetList(ctx, user)
	if err != nil {
		return nil, status.Errorf(getCode(err), "failed to list pdata: "+err.Error())
	}

	var pdataPointers []*pbservice.Data
	for _, item := range data {
		pdataPointers = append(pdataPointers,
			&pbservice.Data{
				Id:       item.ID,
				Title:    item.Title,
				Type:     getPType(item.Type),
				Card:     item.Card,
				Login:    item.Login,
				Password: item.Password,
			})
	}

	return &pbservice.ListDataResponse{Data: pdataPointers}, nil

}
func (s *GRPCServer) DeleteFile(ctx context.Context, in *pbservice.DeleteFileRequest) (*pbservice.UploadStatus, error) {
	uID := jwtrule.GetUserIDFromCTX(ctx)
	s.log.Trace("uID: ", uID)
	user := model.User{
		ID: uID,
	}

	err := s.reposervice.DeleteFile(ctx, in.Filename, &user)
	if err != nil {
		e := fmt.Sprintf("failed to delete file: %v", err)
		s.log.Info(e)
		return nil, status.Errorf(getCode(err), "failed to delete file: "+err.Error())
	}

	return &pbservice.UploadStatus{Success: true, Message: "data was deleted"}, nil
}

// func (s *GRPCServer) getUserID(ctx context.Context) (int64, error) {
// 	md, ok := metadata.FromIncomingContext(ctx)
// 	if !ok {
// 		return 0, status.Errorf(codes.Internal, "failed to retrive metadata from ctx")
// 	}
// 	tokenRaw, ok := md["authorization"]
// 	if !ok || len(tokenRaw[0]) == 0 {
// 		return 0, model.ErrNoToken
// 	}
// 	tokenParsed, err := jwtrule.Validate(tokenRaw[0], s.cfg.SecretKey)
// 	if err != nil {
// 		return 0, model.ErrInvalidToken
// 	}

// 	return tokenParsed.Claims.UserID, nil
// }

// GetFile retrieves a file from MinIO and sends it as a stream of FileChunks
func (s *GRPCServer) GetFile(req *pbservice.GetFileRequest, stream pbservice.DataKeeperService_GetFileServer) error {
	uID := jwtrule.GetUserIDFromCTX(stream.Context())
	s.log.Trace("uID: ", uID)
	user := &model.User{
		ID: uID,
	}
	fileID := req.GetName()

	// Получаем файл из MinIO
	file, err := s.reposervice.GetFile(stream.Context(), fileID, user)
	if err != nil {
		return err
	}
	defer file.Close()

	buffer := make([]byte, 1024*1024) // 1 MB buffer size

	// Читаем файл по частям и отправляем через поток
	for {
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			return err
		}
		if n == 0 {
			break
		}

		// Отправляем часть файла в виде FileChunk
		chunk := &pbservice.FileChunk{
			Data:     buffer[:n],
			Filename: fileID,
		}
		if err := stream.Send(chunk); err != nil {
			return err
		}
	}

	return nil
}

// getCode returns gRCP error code based on error type
func getCode(e error) codes.Code {
	switch e {
	case model.ErrUserAlreadyExists, model.ErrPdataAlreatyEsists:
		return codes.AlreadyExists
	case model.ErrUserAuth, model.ErrInvalidToken, model.ErrNoToken:
		return codes.Unauthenticated
	case model.ErrPdataNotFound:
		return codes.NotFound
	default:
		return codes.Internal
	}

}

func getPType(stype string) pbservice.DataType {
	switch stype {
	case repository.DataTypeCARD:
		return pbservice.DataType_DATA_TYPE_TYPE_CREDIT_CARD
	case repository.DataTypeLOGPASS:
		return pbservice.DataType_DATA_TYPE_TYPE_LOGIN_PASSWORD
	}

	return pbservice.DataType_DATA_TYPE_UNSPECIFIED
}

func getType(stype pbservice.DataType) string {
	switch int(stype.Number()) {
	case int(pbservice.DataType_DATA_TYPE_TYPE_CREDIT_CARD):
		return repository.DataTypeCARD
	case int(pbservice.DataType_DATA_TYPE_TYPE_LOGIN_PASSWORD):
		return repository.DataTypeLOGPASS
	}

	return ""
}
