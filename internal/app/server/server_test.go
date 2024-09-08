package server

import (
	"testing"

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
	// Создаем мок объекта FileRepository
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
	assert.Equal(t, mockUserRepo, *retrievedRepo, "Expected GetUserRepo to return the same repository that was set")
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
	assert.Equal(t, mockDataRepo, *retrievedRepo, "Expected GetDataRepo to return the same repository that was set")
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
