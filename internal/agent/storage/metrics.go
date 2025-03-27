package storage

import "github.com/BazhanovMaxim/metrics/internal/agent/model"

type IMetricStorage interface {
	GetMetrics() map[string]model.Metric
	Update()
	UpdateSys()
}
