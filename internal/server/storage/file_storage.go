package storage

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/file"
	"github.com/BazhanovMaxim/metrics/internal/server/json"
	"github.com/BazhanovMaxim/metrics/internal/server/logger"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"go.uber.org/zap"
	"os"
	"time"
)

type FileStorage struct {
	filePath   string
	config     configs.Config
	memStorage service.IMetricStorage
}

func NewFileStorage(config configs.Config) service.IMetricStorage {
	fs := &FileStorage{
		filePath:   config.FileStoragePath + config.FileStorageName,
		config:     config,
		memStorage: NewMemStorage(),
	}
	fs.init()
	return fs
}

func (s *FileStorage) init() {
	if err := s.loadFile(); err != nil {
		return
	}
	logger.Log.Info("The file for saving metrics has been created successfully or already exists")

	if s.config.Restore {
		logger.Log.Info("Load file storage metrics to memory storage")
		s.loadFileStorageMetricsToMem()
	}

	if s.config.StoreInterval != 0 {
		go s.startPeriodicSave()
	}
}

func (s *FileStorage) Update(metric model.Metrics) (*model.Metrics, error) {
	updatedMetric, err := s.memStorage.Update(metric)
	if err == nil {
		if s.config.StoreInterval == 0 {
			s.updateFileStorageMetric(*updatedMetric)
		}
	}
	return updatedMetric, err
}

func (s *FileStorage) UpdateBatches(metrics []model.Metrics) error {
	for _, metric := range metrics {
		if _, err := s.Update(metric); err != nil {
			return err
		}
	}
	return nil
}

func (s *FileStorage) updateFileStorageMetric(metric model.Metrics) {
	metrics, err := s.readFile()
	if err != nil {
		return
	}

	switch metric.MType {
	case string(model.Gauge):
		metrics.Gauge[metric.ID] = *metric.Value
	default:
		metrics.Counter[metric.ID] = *metric.Delta
	}

	newData, err := json.MarshalIndent(metrics, "", " ")
	if err != nil {
		logger.Log.Error("Failed to indent data", zap.Error(err))
		return
	}

	fl, err := file.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC)
	if err != nil {
		logger.Log.Error("Failed to save metrics to file storage", zap.Error(err))
		return
	}
	defer fl.Close()

	writeError := file.WriteFile(fl, newData)
	if writeError != nil {
		logger.Log.Error("Failed write new data to storage file", zap.Error(err))
	}
}

func (s *FileStorage) GetAllMetrics() []model.Metrics {
	return s.memStorage.GetAllMetrics()
}

func (s *FileStorage) GetMetric(mType, title string) *model.Metrics {
	return s.memStorage.GetMetric(mType, title)
}

func (s *FileStorage) Ping() error {
	return nil
}

func (s *FileStorage) Close() error {
	for _, metric := range s.GetAllMetrics() {
		s.updateFileStorageMetric(metric)
	}
	return nil
}

func (s *FileStorage) readFile() (*model.StorageJSONMetrics, error) {
	data, err := file.ReadFile(s.filePath)
	if err != nil {
		logger.Log.Error("Failed to open metrics file storage", zap.Error(err))
		return nil, err
	}

	var metrics model.StorageJSONMetrics
	// Проверка, является ли файл пустым
	if len(data) != 0 {
		if unmarshalError := json.UnmarshalJSON(data, &metrics); unmarshalError != nil {
			logger.Log.Error("Failed to decode metrics file", zap.Error(unmarshalError))
			return nil, unmarshalError
		}
	}

	if metrics.Gauge == nil {
		metrics.Gauge = make(map[string]float64)
	}

	if metrics.Counter == nil {
		metrics.Counter = make(map[string]int64)
	}
	return &metrics, nil
}

func (s *FileStorage) loadFile() error {
	if err := file.MkdirAll(s.filePath); err != nil {
		logger.Log.Error("Failed to create metrics directories storage", zap.Error(err))
		return err
	}
	if _, openFileErr := file.OpenFile(s.filePath, os.O_RDWR|os.O_CREATE); openFileErr != nil {
		logger.Log.Error("Failed to create metrics file storage", zap.Error(openFileErr))
		return openFileErr
	}
	return nil
}

func (s *FileStorage) loadFileStorageMetricsToMem() {
	metrics, err := s.readFile()
	if err != nil {
		logger.Log.Error("Failed to read file storage", zap.Error(err))
	}
	for key, value := range metrics.Counter {
		_, _ = s.memStorage.Update(model.Metrics{ID: key, MType: string(model.Counter), Delta: &value})
	}
	for key, value := range metrics.Gauge {
		_, _ = s.memStorage.Update(model.Metrics{ID: key, MType: string(model.Gauge), Value: &value})
	}
}

func (s *FileStorage) startPeriodicSave() {
	ticker := time.NewTicker(time.Duration(s.config.StoreInterval) * time.Second)
	defer ticker.Stop()

	// Канал для остановки горутины
	stopChan := make(chan struct{})

	// Горутина для периодического сохранения данных
	for {
		select {
		case <-ticker.C:
			logger.Log.Info("Save memory storage metrics to file storage")
			for _, metric := range s.GetAllMetrics() {
				s.updateFileStorageMetric(metric)
			}
		case <-stopChan:
			return
		}
	}
}
