package client

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Arcadian-Sky/datakkeeper/internal/client"
	"github.com/Arcadian-Sky/datakkeeper/internal/model"
	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"github.com/rivo/tview"
	"github.com/sirupsen/logrus"
	"github.com/sqweek/dialog"
	"google.golang.org/grpc"
)

type FormRegister map[string]string
type Person struct {
	authForm            *tview.Form
	authFormButtons     *FormRegister
	registerForm        *tview.Form
	registerFormButtons *FormRegister
	Form                *tview.Form
}

type Data struct {
	loadForm                 *tview.Form
	loadFormButtons          *FormRegister
	sendLoginPassForm        *tview.Form
	sendLoginPassFormButtons *FormRegister
	sendCardForm             *tview.Form
	sendCardFormButtons      *FormRegister
	list                     *tview.List
	move                     *tview.List
}
type App struct {
	settings            *settings.ClientConfig
	tapp                *tview.Application
	person              Person
	data                Data
	pages               *tview.Pages
	settingsForm        *tview.Form
	settingsFormButtons *FormRegister
	logView             *tview.TextView
	client              client.GRPCClientInterface
	Conn                *grpc.ClientConn
	storage             *client.MemStorage
	log                 *logrus.Logger
}

func NewEmptyApp() *App {
	return &App{
		tapp:                tview.NewApplication(),
		pages:               tview.NewPages(),
		settingsForm:        tview.NewForm(),
		settingsFormButtons: &FormRegister{},
		logView:             tview.NewTextView(),
		person: Person{
			authForm:            tview.NewForm(),
			authFormButtons:     &FormRegister{},
			registerForm:        tview.NewForm(),
			registerFormButtons: &FormRegister{},
			Form:                tview.NewForm(),
		},
		data: Data{
			list:                     tview.NewList(),
			move:                     tview.NewList(),
			loadForm:                 tview.NewForm(),
			loadFormButtons:          &FormRegister{},
			sendLoginPassForm:        tview.NewForm(),
			sendLoginPassFormButtons: &FormRegister{},
			sendCardForm:             tview.NewForm(),
			sendCardFormButtons:      &FormRegister{},
		},
	}
}

func NewClientApp(st *settings.ClientConfig) *App {
	app := NewEmptyApp()
	app.settings = st
	app.log = logrus.New()
	// app.log.SetFormatter(&logrus.JSONFormatter{})
	app.log.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,  // Включить цвета в выводе
		FullTimestamp: false, // Включить полный временной штамп
	})
	app.log.SetLevel(logrus.TraceLevel)
	// app.log.SetReportCaller(true)

	app.storage = client.NewMemStorage()
	app.client, app.Conn = client.NewGclient(*app.settings, app.storage, app.log)

	app.initPersonInterfaces()
	app.initDataInterfaces()
	app.iniSettings()
	app.initPages()
	// Запуск
	if err := app.tapp.Run(); err != nil {
		panic(err)
	}
	return app
}

// Инициализация интерфейсов авторизации и регистрации
func (app *App) initPersonInterfaces() {
	// Создаем форму для авторизации
	app.person.authForm.SetBorder(true).SetTitle("Authorize").SetTitleAlign(tview.AlignLeft)
	app.person.authForm.
		AddInputField("Login", "", 20, nil, nil).
		AddPasswordField("Password", "", 10, '*', nil)

	app.addAction(app.person.authForm, app.person.authFormButtons, "Save", app.actionAuth)
	app.addAction(app.person.authForm, app.person.authFormButtons, "Switch to Register", app.actionSwitchToRegister)
	app.addAction(app.person.authForm, app.person.authFormButtons, "Quit", app.appActionQuit)

	// Создаем форму для регистрации
	app.person.registerForm.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
	app.person.registerForm.
		AddInputField("Login", "", 20, nil, nil).
		AddPasswordField("Password", "", 10, '*', nil)

	app.addAction(app.person.registerForm, app.person.registerFormButtons, "Save", app.actionSaveRegisterForm)
	app.addAction(app.person.registerForm, app.person.registerFormButtons, "Switch to Authorize", app.actionSwitchToAuth)
	app.addAction(app.person.registerForm, app.person.registerFormButtons, "Quit", app.appActionQuit)

	// Создаем формы для авторизации и регистрации
	app.person.Form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
}

func (app *App) actionAuth() {
	login := app.person.authForm.GetFormItem(0).(*tview.InputField).GetText()
	app.storage.Login = ""
	password := app.person.authForm.GetFormItem(1).(*tview.InputField).GetText()
	app.logView.Clear()
	err := app.client.Authenticate(login, password)
	if err != nil {
		app.log.Info("Error client Authentificate: ", err)
		return
	}
	app.storage.Login = login
	app.actionSwitchToMain()
}

func (app *App) actionSaveRegisterForm() {
	login := app.person.registerForm.GetFormItem(0).(*tview.InputField).GetText()
	password := app.person.registerForm.GetFormItem(1).(*tview.InputField).GetText()
	app.storage.Login = ""
	app.logView.Clear()
	err := app.client.Register(login, password)
	if err != nil {
		app.log.Info("Error client Register: ", err)
		return
	}
	app.storage.Login = login
	app.actionSwitchToMain()
}

func (app *App) actionSwitchToAuth() {
	app.logView.Clear()
	app.pages.SwitchToPage("auth")
	app.log.Trace("SwitchToPage auth")
}

func (app *App) actionSwitchToRegister() {
	app.logView.Clear()
	app.pages.SwitchToPage("register")
	app.log.Trace("SwitchToPage register")
}

func (app *App) actionSwitchToMain() {
	app.pages.SwitchToPage("main")
	app.log.Trace("SwitchToPage main")
}

func (app *App) actionSwitchToMainWithClear() {
	app.logView.Clear()
	app.pages.SwitchToPage("main")
	app.log.Trace("SwitchToPage main with clear")
}

func (app *App) actionSwitchToDataListWithClear() {
	app.logView.Clear()
	app.pages.SwitchToPage("datalist")
	app.log.Trace("SwitchToPage datalist with clear")
}

func (app *App) actionSwitchToDataList() {
	app.pages.SwitchToPage("datalist")
	app.log.Trace("SwitchToPage datalist")
}

func (app *App) actionSwitchToFileForm() {
	app.pages.SwitchToPage("fileform")
	app.log.Trace("SwitchToPage fileform")
}

func (app *App) actionSwitchToPerson() {
	app.pages.SwitchToPage("person")
	app.log.Trace("SwitchToPage person")
}

func (app *App) actionSwitchToLogpassForm() {
	app.pages.SwitchToPage("loginpassform")
	app.log.Trace("SwitchToPage loginpassform")
}

func (app *App) actionSwitchToCardForm() {
	app.pages.SwitchToPage("cardform")
	app.log.Trace("SwitchToPage cardform")
}

func (app *App) actionSwitchToSettings() {
	app.pages.SwitchToPage("settings")
	app.log.Trace("SwitchToPage settings")
}

func (app *App) appActionQuit() {
	app.tapp.Stop()
	app.log.Trace("Switch stop app")
}

func (app *App) appActionLoadFiles() {
	app.logView.Clear()
	err := app.loadFiles()
	if err != nil {
		app.log.Info("Error loading data: ", err)
	} else {
		app.log.Trace("ActionLoadFiles: Files loaded")
	}
}

func (app *App) appActionLoadData() {
	app.logView.Clear()
	err := app.loadData()
	if err != nil {
		app.log.Info("Error loading data: ", err)
	} else {
		app.log.Trace("ActionLoadData: Data loaded")
	}
}

func (app *App) appActionSelectFiles(loadForm *tview.Form) func() {
	return func() {
		filename, err := dialog.File().Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting file: %v\n", err)
			return
		}

		// Устанавливаем выбранный файл в InputField
		loadForm.GetFormItem(0).(*tview.InputField).SetText(filename)
	}
}

func (app *App) appActionSendFiles(loadForm *tview.Form) func() {
	return func() {
		app.logView.Clear()
		app.log.Info("Sending data...")
		filePath := loadForm.GetFormItem(0).(*tview.InputField).GetText()
		app.log.Info("File path: ", filePath)
		if err := app.client.UploadFile(filePath); err != nil {
			app.log.Info(fmt.Println("Error uploading file:", err))
		}
	}
}

func (app *App) appActionSendLoginPass() {
	domain := app.data.sendLoginPassForm.GetFormItem(0).(*tview.InputField).GetText()
	login := app.data.sendLoginPassForm.GetFormItem(1).(*tview.InputField).GetText()
	password := app.data.sendLoginPassForm.GetFormItem(2).(*tview.InputField).GetText()

	app.log.Info(fmt.Printf("Login: %s, Password: %s\n, Domain: %s\n", login, "*****", domain))

	if err := app.client.SaveLoginPass(domain, login, password); err != nil {
		app.log.Info(fmt.Println("Error saving login pass:", err))
	}
}

func (app *App) appActionSendCard() {
	domain := app.data.sendCardForm.GetFormItem(0).(*tview.InputField).GetText()
	card := app.data.sendCardForm.GetFormItem(1).(*tview.InputField).GetText()

	app.log.Info(fmt.Printf("Title: %s, Card: %s\n", domain, card))
	if err := app.client.SaveCard(domain, card); err != nil {
		app.log.Info(fmt.Println("Error saving card:", err))
	}
}

func (app *App) checkInputCardField(textToCheck string, lastChar rune) bool {
	// Allow only numbers and spaces in the card number field
	return (lastChar >= '0' && lastChar <= '9') || lastChar == ' '
}

// Инициализация интерфейсов работы с данными
func (app *App) initDataInterfaces() {
	// Создаем интерфейс для отображения данных
	app.data.list.SetTitle("Data List")
	app.data.list.AddItem("Back", "Make click to run action", 'q', app.actionSwitchToMain)
	app.data.list.AddItem("Load Files", "", 'f', app.appActionLoadFiles)
	app.data.list.AddItem("Load Data", "", 'd', app.appActionLoadData)

	// Создаем форму для отправки данных
	app.data.loadForm.
		SetTitle("Send File Entry").
		SetBorder(true)

	app.data.loadForm.
		AddInputField("File Path", "", 40, nil, nil)

	app.addAction(app.data.loadForm, app.data.loadFormButtons, "Select File", app.appActionSelectFiles(app.data.loadForm))
	app.addAction(app.data.loadForm, app.data.loadFormButtons, "Send", app.appActionSendFiles(app.data.loadForm))
	app.addAction(app.data.loadForm, app.data.loadFormButtons, "Cancel", app.actionSwitchToMainWithClear)

	app.data.sendLoginPassForm.
		SetTitle("Send Auth Entry").
		SetBorder(true)

	app.data.sendLoginPassForm.
		AddInputField("Domain", "", 20, nil, nil).
		AddInputField("Login", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil)

	app.addAction(app.data.sendLoginPassForm, app.data.sendLoginPassFormButtons, "Submit", app.appActionSendLoginPass)
	app.addAction(app.data.sendLoginPassForm, app.data.sendLoginPassFormButtons, "Cancel", app.actionSwitchToMain)

	app.data.sendCardForm.
		SetTitle("Send Card Entry").
		SetBorder(true)

	app.data.sendCardForm.
		AddInputField("Title", "", 20, nil, nil).
		AddInputField("Card Number", "", 20, app.checkInputCardField, nil)
	app.addAction(app.data.sendCardForm, app.data.sendCardFormButtons, "Submit", app.appActionSendCard)
	app.addAction(app.data.sendCardForm, app.data.sendCardFormButtons, "Cancel", app.actionSwitchToMain)

}

func (app *App) iniSettings() {
	// Настройки
	app.settingsForm.SetBorder(true).SetTitle("Настройки").SetTitleAlign(tview.AlignLeft)
	app.settingsForm.
		AddInputField("Server Address", app.settings.ServerAddress, 20, nil, func(text string) {
			app.settings.ServerAddress = text
		}).
		// AddTextArea("Bucket", "", 40, 0, 0, nil).
		SetBorder(true).
		SetTitle("Settings").
		SetTitleAlign(tview.AlignLeft)
	app.addAction(app.settingsForm, app.settingsFormButtons, "Cancel", app.actionSwitchToMain)
}

func (app *App) initPages() {

	menu := tview.NewList()

	// Создаем список переходов
	menu.SetBorder(true).SetTitle("Main menu:").SetTitleAlign(tview.AlignLeft)
	menu.
		AddItem("Main", "Go to main page", '1', app.actionSwitchToPerson).
		AddItem("Data list", "View saved data", '2', app.actionSwitchToDataList).
		AddItem("Save file", "Send file", '3', app.actionSwitchToFileForm).
		AddItem("Save auth data", "Send data login and password for domain", '4', app.actionSwitchToLogpassForm).
		AddItem("Save card data", "Send credit card number", '5', app.actionSwitchToCardForm).
		AddItem("Settings", "", 's', app.actionSwitchToSettings).
		AddItem("Quit", "Close application", 'q', app.appActionQuit)

	menu.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		app.logView.Clear()
		app.log.Info("You chose: ", string(shortcut))
	})

	app.pages.AddPage("main", menu, true, true)
	app.pages.AddPage("auth", app.person.authForm, true, false)
	app.pages.AddPage("register", app.person.registerForm, true, false)
	app.pages.AddPage("person", app.person.Form, true, false)
	app.pages.AddPage("datalist", app.data.list, true, false)
	app.pages.AddPage("fileform", app.data.loadForm, true, false)
	app.pages.AddPage("loginpassform", app.data.sendLoginPassForm, true, false)
	app.pages.AddPage("cardform", app.data.sendCardForm, true, false)
	app.pages.AddPage("settings", app.settingsForm, true, false)

	// Check if token is valid
	if isTokenValid() {
		app.actionSwitchToMain()
	} else {
		app.pages.SwitchToPage("auth")
	}

	// Create TextView for logs
	app.logView.
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true)

	// Retranslate logs to TextView
	app.log.SetOutput(app.logView)
	app.log.Info("Welcome to the application!")

	// Flex for logs and pages
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(app.pages, 0, 3, true).
		AddItem(app.logView, 0, 1, false)

	app.tapp.
		SetRoot(flex, true).
		EnableMouse(true)
}

// TODO: check auth prolongation
func isTokenValid() bool {
	return false
}

// Getting data of type data
func (app *App) loadData() error {
	data, err := app.client.GetDataList()
	if err != nil {
		return fmt.Errorf("error client GetDataList: %v", err)
	}

	// Обновление интерфейса на основе полученных данных
	app.updateDatalistPage(data)
	return nil
}

// Render list of type data
func (app *App) updateDatalistPage(data []model.Data) {

	list := tview.NewList()
	// Создание кнопки "Назад"
	list.AddItem("Back", "", 'q', app.actionSwitchToDataListWithClear)

	// Заполнение таблицы данными
	for _, item := range data {
		list.AddItem(item.Title, item.Type, 0, func() {
			app.logView.Clear()
			// Переход к форме с действиями
			app.createDetailForm(item)
		})

	}

	// Добавление Flex как страницы
	app.pages.AddPage("datalistmove", list, true, false)
	app.pages.SwitchToPage("datalistmove")
}

// Getting data with type files
func (app *App) loadFiles() error {
	data, err := app.client.GetFileList()
	if err != nil {
		return fmt.Errorf("error creating stream: %v", err)
	}

	// Обновление интерфейса на основе полученных данных
	app.updateFileDatalistPage(data)
	return nil
}

// List of data with type files
func (app *App) updateFileDatalistPage(data []model.FileItem) {

	list := tview.NewList()
	// Создание кнопки "Назад"
	list.AddItem("Back", "", 'q', app.actionSwitchToDataListWithClear)

	// Заполнение таблицы данными
	for _, item := range data {
		list.AddItem(item.Name, item.Desc, 0, func() {
			app.logView.Clear()
			// Переход к форме с действиями
			app.createMoveForm(item.Hash, item.Name, item.Desc)
		})

	}

	// Добавление Flex как страницы
	app.pages.AddPage("datalistmove", list, true, false)
	app.pages.SwitchToPage("datalistmove")
}

// Detail page of type file with actions
func (app *App) createMoveForm(id string, name, desc string) {
	// Создаем форму с действиями
	actionForm := tview.NewForm()
	actionFormRegister := &FormRegister{}
	actionForm.
		AddTextView("ID", id, 0, 1, false, false).
		AddTextView("Name", name, 0, 1, false, false).
		AddTextView("Description", desc, 0, 1, false, false)

	app.addAction(actionForm, actionFormRegister, "Cancel", app.actionSwitchToDataListWithClear)
	app.addAction(actionForm, actionFormRegister, "Get", app.appActionGetFiles(name, id))
	app.addAction(actionForm, actionFormRegister, "Delete", app.appActionDeleteFiles(name, id))

	// Устанавливаем форму как корневой элемент интерфейса
	app.pages.AddPage("datalistmoveaction", actionForm, true, false)
	app.pages.SwitchToPage("datalistmoveaction")
}

func (app *App) appActionGetFiles(name string, id string) func() {
	return func() {
		app.logView.Clear()
		err := app.client.GetFile(name)
		app.log.Info("Getting started ID: ", id, "\n")
		if err != nil {
			app.log.Info("Error client GetFile: ", err)
			return
		}
		app.log.Info("Got ID: ", id)
	}
}

func (app *App) appActionDeleteFiles(name string, id string) func() {
	return func() {
		app.logView.Clear()
		err := app.client.DeleteFile(name)
		if err != nil {
			app.log.Info("Error client DeleteFile: ", err)
			return
		}
		app.log.Info("Deleted ID: ", id, "\n")
		app.pages.SwitchToPage("datalist")
	}
}

// Detail page of type data with actions
func (app *App) createDetailForm(item model.Data) {
	// Создаем форму с действиями
	actionForm := tview.NewForm()
	actionFormRegister := &FormRegister{}
	actionForm.
		AddTextView("Name", item.Title, 0, 1, false, false).
		AddTextView("ID", strconv.Itoa(int(item.ID)), 0, 1, false, false).
		AddTextView("type", item.Type, 0, 1, false, false).
		AddTextView("Card", item.Card, 0, 1, false, false).
		AddTextView("Login", item.Login, 0, 1, false, false).
		AddTextView("Pass", item.Password, 0, 1, false, false)
	app.addAction(actionForm, actionFormRegister, "Cancel", app.actionSwitchToDataListWithClear)
	app.addAction(actionForm, actionFormRegister, "Delete", app.appActionDeleteData(item.ID))

	// Устанавливаем форму как корневой элемент интерфейса
	app.pages.AddPage("datalistmoveaction", actionForm, true, false)
	app.pages.SwitchToPage("datalistmoveaction")
}

func (app *App) appActionDeleteData(id int64) func() {
	return func() {
		app.logView.Clear()
		app.log.Info("Delete pressed ID: ", id, "\n")
		err := app.client.Delete(id)
		if err != nil {
			app.log.Info("Error delete item: ", err)
			return
		}
		app.pages.SwitchToPage("datalist")
	}
}

func (app *App) addAction(entity *tview.Form, register *FormRegister, title string, action func()) {
	(*register)[title] = title
	entity.AddButton(title, action)
}
