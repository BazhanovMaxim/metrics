package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMetric_ValueAsString(t *testing.T) {
	tests := []struct {
		name           string
		value          interface{}
		metricType     MetricType
		expStringValue interface{}
		error          string
	}{
		{"Positive check Gauge value", 10.0, Gauge, "10.000000", "Gauge value is not equal"},
		{"Positive check Counter value", int64(10), Counter, "10", "Counter value is not equal"},
		{"Positive check Unknown value", int64(10), "Unknown", "10", "Unknown value is not equal"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metric := Metric{Value: test.value, Type: test.metricType}
			assert.Equal(t, test.expStringValue, metric.ValueAsString(), test.error)
		})
	}
}

func TestNegativeMetric_ValueAsString(t *testing.T) {
	tests := []struct {
		name       string
		value      interface{}
		metricType MetricType
	}{
		// panic: interface conversion: interface {} is int, not float64 [recovered]
		{"Negative check Gauge float cast value", 10, Gauge},
		// panic: interface conversion: interface {} is float64, not int64 [recovered]
		{"Negative check Counter int cast value", 10.0, Counter},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			metric := Metric{Value: test.value, Type: test.metricType}
			assert.Panics(t, func() { metric.ValueAsString() })
		})
	}
}
