package client

import (
	"testing"

	pbuser "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/user/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// TestRegister tests the Register method of GRPCClient
func TestGRPCClient_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := pbuser.NewMockUserServiceClient(ctrl)
	mockLogger := logrus.New()
	storage := NewMemStorage()

	client := &GRPCClient{
		User:    mockUserClient,
		log:     mockLogger,
		Storage: storage,
	}

	login := "testUser"
	password := "testPassword"
	authToken := "testAuthToken"

	// Mock the Register method
	mockUserClient.EXPECT().
		Register(gomock.Any(), &pbuser.RegisterRequest{
			Login:    login,
			Password: password,
		}).
		Return(&pbuser.RegisterResponse{
			Success:   true,
			AuthToken: authToken,
		}, nil).
		Times(1)

	// Call the Register method
	err := client.Register(login, password)

	// Verify the result
	assert.NoError(t, err, "Expected no error from Register method")
	assert.Equal(t, authToken, storage.Token, "Expected storage token to match the auth token from the response")
}

func TestGRPCClient_Register_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := pbuser.NewMockUserServiceClient(ctrl)
	mockLogger := logrus.New()
	storage := NewMemStorage()

	client := &GRPCClient{
		User:    mockUserClient,
		log:     mockLogger,
		Storage: storage,
	}

	login := "testUser"
	password := "testPassword"
	errorMessage := "registration error"

	// Mock the Register method to return an error
	mockUserClient.EXPECT().
		Register(gomock.Any(), &pbuser.RegisterRequest{
			Login:    login,
			Password: password,
		}).
		Return(nil, status.Errorf(codes.Unknown, errorMessage)).
		Times(1)

	// Call the Register method
	err := client.Register(login, password)

	// Verify the result
	assert.Error(t, err, "Expected error from Register method")
	assert.Empty(t, storage.Token, "Expected storage token to be empty")
}

// TestGRPCClient_Authenticate tests the Authenticate method of GRPCClient
func TestGRPCClient_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := pbuser.NewMockUserServiceClient(ctrl)
	mockLogger := logrus.New()
	storage := NewMemStorage()

	client := &GRPCClient{
		User:    mockUserClient,
		log:     mockLogger,
		Storage: storage,
	}

	login := "testUser"
	password := "testPassword"
	authToken := "testAuthToken"

	// Mock the Authenticate method
	mockUserClient.EXPECT().
		Authenticate(gomock.Any(), &pbuser.AuthenticateRequest{
			Login:    login,
			Password: password,
		}).
		Return(&pbuser.AuthenticateResponse{
			Success:   true,
			AuthToken: authToken,
		}, nil).
		Times(1)

	// Call the Authenticate method
	err := client.Authenticate(login, password)

	// Verify the result
	assert.NoError(t, err, "Expected no error from Authenticate method")
	assert.Equal(t, authToken, storage.Token, "Expected storage token to match the auth token from the response")
}

func TestGRPCClient_Authenticate_Failure(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockUserClient := pbuser.NewMockUserServiceClient(ctrl)
	mockLogger := logrus.New()
	storage := NewMemStorage()

	client := &GRPCClient{
		User:    mockUserClient,
		log:     mockLogger,
		Storage: storage,
	}

	login := "testUser"
	password := "testPassword"
	errorMessage := "authentication error"

	// Mock the Authenticate method to return an error
	mockUserClient.EXPECT().
		Authenticate(gomock.Any(), &pbuser.AuthenticateRequest{
			Login:    login,
			Password: password,
		}).
		Return(nil, status.Errorf(codes.Unauthenticated, errorMessage)).
		Times(1)

	// Call the Authenticate method
	err := client.Authenticate(login, password)

	// Verify the result
	assert.Error(t, err, "Expected error from Authenticate method")
	assert.Empty(t, storage.Token, "Expected storage token to be empty")
}
