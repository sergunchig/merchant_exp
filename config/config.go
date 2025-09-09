package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type (
	Config struct {
		App     App
		HTTP    HTTP
		Log     Log
		DB      DB
		Storage FileStorage
	}
	App struct {
		Version string
	}
	HTTP struct {
		HOST string
	}
	Log struct {
		PATH string
	}
	DB struct {
		DBCONNECTION string
	}
	FileStorage struct {
		PATH string
	}
)

func NewConfig() (*Config, error) {

	// todo давай использовать логи?
	//logger еще не инициализирован
	err := godotenv.Load()
	if err != nil {
		return nil, fmt.Errorf("can't load environment: %w", err)
	}

	host := os.Getenv("HOST")
	db := os.Getenv("DBCONNECTION")
	logPath := os.Getenv("LOGPATH")
	fileStorage := os.Getenv("FILESTORAGE")

	cfg := &Config{
		HTTP:    HTTP{HOST: host},
		DB:      DB{DBCONNECTION: db},
		Log:     Log{PATH: logPath},
		Storage: FileStorage{PATH: fileStorage},
	}
	return cfg, nil
}
