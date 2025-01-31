package handlers

import (
	"github.com/BazhanovMaxim/metrics/internal/agent/flags"
	"github.com/BazhanovMaxim/metrics/internal/agent/storage"
	"github.com/go-resty/resty/v2"
)

func (h *Handler) sendMetrics(storage *storage.MetricStorage) error {
	for name, metric := range storage.GetMetrics() {
		client := resty.New()

		_, err := client.R().Post("http://" + flags.RunAddress + "/update/" +
			string(metric.Type) +
			"/" + name + "/" +
			metric.ValueAsString())
		if err != nil {
			return err
		}
	}
	return nil
}
