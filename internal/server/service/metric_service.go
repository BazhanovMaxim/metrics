package service

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
)

// IMetricStorage интерфейс для работы с репозиториями
type IMetricStorage interface {
	Update(metric model.Metrics) (*model.Metrics, error) // Обновляет метрику
	UpdateBatches(metrics []model.Metrics) error         // Обновляет метрики
	GetAllMetrics() []model.Metrics                      // Получает и возвращает все метрики из хранилища
	GetMetric(mType, title string) *model.Metrics        // Получает и возвращает метрику из хранилища
	Ping() error                                         // Проверяет работоспособность хранилища
	Close() error                                        // Завершает работу хранилища
}

// MetricService представляет собой обьект для взаимодействия между хендлером и хранилищем
type MetricService struct {
	config  configs.Config
	storage IMetricStorage
}

// NewMetricService создает и возвращает новый экземпляр MetricService
func NewMetricService(config configs.Config, storage IMetricStorage) *MetricService {
	return &MetricService{config: config, storage: storage}
}

// UpdateMetric обновляет метрику в хранилище
func (ms *MetricService) UpdateMetric(metric model.Metrics) (*model.Metrics, error) {
	return ms.storage.Update(metric)
}

// UpdateBatches обновляет коллекцию метрик в хранилище
func (ms *MetricService) UpdateBatches(metrics []model.Metrics) error {
	return ms.storage.UpdateBatches(metrics)
}

// GetMetrics получает и возвращает все метрики из хранилища
func (ms *MetricService) GetMetrics() []model.Metrics {
	return ms.storage.GetAllMetrics()
}

// GetMetricValue получает и возвращает метрику по ее параметрам
func (ms *MetricService) GetMetricValue(metricType, metricTitle string) *model.Metrics {
	return ms.storage.GetMetric(metricType, metricTitle)
}

// CheckDatabaseConnection проверяет доступность к базе данных
// todo: не уверен, что это здесь должно находится
func (ms *MetricService) CheckDatabaseConnection() error {
	return ms.storage.Ping()
}
