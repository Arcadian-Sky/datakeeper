package client

import (
	"context"
	"fmt"
	"time"

	pbsrv "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/service/v1"
	"github.com/Arcadian-Sky/datakkeeper/internal/model"
)

func (gc *GRPCClient) GetDataList() ([]model.Data, error) {
	var data []model.Data
	if gc.Data == nil {
		return data, fmt.Errorf("GRPC client is not initialized")
	}

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	req := &pbsrv.ListDataRequest{}
	// Отправляем запрос на сервер
	res, err := gc.Data.GetDataList(ctx, req)
	if err != nil {
		gc.log.Debug("Error during get list of files : ", err)
		return data, err
	}

	gc.log.Trace(res)
	for _, item := range res.Data {
		data = append(data, model.Data{
			ID:       item.Id,
			Title:    item.Title,
			Type:     item.Type.String(),
			Login:    item.Login,
			Card:     item.Card,
			Password: item.Password,
		})
	}

	return data, nil
}

func (gc *GRPCClient) SaveLoginPass(domain, login, pass string) error {
	if gc.Data == nil {
		return fmt.Errorf("GRPC client is not initialized")
	}
	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := gc.Data.SaveData(ctx, &pbsrv.SaveDataRequest{
		Data: &pbsrv.Data{
			Type:     pbsrv.DataType_DATA_TYPE_TYPE_LOGIN_PASSWORD,
			Password: pass,
			Login:    login,
			Title:    domain,
		},
	})

	if err != nil {
		gc.log.Debug("Error during save file : ", err)
		return err
	}
	gc.log.Trace(res)

	return nil
}

func (gc *GRPCClient) SaveCard(title, card string) error {

	if gc.Data == nil {
		return fmt.Errorf("GRPC client is not initialized")
	}

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := gc.Data.SaveData(ctx, &pbsrv.SaveDataRequest{
		Data: &pbsrv.Data{
			Type:  pbsrv.DataType_DATA_TYPE_TYPE_CREDIT_CARD,
			Card:  card,
			Title: title,
		},
	})

	if err != nil {
		gc.log.Debug("Error during save file : ", err)
		return err
	}
	gc.log.Trace(res)

	return nil
}

func (gc *GRPCClient) Delete(id int64) error {
	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pbsrv.DeleteDataRequest{
		Dataid: id,
	}

	// Отправляем запрос на сервер
	res, err := gc.Data.DeleteData(ctx, req)
	if err != nil {
		gc.log.Debug("Error during get list of files : ", err)
		return err
	}

	gc.log.Trace(res)

	return nil
}
