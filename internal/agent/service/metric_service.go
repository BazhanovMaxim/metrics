package service

import (
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
)

type MetricService struct {
}

func GetMetricService() *MetricService {
	return &MetricService{}
}

func (ms *MetricService) UpdateMetric(storage *storage.MetricStorage) {
	storage.UpdateMetrics()
}
