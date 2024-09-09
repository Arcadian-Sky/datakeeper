package client

import (
	"errors"
	"fmt"
	"testing"

	"github.com/Arcadian-Sky/datakkeeper/internal/client"
	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"github.com/Arcadian-Sky/datakkeeper/mocks"
	"github.com/golang/mock/gomock"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestInitDataInterfaces(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient

	app.initDataInterfaces()

	// Проверяем, что страницы были добавлены
	assert.NotNil(t, app.data.list)
	assert.NotNil(t, app.data.loadForm)
	assert.NotNil(t, app.data.sendLoginPassForm)
	assert.NotNil(t, app.data.sendCardForm)
}

func TestIniSettings(t *testing.T) {
	app := NewEmptyApp()
	app.settings = &settings.ClientConfig{
		ServerAddress: "1111",
	}

	app.iniSettings()

	// Проверяем, что форма настроек была инициализирована
	assert.NotNil(t, app.settingsForm)
}

func TestInitPages(t *testing.T) {
	app := NewEmptyApp()

	app.settings = &settings.ClientConfig{
		ServerAddress: "1111",
	}
	app.log = logrus.New()
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
	app := NewEmptyApp()
	app.client = mockClient

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
	app := NewEmptyApp()
	app.client = mockClient

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
	app := NewEmptyApp()

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
	app := NewEmptyApp()
	app.client = mockClient

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
	app := NewEmptyApp()
	app.client = mockClient

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
	app := NewEmptyApp()
	app.log = logrus.New()

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
	app := NewEmptyApp()
	app.log = logrus.New()

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
	app := NewEmptyApp()

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
	app := NewEmptyApp()
	app.client = mockClient

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

func TestApp_AddAction(t *testing.T) {
	// Create a new instance of App
	app := &App{}

	// Create a new tview.Form and FormRegister (which is a map)
	form := tview.NewForm()
	register := &FormRegister{}

	// Call addAction to add a button and register the action
	title := "Submit"
	app.addAction(form, register, title, func() {})

	// Assert that the action has been added to the register
	assert.Contains(t, *register, title, "register should contain the added action title")
	assert.Equal(t, (*register)[title], title, "The title in the register should match")
}

func TestApp_acgtionSaveRegisterForm_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()

	// Setup the form fields
	app.person.registerForm.AddInputField("Login", "testuser", 20, nil, nil)
	app.person.registerForm.AddInputField("Password", "password", 20, nil, nil)

	// Mock the Register method to succeed
	mockClient.EXPECT().Register("testuser", "password").Return(nil)

	// Call the method
	app.actionSaveRegisterForm()

	// Assert the storage is updated and page is switched to "main"
	assert.Equal(t, "testuser", app.storage.Login)
	// assert.True(t, app.pages.HasPage("main"))
}

func TestApp_acgtionSaveRegisterForm_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()

	// Setup the form fields
	app.person.registerForm.AddInputField("Login", "testuser", 20, nil, nil)
	app.person.registerForm.AddInputField("Password", "password", 20, nil, nil)

	// Expect the Register method to be called and return an error
	mockClient.EXPECT().Register("testuser", "password").Return(errors.New("registration error"))

	// Call the method
	app.actionSaveRegisterForm()

	// Assert storage is not updated and page is not switched to "main"
	assert.Empty(t, app.storage.Login)
	assert.False(t, app.pages.HasPage("main"))
}

func TestApp_actionAuth_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()

	// Setup the form fields (login and password)
	app.person.authForm.AddInputField("Login", "testuser", 20, nil, nil)
	app.person.authForm.AddInputField("Password", "password", 20, nil, nil)

	// Expect the Authenticate method to be called with "testuser" and "password"
	mockClient.EXPECT().Authenticate("testuser", "password").Return(nil)

	// Call the method
	app.actionAuth()

	// Assert the storage is updated and page is switched to "main"
	assert.Equal(t, "testuser", app.storage.Login)
}

func TestApp_actionAuth_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()

	// Setup the form fields (login and password)
	app.person.authForm.AddInputField("Login", "testuser", 20, nil, nil)
	app.person.authForm.AddInputField("Password", "password", 20, nil, nil)

	// Expect the Authenticate method to be called with "testuser" and "password" and return an error
	mockClient.EXPECT().Authenticate("testuser", "password").Return(errors.New("authentication error"))

	// Call the method
	app.actionAuth()

	// Assert storage is not updated and page is not switched to "main"
	assert.Empty(t, app.storage.Login)
}

func TestApp_actionSwitchToAuth(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()

	// Add "auth" page to pages
	app.pages.AddPage("auth", tview.NewTextView(), true, false)

	// Call the method
	app.actionSwitchToAuth()

	// Check if the log view is cleared
	assert.Equal(t, "", app.logView.GetText(true))

	// Check if the page is switched to "auth"
	assert.True(t, app.pages.HasPage("auth"))
}

// type MockApp struct {
// 	*App
// 	mockCtrl *gomock.Controller
// }

// func NewMockApp(t *testing.T) *MockApp {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()

// 	app := NewEmptyApp()
// 	app.client = mocks.NewMockGRPCClientInterface(ctrl)
// 	app.storage = client.NewMemStorage()
// 	app.log = logrus.New()

// 	return &MockApp{
// 		App:      app,
// 		mockCtrl: ctrl,
// 	}
// }

func TestApp_actionSwitchToRegister(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	app := NewEmptyApp()
	app.client = mocks.NewMockGRPCClientInterface(ctrl)
	app.storage = client.NewMemStorage()
	app.log = logrus.New()

	// Add "register" page to pages
	app.pages.AddPage("register", tview.NewTextView(), true, false)

	// Call the method
	app.actionSwitchToRegister()

	// Check if the log view is cleared
	assert.Equal(t, "", app.logView.GetText(true))

	// Check if the page is switched to "register"
	assert.True(t, app.pages.HasPage("register"))
}

func TestApp_appActionSendFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()

	// Set up the form and input field
	form := tview.NewForm()
	inputField := tview.NewInputField().SetText("/path/to/file")
	form.AddFormItem(inputField)
	app.data.loadForm = form

	// Define the behavior of the mock client
	mockClient.EXPECT().UploadFile("/path/to/file").Return(nil).Times(1)

	// Call the method
	app.appActionSendFiles()

	// Check if the log view is cleared
	assert.Equal(t, "", app.logView.GetText(true))
}

func TestApp_appActionSendLoginPass(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()

	// Set up the form and input fields
	form := tview.NewForm()
	form.AddFormItem(tview.NewInputField().SetText("example.com"))
	form.AddFormItem(tview.NewInputField().SetText("user"))
	form.AddFormItem(tview.NewInputField().SetText("password"))
	app.data.sendLoginPassForm = form

	// Define the behavior of the mock client
	mockClient.EXPECT().SaveLoginPass("example.com", "user", "password").Return(nil).Times(1)

	// Call the method
	app.appActionSendLoginPass()

	// Check if the log view is cleared
	assert.Equal(t, "", app.logView.GetText(true))
}

func TestApp_appActionSendCard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()

	// Set up the form and input fields
	form := tview.NewForm()
	form.AddFormItem(tview.NewInputField().SetText("domain"))
	form.AddFormItem(tview.NewInputField().SetText("cardnumber"))
	app.data.sendCardForm = form

	// Define the behavior of the mock client
	mockClient.EXPECT().SaveCard("domain", "cardnumber").Return(nil).Times(1)

	// Call the method
	app.appActionSendCard()

	// Check if the log view is cleared
	assert.Equal(t, "", app.logView.GetText(true))
}

func TestApp_checkInputCardField(t *testing.T) {
	app := &App{}

	tests := []struct {
		textToCheck string
		lastChar    rune
		expected    bool
	}{
		{"1234 5678", ' ', true},
		{"1234 5678", '5', true},
		{"1234 5678", 'a', false},
		{"1234 5678", '\n', false},
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("%s-%c", tt.textToCheck, tt.lastChar), func(t *testing.T) {
			result := app.checkInputCardField(tt.textToCheck, tt.lastChar)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// TestApp_actionDeleteData tests the actionDeleteData method of the App struct
func TestApp_actionDeleteData(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()
	app.log.SetOutput(app.logView)

	// Define the ID to be deleted
	idToDelete := int64(123)

	// Test for successful deletion
	t.Run("Success", func(t *testing.T) {
		// Setup the mock to return no error
		mockClient.EXPECT().Delete(idToDelete).Return(nil).Times(1)

		// Create and call the action
		action := app.appActionDeleteData(idToDelete)
		action()

		logLines := app.logView.GetText(true)
		assert.Contains(t, logLines, "Delete pressed ID: 123")
	})

	// Test for failed deletion
	t.Run("Failure", func(t *testing.T) {
		// Setup the mock to return an error
		mockClient.EXPECT().Delete(idToDelete).Return(fmt.Errorf("delete error")).Times(1)

		// Create and call the action
		action := app.appActionDeleteData(idToDelete)
		action()

		logLines := app.logView.GetText(true)
		assert.Contains(t, logLines, "Delete pressed ID: 123")
		assert.Contains(t, logLines, "Error delete item: delete error")
	})
}

// TestApp_appActionGetFiles tests the appActionGetFiles method of the App struct
func TestApp_appActionGetFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()
	app.log.SetOutput(app.logView)

	name := "testfile"
	id := "123"

	// Test for successful file retrieval
	t.Run("Success", func(t *testing.T) {
		// Setup the mock to return no error
		mockClient.EXPECT().GetFile(name).Return(nil).Times(1)

		// Create and call the action
		action := app.appActionGetFiles(name, id)
		action()

		assert.Contains(t, app.logView.GetText(true), "Getting started ID: 123")
		assert.Contains(t, app.logView.GetText(true), "Got ID: 123")
	})

	// Test for failed file retrieval
	t.Run("Failure", func(t *testing.T) {
		// Setup the mock to return an error
		mockClient.EXPECT().GetFile(name).Return(fmt.Errorf("get file error")).Times(1)

		// Create and call the action
		action := app.appActionGetFiles(name, id)
		action()

		logLines := app.logView.GetText(true)
		assert.Contains(t, logLines, "Getting started ID: 123")
		assert.Contains(t, logLines, "Error client GetFile: get file error")
	})
}

// TestApp_appActionDeleteFiles tests the appActionDeleteFiles method of the App struct
func TestApp_appActionDeleteFiles(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()
	app.log.SetOutput(app.logView)

	name := "testfile"
	id := "123"

	// Test for successful file deletion
	t.Run("Success", func(t *testing.T) {
		// Setup the mock to return no error
		mockClient.EXPECT().DeleteFile(name).Return(nil).Times(1)

		// Create and call the action
		action := app.appActionDeleteFiles(name, id)
		action()

		// Check if the log view is cleared
		assert.Contains(t, app.logView.GetText(true), "Deleted ID: 123")
	})

	// Test for failed file deletion
	t.Run("Failure", func(t *testing.T) {
		// Setup the mock to return an error
		mockClient.EXPECT().DeleteFile(name).Return(fmt.Errorf("delete file error")).Times(1)

		// Create and call the action
		action := app.appActionDeleteFiles(name, id)
		action()

		// Verify the log message
		logLines := app.logView.GetText(true)
		assert.Contains(t, logLines, "Error client DeleteFile: delete file error")
	})
}

func TestApp_actionSwitchToAuth_Log(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()
	app.log.SetOutput(app.logView)
	app.log.SetLevel(logrus.TraceLevel)

	app.actionSwitchToAuth()

	logLines := app.logView.GetText(true)
	assert.Contains(t, logLines, "SwitchToPage auth")
}

func TestApp_actionSwitchToRegister_Log(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()
	app.log.SetOutput(app.logView)
	app.log.SetLevel(logrus.TraceLevel)

	app.actionSwitchToRegister()

	logLines := app.logView.GetText(true)
	assert.Contains(t, logLines, "SwitchToPage register")
}

func TestApp_actionSwitchToMain_Log(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()
	app.log.SetOutput(app.logView)
	app.log.SetLevel(logrus.TraceLevel)

	app.actionSwitchToMain()

	logLines := app.logView.GetText(true)
	assert.Contains(t, logLines, "SwitchToPage main")
}

func TestApp_actionSwitchToDataList_Log(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()
	app.log.SetOutput(app.logView)
	app.log.SetLevel(logrus.TraceLevel)

	app.actionSwitchToDataList()

	logLines := app.logView.GetText(true)
	assert.Contains(t, logLines, "SwitchToPage datalist")
}

func TestApp_actionSwitchToMainWithClear_Log(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаем mock клиента
	mockClient := mocks.NewMockGRPCClientInterface(ctrl)

	app := NewEmptyApp()
	app.client = mockClient
	app.storage = client.NewMemStorage()
	app.log = logrus.New()
	app.log.SetOutput(app.logView)
	app.log.SetLevel(logrus.TraceLevel)

	app.actionSwitchToMainWithClear()

	logLines := app.logView.GetText(true)
	assert.Contains(t, logLines, "SwitchToPage main with clear")
}
