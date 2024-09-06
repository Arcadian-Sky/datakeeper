package repository

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"testing"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/Arcadian-Sky/datakkeeper/mocks"
	"github.com/golang/mock/gomock"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

// func NewMinioStorage() (*minio.Client, error) {
// 	endpoint := "localhost:9000"
// 	accessKeyID := "minioadmin"
// 	secretAccessKey := "minioadminpassword"
// 	creds := credentials.NewStaticV4(accessKeyID, secretAccessKey, "")
// 	useSSL := false
// 	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// 	defer cancel()
// 	// Создание нового клиента MinIO
// 	client, err := minio.New(endpoint, &minio.Options{
// 		Creds:  creds,
// 		Secure: useSSL,
// 	})
// 	if err != nil {
// 		return nil, err
// 	}

//		// Проверка подключения к MinIO
//		_, err = client.ListBuckets(ctx)
//		if err != nil {
//			return nil, err
//		}
//		return client, nil
//	}
// func TestFileRepo_GetFileList(t *testing.T) {

// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
// ctx, cancel := context.WithCancel(context.Background())
// defer cancel()
// user := model.User{
// 	ID: 9,
// }
// lg := logrus.New()
// st, _ := NewMinioStorage()
// repo := NewFileRepository(st, lg, &ctx)

// _, err := repo.GetFileList(ctx, &user)
// assert.NoError(t, err)

// tests := []struct {
// 	name    string
// 	f       *FileRepo
// 	args    args
// 	want    []minio.ObjectInfo
// 	wantErr bool
// }{
// 	// TODO: Add test cases.
// }
// for _, tt := range tests {
// 	t.Run(tt.name, func(t *testing.T) {
// 		got, err := tt.f.GetFileList(tt.args.ctx, tt.args.user, tt.args.data)
// 		if (err != nil) != tt.wantErr {
// 			t.Errorf("FileRepo.GetFileList() error = %v, wantErr %v", err, tt.wantErr)
// 			return
// 		}
// 		if !reflect.DeepEqual(got, tt.want) {
// 			t.Errorf("FileRepo.GetFileList() = %v, want %v", got, tt.want)
// 		}
// 	})
// }
// }

func TestFileRepo_CreateContainer(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type args struct {
		ctx  *context.Context
		user *model.User
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMinio := mocks.NewMockMinioClient(ctrl)
	mockLogger := logrus.New()

	tests := []struct {
		name        string
		f           *FileRepo
		args        args
		setupMocks  func()
		want        model.User
		wantErr     bool
		expectedErr error
	}{
		{
			name: "NoUserID",
			f: &FileRepo{
				db:       mockMinio,
				log:      mockLogger,
				ctx:      &ctx,
				location: "us-east-1",
			},
			args: args{
				ctx:  &ctx,
				user: &model.User{ID: 0},
			},
			setupMocks:  func() {},
			want:        model.User{ID: 0},
			wantErr:     true,
			expectedErr: model.ErrCreateBucketNoUser,
		},
		{
			name: "BucketAlreadyExists",
			f: &FileRepo{
				db:       mockMinio,
				log:      mockLogger,
				ctx:      &ctx,
				location: "us-east-1",
			},
			args: args{
				ctx:  &ctx,
				user: &model.User{ID: 123},
			},
			setupMocks: func() {
				bucketName := "bucketuid" + strconv.Itoa(123)
				mockMinio.EXPECT().MakeBucket(gomock.Any(), bucketName, gomock.Any()).Return(errors.New("bucket exists"))
				mockMinio.EXPECT().BucketExists(gomock.Any(), bucketName).Return(true, nil)
			},
			want:        model.User{ID: 123},
			wantErr:     true,
			expectedErr: model.ErrCreateBucketExists,
		},
		{
			name: "CreateBucketSuccess",
			f: &FileRepo{
				db:       mockMinio,
				log:      mockLogger,
				ctx:      &ctx,
				location: "us-east-1",
			},
			args: args{
				ctx:  &ctx,
				user: &model.User{ID: 123},
			},
			setupMocks: func() {
				bucketName := "bucketuid" + strconv.Itoa(123)
				mockMinio.EXPECT().MakeBucket(gomock.Any(), bucketName, gomock.Any()).Return(nil)
			},
			want: model.User{
				ID:     123,
				Bucket: "bucketuid123",
			},
			wantErr: false,
		},
		{
			name: "BucketCreationFailed",
			f: &FileRepo{
				db:       mockMinio,
				log:      mockLogger,
				ctx:      &ctx,
				location: "us-east-1",
			},
			args: args{
				ctx:  &ctx,
				user: &model.User{ID: 123},
			},
			setupMocks: func() {
				bucketName := "bucketuid" + strconv.Itoa(123)
				mockMinio.EXPECT().MakeBucket(gomock.Any(), bucketName, gomock.Any()).Return(errors.New("failed to create bucket"))
				mockMinio.EXPECT().BucketExists(gomock.Any(), bucketName).Return(false, nil)
			},
			want:        model.User{ID: 123},
			wantErr:     true,
			expectedErr: model.ErrCreateBucketFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			got, err := tt.f.CreateContainer(*tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileRepo.CreateContainer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil && !errors.Is(err, tt.expectedErr) {
				t.Errorf("FileRepo.CreateContainer() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("FileRepo.CreateContainer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFileRepo_GetFile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMinio := mocks.NewMockMinioClient(ctrl)
	mockMinioObject := mocks.NewMockMinioObject(ctrl)
	mockLogger := logrus.New()

	type args struct {
		ctx    *context.Context
		fileID string
		user   *model.User
	}
	tests := []struct {
		name        string
		f           *FileRepo
		args        args
		setupMocks  func()
		want        *os.File
		wantErr     bool
		expectedErr bool
	}{
		{
			name: "NoUserID",
			f: &FileRepo{
				db:       mockMinio,
				log:      mockLogger,
				ctx:      &ctx,
				location: "us-east-1",
			},
			args: args{
				ctx:    &ctx,
				fileID: "file123",
				user:   &model.User{ID: 0},
			},
			setupMocks:  func() {},
			want:        nil,
			wantErr:     true,
			expectedErr: true,
		},
		{
			name: "GetObjectFailure",
			f: &FileRepo{
				db:       mockMinio,
				log:      mockLogger,
				ctx:      &ctx,
				location: "us-east-1",
			},
			args: args{
				ctx:    &ctx,
				fileID: "file123",
				user:   &model.User{ID: 123},
			},
			setupMocks: func() {
				bucketName := "bucketuid" + strconv.Itoa(123)
				mockMinio.EXPECT().GetObject(gomock.Any(), bucketName, "file123", gomock.Any()).Return(nil, errors.New("failed to get object"))
			},
			want:        nil,
			wantErr:     true,
			expectedErr: true,
		},
		// {
		// 	name: "CreateTempFileFailure",
		// 	f: &FileRepo{
		// 		db:       mockMinio,
		// 		log:      mockLogger,
		// 		ctx:      &ctx,
		// 		location: "us-east-1",
		// 	},
		// 	args: args{
		// 		ctx:    &ctx,
		// 		fileID: "file123",
		// 		user:   &model.User{ID: 123},
		// 	},
		// 	setupMocks: func() {
		// 		// Mock the GetObject method to return a mock object
		// 		bucketName := "bucketuid" + strconv.Itoa(123)
		// 		mockObject := io.NopCloser(bytes.NewReader([]byte("test content")))
		// 		mockMinio.EXPECT().GetObject(gomock.Any(), bucketName, "file123", gomock.Any()).Return(mockObject, nil)
		// 	},
		// 	want:        nil,
		// 	wantErr:     true,
		// 	expectedErr: true,
		// },
		{
			name: "Success",
			f: &FileRepo{
				db:       mockMinio,
				log:      mockLogger,
				ctx:      &ctx,
				location: "us-east-1",
			},
			args: args{
				ctx:    &ctx,
				fileID: "file123",
				user:   &model.User{ID: 123},
			},
			setupMocks: func() {
				// Mock the GetObject method to return a mock object
				bucketName := "bucketuid" + strconv.Itoa(123)
				// mock := io.NopCloser(bytes.NewReader([]byte("test content")))
				mockMinio.EXPECT().GetObject(gomock.Any(), bucketName, "file123", gomock.Any()).Return(mockMinioObject, nil)
				mockMinioObject.EXPECT().Read(gomock.Any()).Return(1, nil)
				mockMinioObject.EXPECT().Close()

				// Mock sequential reads for the object
				gomock.InOrder(
					mockMinioObject.EXPECT().Read(gomock.Any()).DoAndReturn(func(p []byte) (int, error) {
						copy(p, "test data")
						return len("test data"), nil // First call returns data
					}),
					mockMinioObject.EXPECT().Read(gomock.Any()).Return(0, io.EOF), // Second call returns EOF
					// mockMinioObject.EXPECT().Close().Return(nil),
				)
				// mockMinioObject.EXPECT().Close().Return(nil)

			},
			want:        nil,
			wantErr:     false,
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			got, err := tt.f.GetFile(*tt.args.ctx, tt.args.fileID, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileRepo.GetFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (err != nil) && !tt.expectedErr {
				t.Errorf("FileRepo.GetFile() error = %v, expectedErr %v", err, tt.expectedErr)
			}
			if !tt.wantErr && got != nil {
				if _, err := os.Stat(got.Name()); os.IsNotExist(err) {
					t.Errorf("FileRepo.GetFile() = file does not exist")
				} else {
					os.Remove(got.Name()) // Clean up after test
				}
			}

		})
	}
}

func TestFileRepo_DeleteFile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMinio := mocks.NewMockMinioClient(ctrl)
	// mockMinioObject := mocks.NewMockMinioObject(ctrl)
	mockLogger := logrus.New()

	type args struct {
		ctx      *context.Context
		fileName string
		user     *model.User
	}
	tests := []struct {
		name       string
		f          *FileRepo
		args       args
		setupMocks func()
		wantErr    bool
	}{
		{
			name: "Success",
			f: &FileRepo{
				db:       mockMinio,  // mockMinio - это мок MinIO клиента
				log:      mockLogger, // mockLogger - это мок логгера logrus
				ctx:      &ctx,
				location: "us-east-1",
			},
			args: args{
				ctx:      &ctx,
				fileName: "file123",
				user:     &model.User{ID: 123},
			},
			setupMocks: func() {
				bucketName := "bucketuid" + strconv.Itoa(123)

				mockMinio.EXPECT().
					RemoveObject(gomock.Any(), bucketName, "file123", gomock.Any()).
					Return(nil) // Успешное удаление
			},
			wantErr: false,
		},
		{
			name: "Error on RemoveObject",
			f: &FileRepo{
				db:       mockMinio,
				log:      mockLogger,
				ctx:      &ctx,
				location: "us-east-1",
			},
			args: args{
				ctx:      &ctx,
				fileName: "file123",
				user:     &model.User{ID: 123},
			},
			setupMocks: func() {
				bucketName := "bucketuid" + strconv.Itoa(123)

				mockMinio.EXPECT().
					RemoveObject(gomock.Any(), bucketName, "file123", gomock.Any()).
					Return(fmt.Errorf("mock remove object error")) // Возврат ошибки
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			err := tt.f.DeleteFile(*tt.args.ctx, tt.args.fileName, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileRepo.DeleteFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFileRepo_GetFileList(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMinio := mocks.NewMockMinioClient(ctrl)
	// mockMinioObject := mocks.NewMockMinioObject(ctrl)
	mockLogger := logrus.New()

	type args struct {
		ctx  *context.Context
		user *model.User
	}
	tests := []struct {
		name       string
		f          *FileRepo
		args       args
		setupMocks func()
		want       []model.FileItem
		wantErr    bool
	}{
		{
			name: "Success",
			f: &FileRepo{
				db:       mockMinio,  // mockMinio - это мок MinIO клиента
				log:      mockLogger, // mockLogger - это мок логгера logrus
				ctx:      &ctx,
				location: "us-east-1",
			},
			args: args{
				ctx:  &ctx,
				user: &model.User{ID: 123},
			},
			setupMocks: func() {
				bucketName := "bucketuid" + strconv.Itoa(123)

				// Создаем канал для эмуляции возвращаемых объектов
				objectCh := make(chan minio.ObjectInfo, 2)
				objectCh <- minio.ObjectInfo{Key: "file1.txt", ETag: "etag1"}
				objectCh <- minio.ObjectInfo{Key: "file2.txt", ETag: "etag2"}
				close(objectCh)

				mockMinio.EXPECT().
					ListObjects(gomock.Any(), bucketName, gomock.Any()).
					Return(objectCh)
			},
			want: []model.FileItem{
				{Hash: "etag1", Name: "file1.txt"},
				{Hash: "etag2", Name: "file2.txt"},
			},
			wantErr: false,
		},
		{
			name: "Error on user ID",
			f: &FileRepo{
				db:       mockMinio,
				log:      mockLogger,
				ctx:      &ctx,
				location: "us-east-1",
			},
			args: args{
				ctx:  &ctx,
				user: &model.User{ID: 0}, // Неверный ID
			},
			setupMocks: func() {},
			want:       nil,
			wantErr:    true,
		},
		{
			name: "Error in ListObjects",
			f: &FileRepo{
				db:       mockMinio,
				log:      mockLogger,
				ctx:      &ctx,
				location: "us-east-1",
			},
			args: args{
				ctx:  &ctx,
				user: &model.User{ID: 123},
			},
			setupMocks: func() {
				bucketName := "bucketuid" + strconv.Itoa(123)

				// Эмулируем ошибку во время получения объекта
				objectCh := make(chan minio.ObjectInfo, 1)
				objectCh <- minio.ObjectInfo{Err: fmt.Errorf("mock test")}
				close(objectCh)

				mockMinio.EXPECT().
					ListObjects(gomock.Any(), bucketName, gomock.Any()).
					Return(objectCh)
			},
			want:    nil,
			wantErr: false, // Ошибки не будет, просто файл не добавится в список
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()

			got, err := tt.f.GetFileList(*tt.args.ctx, tt.args.user)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileRepo.GetFileList() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			// Сравниваем JSON-объекты
			gotJSON, err := json.Marshal(got)
			if err != nil {
				t.Fatalf("Failed to marshal got: %v", err)
			}

			wantJSON, err := json.Marshal(tt.want)
			if err != nil {
				t.Fatalf("Failed to marshal want: %v", err)
			}

			if string(gotJSON) != string(wantJSON) {
				t.Errorf("FileRepo.GetFileList() = %v, want %v", string(gotJSON), string(wantJSON))
			}
		})
	}
}

func TestFileRepo_UploadFile(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockMinioClient := mocks.NewMockMinioClient(ctrl)
	// mockMinioObject := mocks.NewMockMinioObject(ctrl)
	mockLogger := logrus.New()

	mockMinioClient.EXPECT().
		PutObject(
			gomock.Any(),
			"bucketuid1",
			"testfile.txt",
			gomock.Any(),
			int64(-1), //-1 означает неизвестный размер)
			gomock.Any(),
		).Return(minio.UploadInfo{}, nil)

	type args struct {
		ctx        *context.Context
		user       *model.User
		objectName string
		file       *os.File
	}
	tests := []struct {
		name       string
		f          *FileRepo
		args       args
		setupMocks func()
		wantErr    bool
	}{
		{
			name: "Success",
			f: &FileRepo{
				ctx: &ctx,
				db:  mockMinioClient, // предполагаем, что вы используете mock для MinIO клиента
				log: mockLogger,
			},
			args: args{
				ctx:        &ctx,
				user:       &model.User{ID: 1},
				objectName: "testfile.txt",
				file:       createTestFile(t, "testfile.txt"), // Создаем файл для теста
			},
			setupMocks: func() {

			},
			wantErr: false,
		},
		{
			name: "NoUserBucket",
			f: &FileRepo{
				db:  mockMinioClient,
				log: logrus.New(),
			},
			args: args{
				ctx:        &ctx,
				user:       &model.User{ID: 0}, // Пользователь с ID = 0
				objectName: "testfile.txt",
				file:       createTestFile(t, "testfile.txt"),
			},
			setupMocks: func() {

			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMocks()
			err := tt.f.UploadFile(*tt.args.ctx, tt.args.user, tt.args.objectName, tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("FileRepo.UploadFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Вспомогательная функция для создания тестового файла в директории tmp
func createTestFile(t *testing.T, fileName string) *os.File {
	tmpDir := filepath.Join("tmp")
	err := os.MkdirAll(tmpDir, os.ModePerm) // Создаем директорию tmp, если она не существует
	if err != nil {
		t.Fatalf("Ошибка при создании директории tmp: %v", err)
	}

	filePath := filepath.Join(tmpDir, fileName)
	testFile, err := os.Create(filePath)
	if err != nil {
		t.Fatalf("Ошибка при создании тестового файла: %v", err)
	}

	t.Cleanup(func() { os.Remove(filePath) }) // Удаление файла после завершения теста
	return testFile
}
