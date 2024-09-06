package server

import (
	"context"
	"database/sql"
	"time"

	"github.com/Arcadian-Sky/datakkeeper/internal/server/repository"
	minioclient "github.com/Arcadian-Sky/datakkeeper/internal/server/repository/client"
	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"github.com/Arcadian-Sky/datakkeeper/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Workers struct {
	userRepo repository.UserRepository
	dataRepo repository.DataRepository
	fileRepo repository.FileRepository
}
type App struct {
	Logger *logrus.Logger
	DBPG   *sql.DB
	DBMG   *mongo.Client
	// Storage *minio.Client
	Storage minioclient.MinioClient
	Flags   *settings.InitedFlags
	Ctx     context.Context
	CncF    context.CancelFunc
	Workers *Workers
}

func NewApp(ctx *context.Context, ctcf *context.CancelFunc) (*App, error) {
	logg := logrus.New()
	logg.SetLevel(logrus.TraceLevel)
	logg.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,  // Включить цвета в выводе
		FullTimestamp: false, // Включить полный временной штамп
	})
	app := App{
		Logger:  logg,
		Workers: &Workers{},
		Ctx:     *ctx,
		CncF:    *ctcf,
	}
	parsed := settings.Parse()
	logg.Debug("parsed: ", parsed, "\n")
	//set logger
	// logg.SetFormatter(&logrus.JSONFormatter{})

	//set db pg connect
	logg.Debug("parsed.PGDBSettings: ", parsed.DBPGSettings, "\n")
	dbP, err := NewConnectionToPostgresDB(parsed.DBPGSettings, logg)
	if err != nil {
		return &app, err
	}

	// TODO
	// logg.Log(logrus.DebugLevel, "parsed.DBMGSettings: ", parsed.DBMGSettings, "\n")
	// dbM, err := NewСonnectToMongoDB(parsed.DBMGSettings, logg)
	// if err != nil {
	// 	return &app, err
	// }
	// defer func() {
	// 	if err := dbM.Disconnect(context.TODO()); err != nil {
	// 		logg.Fatalf("Failed to disconnect from MongoDB: %v", err)
	// 	}
	// }()

	logg.Debug("parsed.MinIOSettings: ", parsed.Storage, "\n")
	client, err := NewСonnectToMinIO(app.Ctx, parsed.Storage, logg)
	if err != nil {
		logg.Fatalf("Failed to connect to MinIO: %v", err)
	}

	app.DBPG = dbP
	// app.DBMG = dbM
	app.Storage = client
	app.Flags = parsed

	return &app, nil
}

// Подключение к постгрес
func NewConnectionToPostgresDB(dsn string, logg *logrus.Logger) (*sql.DB, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		logrus.Error("failed to create a database connection", err)
		return nil, err
	}

	if err = db.Ping(); err != nil {
		logrus.Error("failed to ping the database", err)
		return nil, err
	}

	logg.Log(logrus.InfoLevel, "Successfully connected to PostgresDB")
	return db, err
}

// Подключение к mongo
func NewСonnectToMongoDB(uri string, logg *logrus.Logger) (*mongo.Client, error) {
	// Создаем контекст с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Настройка параметров подключения
	clientOptions := options.Client().ApplyURI(uri)

	// Подключение к MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Проверка соединения
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	logg.Log(logrus.InfoLevel, "Successfully connected to MongoDB")
	return client, nil
}

// Подключение к minio
func NewСonnectToMinIO(ctx context.Context, settings settings.Storage, logg *logrus.Logger) (minioclient.MinioClient, error) {

	endpoint := settings.Endpoint
	accessKeyID := settings.AccessKeyID
	secretAccessKey := settings.Secret
	creds := credentials.NewStaticV4(accessKeyID, secretAccessKey, "")
	useSSL := false

	// Создание нового клиента MinIO
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  creds,
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// Проверка подключения к MinIO
	_, err = client.ListBuckets(ctx)
	if err != nil {
		return nil, err
	}

	logg.Log(logrus.InfoLevel, "Successfully connected to MinIO")
	return minioclient.NewMinioClient(client), nil
}

func (app *App) SetDataRepo(dR repository.DataRepository) {
	app.Workers.dataRepo = dR
}

func (app *App) GetDataRepo() *repository.DataRepository {
	return &app.Workers.dataRepo
}

// Репозиторий по работе с пользователем
func (app *App) SetUserRepo(uR repository.UserRepository) {
	app.Workers.userRepo = uR
}
func (app *App) GetUserRepo() *repository.UserRepository {
	return &app.Workers.userRepo
}

func (app *App) MigrateDBPG() error {
	goose.SetBaseFS(migrations.Migrations)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := goose.RunContext(ctx, "up", app.DBPG, ".")
	if err != nil {
		return err
	}

	return nil
}

// Репозиторий по работе с документами
func (app *App) SetDBFileRepo(fR repository.FileRepository) {
	app.Workers.fileRepo = fR
}
func (app *App) GetFileRepo() *repository.FileRepository {
	return &app.Workers.fileRepo
}
