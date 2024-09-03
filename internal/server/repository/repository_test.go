package repository

import (
	"testing"
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
func TestFileRepo_GetFileList(t *testing.T) {

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
}
