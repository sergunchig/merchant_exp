package storage

import (
	"fmt"
	"io"
	"os"
)

func SaveFile(in io.Reader, fileName string) error {

	fmt.Println(fileName)

	newFile, err := os.Create(fileName)
	defer newFile.Close()
	if err != nil {
		return fmt.Errorf("can't create file %w", err)
	}
	_, err = io.Copy(newFile, in)
	if err != nil {
		return fmt.Errorf("Can't copy file %w", err)
	}
	return nil
}
