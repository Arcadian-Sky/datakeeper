package router

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"testing"

	pbservice "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/service/v1"
	pbuser "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/user/v1"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/Arcadian-Sky/datakkeeper/internal/server/repository"
	"github.com/Arcadian-Sky/datakkeeper/internal/server/router/jwtrule"
	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"github.com/Arcadian-Sky/datakkeeper/mocks"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func createTestMockServer(t *testing.T) (server *GRPCServer) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoUser := mocks.NewMockUserRepository(ctrl)
	mockRepoFile := mocks.NewMockFileRepository(ctrl)
	mockRepoData := mocks.NewMockDataRepository(ctrl)
	mockLogger := logrus.New()

	server = &GRPCServer{
		reposervice: mockRepoFile,
		repodata:    mockRepoData,
		repouser:    mockRepoUser,
		log:         mockLogger,
		cfg:         &settings.InitedFlags{SecretKey: "test-secret"},
	}

	return server
}

// Тестирование метода Register
func TestGRPCServer_Register(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoUser := mocks.NewMockUserRepository(ctrl)
	mockRepoService := mocks.NewMockFileRepository(ctrl)
	logg := logrus.New()
	// logg.SetLevel(logrus.TraceLevel)
	// logg.SetFormatter(&logrus.TextFormatter{})

	server := &GRPCServer{
		log:         logg,
		repouser:    mockRepoUser,
		reposervice: mockRepoService,
		cfg:         &settings.InitedFlags{SecretKey: "test-secret"},
	}

	tests := []struct {
		name      string
		input     *pbuser.RegisterRequest
		mockSetup func()
		wantErr   bool
		wantResp  *pbuser.RegisterResponse
	}{
		{
			name: "Success",
			input: &pbuser.RegisterRequest{
				Login:    "testuser",
				Password: "password",
			},
			mockSetup: func() {
				mockRepoUser.EXPECT().Register(gomock.Any(), gomock.Any()).Return(int64(1), nil)

				mockRepoService.EXPECT().CreateContainer(gomock.Any(), gomock.Any()).DoAndReturn(func(ctx context.Context, u *model.User) (model.User, error) {
					assert.NotNil(t, u.ID)
					assert.NotNil(t, u.Login)
					bucketName := "bucketuid" + strconv.Itoa(int(u.ID))
					u.Bucket = bucketName
					assert.Equal(t, bucketName, u.Bucket)
					return *u, nil
				})

				// _, _ := jwtrule.Generate(1, "test-secret")
			},
			wantErr:  false,
			wantResp: &pbuser.RegisterResponse{Success: true, Message: "user testuser (userid: 1) was created\nbucket container bucketuid1 (userid: 1) was created\n", AuthToken: "expected-jwt-token"},
		},
		{
			name: "Empty Login",
			input: &pbuser.RegisterRequest{
				Login:    "",
				Password: "password",
			},
			mockSetup: func() {
			},
			wantErr:  true,
			wantResp: nil,
		},
		{
			name: "Register Error",
			input: &pbuser.RegisterRequest{
				Login:    "testuser",
				Password: "password",
			},
			mockSetup: func() {
				mockRepoUser.EXPECT().Register(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("registration error"))
			},
			wantErr:  true,
			wantResp: nil,
		},
		{
			name: "Create Container Error",
			input: &pbuser.RegisterRequest{
				Login:    "testuser",
				Password: "password",
			},
			mockSetup: func() {
				mockRepoUser.EXPECT().Register(gomock.Any(), gomock.Any()).Return(int64(1), nil)
				mockRepoService.EXPECT().CreateContainer(gomock.Any(), gomock.Any()).Return(model.User{}, errors.New("create container error"))
			},
			wantErr:  true,
			wantResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			gotResp, err := server.Register(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantResp != nil && gotResp != nil {
				// Here you might need to verify the JWT token generated if necessary
				assert.Equal(t, tt.wantResp.Success, gotResp.Success)
				assert.Equal(t, tt.wantResp.Message, gotResp.Message)
				// assert.Equal(t, tt.wantResp.AuthToken, gotResp.AuthToken)
			}
		})
	}
}

func TestGRPCServer_Authenticate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepoUser := mocks.NewMockUserRepository(ctrl)
	mockRepoFile := mocks.NewMockFileRepository(ctrl)
	mockRepoData := mocks.NewMockDataRepository(ctrl)
	mockLogger := logrus.New()

	server := &GRPCServer{
		reposervice: mockRepoFile,
		repodata:    mockRepoData,
		repouser:    mockRepoUser,
		log:         mockLogger,
		cfg:         &settings.InitedFlags{SecretKey: "test-secret"},
	}

	tests := []struct {
		name      string
		input     *pbuser.AuthenticateRequest
		mockSetup func()
		wantErr   bool
		wantResp  *pbuser.AuthenticateResponse
	}{
		{
			name: "Success",
			input: &pbuser.AuthenticateRequest{
				Login:    "testuser",
				Password: "password",
			},
			mockSetup: func() {
				user := &model.User{ID: 1, Login: "testuser", Password: "password"}

				// Mock Auth method to return a user
				mockRepoUser.EXPECT().
					Auth(gomock.Any(), gomock.Any()).
					DoAndReturn(func(ctx context.Context, u *model.User) (*model.User, error) {
						assert.Equal(t, "testuser", u.Login)
						assert.Equal(t, "password", u.Password)
						return user, nil
					})
			},
			wantErr: false,
			wantResp: &pbuser.AuthenticateResponse{
				Success:   true,
				AuthToken: "jwt-token",
				Message:   "authorized as userID: 1 token generated",
			},
		},
		{
			name: "InvalidArgument",
			input: &pbuser.AuthenticateRequest{
				Login:    "",
				Password: "",
			},
			mockSetup: func() {},
			wantErr:   true,
			wantResp:  nil,
		},
		{
			name: "AuthFailure",
			input: &pbuser.AuthenticateRequest{
				Login:    "testuser",
				Password: "password",
			},
			mockSetup: func() {
				// Mock Auth method to return an error
				mockRepoUser.EXPECT().
					Auth(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("auth error"))
			},
			wantErr:  true,
			wantResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			gotResp, err := server.Authenticate(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Authenticate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotResp != nil {
				assert.Equal(t, tt.wantResp.Success, gotResp.Success)
				// assert.Equal(t, tt.wantResp.AuthToken, gotResp.AuthToken)
				assert.Equal(t, tt.wantResp.Message, gotResp.Message)
			}
		})
	}
}

func TestGRPCServer_GetFileList(t *testing.T) {
	server := createTestMockServer(t)

	ctx := context.Background()
	ctx = jwtrule.SetUserIDToCTX(ctx, 1) // Assuming you have a way to set user ID in context

	tests := []struct {
		name      string
		input     *pbservice.ListFileRequest
		mockSetup func()
		wantErr   bool
		wantResp  *pbservice.ListFileResponse
	}{
		{
			name:  "Success",
			input: &pbservice.ListFileRequest{},
			mockSetup: func() {
				mockFileItems := []model.FileItem{
					{Hash: "hash1", Name: "file1"},
					{Hash: "hash2", Name: "file2"},
				}

				mockRepoFile := server.reposervice.(*mocks.MockFileRepository)
				mockRepoFile.EXPECT().
					GetFileList(gomock.Any(), &model.User{ID: 1}).
					Return(mockFileItems, nil).
					Times(1)
			},
			wantErr: false,
			wantResp: &pbservice.ListFileResponse{
				Fileitem: []*pbservice.FileItem{
					{Key: "hash1", Name: "file1"},
					{Key: "hash2", Name: "file2"},
				},
			},
		},
		{
			name:  "Error",
			input: &pbservice.ListFileRequest{},
			mockSetup: func() {
				mockRepoFile := server.reposervice.(*mocks.MockFileRepository)
				mockRepoFile.EXPECT().
					GetFileList(gomock.Any(), &model.User{ID: 1}).
					Return(nil, fmt.Errorf("db error")).
					Times(1)
			},
			wantErr:  true,
			wantResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			gotResp, err := server.GetFileList(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetFileList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotResp != nil {
				assert.Equal(t, tt.wantResp.Fileitem, gotResp.Fileitem)
			}
		})
	}
}

func TestGRPCServer_SaveData(t *testing.T) {
	server := createTestMockServer(t)

	ctx := context.Background()
	ctx = jwtrule.SetUserIDToCTX(ctx, 1) // Assuming you have a way to set user ID in context

	tests := []struct {
		name      string
		input     *pbservice.SaveDataRequest
		mockSetup func()
		wantErr   bool
		wantResp  *pbservice.UploadStatus
	}{
		{
			name: "Success",
			input: &pbservice.SaveDataRequest{
				Data: &pbservice.Data{
					Type:     pbservice.DataType_DATA_TYPE_TYPE_CREDIT_CARD,
					Title:    "test title",
					Card:     "1234-5678-9012-3456",
					Login:    "testuser",
					Password: "password",
				},
			},
			mockSetup: func() {
				mockData := model.Data{
					Type:     repository.DataTypeCARD,
					UserID:   1,
					Title:    "test title",
					Card:     "1234-5678-9012-3456",
					Login:    "testuser",
					Password: "password",
				}

				mockRepoData := server.repodata.(*mocks.MockDataRepository)
				mockRepoData.EXPECT().
					Save(gomock.Any(), &mockData).
					Return(int64(1), nil).
					Times(1)

				mockRepoUser := server.repouser.(*mocks.MockUserRepository)
				mockRepoUser.EXPECT().
					SetLastUpdate(gomock.Any(), gomock.Any()).
					Return(&model.User{}, nil).
					Times(1)
			},
			wantErr: false,
			wantResp: &pbservice.UploadStatus{
				Success: true,
				Message: "empty",
			},
		},
		{
			name: "SaveDataError",
			input: &pbservice.SaveDataRequest{
				Data: &pbservice.Data{
					Type:     pbservice.DataType_DATA_TYPE_TYPE_CREDIT_CARD,
					Title:    "test title",
					Card:     "1234-5678-9012-3456",
					Login:    "testuser",
					Password: "password",
				},
			},
			mockSetup: func() {
				mockData := model.Data{
					Type:     repository.DataTypeCARD,
					UserID:   1,
					Title:    "test title",
					Card:     "1234-5678-9012-3456",
					Login:    "testuser",
					Password: "password",
				}

				mockRepoData := server.repodata.(*mocks.MockDataRepository)
				mockRepoData.EXPECT().
					Save(gomock.Any(), &mockData).
					Return(int64(0), fmt.Errorf("db error")).
					Times(1)
			},
			wantErr:  true,
			wantResp: nil,
		},
		{
			name: "SetLastUpdateError",
			input: &pbservice.SaveDataRequest{
				Data: &pbservice.Data{
					Type:     pbservice.DataType_DATA_TYPE_TYPE_CREDIT_CARD,
					Title:    "test title",
					Card:     "1234-5678-9012-3456",
					Login:    "testuser",
					Password: "password",
				},
			},
			mockSetup: func() {
				mockData := model.Data{
					Type:     repository.DataTypeCARD,
					UserID:   1,
					Title:    "test title",
					Card:     "1234-5678-9012-3456",
					Login:    "testuser",
					Password: "password",
				}

				mockRepoData := server.repodata.(*mocks.MockDataRepository)
				mockRepoData.EXPECT().
					Save(gomock.Any(), &mockData).
					Return(int64(1), nil).
					Times(1)

				mockRepoUser := server.repouser.(*mocks.MockUserRepository)
				mockRepoUser.EXPECT().
					SetLastUpdate(gomock.Any(), gomock.Any()).
					Return(&model.User{}, fmt.Errorf("update error")).
					Times(1)
			},
			wantErr:  true,
			wantResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			gotResp, err := server.SaveData(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveData() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotResp != nil {
				assert.Equal(t, tt.wantResp.Success, gotResp.Success)
				assert.Equal(t, tt.wantResp.Message, gotResp.Message)
			}
		})
	}
}

func TestGRPCServer_GetDataList(t *testing.T) {
	server := createTestMockServer(t)

	ctx := context.Background()
	ctx = jwtrule.SetUserIDToCTX(ctx, 1) // Assuming you have a way to set user ID in context

	tests := []struct {
		name      string
		input     *pbservice.ListDataRequest
		mockSetup func()
		wantErr   bool
		wantResp  *pbservice.ListDataResponse
	}{
		{
			name:  "Success",
			input: &pbservice.ListDataRequest{},
			mockSetup: func() {
				data := []model.Data{
					{
						ID:       1,
						Title:    "Title 1",
						Type:     repository.DataTypeCARD,
						Card:     "1234-5678-9012-3456",
						Login:    "user1",
						Password: "pass1",
					},
					{
						ID:       2,
						Title:    "Title 2",
						Type:     repository.DataTypeLOGPASS,
						Card:     "9876-5432-1098-7654",
						Login:    "user2",
						Password: "pass2",
					},
				}

				mockRepoData := server.repodata.(*mocks.MockDataRepository)
				mockRepoData.EXPECT().
					GetList(gomock.Any(), gomock.Any()).
					Return(data, nil).
					Times(1)
			},
			wantErr: false,
			wantResp: &pbservice.ListDataResponse{
				Data: []*pbservice.Data{
					{
						Id:       1,
						Title:    "Title 1",
						Type:     pbservice.DataType_DATA_TYPE_TYPE_CREDIT_CARD,
						Card:     "1234-5678-9012-3456",
						Login:    "user1",
						Password: "pass1",
					},
					{
						Id:       2,
						Title:    "Title 2",
						Type:     pbservice.DataType_DATA_TYPE_TYPE_LOGIN_PASSWORD,
						Card:     "9876-5432-1098-7654",
						Login:    "user2",
						Password: "pass2",
					},
				},
			},
		},
		{
			name:  "GetListError",
			input: &pbservice.ListDataRequest{}, // Adjust input if necessary
			mockSetup: func() {
				mockRepoData := server.repodata.(*mocks.MockDataRepository)
				mockRepoData.EXPECT().
					GetList(gomock.Any(), gomock.Any()).
					Return(nil, fmt.Errorf("db error")).
					Times(1)
			},
			wantErr:  true,
			wantResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			gotResp, err := server.GetDataList(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetDataList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotResp != nil {
				assert.ElementsMatch(t, tt.wantResp.Data, gotResp.Data)
			}
		})
	}
}

func TestGRPCServer_DeleteFile(t *testing.T) {
	server := createTestMockServer(t)

	ctx := context.Background()
	ctx = jwtrule.SetUserIDToCTX(ctx, 1) // Assuming you have a way to set user ID in context

	tests := []struct {
		name      string
		input     *pbservice.DeleteFileRequest
		mockSetup func()
		wantErr   bool
		wantResp  *pbservice.UploadStatus
	}{
		{
			name: "Success",
			input: &pbservice.DeleteFileRequest{
				Filename: "file-to-delete.txt",
			},
			mockSetup: func() {
				// Mock DeleteFile method to return no error
				server.reposervice.(*mocks.MockFileRepository).EXPECT().
					DeleteFile(gomock.Any(), "file-to-delete.txt", gomock.Any()).
					Return(nil).
					Times(1)
			},
			wantErr: false,
			wantResp: &pbservice.UploadStatus{
				Success: true,
				Message: "data was deleted",
			},
		},
		{
			name: "DeleteFileError",
			input: &pbservice.DeleteFileRequest{
				Filename: "file-to-delete.txt",
			},
			mockSetup: func() {
				// Mock DeleteFile method to return an error
				server.reposervice.(*mocks.MockFileRepository).EXPECT().
					DeleteFile(gomock.Any(), "file-to-delete.txt", gomock.Any()).
					Return(fmt.Errorf("delete error")).
					Times(1)
			},
			wantErr:  true,
			wantResp: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mockSetup()
			gotResp, err := server.DeleteFile(ctx, tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && gotResp != nil {
				assert.Equal(t, tt.wantResp.Success, gotResp.Success)
				assert.Equal(t, tt.wantResp.Message, gotResp.Message)
			}
		})
	}
}

func TestGetPType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected pbservice.DataType
	}{
		{
			name:     "DataTypeCARD",
			input:    repository.DataTypeCARD,
			expected: pbservice.DataType_DATA_TYPE_TYPE_CREDIT_CARD,
		},
		{
			name:     "DataTypeLOGPASS",
			input:    repository.DataTypeLOGPASS,
			expected: pbservice.DataType_DATA_TYPE_TYPE_LOGIN_PASSWORD,
		},
		{
			name:     "UnknownType",
			input:    "unknown",
			expected: pbservice.DataType_DATA_TYPE_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getPType(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

func TestGetType(t *testing.T) {
	tests := []struct {
		name     string
		input    pbservice.DataType
		expected string
	}{
		{
			name:     "DataTypeCARD",
			input:    pbservice.DataType_DATA_TYPE_TYPE_CREDIT_CARD,
			expected: repository.DataTypeCARD,
		},
		{
			name:     "DataTypeLOGPASS",
			input:    pbservice.DataType_DATA_TYPE_TYPE_LOGIN_PASSWORD,
			expected: repository.DataTypeLOGPASS,
		},
		{
			name:     "UnspecifiedType",
			input:    pbservice.DataType_DATA_TYPE_UNSPECIFIED,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := getType(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
