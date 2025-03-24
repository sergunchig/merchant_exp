package logger

import (
	"log"
	"os"
)

type AppLogger struct {
	log *log.Logger
}

func NewLogger(logfile string) (*AppLogger, error) {
	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	lg := log.New(file, "log", log.Ldate|log.Ltime|log.Lshortfile)
	return &AppLogger{
		log: lg,
	}, nil
}
