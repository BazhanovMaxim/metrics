package service

import (
	"fmt"
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"strconv"
)

type MetricService struct {
}

func NewMetricService() *MetricService {
	return &MetricService{}
}

func (ms *MetricService) FindService(metricType string) (func(string, string, storage.IMetricStorage) error, bool) {
	metrics := map[string]func(string, string, storage.IMetricStorage) error{
		"gauge":   ms.updateGauge,
		"counter": ms.updateCounter,
	}
	metric, ok := metrics[metricType]
	return metric, ok
}

func (ms *MetricService) GetMetricValue(metricType, metricTitle string, metricStorage storage.IMetricStorage) (string, bool) {
	switch metricType {
	case "gauge":
		val, ok := metricStorage.GetGauge(metricTitle)
		return strconv.FormatFloat(val, 'f', -1, 64), ok
	case "counter":
		val, ok := metricStorage.GetCounter(metricTitle)
		return fmt.Sprintf("%d", val), ok
	default:
		return "", false
	}
}

func (ms *MetricService) GetCounters(metricStorage storage.IMetricStorage) map[string]int64 {
	return metricStorage.GetCounters()
}

func (ms *MetricService) GetGauges(metricStorage storage.IMetricStorage) map[string]float64 {
	return metricStorage.GetGauges()
}

func (ms *MetricService) updateGauge(metricTitle, metricValue string, storage storage.IMetricStorage) error {
	value, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return err
	}
	storage.UpdateGauge(metricTitle, value)
	return nil
}

func (ms *MetricService) updateCounter(metricTitle, metricValue string, storage storage.IMetricStorage) error {
	value, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		return err
	}
	storage.UpdateCounter(metricTitle, value)
	return nil
}
