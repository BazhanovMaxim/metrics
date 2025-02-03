package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/agent/configs"
	"github.com/BazhanovMaxim/metrics/internal/agent/service"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
	"github.com/go-resty/resty/v2"
	"time"
)

type Handler struct {
	storage storage.MetricStorage
	config  configs.Config
}

func NewHandler(config configs.Config, storage storage.MetricStorage) *Handler {
	return &Handler{config: config, storage: storage}
}

func (h *Handler) Start() error {
	client := resty.New()
	pollTicker := time.NewTicker(time.Duration(h.config.PollInterval) * time.Second)
	reportTicker := time.NewTicker(time.Duration(h.config.ReportInterval) * time.Second)
	defer pollTicker.Stop()
	defer reportTicker.Stop()

	for {
		select {
		case <-pollTicker.C:
			service.GetMetricService().UpdateMetric(&h.storage)
		case <-reportTicker.C:
			if err := h.sendMetrics(&h.storage, client); err != nil {
				return err
			}
		}
	}
}
