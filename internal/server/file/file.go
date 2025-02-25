package file

import (
	"os"
	"path/filepath"
)

func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func OpenFile(filePath string, flag int) (*os.File, error) {
	return os.OpenFile(filePath, flag, 0666)
}

func MkdirAll(filePath string) error {
	return os.MkdirAll(filepath.Dir(filePath), 0755)
}

func WriteFile(file *os.File, data []byte) error {
	_, err := file.Write(data)
	return err
}
