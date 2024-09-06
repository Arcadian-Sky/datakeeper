package interceptor

import (
	"bytes"
	"context"
	"reflect"
	"testing"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/Arcadian-Sky/datakkeeper/internal/server/router/jwtrule"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func TestUnaryInterceptor_NoAuthSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := logrus.New()
	secretKey := "test-secret"

	interceptor := UnaryInterceptor(log, secretKey)
	// Создаем мокаем контекст с JWT токеном
	jwToken, err := jwtrule.Generate(123, secretKey)
	assert.NoError(t, err)

	ctx := jwtrule.SetUserIDToCTX(context.Background(), int(jwToken.Claims.UserID))
	info := &grpc.UnaryServerInfo{
		FullMethod: "/proto.api.user.v1.UserService/Register",
	}

	// Мокаем функцию обработчика
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "test response", nil
	}
	// Вызов интерцептора
	resp, err := interceptor(ctx, nil, info, handler)
	assert.NoError(t, err)

	assert.Equal(t, "test response", resp)
}

func TestUnaryInterceptor_AuthCheckSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := logrus.New()
	// log.SetLevel(logrus.TraceLevel)
	secretKey := "test-secret"

	jwToken, err := jwtrule.Generate(123, secretKey)
	assert.NoError(t, err)

	// Создаем метаданные и добавляем туда заголовок Authorization с типом Bearer
	md := metadata.New(map[string]string{"authorization": "bearer " + jwToken.Token})
	ctx := metadata.NewIncomingContext(context.Background(), md)

	interceptor := UnaryInterceptor(log, secretKey)

	info := &grpc.UnaryServerInfo{
		FullMethod: "/proto.api.service.v1.DataKeeperService/GetFile",
	}

	// Мокаем функцию обработчика
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return "test response", nil
	}
	// Вызов интерцептора
	resp, err := interceptor(ctx, nil, info, handler)
	assert.NoError(t, err)

	assert.Equal(t, "test response", resp)
}

func TestUnaryInterceptor_AuthFail(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	log := logrus.New()
	secretKey := "test-secret"

	interceptor := UnaryInterceptor(log, secretKey)

	// Мокаем контекст без аутентификации
	ctx := context.Background()
	info := &grpc.UnaryServerInfo{
		FullMethod: "/proto.api.service.v1.DataKeeperService/GetFile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, nil
	}

	_, err := interceptor(ctx, nil, info, handler)

	assert.Error(t, err)
	assert.Equal(t, codes.Unauthenticated, status.Code(err))
}

func Test_checkAuth(t *testing.T) {
	// Mock JWT token validation
	mockedToken := model.Jtoken{
		Claims: model.Claims{
			UserID: 123,
		},
	}

	tests := []struct {
		name string
		args struct {
			ctx          *context.Context
			log          *logrus.Logger
			secretKey    string
			method       string
			validateFunc func(tokenString string, key string) (model.Jtoken, error)
		}
		want    *model.Jtoken
		wantErr bool
	}{
		{
			name: "Success",
			args: struct {
				ctx          *context.Context
				log          *logrus.Logger
				secretKey    string
				method       string
				validateFunc func(tokenString string, key string) (model.Jtoken, error)
			}{
				ctx:       contextWithToken("bearer valid-token"),
				log:       logrus.New(),
				secretKey: "test-secret",
				method:    "/proto.api.service.v1.DataKeeperService/GetFile",
				validateFunc: func(tokenString, secretKey string) (model.Jtoken, error) {
					if tokenString == "valid-token" {
						return mockedToken, nil
					}
					return model.Jtoken{}, status.Error(codes.Unauthenticated, "invalid token")
				},
			},
			want:    &mockedToken,
			wantErr: false,
		},
		{
			name: "NoToken",
			args: struct {
				ctx          *context.Context
				log          *logrus.Logger
				secretKey    string
				method       string
				validateFunc func(tokenString string, key string) (model.Jtoken, error)
			}{
				ctx:       contextWithToken(""),
				log:       logrus.New(),
				secretKey: "test-secret",
				method:    "/proto.api.service.v1.DataKeeperService/GetFile",
				validateFunc: func(tokenString, key string) (model.Jtoken, error) {
					return model.Jtoken{}, status.Error(codes.Unauthenticated, "no token")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "InvalidToken",
			args: struct {
				ctx          *context.Context
				log          *logrus.Logger
				secretKey    string
				method       string
				validateFunc func(tokenString string, key string) (model.Jtoken, error)
			}{
				ctx:       contextWithToken("bearer invalid-token"),
				log:       logrus.New(),
				secretKey: "test-secret",
				method:    "/proto.api.service.v1.DataKeeperService/GetFile",
				validateFunc: func(tokenString, secretKey string) (model.Jtoken, error) {
					return model.Jtoken{}, status.Error(codes.Unauthenticated, "invalid token")
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "SkipMethod",
			args: struct {
				ctx          *context.Context
				log          *logrus.Logger
				secretKey    string
				method       string
				validateFunc func(tokenString string, key string) (model.Jtoken, error)
			}{
				ctx:          contextWithToken("bearer some-token"),
				log:          logrus.New(),
				secretKey:    "test-secret",
				method:       "/proto.api.user.v1.UserService/Register",
				validateFunc: nil, // Use default validation function
			},
			want:    nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := checkAuth(tt.args.ctx, tt.args.log, tt.args.secretKey, tt.args.method, tt.args.validateFunc)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkAuth() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("checkAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

// Helper function to create a context with metadata
func contextWithToken(token string) *context.Context {
	md := metadata.New(map[string]string{"authorization": token})
	ctx := metadata.NewIncomingContext(context.Background(), md)
	return &ctx
}

func Test_preProcess(t *testing.T) {
	type args struct {
		ctx  context.Context
		info string
	}
	log := logrus.New()
	log.SetLevel(logrus.TraceLevel)

	tests := []struct {
		name string
		args args
	}{
		{
			name: "with userID in context",
			args: args{
				ctx:  context.WithValue(context.Background(), jwtrule.CtxKeyUserID, 123),
				info: "/proto.api.service.v1.DataKeeperService/GetFile",
			},
		},
		{
			name: "without userID in context",
			args: args{
				ctx:  context.Background(),
				info: "/proto.api.service.v1.DataKeeperService/GetFile",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var logOutput bytes.Buffer
			log.SetOutput(&logOutput)

			preProcess(tt.args.ctx, tt.args.info, log, "")

			// Check log output
			assert.Contains(t, logOutput.String(), "--> interceptor: before executing:")
			if tt.args.ctx.Value("userID") != nil {
				assert.Contains(t, logOutput.String(), "--> pre-processing for user ID: ")
			}
		})
	}
}

func Test_postProcess(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.TraceLevel)
	type args struct {
		ctx  context.Context
		info string
		err  error
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "with error and userID in context",
			args: args{
				ctx:  context.WithValue(context.Background(), jwtrule.CtxKeyUserID, 123),
				info: "/proto.api.service.v1.DataKeeperService/GetFile",
				err:  assert.AnError,
			},
		},
		{
			name: "with error and without userID in context",
			args: args{
				ctx:  context.Background(),
				info: "/proto.api.service.v1.DataKeeperService/GetFile",
				err:  assert.AnError,
			},
		},
		{
			name: "without error and with userID in context",
			args: args{
				ctx:  context.WithValue(context.Background(), jwtrule.CtxKeyUserID, 123),
				info: "/proto.api.service.v1.DataKeeperService/GetFile",
				err:  nil,
			},
		},
		{
			name: "without error and without userID in context",
			args: args{
				ctx:  context.Background(),
				info: "/proto.api.service.v1.DataKeeperService/GetFile",
				err:  nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var logOutput bytes.Buffer
			log.SetOutput(&logOutput)

			postProcess(tt.args.ctx, tt.args.info, log, tt.args.err)

			// Check log output
			assert.Contains(t, logOutput.String(), "--> interceptor: after executing:")
			if tt.args.err != nil {
				assert.Contains(t, logOutput.String(), "--> interceptor: error occurred: ")
			}
			if tt.args.ctx.Value("userID") != nil {
				assert.Contains(t, logOutput.String(), "--> interceptor: post-processing for user ID: ")
			}
		})
	}
}

// func Test_serverStreamWithContext_Context(t *testing.T) {
// 	// Define the test context and create an instance of serverStreamWithContext
// 	expectedCtx := context.WithValue(context.Background(), jwtrule.CtxKeyUserID, "testValue")

// 	// Mocking ServerStream
// 	mockServerStream := &mockServerStream{}

// 	// Create instance of serverStreamWithContext with the test context
// 	s := &serverStreamWithContext{
// 		ServerStream: mockServerStream,
// 		ctx:          expectedCtx,
// 	}

// 	// Call the Context method
// 	actualCtx := s.Context()

// 	// Check if the returned context is the same as the expected context
// 	assert.Equal(t, expectedCtx, actualCtx, "Expected context does not match the actual context")

// 	// Optional: Verify that the value set in the context is present
// 	value := actualCtx.Value("testKey")
// 	assert.Equal(t, "testValue", value, "Context value does not match")
// }

// MockServerStream is a mock implementation of grpc.ServerStream for testing
// type mockServerStream struct {
// 	grpc.ServerStream
// }

// serverStreamWithContext - структура для тестирования

// NewServerStreamWithContext - функция для создания экземпляра serverStreamWithContext
func NewServerTestStreamWithContext(ctx context.Context) *serverStreamWithContext {
	return &serverStreamWithContext{ctx: ctx}
}

// TestContext проверяет, что метод Context возвращает правильный контекст.
func TestContext(t *testing.T) {
	// Создаем новый контекст
	expectedCtx := context.WithValue(context.Background(), "key", "value")

	// Создаем экземпляр serverStreamWithContext с этим контекстом
	stream := NewServerTestStreamWithContext(expectedCtx)

	// Получаем контекст через метод Context
	actualCtx := stream.Context()

	// Проверяем, что возвращаемый контекст совпадает с ожидаемым
	assert.Equal(t, expectedCtx, actualCtx)

	// Дополнительно проверяем значение контекста, если это необходимо
	assert.Equal(t, "value", actualCtx.Value("key"))
}

// func TestStreamInterceptor_AuthorizationFailure(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	mockLogger := logrus.New()
// 	mockLogger.SetLevel(logrus.TraceLevel)

// 	// Интерцептор, который мы тестируем
// 	interceptor := StreamInterceptor(mockLogger, "test-secret-key")

// 	// Создаем mock серверный поток
// 	mockServerStream := pbservice.NewMockDataKeeperServiceServer(ctrl)
// 	mockServerStream.EXPECT().Context().Return(context.Background())

// 	// Настроим mock для checkAuth, чтобы он возвращал ошибку
// 	mockCheckAuth := func(ctx *context.Context, log *logrus.Logger, secretKey string, fullMethod string, headers interface{}) (*jwtrule.Token, error) {
// 		return nil, status.Errorf(codes.Unauthenticated, "auth error")
// 	}

// 	// Выполним тест
// 	err := interceptor(
// 		nil,              // srv
// 		mockServerStream, // ss
// 		&grpc.StreamServerInfo{FullMethod: "/test/method"},               // info
// 		func(srv interface{}, ss grpc.ServerStream) error { return nil }, // handler
// 	)

// 	// Проверяем результаты
// 	assert.Error(t, err)
// 	assert.Equal(t, codes.Unauthenticated, status.Code(err))
// }
