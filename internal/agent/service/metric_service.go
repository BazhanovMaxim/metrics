package service

import (
	"context"
	"encoding/json"
	"github.com/BazhanovMaxim/metrics/internal/agent/compress"
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/BazhanovMaxim/metrics/internal/agent/logger"
	"github.com/BazhanovMaxim/metrics/internal/agent/model"
	"github.com/BazhanovMaxim/metrics/internal/agent/router"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
	"go.uber.org/zap"
	"time"
)

type MetricService struct {
	config   configs.Config
	handlers router.Router
	storage  storage.MetricStorage
}

func NewMetricService(config configs.Config, storage storage.MetricStorage, handlers router.Router) *MetricService {
	return &MetricService{config: config, handlers: handlers, storage: storage}
}

func (ms *MetricService) Start(ctx context.Context) {
	pollTicker := time.NewTicker(time.Duration(ms.config.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(ms.config.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		// Выходим из метода по завершению работы контекста
		case <-ctx.Done():
			return
		case <-pollTicker.C:
			ms.UpdateMetric()
		case <-reportTicker.C:
			ms.SendMetricsToServer()
		}
	}
}

func (ms *MetricService) UpdateMetric() {
	ms.storage.Update()
}

func (ms *MetricService) SendMetricsToServer() {
	for id, metric := range ms.storage.GetMetrics() {
		metricsPojo := model.Metrics{MType: string(metric.Type), ID: id}
		switch v := metric.Value.(type) {
		case float64:
			metricsPojo.Value = &v
		case int64:
			metricsPojo.Delta = &v
		}

		marshPojo, marshErr := json.Marshal(metricsPojo)
		if marshErr != nil {
			logger.Log.Error("Failed to marshal POJO", zap.Error(marshErr))
			return
		}
		buf, compressErr := compress.GzipCompress(marshPojo)
		if compressErr != nil {
			logger.Log.Error("Failed to compress data", zap.Error(compressErr))
			return
		}
		ms.handlers.SendMetrics(ms.config, buf.Bytes())
	}
}
