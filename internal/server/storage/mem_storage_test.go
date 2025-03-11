package storage

import (
	"github.com/BazhanovMaxim/metrics/internal/server/model"
	"github.com/BazhanovMaxim/metrics/internal/server/service"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestMetricStorage_UpdateCounter(t *testing.T) {
	tests := []struct {
		name        string
		error       string
		key         string
		counter     int64
		expectedSum int64
	}{
		{"Positive check update counter", "The value has not changed", "first", 1, 1},
		{"Positive check update counter", "The value has not changed", "first", 99, 100},
		{"Positive check update counter", "The value has not changed", "first", -50, 50},
	}
	storage := NewMemStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _ = storage.Update(model.Metrics{ID: test.key, MType: "counter", Delta: &test.counter})
			mapValues := readCounterData(storage)
			assert.Equal(t, test.expectedSum, mapValues[test.key], test.error)
		})
	}
}

func TestMetricStorage_UpdateGauge(t *testing.T) {
	tests := []struct {
		name     string
		error    string
		key      string
		value    float64
		expected float64
	}{
		{"Positive check update gauge", "The value has not changed", "first", 1, 1},
		{"Positive check update gauge", "The value has not changed", "first", -50, -50},
	}
	storage := NewMemStorage()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, _ = storage.Update(model.Metrics{ID: test.key, MType: "gauge", Value: &test.value})
			mapValues := readGaugeData(storage)
			assert.Equal(t, test.expected, mapValues[test.key], test.error)
		})
	}
}

func readCounterData(metricStorage service.IMetricStorage) map[string]int64 {
	return readField("Counter", metricStorage).Interface().(map[string]int64)
}

func readGaugeData(metricStorage service.IMetricStorage) map[string]float64 {
	return readField("Gauge", metricStorage).Interface().(map[string]float64)
}

func readField(fieldName string, metricStorage service.IMetricStorage) reflect.Value {
	// Получаем значение поля Counter
	field := reflect.ValueOf(metricStorage).Elem().FieldByName(fieldName)
	// Получаем конкретное значение, на которое указывает интерфейс
	return reflect.ValueOf(field.Interface())
}
