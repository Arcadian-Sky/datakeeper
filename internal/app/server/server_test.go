package server

import (
	"bytes"
	"database/sql"
	"testing"

	"github.com/Arcadian-Sky/datakkeeper/internal/settings"
	"github.com/Arcadian-Sky/datakkeeper/mocks"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/golang/mock/gomock"
	"github.com/sirupsen/logrus"
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

// func TestSetStorage(t *testing.T) {
// 	ctrl := gomock.NewController(t)
// 	defer ctrl.Finish()
// 	// Создаем мок объекта FileRepository
// 	mockStorage := mocks.NewMockFileRepository(ctrl)

// 	app := &App{}

// 	// Set Storage
// 	app.SetStorage(mockStorage)

// 	// Check if Storage is set correctly
// 	assert.Equal(t, mockStorage, app.Storage)
// }

func TestSetFlags(t *testing.T) {
	// Initialize App and Settings
	app := &App{}
	flags := &settings.InitedFlags{
		SecretKey: "test-secret",
	}

	// Create a custom logger to capture debug logs
	logger := logrus.New()
	buf := new(bytes.Buffer)
	logger.SetOutput(buf)

	// Set Logger and Flags
	app.SetLogger(logger)
	app.SetFlags(flags)

	// Check if Flags are set correctly
	assert.Equal(t, flags, app.Flags)

	// Check if debug log contains expected output
	// expectedLog := "parsed: &{SecretKey:test-secret} \n"
	// fmt.Printf("buf.String(): %v\n", buf.String())
	// assert.Contains(t, buf.String(), expectedLog)
}
