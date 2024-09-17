package client

import (
	context "context"
	"io"

	"github.com/minio/minio-go/v7"
)

// MinioClient интерфейс для MinIO клиента
// *minio.Client
type MinioClient interface {
	ListBuckets(ctx context.Context) ([]minio.BucketInfo, error)
	MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) (err error)
	BucketExists(ctx context.Context, bucketName string) (bool, error)
	GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (MinioObject, error)
	RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error
	ListObjects(ctx context.Context, bucketName string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo
	PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,
		opts minio.PutObjectOptions,
	) (info minio.UploadInfo, err error)
}

// *minio.Object
type MinioObject interface {
	io.ReadCloser
	ReadAt(b []byte, offset int64) (n int, err error)
}

func NewMinioClient(st *minio.Client) *MinioClientWrapper {
	p := &MinioClientWrapper{
		client: st,
	}
	return p
}

type MinioClientWrapper struct {
	client *minio.Client
}

func (m *MinioClientWrapper) ListBuckets(ctx context.Context) ([]minio.BucketInfo, error) {
	return m.client.ListBuckets(ctx)
}

func (m *MinioClientWrapper) MakeBucket(ctx context.Context, bucketName string, opts minio.MakeBucketOptions) (err error) {
	return m.client.MakeBucket(ctx, bucketName, opts)
}

func (m *MinioClientWrapper) BucketExists(ctx context.Context, bucketName string) (bool, error) {
	return m.client.BucketExists(ctx, bucketName)
}

func (m *MinioClientWrapper) GetObject(ctx context.Context, bucketName, objectName string, opts minio.GetObjectOptions) (MinioObject, error) {
	return m.client.GetObject(ctx, bucketName, objectName, opts)
}

func (m *MinioClientWrapper) RemoveObject(ctx context.Context, bucketName, objectName string, opts minio.RemoveObjectOptions) error {
	return m.client.RemoveObject(ctx, bucketName, objectName, opts)
}

func (m *MinioClientWrapper) ListObjects(ctx context.Context, bucketName string, opts minio.ListObjectsOptions) <-chan minio.ObjectInfo {
	return m.client.ListObjects(ctx, bucketName, opts)
}

func (m *MinioClientWrapper) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader, objectSize int64,
	opts minio.PutObjectOptions,
) (info minio.UploadInfo, err error) {
	return m.client.PutObject(ctx, bucketName, objectName, reader, objectSize, opts)
}
