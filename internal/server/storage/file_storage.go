package storage

import (
	"github.com/BazhanovMaxim/metrics/internal/server/file"
	"github.com/BazhanovMaxim/metrics/internal/server/logger"
	"go.uber.org/zap"
	"os"
)

type IFileStorage interface {
	IMetricStorage
	WriteFile(data []byte)
	ReadFile() []byte
}

type FileStorage struct {
	filePath string
}

func NewFileStorage(filePath string) IFileStorage {
	fs := &FileStorage{filePath: filePath}
	fs.init()
	return fs
}

func (s *FileStorage) WriteFile(data []byte) {
	fl, err := file.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		logger.Log.Error("Failed to save metrics to file storage", zap.Error(err))
		return
	}
	defer fl.Close()

	writeError := file.WriteFile(fl, data)
	if writeError != nil {
		logger.Log.Error("Failed write new data to storage file", zap.Error(err))
	}
}

func (s *FileStorage) ReadFile() []byte {
	data, err := file.ReadFile(s.filePath)
	if err != nil {
		logger.Log.Error("Failed to open metrics file storage", zap.Error(err))
		return nil
	}
	return data
}

func (s *FileStorage) init() {
	if err := file.MkdirAll(s.filePath); err != nil {
		logger.Log.Error("Failed to create metrics directories storage", zap.Error(err))
		return
	}
	if _, openFileErr := file.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE); openFileErr != nil {
		logger.Log.Error("Failed to create metrics file storage", zap.Error(openFileErr))
		return
	}
	logger.Log.Info("The file for saving metrics has been created successfully or already exists")
}
