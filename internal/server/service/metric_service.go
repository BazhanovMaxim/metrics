package service

import (
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"strconv"
)

type MetricService struct {
}

func NewMetricService() *MetricService {
	return &MetricService{}
}

func (ms *MetricService) FindService(metricType string) func(string, interface{}, storage.MetricStorage) *model.Metrics {
	switch metricType {
	case string(model.Gauge):
		return ms.updateGauge
	case string(model.Counter):
		return ms.updateCounter
	default:
		return nil
	}
}

func (ms *MetricService) GetMetricValue(metricType, metricTitle string, metricStorage storage.MetricStorage) (*model.Metrics, bool) {
	switch metricType {
	case string(model.Gauge):
		if value, find := ms.GetGauges(metricStorage)[metricTitle]; find {
			return &model.Metrics{ID: metricTitle, MType: string(model.Gauge), Value: &value}, true
		}
		return nil, false
	case string(model.Counter):
		if value, find := ms.GetCounters(metricStorage)[metricTitle]; find {
			return &model.Metrics{ID: metricTitle, MType: string(model.Counter), Delta: &value}, true
		}
		return nil, false
	default:
		return nil, false
	}
}

func (ms *MetricService) GetCounters(metricStorage storage.MetricStorage) map[string]int64 {
	return metricStorage.Counter.GetAll()
}

func (ms *MetricService) GetGauges(metricStorage storage.MetricStorage) map[string]float64 {
	return metricStorage.Gauge.GetAll()
}

func (ms *MetricService) updateGauge(metricTitle string, metricValue interface{}, storage storage.MetricStorage) *model.Metrics {
	switch val := metricValue.(type) {
	case float64:
		value := storage.Gauge.Update(metricTitle, val)
		return &model.Metrics{ID: metricTitle, MType: string(model.Gauge), Value: &value}
	case string:
		floatValue, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return nil
		}
		value := storage.Gauge.Update(metricTitle, floatValue)
		return &model.Metrics{ID: metricTitle, MType: string(model.Gauge), Value: &value}
	default:
		return nil
	}
}

func (ms *MetricService) updateCounter(metricTitle string, metricValue interface{}, storage storage.MetricStorage) *model.Metrics {
	switch val := metricValue.(type) {
	case int64:
		value := storage.Counter.Update(metricTitle, val)
		return &model.Metrics{ID: metricTitle, MType: string(model.Counter), Delta: &value}
	case string:
		intValue, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return nil
		}
		value := storage.Counter.Update(metricTitle, intValue)
		return &model.Metrics{ID: metricTitle, MType: string(model.Counter), Delta: &value}
	default:
		return nil
	}
}
