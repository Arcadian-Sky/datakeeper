package client

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"

	pbservice "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/service/v1"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetFileList_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Define expected response
	expectedFileItems := []*pbservice.FileItem{
		{
			Key:  "fileHash1",
			Name: "fileName1",
		},
		{
			Key:  "fileHash2",
			Name: "fileName2",
		},
	}

	mockDataClient.EXPECT().
		GetFileList(gomock.Any(), gomock.Any()).
		Return(&pbservice.ListFileResponse{Fileitem: expectedFileItems}, nil).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	fileList, err := client.GetFileList()

	// Assertions
	assert.NoError(t, err)
	assert.Len(t, fileList, 2)
	assert.Equal(t, "fileHash1", fileList[0].Hash)
	assert.Equal(t, "fileName1", fileList[0].Name)
	assert.Equal(t, "fileHash2", fileList[1].Hash)
	assert.Equal(t, "fileName2", fileList[1].Name)
}

func TestGetFileList_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Define expected behavior
	mockDataClient.EXPECT().
		GetFileList(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("test error")).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	fileList, err := client.GetFileList()

	// Assertions
	assert.Error(t, err)
	assert.Empty(t, fileList)
	assert.EqualError(t, err, "test error")
}

func TestDeleteFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Define expected behavior
	mockDataClient.EXPECT().
		DeleteFile(gomock.Any(), gomock.Any()).
		Return(&pbservice.UploadStatus{}, nil).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	err := client.DeleteFile("test-file")

	// Assertions
	assert.NoError(t, err)
}

func TestDeleteFile_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Define expected behavior
	mockDataClient.EXPECT().
		DeleteFile(gomock.Any(), gomock.Any()).
		Return(nil, fmt.Errorf("test error")).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	err := client.DeleteFile("test-file")

	// Assertions
	assert.Error(t, err)
	assert.EqualError(t, err, "test error")
}

func TestUploadFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockStream := pbservice.NewMockDataKeeperService_UploadFileClient(ctrl)
	mockLogger := logrus.New()

	// Mock behaviors
	mockDataClient.EXPECT().
		UploadFile(gomock.Any()).
		Return(mockStream, nil).
		Times(1)

	mockStream.EXPECT().
		Send(gomock.Any()).
		Return(nil).
		Times(1)
	mockStream.EXPECT().
		CloseAndRecv().
		Return(&pbservice.UploadStatus{Success: true, Message: "Upload successful"}, nil).
		Times(1)

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("test content")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test
	err = client.UploadFile(tempFile.Name())

	// Assertions
	assert.NoError(t, err)
}

func TestUploadFile_OpenFileError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Call the method to test with an invalid file path
	err := client.UploadFile("/invalid/path/to/file")

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "could not open file")
}

func TestUploadFile_StreamError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockLogger := logrus.New()

	// Mock behaviors
	mockDataClient.EXPECT().
		UploadFile(gomock.Any()).
		Return(nil, fmt.Errorf("stream error")).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: nil, // Update if needed
	}

	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "testfile")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tempFile.Name())

	_, err = tempFile.WriteString("test content")
	if err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	// Call the method to test
	err = client.UploadFile(tempFile.Name())

	// Assertions
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "error creating stream")
}

func TestGetFile_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockStream := pbservice.NewMockDataKeeperService_GetFileClient(ctrl)
	mockLogger := logrus.New()

	// Mock behaviors
	mockDataClient.EXPECT().
		GetFile(gomock.Any(), gomock.Any()).
		Return(mockStream, nil).
		Times(1)

	mockStream.EXPECT().
		Recv().
		Return(&pbservice.FileChunk{Data: []byte("chunk1")}, nil).
		Times(1)

	mockStream.EXPECT().
		Recv().
		Return(&pbservice.FileChunk{Data: []byte("chunk2")}, nil).
		Times(1)

	mockStream.EXPECT().
		Recv().
		Return(nil, io.EOF).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: &MemStorage{PfilesDir: t.TempDir()},
	}

	// Create a file for testing
	fileName := "testfile.txt"
	filePath := filepath.Join(client.Storage.PfilesDir, fileName)
	err := client.GetFile(fileName)

	// Assertions
	assert.NoError(t, err)

	// Verify file content
	content, err := os.ReadFile(filePath)
	assert.NoError(t, err)
	assert.Equal(t, "chunk1chunk2", string(content))
}

func TestGetFile_RecvError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock objects
	mockDataClient := pbservice.NewMockDataKeeperServiceClient(ctrl)
	mockStream := pbservice.NewMockDataKeeperService_GetFileClient(ctrl)
	mockLogger := logrus.New()

	mockDataClient.EXPECT().
		GetFile(gomock.Any(), gomock.Any()).
		Return(mockStream, nil).
		Times(1)

	mockStream.EXPECT().
		Recv().
		Return(nil, fmt.Errorf("recv error")).
		Times(1)

	client := &GRPCClient{
		log:     mockLogger,
		Data:    mockDataClient,
		Storage: &MemStorage{PfilesDir: t.TempDir()},
	}

	err := client.GetFile("testfile.txt")

	// Assertions
	assert.Error(t, err)
}
