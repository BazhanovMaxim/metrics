package service

import (
	"fmt"
	"github.com/BazhanovMaxim/metrics/internal/agent/model"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"strconv"
)

type MetricService struct {
}

func NewMetricService() *MetricService {
	return &MetricService{}
}

func (ms *MetricService) FindService(metricType string) (func(string, string, storage.MetricStorage) error, bool) {
	metrics := map[string]func(string, string, storage.MetricStorage) error{
		"gauge":   ms.updateGauge,
		"counter": ms.updateCounter,
	}
	metric, ok := metrics[metricType]
	return metric, ok
}

func (ms *MetricService) GetMetricValue(metricType, metricTitle string, metricStorage storage.MetricStorage) (string, bool) {
	switch metricType {
	case string(model.Gauge):
		val, ok := metricStorage.Gauge.Get(metricTitle)
		return strconv.FormatFloat(val, 'f', -1, 64), ok
	case string(model.Counter):
		val, ok := metricStorage.Counter.Get(metricTitle)
		return fmt.Sprintf("%d", val), ok
	default:
		return "", false
	}
}

func (ms *MetricService) GetCounters(metricStorage storage.MetricStorage) map[string]int64 {
	return metricStorage.Counter.GetAll()
}

func (ms *MetricService) GetGauges(metricStorage storage.MetricStorage) map[string]float64 {
	return metricStorage.Gauge.GetAll()
}

func (ms *MetricService) updateGauge(metricTitle, metricValue string, storage storage.MetricStorage) error {
	value, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return err
	}
	storage.Gauge.Update(metricTitle, value)
	return nil
}

func (ms *MetricService) updateCounter(metricTitle, metricValue string, storage storage.MetricStorage) error {
	value, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		return err
	}
	storage.Counter.Update(metricTitle, value)
	return nil
}
