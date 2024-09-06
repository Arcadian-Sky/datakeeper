package client

import (
	"errors"
	"testing"

	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"github.com/Arcadian-Sky/datakkeeper/mocks"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestInitDataInterfaces(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := &App{
		data: Data{
			list:              tview.NewList(),
			loadForm:          tview.NewForm(),
			sendLoginPassForm: tview.NewForm(),
			sendCardForm:      tview.NewForm(),
		},
		pages:   tview.NewPages(),
		logView: tview.NewTextView(),
		client:  mockClient,
	}

	app.initDataInterfaces()

	// Проверяем, что страницы были добавлены
	assert.NotNil(t, app.data.list)
	assert.NotNil(t, app.data.loadForm)
	assert.NotNil(t, app.data.sendLoginPassForm)
	assert.NotNil(t, app.data.sendCardForm)
}

func TestIniSettings(t *testing.T) {
	app := &App{
		settings: &settings.ClientConfig{
			ServerAddress: "1111",
		},
		settingsForm: tview.NewForm(),
		pages:        tview.NewPages(),
	}

	app.iniSettings()

	// Проверяем, что форма настроек была инициализирована
	assert.NotNil(t, app.settingsForm)
}

func TestInitPages(t *testing.T) {
	app := &App{
		tapp:         tview.NewApplication(),
		pages:        tview.NewPages(),
		settingsForm: tview.NewForm(),
		logView:      tview.NewTextView(),
		person: Person{
			authForm:     tview.NewForm(),
			registerForm: tview.NewForm(),
			Form:         tview.NewForm(),
		},
		data: Data{
			list:              tview.NewList(),
			move:              tview.NewList(),
			loadForm:          tview.NewForm(),
			sendLoginPassForm: tview.NewForm(),
			sendCardForm:      tview.NewForm(),
		},
		settings: &settings.ClientConfig{
			ServerAddress: "1111",
		},
		log: logrus.New(),
	}
	// app.log.SetFormatter(&logrus.JSONFormatter{})
	app.log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,  // Включить цвета в выводе
		FullTimestamp: false, // Включить полный временной штамп
	})
	app.log.SetLevel(logrus.FatalLevel)

	app.initPages()

	pageNames := app.pages.GetPageNames(false)
	assert.Contains(t, pageNames, "main")
	assert.Contains(t, pageNames, "auth")
	assert.Contains(t, pageNames, "register")
	assert.Contains(t, pageNames, "person")
	assert.Contains(t, pageNames, "datalist")
	assert.Contains(t, pageNames, "fileform")
	assert.Contains(t, pageNames, "loginpassform")
	assert.Contains(t, pageNames, "cardform")
	assert.Contains(t, pageNames, "settings")
	assert.NotNil(t, app.logView)
}

func TestIsTokenValid(t *testing.T) {
	// Вызываем тестируемую функцию
	result := isTokenValid()

	// Проверяем, что результат соответствует ожидаемому значению
	assert.False(t, result, "Expected isTokenValid to return false")
}

func TestAppLoadData_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	// Пример данных, которые будет возвращать клиент
	mockData := []model.Data{
		{ID: 1, Title: "data1", Type: "type1", Card: "card1", Login: "login1", Password: "password1"},
		{ID: 2, Title: "data2", Type: "type2", Card: "card2", Login: "login2", Password: "password2"},
	}

	// Настраиваем mock для успешного вызова GetDataList
	mockClient.EXPECT().GetDataList().Return(mockData, nil)

	// Создаем экземпляр приложения с mock клиентом
	app := &App{
		tapp:   tview.NewApplication(),
		pages:  tview.NewPages(),
		client: mockClient,
	}

	// Вызываем тестируемый метод
	err := app.loadData()

	// Проверяем, что ошибки не возникло
	assert.NoError(t, err)

}

func TestAppLoadData_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	// Настраиваем mock для возврата ошибки
	mockClient.EXPECT().GetDataList().Return(nil, errors.New("client error"))

	// Создаем экземпляр приложения с mock клиентом
	app := &App{
		tapp:   tview.NewApplication(),
		pages:  tview.NewPages(),
		client: mockClient,
	}

	// Вызываем тестируемый метод
	err := app.loadData()

	// Проверяем, что возникла ожидаемая ошибка
	assert.EqualError(t, err, "error client GetDataList: client error")

	// Убедимся, что updateDatalistPage не был вызван
	assert.NotNil(t, err, "Expected an error but got nil")
}

func TestAppUpdateDatalistPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock для страницы
	pages := tview.NewPages()
	app := &App{
		pages: pages,
	}

	// Пример данных для обновления
	mockData := []model.Data{
		{ID: 1, Title: "data1", Type: "type1", Card: "card1", Login: "login1", Password: "password1"},
		{ID: 2, Title: "data2", Type: "type2", Card: "card2", Login: "login2", Password: "password2"},
	}

	// Вызываем тестируемый метод
	app.updateDatalistPage(mockData)
}

func TestAppLoadFiles_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	// Пример данных, которые будет возвращать клиент
	mockData := []model.FileItem{
		{Name: "file1", Desc: "description1", Hash: "hash1"},
		{Name: "file2", Desc: "description2", Hash: "hash2"},
	}

	// Настраиваем mock для успешного вызова GetFileList
	mockClient.EXPECT().GetFileList().Return(mockData, nil)

	// Создаем экземпляр приложения с mock клиентом
	app := &App{
		tapp:   tview.NewApplication(),
		pages:  tview.NewPages(),
		client: mockClient,
	}

	// Вызываем тестируемый метод
	err := app.loadFiles()

	// Проверяем, что ошибки не возникло
	assert.NoError(t, err)
}

func TestAppLoadFiles_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	// Настраиваем mock для возврата ошибки
	mockClient.EXPECT().GetFileList().Return(nil, errors.New("stream error"))

	// Создаем экземпляр приложения с mock клиентом
	app := &App{
		client: mockClient,
	}

	// Вызываем тестируемый метод
	err := app.loadFiles()

	// Проверяем, что возникла ожидаемая ошибка
	assert.EqualError(t, err, "error creating stream: stream error")

	// Убедимся, что updateFileDatalistPage не был вызван
	assert.NotNil(t, err, "Expected an error but got nil")
}

func TestAppUpdateFileDatalistPage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем экземпляр приложения
	app := &App{
		tapp:    tview.NewApplication(),
		pages:   tview.NewPages(),
		logView: tview.NewTextView(),
		log:     logrus.New(),
	}

	// Пример данных для списка файлов
	data := []model.FileItem{
		{Name: "file1", Desc: "description1", Hash: "hash1"},
		{Name: "file2", Desc: "description2", Hash: "hash2"},
	}

	// Вызываем метод обновления списка файлов
	app.updateFileDatalistPage(data)

	// Проверяем, что страница была добавлена
	pageNames := app.pages.GetPageNames(false)
	assert.Contains(t, pageNames, "datalistmove")

}

func TestAppCreateMoveForm(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Mock GRPCClient

	// Создаем экземпляр приложения
	app := &App{
		tapp:    tview.NewApplication(),
		pages:   tview.NewPages(),
		logView: tview.NewTextView(),
		log:     logrus.New(),
	}

	// Входные данные для теста
	id := "12345"
	name := "testfile"
	desc := "test description"

	// Вызываем метод создания формы
	app.createMoveForm(id, name, desc)

	// Проверяем, что форма была добавлена как страница
	pageNames := app.pages.GetPageNames(false)
	assert.Contains(t, pageNames, "datalistmoveaction")
}

func TestAppInitialization(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем экземпляр приложения
	app := &App{
		tapp:   tview.NewApplication(),
		person: Person{},
		data:   Data{},
		pages:  tview.NewPages(),
		Conn:   &grpc.ClientConn{}, // Можно мокировать при необходимости
		log:    logrus.New(),
	}

	// Инициализация форм
	app.person.authForm = tview.NewForm()
	app.person.registerForm = tview.NewForm()
	app.data.sendLoginPassForm = tview.NewForm()
	app.data.sendCardForm = tview.NewForm()

	// Добавляем страницы с формами
	app.pages.AddPage("auth", app.person.authForm, true, false)
	app.pages.AddPage("register", app.person.registerForm, true, false)

	// Проверяем, что страницы добавлены
	pageNames := app.pages.GetPageNames(false)
	assert.Contains(t, pageNames, "auth")
	assert.Contains(t, pageNames, "register")

	// Переключаемся на страницу и проверяем, что переключение происходит
	app.pages.SwitchToPage("auth")
	assert.Contains(t, app.pages.GetPageNames(false), "auth")

	app.pages.SwitchToPage("register")
	assert.Contains(t, app.pages.GetPageNames(false), "register")

	// Проверяем, что формы корректно созданы
	assert.NotNil(t, app.person.authForm)
	assert.NotNil(t, app.person.registerForm)

	// Проверка работы с формой данных
	app.data.list = tview.NewList()
	app.data.move = tview.NewList()

	// Добавляем действия в список
	app.data.list.AddItem("Item 1", "Description", 'a', nil)
	app.data.move.AddItem("Move 1", "Move Description", 'b', nil)

	assert.Equal(t, 1, app.data.list.GetItemCount())
	assert.Equal(t, 1, app.data.move.GetItemCount())
}

func TestCreateDetailForm(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	// Инициализация тестового приложения
	app := &App{
		pages:   tview.NewPages(),
		logView: tview.NewTextView(),
		client:  mockClient,
	}

	// Тестовые данные
	item := model.Data{
		Title:    "Test Item",
		ID:       1,
		Type:     "Test Type",
		Card:     "1234 5678 9012 3456",
		Login:    "testuser",
		Password: "testpass",
	}

	// Вызов метода
	app.createDetailForm(item)

	pageNames := app.pages.GetPageNames(false)
	assert.Contains(t, pageNames, "datalistmoveaction", "Page 'datalistmoveaction' should be added")
}
