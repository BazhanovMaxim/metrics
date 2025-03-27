package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/BazhanovMaxim/metrics/internal/agent/logger"
	"github.com/BazhanovMaxim/metrics/internal/agent/model"
	"github.com/BazhanovMaxim/metrics/internal/agent/router"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
	"go.uber.org/zap"
	"net"
	"strings"
	"sync"
	"time"
)

type MetricService struct {
	config   configs.Config
	handlers router.Router
	storage  storage.IMetricStorage
}

// NewMetricService создает и возвращает новый экземпляр MetricService
func NewMetricService(config configs.Config, storage storage.IMetricStorage, handlers router.Router) *MetricService {
	return &MetricService{config: config, handlers: handlers, storage: storage}
}

func (ms *MetricService) Start(ctx context.Context) {
	if ms.config.RateLimit == 0 {
		ms.sendMetrics(ctx)
		return
	}
	ms.sendMetricsWithRateLimit(ctx)
}

// sendMetricsWithRateLimit обновляет метрики на сервере с ограничением Rate Limit
func (ms *MetricService) sendMetricsWithRateLimit(ctx context.Context) {
	var wg sync.WaitGroup

	channel := make(chan []model.Metrics, ms.config.RateLimit)

	wg.Add(1)
	go ms.collectMetrics(ctx, channel, &wg)

	wg.Add(1)
	go ms.collectSysAndMemMetrics(ctx, channel, &wg)

	for i := 0; i < ms.config.RateLimit; i++ {
		wg.Add(1)
		go ms.send(ctx, channel, &wg)
	}

	wg.Wait()
	close(channel)
}

// collectMetrics собирает основные метрики
func (ms *MetricService) collectMetrics(ctx context.Context, ch chan<- []model.Metrics, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		// Контекст отменён, завершаем работу
		case <-ctx.Done():
			return
		// Отправляем метрику в канал
		default:
			ms.storage.Update()
			ch <- ms.convertMetrics()
		}
	}
}

// collectSysAndMemMetrics собирает дополнительные метрики(TotalMemory, FreeMemory, CPUutilization1)
func (ms *MetricService) collectSysAndMemMetrics(ctx context.Context, ch chan<- []model.Metrics, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		// Контекст отменён, завершаем работу
		case <-ctx.Done():
			return
		// Отправляем метрику в канал
		default:
			ms.storage.UpdateSys()
			ch <- ms.convertMetrics()
		}
	}
}

// send отправляет метрики на сервер
func (ms *MetricService) send(ctx context.Context, ch <-chan []model.Metrics, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			// Контекст отменён, завершаем работу
			return
		case val, ok := <-ch:
			if !ok {
				// Канал закрыт, завершаем работу
				return
			}
			body, marshErr := json.Marshal(val)
			if marshErr != nil {
				logger.Log.Error("Failed to marshal POJO", zap.Error(marshErr))
				return
			}
			err := ms.handlers.SendMetrics(ms.config, body)
			if err != nil {
				logger.Log.Error("Failed to send metrics", zap.Error(err))
				continue // Продолжаем работу, если возникла ошибка отправки
			}
		}
	}
}

// sendMetrics обновляет метрики на сервере с интервалами PollInterval и ReportInterval
func (ms *MetricService) sendMetrics(ctx context.Context) {
	pollTicker := time.NewTicker(time.Duration(ms.config.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(ms.config.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		// обновляет метрики в репозитории
		case <-pollTicker.C:
			ms.storage.Update()
		// отправляет метрики на сервер
		case <-reportTicker.C:
			ms.sendMetricsToServer()
		// Выходим из метода по завершению работы контекста
		case <-ctx.Done():
			return
		}
	}
}

func (ms *MetricService) sendMetricsToServer() {
	body, marshErr := json.Marshal(ms.convertMetrics())
	if marshErr != nil {
		logger.Log.Error("Failed to marshal POJO", zap.Error(marshErr))
		return
	}
	if err := ms.handlers.SendMetrics(ms.config, body); err != nil {
		var opErr *net.OpError
		if errors.As(err, &opErr) && strings.Contains(err.Error(), "connect: connection refused") {
			logger.Log.Info("The server is unavailable. Retry sending metrics to the server")
			ticker := time.NewTicker(2 * time.Second)
			for i := 0; i < 3; i++ {
				<-ticker.C
				if err = ms.handlers.SendMetrics(ms.config, body); err == nil ||
					(!errors.As(err, &opErr) && !strings.Contains(err.Error(), "connect: connection refused")) {
					ticker.Stop()
					logger.Log.Error("Failed to send message to server", zap.Error(err))
					return
				}
			}
			ticker.Stop()
		}
	}
}

func (ms *MetricService) convertMetrics() []model.Metrics {
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
	return metrics
}
