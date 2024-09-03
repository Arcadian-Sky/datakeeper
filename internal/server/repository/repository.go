package repository

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type FileRepository interface {
	GetFile(ctx context.Context, fileID string, user *model.User) (*os.File, error)
	GetFileList(ctx context.Context, user *model.User) ([]model.FileItem, error)
	DeleteFile(ctx context.Context, fileID string, user *model.User) error
	UploadFile(ctx context.Context, user *model.User, objectName string, file *os.File) error

	Save(ctx context.Context, user model.User, data model.Data) (int64, error)
	CreateContainer(ctx context.Context, user *model.User) (model.User, error)
}

type FileRepo struct {
	db       *minio.Client
	log      *logrus.Logger
	ctx      *context.Context
	location string
}

func NewFileRepository(st *minio.Client, lg *logrus.Logger, ct *context.Context) *FileRepo {
	p := &FileRepo{
		db:       st,
		log:      lg,
		ctx:      ct,
		location: "us-east-1",
	}
	return p
}

func (f *FileRepo) CreateContainer(ctx context.Context, user *model.User) (model.User, error) {
	if user.ID == 0 {
		return *user, model.ErrCreateBucketNoUser
	}

	bucketName := "bucketuid" + strconv.Itoa(int(user.ID))

	err := f.db.MakeBucket(*f.ctx, bucketName, minio.MakeBucketOptions{Region: f.location})
	if err != nil {
		exists, errBucketExists := f.db.BucketExists(*f.ctx, bucketName)
		if errBucketExists == nil && exists {
			f.log.Error("FileRepo: Bucket " + bucketName + " already exists\n")
			return *user, model.ErrCreateBucketExists
		} else {
			f.log.Log(logrus.ErrorLevel, "FileRepo: Failed to create bucket ", bucketName, ":", err, "\n")
			return *user, model.ErrCreateBucketFailed
		}
	} else {
		user.Bucket = bucketName
		f.log.Log(logrus.DebugLevel, "FileRepo: Successfully created bucket ", bucketName, "\n")
	}

	return *user, nil
}

func (f *FileRepo) GetFile(ctx context.Context, fileID string, user *model.User) (*os.File, error) {
	if user.ID == 0 {
		return nil, model.ErrCreateBucketNoUser
	}

	bucketName := "bucketuid" + strconv.Itoa(int(user.ID))

	// Create a temporary file to store the downloaded file
	tempFile, err := os.CreateTemp("", "minio_file_*.tmp")
	if err != nil {
		return nil, fmt.Errorf("failed to create temp file: %v", err)
	}

	defer func() {
		if err != nil {
			os.Remove(tempFile.Name()) // Clean up on error
		}
	}()

	// Get the object from MinIO
	object, err := f.db.GetObject(ctx, bucketName, fileID, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object from minio: %v", err)
	}
	defer object.Close()

	// Write the object to the temporary file
	if _, err = io.Copy(tempFile, object); err != nil {
		return nil, fmt.Errorf("failed to write object to temp file: %v", err)
	}

	// Ensure the file is properly closed and ready for use
	if err := tempFile.Close(); err != nil {
		return nil, fmt.Errorf("failed to close temp file: %v", err)
	}

	// Reopen the file for reading
	reopenedFile, err := os.Open(filepath.Clean(tempFile.Name()))
	if err != nil {
		return nil, fmt.Errorf("failed to reopen temp file: %v", err)
	}

	return reopenedFile, nil
}

// Операции с объектами
// defer func() {
// 	if err := client.RemoveBucket(app.Ctx, bucketName); err != nil {
// 		log.Fatalf("Failed to remove bucket: %v", err)
// 	}
// }()

func (f *FileRepo) DeleteFile(ctx context.Context, fileName string, user *model.User) error {
	bucketName := "bucketuid" + strconv.Itoa(int(user.ID))

	err := f.db.RemoveObject(ctx, bucketName, fileName, minio.RemoveObjectOptions{
		ForceDelete: true,
	})
	if err != nil {
		return fmt.Errorf("failed to delete bucket: %w", err)
	}

	return nil
}

func (f *FileRepo) GetFileList(ctx context.Context, user *model.User) ([]model.FileItem, error) {
	if user.ID <= 0 {
		return nil, model.ErrCreateBucketNoUser
	}

	bucketName := "bucketuid" + strconv.Itoa(int(user.ID))

	objectCh := f.db.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})

	var objects []model.FileItem
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			continue
		}
		objects = append(objects, model.FileItem{
			Hash: object.ETag,
			Name: object.Key,
		})
	}

	f.log.Info(objects)

	return objects, nil
}

// UploadFile uploads a file to a MinIO bucket.
func (f *FileRepo) UploadFile(ctx context.Context, user *model.User, objectName string, file *os.File) error {
	if user.ID <= 0 {
		return model.ErrNoUserBucket
	}

	bucketName := "bucketuid" + strconv.Itoa(int(user.ID))

	// Upload the file
	_, err := f.db.PutObject(ctx, bucketName, objectName, file, -1, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload file to MinIO: %w", err)
	}

	return nil
}

func (f *FileRepo) Save(ctx context.Context, user model.User, data model.Data) (int64, error) {

	// Проверка существования бакета, создание, если он не существует
	exists, err := f.db.BucketExists(context.Background(), user.Bucket)
	if err != nil {
		f.log.Log(logrus.InfoLevel, "failed to check bucket existence: ", err)
		return 0, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = f.db.MakeBucket(context.Background(), user.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return 0, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	// Преобразуем строку данных в байтовый поток
	dataReader := strings.NewReader(data.Card)

	// Загрузка данных в бакет
	_, err = f.db.PutObject(context.Background(), user.Bucket, data.Title, dataReader, int64(dataReader.Len()), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return 0, fmt.Errorf("failed to save data to MinIO: %w", err)
	}

	f.log.Printf("Successfully saved %s to bucket %s\n", data.Title, user.Bucket)

	return 1, nil
}
