package router

import (
	"context"
	"fmt"
	"net"

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
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var (
	newerVersionDetected   = "Current / newer version found in database. Please synchronize you app to get the most actual data."
	failedToSaveNewVersion = "Failed to save new record. Please try again."
	failedDBQuery          = "failed to obtain database data"
)

type GRPCServer struct {
	cfg         *settings.InitedFlags
	log         *logrus.Logger
	reposervice repository.FileRepository
	repouser    repository.UserRepository
	serv        *grpc.Server
	// tokenKey
	pbservice.UnimplementedDataKeeperServiceServer
	pbuser.UnimplementedUserServiceServer
}

// InitGRPCServer initializes a new gRPC server.
func InitGRPCServer(cf *settings.InitedFlags, lg *logrus.Logger, rs *repository.FileRepository, ru *repository.UserRepository) (*GRPCServer, error) {
	// creates a gRPC server
	s := grpc.NewServer(
		grpc.UnaryInterceptor(interceptor.AuthCheckGRPC(lg, cf.SecretKey)),
	)

	ob := &GRPCServer{
		cfg:         cf,
		log:         lg,
		reposervice: *rs,
		repouser:    *ru,
		serv:        s,
	}
	// register the service
	pbservice.RegisterDataKeeperServiceServer(s, ob)
	pbuser.RegisterUserServiceServer(s, ob)

	return ob, nil
}

// Регистрация нового пользователя.
func (s *GRPCServer) Register(ctx context.Context, in *pbuser.RegisterRequest) (*pbuser.RegisterResponse, error) {
	if in.Login == `` {
		return nil, status.Errorf(codes.InvalidArgument, "user is not set")
	}
	r := ""
	// encPass := service.EncryptPass(in.Password)
	user := model.User{Login: in.Login, Password: in.Password}

	id, err := s.repouser.Register(ctx, user)
	if err != nil {
		e := fmt.Sprintf("failed to register user (Register): %s", err.Error())
		return nil, status.Errorf(getCode(err), e)
	}
	str := fmt.Sprintf("user %s (userid: %d) was created\n", user.Login, id)
	r += str
	s.log.Info(s)

	user.ID = id
	_, err = s.reposervice.CreateContainer(ctx, &user)
	if err != nil {
		e := fmt.Sprintf("failed to register user (CreateContainer): %s", err.Error())
		return nil, status.Errorf(getCode(err), e)
	}
	str = fmt.Sprintf("bucket container %s (userid: %d) was created\n", user.Bucket, id)
	r += str
	s.log.Info(str)

	// generate JWT
	userJWT, err := jwtrule.Generate(user.ID, s.cfg.SecretKey)
	if err != nil {
		e := fmt.Sprintf("cant generate token: %s", err.Error())
		return nil, status.Errorf(codes.Internal, e)
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
	user := model.User{
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
		return nil, status.Errorf(codes.Internal, e)
	}

	mess += fmt.Sprintf("token generated")

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

func (s *GRPCServer) AddData(ctx context.Context, in *pbservice.AddDataRequest) (*pbservice.AddDataResponse, error) {
	return nil, nil
}
func (s *GRPCServer) UploadFile(in grpc.ClientStreamingServer[pbservice.FileChunk, pbservice.UploadStatus]) error {
	// for {
	// 	req, err := stream.Recv()
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		return status.Errorf(codes.Internal, "failed to receive data: %v", err)
	// 	}

	// 	// Get user information (you would need to have a way to retrieve the user)
	// 	user := model.User{
	// 		ID:     req.GetUserId(),
	// 		Login:  req.GetUserLogin(),
	// 		Bucket: req.GetUserBucket(),
	// 	}

	// 	// Prepare the data to be saved
	// 	data := model.Data{
	// 		Name:        req.GetName(),
	// 		Type:        req.GetType(),
	// 		KeyHash:     req.GetKeyHash(),
	// 		PrivateData: req.GetPrivateData(),
	// 	}

	// 	// Save the data using the FileRepo's Save method
	// 	dataID, err := s.reposervice.Save(context.Background(), user, data)
	// 	if err != nil {
	// 		return status.Errorf(codes.Internal, "failed to save data: %v", err)
	// 	}

	// 	// Send response back to the client
	// 	resp := &pbservice.AddDataResponse{
	// 		DataId: dataID,
	// 	}
	// 	if err := stream.Send(resp); err != nil {
	// 		return status.Errorf(codes.Internal, "failed to send response: %v", err)
	// 	}
	// }

	return nil
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

	return &pbservice.ListFileResponse{FileItem: resp}, nil
}

// func (s *GRPCServer) UpdateData(grpc.BidiStreamingServer[pbservice.UpdateDataRequest, pbservice.UpdateDataResponse]) error {
// 	return nil
// }

// func (s *GRPCServer) AddData(ctx context.Context, in *pb.AddDataRequest) (*pb.AddDataResponse, error) {
// 	if in.GetData() == nil {
// 		return nil, status.Errorf(codes.InvalidArgument, "data is nil")
// 	}

// 	data := model.Data{
// 		Name:        in.Data.Name,
// 		Type:        getType(in.Data.Type),
// 		KeyHash:     base64.StdEncoding.EncodeToString(in.Data.KeyHash),
// 		PrivateData: base64.StdEncoding.EncodeToString(in.Data.Data)}

//		userID, err := s.getUserID(ctx)
//		if err != nil {
//			e := fmt.Sprintf("failed to get userid: %s", err.Error())
//			return nil, status.Errorf(getCode(err), e)
//		}
//		_, err = s.db.PrivateAdd(ctx, userID, data)
//		if err != nil {
//			e := fmt.Sprintf("failed to add data to database: %s", err.Error())
//			return nil, status.Errorf(getCode(err), e)
//		}
//		r := fmt.Sprintf("data %s was added sucsessfully", data.Name)
//		return &pb.AddDataResponse{Response: r}, nil
//	}
func (s *GRPCServer) GetData(grpc.BidiStreamingServer[pbservice.GetDataRequest, pbservice.GetDataResponse]) error {
	return nil
}

// func (s *GRPCServer) GetData(ctx context.Context, in *pb.GetDataRequest) (*pb.GetDataResponse, error) {
// 	// if in.DataID == 0 {
// 	// 	return nil, status.Errorf(codes.InvalidArgument, "pname is not set")
// 	// }
// 	// userID, err := s.getUserID(ctx)
// 	// if err != nil {
// 	// 	e := fmt.Sprintf("failed to get userid: %s", err.Error())
// 	// 	return nil, status.Errorf(getCode(err), e)
// 	// }

// 	// data, err := s.repo.Get(ctx, userID, in.DataID)
// 	// if err != nil {
// 	// 	e := fmt.Sprintf("failed to get pdata: %s", err.Error())
// 	// 	return nil, status.Errorf(getCode(err), e)
// 	// }

// 	// keyHash, err := base64.StdEncoding.DecodeString(data.KeyHash)
// 	// if err != nil {
// 	// 	e := fmt.Sprintf("cant decode khash_base64 to byte array: %s", err.Error())
// 	// 	return nil, status.Error(codes.Internal, e)
// 	// }
// 	// privateData, err := base64.StdEncoding.DecodeString(data.PrivateData)
// 	// if err != nil {
// 	// 	e := fmt.Sprintf("cant decode data_base64 to byte array: %s", err.Error())
// 	// 	return nil, status.Error(codes.Internal, e)
// 	// }
// 	// pbData := pb.Data{
// 	// 	ID: in.DataID,
// 	// 	Name:        data.Name,
// 	// 	Type:        getTypeCode(data.Type),
// 	// 	KeyHash:     keyHash,
// 	// 	PrivateData: []byte(privateData),
// 	// }
// 	// &pbData
// 	// pbData := pb.Data{}

// 	return &pb.GetDataResponse{Data: &pb.Data{}}, nil
// }

// func (s *GRPCServer) UpdateData(ctx context.Context, in *pb.UpdateDataRequest) (*pb.UpdateDataResponse, error) {
// 	if in.Data == nil {
// 		return nil, status.Errorf(codes.InvalidArgument, "pname is not set")
// 	}
// 	userID, err := s.getUserID(ctx)
// 	if err != nil {
// 		e := fmt.Sprintf("failed to get userid: %s", err.Error())
// 		return nil, status.Errorf(getCode(err), e)
// 	}
// 	data := model.Data{
// 		ID:          in.Data.ID,
// 		Name:        in.Data.Name,
// 		Type:        getType(in.Data.Type),
// 		KeyHash:     base64.StdEncoding.EncodeToString(in.Data.KeyHash),
// 		PrivateData: base64.StdEncoding.EncodeToString(in.Data.Data)}

// 	err = s.repo.Update(ctx, userID, data)
// 	if err != nil {
// 		e := fmt.Sprintf("failed to update pdata: %s", err.Error())
// 		return nil, status.Errorf(getCode(err), e)
// 	}
// 	r := fmt.Sprintf("data %s was updated sucsessfully", data.Name)
// 	return &pb.UpdateDataResponse{Response: r}, nil
// }

func (s *GRPCServer) ListData(ctx context.Context, in *pbservice.ListDataRequest) (*pbservice.ListDataResponse, error) {
	// userID, err := s.getUserID(ctx)
	// if err != nil {
	// 	e := fmt.Sprintf("failed to get userid: %s", err.Error())
	// 	return nil, status.Errorf(getCode(err), e)
	// }

	// data, err := s.repo.GetList(ctx, userID, in.Type)
	// if err != nil {
	// 	e := fmt.Sprintf("failed to list pdata: %s", err.Error())
	// 	return nil, status.Errorf(getCode(err), e)
	// }

	// var pdataPointers []*pb.DataItem
	// for _, p := range data {
	// 	pdataPointers = append(pdataPointers,
	// 		&pb.DataItem{Name: p.Name, ID: p.ID})
	// }

	// return &pb.ListDataResponse{DataItem: pdataPointers}, nil
	return &pbservice.ListDataResponse{DataItem: nil}, nil

}
func (s *GRPCServer) DeleteFile(ctx context.Context, in *pbservice.DeleteFileRequest) (*pbservice.UploadStatus, error) {
	uID := jwtrule.GetUserIDFromCTX(ctx)
	s.log.Trace("uID: ", uID)
	user := model.User{
		ID: uID,
	}

	err := s.reposervice.DeleteFile(ctx, in.FileName, &user)
	if err != nil {
		e := fmt.Sprintf("failed to delete file: %v", err)
		s.log.Info(e)
		return nil, status.Errorf(getCode(err), e)
	}

	return &pbservice.UploadStatus{Success: true, Message: "data was deleted"}, nil
}

func (s *GRPCServer) getUserID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, status.Errorf(codes.Internal, "failed to retrive metadata from ctx")
	}
	tokenRaw, ok := md["authorization"]
	if !ok || len(tokenRaw[0]) == 0 {
		return 0, model.ErrNoToken
	}
	tokenParsed, err := jwtrule.Validate(tokenRaw[0], s.cfg.SecretKey)
	if err != nil {
		return 0, model.ErrInvalidToken
	}

	return tokenParsed.Claims.UserID, nil

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

func getType(_ pbservice.DataType) string {
	return ""
}

func getTypeCode(_ string) pbservice.DataType {
	return pbservice.DataType_DATA_TYPE_UNSPECIFIED
}
