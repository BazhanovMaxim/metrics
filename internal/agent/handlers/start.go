package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/BazhanovMaxim/metrics/internal/agent/service"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
	"time"
)

type Handler struct {
	storage storage.MetricStorage
}

func NewHandler(storage *storage.MetricStorage) *Handler {
	return &Handler{storage: *storage}
}

func (h *Handler) Start() error {
	pollTicker := time.NewTicker(time.Duration(configs.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(configs.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			service.GetMetricService().UpdateMetric(&h.storage)
		case <-reportTicker.C:
			if err := h.sendMetrics(&h.storage); err != nil {
				return err
			}
		}
	}
}
