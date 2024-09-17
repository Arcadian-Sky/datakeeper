package client

import (
	"context"
	"errors"
	"fmt"
	"time"

	pb "github.com/Arcadian-Sky/datakkeeper/gen/proto/api/user/v1"
)

// Регистрация нового пользователя.
func (gc *GRPCClient) Register(login, password string) error {

	if gc.User == nil {
		return fmt.Errorf("GRPC client is not initialized")
	}
	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Формируем запрос на регистрацию
	req := &pb.RegisterRequest{
		Login:    login,
		Password: password,
	}

	// Отправляем запрос на сервер
	res, err := gc.User.Register(ctx, req)
	if err != nil {
		gc.log.Debug("Error during registration:", err)
		return err
	}
	gc.Storage.SetToken(res.AuthToken)

	// Обрабатываем ответ сервера
	if res.Success {
		gc.log.Info("Registration successful, auth token:", res.AuthToken)
	} else {
		gc.log.Info("Registration failed:", res.Message)
	}

	return nil
}

// Аутентификация пользователя.
func (gc *GRPCClient) Authenticate(login, password string) error {
	if gc.User == nil {
		return fmt.Errorf("GRPC client is not initialized")
	}

	// Создаем контекст с таймаутом для запроса
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// Формируем запрос на регистрацию
	req := &pb.AuthenticateRequest{
		Login:    login,
		Password: password,
	}

	// Отправляем запрос на сервер
	res, err := gc.User.Authenticate(ctx, req)
	if err != nil {
		gc.log.Debug("Error during authrntificate: ", err)
		return err
	}
	gc.Storage.SetToken(res.AuthToken)

	// Обрабатываем ответ сервера
	if res.Success {
		gc.log.Info("Authrntificate successful, auth login: ", gc.Storage.Login)
		gc.log.Info("Authrntificate successful, auth token: ", res.AuthToken)
	} else {
		gc.log.Info("Authrntificate failed: ", res.Message)
		return errors.New("Authrntificate failed: " + res.Message)
	}

	return nil
}
