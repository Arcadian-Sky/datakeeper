package settings

import (
	"crypto/rand"
	"encoding/hex"
	"flag"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// - адрес и порт запуска сервиса: переменная окружения ОС `RUN_ADDRESS` или флаг `-a`
// - адрес подключения к базе данных: переменная окружения ОС `PG_DATABASE_URI` или флаг `-d`

type Storage struct {
	Endpoint    string
	AccessKeyID string
	Secret      string
}

type InitedFlags struct {
	Endpoint     string
	DBPGSettings string
	DBMGSettings string
	SecretKey    string
	Storage      Storage
}

func Parse() *InitedFlags {
	end := flag.String("a", ":8080", "endpoint address")
	flagPBDBSettings := flag.String("dp", "", "Адрес подключения к БД")
	flagMBDBSettings := flag.String("dm", "", "Адрес подключения к БД")

	flag.Parse()
	_ = godotenv.Load()

	endpoint := *end
	if envRunAddr := os.Getenv("DATAKEEPER_RUN_ADDRESS"); envRunAddr != "" {
		endpoint = envRunAddr
	}

	dbSettings := *flagPBDBSettings
	if envRunDBSettings := os.Getenv("PG_DATABASE_URI"); envRunDBSettings != "" {
		dbSettings = envRunDBSettings
	}

	dbMdSettings := *flagMBDBSettings
	if envRunDBSettings := os.Getenv("MG_DATABASE_URI"); envRunDBSettings != "" {
		dbMdSettings = envRunDBSettings
	}

	envRunFileStorageURI := os.Getenv("FILE_DATABASE_URI")
	envRunFileStorageAccKeyID := os.Getenv("FILE_DATABASE_ACCESS_KEY")
	envRunFileStorageSecret := os.Getenv("FILE_DATABASE_SECRET")

	// Длина ключа в байтах (например, 32 байта = 256 бит)
	secretKey, err := GenerateSecretKey(32)
	if err != nil {
		fmt.Print("parse err:", err)
	}

	return &InitedFlags{
		Endpoint:     endpoint,
		DBPGSettings: dbSettings,
		DBMGSettings: dbMdSettings,
		SecretKey:    secretKey,
		Storage: Storage{
			Endpoint:    envRunFileStorageURI,
			AccessKeyID: envRunFileStorageAccKeyID,
			Secret:      envRunFileStorageSecret,
		},
	}

}

// GenerateSecretKey генерирует криптографически безопасный случайный ключ.
func GenerateSecretKey(length int) (string, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
