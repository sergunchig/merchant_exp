package storage

import (
	"fmt"
	"io"
	"os"
)

type loggerI interface {
	Error(msg string)
}
type StorageService struct {
	Log loggerI
}

func New(log loggerI) StorageService {
	return StorageService{
		Log: log,
	}
}
func (s StorageService) SaveFile(in io.Reader, fileName string) error {

	fmt.Println(fileName)

	newFile, err := os.Create(fileName)
	defer func() {
		if err = newFile.Close(); err != nil {
			s.Log.Error("error close file")
		}
	}()
	if err != nil {
		return fmt.Errorf("can't create file %w", err)
	}
	_, err = io.Copy(newFile, in)
	if err != nil {
		return fmt.Errorf("can't copy file %w", err)
	}
	return nil
}
