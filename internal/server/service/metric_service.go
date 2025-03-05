package service

import (
	"github.com/BazhanovMaxim/metrics/internal/server/configs"
	"github.com/BazhanovMaxim/metrics/internal/server/json"
	"github.com/BazhanovMaxim/metrics/internal/server/logger"
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"go.uber.org/zap"
	"strconv"
)

// MetricService
// todo: подумать, можно ли упростить передачу параметров
// из вариантов, разделить MetricService на fileMetricService, memMetricService.
// При таком подходе общение между сервисами усложняется, но сервис не будет заботится, что за тип сторадж в нем лежит
type MetricService struct {
	config      configs.Config
	memStorage  storage.IMemStorage
	fileStorage storage.IFileStorage
	dbStorage   storage.IDbStorage
}

func NewMetricService(config configs.Config, memStorage storage.IMemStorage, fileStorage storage.IFileStorage, dbStorage storage.IDbStorage) *MetricService {
	return &MetricService{
		config:      config,
		memStorage:  memStorage,
		fileStorage: fileStorage,
		dbStorage:   dbStorage,
	}
}

func (ms *MetricService) FindService(metricType string) func(string, interface{}) *model.Metrics {
	switch metricType {
	case string(model.Gauge):
		return ms.updateGauge
	case string(model.Counter):
		return ms.updateCounter
	default:
		return nil
	}
}

func (ms *MetricService) GetMetricValue(metricType, metricTitle string) *model.Metrics {
	switch metricType {
	case string(model.Gauge):
		if value, find := ms.GetGauges()[metricTitle]; find {
			return &model.Metrics{ID: metricTitle, MType: string(model.Gauge), Value: &value}
		}
		return nil
	case string(model.Counter):
		if value, find := ms.GetCounters()[metricTitle]; find {
			return &model.Metrics{ID: metricTitle, MType: string(model.Counter), Delta: &value}
		}
		return nil
	default:
		return nil
	}
}

func (ms *MetricService) GetCounters() map[string]int64 {
	return ms.memStorage.GetCounters()
}

func (ms *MetricService) GetGauges() map[string]float64 {
	return ms.memStorage.GetGauges()
}

func (ms *MetricService) updateGauge(metricTitle string, metricValue interface{}) *model.Metrics {
	switch val := metricValue.(type) {
	case float64:
		value := ms.memStorage.UpdateGauge(metricTitle, val)
		return &model.Metrics{ID: metricTitle, MType: string(model.Gauge), Value: &value}
	case string:
		floatValue, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil
		}
		value := ms.memStorage.UpdateGauge(metricTitle, floatValue)
		return &model.Metrics{ID: metricTitle, MType: string(model.Gauge), Value: &value}
	default:
		return nil
	}
}

func (ms *MetricService) updateCounter(metricTitle string, metricValue interface{}) *model.Metrics {
	switch val := metricValue.(type) {
	case int64:
		value := ms.memStorage.UpdateCounter(metricTitle, val)
		return &model.Metrics{ID: metricTitle, MType: string(model.Counter), Delta: &value}
	case string:
		intValue, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil
		}
		value := ms.memStorage.UpdateCounter(metricTitle, intValue)
		return &model.Metrics{ID: metricTitle, MType: string(model.Counter), Delta: &value}
	default:
		return nil
	}
}

func (ms *MetricService) LoadStorageMetrics() {
	logger.Log.Info("Load storage metrics")
	data := ms.fileStorage.ReadFile()

	var metrics model.StorageJSONMetrics
	// Проверка, является ли файл пустым
	if len(data) != 0 {
		if marshalError := json.UnmarshalJSON(data, &metrics); marshalError != nil {
			logger.Log.Error("Failed to decode metrics file", zap.Error(marshalError))
			return
		}
	}

	if metrics.Gauge != nil {
		for key, value := range metrics.Gauge {
			ms.memStorage.UpdateGauge(key, value)
		}
	}

	if metrics.Counter != nil {
		for key, value := range metrics.Counter {
			ms.memStorage.UpdateCounter(key, value)
		}
	}
}

func (ms *MetricService) SaveMetricsToStorage() {
	logger.Log.Info("Save metrics to storage")
	for key, value := range ms.memStorage.GetGauges() {
		ms.SaveMetricToStorage(&model.Metrics{ID: key, MType: string(model.Gauge), Value: &value})
	}
	for key, value := range ms.memStorage.GetCounters() {
		ms.SaveMetricToStorage(&model.Metrics{ID: key, MType: string(model.Counter), Delta: &value})
	}
}

func (ms *MetricService) SaveMetricToStorage(metric *model.Metrics) {
	data := ms.fileStorage.ReadFile()

	var metrics model.StorageJSONMetrics
	// Проверка, является ли файл пустым
	if len(data) != 0 {
		if err := json.UnmarshalJSON(data, &metrics); err != nil {
			logger.Log.Error("Failed to decode metrics file", zap.Error(err))
			return
		}
	}

	if metrics.Gauge == nil {
		metrics.Gauge = make(map[string]float64)
	}

	if metrics.Counter == nil {
		metrics.Counter = make(map[string]int64)
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

	ms.fileStorage.WriteFile(newData)
}

func (ms *MetricService) PingConnection() error {
	return ms.dbStorage.Ping()
}
