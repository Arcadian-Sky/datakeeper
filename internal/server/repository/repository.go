package repository

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/minio/minio-go/v7"
	"github.com/sirupsen/logrus"
)

type FileRepository interface {
	Get(ctx context.Context, userID int64, data model.Data) (int64, error)
	GetFileList(ctx context.Context, user *model.User) ([]minio.ObjectInfo, error)
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

// Операции с объектами
// defer func() {
// 	if err := client.RemoveBucket(app.Ctx, bucketName); err != nil {
// 		log.Fatalf("Failed to remove bucket: %v", err)
// 	}
// }()

func (f *FileRepo) Get(ctx context.Context, userID int64, data model.Data) (int64, error) {
	return 1, nil
}

func (f *FileRepo) GetFileList(ctx context.Context, user *model.User) ([]minio.ObjectInfo, error) {
	if user.ID == 0 {
		return nil, model.ErrCreateBucketNoUser
	}

	bucketName := "bucketuid" + strconv.Itoa(int(user.ID))

	objectCh := f.db.ListObjects(ctx, bucketName, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})

	var objects []minio.ObjectInfo
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
			continue
		}
		objects = append(objects, object)
	}

	f.log.Info(objects)

	return objects, nil
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
	dataReader := strings.NewReader(data.PrivateData)

	// Загрузка данных в бакет
	_, err = f.db.PutObject(context.Background(), user.Bucket, data.Name, dataReader, int64(dataReader.Len()), minio.PutObjectOptions{
		ContentType: "application/octet-stream",
	})
	if err != nil {
		return 0, fmt.Errorf("failed to save data to MinIO: %w", err)
	}

	f.log.Printf("Successfully saved %s to bucket %s\n", data.Name, user.Bucket)

	return 1, nil
}
