package server

import (
	"context"
	"database/sql"
	"time"

	"github.com/Arcadian-Sky/datakkeeper/internal/server/repository"
	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"github.com/Arcadian-Sky/datakkeeper/migrations"
	minioclient "github.com/Arcadian-Sky/datakkeeper/tools/client"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/pressly/goose/v3"
	"github.com/sirupsen/logrus"
)

type Workers struct {
	userRepo repository.UserRepository
	dataRepo repository.DataRepository
	fileRepo repository.FileRepository
}
type App struct {
	Logger *logrus.Logger
	DBPG   *sql.DB
	// DBMG   *mongo.Client
	// Storage *minio.Client
	Storage minioclient.MinioClient
	Flags   *settings.InitedFlags
	Ctx     context.Context
	CncF    context.CancelFunc
	Workers *Workers
}

func NewApp(ctx *context.Context, ctcf *context.CancelFunc) (*App, error) {
	app := App{
		Workers: &Workers{},
		Ctx:     *ctx,
		CncF:    *ctcf,
	}
	return &app, nil
}

func NewLogger() *logrus.Logger {
	logg := logrus.New()
	logg.SetLevel(logrus.TraceLevel)
	logg.SetFormatter(&logrus.TextFormatter{
		ForceColors:   true,  // Включить цвета в выводе
		FullTimestamp: false, // Включить полный временной штамп
	})
	return logg
}

func (ap *App) SetLogger(logg *logrus.Logger) {
	ap.Logger = logg
}

func (ap *App) SetDBPG(db *sql.DB) {
	ap.DBPG = db
}

func (ap *App) SetStorage(db minioclient.MinioClient) {
	ap.Storage = db
}

func (ap *App) SetFlags(st *settings.InitedFlags) {
	ap.Logger.Debug("parsed: ", st, "\n")
	ap.Flags = st
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

func (app *App) GetDataRepo() repository.DataRepository {
	return app.Workers.dataRepo
}

// Репозиторий по работе с пользователем
func (app *App) SetUserRepo(uR repository.UserRepository) {
	app.Workers.userRepo = uR
}
func (app *App) GetUserRepo() repository.UserRepository {
	return app.Workers.userRepo
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
func (app *App) GetFileRepo() repository.FileRepository {
	return app.Workers.fileRepo
}
