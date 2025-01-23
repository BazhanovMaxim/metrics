package service

import (
	"github.com/BazhanovMaxim/metrics/internal/server/repository"
	"strconv"
)

type MetricService struct {
}

func GetMetricService() *MetricService {
	return &MetricService{}
}

func (ms *MetricService) FindService(metricType string) (func([]string, repository.MetricStorage) error, bool) {
	metrics := map[string]func([]string, repository.MetricStorage) error{
		"gauge":   ms.updateGauge,
		"counter": ms.updateCounter,
	}
	metric, ok := metrics[metricType]
	return metric, ok
}

func (ms *MetricService) updateGauge(path []string, storage repository.MetricStorage) error {
	metricTitle := path[3]
	metricValue := path[4]
	value, err := strconv.ParseFloat(metricValue, 64)
	if err != nil {
		return err
	}
	storage.UpdateGauge(metricTitle, value)
	return nil
}

func (ms *MetricService) updateCounter(path []string, storage repository.MetricStorage) error {
	metricTitle := path[3]
	metricValue := path[4]
	value, err := strconv.ParseInt(metricValue, 10, 64)
	if err != nil {
		return err
	}
	storage.UpdateCounter(metricTitle, value)
	return nil
}
