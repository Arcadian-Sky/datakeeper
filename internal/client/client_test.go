package client

import (
	"context"
	"testing"

	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

func TestNewGclient_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mock storage and logger
	mockStorage := &MemStorage{}
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Mock the server configuration
	clientConfig := settings.ClientConfig{
		ServerAddress: "localhost:50051",
		UseTLS:        false,
	}

	// Mock grpc.NewClientConn to simulate a successful connection
	conn, err := grpc.Dial(clientConfig.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)
	defer conn.Close()

	// Test the NewGclient function
	client, grpcConn := NewGclient(clientConfig, mockStorage, logger)

	assert.NotNil(t, client)
	assert.NotNil(t, grpcConn)
	// assert.Equal(t, mockStorage, client.Storage)
	// assert.Equal(t, logger, client.log)

	// Ensure connection is established
	assert.Equal(t, grpcConn.Target(), clientConfig.ServerAddress)
}

// Test that the interceptor adds the authorization token when the method is not skipped
func TestUnaryClientInterceptor_AddsToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mock storage with a token
	mockStorage := &MemStorage{
		Token: "test-token",
	}

	// Define method that is not in SkipCheckMethods
	testMethod := "/TestService/Method"
	SkipCheckMethods = map[string]struct{}{} // Ensure no methods are skipped

	// Define a mock invoker
	mockInvoker := func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		opts ...grpc.CallOption,
	) error {
		// Check if the "authorization" metadata has been set
		md, ok := metadata.FromOutgoingContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, []string{"Bearer test-token"}, md["authorization"])
		return nil
	}

	// Get the unary client interceptor
	interceptor := getUnaryClientInterceptor(mockStorage)

	// Call the interceptor and check if token was added to the context
	err := interceptor(context.Background(), testMethod, nil, nil, nil, mockInvoker)
	assert.NoError(t, err)
}

// Test that the interceptor does not add the authorization token when the method is skipped
func TestUnaryClientInterceptor_SkipToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mock storage with a token
	mockStorage := &MemStorage{
		Token: "test-token",
	}

	// Define method that is in SkipCheckMethods
	testMethod := "/TestService/Method"
	SkipCheckMethods = map[string]struct{}{
		testMethod: {},
	}

	// Define a mock invoker
	mockInvoker := func(
		ctx context.Context,
		method string,
		req, reply interface{},
		cc *grpc.ClientConn,
		opts ...grpc.CallOption,
	) error {
		// Ensure no "authorization" metadata is set
		md, ok := metadata.FromOutgoingContext(ctx)
		assert.False(t, ok)
		assert.Nil(t, md)
		return nil
	}

	// Get the unary client interceptor
	interceptor := getUnaryClientInterceptor(mockStorage)

	// Call the interceptor and check if token was skipped
	err := interceptor(context.Background(), testMethod, nil, nil, nil, mockInvoker)
	assert.NoError(t, err)
}

// Test that the stream interceptor adds the authorization token when the method is not skipped
func TestStreamClientInterceptor_AddsToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mock storage with a token
	mockStorage := &MemStorage{
		Token: "test-token",
	}

	// Define method that is not in SkipCheckMethods
	testMethod := "/TestService/StreamMethod"
	SkipCheckMethods = map[string]struct{}{} // Ensure no methods are skipped

	// Define a mock streamer
	mockStreamer := func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		// Check if the "authorization" metadata has been set
		md, ok := metadata.FromOutgoingContext(ctx)
		assert.True(t, ok)
		assert.Equal(t, []string{"Bearer test-token"}, md["authorization"])

		// Simulate creating a stream (return nil in this simple mock case)
		return nil, nil
	}

	// Get the stream client interceptor
	interceptor := getStreamClientInterceptor(mockStorage)

	// Call the interceptor and check if token was added to the context
	_, err := interceptor(context.Background(), &grpc.StreamDesc{}, nil, testMethod, mockStreamer)
	assert.NoError(t, err)
}

// Test that the stream interceptor does not add the authorization token when the method is skipped
func TestStreamClientInterceptor_SkipToken(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Setup mock storage with a token
	mockStorage := &MemStorage{
		Token: "test-token",
	}

	// Define method that is in SkipCheckMethods
	testMethod := "/TestService/StreamMethod"
	SkipCheckMethods = map[string]struct{}{
		testMethod: {},
	}

	// Define a mock streamer
	mockStreamer := func(
		ctx context.Context,
		desc *grpc.StreamDesc,
		cc *grpc.ClientConn,
		method string,
		opts ...grpc.CallOption,
	) (grpc.ClientStream, error) {
		// Ensure no "authorization" metadata is set
		md, ok := metadata.FromOutgoingContext(ctx)
		assert.False(t, ok)
		assert.Nil(t, md)

		// Simulate creating a stream (return nil in this simple mock case)
		return nil, nil
	}

	// Get the stream client interceptor
	interceptor := getStreamClientInterceptor(mockStorage)

	// Call the interceptor and check if token was skipped
	_, err := interceptor(context.Background(), &grpc.StreamDesc{}, nil, testMethod, mockStreamer)
	assert.NoError(t, err)
}
