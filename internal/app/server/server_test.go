package server

import (
	"context"
	"database/sql"
	"testing"

	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"github.com/Arcadian-Sky/datakkeeper/mocks"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/test"
	"github.com/stretchr/testify/assert"
)

func TestGetFileRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Создаем мок объекта FileRepository
	mockFileRepo := mocks.NewMockFileRepository(ctrl)

	// Создаем экземпляр App и инициализируем его
	app := &App{
		Logger: logrus.New(),
		Workers: &Workers{
			fileRepo: mockFileRepo,
		},
	}

	// Получаем FileRepository через метод GetFileRepo
	_ = app.GetFileRepo()
}

func TestNewLogger(t *testing.T) {
	logger := NewLogger()

	assert.NotNil(t, logger)

	assert.Equal(t, logrus.TraceLevel, logger.GetLevel())

	textFormatter, ok := logger.Formatter.(*logrus.TextFormatter)
	assert.True(t, ok, "Expected TextFormatter")

	assert.True(t, textFormatter.ForceColors)

	assert.False(t, textFormatter.FullTimestamp)
}

func TestSetDBFileRepo(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Создаем мок объекта FileRepository
	mockFileRepo := mocks.NewMockFileRepository(ctrl)

	app := &App{
		Logger:  logrus.New(),
		Workers: &Workers{},
	}

	// Call SetDBFileRepo with the mock repository
	app.SetDBFileRepo(mockFileRepo)

	assert.NotNil(t, app.Workers.fileRepo)

	// Verify that Workers.fileRepo is set to the mock repository
	assert.Same(t, mockFileRepo, app.Workers.fileRepo, "Expected Workers.fileRepo to be set to mockFileRepo")
}

func TestSetAndGetUserRepo(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockUserRepo := mocks.NewMockUserRepository(ctrl)

	// Create the app instance
	app := &App{
		Workers: &Workers{},
	}

	// Set the UserRepository
	app.SetUserRepo(mockUserRepo)

	// Get the UserRepository
	retrievedRepo := app.GetUserRepo()

	// Assert that the returned repository is the same as the one set
	assert.Equal(t, mockUserRepo, retrievedRepo, "Expected GetUserRepo to return the same repository that was set")
}

func TestSetAndGetDataRepo(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	// Создаем мок объекта FileRepository
	mockDataRepo := mocks.NewMockDataRepository(ctrl)

	// Create the app instance
	app := &App{
		Workers: &Workers{},
	}

	// Set the DataRepository
	app.SetDataRepo(mockDataRepo)

	// Get the DataRepository
	retrievedRepo := app.GetDataRepo()

	// Assert that the returned repository is the same as the one set
	assert.Equal(t, mockDataRepo, retrievedRepo, "Expected GetDataRepo to return the same repository that was set")
}

func TestNewConnectionToPostgresDB(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Mocking the database
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	// Mock ping expectation
	mock.ExpectPing()

	// Call the function
	_, err = NewConnectionToPostgresDB("postgres://mockuser:mockpass@localhost/db", logger)
	assert.Error(t, err)

}

func TestSetLogger(t *testing.T) {
	// Initialize App and Logger
	app := &App{}
	logger := logrus.New()

	// Set Logger
	app.SetLogger(logger)

	// Check if Logger is set correctly
	assert.Equal(t, logger, app.Logger)
}

func TestSetDBPG(t *testing.T) {
	// Initialize App and DB
	app := &App{}
	db, err := sql.Open("pgx", "mock-dsn")
	if err != nil {
		t.Fatalf("failed to open database connection: %v", err)
	}

	// Set DB
	app.SetDBPG(db)

	// Check if DB is set correctly
	assert.Equal(t, db, app.DBPG)
}

func TestApp_SetStorage(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockStorage := mocks.NewMockMinioClient(ctrl)

	ap := &App{}

	ap.SetStorage(mockStorage)

	assert.Equal(t, mockStorage, ap.Storage, "Storage was not set correctly")
}

func TestApp_SetFlags(t *testing.T) {
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	flags := &settings.InitedFlags{
		Endpoint: "test-flag",
	}

	ap := &App{
		Logger: logger,
	}

	hook := test.NewLocal(logger)

	ap.SetFlags(flags)

	assert.Equal(t, flags, ap.Flags, "Flags were not set correctly")

	entries := hook.AllEntries()
	assert.Len(t, entries, 1, "Expected one log entry")
	assert.Contains(t, entries[0].Message, "parsed", "Log message does not contain 'parsed'")
	assert.Contains(t, entries[0].Message, "test-flag", "Log message does not contain 'test-flag'")
}

func TestNewConnectToMinIO_Error(t *testing.T) {
	logg := logrus.New()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// mockMinioClient := mocks.NewMockMinioClient(ctrl)
	// mockMinioClient.EXPECT().ListBuckets(gomock.Any()).Return(nil, nil).Times(1)

	// Settings for MinIO
	storageSettings := settings.Storage{
		Endpoint:    "localhost:9000000",
		AccessKeyID: "testAccessKey",
		Secret:      "testSecret",
	}

	ctx := context.Background()

	got, err := NewСonnectToMinIO(ctx, storageSettings, logg)

	assert.Error(t, err)
	assert.Nil(t, got)
}
