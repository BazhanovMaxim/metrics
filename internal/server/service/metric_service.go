package service

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/logger"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"go.uber.org/zap"
)

// IMetricStorage интерфейс для работы с репозиториями
type IMetricStorage interface {
	Update(metric model.Metrics) (*model.Metrics, error)
	GetAllMetrics() []model.Metrics
	GetMetric(mType, title string) *model.Metrics
	Ping() error
	Close() error
}

type MetricService struct {
	config  configs.Config
	storage IMetricStorage
}

func NewMetricService(config configs.Config, storage IMetricStorage) *MetricService {
	return &MetricService{config: config, storage: storage}
}

func (ms *MetricService) UpdateMetric(metric model.Metrics) (*model.Metrics, error) {
	return ms.storage.Update(metric)
}

func (ms *MetricService) GetMetrics() []model.Metrics {
	return ms.storage.GetAllMetrics()
}

func (ms *MetricService) GetMetricValue(metricType, metricTitle string) *model.Metrics {
	return ms.storage.GetMetric(metricType, metricTitle)
}

func (ms *MetricService) SaveMetricToStorage(metric *model.Metrics) {
	_, err := ms.storage.Update(*metric)
	if err != nil {
		logger.Log.Error("Error", zap.Error(err))
	}
}

func (ms *MetricService) CheckDatabaseConnection() error {
	return ms.storage.Ping()
}
