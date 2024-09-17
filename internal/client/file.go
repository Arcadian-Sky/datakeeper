package client

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	pbsrv "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/service/v1"
	"github.com/Arcadian-Sky/datakkeeper/internal/model"
)

// Получение списка файлов
func (gc *GRPCClient) GetFileList() ([]model.FileItem, error) {
	var data []model.FileItem
	if gc.Data == nil {
		return data, fmt.Errorf("GRPC client is not initialized")
	}

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pbsrv.ListFileRequest{}
	// Отправляем запрос на сервер
	res, err := gc.Data.GetFileList(ctx, req)
	if err != nil {
		gc.log.Debug("Error during get list of files : ", err)
		return data, err
	}
	gc.log.Trace(res)
	for _, item := range res.Fileitem {
		data = append(data, model.FileItem{
			Hash: item.Key,
			Name: item.Name,
		})
	}

	return data, nil
}

// Удаление файла
func (gc *GRPCClient) DeleteFile(fileName string) error {
	if gc.Data == nil {
		return fmt.Errorf("GRPC client is not initialized")
	}

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pbsrv.DeleteFileRequest{
		Filename: fileName,
	}
	// Отправляем запрос на сервер
	res, err := gc.Data.DeleteFile(ctx, req)
	if err != nil {
		gc.log.Debug("Error during delete file : ", err)
		return err
	}
	gc.log.Trace(res)

	return nil
}

// Отправка файлов на сервер
func (gc *GRPCClient) UploadFile(filePath string) error {
	fileName := filepath.Base(filePath)
	gc.log.Info("File name: ", fileName)
	// Открываем файл для чтения
	file, err := os.Open(filePath)
	if err != nil {
		gc.log.Trace("could not open file: ", err)
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	stream, err := gc.Data.UploadFile(context.Background())
	if err != nil {
		gc.log.Trace("error creating stream: ", err)
		return fmt.Errorf("error creating stream: %v", err)
	}

	// Read the file and send chunks
	buffer := make([]byte, 1024)
	for {
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}

		if err != nil {
			gc.log.Trace("Failed to read file: ", err)
		}

		if err := stream.Send(&pbsrv.FileChunk{
			Data:     buffer[:n],
			Filename: fileName,
		}); err != nil {
			gc.log.Trace("Failed to send chunk: ", err)
		}
	}
	// Close the stream and get the response
	status, err := stream.CloseAndRecv()
	if err != nil {
		gc.log.Trace("Failed to receive response: ", err)
	}

	gc.log.Info("Upload status:", status.Success, ", message: ", status.Message)

	return nil
}

func (gc *GRPCClient) GetFile(fileName string) error {

	stream, err := gc.Data.GetFile(context.Background(), &pbsrv.GetFileRequest{Name: fileName})
	if err != nil {
		gc.log.Fatalf("Ошибка при вызове GetFile: %v", err)
	}

	filePath := filepath.Join(gc.Storage.PfilesDir, fileName)
	// Создаём файл для записи полученных данных
	file, err := os.Create(filePath)
	if err != nil {
		gc.log.Trace("Не удалось создать файл: ", err)
		return err
	}
	defer file.Close()

	// Читаем поток данных и записываем в файл
	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			gc.log.Trace("Ошибка при получении данных: ", err)
			return err
		}
		_, err = file.Write(chunk.Data)
		if err != nil {
			gc.log.Trace("Не удалось записать в файл: ", err)
			return err
		}
	}
	gc.log.Info("Файл успешно получен и сохранён:", filePath)
	return nil
}
