package service

import (
	"github.com/BazhanovMaxim/metrics/internal/server/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricService_FindService(t *testing.T) {
	metriceService := GetMetricService()
	tests := []struct {
		name         string
		metricType   string
		expectedFunc func([]string, storage.IMetricStorage) error
		expectedOk   bool
		errorText    string
	}{
		{"Positive get Gauge service", "gauge", metriceService.updateGauge, true, "Gauge service is not founded"},
		{"Positive get Counter service", "counter", metriceService.updateCounter, true, "Counter service is not founded"},
		{"Negative unknown service", "unknown", nil, false, "Unknown service is founded"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, ok := metriceService.FindService(test.metricType)
			assert.Equal(t, test.expectedOk, ok, test.errorText)
			// todo: В Go нельзя сравнивать функции, подумать над тестом
			//	assert.Equal(t, test.expectedFunc, service, "Expected function mismatch")
		})
	}
}
