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
		FILE string
	}
	DB struct {
		DBCONNECTION string
	}
	FileStorage struct {
		PATH string
	}
)

func NewConfig() (*Config, error) {

	fmt.Println("init  configuration")
	err := godotenv.Load()
	if err != nil {
		fmt.Println("can't load environment")
		panic("can't load environment")
	}

	host := os.Getenv("HOST")
	db := os.Getenv("DBCONNECTION")
	logfile := os.Getenv("LOGFILE")
	fileStorage := os.Getenv("FILESTORAGE")

	cfg := &Config{
		HTTP:    HTTP{HOST: host},
		DB:      DB{DBCONNECTION: db},
		Log:     Log{FILE: logfile},
		Storage: FileStorage{PATH: fileStorage},
	}
	return cfg, nil
}
