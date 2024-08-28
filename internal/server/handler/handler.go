package handler

import (
	"context"
	"time"

	pb "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/user/v1"
	app "github.com/Arcadian-Sky/datakkeeper/internal/app/server"
	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/Arcadian-Sky/datakkeeper/internal/server/repository"
	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Handler struct {
	repo *repository.UserRepository
	log  *logrus.Logger
	cfg  *settings.InitedFlags
	// pb   pb.UnimplementedUserServiceServer
}

// NewHandler создает экземпляр Handler
// storage repository.Repository,
func NewHandler(ap *app.App) {
	s := grpc.NewServer()
	pb.RegisterUserServiceServer(s, &Handler{
		repo: ap.GetUserRepo(),
		log:  ap.Logger,
		cfg:  ap.Flags,
		// pb:   pb.UnimplementedUserServiceServer,
	})
}

// Authenticate выполняет аутентификацию пользователя.
func (s *Handler) Authenticate(ctx context.Context, req *pb.AuthenticateRequest) (*pb.AuthenticateResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Invalid request: %v", err)
	}

	var hashedPassword string
	var userID int64

	// req.Email
	// req.Password

	// Запрос на получение хеша пароля и ID пользователя
	// err := s.db.QueryRowContext(ctx, "SELECT id, password FROM users WHERE email = $1", req.Email).Scan(&userID, &hashedPassword)
	// if err != nil {
	// 	if err == sql.ErrNoRows {
	// 		return &pb.AuthenticateResponse{
	// 			Success: false,
	// 			Message: "Invalid email or password",
	// 		}, nil
	// 	}
	// 	return nil, err
	// }

	// Сравнение пароля
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		return &pb.AuthenticateResponse{
			Success: false,
			Message: "Invalid email or password",
		}, nil
	}

	// Создание JWT
	tokenString, err := s.BuildJWTString(userID)
	if err != nil {
		return nil, err
	}

	return &pb.AuthenticateResponse{
		Success:   true,
		AuthToken: tokenString,
		Message:   "Authentication successful",
	}, nil
}
func (s *Handler) Register(context.Context, *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	return &pb.RegisterResponse{}, nil
}

const TokenExp = time.Hour * 3

// BuildJWTString создаёт токен и возвращает его в виде строки.
func (h *Handler) BuildJWTString(userID int64) (string, error) {
	// создаём новый токен с алгоритмом подписи HS256 и утверждениями — Claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, model.Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			// когда создан токен
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})

	// создаём строку токена
	tokenString, err := token.SignedString([]byte(h.cfg.SecretKey))
	if err != nil {
		return "", err
	}

	// возвращаем строку токена
	return tokenString, nil
}

// // StreamInterceptor - интерцептор для обработки стриминговых вызовов
// func StreamInterceptor(
// 	srv interface{},
// 	ss grpc.ServerStream,
// 	info *grpc.StreamServerInfo,
// 	handler grpc.StreamHandler,
// ) error {
// 	// Логирование входящих запросов
// 	log.Printf("Received stream for method: %s", info.FullMethod)

// 	// Вызов основного обработчика
// 	err := handler(srv, ss)

// 	// Логирование завершения
// 	if err != nil {
// 		st, _ := status.FromError(err)
// 		log.Printf("Stream %s returned error: %s", info.FullMethod, st.Message())
// 	} else {
// 		log.Printf("Stream %s returned successfully", info.FullMethod)
// 	}

// 	return err
// }

// // Пример сервиса
// type server struct{}

// func (s *server) SayHelloStream(req *HelloRequest, stream Greeter_SayHelloStreamServer) error {
// 	for {
// 		if err := stream.Send(&HelloResponse{Message: "Hello " + req.Name}); err != nil {
// 			return err
// 		}
// 	}
// }
