package handlers

import (
	"fmt"
	"github.com/BazhanovMaxim/metrics/internal/agent/model"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
	"github.com/go-resty/resty/v2"
)

func (h *Handler) sendMetrics(storage *storage.MetricStorage, client *resty.Client) error {
	for name, metric := range storage.GetMetrics() {
		var valueString string
		switch metric.Type {
		case model.Gauge:
			valueString = fmt.Sprintf("%f", metric.Value.(float64))
		default:
			valueString = fmt.Sprintf("%d", metric.Value.(int64))
		}

		url := fmt.Sprintf("http://%s/update/%s/%s/%s", h.config.RunAddress, metric.Type, name, valueString)
		_, err := client.R().Post(url)
		if err != nil {
			return err
		}
	}
	return nil
}
