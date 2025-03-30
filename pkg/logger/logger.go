package logger

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

type AppLogger struct {
	log *log.Logger
}

func NewLogger(path string) (*AppLogger, error) {

	//date := time.Now().Format(time.DateOnly)
	time := time.Now().Format(time.Stamp)
	time = strings.ReplaceAll(time, ":", "_")

	logfile := fmt.Sprintf("%slog_%s.log", path, time)

	file, err := os.OpenFile(logfile, os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		return nil, err
	}
	lg := log.New(file, "log", log.Ldate|log.Ltime|log.Lshortfile)
	return &AppLogger{
		log: lg,
	}, nil
}

func (l *AppLogger) Info(msg string) {
	l.log.Printf("[Info] %s \n", msg)
}

func (l *AppLogger) Error(msg string) {
	l.log.Printf("[Error] %s \n", msg)
}

func (l *AppLogger) Warn(msg string) {
	l.log.Printf("[Warn] %s \n", msg)
}
