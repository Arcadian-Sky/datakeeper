package client

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Arcadian-Sky/datakkeeper/internal/client"
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
	token        string
	bucket       string
}

type Data struct {
	loadForm *tview.Form
	list     *tview.List
	move     *tview.List
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
			list:     tview.NewList(),
			move:     tview.NewList(),
			loadForm: tview.NewForm(),
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
	app.initMenu()

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

		// filePath := app.data.loadForm.GetFormItem(0).(*tview.InputField).GetText()
		// if err := app.uploadFile(filePath); err != nil {
		// 	fmt.Println("Error uploading file:", err)
		// }
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
	app.data.list.AddItem("Back", "", 'q', func() {
		app.pages.SwitchToPage("main")
	})
	app.data.list.AddItem("Load Files", "", 'f', func() {
		app.logView.Clear()
		err := app.loadData()
		if err != nil {
			app.log.Info("Error loading data: ", err)
		}
	})
	app.data.list.AddItem("Load Data", "", 'f', func() {
		app.pages.SwitchToPage("datalist")
		// err := loadData()
		// if err != nil {
		// 	// Обработка ошибки
		// 	fmt.Fprintf(app.logView, "Error loading data: %v\n", err)
		// } else {
		// 	// Данные успешно загружены, переключаемся на страницу
		// 	app.pages.SwitchToPage("datalist")
		// }
	})
	// app.data.list.AddItem("List item 1", "Some explanatory text", 'a', nil)
	// for _, item := range data {
	// 	app.data.list.AddItem(item, "", 0, nil)
	// }

	// Создаем форму для отправки данных
	app.data.loadForm.
		SetTitle("Send Data").
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
			// Добавить логику отправки данных здесь
			fmt.Println("Sending data...")

			filePath := app.data.loadForm.GetFormItem(0).(*tview.InputField).GetText()
			if err := app.uploadFile(filePath); err != nil {
				fmt.Println("Error uploading file:", err)
			}

		}).
		AddButton("Quit", func() {
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

func (app *App) initMenu() {

	menu := tview.NewList()

	// Создаем список переходов
	menu.SetBorder(true).SetTitle("Main menu:").SetTitleAlign(tview.AlignLeft)
	menu.
		AddItem("Главная", "Перейти на главную страницу", '1', func() {
			app.pages.SwitchToPage("person")
		}).
		// AddItem("Авторизация", "Авторизация", '2', func() {
		// 	app.pages.SwitchToPage("auth")
		// }).
		// AddItem("Авторизация", "Регистрация", '2', func() {
		// 	app.pages.SwitchToPage("register")
		// }).
		AddItem("Список данных", "Посмотреть сохраненные данные", '3', func() {
			app.pages.SwitchToPage("datalist")
		}).
		AddItem("Загрузка данных", "Отправить данные", '4', func() {
			app.pages.SwitchToPage("dataform")
		}).
		AddItem("Настройки", "Настройки", 's', func() {
			app.pages.SwitchToPage("settings")
		}).
		AddItem("Выход", "Закрыть приложение", 'q', func() {
			app.tapp.Stop()
		})

	// Устанавливаем обработчик выбора элемента
	menu.SetSelectedFunc(func(index int, mainText string, secondaryText string, shortcut rune) {
		fmt.Print(" Вы выбрали: ", string(shortcut))
	})

	app.pages.AddPage("main", menu, true, true)
	app.pages.AddPage("auth", app.person.authForm, true, false)
	app.pages.AddPage("register", app.person.registerForm, true, false)
	app.pages.AddPage("person", app.person.Form, true, false)
	app.pages.AddPage("datalist", app.data.list, true, false)
	app.pages.AddPage("dataform", app.data.loadForm, true, false)
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
	// Пример вывода в лог
	// log.Println("Application started")
	// fmt.Fprintf(app.logView, "Welcome to the application!\n")

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

// Симуляция сохранения нового токена после успешной авторизации
func saveNewToken() {
	// Логика сохранения токена (например, запись в файл/БД)
	// Здесь просто пример с сохранением времени
	fmt.Println("Новый токен сохранен:", time.Now())
}

func (app *App) uploadFile(filePath string) error {

	// Открываем файл для чтения
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	stream, err := app.client.UploadFile(context.Background())
	if err != nil {
		return fmt.Errorf("error creating stream: %v", err)
	}

	// Отправляем части файла
	// buffer := make([]byte, 1024)
	// for {
	// 	n, err := file.Read(buffer)
	// 	if err == io.EOF {
	// 		break
	// 	}
	// 	if err != nil {
	// 		return fmt.Errorf("error reading file: %v", err)
	// 	}

	// 	err = stream.Send(&pb.FileChunk{
	// 		Content:  buffer[:n],
	// 		Filename: filePath, // Передаем имя файла
	// 	})
	// 	if err != nil {
	// 		return fmt.Errorf("error sending chunk: %v", err)
	// 	}
	// }

	// Завершаем отправку и получаем ответ
	status, err := stream.CloseAndRecv()
	if err != nil {
		return fmt.Errorf("error receiving response: %v", err)
	}

	fmt.Printf("Upload finished with status: %v\n", status.Message)
	return nil
}

type Item struct {
	ID   int64
	Name string
	Desc string
}

func (app *App) loadData() error {
	// Вызов клиента для получения данных

	err := app.client.GetFileList()
	if err != nil {
		return fmt.Errorf("error creating stream: %v", err)
	}

	data := []Item{
		{ID: 1, Name: "Item 1"},
		{ID: 2, Name: "Item 2"},
		{ID: 3, Name: "Item 3"},
		{ID: 4, Name: "Item 4"},
		{ID: 54, Name: "Item 54"},
		{ID: 64, Name: "Item 64"},
		{ID: 674, Name: "Item 674"},
		{ID: 34, Name: "Item 34"},
		{ID: 84, Name: "Item 84"},
		{ID: 49, Name: "Item 84"},
	}

	// Обновление интерфейса на основе полученных данных
	app.updateDatalistPage(data)
	return nil
}

func (app *App) updateDatalistPage(data []Item) {

	list := tview.NewList()
	// Создание кнопки "Назад"
	list.AddItem("Назад", "", 'q', func() {
		app.pages.SwitchToPage("datalist")
	})

	// Заполнение таблицы данными
	for _, item := range data {
		list.AddItem(item.Name, item.Desc, 0, func() {
			// Переход к форме с действиями
			app.createMoveForm(int(item.ID), item.Name, item.Desc)
		})

	}

	// Добавление Flex как страницы
	app.pages.AddPage("datalistmove", list, true, false)
	app.pages.SwitchToPage("datalistmove")
}

// createActionForm создает форму с действиями, используя переданный ID

func (app *App) createMoveForm(id int, name, desc string) {
	// Создаем форму с действиями
	actionForm := tview.NewForm().
		AddTextView("ID", strconv.Itoa(id), 0, 1, false, false). // Показываем ID
		AddTextView("Name", name, 0, 1, false, false).           // Показываем имя
		AddTextView("Description", desc, 0, 1, false, false).    // Показываем описание
		AddButton("Назад", func() {
			app.pages.SwitchToPage("datalist") // Возвращаемся к списку
		}).
		AddButton("Получить", func() {
			fmt.Printf("Получить нажато для ID: %d\n", id)
		}).
		AddButton("Удалить", func() {
			fmt.Printf("Удалить нажато для ID: %d\n", id)
		})

	// Устанавливаем форму как корневой элемент интерфейса
	app.pages.AddPage("datalistmoveaction", actionForm, true, false)
	app.pages.SwitchToPage("datalistmoveaction")
}
