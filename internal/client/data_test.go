package client

import (
	"fmt"
	"testing"

	pbservice "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/service/v1"
	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetDataList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Define expected behavior
	mockDataClient.EXPECT().
		GetDataList(gomock.Any(), gomock.Any()).
		Return(&pbservice.ListDataResponse{
			Data: []*pbservice.Data{
				{
					Id:       1,
					Title:    "Test Data",
					Type:     pbservice.DataType_DATA_TYPE_TYPE_CREDIT_CARD,
					Login:    "user123",
					Card:     "card123",
					Password: "pass123",
				},
			},
		}, nil).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	dataList, err := client.GetDataList()

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, dataList, 1)
	assert.Equal(t, model.Data{
		ID:       1,
		Title:    "Test Data",
		Type:     "DATA_TYPE_TYPE_CREDIT_CARD", // Update based on your enum mapping
		Login:    "user123",
		Card:     "card123",
		Password: "pass123",
	}, dataList[0])
}

func TestGetDataList_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Define expected behavior
	mockDataClient.EXPECT().
		GetDataList(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("test error")).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	dataList, err := client.GetDataList()

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "test error")
	assert.Len(t, dataList, 0)
}

func TestSaveLoginPass_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Define expected behavior
	mockDataClient.EXPECT().
		SaveData(gomock.Any(), gomock.Any()).
		Return(&pbservice.UploadStatus{}, nil).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	err := client.SaveLoginPass("example.com", "user123", "pass123")

	// Assertions
	assert.NoError(t, err)
}

func TestSaveLoginPass_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Define expected behavior
	mockDataClient.EXPECT().
		SaveData(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("test error")).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	err := client.SaveLoginPass("example.com", "user123", "pass123")

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "test error")
}

func TestSaveCard_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Define expected behavior
	mockDataClient.EXPECT().
		SaveData(gomock.Any(), gomock.Any()).
		Return(&pbservice.UploadStatus{}, nil).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	err := client.SaveCard("Visa", "4111111111111111")

	// Assertions
	assert.NoError(t, err)
}

func TestSaveCard_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Define expected behavior
	mockDataClient.EXPECT().
		SaveData(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("test error")).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	err := client.SaveCard("Visa", "4111111111111111")

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "test error")
}

func TestDelete_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Define expected behavior
	mockDataClient.EXPECT().
		DeleteData(gomock.Any(), gomock.Any()).
		Return(&pbservice.UploadStatus{}, nil).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	err := client.Delete(12345)

	// Assertions
	assert.NoError(t, err)
}

func TestDelete_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Define expected behavior
	mockDataClient.EXPECT().
		DeleteData(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("test error")).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	err := client.Delete(12345)

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "test error")
}
