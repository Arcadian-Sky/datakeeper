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

type Person struct {
	authForm     *tview.Form
	registerForm *tview.Form
	Form         *tview.Form
}

type Data struct {
	loadForm          *tview.Form
	sendLoginPassForm *tview.Form
	sendCardForm      *tview.Form
	list              *tview.List
	move              *tview.List
}
type App struct {
	settings     *settings.ClientConfig
	tapp         *tview.Application
	person       Person
	data         Data
	pages        *tview.Pages
	settingsForm *tview.Form
	logView      *tview.TextView
	client       client.GRPCClient
	Conn         *grpc.ClientConn
	storage      *client.MemStorage
	log          *logrus.Logger
}

func NewClientApp(st *settings.ClientConfig) *App {
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
		settings: st,
	}
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
		AddPasswordField("Password", "", 10, '*', nil).
		AddButton("Save", func() {
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
		}).
		AddButton("Switch to Authorize", func() {
			app.logView.Clear()
			app.pages.SwitchToPage("auth")
		}).
		AddButton("Quit", func() {
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
		AddInputField("File Path", "", 40, nil, nil).
		AddButton("Select File", func() {
			filename, err := dialog.File().Load()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Error selecting file: %v\n", err)
				return
			}

			// Устанавливаем выбранный файл в InputField
			app.data.loadForm.GetFormItem(0).(*tview.InputField).SetText(filename)
		}).
		AddButton("Send", func() {
			app.logView.Clear()
			app.log.Info("Sending data...")
			filePath := app.data.loadForm.GetFormItem(0).(*tview.InputField).GetText()
			app.log.Info("File path: ", filePath)
			if err := app.client.UploadFile(filePath); err != nil {
				app.log.Info(fmt.Println("Error uploading file:", err))
			}

		}).
		AddButton("Quit", func() {
			app.logView.Clear()
			app.pages.SwitchToPage("main")
		})

	app.data.sendLoginPassForm.
		SetTitle("Send Auth Entry").
		SetBorder(true)

	app.data.sendLoginPassForm.
		AddInputField("Domain", "", 20, nil, nil).
		AddInputField("Login", "", 20, nil, nil).
		AddPasswordField("Password", "", 20, '*', nil).
		AddButton("Submit", func() { // Кнопка "Submit"
			domain := app.data.sendLoginPassForm.GetFormItem(0).(*tview.InputField).GetText()
			login := app.data.sendLoginPassForm.GetFormItem(1).(*tview.InputField).GetText()
			password := app.data.sendLoginPassForm.GetFormItem(2).(*tview.InputField).GetText()

			app.log.Info(fmt.Printf("Login: %s, Password: %s\n, Domain: %s\n", login, "*****", domain))

			if err := app.client.SaveLoginPass(domain, login, password); err != nil {
				app.log.Info(fmt.Println("Error saveing login pass:", err))
			}
		}).
		AddButton("Quit", func() {
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
		}, nil).
		AddButton("Submit", func() {
			domain := app.data.sendCardForm.GetFormItem(0).(*tview.InputField).GetText()
			card := app.data.sendCardForm.GetFormItem(1).(*tview.InputField).GetText()

			app.log.Info(fmt.Printf("Title: %s, Card: %s\n", domain, card))
			if err := app.client.SaveCard(domain, card); err != nil {
				app.log.Info(fmt.Println("Error saveing login pass:", err))
			}
		}).
		AddButton("Cancel", func() {
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
		AddItem("Главная", "Перейти на главную страницу", '1', func() {
			app.pages.SwitchToPage("person")
		}).
		AddItem("Список данных", "Посмотреть сохраненные данные", '2', func() {
			app.pages.SwitchToPage("datalist")
		}).
		AddItem("Загрузка файла", "Отправить файл", '3', func() {
			app.pages.SwitchToPage("fileform")
		}).
		AddItem("Загрузка авторизации", "Отправить данные пары логин пароль к домены", '4', func() {
			app.pages.SwitchToPage("loginpassform")
		}).
		AddItem("Загрузка карты", "Отправить номер кредитки", '5', func() {
			app.pages.SwitchToPage("cardform")
		}).
		AddItem("Настройки", "Настройки", 's', func() {
			app.pages.SwitchToPage("settings")
		}).
		AddItem("Выход", "Закрыть приложение", 'q', func() {
			app.tapp.Stop()
		})

	// Устанавливаем обработчик выбора элемента
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

	// Проверяем валидность токена
	if isTokenValid() {
		// Если токен действителен, показываем главное меню
		app.pages.SwitchToPage("main")
	} else {
		// Если токен недействителен, показываем страницу авторизации
		app.pages.SwitchToPage("auth")
	}

	// Создаём TextView для логов
	app.logView.
		SetDynamicColors(true).
		SetRegions(true).
		SetScrollable(true)

	// Перенаправляем вывод логов в TextView
	app.log.SetOutput(app.logView)
	app.log.Info("Welcome to the application!")

	// Flex для размещения логов и страниц
	flex := tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(app.pages, 0, 3, true).   // Страницы приложения будут ниже
		AddItem(app.logView, 0, 1, false) // Логи будут в нижней части

		// if err := app.tapp.SetRoot(app.pages, true).EnableMouse(true).Run(); err != nil {
		// 	panic(err)
		// }
	// Устанавливаем Flex как корневой элемент
	if err := app.tapp.SetRoot(flex, true).EnableMouse(true).Run(); err != nil {
		panic(err)
	}
}

// Симуляция проверки токена авторизации
func isTokenValid() bool {
	return false
}

// Вызов клиента для получения данных из бд
func (app *App) loadData() error {
	data, err := app.client.GetDataList()
	if err != nil {
		return fmt.Errorf("error client GetDataList: %v", err)
	}

	// Обновление интерфейса на основе полученных данных
	app.updateDatalistPage(data)
	return nil
}
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

// Вызов клиента для получения данных файлов
func (app *App) loadFiles() error {
	data, err := app.client.GetFileList()
	if err != nil {
		return fmt.Errorf("error creating stream: %v", err)
	}

	// Обновление интерфейса на основе полученных данных
	app.updateFileDatalistPage(data)
	return nil
}

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

// createActionForm создает форму с действиями, используя переданный ID
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
