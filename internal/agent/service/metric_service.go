package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/BazhanovMaxim/metrics/internal/agent/compress"
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/BazhanovMaxim/metrics/internal/agent/logger"
	"github.com/BazhanovMaxim/metrics/internal/agent/model"
	"github.com/BazhanovMaxim/metrics/internal/agent/router"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
	"go.uber.org/zap"
	"net"
	"strings"
	"time"
)

type MetricService struct {
	config   configs.Config
	handlers router.Router
	storage  storage.IMetricStorage
}

func NewMetricService(config configs.Config, storage storage.IMetricStorage, handlers router.Router) *MetricService {
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
	var metrics []model.Metrics
	for id, metric := range ms.storage.GetMetrics() {
		metricsPojo := model.Metrics{MType: string(metric.Type), ID: id}
		switch v := metric.Value.(type) {
		case float64:
			metricsPojo.Value = &v
		case int64:
			metricsPojo.Delta = &v
		}
		metrics = append(metrics, metricsPojo)
	}
	marshPojo, marshErr := json.Marshal(metrics)
	if marshErr != nil {
		logger.Log.Error("Failed to marshal POJO", zap.Error(marshErr))
		return
	}
	buf, compressErr := compress.GzipCompress(marshPojo)
	if compressErr != nil {
		logger.Log.Error("Failed to compress data", zap.Error(compressErr))
		return
	}
	if err := ms.handlers.SendMetrics(ms.config, buf.Bytes()); err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) && strings.Contains(err.Error(), "connect: connection refused") {
			logger.Log.Info("The server is unavailable. Retry sending metrics to the server")
			ticker := time.NewTicker(2 * time.Second)
			for i := 0; i < 3; i++ {
				<-ticker.C
				if err = ms.handlers.SendMetrics(ms.config, buf.Bytes()); err == nil ||
					(!errors.As(err, &opErr) && !strings.Contains(err.Error(), "connect: connection refused")) {
					ticker.Stop()
					return
				}
			}
			ticker.Stop()
		}
	}
}
