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
	settings     *settings.ClientConfig
	tapp         *tview.Application
	person       Person
	data         Data
	pages        *tview.Pages
	settingsForm *tview.Form
	logView      *tview.TextView
	client       client.GRPCClientInterface
	Conn         *grpc.ClientConn
	storage      *client.MemStorage
	log          *logrus.Logger
}

func NewEmptyApp() *App {
	return &App{
		tapp:         tview.NewApplication(),
		pages:        tview.NewPages(),
		settingsForm: tview.NewForm(),
		logView:      tview.NewTextView(),
		person: Person{
			authForm:            tview.NewForm(),
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
	// Создаем формы для авторизации
	app.person.authForm.SetBorder(true).SetTitle("Authorize").SetTitleAlign(tview.AlignLeft)
	app.person.authForm.
		AddInputField("Login", "", 20, nil, nil).
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Save", func() {
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
			app.pages.SwitchToPage("main")
		}).
		AddButton("Switch to Register", func() {
			app.logView.Clear()
			app.pages.SwitchToPage("register")
		}).
		AddButton("Quit", func() {
			app.tapp.Stop()
		})

	// Создаем формы для регистрации
	app.person.registerForm.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
	app.person.registerForm.
		AddInputField("Login", "", 20, nil, nil).
		AddPasswordField("Password", "", 10, '*', nil)

	app.addAction(app.person.registerForm, app.person.registerFormButtons, "Save", func() {
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
		app.pages.SwitchToPage("main")
	})
	app.addAction(app.person.registerForm, app.person.registerFormButtons, "Switch to Authorize", func() {
		app.logView.Clear()
		app.pages.SwitchToPage("auth")
	})
	app.addAction(app.person.registerForm, app.person.registerFormButtons, "Quit", func() {
		app.tapp.Stop()
	})

	// Создаем формы для авторизации и регистрации
	app.person.Form.SetBorder(true).SetTitle("Enter some data").SetTitleAlign(tview.AlignLeft)
}

// Инициализация интерфейсов работы с данными
func (app *App) initDataInterfaces() {
	// Создаем интерфейс для отображения данных
	app.data.list.SetTitle("Data List")
	app.data.list.AddItem("Back", "Make click to run action", 'q', func() {
		app.pages.SwitchToPage("main")
	})
	app.data.list.AddItem("Load Files", "", 'f', func() {
		app.logView.Clear()
		err := app.loadFiles()
		if err != nil {
			app.log.Info("Error loading data: ", err)
		}
	})
	app.data.list.AddItem("Load Data", "", 'f', func() {
		app.logView.Clear()
		err := app.loadData()
		if err != nil {
			app.log.Info("Error loading data: ", err)
		}
	})

	// Создаем форму для отправки данных
	app.data.loadForm.
		SetTitle("Send File Entry").
		SetBorder(true)

	app.data.loadForm.
		AddInputField("File Path", "", 40, nil, nil)

	app.addAction(app.data.loadForm, app.data.loadFormButtons, "Select File", func() {
		filename, err := dialog.File().Load()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error selecting file: %v\n", err)
			return
		}

		// Устанавливаем выбранный файл в InputField
		app.data.loadForm.GetFormItem(0).(*tview.InputField).SetText(filename)
	})

	app.addAction(app.data.loadForm, app.data.loadFormButtons, "Send", func() {
		app.logView.Clear()
		app.log.Info("Sending data...")
		filePath := app.data.loadForm.GetFormItem(0).(*tview.InputField).GetText()
		app.log.Info("File path: ", filePath)
		if err := app.client.UploadFile(filePath); err != nil {
			app.log.Info(fmt.Println("Error uploading file:", err))
		}

	})

	app.addAction(app.data.loadForm, app.data.loadFormButtons, "Quit", func() {
		app.logView.Clear()
		app.pages.SwitchToPage("main")
	})

	app.data.sendLoginPassForm.
		SetTitle("Send Auth Entry").
		SetBorder(true)

	app.data.sendLoginPassForm.
		AddInputField("Domain", "", 20, nil, nil).
		AddInputField("Login", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil)
	app.addAction(app.data.sendLoginPassForm, app.data.sendLoginPassFormButtons, "Submit", func() {
		domain := app.data.sendLoginPassForm.GetFormItem(0).(*tview.InputField).GetText()
		login := app.data.sendLoginPassForm.GetFormItem(1).(*tview.InputField).GetText()
		password := app.data.sendLoginPassForm.GetFormItem(2).(*tview.InputField).GetText()

		app.log.Info(fmt.Printf("Login: %s, Password: %s\n, Domain: %s\n", login, "*****", domain))

		if err := app.client.SaveLoginPass(domain, login, password); err != nil {
			app.log.Info(fmt.Println("Error saveing login pass:", err))
		}
	})

	app.addAction(app.data.sendLoginPassForm, app.data.sendLoginPassFormButtons, "Quit", func() {
		app.pages.SwitchToPage("main")
	})

	app.data.sendCardForm.
		SetTitle("Send Card Entry").
		SetBorder(true)

	app.data.sendCardForm.
		AddInputField("Title", "", 20, nil, nil).
		AddInputField("Card Number", "", 20, func(textToCheck string, lastChar rune) bool {
			// Allow only numbers and spaces in the card number field
			return (lastChar >= '0' && lastChar <= '9') || lastChar == ' '
		}, nil)
	app.addAction(app.data.sendCardForm, app.data.sendCardFormButtons, "Submit", func() {
		domain := app.data.sendCardForm.GetFormItem(0).(*tview.InputField).GetText()
		card := app.data.sendCardForm.GetFormItem(1).(*tview.InputField).GetText()

		app.log.Info(fmt.Printf("Title: %s, Card: %s\n", domain, card))
		if err := app.client.SaveCard(domain, card); err != nil {
			app.log.Info(fmt.Println("Error saveing login pass:", err))
		}
	})

	app.addAction(app.data.sendCardForm, app.data.sendCardFormButtons, "Cancel", func() {
		app.pages.SwitchToPage("main")
	})

}

func (app *App) iniSettings() {
	// Настройки
	app.settingsForm.SetBorder(true).SetTitle("Настройки").SetTitleAlign(tview.AlignLeft)
	app.settingsForm.
		AddInputField("Server Address", app.settings.ServerAddress, 20, nil, func(text string) {
			app.settings.ServerAddress = text
		}).
		// AddTextArea("Bucket", "", 40, 0, 0, nil).
		AddButton("Back", func() {
			app.pages.SwitchToPage("main")
		}).
		SetBorder(true).
		SetTitle("Settings").
		SetTitleAlign(tview.AlignLeft)
}

func (app *App) initPages() {

	menu := tview.NewList()

	// Создаем список переходов
	menu.SetBorder(true).SetTitle("Main menu:").SetTitleAlign(tview.AlignLeft)
	menu.
		AddItem("Main", "Go to main page", '1', func() {
			app.pages.SwitchToPage("person")
		}).
		AddItem("Data list", "View saved data", '2', func() {
			app.pages.SwitchToPage("datalist")
		}).
		AddItem("Save file", "Send file", '3', func() {
			app.pages.SwitchToPage("fileform")
		}).
		AddItem("Save auth data", "Send data login and password for domain", '4', func() {
			app.pages.SwitchToPage("loginpassform")
		}).
		AddItem("Save card data", "Send credit card number", '5', func() {
			app.pages.SwitchToPage("cardform")
		}).
		AddItem("Settings", "", 's', func() {
			app.pages.SwitchToPage("settings")
		}).
		AddItem("Quit", "Close application", 'q', func() {
			app.tapp.Stop()
		})

	menu.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		app.logView.Clear()
		app.log.Info(" Вы выбрали: ", string(shortcut))
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
		app.pages.SwitchToPage("main")
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
	list.AddItem("Назад", "", 'q', func() {
		app.logView.Clear()
		app.pages.SwitchToPage("datalist")
	})

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
	list.AddItem("Назад", "", 'q', func() {
		app.logView.Clear()
		app.pages.SwitchToPage("datalist")
	})

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
	actionForm.
		AddTextView("ID", id, 0, 1, false, false).            // Показываем ID
		AddTextView("Name", name, 0, 1, false, false).        // Показываем имя
		AddTextView("Description", desc, 0, 1, false, false). // Показываем описание
		AddButton("Назад", func() {
			app.logView.Clear()
			app.pages.SwitchToPage("datalist") // Возвращаемся к списку
		}).
		AddButton("Получить", func() {
			app.logView.Clear()
			err := app.client.GetFile(name)
			if err != nil {
				app.log.Info("Error client GetFile: ", err)
				return
			}
			fmt.Printf("Получить нажато для ID: %v\n", id)
		}).
		AddButton("Удалить", func() {
			app.logView.Clear()
			err := app.client.DeleteFile(name)
			if err != nil {
				app.log.Info("Error client DeleteFile: ", err)
				return
			}
			fmt.Printf("Удалить нажато для ID: %v\n", id)
			app.pages.SwitchToPage("datalist") // Возвращаемся к списку
		})

	// Устанавливаем форму как корневой элемент интерфейса
	app.pages.AddPage("datalistmoveaction", actionForm, true, false)
	app.pages.SwitchToPage("datalistmoveaction")
}

// Detail page of type data with actions
func (app *App) createDetailForm(item model.Data) {
	// Создаем форму с действиями
	actionForm := tview.NewForm()
	actionForm.
		AddTextView("Name", item.Title, 0, 1, false, false).               // Показываем имя
		AddTextView("ID", strconv.Itoa(int(item.ID)), 0, 1, false, false). // Показываем ID
		AddTextView("type", item.Type, 0, 1, false, false).                // Показываем ID
		AddTextView("Card", item.Card, 0, 1, false, false).                // Показываем ID
		AddTextView("Login", item.Login, 0, 1, false, false).              // Показываем ID
		AddTextView("Pass", item.Password, 0, 1, false, false).            // Показываем ID
		AddButton("Назад", func() {
			app.logView.Clear()
			app.pages.SwitchToPage("datalist") // Возвращаемся к списку
		}).
		AddButton("Удалить", func() {
			app.logView.Clear()
			fmt.Printf("Удалить нажато для ID: %v\n", item.ID)
			err := app.client.Delete(item.ID)
			if err != nil {
				app.log.Info("Error delete item: ", err)
				return
			}
			app.pages.SwitchToPage("datalist") // Возвращаемся к списку
		})

	// Устанавливаем форму как корневой элемент интерфейса
	app.pages.AddPage("datalistmoveaction", actionForm, true, false)
	app.pages.SwitchToPage("datalistmoveaction")
}

func (app *App) addAction(entity *tview.Form, register *FormRegister, title string, action func()) {
	(*register)[title] = title
	entity.AddButton(title, action)
}
