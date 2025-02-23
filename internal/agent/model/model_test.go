package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetric_ValueAsString(t *testing.T) {
	tests := []struct {
		name       string
		value      interface{}
		metricType MetricType
		expValue   interface{}
		error      string
	}{
		{"Positive check Gauge value", 10.0, Gauge, 10.0, "Gauge value is not equal"},
		{"Positive check Counter value", int64(10), Counter, int64(10), "Counter value is not equal"},
		{"Positive check Unknown value", int64(10), "Unknown", int64(10), "Unknown value is not equal"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metric := Metric{Value: test.value, Type: test.metricType}
			assert.Equal(t, test.expValue, metric.Value, test.error)
			assert.Equal(t, test.metricType, metric.Type, "Metric type is not equal")
		})
	}
}
