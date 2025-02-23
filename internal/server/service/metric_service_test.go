package service

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetricService_FindService(t *testing.T) {
	metricService := NewMetricService()
	tests := []struct {
		name       string
		metricType string
		expectedOk bool
		errorText  string
	}{
		{"Positive get Gauge service", "gauge", true, "Gauge service is not founded"},
		{"Positive get Counter service", "counter", true, "Counter service is not founded"},
		{"Negative unknown service", "unknown", false, "Unknown service is founded"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := metricService.FindService(test.metricType)
			if test.expectedOk {
				assert.NotNil(t, service)
				return
			}
			assert.Nil(t, service)
		})
	}
}
